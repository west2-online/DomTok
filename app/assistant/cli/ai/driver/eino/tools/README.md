# 在此留言

你可以看到这个文件夹有两个子目录: `local`和`remote`

- `local`对应本地调用，AI模型调用`function calling`时，不会触发对其他服务器的调用。
- `remote`对应远程调用，AI模型调用`function calling`时，会触发对其他服务器的调用。

尽管`local`能够通过共用包解决，但鉴于不同AI模型的接口参数与响应参数不同。依然需要专门实现一个`local`目录。

> 如果AI服务共用一套规范，比如`openapi 3.0`, 那么不管是`local`还是`remote`都可以通过共用包解决。
> 这当然是以后的事情了。现在我们只需要专注于`local`和`remote`的实现即可。
