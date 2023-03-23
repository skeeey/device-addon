#!/usr/bin/env bash

REPO_DIR="$(cd "$(dirname ${BASH_SOURCE[0]})/../.." ; pwd -P)"

demo_dir=${REPO_DIR}/contrib/demo
cluster=edge-node

source ${demo_dir}/demo_magic

comment "managed cluster and device addon"
pe "kubectl get managedclusters"
pe "kubectl get clustermanagementaddons"
pe "kubectl -n ${cluster} get managedclusteraddons"

comment "device data model"
pe "kubectl apply -f ${demo_dir}/resources/devicedatamodel.yaml"
pe "kubectl get devicedatamodels thermometer -oyaml"

comment "create devices"
pe "cat ${demo_dir}/resources/device-a.yaml"
pe "kubectl apply -n ${cluster} -f ${demo_dir}/resources/device-a.yaml"

pe "cat ${demo_dir}/resources/device-b.yaml"
pe "kubectl apply -n ${cluster} -f ${demo_dir}/resources/device-b.yaml"

comment "get devices from hub"
pe "kubectl -n ${cluster} get devices"

pe "kubectl -n ${cluster} get device room-a -oyaml"

pe "kubectl -n ${cluster} get device room-b -oyaml"
