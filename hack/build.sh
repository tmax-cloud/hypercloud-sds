# include
. $(dirname "$0")/common.sh

function build_go() {
  print_red "========================== build go =========================="
  (
    #TODO if go is not installed, install and pre-set ??
    checkIfGoInstalledCmd="go version"

    if ! $checkIfGoInstalledCmd; then
      echo "go need to be installed first"
      exit 1
    fi

    source /etc/profile #TODO this line is necessary to use go env, But this style is temporary!!!

    echo "enable module mode on for run go file outside of \$GOPATH"
    enableModuleMode="GO111MODULE=on"
    export $enableModuleMode
  )
  print_red "========================== ok build go =========================="
}

function build_prerequisites() {
  print_red "========================== build prerequisites =========================="
  (
    checkIfGoLintInstalledCmd="golangci-lint --version"

    if ! $checkIfGoLintInstalledCmd; then
      echo "golangci-lint need to be installed first by following command"
      echo "curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$(go env GOPATH)/bin v1.23.1"
      exit 1
    fi

    # dependent packages install
    pkgDownloadCmd="go get ./..."
    cd $testDir && $pkgDownloadCmd
  )
  print_red "========================== ok build prerequisites =========================="
}

function main() {
  case "${1:-}" in
  build_go)
    build_go
    ;;
  build_prerequisites)
    build_prerequisites
    ;;
  *)
    echo "usage:" >&2
    echo "  $0 build_go" >&2
    echo "  $0 build_prerequisites" >&2
    echo "  $0 help" >&2
    ;;
  esac
}

main $1
