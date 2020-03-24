# How to install Kubespray's required packages via pip3 offline
> pip is the package installer for Python. You can us pip to install packages
from the `Python Package Index` and other indexes. [read more](pip_install_guide.md)

## Kubespray requirement packages
* List of required packages in `requirements.txt`
  ```shell
  $ cat ./kubespray/requirements.txt
  ansible==2.7.16
  jinja2==2.10.1
  netaddr==0.7.19
  pbr==5.2.0
  hvac==0.8.2
  jmespath==0.9.4
  ruamel.yaml==0.15.96
  ```
* To deploy the cluster with Ansible, you will need to install dependencies from
`requirements.txt`
  ```shell
  sudo pip3 install -r requirements.txt
  ```
  - Make sure that `pip3` is installed, and you can run from command line.
    ```shell
    pip3 --version
    ```
  - However required packages won't be successfully installed if you don't have Internet connection.

## Installing pip3
* With Internet connection (Online), you can easily install it by:
  ```shell
  sudo apt-get update
  sudo apt-get install -y python3-pip
  ```
* Without Internet connection (Offline), you have to install it manually.
On Ubuntu, you can download deb package file and its dependencies, then install
by using `dpkg -i *.deb` or create you own [private repository](create_ubuntu_repository.md).
  - Based on Ubuntu 18.04.3 LTS, these deb packages need to be installed.
    ```
    binutils_2.30-21ubuntu1~18.04.2_amd64.deb
    binutils-common_2.30-21ubuntu1~18.04.2_amd64.deb
    binutils-x86-64-linux-gnu_2.30-21ubuntu1~18.04.2_amd64.deb
    build-essential_12.4ubuntu1_amd64.deb
    cpp_4%3a7.4.0-1ubuntu2.3_amd64.deb
    cpp-7_7.5.0-3ubuntu1~18.04_amd64.deb
    dh-python_3.20180325ubuntu2_all.deb
    dpkg-dev_1.19.0.5ubuntu2.3_all.deb
    fakeroot_1.22-2ubuntu1_amd64.deb
    g++_4%3a7.4.0-1ubuntu2.3_amd64.deb
    g++-7_7.5.0-3ubuntu1~18.04_amd64.deb
    gcc_4%3a7.4.0-1ubuntu2.3_amd64.deb
    gcc-7_7.5.0-3ubuntu1~18.04_amd64.deb
    gcc-7-base_7.5.0-3ubuntu1~18.04_amd64.deb
    libalgorithm-diff-perl_1.19.03-1_all.deb
    libalgorithm-diff-xs-perl_0.04-5_amd64.deb
    libalgorithm-merge-perl_0.08-3_all.deb
    libasan4_7.5.0-3ubuntu1~18.04_amd64.deb
    libatomic1_8.3.0-26ubuntu1~18.04_amd64.deb
    libbinutils_2.30-21ubuntu1~18.04.2_amd64.deb
    libc6-dev_2.27-3ubuntu1_amd64.deb
    libcc1-0_8.3.0-26ubuntu1~18.04_amd64.deb
    libc-dev-bin_2.27-3ubuntu1_amd64.deb
    libcilkrts5_7.5.0-3ubuntu1~18.04_amd64.deb
    libdpkg-perl_1.19.0.5ubuntu2.3_all.deb
    libexpat1-dev_2.2.5-3ubuntu0.2_amd64.deb
    libfakeroot_1.22-2ubuntu1_amd64.deb
    libfile-fcntllock-perl_0.22-3build2_amd64.deb
    libgcc-7-dev_7.5.0-3ubuntu1~18.04_amd64.deb
    libgomp1_8.3.0-26ubuntu1~18.04_amd64.deb
    libisl19_0.19-1_amd64.deb
    libitm1_8.3.0-26ubuntu1~18.04_amd64.deb
    liblsan0_8.3.0-26ubuntu1~18.04_amd64.deb
    libmpc3_1.1.0-1_amd64.deb
    libmpx2_8.3.0-26ubuntu1~18.04_amd64.deb
    libpython3.6_3.6.9-1~18.04_amd64.deb
    libpython3.6-dev_3.6.9-1~18.04_amd64.deb
    libpython3.6-minimal_3.6.9-1~18.04_amd64.deb
    libpython3.6-stdlib_3.6.9-1~18.04_amd64.deb
    libpython3-dev_3.6.7-1~18.04_amd64.deb
    libquadmath0_8.3.0-26ubuntu1~18.04_amd64.deb
    libstdc++-7-dev_7.5.0-3ubuntu1~18.04_amd64.deb
    libtsan0_8.3.0-26ubuntu1~18.04_amd64.deb
    libubsan0_7.5.0-3ubuntu1~18.04_amd64.deb
    linux-libc-dev_4.15.0-91.92_amd64.deb
    make_4.1-9.1ubuntu1_amd64.deb
    manpages-dev_4.15-1_all.deb
    python3.6_3.6.9-1~18.04_amd64.deb
    python3.6-dev_3.6.9-1~18.04_amd64.deb
    python3.6-minimal_3.6.9-1~18.04_amd64.deb
    python3-crypto_2.6.1-8ubuntu2_amd64.deb
    python3-dev_3.6.7-1~18.04_amd64.deb
    python3-distutils_3.6.9-1~18.04_all.deb
    python3-keyring_10.6.0-1_all.deb
    python3-keyrings.alt_3.0-1_all.deb
    python3-lib2to3_3.6.9-1~18.04_all.deb
    python3-pip_9.0.1-2.3~ubuntu1.18.04.1_all.deb
    python3-secretstorage_2.3.1-2_all.deb
    python3-setuptools_39.0.1-2_all.deb
    python3-wheel_0.30.0-0.2_all.deb
    python3-xdg_0.25-4ubuntu1_all.deb
    python-pip-whl_9.0.1-2.3~ubuntu1.18.04.1_all.deb
    ```

