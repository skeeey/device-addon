#!/usr/bin/env bash

REPO_DIR="$(cd "$(dirname ${BASH_SOURCE[0]})/../.." ; pwd -P)"

demo_cluster="edge-demo"
demo_cluster_context="kind-${demo_cluster}"
demo_cluster_config=${REPO_DIR}/contrib/demo/kind/config.yaml

addon_namespace="open-cluster-management-agent-addon"
addon_deploy=${REPO_DIR}/contrib/deploy

opcua_server=${REPO_DIR}/contrib/demo/opcuaserver

rm -rf ${REPO_DIR}/_output

clusters_path=${REPO_DIR}/_output/clusters
kubeconfig=${clusters_path}/${demo_cluster}-kind.kubeconfig

mkdir -p ${clusters_path}

kind delete clusters ${demo_cluster}
kind create cluster --name=${demo_cluster} --config ${demo_cluster_config} --kubeconfig ${kubeconfig}
kind load docker-image quay.io/skeeey/device-addon --name=${demo_cluster}
kind load docker-image quay.io/skeeey/opcua-server --name=${demo_cluster}

export KUBECONFIG=${kubeconfig}

clusteradm init --context=${demo_cluster_context} --wait --feature-gates=AddonManagement=true --output-join-command-file join.sh
sh -c "$(cat ${REPO_DIR}/join.sh) cluster1 --feature-gates=AddonManagement=true --force-internal-endpoint-lookup --context ${demo_cluster_context}"
sleep 30
clusteradm accept --clusters cluster1 --context ${demo_cluster_context}

sleep 30
kubectl apply -k ${addon_deploy}
kubectl apply -k ${opcua_server}

unset KUBECONFIG
