# include
. $(dirname "$0")/common.sh

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