## Downloading Kubespray's required packages
* On node with Internet connection, download required packages from `requirements.txt`
  ```shell
  sudo pip3 download -r requirements.txt
  ```
  Then, copy downloaded files to offline node for installation.
* List of downloaded Packages
  ```shell
  $ tree /tmp/pip_download_pkg_dir/
  .
  ├── ansible-2.7.16.tar.gz
  ├── bcrypt-3.1.7-cp34-abi3-manylinux1_x86_64.whl
  ├── certifi-2019.11.28-py2.py3-none-any.whl
  ├── cffi-1.14.0-cp36-cp36m-manylinux1_x86_64.whl
  ├── chardet-3.0.4-py2.py3-none-any.whl
  ├── configparser-4.0.2-py2.py3-none-any.whl
  ├── cryptography-2.8-cp34-abi3-manylinux1_x86_64.whl
  ├── hvac-0.8.2-py2.py3-none-any.whl
  ├── idna-2.9-py2.py3-none-any.whl
  ├── ipaddress-1.0.23-py2.py3-none-any.whl
  ├── Jinja2-2.10.1-py2.py3-none-any.whl
  ├── jmespath-0.9.4-py2.py3-none-any.whl
  ├── MarkupSafe-1.1.1-cp36-cp36m-manylinux1_x86_64.whl
  ├── netaddr-0.7.19-py2.py3-none-any.whl
  ├── paramiko-2.7.1-py2.py3-none-any.whl
  ├── pbr-5.2.0-py2.py3-none-any.whl
  ├── pycparser-2.20-py2.py3-none-any.whl
  ├── PyNaCl-1.3.0-cp34-abi3-manylinux1_x86_64.whl
  ├── PyYAML-5.3.1.tar.gz
  ├── requests-2.23.0-py2.py3-none-any.whl
  ├── ruamel.yaml-0.15.96-cp36-cp36m-manylinux1_x86_64.whl
  ├── ruamel.yaml-0.16.10-py2.py3-none-any.whl
  ├── ruamel.yaml.clib-0.2.0-cp36-cp36m-manylinux1_x86_64.whl
  ├── setuptools-46.1.1-py3-none-any.whl
  ├── six-1.14.0-py2.py3-none-any.whl
  └── urllib3-1.25.8-py2.py3-none-any.whl

  0 directories, 26 files
  ```

## Installing Kubespray's required packages
* To install packages from local directory:
  ```shell
  sudo pip3 install --no-index --find-links=/tmp/pip_download_pkg_dir/ -r requirement.txt
  ```
* To install packages from private index (repository):
  ```shell
  sudo pip3 install --index-url http:/192.168.0.200/hcs-repo/pip/ -r requirements.txt
  ```
