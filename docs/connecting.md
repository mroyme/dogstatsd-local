# Connecting Your Service

dogstatsd-local speaks the same DogStatsD protocol as the Datadog agent. Any library or service that can send metrics to Datadog can point to dogstatsd-local instead — just change the host and port.

## Datadog Agent (dd-agent)

If you're running the Datadog agent, set the DogStatsD socket in your agent config:

```yaml
# datadog.yaml
dogstatsd_socket: /var/run/datadog/dsd.sock
dogstatsd_port: 8126
```

Then point dogstatsd-local at the agent:

```bash
dogstatsd-local -port 8125 -forward 127.0.0.1:8126
```

## DogStatsD Client Libraries

Most Datadog client libraries accept a host and port. Point them at dogstatsd-local:

### Go (DataDog/datadog-go)

```go
client, _ := statsd.New("127.0.0.1:8125")
client.Gauge("my.metric", 42, []string{"env:dev"}, 1)
```

### Python (datadog-python)

```python
from datadog import DogStatsd

statsd = DogStatsd(host="127.0.0.1", port=8125)
statsd.gauge("my.metric", 42, tags=["env:dev"])
```

### Ruby (dogstatsd-ruby)

```ruby
statsd = Datadog::Statsd.new("127.0.0.1", 8125)
statsd.gauge("my.metric", 42, tags: ["env:dev"])
```

### Java (java-dogstatsd-client)

```java
StatsDClient client = new NonBlockingStatsDClientBuilder()
    .hostname("127.0.0.1")
    .port(8125)
    .build();
client.gauge("my.metric", 42, "env:dev");
```

### Node.js (hot-shots)

```javascript
const StatsD = require("hot-shots");
const client = new StatsD({ host: "127.0.0.1", port: 8125 });
client.gauge("my.metric", 42, ["env:dev"]);
```

## Environment Variables

Many services and libraries read from standard environment variables:

| Variable                         | Example     | Description                             |
| -------------------------------- | ----------- | --------------------------------------- |
| `DD_AGENT_HOST`                  | `127.0.0.1` | DogStatsD host                          |
| `DD_DOGSTATSD_PORT`              | `8125`      | DogStatsD port                          |
| `DD_ENTITY_ID`                   | _(auto)_    | Datadog entity ID for container tagging |
| `DD_DOGSTATSD_NON_LOCAL_TRAFFIC` | `true`      | Allow non-local traffic (agent only)    |

To redirect a service to dogstatsd-local:

```bash
export DD_AGENT_HOST=127.0.0.1
export DD_DOGSTATSD_PORT=8125
```

## Docker Compose

Point your app container at the dogstatsd-local container:

```yaml
services:
  dogstatsd-local:
    image: mroyme/dogstatsd-local
    ports:
      - "8125:8125/udp"

  app:
    image: my-app
    environment:
      DD_AGENT_HOST: dogstatsd-local
      DD_DOGSTATSD_PORT: 8125
```

## Kubernetes

Set the environment variables in your pod spec to point at dogstatsd-local:

```yaml
env:
  - name: DD_AGENT_HOST
    value: "dogstatsd-local.default.svc.cluster.local"
  - name: DD_DOGSTATSD_PORT
    value: "8125"
```

Or run dogstatsd-local as a sidecar:

```yaml
sidecars:
  - name: dogstatsd-local
    image: mroyme/dogstatsd-local
    ports:
      - containerPort: 8125
        protocol: UDP
    env:
      - name: DD_AGENT_HOST
        value: "dogstatsd-local"
      - name: DD_DOGSTATSD_PORT
        value: "8125"
```

## Testing with netcat

For quick manual testing, you can send raw DogStatsD datagrams via netcat:

```bash
# Metric
printf "page.views:1|c|#env:dev" | nc -u -w1 localhost 8125

# Service check
printf "_sc|Redis connection|2|#env:dev|m:Timeout" | nc -u -w1 localhost 8125

# Event
printf "_e{21,21}:An exception occurred|Cannot parse CSV file|t:warning|#err_type:bad_file" | nc -u -w1 localhost 8125
```

This is useful for verifying that dogstatsd-local is running and parsing correctly, but for real services, use a proper DogStatsD client library or environment variables as described above.
