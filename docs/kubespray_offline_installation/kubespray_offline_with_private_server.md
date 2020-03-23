# How to run kubespray in offline environment with private server
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
* Edit environment variables in `k8s-cluster.yml`
```yml
# Default cluster component
container_manager: docker
kube_network_plugin: calico
# Version
kube_version: v1.17.2
etcd_version: v3.3.12
cni_version: "v0.8.3"
calico_version: "v3.11.1"
calico_ctl_version: "v3.11.1"
calico_cni_version: "v3.11.1"
calico_policy_version: "v3.11.1"
nodelocaldns_version: "1.15.8"
coredns_version: "1.6.0"
dnsautoscaler_version: 1.6.0
pod_infra_version: 3.1
nginx_image_tag: 1.17
dashboard_image_tag: "v1.10.1"
# Download binary
hcs_url: "http://192.168.0.200/binary"
kubelet_download_url: "{{ hcs_url }}/storage.googleapis.com/kubernetes-release/release/{{ kube_version }}/bin/linux/{{ image_arch }}/kubelet"
kubectl_download_url: "{{ hcs_url }}/storage.googleapis.com/kubernetes-release/release/{{ kube_version }}/bin/linux/{{ image_arch }}/kubectl"
kubeadm_download_url: "{{ hcs_url }}/storage.googleapis.com/kubernetes-release/release/{{ kube_version }}/bin/linux/{{ image_arch }}/kubeadm"
etcd_download_url: "{{ hcs_url }}/github.com/coreos/etcd/releases/download/{{ etcd_version }}/etcd-{{ etcd_version }}-linux-{{ image_arch }}.tar.gz"
cni_download_url: "{{ hcs_url }}/github.com/containernetworking/plugins/releases/download/{{ cni_version }}/cni-plugins-linux-{{ image_arch }}-{{ cni_version }}.tgz"
calicoctl_download_url: "{{ hcs_url }}/github.com/projectcalico/calicoctl/releases/download/{{ calico_ctl_version }}/calicoctl-linux-{{ image_arch }}"
crictl_download_url: "{{ hcs_url }}/github.com/kubernetes-sigs/cri-tools/releases/download/{{ crictl_version }}/crictl-{{ crictl_version }}-{{ ansible_system | lower }}-{{ image_arch }}.tar.gz"
# Checksum [These values will be update automatically when execute download_binary.sh
etcd_binary_checksum: dc5d82df095dae0a2970e4d870b6929590689dd707ae3d33e7b86da0f7f211b6
cni_binary_checksum: 29a092bef9cb6f26c8d5340f3d56567b62c7ebdb1321245d94b1842c80ba20ba
kubelet_binary_checksum: 680d6afa09cd51061937ebb33fd5c9f3ff6892791de97b028b1e7d6b16383990
kubectl_binary_checksum: 4475f68c51af23925d7bd7fc3d1bd01bedd3d4ccbb64503517d586e31d6f607c
kubeadm_binary_checksum: 366a7f260cbd1aaa2661b1e3b83a7fc8781c8a8b07c71944bdaf66d49ff5abae
calicoctl_binary_checksum: 045fdbfdb30789194c499ba17c8eac6d1704fe20d05e3c10027eb570767386db
crictl_binary_checksum: c3b71be1f363e16078b51334967348aab4f72f46ef64a61fe7754e029779d45a
# Docker registry
docker_insecure_registries:
  - 192.168.0.200:5000
hcs_image_repo: "192.168.0.200:5000"
docker_image_repo: "{{ hcs_image_repo }}"
quay_image_repo: "{{ hcs_image_repo }}"
gcr_image_repo: "{{ hcs_image_repo }}"
kube_image_repo: "{{ gcr_image_repo }}/google-containers"
```

## Add required packages into our local repository
> suppose that `192.168.0.200` is server IP.

* Create Ubuntu repository  
see [here](create_ubuntu_repository.md)  
* Check required packages
> ~/kubespray/roles/kubernetes/preinstall/vars/ubuntu.yml  
 required_pks: python-apt, aufs-tools, apt-transport-https, software-properties-common, ebtables

* Download required packages and dependencies
```shell
sudo apt-get install --download-only <package_name>
```
All downloaded deb files will be saved in `/var/cache/apt/archives` directory
* [optional] use `apt-rdepends` to get all packages
```shell
sudo apt install apt-rdepends
sudo apt download $(apt-rdepends vim | grep -v "^ ")
```
* Copy downloaded deb packages into our repository directory
```shell
cp *.deb /var/www/html/hcs-repo/bionic/
```
* Update `Release` and index files
```shell
cd /var/www/html/hcs-repo/bionic
apt-ftparchive packages . > Packages
gzip -c Packages > Packages.gz
apt-ftparchive release . > Release
gpg --yes --clearsign -o InRelease Release
gpg --yes -abs -o Release.gpg Release
```

## Add required binary files into our local web Server
> suppose that `192.168.0.200` is server IP.

* Create local web server using apache2.  
ref: [create_ubuntu_repository](create_ubuntu_repository.md)
* Check required binary files
> ~/kubespray/roles/download/defaults/main.yml  
required_file: kubeadm, kubelet, kubectl, etcd, cni, calicoctl, crictl

* Download required binary files from the original url put into our local web server.
```shell
sudo ./download_binary.sh ~/kubespray/inventory/mycluster/group_vars/k8s-cluster/k8s-cluster.yml
```
script: [download_binary.sh](download_binary.sh)  
This will download all required binary files with defined version in `k8s-cluster.yml` and store them into our local web server located in `/var/www/html/binary/`.

## Add required docker images into our local Docker registry
> suppose that `192.168.0.200:5000` is docker registry url.

* Create local Docker registry  
see [here](create_docker_registry.md)

* Pull and Push required docker images into local registry
```shell
sudo ./push_docker_image.sh 192.168.0.200:5000 \
    ~/kubespray/inventory/mycluster/group_vars/k8s-cluster/k8s-cluster.yml
```
script: [push_docker_image.sh](push_docker_image.sh)  

## Run kubespray's ansible-playbook to deploy Kubernetes mycluster
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
