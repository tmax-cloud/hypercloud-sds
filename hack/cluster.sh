# include
. $(dirname "$0")/common.sh

function wait_for_ssh() {
  local tries=100
  while ((tries > 0)); do
    if minikube ssh echo connected &>/dev/null; then
      return 0
    fi
    tries=$((tries - 1))
    sleep 0.1
  done
  echo ERROR: ssh did not come up >&2
  exit 1
}

function clusterUp() {
  # TODO: minikube version
  # TODO: multinode cluster
  # TODO: minikube profile

  minikube start
  wait_for_ssh

  # Rook에서 사용할 디렉토리를 마운트
  minikube ssh "sudo mkdir -p /mnt/sda1/${PWD}; sudo mkdir -p $(dirname $PWD); sudo ln -s /mnt/sda1/${PWD} $(dirname $PWD)/"
  minikube ssh "sudo mkdir -p /mnt/sda1/var/lib/rook;sudo ln -s /mnt/sda1/var/lib/rook /var/lib/rook"
}

function clusterClean() {
  minikube delete
}
