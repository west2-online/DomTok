# 1. 目录结构

```
.
├── cli // 调用远程服务的客户端
│   ├── ai // 调用AI服务的客户端
│   │   ├── adapter // 适配器
│   │   └── driver  // 驱动
│   │       └── eino // 基于Eino实现的AI服务
│   │           ├── model // 定义一些行为模型，比如传入参数、返回参数以及构造规则
│   │           └── tools // 用于Function Calling的工具函数
│   │               ├── local  // 数据位于本地的工具函数
│   │               └── remote // 数据位于远程的工具函数
│   └── server // 调用服务器的客户端
│       ├── adapter // 适配器
│       └── driver  // 驱动
│           └── http  // 基于HTTP协议的服务器
├── handler // 用于处理请求的Handler(将协议升级为Websocket)
├── model   // 用于定义一些数据模型，比如请求、响应、错误、数据流等
├── pack    // 用于打包或生成数据的工具函数
├── router  // 注册路由
└── service // 用于处理业务逻辑的Service
```

# 2. 如果我想要添加一个新的功能，我应该在哪里添加？

## 2.1 添加一个新的工具函数

- 如果你想要添加一个新的`local`工具函数，你应该在`cli/ai/driver/eino/tools/local`目录下添加一个新的文件。
- 如果你想要添加一个新的`remote`工具函数，你应该在`cli/ai/driver/eino/tools/remote`目录下添加一个新的文件。

### 2.1.1 构建
请按照如下模板构建一个新的工具函数:

``` go
package local

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"

	"github.com/west2-online/DomTok/app/assistant/cli/ai/driver/eino/tools"
)

type ToolXXX struct {
	tool.InvokableTool // 可加可不加，但你必须实现InvokableRun和Info方法
}

const (
	ToolXXXName = "name" // 你的工具函数的名字
	ToolXXXDesc = "desc" // 你对工具函数的描述
)

type ToolXXXArgs struct {
	A1 string `json:"a1" yaml:"a1" desc:"a1" required:"true"`
} // 你的工具函数的参数，对于tag的解释，请查看cli/ai/driver/eino/tools/README.md

// 建立一个映射关系，用于描述你的工具函数的参数
// 注意，tool.Reflect的参数既可以是一个结构体本身，也可以是结构体的指针
// 你可以通过tool.Reflect(ToolXXXArgs{})或tool.Reflect(&ToolXXXArgs{})来得到这个映射关系
// 参数是一个指针时，tool.Reflect会认为尽可能展开这个指针，直到找到一个非指针的结构体
var ToolXXXRequestBody = schema.NewParamsOneOfByParams(*tools.Reflect(ToolXXXArgs{}))

// 工具构建函数，用于导出你的工具函数
func XXX() *ToolXXX {
	return &ToolXXX{}
}

// InvokableRun方法，用于实现你的工具函数，AI模型中的Function Calling会调用这个方法
func (t *ToolXXX) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	args := &ToolXXXArgs{}
	err := json.Unmarshal([]byte(argumentsInJSON), args)
	if err != nil {
		return "", err
	}
	return args.Message, nil
}

// Info方法，用于描述你的工具函数，让AI模型知道你的工具函数的名字、描述和参数
func (t *ToolXXX) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name:        ToolXXXName,
		Desc:        ToolXXXDesc,
		ParamsOneOf: ToolXXXRequestBody,
	}, nil
}
```

### 2.1.2 添加

在`cli/ai/driver/eino/request_props.go`中的`BuildTools`方法中添加你的工具函数:

``` go
// BuildTools builds the tools
func BuildTools(caller model.GetServerCaller) []tool.BaseTool {
	return []tool.BaseTool{
		local.XXX(),
	}
}
```

## 2.2 如果我想换一个AI模型平台，我应该在哪里修改？

通过实现`cli/ai/driver/eino/model.go`中的`BuildChatModel`方法来更换AI模型。

以下是Volcengine的一个例子:
``` go
ai := eino.NewClient()
ai.SetBuilder(func(ctx context.Context) (model.ChatModel, error) {
    return ark.NewChatModel(ctx, &ark.ChatModelConfig{
        APIKey:  config.Volcengine.ApiKey,
        BaseURL: config.Volcengine.BaseUrl,
        Region:  config.Volcengine.Region,
        Model:   config.Volcengine.Model,
    })
})
```

当然，代码不属于`app/assistant`，其位于[cmd/assistant/main.go](../../cmd/assistant/main.go)中。

