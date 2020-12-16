#!/bin/bash

declare -A PL_VMs
PL_VMs=( ["pl-node1"]="172.22.4.221" \
         ["pl-node2"]="172.22.4.222" \
         ["pl-node3"]="172.22.4.223" )
VM="VirtualBox"
DEPLOYED_KUBE=( "prolinux_v1.15.3_calico" \
                "prolinux_v1.17.9_calico" \
                "prolinux_v1.15.3_flannel" \
                "prolinux_v1.17.9_flannel" \
	        "prolinux_v1.19.4_calico" )

# Check if VirtualBox vm is running
is_vm_running(){
    VBoxManage list runningvms | grep "$1" &> /dev/null
}

# Check if VirtualBox vm exists
is_vm_exist(){
    VBoxManage list vms | grep "$1" &> /dev/null
}

# VirtualBox: Start vm
virtualbox_startvm(){
    echo "[$VM] start vm: '$1'..."
    if is_vm_exist "$1" && ! is_vm_running "$1"; then
        if ! VBoxManage startvm "$1" --type headless; then
            echo "[$VM] [ERROR] Unable to start '$1'!"
            exit 1
        fi
    else
        echo "[$VM] [ERROR] '$1' is NOT exist!"
        exit 1
    fi
}

# VirtualBox: stop vm
virtualbox_stopvm(){
    echo "[$VM] Poweroff vm: '$1'..."
    if is_vm_running "$1"; then
        if ! VBoxManage controlvm "$1" poweroff; then
            echo "[$VM] [ERROR] Unable to stop '$1'!"
            exit 1
        fi
    else
        echo "[$VM] $1 is not running!"
    fi
}

# VirtualBox: Check if snapshot exists
virtualbox_check_snapshot(){
    VBoxManage snapshot "$1" list | grep "$2" &> /dev/null
}

# VirtualBOx: restore snapshot
virtualbox_restore_snapshot(){
    echo "[$VM] '$1' restore snapshot: '$2'"
    if virtualbox_check_snapshot "$@"; then
        if is_vm_running "$1"; then
            virtualbox_stopvm "$1"
        fi

        if VBoxManage snapshot "$1" restore "$2"; then
            echo "[$VM] '$1' has been restored to: $2"
        else
            echo "[$VM] '$1' Unable to restore snapshot: $2"
            exit 1
        fi
    else
        echo "[$VM] [ERROR] '$1' snapshot ($2) is not exist!"
        exit 1
    fi
}

# Check if Kubernetes cluster is alread deployed
is_deployed_kube(){
    for kube in "${DEPLOYED_KUBE[@]}"; do
        if [ "$1" == "$kube" ]; then
            return 0
        fi
    done
    return 1
}

# Check if SSH to vm nodes is success
is_ssh_success(){
    local tries=1
    while ((tries < 30)); do
        echo "[SSH] Checking SSH connection to $1 : try ($tries)"
        if ssh -q root@"$1" exit; then
            echo "[SSH] SSH connection to $1 : SUCCESS"
            return 0
        fi
        tries=$((tries + 1))
        sleep 1
    done
    echo "[SSH] SSH connection to $1 : FAIL"
    return 1
}

is_all_pod_running(){
    local tries=1
    while ((tries < 30)); do
		echo "[Kubernetes] Waiting all pod to be Running..."
		if kubectl get node &> /dev/null; then
			if ! kubectl get pod -A | grep -E 'Error|CrashLoopBackOff' &> /dev/null; then
				kubectl cluster-info
				kubectl get node,pod -A -o wide
				return 0
			fi
		fi
        tries=$((tries + 1))
        sleep 3
    done
    return 1
}

# Copy Kubernetes admin.conf from master node to local node
get_kube_config(){
    KUBE_DIR=$1
    [ ! -d "$KUBE_DIR" ] && mkdir -p "$KUBE_DIR"
    for ip in "${PL_VMs[@]}"; do
        if scp root@"${ip}":/etc/kubernetes/admin.conf "$KUBE_DIR"/config &> /dev/null; then
            sudo chown "$(id -u)":"$(id -g)" "$KUBE_DIR"/config &> /dev/null
            return 0
        fi
    done
    return 1
}

# Restart ntp service : solve ceph HEALTH_WARN (clock skew)
restart_ntp(){
    for ip in "${PL_VMs[@]}"; do
		ssh root@"${ip}" 'systemctl restart ntpdate' &> /dev/null
		ssh root@"${ip}" 'timedatectl set-ntp true' &> /dev/null
    done
}

