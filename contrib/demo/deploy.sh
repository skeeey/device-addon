#!/usr/bin/env bash

REPO_DIR="$(cd "$(dirname ${BASH_SOURCE[0]})/../.." ; pwd -P)"

hub="edge-hub"
spoke="edge-node"

rm -rf ${REPO_DIR}/_output

controlplane_path=${REPO_DIR}/_output/multicluster-controlplane
agent_deploy_path=${controlplane_path}/hack/deploy/agent
mosquitto_deploy_path=${REPO_DIR}/contrib/mosquitto
addon_deploy_path=${REPO_DIR}/contrib/deploy

hub_kubeconfig=${REPO_DIR}/_output/hub.kubeconfig

agent_namespace="multicluster-controlplane-agent"

kind delete clusters ${hub} ${spoke}
kind create cluster --name=${hub}
kind create cluster --name=${spoke}

echo "##### Clone multicluster-controlplane"
mkdir -p ${REPO_DIR}/_output/multicluster-controlplane
git clone --depth=1 https://github.com/open-cluster-management-io/multicluster-controlplane.git $REPO_DIR/_output/multicluster-controlplane

echo "##### Prepare hub..."
kubectl --context kind-${hub} config view --minify --flatten > ${hub_kubeconfig}
export KUBECONFIG=${hub_kubeconfig}
clusteradm init
unset KUBECONFIG

echo "##### Prepare spoke..."
kubectl --context kind-${hub} config view --minify --flatten > ${agent_deploy_path}/hub-kubeconfig
kubectl --kubeconfig ${agent_deploy_path}/hub-kubeconfig config set-cluster kind-${hub} --server=https://edge-hub-control-plane:6443

cp -f ${REPO_DIR}/contrib/demo/multicluster-controlplane-agent/deployment.yaml ${agent_deploy_path}/deployment.yaml

kubectl --context kind-${spoke} delete namespace ${agent_namespace} --ignore-not-found
kubectl --context kind-${spoke} create namespace ${agent_namespace}
kubectl --context kind-${spoke} -n ${agent_namespace} apply -k ${agent_deploy_path}

sleep 60

kubectl --context kind-${hub} patch managedcluster edge-node -p='{"spec":{"hubAcceptsClient":true}}' --type=merge
kubectl --context kind-${hub} get csr -l open-cluster-management.io/cluster-name=edge-node | grep Pending | awk '{print $1}' | xargs kubectl certificate approve

echo "##### Prepare addon..."
kubectl --context kind-${hub} apply -k ${addon_deploy_path}
