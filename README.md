# device-addon

60n32XY8SXaTjUiRCAEb
nZ7EgGjpI7DJKm6qSC2N

docker run -it \
    --env HTTP_PROXY="http://squid.corp.redhat.com:3128" \
    --env HTTPS_PROXY="http://squid.corp.redhat.com:3128" \
    --env NO_PROXY="localhost,127.0.0.0/8" \
    --network host \
    -v ~/.tb-gateway/logs:/thingsboard_gateway/logs \
    -v ~/.tb-gateway/extensions:/thingsboard_gateway/extensions \
    -v ~/.tb-gateway/config:/thingsboard_gateway/config \
    --name tb-gateway \
    --restart always thingsboard/tb-gateway

docker run -it \
    --network host \
    -v ~/.tb-gateway/logs:/thingsboard_gateway/logs \
    -v ~/.tb-gateway/extensions:/thingsboard_gateway/extensions \
    -v ~/.tb-gateway/config:/thingsboard_gateway/config \
    --name tb-gateway \
    --restart always thingsboard/tb-gateway
