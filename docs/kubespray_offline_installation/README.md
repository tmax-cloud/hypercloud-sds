# How to create Kubernetes cluster in offline environment with Kubespray
외부 인터넷이 접속 안 되는 환경에서는 Kubernetes cluster를 설치하려면 kubespray를 통해서 할 수 있습니다. kubespray offline 설치는 2가지 방식으로 할 수 있습니다.

  1. [Private local environment](kubespray_offline_with_private_server.md)
     - Create package repository (install deb package)
     - Create Webserver (Download binary files)
     - Create Docker registry (Download images)
     - Modify some variables to download and install from local url
  2. [Install from Cache](kubespray_offline_with_cache.md)
     - Manually install all required packages or create local repository
     - Copy all required binary files into cache_dir
     - Copy all required Docker images into cache_dir/images
     - Modify some variables to download and install from cache
