# include
. $(dirname "$0")/common.sh

func lint() {
#  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.23.1
# 'golangci-lint --version' check

# cd srcDir
# go fmt
# golangci-lint run --timeout=30m \
#  --disable-all --enable=deadcode  --enable=gocyclo --enable=golint --enable=varcheck \
#  --enable=structcheck --enable=maligned --enable=errcheck --enable=dupl --enable=ineffassign \
#  --enable=interfacer --enable=unconvert --enable=goconst --enable=gosec --enable=megacheck --enable=lll --enable=whitespace --enable=gomnd
# => $? check

}
function run() {
  print_red "========================== run test =========================="
  (
    $build_out1
    $build_out2
  )
  print_red "========================== ok run test =========================="
}

function result() {
  print_red "========================== result of test =========================="
  (
   #TODO use some nice tool
    ls
  )
  print_red "========================== result of test =========================="

}


function main() {
  case "${1:-}" in
  run)
    run
    ;;
  result)
    result
    ;;
  *)
    echo "usage:" >&2
    echo "  $0 run" >&2
    echo "  $0 result" >&2
    echo "  $0 help" >&2
    ;;
  esac
}

main $1
