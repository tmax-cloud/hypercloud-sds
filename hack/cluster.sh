# include
. $(dirname "$0")/common.sh

function clusterUp() {
  # TODO: minikube version
  # TODO: multinode cluster
  # TODO: minikube profile

  minikube start
}

function clusterClean() {
  minikube delete
}
