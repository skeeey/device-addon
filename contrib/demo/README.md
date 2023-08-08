# Demo

1. Run `make deploy` to prepare demo env
2. Run `make demo` to start the demo

```
#!/usr/bin/env bash

# https://demo.thingsboard.io/
# host.docker.internal

docker run -it -v ~/.tb-gateway/logs:/thingsboard_gateway/logs \
    -v ~/.tb-gateway/extensions:/thingsboard_gateway/extensions \
    -v ~/.tb-gateway/config:/thingsboard_gateway/config \
    --name tb-gateway \
    --restart always thingsboard/tb-gateway

# docker stop tb-gateway
# docker start tb-gateway

#mosquitto_pub -h 127.0.0.1 -p 1883 -t "/sensor/OPCUA-001/connect" -m ''
#mosquitto_pub -h 127.0.0.1 -p 1883 -t "/ocm/devices/OPCUA-001/data" -m '{"counter": 42, "random": 58}'
#mosquitto_pub -h 127.0.0.1 -p 1883 -t "/sensor/OPCUA-001/disconnect" -m ''
```
