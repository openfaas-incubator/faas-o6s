#!/usr/bin/env bash

set -o errexit

REPO_ROOT=$(git rev-parse --show-toplevel)

logs() {
  kubectl -n openfaas logs deployment/gateway -c operator
}
trap "logs" EXIT SIGINT

echo ">>> Create test secrets"
kubectl -n openfaas-fn create secret generic faas-token --from-literal=faas-token=token
kubectl -n openfaas-fn create secret generic faas-key --from-literal=faas-key=key

echo ">>> Create test function"
kubectl apply -f ${REPO_ROOT}/artifacts/nodeinfo.yaml

echo '>>> Waiting for function deployment'
retries=10
count=0
ok=false
until ${ok}; do
    kubectl -n openfaas-fn get deployment/nodeinfo | grep 'nodeinfo' && ok=true || ok=false
    sleep 5
    count=$(($count + 1))
    if [[ ${count} -eq ${retries} ]]; then
        echo "No more retries left"
        exit 1
    fi
done

echo '>>> Waiting for function deployment to be ready'
kubectl -n openfaas-fn rollout status deployment/nodeinfo
