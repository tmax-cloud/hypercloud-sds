## ENV
ubuntu 18.04.3 LTS

## Setting up the NFS host server
```shell
# Step 1: Install NFS kernel Server
sudo apt update
sudo apt install nfs-kernel-server

# Step 2: Create the Export Directory
sudo mkdir -p /mnt/nfs-shared-dir
# Remove restrictive permissions of export folder
sudo chown nobody:nogroup /mnt/nfs-shared-dir
sudo chmod 777 /mnt/nfs-shared-dir

# Step 3: Assign server access to client(s) through NFS export file
sudo vi /etc/exports
# A single client: 
/mnt/nfs-shared-dir 192.168.0.101(rw,sync,no_subtree_check)
# Multiple clients:
/mnt/nfs-shared-dir 192.168.0.102(rw,sync,no_subtree_check)
/mnt/nfs-shared-dir 192.168.0.103(rw,sync,no_subtree_check)
# Entire subnet that the clients belong to:
/mnt/nfs-shared-dir 192.168.0.1/24(rw,sync,no_subtree_check)
# Accept from all IPs
/mnt/nfs-shared-dir *(rw,sync,no_subtree_check,insecure)

# Step 4: Export the shared directory
sudo exportfs -rav
sudo systemctl restart nfs-kernel-server

# Step 5: Check if firewall is open for the client(s)
sudo ufw allow from 192.168.0.1/24 to any port nfs
sudo ufw status
```

## Configuring the Client Machine
```shell
# Step 1: Install NFS common
sudo apt update
sudo apt install nfs-common

# Step 2: Create a mount point for the NFS host's shared folder
sudo mkdir -p $HOME/nfs-shared-folder

# Step 3: Mount the shared directory on the client
sudo mount ServerIP:/mnt/nfs-shared-dir $HOME/nfs-shared-folder
# NFS version 3
sudo mount -t nfs -o nfsvers=3 ServerIP:/mnt/nfs-shared-dir $HOME/nfs-shared-folder

# Step 4: Set automatically mount at boot
sudo vi /etc/fstab
SeverIP:/mnt/nfs-shared-dir $HOME/nfs-shared-folder nfs auto,nofail,noatime,nolock,intr,tcp,actimeo=1800,vers=3 0 0

$ Step 5: Unmounting an NFS remote share
sudo umount $HOME/nfs-shared-folder
```
## minikube setup
```shell
# Create minikube on Virtualbox
minikube start --vm-driver=virtualbox
# By default "/home" is shared between host and minikube vm
# If shared folder is not in "/home/" directory:
minikube start --vm-driver=virtualbox --mount --mount-string="/mnt/nfs-shared-folder:/minikube-nfs"
```

## Test 1: Create pod with hostpath volume
```yml
apiVersion: v1
kind: Pod
metadata:
  name: busybox1
  labels:
    app: busybox1
spec:
  volumes:
  - name: nfs-hostpath
    hostPath:
      path: /hosthome/k8s/nfs-shared-folder
  containers:
  - image: busybox
    command:
      - sleep
      - "3600"
    imagePullPolicy: IfNotPresent
    name: busybox
    volumeMounts:
    - mountPath: /nfs-shared
      name: nfs-hostpath

```

## Test 2: Create PV, PVC, POD
### Case 1: Default - no specific binding
```yml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: nfs-pv
spec:
  capacity:
    storage: 1Gi
  storageClassName: standard
  persistentVolumeReclaimPolicy: Retain
  accessModes:
  - ReadWriteOnce
  nfs:
    server: 192.168.0.101 #NFS ServerIP
    path: '/mnt/nfs-shared-dir'
    readOnly: false
#  mountOptions:
#  - nfsvers=3
---
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: nfs-pvc
spec:
  storageClassName: standard
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
```
### Case 2: Create pv with `claimRef` to specific pvc
```yml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: nfs-pv-arirang
spec:
  capacity:
    storage: 2Gi
  storageClassName: standard
  persistentVolumeReclaimPolicy: Retain
  accessModes:
  - ReadWriteOnce
  claimRef:
    namespace: default
    name: nfs-pvc-arirang
  nfs:
    server: 192.168.0.101
    path: '/mnt/nfs-shared-dir'
    readOnly: false
---
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: nfs-pvc-arirang
spec:
  storageClassName: standard
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 2Gi
```

### Case 3: Create pvc binding to specific pv using `volumeName`
```yml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: nfs-pv-bibimbap
spec:
  capacity:
    storage: 3Gi
  storageClassName: standard
  persistentVolumeReclaimPolicy: Retain
  accessModes:
  - ReadWriteOnce
  nfs:
    server: 192.168.0.101
    path: '/mnt/nfs-shared-dir'
    readOnly: false
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: nfs-pv-chueotang
spec:
  capacity:
    storage: 1Gi
  storageClassName: standard
  persistentVolumeReclaimPolicy: Retain
  accessModes:
  - ReadWriteOnce
  nfs:
    server: 192.168.0.101
    path: '/mnt/nfs-shared-dir'
    readOnly: false
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: nfs-pv-dakgalbi
spec:
  capacity:
    storage: 2Gi
  storageClassName: standard
  persistentVolumeReclaimPolicy: Retain
  accessModes:
  - ReadWriteOnce
  nfs:
    server: 192.168.0.101
    path: '/mnt/nfs-shared-dir'
    readOnly: false
---
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: nfs-pvc-dakgalbi
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
  storageClassName: standard
  volumeName: nfs-pv-dakgalbi
```

# Create Pod
```yml
apiVersion: v1
kind: Pod
metadata:
  name: busybox-nfs
  labels:
    app: busybox-nfs
spec:
  volumes:
  - name: nfs-pvc
    persistentVolumeClaim:
      claimName: nfs-pvc
  containers:
  - image: busybox
    command:
      - sleep
      - "3600"
    imagePullPolicy: IfNotPresent
    name: busybox-pvc
    volumeMounts:
    - mountPath: /nfs-shared
      name: nfs-pvc
```
