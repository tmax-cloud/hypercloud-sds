CLUSTER=./hack/cluster.sh
HYPERCLOUD_STORAGE=./hack/hypercloud-storage.sh

yaml:
	${HYPERCLOUD_STORAGE} yaml

install:
	${HYPERCLOUD_STORAGE} install

uninstall:
	${HYPERCLOUD_STORAGE} uninstall

clusterUp:
	${HYPERCLOUD_STORAGE} clusterUp
	kubectl get nodes

clusterClean:
	${HYPERCLOUD_STORAGE} clusterClean

test:
	${HYPERCLOUD_STORAGE} test

help:
	@echo "Usage: make [Target ...]"
	@echo "  yaml           cluster_config로 부터 설치를 위한 yaml 파일을 생성합니다."
	@echo "  install        생성된 yaml 파일로 hypercloud-storage를 설치합니다."
	@echo "  test           hypercloud-storage가 잘 설치되었는지 확인합니다. (end-to-end 테스트 수행)"
	@echo "  uninstall      hypercloud-storage를 제거합니다."
	@echo "  clusterUp      테스트를 위한 로컬 가상 환경을 생성합니다."
	@echo "  clusterClean   로컬 가상 환경을 삭제합니다."

.DEFAULT_GOAL := help
