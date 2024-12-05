#!/bin/bash

set -e

# clone project
lxc file push -rp . lxd-cm/lxd-cluster-manager/

# Setup k8s cluster
snap install microk8s --classic --channel=1.31/stable
echo "export PATH=$PATH:/snap/bin" >> ~/.bashrc
source ~/.bashrc
microk8s status --wait-ready
microk8s enable registry
alias kubectl="microk8s kubectl"

# Setup lxd for rockcraft
lxd init --auto

# Install rockcraft
snap install rockcraft --classic

# Install and configure docker for running the rock
# assuming we have root access to the host, so no need to allow user with sudo privileges
snap install docker

# build the rock and push it to the docker daemon
# ensure the docker image is tagged with the microk8s built-in registry address
# load the image into the microk8s built-in registry
# NOTE: for production the image should be hosted on a public/private registry, we wouldn't need to build the images here
export IMAGE=localhost:32000/lxd-cluster-manager:1.0.0
deployment/rock/build.sh
docker push $IMAGE

# clean up build environment to free up resources
# NOTE: standard github hosted runners for private repos only have 7GB of memory
# https://docs.github.com/en/actions/using-github-hosted-runners/using-github-hosted-runners/about-github-hosted-runners
snap remove lxd --purge
snap remove rockcraft --purge

# ================================================================
# The following is the deployment script for the k8s cluster
# ================================================================

# 1. deploy cert-manager for the cluster and create certificate secrets
# Wait for Cert-Manager to be deployed and ready
echo "Installing cert-manager.."
kubectl apply -f deployment/k8s/cicd/cert/cert-manager.yaml
echo "Waiting for Cert-Manager deployment to become available..."
kubectl wait --for=condition=available --timeout=300s deployment --all -n cert-manager

# Apply the ClusterIssuer if not already created
echo "Applying ClusterIssuer..."
kubectl apply -f deployment/k8s/cicd/cert/cert-issuer.yaml

# Apply the Certificate resource
echo "Applying Certificates..."
kubectl apply -f deployment/k8s/cicd/cert/management-api-cert.yaml
kubectl apply -f deployment/k8s/cicd/cert/cluster-connector-cert.yaml

# Wait for the cert Secret to be created (it contains the certificate)
echo "Waiting for the certificate Secrets to be created..."
kubectl wait --for=create --timeout=600s secret/management-api-cert-secret -n default
kubectl wait --for=create --timeout=600s secret/cluster-connector-cert-secret -n default

echo "Certificates is ready!"
echo "Proceeding to application deployment..."

# 2. deploy Postgres database to the cluster
echo "Deploying Postgres database..."
kubectl apply -f deployment/k8s/cicd/db/config.yaml
kubectl apply -f deployment/k8s/cicd/db/pv.yaml
kubectl apply -f deployment/k8s/cicd/db/pvc.yaml
kubectl apply -f deployment/k8s/cicd/db/svc.yaml
kubectl apply -f deployment/k8s/cicd/db/ss.yaml
kubectl rollout status --watch --timeout=600s statefulset/db-ss
echo "Postgres database is ready!"

# 3. deploy configs to the cluster
echo "Deploying configs..."
kubectl apply -f deployment/k8s/cicd/config/config.yaml
kubectl wait --for=create --timeout=600s config/config -n default
echo "Configs is ready!"

# 4. deploy management-api to the cluster
echo "Deploying management-api..."
kubectl apply -f deployment/k8s/cicd/management-api/svc.yaml
kubectl apply -f deployment/k8s/cicd/management-api/depl.yaml
kubectl rollout status --watch --timeout=600s deployment/management-api-depl
echo "Management-api is ready!"

# 5. deploy cluster-connector to the cluster
echo "Deploying cluster-connector..."
kubectl apply -f deployment/k8s/cicd/cluster-connector/svc.yaml
kubectl apply -f deployment/k8s/cicd/cluster-connector/depl.yaml
kubectl rollout status --watch --timeout=600s deployment/cluster-connector-depl
echo "Cluster-connector is ready!"

# 6. Setup port-forwarding for the management-api and cluster-connector
# This is not production ready approach, we would have a single ingress setup with ssl passthrough ideally
# Since the backend e2e tests will restart the pods, we need to establish a persistent port-forwarding
echo "Exposing management-api and cluster-connector..."
(
    services=(
    "svc/management-api-svc:9000:management-api"
    "svc/cluster-connector-svc:9001:cluster-conn"
    )

    while true; do
    for svc in "${services[@]}"; do
        svc_name=$(echo $svc | cut -d':' -f1)
        local_port=$(echo $svc | cut -d':' -f2)
        target_port=$(echo $svc | cut -d':' -f3)
        
        # Check if port-forwarding is already active
        if ! lsof -i :$local_port > /dev/null; then
        echo "Reconnecting to $svc_name..."
        kubectl port-forward $svc_name $local_port:$target_port &
        fi
    done
    sleep 5
    done
) &
echo "management-api and cluster-connector is exposed on localhost:9000 and localhost:9001 respectively"

# ================================================================
# Run tests against the k8s cluster
# ================================================================

# 1. Setup dependencies for backend e2e tests
apt update
apt install -y make
snap install go --classic --channel=1.23/stable
make tidy-gomod

# 2. Run backend e2e tests
sed -i 's/kubectl/microk8s kubectl/g' makefile
microk8s config > ~/.kube/config # e2e tests needs to establish a client connection to the k8s cluster using the kubeconfig file
make test-e2e

# 3. Setup dependencies for frontend e2e tests
snap install node --channel=22/stable --classic
cd ui && yarn install --frozen=lockfile
npx playwright install --with-deps

# 4. Run frontend e2e tests
cd ..
# MUST setup oidc credentials for e2e tests
make test-ui-e2e CI=true