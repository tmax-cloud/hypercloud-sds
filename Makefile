CLUSTER=./hack/cluster.sh
HYPERCLOUD_STORAGE=./hack/hypercloud-storage.sh

install:
	${HYPERCLOUD_STORAGE} install

uninstall:
	${HYPERCLOUD_STORAGE} uninstall

unit:
	echo "unit test"

e2e:
	${HYPERCLOUD_STORAGE} e2e

build:
	echo "build"

clusterUp:
	${HYPERCLOUD_STORAGE} clusterUp
	kubectl get nodes

clusterClean:
	${HYPERCLOUD_STORAGE} clusterClean
