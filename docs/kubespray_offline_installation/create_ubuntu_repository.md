# How to Create an Authenticated Repository
> Tested on Ubuntu-server 18.04.3 LTS (bionic)

## Requirements
  1. Packages: apt-utils (기본적으로 설치됨), dpkg-dev, a web server (apache2 or
    nginx), and dpkg-sig  
  2. Base Directory for Repository
  3. .deb file (deb 패키지 파일)  

## Installing the Required Packages
```shell
sudo apt-get install dpkg-dev
sudo apt-get install apache2
sudo apt-get install dpkg-sig
```

## Create the Repository Directory Structure
* apache2 web server 사용할 때 `/var/www/html` Directory에서 만들어야 합니다.
```shell
sudo mkdir -p /var/www/html/hcs-repo/binary
```
* [OR] symbolic link 사용하셔도 됩니다.
```shell
sudo ln -s ~/repo-dir /var/www/html/repo-dir
```
* .deb 패키지 파일을 binary Directory에 이동
```shell
sudo mv /path/to/my/Packages.deb /var/www/html/hcs-repo/binary
```

## Authenticating Repository and Packages
* GPG key pair를 만들기
```shell
gpg --gen-key
# Input Real name, Email, and passphrase
```
* GPG key를 확인하기
```shell
gpg --list-keys
# Output:
#/home/tmax/.gnupg/pubring.kbx
#-----------------------------
#pub   rsa3072 2020-03-12 [SC] [expires: 2022-03-12]
#      9359DA7C2594A5C90E90421E1965FFAEB9D75E4B
#uid           [ultimate] hcs-ck34 <ck3-4team@tmax.co.kr>
#sub   rsa3072 2020-03-12 [E] [expires: 2022-03-12]
```
* Public key를 가져오기
```shell
sudo gpg --output GPGkey --armor --export 9359DA7C2594A5C90E90421E1965FFAEB9D75E4B
```
* Public key (GPGkey) Repository Directory에 복사하기
```shell
sudo mv GPGkey /var/www/html/hcs-repo/GPGkey
```
* Change the ownership of the Directory Structure
```shell
sudo chown -R tmax:tmax -R /var/www/html/hcs-repo
```
* .deb 파일이랑 같은 Directory에서 `Packages` `Packages.gz` index file 만들기
```shell
cd /var/www/html/hcs-repo/binary
apt-ftparchive packages . > Packages
gzip -c Packages > Packages.gz
```
* .deb 파일이랑 같은 Directory에서 `Release` `InRelease` `Release.gpg` file 만들기
```shell
cd /var/www/html/hcs-repo/binary
apt-ftparchive release . > Release
gpg --yes --clearsign -o InRelease Release
gpg --yes -abs -o Release.gpg Release
```

## On Client node
* 우리 Repository url 추가: (Repository server IP는 `192.168.0.200`로 가정합니다.)
```shell
echo 'deb http://192.168.0.200/hcs-repo/binary /' >> /etc/apt/source.list
```
* Download and add Repository's public key (GPGkey)
```shell
wget -O - http://192.168.0.200/hcs-repo/GPGkey | sudo apt-key add -
```
* 추가한 public key 확인하기
```shell
apt-key list
```
* Update and install packages
```shell
sudo apt update
sudo apt install my_package.deb
```
