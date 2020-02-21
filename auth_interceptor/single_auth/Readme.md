##	普通方法:一元拦截器 grpc.UnaryInterceptor

```
grpc.UnaryInterceptor
func UnaryInterceptor(i UnaryServerInterceptor) ServerOption {
    return func(o *options) {
        if o.unaryInt != nil {
            panic("The unary server interceptor was already set and may not be reset.")
        }
        o.unaryInt = i
    }
}
type UnaryServerInterceptor func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (resp interface{}, err error)
```

通过查看源码可得知,要完成一个拦截器需要实现 UnaryServerInterceptor 方法.形参如下:

-	ctx context.Context:请求上下文
-	req interface{}:RPC 方法的请求参数
-	info *UnaryServerInfo:RPC 方法的所有信息
-	handler UnaryHandler:RPC 方法本身


##	实现多个拦截器

https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/chain.go 已经有了实现

```
import "github.com/grpc-ecosystem/go-grpc-middleware"

myServer := grpc.NewServer(
    grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
        ...
    )),
    grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
       ...
    )),
)
```

##	思考:

如何理解Interceptor方法中的handler?