# Loading requested cluster from snapshot
loading_kube(){
    for vm in "${!PL_VMs[@]}"
    do
        virtualbox_restore_snapshot "$vm" "$1"
        echo
        virtualbox_startvm "$vm"
        echo
    done

    # Check SSH connection
    for ip in "${PL_VMs[@]}"; do
        if ! is_ssh_success "$ip"; then
            exit 1
        fi
    done
}

# Deploy Kubernetes cluster with kubespray
deploying_kube(){
    echo "====================Kubespray Deploy Cluster (Start)==================="
    echo "[DEPLOY] KUBE_VERSION=$1, KUBE_NETWORK=$2"
    echo "[DEPLOY] IPS=(${PL_VMs[*]})"

    e2eDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    kubesprayDir="${e2eDir}/kubespray"

    # remove old ansible fact
    sudo rm -rf /tmp/node*

    set -eo pipefail
    pushd "$kubesprayDir"
    sudo pip3 install -r requirements.txt
    cp -rfp inventory/sample inventory/mycluster
    CONFIG_FILE=inventory/mycluster/hosts.yaml python3 contrib/inventory_builder/inventory.py "${PL_VMs[@]}"
    ansible-playbook -i inventory/mycluster/hosts.yaml --become --become-user=root cluster.yml -e \
        "ansible_user=root kube_version_min_required=v1.15.0 kube_version=$1 kube_network_plugin=$2" -e ansible_os_family=RedHat -e ansible_distribution=CentOS
    popd
    set +eo pipefail
    echo "====================Kubespray Deploy Cluster (Finish)==================="
    echo
}

# Loading requestd cluster if exists, or deploy if not
prolinuxKubeUp(){
    echo "[ClusterUP] Deploy Kubernetes cluster (Start)"
    kubeVer=$1
    kubeNet=$2

    if [ "$BOX_OS" != "prolinux" ] ; then
        echo "[ERROR] Need env BOX_OS=prolinux"
        exit 1
    fi

    if [ "$kubeVer" == "" ] ; then
        echo "[ERROR] Need env KUBE_VERSION=v1.xx.x"
        exit 1
    fi

    if [ "$kubeNet" == "" ] ; then
        echo "[INFO] Default network plugin (calico) is used!"
        kubeNet=calico
    fi

    requestKube="${BOX_OS}_${kubeVer}_${kubeNet}"
    if is_deployed_kube "$requestKube"; then
        echo "##### ($requestKube) is already deployed. #####"
        echo
        loading_kube "$requestKube"
    else
        echo "***** ($requestKube) is being deployed. *****"
        echo
        prolinuxKubeClean
        deploying_kube "$kubeVer" "$kubeNet"
    fi

	restart_ntp

    if get_kube_config "$HOME/.kube"; then
        export KUBECONFIG=$KUBE_DIR/config
		if ! is_all_pod_running; then
			echo "[ClusterUP] Cluster is not healthy."
			exit 1
		fi

    else
        echo "[ClusterUp] Unable to kube config 'admin.conf' from master node!"
        exit 1
    fi

    echo "[ClusterUP] Deploy Kubernetes cluster (Finish)"
    echo
}

# Clean kubernetes cluster: restore vm snapshot to No_k8s stage
prolinuxKubeClean(){
    loading_kube "prolinux_no_k8s"
}

# Poweroff vm nodes
prolinuxKubeDown(){
    for vm in "${!PL_VMs[@]}"
    do
        virtualbox_stopvm "$vm"
    done
}

# main
case "$1" in
    up)
        #export BOX_OS=prolinux KUBE_VERSION=v1.15.3 KUBE_NETWORK=calico
        prolinuxKubeUp "$KUBE_VERSION" "$KUBE_NETWORK"
        ;;
    clean)
        prolinuxKubeClean
        ;;
    down)
        prolinuxKubeDown
        ;;
    *)
        echo "Usage: $0 COMMAND"
        echo "COMMAND:"
        echo "  up      : Deploy ProLinux kubernetes cluster with env KUBE_VERSION, KUBE_NETWORK"
        echo "  clean   : Clean ProLinux kubernetes cluster from nodes (no_k8s stage)"
        echo "  down    : Poweroff ProLinux kubernetes cluster nodes"
        echo 
        echo "Prerequisite:"
        echo "export BOX_OS=prolinux KUBE_VERSION=v1.16.8 KUBE_NETWORK=calico"
        echo
        ;;
esac

