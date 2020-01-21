# include
. $(dirname "$0")/common.sh

function build_go() {
  print_red "========================== build go =========================="
  (
   #TODO if go is not installed, install and pre-set ??
   
   source /etc/profile #TODO this line is necessary to use go env, But this style is temporary!!!
  )
  print_red "========================== ok build go =========================="
}

function build_test() {
  print_red "========================== build test =========================="
  (
  #TODO temporaily hardcoding to test simple case
   mkdir -p $testDir

   # 1st test case
   cd $pkgDir1 && go build -o $build_out1 .

  # 2nd test case
   cd $pkgDir2 && go build -o $build_out2 .
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
