#!/bin/bash

# pre-requisites
# 1. clone project
lxc file push -rp . lxd-cm/home/lxd-cluster-manager/
# 2. Setup k8s cluster
snap install microk8s --classic --channel=1.31/stable
microk8s status --wait-ready
alias k="microk8s kubectl"
# 3. Setup lxd for rockcraft
lxd init --auto
# 4. Install rockcraft
snap install rockcraft --classic
# 5. Install and configure docker for running the rock
# assuming we have root access to the host, so no need to allow user with sudo privileges
snap install docker

# build and run rock with docker
deployment/rock/build.sh

set -e
