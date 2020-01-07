#!/bin/bash

# include
. $(dirname "$0")/common.sh

function yaml() {
  $(dirname "$0")/rook.sh yaml
  $(dirname "$0")/cdi.sh yaml

  echo "Generated yaml in ./_deploy"
}

function install() {
  $(dirname "$0")/rook.sh install
  $(dirname "$0")/cdi.sh install
}

function uninstall() {
  $(dirname "$0")/cdi.sh uninstall
  $(dirname "$0")/rook.sh uninstall

  print_red "Warning: 모든 물리 노드의 /var/lib/rook 디렉토리를 삭제해야 합니다."
}

function test() {
  echo "e2e test ok"
}

function main() {
  case "${1:-}" in
  yaml)
    yaml
    ;;
  install)
    install
    ;;
  test)
    test
    ;;
  uninstall)
    uninstall
    ;;
  minikubeUp)
    $(dirname "$0")/cluster.sh minikubeUp
    ;;
  minikubeClean)
    $(dirname "$0")/cluster.sh minikubeClean
    ;;
  clusterUp)
    $(dirname "$0")/cluster.sh clusterUp
    ;;
  clusterClean)
    $(dirname "$0")/cluster.sh clusterClean
    ;;
  *)
    set +x
    echo "usage:" >&2
    echo "  $0 yaml" >&2
    echo "  $0 install" >&2
    echo "  $0 test" >&2
    echo "  $0 uninstall" >&2
    echo "  $0 minikubeUp" >&2
    echo "  $0 minikubeClean" >&2
    echo "  $0 clusterUp" >&2
    echo "  $0 clusterClean" >&2
    echo "  $0 help" >&2
    ;;
  esac
}

main $1
