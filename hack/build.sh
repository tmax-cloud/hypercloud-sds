# include
. $(dirname "$0")/common.sh

function build_go() {
  print_red "========================== build go =========================="
  (
  ls
  )
  print_red "========================== ok build go =========================="
}

function build_test() {
  print_red "========================== build test =========================="
  (
  # temporaily hardcoding to test simple case
   mkdir -p $testDir
   ls $testDir
 
   ls $pkgDir

   cd $pkgDir && go build -o $build_out .
  )
  print_red "========================== ok build test =========================="
}

function main() {
  case "${1:-}" in
  build_go)
    build_go
    ;;
  build_test)
    build_test
    ;;
  *)
    echo "usage:" >&2
    echo "  $0 build_go" >&2
    echo "  $0 build_test" >&2
    echo "  $0 help" >&2
    ;;
  esac
}

main $1