Eino框架支持多种AI模型平台，通过这个链接查看更多的AI模型: [Eino](https://github.com/cloudwego/eino-ext/tree/main/components/model)

## 2.3 我发现Eino框架中没有我想要的AI模型平台适配，我应该怎么办？

你可以试着实现一个model.ChatModel接口，然后通过`cli/ai/driver/eino/model.go`中的`BuildChatModel`方法来实现你的AI客户端。

> 实在没有办法的话，在`cli/ai/driver`中添加一个新的目录，然后实现你的AI客户端。
> 
> 此时你会发现一个叫做`IDialog`的接口，这个接口是用作AI模型流输出与Websocket输出的桥梁；同时，在`model/dialog.go`中，你会发现一个叫做`Dialog`的结构体，这个结构体是一个`IDialog`的具体实现。
> 
> 需要注意以下几点:
>   - `IDialog`接口中
>     - `Unique()` 方法用于获取一个唯一的标识符，用于标识一个对话。
>     - `Message()` 方法用于获取一个消息，通常这个消息是用户在Websocket中一次性发送完成的消息。
>     - `Send(msg string)` 方法用于发送一个消息，通常这个消息是AI模型的输出，可以是多次发送。
>     - `Close()` 方法用于关闭一个对话，通常这个方法会在AI模型输出结束后调用。
>   -  `Dialog`结构体中
>     - `NotifyOnClosed()` 方法用于通知一个对话已经关闭，是一次性的单播通知。
>     - `NotifyOnMessage()` 方法用于通知一个消息已经发送，是一次性的单播通知。
>   - 需要注意的是
>     - 对`IDialog`的定义是单向的，即AI模型只能通过`IDialog`来发送消息，而Websocket只能通过`IDialog`来接收消息。
>       - 这个过程中，只能是AI模型"回复"用户的消息
>     - `Dialog`结构体是`IDialog`的具体实现，同时也是`IDialog`的一个代理，用于处理`IDialog`的通知。
>       - 这个过程中，用户需要"接收"AI模型的消息
>     - 你可以通过`IDialog`来实现一个双向的通信，但是这个过程中需要你自己来处理消息的发送和接收。
>     - 其实也不难理解，在现实对话中，为了保证对话的连贯性，通常是一问一答的，而不是同时问答。
>       - 对于问的问题，我们需要等待其完整的输出，然后再回答。
>       - 对于答的答案，我们不在乎其输出是一次性的还是流式的，毕竟AI的答案通常不会问我们问题。

## 2.4 远程服务可以对应不止一个服务器吗？

可以的，在AIClient的适配器中，可以看到`SetServerStrategy`方法，这个方法用于设置服务器的策略。

``` go
// AIClient is the interface for calling the AI
// It is used by the service to call the AI
type AIClient interface {
	// Call calls the AI with the dialog
	Call(ctx context.Context, dialog model.IDialog) error
	// ForgetDialog tells the AI to forget the dialog
	// This is used when the user logs out
	ForgetDialog(dialog model.IDialog)
	// SetServerStrategy sets the server category
	SetServerStrategy(strategy strategy.GetServerCaller) 
}
```

`strategy.GetServerCaller`是一个函数类型，其定义如下:

``` go
type GetServerCaller func(functionName string) adapter.ServerCaller
```

设置调用策略的代码不属于`app/assistant`，其位于[cmd/assistant/main.go](../../cmd/assistant/main.go)中，以下是个例子：

``` go
ai.SetServerStrategy(func(functionName string) adapter.ServerCaller {
    return http.NewClient(&http.ClientConfig{
        BaseUrl: ``,
    })
})
```

你可以通过这个函数类型来实现对不同服务器的调用策略，比如负载均衡、随机选择等。

## 2.5 如果服务器实现了一个新的功能，我应该在哪里添加？

你应该在`cli/server/adpater/adpater.go`中的`ServerCaller`中添加一个新的函数。
> 注意，AI Client使用的是`ServerCaller`接口，所以你应该在这个接口中添加新的函数。
> 
> 在这之后，你需要按照上文中的方法来把这个函数添加到`cli/ai/driver/eino/tools/remote`中。

``` go
type ServerCaller interface {
	Ping(ctx context.Context) ([]byte, error)

	// TODO: add more methods here
}
```

**不要忘记了在`cli/server/driver/http/server.go`中实现这个函数。**

## 2.6 如果我想用其他协议来实现服务器，我应该怎么做？

在`cli/server/driver`中添加一个新的目录，比如`grpc`，然后在这个目录中实现你的服务器。

你可以参考`cli/server/driver/http`中的实现，但是请注意，你需要实现`adapter.ServerCaller`接口。
