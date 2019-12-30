CLUSTER=./hack/cluster.sh

cluster-up:
	${CLUSTER} up

cluster-status:
	${CLUSTER} status

cluster-clean:
	${CLUSTER} clean