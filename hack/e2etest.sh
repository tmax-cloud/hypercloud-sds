#!/bin/bash

# shellcheck source=common.sh
. "$(dirname "$0")/common.sh"

function lint() {
  print_red "========================== run lint =========================="
  lint="golangci-lint run --timeout=30m --disable-all --enable=deadcode  --enable=gocyclo --enable=golint --enable=varcheck --enable=structcheck --enable=maligned --enable=errcheck --enable=dupl --enable=ineffassign --enable=interfacer --enable=unconvert --enable=goconst --enable=gosec --enable=megacheck --enable=lll --enable=whitespace --enable=gomnd"
  (
    cd "$testDir" && $lint
  )
  print_red "========================== ok run lint =========================="
}

function run() {
  print_red "========================== run test =========================="
  (
    cd "$testDir" && ginkgo
  )
  print_red "========================== ok run test =========================="
}

function main() {
  case "${1:-}" in
  run)
    run
    ;;
  lint)
    lint
    ;;
  *)
    echo "usage:" >&2
    echo "  $0 run" >&2
    echo "  $0 lint" >&2
    echo "  $0 help" >&2
    ;;
  esac
}

main "$1"
