# Host服务模块


## IMPL

这个模块实现完成， Host Service的具体实现，上层业务就基于Service进行编程，面向接口开发

```bash
http
  |
  Host Service
    |
    impl(基于MySQL实现)
```

Host Service定义并实现之后，有四种用途：
- 用于内部模块调用，基于它封装更高一层的业务逻辑，比如：服务发布
- Host Service对外暴露：http协议（暴露给用户）
- Host Service对内暴露：gRPC框架（暴露给内部服务）
- ...