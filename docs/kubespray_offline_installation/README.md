# How to deploy Kubernetes cluster with Kubespray
> Kubespray is an open-source project used to deploy a production ready Kubernetes cluster by using Ansible. [Read more](https://kubespray.io/#/)

## Ansible prerequisites
Since the **Kubespray** is using **Ansible** to deploy a Kubernetes cluster, you will
need to setup some requirements to enable Ansible to have permission to access
to each node of the cluster.
* Enable **SSH** access without passphrase
  - Creating authentication key pairs for SSH on all nodes
    ```
    [tmax@c30 ~]$ ssh-keygen
    Generating public/private rsa key pair.
    Enter file in which to save the key (/home/tmax/.ssh/id_rsa):
    Created directory '/home/tmax/.ssh'.
    Enter passphrase (empty for no passphrase):
    Enter same passphrase again:
    Your identification has been saved in /home/tmax/.ssh/id_rsa.
    Your public key has been saved in /home/tmax/.ssh/id_rsa.pub.
    The key fingerprint is:
    SHA256:FUnSTSKhFhz//IjmvI+JIqtH2QIyuYSIpx7VyNRIowI tmax@c30
    The key's randomart image is:
    +---[RSA 2048]----+
    |E .oo.o.++++.    |
    |. .o...+ ooo.    |
    |=oo o o . .      |
    |X..+ o   +       |
    |o*.o    S o      |
    |o.+ .    . o     |
    |....    o . .    |
    | .o .  = o       |
    |.o.o .. *o.      |
    +----[SHA256]-----+
    ```

  - Copying the public key to each node
    ```
    [tmax@c30 ~]$ ssh-copy-id -i ~/.ssh/id_rsa tmax@192.168.0.31
    /usr/bin/ssh-copy-id: INFO: Source of key(s) to be installed: "/home/tmax/.ssh/id_rsa.pub"
    The authenticity of host '192.168.0.31 (192.168.0.31)' can't be established.
    ECDSA key fingerprint is SHA256:BYEB4B9fYUgzYs52jih/eDn3GkibncvdcTM6kUvha+s.
    ECDSA key fingerprint is MD5:69:6f:cf:50:4f:13:1a:91:1a:e1:8f:0d:7a:5f:69:ee.
    Are you sure you want to continue connecting (yes/no)? yes
    /usr/bin/ssh-copy-id: INFO: attempting to log in with the new key(s), to filter out any that are already installed
    /usr/bin/ssh-copy-id: INFO: 1 key(s) remain to be installed -- if you are prompted now it is to install the new keys
    tmax@192.168.0.31's password:

    Number of key(s) added: 1

    Now try logging into the machine, with:   "ssh 'tmax@192.168.0.31'"
    and check to make sure that only the key(s) you wanted were added.
    ```

* Enable **sudo** privileges without password on all nodes
  ```shell
  sudo -i
  echo 'tmax ALL=(ALL) NOPASSWD:ALL' > /etc/sudoers.d/tmax
  # change 'tmax' to your username
  ```

# Online Environment
In **online** environment where all of your cluster nodes can access to the
Internet,deploying a Kubernetes cluster can be done easily with just a few
commands.

* Get Kubespray from Github
  ```shell
  git clone https://github.com/kubernetes-sigs/kubespray.git
  cd kubespray
  ```

* Install dependencies from `requirements.txt`
  ```shell
  sudo pip3 install -r requirements.txt
  ```

* Copy `inventory/sample` as `inventory/mycluster`
  ```shell
  cp -rfp inventory/sample inventory/mycluster
  ```

* Update Ansible inventory file with inventory builder
  ```shell
  declare -a IPS=(192.168.0.31 192.168.0.32 192.168.0.33)
  CONFIG_FILE=inventory/mycluster/hosts.yaml python3 contrib/inventory_builder/inventory.py ${IPS[@]}
  ```

* Review and change parameters under `inventory/mycluster/group_vars` to deploy your desired cluster
  ```shell
  cat inventory/mycluster/group_vars/all/all.yml
  cat inventory/mycluster/group_vars/k8s-cluster/k8s-cluster.yml
  ```

* Deploy a Kubernetes cluster
  ```shell
  ansible-playbook -i inventory/mycluster/hosts.yaml  --become --become-user=root cluster.yml
  ```

# Offline Environment
In **offline** environment where all of your cluster nodes are not able to access to
the Internet, deploying a Kubernetes cluster can be complicated and troublesome.

* To deploy Kubernetes cluster, you will need to:
  1. **Install Kubespray's requirement packages** (Python packages)
  > Packages (requirements.txt):  ansible==2.7.16, jinja2==2.10.1, netaddr==0.7.19,
  pbr==5.2.0, hvac==0.8.2, jmespath==0.9.4, ruamel.yaml==0.15.96

  2. **Install Kubernetes's required packages** (deb for Ubuntu, rpm for CentOS)
  > Packages: docker, python-apt, aufs-tools, apt-transport-https, software-properties-common, ebtables, etc...

  3. **Download binary files and Docker images**
  > Binary: kubeadm, kubelet, kubectl, etcd, etcdctl, calicoctl, etc...

  > Images: kube-proxy, kube-apiserver, kube-controller-manager, kube-scheduler,
  etcd, pause, calico_node, calico_cni, calico_kuber-controllers, k8s-dns-node-cache, coredns, nginx, etc...

  4. **Modify parameter variables**
  > Modify some parameters to install from cache or private server

  5. **Execute Ansible-playbook**
  > Deploy or Reset the Kubernetes cluster

## 1. Install Kubespray's required packages (Python packages)
In order to install Python packages on offline node, you first need to download all Python required packages and dependencies, copy to the node you wish to run *ansible-playbook* to deploy the Kubernetes cluster, and install them manually.
> Assume that all required Python packages and dependencies are already downloaded.
[Read more](pip_install_kubespray_requirements.md)

* Installing required packages in `requirements.txt` from local directory
  ```shell
  sudo pip3 install --no-index --find-links=/path/to/pkg/ -r requirements.txt
  ```
* In case you want to install from private WebServer (example: 192.168.0.200)
  ```shell
  sudo pip3 install --index-url http://192.168.0.200/pip-pkg/ -r requirements.txt
  ```

## 2. Install Kubernetes's required packages
In order to install packages on offline node, you first need to download all required packages and dependencies, copy to all nodes, and then install them manually. However, I personally recommend creating **local repository**.
> Assume that all required packages and dependencies are already downloaded.

* Installing **deb** packages on Ubuntu
  ```shell  
  sudo dpkg -i *.deb
  ```
  If you want to create Ubuntu local repository, [read this](create_ubuntu_repository.md)

* Installing **rpm** packages on CentOS
  ```shell
  sudo rpm -ivh *.rpm
  ```
  If you want to create CentOS local repository, [read this](create_centos_repository.md)

## 3. Download binary files and Docker images
When we execute *ansible-playbook* to deploy the Kubernetes cluster, Kubespray will download binary files (`kubeadm`, `kubelet`, `kubectl`, `etcdctl`, `calicoctl`, ...), also pull container images (`kube-proxy`, `kube-apiserver`, `kube-controller-manager`, `kube-scheduler`, `etcd`, ...) to initialize the cluster. However, this will be failed because you don't have the Internet connection to download them from the official URL address.

To solve this problem, you have two choices:
  1. [**Cache**] Download all required *binary* files, and then store them in `download_cache_dir` location. For *Docker images*, pull and save container images as `tar` file, and then store them in `download_cache_dir`/images location. [Read more](kubespray_offline_with_cache.md)

  2. [**Local Webserver**] Download all required *binary* files and then upload to local WebServer. For *Docker* images, create local Docker registry,
  pull required images from official registry and then push them to local registry.
  [Read more](kubespray_offline_with_private_server.md)


### Binary files
Ensure that you have downloaded the correct **version** of the binary which is declared in `roles/download/defaults/main.yml`.
```yaml
download_cache_dir: /tmp/kubespray_cache

# Arch of Docker images and needed packages
image_arch: "{{host_architecture | default('amd64')}}"

# Versions
kube_version: v1.17.2
etcd_version: v3.3.12
# More...

# Binary File Download URLs
kubelet_download_url: "https://storage.googleapis.com/kubernetes-release/release/{{ kube_version }}/bin/linux/{{ image_arch }}/kubelet"
kubectl_download_url: "https://storage.googleapis.com/kubernetes-release/release/{{ kube_version }}/bin/linux/{{ image_arch }}/kubectl"
kubeadm_download_url: "https://storage.googleapis.com/kubernetes-release/release/{{ kube_version }}/bin/linux/{{ image_arch }}/kubeadm"
etcd_download_url: "https://github.com/coreos/etcd/releases/download/{{ etcd_version }}/etcd-{{ etcd_version }}-linux-{{ image_arch }}.tar.gz"
cni_download_url: "https://github.com/containernetworking/plugins/releases/download/{{ cni_version }}/cni-plugins-linux-{{ image_arch }}-{{ cni_version }}.tgz"
calicoctl_download_url: "https://github.com/projectcalico/calicoctl/releases/download/{{ calico_ctl_version }}/calicoctl-linux-{{ image_arch }}"
crictl_download_url: "https://github.com/kubernetes-sigs/cri-tools/releases/download/{{ crictl_version }}/crictl-{{ crictl_version }}-{{ ansible_system | lower }}-{{ image_arch }}.tar.gz"
```

### Docker images
Ensure that you have downloaded the correct **tag version** of image which is
declared in `roles/download/defaults/main.yml`.
```yaml
# image repo define
gcr_image_repo: "gcr.io"
kube_image_repo: "k8s.gcr.io"
docker_image_repo: "docker.io"
quay_image_repo: "quay.io"

# Container image name and tag
kube_proxy_image_repo: "{{ kube_image_repo }}/kube-proxy"
kube_proxy_image_tag: "{{ kube_version }}"
etcd_image_repo: "{{ quay_image_repo }}/coreos/etcd"
etcd_image_tag: "{{ etcd_version }}{%- if image_arch != 'amd64' -%}-{{ image_arch }}{%- endif -%}"
calico_node_image_repo: "{{ docker_image_repo }}/calico/node"
calico_node_image_tag: "{{ calico_version }}"
coredns_image_repo: "{{ docker_image_repo }}/coredns/coredns"
coredns_image_tag: "1.6.7"
# More...
```
## 4. Modify parameter variables
After you have downloaded binary files and container images, you will need to
change some parameter **variables** to tell Kubespray to install in offline mode.
Those variables can be modified in `roles/download/defaults/main.yml`.

> It is recommended to add or modify parameter variables in `inventory/mycluster/group_vars/k8s-cluster/k8s-cluster.yml` because this will overwrite parameter variables in `roles/download/defaults/main.yml`.

* **In case, you choose to install from cache**, [Read more](kubespray_offline_with_cache.md)
  ```yaml
  # Download cache directory
  download_cache_dir: /tmp/kubespray_cache
  # Run download binary files and container images only once
  download_run_once: true
  # Use the local_host for download_run_once mode
  download_localhost: true
  ```

* **In case, you choose to install from local WebServer**, [Read more](kubespray_offline_with_private_server.md)
  ```yaml
  # Download binary URL
  hcs_url: "http://192.168.0.200/binary"
  kubelet_download_url: "{{ hcs_url }}/storage.googleapis.com/kubernetes-release/release/{{ kube_version }}/bin/linux/{{ image_arch }}/kubelet"
  kubectl_download_url: "{{ hcs_url }}/storage.googleapis.com/kubernetes-release/release/{{ kube_version }}/bin/linux/{{ image_arch }}/kubectl"
  kubeadm_download_url: "{{ hcs_url }}/storage.googleapis.com/kubernetes-release/release/{{ kube_version }}/bin/linux/{{ image_arch }}/kubeadm"
  etcd_download_url: "{{ hcs_url }}/github.com/coreos/etcd/releases/download/{{ etcd_version }}/etcd-{{ etcd_version }}-linux-{{ image_arch }}.tar.gz"
  cni_download_url: "{{ hcs_url }}/github.com/containernetworking/plugins/releases/download/{{ cni_version }}/cni-plugins-linux-{{ image_arch }}-{{ cni_version }}.tgz"
  calicoctl_download_url: "{{ hcs_url }}/github.com/projectcalico/calicoctl/releases/download/{{ calico_ctl_version }}/calicoctl-linux-{{ image_arch }}"
  crictl_download_url: "{{ hcs_url }}/github.com/kubernetes-sigs/cri-tools/releases/download/{{ crictl_version }}/crictl-{{ crictl_version }}-{{ ansible_system | lower }}-{{ image_arch }}.tar.gz"

  # Docker registry
  docker_insecure_registries:
    - 192.168.0.200:5000
  hcs_image_repo: "192.168.0.200:5000"
  docker_image_repo: "{{ hcs_image_repo }}"
  quay_image_repo: "{{ hcs_image_repo }}"
  gcr_image_repo: "{{ hcs_image_repo }}"
  kube_image_repo: "{{ gcr_image_repo }}/google-containers"
  ```

## 5. Execute Ansible-playbook
* To deploy cluster:
  ```shell
  ansible-playbook -i inventory/mycluster/hosts.yaml \
    --become --become-user=root cluster.yml
  ```

* To reset cluster:
  ```shell
  ansible-playbook -i inventory/mycluster/hosts.yaml \
    --become --become-user=root reset.yml
  ```

## Related documents
* [How to deploy Kubernetes cluster in offline environment with cache (Kubespray)](kubespray_offline_with_cache.md)
* [How to deploy Kubernetes cluster in offline environment with local WebServer (kubespray)](kubespray_offline_with_private_server.md)
* [How to create local repository (Ubuntu)](create_ubuntu_repository.md)
* [How to create local repository (CentOS)](create_centos_repository.md)
* [How to create local registry (Docker)](create_docker_registry.md)
* [How to install Python packages in offline environment (pip3 install)](pip_install_kubespray_requirements.md)
