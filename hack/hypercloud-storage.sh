#!/bin/bash

# include
. $(dirname "$0")/common.sh
. $(dirname "$0")/cluster.sh

function install() {
  $(dirname "$0")/cdi.sh install
}

function uninstall() {
  $(dirname "$0")/cdi.sh uninstall
}

function e2e() {
  echo "e2e test ok"
}

function main() {
  case "${1:-}" in
  install)
    install
    ;;
  uninstall)
    uninstall
    ;;
  clusterUp)
    clusterUp
    ;;
  clusterClean)
    clusterClean
    ;;
  e2e)
    e2e
    ;;
  *)
    set +x
    echo "usage:" >&2
    echo "  $0 install" >&2
    echo "  $0 uninstall" >&2
    echo "  $0 clusterUp" >&2
    echo "  $0 clusterClean" >&2
    echo "  $0 e2e" >&2
    echo "  $0 help" >&2
    ;;
  esac
}

main $1
