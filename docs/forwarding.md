# Forwarding

Use `-forward` to proxy all datagrams to an upstream DogStatsD server while still inspecting them locally:

```bash
dogstatsd-local -forward 127.0.0.1:8126
```

Datagrams are forwarded as-is — no parsing, no modification. This means upstream receives the exact same bytes your application sends, while you see the parsed output locally.

## Middleware Pattern

Run dogstatsd-local as a **middleware** between your application and the Datadog agent. Point your app at dogstatsd-local, and it transparently forwards everything upstream:

```
Application → dogstatsd-local (8125) → Datadog agent (8126)
```

```bash
dogstatsd-local -port 8125 -forward 127.0.0.1:8126
```

Your metrics continue flowing to Datadog, while you inspect them locally. No metrics are lost.

## Chaining Instances

You can chain multiple instances with different output formats — for example, one for pretty printing and one for raw passthrough:

```
Application → instance 1 (8125, pretty) → instance 2 (8126, raw) → Datadog agent (8127)
```

```bash
# Instance 1: pretty output, forwards to instance 2
dogstatsd-local -port 8125 -out pretty -forward 127.0.0.1:8126 &

# Instance 2: raw output, forwards to Datadog agent
dogstatsd-local -port 8126 -out raw -forward 127.0.0.1:8127 &
```

## Performance

Forwarding is non-blocking and buffered. If the upstream server is slow or unavailable, datagrams are buffered up to 10,000 messages. When the buffer is full, new datagrams are dropped and an error is logged. This ensures the local inspector is never blocked by a slow upstream.