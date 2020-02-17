#!/bin/bash
set -euo pipefail
shopt -s inherit_errexit

function install {
  minikube ssh "sudo mkdir -p /mnt/sda1/var/lib/rook" && helm install hypercloud-storage .
}

function uninstall {
  helm delete hypercloud-storage; minikube ssh "sudo rm -rf /var/lib/rook"; minikube ssh "sudo rm -rf /mnt/sda1/var/lib/rook"
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
*)
  echo "usage:" >&2
  echo "  $0 i   install" >&2
  echo "  $0 u   uninstall" >&2
  echo "  $0 t   template" >&2
  ;;
esac
