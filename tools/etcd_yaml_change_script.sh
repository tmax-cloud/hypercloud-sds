#!/bin/bash

if  [ "$#" -ne 2 ]; then
	echo "usage : ./etcd_yaml_change_script.sh <registry_endpoint> <master_node>"
	echo "example : ./etcd_yaml_change_script.sh 192.168.50.90:5000 master1"
	exit 0
fi

registry_endpoint=$1
master_node=$2

sed -i 's/{registry_endpoint}/'$registry_endpoint'/g' etcd_snapshot.yaml
sed -i 's/{master_node}/'$master_node'/g' etcd_snapshot.yaml
