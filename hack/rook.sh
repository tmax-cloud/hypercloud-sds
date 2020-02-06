#!/bin/bash

# shellcheck source=common.sh
. "$(dirname "$0")/common.sh"

rookDeployDir="${deployDir}/rook"

function rook_yaml() {
  rookDockerRegistry=${dockerRegistry:-rook}
  cephCsiDockerRegistry=${dockerRegistry:-quay.io/cephcsi}
  csiDockerRegistry=${dockerRegistry:-quay.io/k8scsi}
  cephDockerRegistry=${dockerRegistry:-ceph}

  mkdir -p "${rookDeployDir}"
  cp -r "${templatesDir}"/rook "${deployDir}"

  sed -i -- "s|{{.rookDockerRegistry}}|${rookDockerRegistry}|g" "${rookDeployDir}"/cluster/operator.yaml.in
  sed -i -- "s|{{.cephCsiDockerRegistry}}|${cephCsiDockerRegistry}|g" "${rookDeployDir}"/cluster/operator.yaml.in
  sed -i -- "s|{{.csiDockerRegistry}}|${csiDockerRegistry}|g" "${rookDeployDir}"/cluster/operator.yaml.in
  sed -i -- "s|{{.cephDockerRegistry}}|${cephDockerRegistry}|g" "${rookDeployDir}"/cluster/cluster-test.yaml.in
  sed -i -- "s|{{.rookDockerRegistry}}|${rookDockerRegistry}|g" "${rookDeployDir}"/cluster/toolbox.yaml.in

  mv "${rookDeployDir}"/cluster/operator.yaml.in "${rookDeployDir}"/cluster/operator.yaml
  mv "${rookDeployDir}"/cluster/cluster-test.yaml.in "${rookDeployDir}"/cluster/cluster-test.yaml
  mv "${rookDeployDir}"/cluster/toolbox.yaml.in "${rookDeployDir}"/cluster/toolbox.yaml
}

function rook_install() {
  print_red "========================== install rook =========================="
  (
    set -x

    # Deploy cluster
    kubectl create -f "${rookDeployDir}"/cluster/common.yaml
    kubectl create -f "${rookDeployDir}"/cluster/operator.yaml
    # TODO wait ?
    kubectl create -f "${rookDeployDir}"/cluster/cluster-test.yaml
    kubectl create -f "${rookDeployDir}"/cluster/toolbox.yaml

    # TODO change to use go client ?
    kubectl_wait_avail rook-ceph deployment/rook-ceph-operator 600
    kubectl_wait_avail rook-ceph deployment/csi-cephfsplugin-provisioner 600
    kubectl_wait_avail rook-ceph deployment/csi-rbdplugin-provisioner 600

    # Deploy rbd, cephfs
    kubectl create -f "${rookDeployDir}"/rbd/storageclass-test.yaml
    kubectl create -f "${rookDeployDir}"/rbd/snapshotclass.yaml
    kubectl create -f "${rookDeployDir}"/cephfs/storageclass.yaml
    kubectl create -f "${rookDeployDir}"/cephfs/filesystem-test.yaml

    # set default sc
    kubectl patch storageclass csi-cephfs -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'

    # TODO: wait ceph -s
  )
  print_red "========================== ok install rook =========================="
}

function rook_uninstall() {
  print_red "========================== uninstall rook =========================="
  (
    set +eo pipefail
    set -x

    kubectl delete --ignore-not-found=true -f "${rookDeployDir}"/cephfs/filesystem-test.yaml
    kubectl delete --ignore-not-found=true -f "${rookDeployDir}"/cephfs/storageclass.yaml
    kubectl delete --ignore-not-found=true -f "${rookDeployDir}"/rbd/snapshotclass.yaml
    kubectl delete --ignore-not-found=true -f "${rookDeployDir}"/rbd/storageclass-test.yaml
    kubectl delete --ignore-not-found=true -f "${rookDeployDir}"/cluster/toolbox.yaml
    kubectl delete --ignore-not-found=true -f "${rookDeployDir}"/cluster/cluster-test.yaml
    kubectl delete --ignore-not-found=true -f "${rookDeployDir}"/cluster/operator.yaml

    kubectl_wait_delete rook-ceph deployment/rook-ceph-operator 600

    kubectl delete --ignore-not-found=true -f "${rookDeployDir}"/cluster/common.yaml

    kubectl_wait_delete rook-ceph deployment/csi-cephfsplugin-provisioner 600
    kubectl_wait_delete rook-ceph deployment/csi-rbdplugin-provisioner 600
  )
  print_red "========================== ok uninstall rook =========================="
}

function main() {
  case "${1:-}" in
  yaml)
    rook_yaml
    ;;
  install)
    rook_install
    ;;
  uninstall)
    rook_uninstall
    ;;
  *)
    echo "usage:" >&2
    echo "  $0 yaml" >&2
    echo "  $0 install" >&2
    echo "  $0 uninstall" >&2
    echo "  $0 help" >&2
    ;;
  esac
}

main "$1"
