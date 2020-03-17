#!/bin/bash
set -eo pipefail

if [ "$EUID" -ne 0 ]; then
	echo "Please run as root"
	exit 1
fi

if [ -z "$1" ]; then
  echo "USAGE: $0 /path/to/k8s-cluster.yml"
  exit 1
fi

BINARY_DIR="/var/www/html/binary"
CLUSTER_YML="$1"

echo "Download location: ${BINARY_DIR}"

# download kubelet kubectl kubeadm
ARCH="amd64"
VERSION=$(grep ^kube_version ${CLUSTER_YML} | cut -d ' ' -f2 | tr -d '"')
LOC="storage.googleapis.com/kubernetes-release/release/${VERSION}/bin/linux/${ARCH}"
for FILE in "kubelet" "kubectl" "kubeadm"; do
	echo -e "\n[Download] ${FILE}:${VERSION}"
	mkdir -p ${BINARY_DIR}/${LOC}
	TARGET="${BINARY_DIR}/${LOC}/${FILE}"
	wget -q -O ${TARGET} "https://${LOC}/${FILE}"
	ls -lh ${TARGET}
	CHKSUM=$(sha256sum ${TARGET} | awk '{print $1}')
	echo -e "[SHA256SUM] ${FILE}:${VERSION} - ${CHKSUM}"
	case ${FILE} in
		"kubelet" )
			sed -i "s/^\(\s*kubelet_binary_checksum\s*:\s*\).*/\1${CHKSUM}/" ${CLUSTER_YML}
			;;
		"kubectl" )
			sed -i "s/^\(\s*kubectl_binary_checksum\s*:\s*\).*/\1${CHKSUM}/" ${CLUSTER_YML}
			;;
		"kubeadm" )
			sed -i "s/^\(\s*kubeadm_binary_checksum\s*:\s*\).*/\1${CHKSUM}/" ${CLUSTER_YML}
			;;
		* )
			echo "[ERROR] Not match binary file: ${FILE}"
			exit 1
			;;
	esac
	echo "------------------------------------------------------------"
done

# download etcd binary
VERSION=$(grep ^etcd_version ${CLUSTER_YML} | cut -d ' ' -f2 | tr -d '"')
LOC="github.com/coreos/etcd/releases/download/${VERSION}"
FILE="etcd-${VERSION}-linux-${ARCH}.tar.gz"
echo -e "\n[Download] ${FILE}:${VERSION}"
mkdir -p ${BINARY_DIR}/${LOC}
TARGET="${BINARY_DIR}/${LOC}/${FILE}"
wget -q -O ${TARGET} "https://${LOC}/${FILE}"
ls -lh ${TARGET}
CHKSUM=$(sha256sum ${TARGET} | awk '{print $1}')
echo -e "[SHA256SUM] ${FILE}:${VERSION} - ${CHKSUM}"
sed -i "s/^\(\s*etcd_binary_checksum\s*:\s*\).*/\1${CHKSUM}/" ${CLUSTER_YML}
echo "------------------------------------------------------------"

# download cni-plugins
VERSION=$(grep ^cni_version ${CLUSTER_YML} | cut -d ' ' -f2 | tr -d '"')
LOC="github.com/containernetworking/plugins/releases/download/${VERSION}"
FILE="cni-plugins-linux-${ARCH}-${VERSION}.tgz"
echo -e "\n[Download] ${FILE}:${VERSION}"
mkdir -p ${BINARY_DIR}/${LOC}
TARGET="${BINARY_DIR}/${LOC}/${FILE}"
wget -q -O ${TARGET} "https://${LOC}/${FILE}"
ls -lh ${TARGET}
CHKSUM=$(sha256sum ${TARGET} | awk '{print $1}')
echo -e "[SHA256SUM] ${FILE}:${VERSION} - ${CHKSUM}"
sed -i "s/^\(\s*cni_binary_checksum\s*:\s*\).*/\1${CHKSUM}/" ${CLUSTER_YML}
echo "------------------------------------------------------------"

# download calicoctl
VERSION=$(grep ^calico_ctl_version ${CLUSTER_YML} | cut -d ' ' -f2 | tr -d '"')
LOC="github.com/projectcalico/calicoctl/releases/download/${VERSION}"
FILE="calicoctl-linux-${ARCH}"
echo -e "\n[Download] ${FILE}:${VERSION}"
mkdir -p ${BINARY_DIR}/${LOC}
TARGET="${BINARY_DIR}/${LOC}/${FILE}"
wget -q -O ${TARGET} "https://${LOC}/${FILE}"
ls -lh ${TARGET}
CHKSUM=$(sha256sum ${TARGET} | awk '{print $1}')
echo -e "[SHA256SUM] ${FILE}:${VERSION} - ${CHKSUM}"
sed -i "s/^\(\s*calicoctl_binary_checksum\s*:\s*\).*/\1${CHKSUM}/" ${CLUSTER_YML}
echo "------------------------------------------------------------"

# download crictl
crictl_versions="v1.15.0 v1.16.1 v1.17.0"
kube_version=$(grep ^kube_version ${CLUSTER_YML} | cut -d ' ' -f2 | tr -d '"')
major_version=${kube_version::-2} # v1.15.3 -> v1.15
for VERSION in ${crictl_versions}; do
	if [[ ${VERSION} =~ ${major_version} ]]; then
		LOC="github.com/kubernetes-sigs/cri-tools/releases/download/${VERSION}"
		FILE="crictl-${VERSION}-linux-${ARCH}.tar.gz"
		echo -e "\n[Download] ${FILE}"
		mkdir -p ${BINARY_DIR}/${LOC}
		TARGET="${BINARY_DIR}/${LOC}/${FILE}"
		wget -q -O ${TARGET} "https://${LOC}/${FILE}"
		ls -lh ${TARGET}
		CHKSUM=$(sha256sum ${TARGET} | awk '{print $1}')
		echo -e "[SHA256SUM] ${FILE}:${VERSION} - ${CHKSUM}"
		sed -i "s/^\(\s*crictl_binary_checksum\s*:\s*\).*/\1${CHKSUM}/" ${CLUSTER_YML}
		echo "------------------------------------------------------------"
		break
	fi
done
