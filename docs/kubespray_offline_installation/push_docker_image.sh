#!/bin/bash

if [ "$EUID" -ne 0 ]; then
	echo "Please run as root"
	exit 1
fi

if [ "$#" -ne 2 ]; then
  echo "[Usage]: $0 {Docker Registry URL} {Path/to/k8s-cluster.yml}"
	echo "Example: $0 192.168.0.200:5000 inventory/mycluster/group-vars/k8s-cluster/k8s-cluster.yml"
  exit 1
fi

REG_URL="$1"
CLUSTER_YML="$2"
IMAGE_LIST=()

status_code=$(curl -I -k -s "${REG_URL}" | head -n 1 | cut -d ' ' -f 2)
if [[ "$status_code" != "200" ]]; then
	echo "[ERROR] Docker registry (${REG_URL}) is not running."
	exit 1
fi

set -eo pipefail
# Kubespary-based container image
TAG=$(grep ^nginx_image_tag "${CLUSTER_YML}" | cut -d ' ' -f2 | tr -d '"')
IMAGE_LIST+=("docker.io/library/nginx:${TAG}")

TAG=$(grep ^calico_version "${CLUSTER_YML}" | cut -d ' ' -f2 | tr -d '"')
IMAGE_LIST+=("docker.io/calico/node:${TAG}")

TAG=$(grep ^calico_cni_version "${CLUSTER_YML}" | cut -d ' ' -f2 | tr -d '"')
IMAGE_LIST+=("docker.io/calico/cni:${TAG}")

TAG=$(grep ^calico_policy_version "${CLUSTER_YML}" | cut -d ' ' -f2 | tr -d '"')
IMAGE_LIST+=("docker.io/calico/kube-controllers:${TAG}")

TAG=$(grep ^coredns_version "${CLUSTER_YML}" | cut -d ' ' -f2 | tr -d '"')
IMAGE_LIST+=("docker.io/coredns/coredns:${TAG}")

TAG=$(grep ^nodelocaldns_version "${CLUSTER_YML}" | cut -d ' ' -f2 | tr -d '"')
IMAGE_LIST+=("gcr.io/google-containers/k8s-dns-node-cache:${TAG}")

TAG=$(grep ^kube_version "${CLUSTER_YML}" | cut -d ' ' -f2 | tr -d '"')
IMAGE_LIST+=("gcr.io/google-containers/kube-proxy:${TAG}")
IMAGE_LIST+=("gcr.io/google-containers/kube-apiserver:${TAG}")
IMAGE_LIST+=("gcr.io/google-containers/kube-scheduler:${TAG}")
IMAGE_LIST+=("gcr.io/google-containers/kube-controller-manager:${TAG}")

TAG=$(grep ^dnsautoscaler_version "${CLUSTER_YML}" | cut -d ' ' -f2 | tr -d '"')
IMAGE_LIST+=("gcr.io/google-containers/cluster-proportional-autoscaler-amd64:${TAG}")

TAG=$(grep ^dashboard_image_tag "${CLUSTER_YML}" | cut -d ' ' -f2 | tr -d '"')
IMAGE_LIST+=("gcr.io/google_containers/kubernetes-dashboard-amd64:${TAG}")

TAG=$(grep ^pod_infra_version "${CLUSTER_YML}" | cut -d ' ' -f2 | tr -d '"')
IMAGE_LIST+=("gcr.io/google-containers/pause:${TAG}")
IMAGE_LIST+=("gcr.io/google_containers/pause-amd64:${TAG}")

TAG=$(grep ^etcd_version "${CLUSTER_YML}" | cut -d ' ' -f2 | tr -d '"')
IMAGE_LIST+=("quay.io/coreos/etcd:${TAG}")


for image in "${IMAGE_LIST[@]}"
do
	echo ""
	echo -e "\n[PULL]<- ${image}"
	docker pull "${image}"

	ORI_URL=$(echo "${image}" | cut -d '/' -f1)
	NEW_IMAGE="${image//${ORI_URL}/${REG_URL}}"

	echo -e "\n[PUSH]-> ${NEW_IMAGE}"
	docker tag "${image}" "${NEW_IMAGE}"
	docker push "${NEW_IMAGE}"
done
exit

