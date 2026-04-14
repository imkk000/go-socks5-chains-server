# Socks5 Proxy Forwarder

## Why?

I want to understand how proxy chains work and socks5 proxy forwarding. Try shuffle nodes, then exit node is tor proxy.

## Topology

```
Local -> Entry Node -> Middle Nodes -> Exit Node
```

## Shuffle Chains

```
# entry 9001
from: 127.0.0.1:9001 -> 127.0.0.1:9004

# random
from: 127.0.0.1:9004 -> 127.0.0.1:9003
from: 127.0.0.1:9003 -> 127.0.0.1:9005
from: 127.0.0.1:9005 -> 127.0.0.1:9002

# exit 9050 (tor)
from: 127.0.0.1:9002 -> 127.0.0.1:9050
```
