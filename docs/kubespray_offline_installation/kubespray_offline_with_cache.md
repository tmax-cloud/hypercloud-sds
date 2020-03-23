# How to run kubespray in offline environment (cache)
> tested on ubuntu-server 18.04.3 LTS (bionic)

## Get kubespray
* clone kubespray github
```shell
git clone https://github.com/kubernetes-sigs/kubespray.git
cd kubespray
```
* Install dependencies from `requirements.txt`
```shell
sudo pip install -r requirements.txt
```
* Copy `inventory/sample` as `inventory/mycluster`
```shell
cp -rfp inventory/sample inventory/mycluster
```
* Update Ansible inventory file with inventory builder
```shell
declare -a IPS=(192.168.0.201 192.168.0.202 192.168.0.203)
CONFIG_FILE=inventory/mycluster/hosts.yaml python3 contrib/inventory_builder/inventory.py ${IPS[@]}
```
* Review and change parameters under `inventory/mycluster/group_vars`
```shell
cat inventory/mycluster/group_vars/all/all.yml
cat inventory/mycluster/group_vars/k8s-cluster/k8s-cluster.yml
```
* Modify or Add environment variables in `k8s-cluster.yml`
```yml
# Download cache directory
download_cache_dir: /tmp/kubespray_cache
# Run download binary files and container images only once
download_run_once: true
# Use the local_host for download_run_once mode
download_localhost: true
```

* Make sure that `download_cache_dir` (/tmp/kubespray_cache) contains all required binary files and images
```shell
$ tree /tmp/kubespray_cache
.
├── calicoctl
├── cni-plugins-linux-amd64-v0.8.5.tgz
├── images
│   ├── docker.io_calico_cni_v3.11.1.tar
│   ├── docker.io_calico_kube-controllers_v3.11.1.tar
│   ├── docker.io_calico_node_v3.11.1.tar
│   ├── docker.io_coredns_coredns_1.6.7.tar
│   ├── docker.io_library_nginx_1.17.tar
│   ├── gcr.io_google_containers_kubernetes-dashboard-amd64_v1.10.1.tar
│   ├── k8s.gcr.io_cluster-proportional-autoscaler-amd64_1.6.0.tar
│   ├── k8s.gcr.io_k8s-dns-node-cache_1.15.8.tar
│   ├── k8s.gcr.io_kube-apiserver_v1.17.2.tar
│   ├── k8s.gcr.io_kube-controller-manager_v1.17.2.tar
│   ├── k8s.gcr.io_kube-proxy_v1.17.2.tar
│   ├── k8s.gcr.io_kube-scheduler_v1.17.2.tar
│   ├── k8s.gcr.io_pause_3.1.tar
│   └── quay.io_coreos_etcd_v3.3.12.tar
├── kubeadm-v1.17.2-amd64
├── kubectl-v1.17.2-amd64
└── kubelet-v1.17.2-amd64

1 directory, 19 files
```

* Make sure required packages are already installed, or create local repository  
see [create_ubuntu_repository](create_ubuntu_repository)
> ~/kubespray/roles/kubernetes/preinstall/vars/ubuntu.yml  
 required_pkgs: docker, python-minimal, python-apt, aufs-tools, apt-transport-https, software-properties-common, ebtables, etc...  

* Edit Ubuntu docker repository url and GPGkey  
File: `kubespray/roles/container-engine/docker/vars/ubuntu-amd64.yml`
 ```shell
 docker_repo_key_info:
  pkg_key: apt_key
  url: 'http://192.168.0.200/hcs-repo/GPGkey'
  repo_keys:
    - 65B6A53C5C3D875175F3F59F906CC4D42DF1595D

docker_repo_info:
  pkg_repo: apt_repository
  repos:
    - >
       deb http://192.168.0.200/hcs-repo/deb /
 ```


## Run ansible-playbook to create/reset Kubernetes cluster
* To create cluster:
```shell
ansible-playbook -i inventory/mycluster/hosts.yaml \
    --become --become-user=root cluster.yml
```
* To reset cluster:
```shell
ansible-playbook -i inventory/mycluster/hosts.yaml \
    --become --become-user=root reset.yml
```
