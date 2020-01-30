#!/bin/bash

# shellcheck source=common.sh
. "$(dirname "$0")/common.sh"

cdiDeployDir="${deployDir}/cdi"

function cdi_yaml() {
  cdiDockerRegistry=${dockerRegistry:-kubevirt}

  mkdir -p "${cdiDeployDir}"
  # shellcheck disable=SC2086
  cp -r ${templatesDir}/cdi ${deployDir}

  sed -i -- "s|{{.cdiDockerRegistry}}|${cdiDockerRegistry}|g" "${cdiDeployDir}"/cdi-operator.yaml.in

  mv "${cdiDeployDir}"/cdi-operator.yaml.in "${cdiDeployDir}"/cdi-operator.yaml
}

function cdi_install() {
  print_red "========================== install cdi =========================="
  (
    set -x
    kubectl create -f "${cdiDeployDir}"/cdi-operator.yaml
    kubectl create -f "${cdiDeployDir}"/cdi-cr.yaml

    # TODO change to use go client ?
    kubectl_wait_avail cdi deployment/cdi-apiserver 300
    kubectl_wait_avail cdi deployment/cdi-deployment 300
    kubectl_wait_avail cdi deployment/cdi-operator 300
    kubectl_wait_avail cdi deployment/cdi-uploadproxy 300
    kubectl_wait_avail cdi cdi.cdi.kubevirt.io/cdi 300
  )
  print_red "========================== ok install cdi =========================="
}

function cdi_uninstall() {
  print_red "========================== uninstall cdi =========================="
  (
    set +eo pipefail
    set -x
    kubectl delete --ignore-not-found=true --wait=true -f "${cdiDeployDir}"/cdi-cr.yaml
    kubectl_wait_delete cdi deployment/cdi-apiserver 300
    kubectl_wait_delete cdi deployment/cdi-deployment 300
    kubectl_wait_delete cdi deployment/cdi-uploadproxy 300

    kubectl delete --ignore-not-found=true -f "${cdiDeployDir}"/cdi-operator.yaml
    kubectl_wait_delete cdi deployment/cdi-operator 300
  )
  print_red "========================== ok uninstall cdi =========================="
}

function main() {
  case "${1:-}" in
  yaml)
    cdi_yaml
    ;;
  install)
    cdi_install
    ;;
  uninstall)
    cdi_uninstall
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
