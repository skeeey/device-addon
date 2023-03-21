#!/usr/bin/env bash

REPO_DIR="$(cd "$(dirname ${BASH_SOURCE[0]})/../.." ; pwd -P)"

cluster=cluster1

kubectl apply -f ${REPO_DIR}/contrib/demo/resources/device_data_model.yaml

kubectl get clustermanagementaddons
kubectl get devicedatamodels

kubectl apply -f ${REPO_DIR}/contrib/demo/resources/device_data_model.yaml

kubectl -n ${cluster} get managedclusteraddons

kubectl -n ${cluster} apply -f ${REPO_DIR}/contrib/demo/resources/device.yaml


kubectl -n mosquitto get svc mosquitto -ojsonpath='{.spec.ports[0].nodePort}'
kubectl get node edge-control-plane -ojsonpath='{.status.addresses[0].address}'