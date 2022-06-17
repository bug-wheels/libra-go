# libra-go

可以在需要的平台运行 `go build` 来生成二进制包，或者通过 release 下载。

具体使用命令如下, 其中你要创建一个文件，用来描述要监控的服务信息和告警要求，具体的文件格式见下面。
最简单的使用方式为 `libra -c 文件路径`, 如果配置文件和程序同名且为 `libra.yaml` 则可不写

```
USAGE:
   libra [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config FILE, -c FILE         Load configuration from FILE (default: "libra.yaml")
   --help, -h                     show help (default: false)
   --intervalMS value             循环间隔时间，毫秒 (default: 10000)
   --loop                         是否循环执行 (default: true)
   --loopcount value, --lc value  循环次数，小于 1 则一直执行 (default: 0)
   --ui                           with ui (default: false)

```

你需要创建如下文件

libra.yaml

```yaml
registration-centers: # 注册中心的配置
  - name: dailynacos  # 名称，要求唯一
    type: nacos       # 注册中心类型，暂时只支持 nacos
    nacos: # nacos 的配置
      namespace:
      server-addr:
      group: DEFAULT_GROUP
      username:
      password:
services: # 要监控的服务
  - name: 用户服务                    # 告警时用于输出的名称，建议唯一 
    service-id: user-server          # 注册中心上对应的 serviceId 
    registration-center: dailynacos  # 注册中心的名称，对应上面的配置
    alarm: # 告警的配置，可不配，如果不配则不告警
      check: true                    # 是否进行检测，默认为 true, 可以不配置，当为 false 时，则不进行监控
      min-size: 1                    # 存活服务的的最小数量，低于最小数量则告警
      mustHosts: # 必须匹配上的 hosts 列表，如果配置的 hosts 中不在存活的实例中，则告警，可以不配 
        - 127.0.0.1
```

# 交叉编译

自己平台编译

```shell
go build -o libra-go main.go
```

编译为 Mac 平台

```shell
GOOS=darwin GOARCH=amd64 go build -o libra-go main.go
```

编译为 linux 平台

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o libra-go main.go
```

编译为 windows 平台

```shell
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o libra-go.exe main.go
```

