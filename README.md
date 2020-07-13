# nacos-grpc

实现基础的nacas resolver, 定制LoadBalance

## 服务端

在grpc基础创建完以后, 按照约定注册到nacos即可

## 客户端

主要是要处理下target

```golang
con, err := grpc.Dial(nacosgrpc.Target("http://nacos:nacos@127.0.0.1:8848/nacos", "hello"), grpc.WithInsecure(), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, weightedroundrobin.Name)))
```

```golang
func Target(nacosAddr string, serviceName string, ops ...Option) string
```

Options对应nacos注册

- OptionGroupName 分组 (DEFAULT_GROUP)
- OptionNameSpaceID 命名空间 (pubic)
- OptionClusters 集群()
- OptionModeHeartBeat 更新源使用心跳方式(默认）
- OptionModeSubscribe 更新源使用订阅模式

## 测试

待补

weighted_load_balance

```text
50003-25.000000 count: 157
50002-25.000000 count: 158
50001-50.000000 count: 684
```
