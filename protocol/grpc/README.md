# Introduction


## Registration with the core service

```sh
go get github.com/enesanbar/go-service/protocol/grpc
```

```go
package main

import (
	"github.com/enesanbar/go-service/core/service"
	"github.com/enesanbar/go-service/protocol/grpc"
)

func main() {
	service.New(
		"my-service",
		grpc.Option(),
	).Run()
}
```

## Configuration
### Configuration Options

| Option                                               | Default Value     | Description                                                                                  |
| ---------------------------------------------------- | ----------------- | -------------------------------------------------------------------------------------------- |
| `tls.certFile`                                       | /etc/tls/tls.crt  | Path to TLS certificate file                                                                 |
| `tls.keyFile`                                        | /etc/tls/tls.key  | Path to TLS key file                                                                         |
| `tls.caFile`                                         | /etc/tls/ca.crt   | Path to TLS CA certificate file                                                              |
| `server.grpc.port`                                   | `50051`           | gRPC server port                                                                             |
| `server.grpc.gracefulStopTimeoutSeconds`             | `10`              | Graceful stop timeout (seconds)                                                              |
| `server.grpc.keepalive.timeSeconds`                  | `7200`            | Ping the client if it is idle for specified seconds to ensure the connection is still active |
| `server.grpc.keepalive.timeoutSeconds`               | `20`              | Wait for the ping ack before assuming the connection is dead                                 |
| `server.grpc.keepalive.minTimeSeconds`               | `300`             | Minimum time between pings (seconds)                                                         |
| `server.grpc.keepalive.permitWithoutStream`          | `false`           | Allow keepalive pings without active streams                                                 |
| `server.grpc.keepalive.maxConnectionIdleSeconds`     | `15`              | Max idle time before connection is closed                                                    |
| `server.grpc.keepalive.maxConnectionAgeSeconds`      | `30`              | Max age of a connection before forced closure                                                |
| `server.grpc.keepalive.maxConnectionAgeGraceSeconds` | `5`               | Additional grace period after max connection age                                             |
| `server.grpc.tls.enabled`                            | `false`           | Enable TLS for gRPC server                                                                   |
| `client.grpc.tls.enabled`                            | `false`           | Enable TLS for gRPC client                                                                   |
| `client.grpc.someservice.address`                    | `localhost:50505` | Address of the gRPC service                                                                  |
| `client.grpc.someservice.serviceConfig`              | See below         | gRPC client service config                                                                   |

**Default `serviceConfig`:**
```json
{
    "methodConfig": [
        {
            "name": [],
            "retryPolicy": {
                "maxAttempts": 3,
                "initialBackoff": "0.1s",
                "maxBackoff": "10s",
                "backoffMultiplier": 1.5,
                "retryableStatusCodes": ["UNAVAILABLE", "DEADLINE_EXCEEDED"]
            },
            "timeout": "15s"
        }
    ]
}
```

Configuration settings may be supplied either through a configuration file or via environment variables, with environment variables taking precedence. 

To translate entries from the configuration file into environment variables, please convert keys to uppercase and separate hierarchical levels with underscores. 

For example, the server can be configured using the following environment variables:
```sh
export SERVER_GRPC_PORT=8888
export SERVER_TLS_ENABLED=true
```

config/base.yaml
```yaml
tls:
  certFile: /certs/tls.crt
  keyFile: /certs/tls.key
  caFile: /certs/ca.crt

server:
  grpc:
    port: 50001
    gracefulStopTimeoutSeconds: 10
    keepalive:
      permitWithoutStream: true
      maxConnectionIdleSeconds: 3000
      maxConnectionAgeSeconds: 3000
      maxConnectionAgeGraceSeconds: 3000
    tls:
      enabled: false

client:
  grpc:
    tls:
      enabled: true
    someservice:
      address: localhost:50505
# default if not provided: "{\"methodConfig\":[{\"name\":[],\"retryPolicy\":{\"maxAttempts\":3,\"initialBackoff\":\"0.1s\",\"maxBackoff\":\"10s\",\"backoffMultiplier\":1.5,\"retryableStatusCodes\":[\"UNAVAILABLE\",\"DEADLINE_EXCEEDED\"]},\"timeout\":\"15s\"}]}" 
      serviceConfig:
        methodConfig:
          - name:
              - service: helloworld.Greeter
            timeout: 15s
            retryPolicy:
              maxAttempts: 10
              initialBackoff: 1s
              maxBackoff: 10s
              backoffMultiplier: 2.0
              retryableStatusCodes:
                - UNAVAILABLE
                - DEADLINE_EXCEEDED
```
