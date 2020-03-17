# Create Docker registry
> Tested on Ubuntu-server 18.04.3 LTS (bionic)

## Setup Docker environment
* Install packages to allow apt to use a repository over HTTP/HTTPS
```shell
apt-get update && apt-get install -y \
          apt-transport-https ca-certificates curl software-properties-common gnupg2
```
* Add Docker’s official GPG key
```shell
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
```

* Add Docker apt repository.
```shell
add-apt-repository \
    "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
    $(lsb_release -cs) stable"
```
* Install Docker CE.
```shell
apt-get update && apt-get install -y \
        containerd.io=1.2.10-3 \
        docker-ce=6:19.03.4~3-0~ubuntu-$(lsb_release -cs) \
        docker-ce-cli=5:19.03.4~3-0~ubuntu-$(lsb_release -cs)
```
* Setup daemon.
```shell
IP=$(hostname -I | cut -d' ' -f1)
cat > /etc/docker/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2",
  "insecure-registries": ["${IP}:5000"]
}
EOF
```

* Restart docker.
```shell
mkdir -p /etc/systemd/system/docker.service.d
systemctl daemon-reload
systemctl restart docker
```

## Setup Docker registry
> Server IP는 `192.168.0.200`로 가정합니다.

* Start docker registry
```shell
sudo docker run -it -d -p 5000:5000 \
            -v ~/docker_images:/var/lib/registry \
            --name hcs-registry registry:latest
```
* Docker registry 확인하기
```shell
curl -X GET http://192.168.0.200:5000/v2/_catalog
```
* Pull some images from the original hub
```shell
sudo docker pull gcr.io/google-containers/kube-proxy:v1.17.2
```
* Tag the image so that it points to our registry
```shell
sudo docker tag gcr.io/google-containers/kube-proxy:v1.17.2 \
            192.168.0.200:5000/google-containers/kube-proxy:v1.17.2
```
* Push image to our registry
```shell
sudo docker push 192.168.0.200:5000/google-containers/kube-proxy:v1.17.2
```
* Docker registry 확인하기
```shell
curl -X GET 192.168.0.200:5000/v2/_catalog
curl -X GET 192.168.0.200:5000/v2/google-containers/kube-proxy/tags/list
```
