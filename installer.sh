#!/bin/bash
set -euo pipefail
shopt -s inherit_errexit

installDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

function wait_condition {
  cond=$1
  timeout=$2

  for ((i=0; i<timeout; i+=5)) do
    echo "Waiting for ${i}s condition: \"$cond\""
    if eval $cond > /dev/null 2>&1; then echo "Conditon met"; return 0; fi;
    sleep 5
  done

  echo "Condition timeout"
  return 1
}

function install {
  inventory=$1

  kustomize build inventory/$inventory/operators | kubectl apply -f -
  sleep 10
  kubectl wait --for=condition=available deployment rook-ceph-operator -n rook-ceph --timeout=180s
  kubectl wait --for=condition=available deployment cdi-operator -n cdi --timeout=180s

  kustomize build inventory/$inventory/resources | kubectl apply -f -
  sleep 10
  wait_condition "kubectl get cephclusters.ceph.rook.io -n rook-ceph | grep Created" 360
  wait_condition "kubectl get pod -n rook-ceph | grep osd" 120
  kubectl wait --for=condition=available deployment cdi-apiserver --timeout=60s -n cdi
  kubectl wait --for=condition=available deployment cdi-uploadproxy --timeout=60s -n cdi
  kubectl wait --for=condition=available deployment cdi-deployment --timeout=60s -n cdi
}

function uninstall {
  inventory=$1

  kustomize build inventory/$inventory/resources | kubectl delete -f -
  sleep 10
  wait_condition "! kubectl get cephclusters.ceph.rook.io -n rook-ceph rook-ceph" 180
  wait_condition "! kubectl get cdis.cdi.kubevirt.io -n cdi cdi" 180

  kustomize build inventory/$inventory/operators | kubectl delete -f -
  wait_condition "! kubectl get ns | grep cdi" 180
  wait_condition "! kubectl get ns | grep rook-ceph" 180

  if [ "$inventory" == "minikube" ]; then
    minikube ssh "sudo rm -rf /var/lib/rook" && minikube ssh "sudo rm -rf /data/osd*"
  fi

  echo "Warning: rook의 데이터가 남아있을 수 있습니다. 문서를 참조하세요."
}

function test {
  cd ./e2e && ginkgo # go test
}

case "${1:-}" in
install)
  install ${2:?error: inventory가 지정되지 않았습니다.}
  ;;
uninstall)
  uninstall ${2:?error: inventory가 지정되지 않았습니다.}
  ;;
test)
  test
  ;;
*)
  echo "usage:" >&2
  echo "  $0 install [inventory]" >& 2
  echo "  $0 uninstall [inventory]" >& 2
  echo "  $0 test" >& 2
  ;;
esac
