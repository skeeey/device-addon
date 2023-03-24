#!/usr/bin/env bash

REPO_DIR="$(cd "$(dirname ${BASH_SOURCE[0]})/../.." ; pwd -P)"

demo_dir=${REPO_DIR}/contrib/demo

hub="edge-hub"
spoke="edge-node"

cluster="edge-node"

source ${demo_dir}/demo_magic

kubectl config set current-context kind-${hub}

comment "managed cluster and device addon on hub cluster"
pe "kubectl get managedclusters"
pe "kubectl get clustermanagementaddons"
pe "kubectl -n ${cluster} get managedclusteraddons"

comment "create device data model on hub"
pe "kubectl apply -f ${demo_dir}/resources/devicedatamodel.yaml"
pe "kubectl get devicedatamodels thermometer -oyaml"

comment "create devices on hub"
pe "cat ${demo_dir}/resources/device-a.yaml"
pe "kubectl apply -n ${cluster} -f ${demo_dir}/resources/device-a.yaml"

pe "cat ${demo_dir}/resources/device-b.yaml"
pe "kubectl apply -n ${cluster} -f ${demo_dir}/resources/device-b.yaml"

comment "get devices from hub"
pe "kubectl -n ${cluster} get devices"

pe "kubectl -n ${cluster} get device room-a -oyaml"

pe "kubectl -n ${cluster} get device room-b -oyaml"

comment "mqtt broker on the edge node"
pe "kubectl --context kind-${spoke} -n mosquitto get svc"

comment "ocm agents on the edge node"
pe "kubectl --context kind-${spoke} -n multicluster-controlplane-agent get pods"
