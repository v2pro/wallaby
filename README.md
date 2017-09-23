# wallaby

service mesh tackles non-functional challenges in one place.

![routing process](https://docs.google.com/drawings/d/e/2PACX-1vRAsCkJMutbN8DfH1atFLET15yzYGOwMQ0JLFQrvbg3tuXq71fCk5WF56xR0rBoXTVxtAYavD9fVJM_/pub?w=1011&h=764)

Overall Proxy Sequence

ServerConn => ConnForwardingDecision => ServerRequest => ClientRequest => ServiceKinds => ServiceInstance => RoutingDecision

there are three routing modes

* per connection routing: RoutingDecision is determined by ServerConn. 
This mode is most generic, can handle any kind of tcp stream without knowing the protocol
* per stream routing: RoutingDecision is determined by first request packet in the connection 
or stream (when protocol is multiplex, there might be multiple streams on one connection)
this mode do not need to do stateful protocol handling, and can route with more information
* per packet routing (a.k.a RPC mode): RoutingDecision might be different for different request packet
this mode is most powerful and most costly, need complete implementation of protocol
including encoding/decoding/stateful action sequences

the routing decision process works like this:

* ServerConn => ConnForwardingDecision: when a tcp connection is established,
how to forward the connection (routing mode/protocol) is determined
* ConnForwardingDecision => ServerRequest: parse request arrived server
* ServerRequest => ClientRequest: by parsing the request, we know what is the target service
* ClientRequest => ServiceKinds: one service have many clusters, filter out feasible clusters by cluster routing table.
* ServiceKinds => ServiceInstance: choose one most optimal service cluster from many clusters, 
choose one most optimal service instance from that cluster
* ServiceInstance => RoutingDecision: given the service status, should we accept/reject/wait the request. If accept,
handle the request by chosen service instance.

# User Guide

## 1 get dependencies (require dep)

```bash
dep ensure
```