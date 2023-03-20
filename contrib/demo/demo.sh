#!/usr/bin/env bash

REPO_DIR="$(cd "$(dirname ${BASH_SOURCE[0]})/../.." ; pwd -P)"

kubectl apply -f ${REPO_DIR}/contrib/demo/resources/device_data_model.yaml

kubectl get clustermanagementaddons
kubectl get devicedatamodels

kubectl apply -f ${REPO_DIR}/contrib/demo/resources/device_data_model.yaml

kubectl -n cluster1 get managedclusteraddons

kubectl -n cluster1 apply -f ${REPO_DIR}/contrib/demo/resources/device.yaml
