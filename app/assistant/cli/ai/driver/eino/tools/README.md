# 在此留言

## 1. 为什么有两个子目录？
你可以看到这个文件夹有两个子目录: `local`和`remote`

- `local`对应本地调用，AI模型调用`function calling`时，不会触发对其他服务器的调用。
- `remote`对应远程调用，AI模型调用`function calling`时，会触发对其他服务器的调用。

尽管`local`能够通过共用包解决，但鉴于不同AI模型的接口参数与响应参数不同。依然需要专门实现一个`local`目录。

> 如果AI服务共用一套规范，比如`openapi 3.0`, 那么不管是`local`还是`remote`都可以通过共用包解决。
> 这当然是以后的事情了。现在我们只需要专注于`local`和`remote`的实现即可。

## 2. `tool.Reflect`

### 2.1 为什么要有`tool.Reflect`？

目前，你应该可以发现其中有一个对外暴露的函数`Reflect`。这个函数是用于将Golang的结构体映射到`map[string]*schema.ParameterInfo`的一个工具函数。

没有这个工具函数，如果我们要约束AI模型的输入输出参数，我们会这样实现一个对Function Calling的Tool的描述:

``` go
type ToolRepeat struct {
	tool.InvokableTool
}

const (
	ToolRepeatName = "repeat"
	ToolRepeatDesc = "重复用户的输入"
)

type RepeatArgs struct {
	Message string `json:"message"`
}

func Repeat() *ToolRepeat {
	return &ToolRepeat{}
}

func (t *ToolRepeat) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	args := &RepeatArgs{}
	err := json.Unmarshal([]byte(argumentsInJSON), args)
	if err != nil {
		return "", err
	}
	return args.Message, nil
}

func (t *ToolRepeat) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: ToolRepeatName,
		Desc: ToolRepeatDesc,
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"message": {
				Type:     schema.String,
				Desc:     "要重复的消息",
				Required: true,
			},
		}),
	}, nil
}
```

很明显，这样的实现是不够优雅的。尽管`RepeatArgs`是一个结构体，但是我们还是需要手动的将其映射到`map[string]*schema.ParameterInfo`。
这个过程是重复的，而且容易出错。同时会产生大量的魔法字符串。

`tool.go`的目的就是为了解决这个问题。我们只需要定义一个结构体，然后调用`Reflect`函数，就可以得到一个`map[string]*schema.ParameterInfo`。

``` go
type ToolRepeat struct {
	tool.InvokableTool
}

const (
	ToolRepeatName = "repeat"
	ToolRepeatDesc = "重复用户的输入"
)

type ToolRepeatArgs struct {
	Message string `json:"message" yaml:"message" desc:"要重复的消息" required:"true"`
}

var ToolRepeatRequestBody = schema.NewParamsOneOfByParams(*tools.Reflect(ToolRepeatArgs{}))

func Repeat() *ToolRepeat {
	return &ToolRepeat{}
}

func (t *ToolRepeat) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	args := &ToolRepeatArgs{}
	err := json.Unmarshal([]byte(argumentsInJSON), args)
	if err != nil {
		return "", err
	}
	return args.Message, nil
}

func (t *ToolRepeat) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name:        ToolRepeatName,
		Desc:        ToolRepeatDesc,
		ParamsOneOf: ToolRepeatRequestBody,
	}, nil
}
```

这样，我们就可以避免手动的将结构体映射到`map[string]*schema.ParameterInfo`。

### 2.2 `tool.Reflect`的实现原理

`tool.Reflect`的原理是通过反射的方式，将结构体的字段映射到`map[string]*schema.ParameterInfo`。

我们约定，结构体的字段的标签中，`json`和`yaml`是用于序列化和反序列化的，`desc`是用于描述字段的，`required`是用于标记字段是否必须的，`enum`是用于标记字段的枚举值。

很显然，这与`openapi 3.0`的规范是一致的。对于`schema.ParameterInfo`的`Type`字段，我们会根据字段的类型，映射到`schema.ParameterInfo`的`Type`字段。

目前支持的类型有:`string`, `number`, `integer`, `boolean`, `array`, `object`

- 对于`array`和`object`类型，我们会递归的映射到`schema.ParameterInfo`的`ElemInfo`和`SubParams`字段。
- 根据`string`的`enum`标签，映射到`schema.ParameterInfo`的`Enum`字段。
  > `enum`标签的值是用逗号分隔的字符串。
  > 
  > 比如`enum:"success,failed"`，我们会将其映射到`schema.ParameterInfo`的`Enum`字段为`[]string{"success", "failed"}`。
  > 
  > 不会处理其中的空格，比如`enum:"success, failed"`，我们会将其映射到`schema.ParameterInfo`的`Enum`字段为`[]string{"success", " failed"}`。
- 根据`required`标签，映射到`schema.ParameterInfo`的`Required`字段。
- 根据`desc`标签，映射到`schema.ParameterInfo`的`Desc`字段。
- 根据`json`标签，映射到`schema.ParameterInfo`的`Name`字段。
  > 当然，如果`json`标签不存在，我们会使用字段的小写名称。

### 2.3 `tool.Reflect`的使用例子

比如我们有一个结构体:
``` go
type ToolRepeatArgs struct {
    Message string `json:"message" yaml:"message" desc:"要重复的消息" required:"true"`
}
```
`tool.Reflect`会将`ToolRepeatArgs`的字段`Message`映射到`map[string]*schema.ParameterInfo`:
``` go
{
    "message": {
        Type:     schema.String,
        Desc:     "要重复的消息",
        Required: true,
    },
}
```

同时支持结构体、数组的嵌套，以及字符串的枚举：
``` go
type Args struct {
    Status  string   `json:"status" yaml:"status" desc:"状态"      required:"true" enum:"success,failed"`
    StrArr  []string `json:"strArr" yaml:"strArr" desc:"字符串数组" required:"true"`
    Struct  struct {
        Name string `json:"name" yaml:"name" desc:"名字" required:"true"`
    } `json:"struct" yaml:"struct" desc:"结构体" required:"true"`
}
```
将会映射到:
``` go
{
    "status": {
        Type:     schema.String,
        Desc:     "状态",
        Required: true,
        Enum:     []string{"success", "failed"},
    },
    "strArr": {
        Type:     schema.Array,
        Desc:     "字符串数组",
        Required: true,
        Items:    &schema.ParameterInfo{
            Type: schema.String,
        },
    },
    "struct": {
        Type:     schema.Object,
        Desc:     "结构体",
        Required: true,
        Properties: map[string]*schema.ParameterInfo{
            "name": {
                Type:     schema.String,
                Desc:     "名字",
                Required: true,
            },
        },
    },
}
```
