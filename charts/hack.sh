#!/bin/bash
set -euo pipefail
shopt -s inherit_errexit

srcDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/../"

. "${srcDir}/hack/util.sh"

function install() {
  minikube ssh "sudo mkdir -p /mnt/sda1/var/lib/rook" && helm install hypercloud-storage .
}

function uninstall() {
  helm delete hypercloud-storage
  minikube ssh "sudo rm -rf /var/lib/rook"
  minikube ssh "sudo rm -rf /mnt/sda1/var/lib/rook"
}

function wait_hypercloud_storage_available() {
  kubectl_wait_avail rook-ceph deployment/rook-ceph-operator 600
  kubectl_wait_avail rook-ceph deployment/csi-cephfsplugin-provisioner 600
  kubectl_wait_avail rook-ceph deployment/csi-rbdplugin-provisioner 600
  # TODO mon, mgr, osd check

  kubectl_wait_avail cdi deployment/cdi-apiserver 300
  kubectl_wait_avail cdi deployment/cdi-deployment 300
  kubectl_wait_avail cdi deployment/cdi-operator 300
  kubectl_wait_avail cdi deployment/cdi-uploadproxy 300
  kubectl_wait_avail cdi cdi.cdi.kubevirt.io/cdi 300

  echo "hypercloud_storage is now available"
}

function wait_hypercloud_storage_deleted() {
  kubectl_wait_delete cdi deployment/cdi-apiserver 300
  kubectl_wait_delete cdi deployment/cdi-deployment 300
  kubectl_wait_delete cdi deployment/cdi-uploadproxy 300
  kubectl_wait_delete cdi deployment/cdi-operator 300

  kubectl_wait_delete rook-ceph deployment/rook-ceph-operator 600
  kubectl_wait_delete rook-ceph deployment/csi-cephfsplugin-provisioner 600
  kubectl_wait_delete rook-ceph deployment/csi-rbdplugin-provisioner 600

  echo "hypercloud_storage is now deleted"
}

case "${1:-}" in
i)
  install
  ;;
u)
  uninstall
  ;;
t)
  helm template .
  ;;
wait_hypercloud_storage_available)
  wait_hypercloud_storage_available
  ;;
wait_hypercloud_storage_deleted)
  wait_hypercloud_storage_deleted
  ;;
*)
  echo "usage:" >&2
  echo "  $0 i                                     install chart on minikube" >&2
  echo "  $0 u                                     uninstall chart on minikube" >&2
  echo "  $0 t                                     template" >&2
  echo "  $0 wait_hypercloud_storage_available     wait until hypercloud_storage become available" >&2
  echo "  $0 wait_hypercloud_storage_deleted       wait until hypercloud_storage deleted" >&2
  ;;
esac
