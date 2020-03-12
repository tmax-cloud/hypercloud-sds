#!/bin/bash

TIMEOUT=300
OBS_NAME="my-store"

function kubectl_objstore_yaml(){
cat <<YAML | kubectl "$1" -f -
apiVersion: ceph.rook.io/v1
kind: CephObjectStore
metadata:
  name: "$OBS_NAME"
  namespace: rook-ceph
spec:
  metadataPool:
    replicated:
      size: 1
  dataPool:
    replicated:
      size: 1
  gateway:
    type: s3
    port: 80
    securePort:
    instances: 1
YAML
}

function kubectl_sc_bucket_yaml(){
cat <<YAML | kubectl "$1" -f -
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
   name: rook-ceph-delete-bucket
provisioner: ceph.rook.io/bucket
reclaimPolicy: Delete
parameters:
  objectStoreName: "$OBS_NAME"
  objectStoreNamespace: rook-ceph
  region: us-east-1
YAML
}

function kubectl_obc_yaml(){
cat <<YAML | kubectl "$1" -f -
apiVersion: objectbucket.io/v1alpha1
kind: ObjectBucketClaim
metadata:
  name: ceph-delete-bucket
spec:
  generateBucketName: ceph-bkt
  storageClassName: rook-ceph-delete-bucket
YAML
}

case "$1" in
	create|c)
		if kubectl -n rook-ceph get cephobjectstore -oyaml | grep 'items: \[\]' &>/dev/null; then
			kubectl_objstore_yaml apply
		else
			OBS_NAME=$(kubectl -n rook-ceph get cephobjectstore -oname | head -1 | cut -d '/' -f2)
		fi

		for ((retry = 0; retry <= TIMEOUT; retry = retry + 5)); do
			if kubectl -n rook-ceph get pod -l app=rook-ceph-rgw | grep Running &>/dev/null; then
				break;
			fi
			echo "Waiting for cephObjectStore '$OBS_NAME' to be created... ${retry}s"
			sleep 5
		done

		if [ "$retry" -gt "$TIMEOUT" ]; then
			echo "[ERROR] Failed to create cephObjectStore: rook-ceph-rgw-my-store is not running."
			exit
		fi

		if ! kubectl get sc rook-ceph-delete-bucket &>/dev/null; then
			kubectl_sc_bucket_yaml apply
		fi
		if ! kubectl get obc ceph-delete-bucket &>/dev/null; then
			kubectl_obc_yaml apply
		fi
		;;

	delete|d)
		if kubectl get sc rook-ceph-delete-bucket &>/dev/null; then
			kubectl_sc_bucket_yaml delete
		fi
		if kubectl get obc ceph-delete-bucket &>/dev/null; then
			kubectl_obc_yaml delete
		fi
		exit
		;;
	*)
		echo "[Create]: $0 create|c"
		echo "[Delete]: $0 delete|d"
		exit
		;;
esac

# Test case1: rgw service type as ClusterIP
OBC="ceph-delete-bucket"
until (kubectl -n default get cm | grep "$OBC" &>/dev/null);
do
	echo "Waiting configmap '$OBC' to be created..."
	sleep 1
done

AWS_HOST=$(kubectl -n default get cm "$OBC" -oyaml | grep BUCKET_HOST | awk '{print $2}')
echo "AWS_HOST: $AWS_HOST"

AWS_BUCKET=$(kubectl -n default get cm "$OBC" -oyaml | grep BUCKET_NAME | awk '{print $2}')
echo "AWS_BUCKET: $AWS_BUCKET"

AWS_ACCESS_KEY=$(kubectl -n default get secrets "$OBC" -oyaml | grep AWS_ACCESS_KEY_ID | awk '{print $2}' | base64 --decode)
echo "AWS_ACCESS_KEY: $AWS_ACCESS_KEY"

AWS_SECRET_KEY=$(kubectl -n default get secrets "$OBC" -oyaml | grep AWS_SECRET_ACCESS_KEY | awk '{print $2}' | base64 --decode)
echo "AWS_SECRET_KEY: $AWS_SECRET_KEY"

IP=$(kubectl -n rook-ceph get svc -l app=rook-ceph-rgw -oyaml | grep clusterIP | awk '{print $2}')
PORT=$(kubectl -n rook-ceph get svc -l app=rook-ceph-rgw -oyaml | grep port: | awk '{print $2}')
echo "AWS_ENDPOINT: $IP:$PORT"

TOOL_BOX_POD=$(kubectl -n rook-ceph get pod -l app=rook-ceph-tools -oname)

echo -e "\n[Info] Install 's3cmd' tool into '$TOOL_BOX_POD' ..."
if kubectl -n rook-ceph exec "$TOOL_BOX_POD" -it -- yum --assumeyes install s3cmd &>/tmp/install_s3cmd.log
then
	echo "[OK] 's3cmd' tool is installed"
else
	echo "[ERROR] Failed to install 's3cmd' tool. view log: /tmp/install_s3cmd.log"
   	exit
fi

echo -e "\n[Info] List bucket:"
if kubectl -n rook-ceph exec "$TOOL_BOX_POD" -it -- s3cmd ls --no-ssl --host="$AWS_HOST" --access_key="$AWS_ACCESS_KEY" --secret_key="$AWS_SECRET_KEY" --host-bucket= &>/tmp/list_bucket.log
then
	cat /tmp/list_bucket.log
else
	echo "[ERROR] Unable to list bucket. view log: /tmp/list_bucket.log"
	exit
fi

echo -e "\n[Info] Create test file '/tmp/rookObj' in $TOOL_BOX_POD"
if kubectl -n rook-ceph exec "$TOOL_BOX_POD" -it -- touch /tmp/rookObj
then
	echo "[OK] check created test file:"
	kubectl -n rook-ceph exec "$TOOL_BOX_POD" -it -- ls -l /tmp/rookObj
else
	echo "[ERROR] Unable to create test file '/tmp/rookObj'"
	exit
fi

echo -e "\n[Info] Upload test file into cephObjectStore"
if kubectl -n rook-ceph exec "$TOOL_BOX_POD" -it -- s3cmd put /tmp/rookObj --no-ssl --host="$AWS_HOST" --access_key="$AWS_ACCESS_KEY" --secret_key="$AWS_SECRET_KEY" --host-bucket= s3://"$AWS_BUCKET"
then
	echo "[OK] Upload success!"
else
	echo "[ERROR] Upload fail!"
	exit
fi

echo -e "\n[Info] Download test file from cephObjectStore"
if kubectl -n rook-ceph exec "$TOOL_BOX_POD" -it -- s3cmd get s3://"$AWS_BUCKET"/rookObj /tmp/rookObj-dl_"$(date +%s)" --no-ssl --host="$AWS_HOST" --access_key="$AWS_ACCESS_KEY" --secret_key="$AWS_SECRET_KEY" --host-bucket=
then
	echo "[OK] Download success!"
else
	echo "[ERROR] Download fail!"
	exit
fi

echo -e "\n[Info] Check downloaded test file (rookObj-dl)"
if ! kubectl -n rook-ceph exec "$TOOL_BOX_POD" -it -- /bin/sh -c "ls -l /tmp/rookObj-dl*"
then
	echo "[ERROR] No downloaded file!"
	exit
fi
