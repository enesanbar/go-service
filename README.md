[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/enesanbar/go-service)

# Go Service
## Introduction
Go Service is a base service project that supports
* **routing** and **middleware** support with [echo](https://echo.labstack.com/)
* ready-to-use and configurable http server with health checker 
* ready-to-use and configurable grpc server with health checker
* context-aware **structured logging** with [uber-go/zap](https://github.com/uber-go/zap)
* remote and local **configuration** with [viper](https://github.com/spf13/viper) 
* easy **dependency injection** and **application lifecycle management** with [uber-go/fx](https://github.com/uber-go/fx)
* automatic **prometheus** metrics and open telemetry tracing
* automatic instana sensor
* automatic profiling server for debugging purposes
* automatic **swagger ui** (given code is annotated and a swagger doc is generated at an expected location)
* in-memory **caching** with [go-cache](github.com/patrickmn/go-cache), redis and memcached support is coming soon...
* running **cron** tasks [BETA]
* connecting common **database** [BETA]
* **validation** mechanism [BETA]

go-service provides a ready-to-use, batteries included web and grpc server and lets you develop your service
without worrying about various pieces of common microservice components.

## ROADMAP
* Add common database connection methods with external pool configuration
* Create a generic pipeline to build go and docker images 

## Configuration
Viper is used for configuration, which means both config files and environment variables are supported with prefix to avoid conflict with other packages.

Given the following configuration
```yaml
debug: true
server:
  http:
	port: 8080
	context-path: /my-service
	healthcheck-path: /_hc
  grpc:
    port: 50231
```

the same configuration can be supplied via environment variables
```bash
DEBUG=true
SERVER_PORT=8080
SERVER_CONTEXT-PATH=/my-service
SERVER_HEALTHCHECK-PATH=/_hc
```

or using environment prefix to avoid collision with other packages
```bash
ENV_PREFIX=MY_ORG
MY_ORG_DEBUG=true
MY_ORG_SERVER_PORT=8080
MY_ORG_SERVER_CONTEXT-PATH=/my-service
MY_ORG_SERVER_HEALTHCHECK-PATH=/_hc
```

### Config source
go-service supports **consul** and **file (default)** as config sources

Expose CONFIG_SOURCE as an environment variable to use one of the config source. 

```shell
CONFIG_SOURCE=consul go run *.go
CONFIG_SOURCE=file go run *.go
```

### Environment 
Supported environments
* dev
* test
* staging
* prod

By default projects run in dev environment. You can specify an environment with
```shell
DEPLOY_TYPE=prod CONFIG_SOURCE=consul go run *.go
```

### Configuration file locations
> NOTE: It is not required to create config files

```yaml
debug: true 				# default: false
server:
  port: 8080 				# default: 9090
  context-path: /my-service # default: ""
  healthcheck-path: /_hc 	# default: health
```

If you specify 'test' as an environment variable with **file** as a config source, you can create following files to supply your configuration using yaml files
* ${PROJECT_DIR}/config/base.yaml
* ${PROJECT_DIR}/config/test.yaml

If you specify 'test' as an environment variable with **consul** as a config source, you can create the following at consul to supply your configuration remotely.
* go-config/my-service.yaml
* go-config/my-service_test.yaml

## Example projects

### Minimal example project (REST)
Following is a minimal unstructured code in a single file that is required to run a go-service project.

* Create an empty file in a directory of your choice called main.go
```shell
mkdir my-go-service
cd my-go-service
touch main.go
```

* Create your required configuration files

* Paste the following code block to main.go
```go
package main

import (
	"net/http"

	"github.com/enesanbar/go-service/router"
	"github.com/enesanbar/go-service/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func main() {
	service.New("my-service")
		WithRestAdapter().
		WithConstructor(
			NewMyHealthChecker,
			NewHandler,
			fx.Annotated{
				Group:  "routes",
				Target: RegisterRoutes,
			},
		).Build().Run()
}

// RegisterRoutes registers routes in the echo router
func RegisterRoutes(adapter *Handler) router.RouteConfig {
	return router.RouteConfig{
		Path: "/test",
		Router: func(group *echo.Group) {
			group.GET("/hello", adapter.Handle) // GET localhost:9090/test/hello
		},
	}
}

// MyHealthChecker runs at /service-name/_hc
type MyHealthChecker struct {
}

func NewMyHealthChecker() router.HealthChecker {
	return &MyHealthChecker{}
}

func (c *MyHealthChecker) Check() error {
	return nil
}

// Handler is a sample handler/controller
type Handler struct {
	router.BaseHandler
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(c echo.Context) error {
	return h.NewSuccess(c, "OK", http.StatusOK)
}
```

* Create go mod file and fetch dependencies
```shell
go mod init github.com/my-org/my-service
go mod tidy
```

* At this point, you should have the following directory structure
```
├── go.mod
├── go.sum
└── main.go
```

* Run the project depending on your configuration source choice
```shell
DEPLOY_TYPE=dev go run main.go
```

* Access the endpoints
```shell
curl -i localhost:9090/_hc
curl -i localhost:9090/test/hello
curl -i localhost:9092/metrics
```
