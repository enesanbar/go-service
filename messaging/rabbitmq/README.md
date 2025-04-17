# Create RabbitMQ Object using config

When the constructor functions run, they create and return map of all objects.

Here's example configuration

```yaml
datasources:
  rabbitmq:
    connections:
      default:
        # host: localhost # default
        # port: "5672" # default
        username: guest
        password: guest
    
    channels:
      default:
        connection: default
        mandatory: false
        immediate: false

    exchanges:
      default:
        channel: default
        exclusive: false
        type: direct
        durable: true
        auto-delete: false
        internal: false
        no-wait: false

    queues:
      default:
        channel: default
        durable: true
        auto-delete: false
        exclusive: false
        no-wait: false

    bindings:
      - exchange: default
        queue: default
        no-wait: false
        routing-keys:
          - EventName1
          - EventName2
          - EventName3
```

This configuration generates following objects, which can be injected into any constructors. The binding are automatically created.
```go
map[string]*Queue{
    "default": *Queue{}
}

map[string]*Exchange{
    "default": *Exchange{}
}

map[string]*Connection{
    "default": *Connection{}
}
```
