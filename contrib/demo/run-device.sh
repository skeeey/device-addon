#!/usr/bin/env bash

REPO_DIR="$(cd "$(dirname ${BASH_SOURCE[0]})/../.." ; pwd -P)"

spoke="edge-node"

port=$(kubectl --context kind-${spoke} -n mosquitto get svc mosquitto -ojsonpath='{.spec.ports[0].nodePort}')
addr=$(kubectl --context kind-${spoke} get node ${spoke}-control-plane -ojsonpath='{.status.addresses[0].address}')

${REPO_DIR}/bin/thermometer "$1" "tcp://$addr:$port"
