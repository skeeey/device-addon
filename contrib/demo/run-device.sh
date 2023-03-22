#!/usr/bin/env bash

REPO_DIR="$(cd "$(dirname ${BASH_SOURCE[0]})/../.." ; pwd -P)"

KUBECONTEXT=${KUBECONTEXT:-"kind-edge-demo"}

port=$(kubectl --context ${KUBECONTEXT} -n mosquitto get svc mosquitto -ojsonpath='{.spec.ports[0].nodePort}')
addr=$(kubectl --context ${KUBECONTEXT} get node edge-demo-control-plane -ojsonpath='{.status.addresses[0].address}')

${REPO_DIR}/bin/thermometer "$1" "tcp://$addr:$port"
