# How to create local Yum Repository
> Tested on Centos 7

## Requirements
  1. Access to a user account with **root** or **sudo** privileges
  2. Packages:
     * yum: Yellowdog Updater Modified (installed by default)
     * yum-utils: Utilities based around the yum package manager
     * httpd: web server (apache)
     * createrepo: A tool used to create yum repository
  3. rpm package files ({package_name}.rpm)

## Installing the required packages
  ```shell
  sudo yum install yum-utils
  sudo yum install httpd
  sudo yum install createrepo
  ```

## Create the repository directory Structure
* Create a directory for an HTTP repository using:
  ```shell
  sudo mkdir -p /var/www/html/hcs-yum/packages
  ```

* Move your rpm files (.rpm) to repository package directory
  ```shell
  sudo mv /path/to/my-packages/*.rpm /var/www/html/hcs-yum/packages
  ```

## [Optional] Synchronize HTTP repositories (mirror)
> You can download a local copy of the original official CentOS
repositories to your server by using `reposync` command.

* Create directory to store repositories
  ```shell
  sudo mkdir –p /var/www/html/hcs-yum/{base,centosplus,extras,updates}
  ```

* To download the official CentOS **base** repository:
  ```shell
  sudo reposync -g -l -d -m --repoid=base --newest-only --download-metadata --download_path=/var/www/html/hcs-yum/
  ```

* To download the official CentOS **centosplus** repository:
  ```shell
  sudo reposync -g -l -d -m --repoid=centosplus --newest-only --download-metadata --download_path=/var/www/html/hcs-yum/
  ```

* To download the official CentOS **extras** repository:
  ```shell
  sudo reposync -g -l -d -m --repoid=extras --newest-only --download-metadata --download_path=/var/www/html/hcs-yum/
  ```

* To download the official CentOS **updates** repository:
  ```shell
  sudo reposync -g -l -d -m --repoid=updates --newest-only --download-metadata --download_path=/var/www/html/hcs-yum/
  ```

* In the previous commands, the options are as follows:
  ```shell
  -g – lets you remove or uninstall packages on CentOS that fail a GPG check
  -l – yum plugin support
  -d – lets you delete local packages that no longer exist in the repository
  -m – lets you download comps.xml files, useful for bundling groups of packages by function
  --repoid – specify repository ID
  --newest-only – only download the latest package version, helps manage the size of the repository
  --download-metadata – download non-default metadata
  --download-path – specifies the location to save the packages
  ```

## Create the Repository
* Using **createrepo** utility to create a repository
  ```shell
  sudo createrepo /var/www/html/hcs-yum
  ```

## Make the Apache HTTP Server accessible
* Start HTTP service
  ```shell
  sudo systemctl restart httpd
  ```

* Enable HTTP serive start automatically on system boot
  ```shell
  sudo systemctl enable httpd
  ```

* check all the allowed services
  ```shell
  sudo firewall-cmd --list-all
  ```

* check if http service is enabled
  ```shell
  [root@c30 html]# sudo firewall-cmd --list-all
  public (active)
    target: default
    icmp-block-inversion: no
    interfaces: enp0s3
    sources:
    services: dhcpv6-client ssh
    ports:
    protocols:
    masquerade: no
    forward-ports:
    source-ports:
    icmp-blocks:
    rich rules:
  ```
  There are only **dhcpv6-client ssh** services are enabled.  
  **http** service or port **80** need be to enabled.

* add HTTP service and/or port 80
  ```shell
  sudo firewall-cmd --add-service=http --permanent
  sudo firewall-cmd --add-port=80/tcp --permanent
  ```

* restart firewalld to apply changes
  ```shell
  sudo firewall-cmd --reload
  ```
* check if http service / port 80 is enabled
  ```shell
  [root@c30 html]# sudo firewall-cmd --list-all
  public (active)
    target: default
    icmp-block-inversion: no
    interfaces: enp0s3
    sources:
    services: dhcpv6-client http ssh
    ports: 80/tcp
    protocols:
    masquerade: no
    forward-ports:
    source-ports:
    icmp-blocks:
    rich rules:
  ```

## On Client system, setup local yum repository
* Preventing **yum** from downloading from wrong location
  ```shell
  sudo mv /etc/yum.repos.d/*.repo /tmp/yum_repo_backup/
  ```

* Creating new repository config file:
  ```shell
  cat <<EOF > /etc/yum.repos.d/hcs.repo
  [hcs]
  name=HyperCloud-storage Repository
  baseurl=http://192.168.0.30/hcs-yum
  enabled=1
  gpgcheck=0
  EOF

  # change '192.168.0.30' to your server IP address
  ```

* Installing new packages
  ```shell
  sudo yum clean all
  sudo yum update
  sudo yum install MYPAKCAGE_NAME
  ```
