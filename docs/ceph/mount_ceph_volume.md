## Mount Rook Ceph Volume
> 본 문서는 특정 host 에서 ceph client 를 통해 rook ceph cluster 의 rbd image 와 cephfs volume 를 해당 host 에 mount 하는 방법에 대해서 다룹니다.

### Rook Ceph Volume 를 Mount 하는 방법
1. <strong>Ceph client 설치</strong>
    - sudo apt install ceph-common
	
2. <strong>Rook ceph cluster 와의 통신을 위해서 rook ceph 의 keyring 과 config file 를 host 의 /etc/ceph dir 로 가져옵니다.</strong>
    - Rook ceph cluster 관련 pod 이 deploy 된 host 의 경우, host 의 /var/lib/rook/rook-ceph dir 에 rook ceph cluster 에 대한 keyring 과 config file 이 존재합니다.
    - Rook ceph cluster 의 config 을 수정하여 ceph client가 사용할 ceph config 파일을 만듭니다.
        - Rook ceph config 파일에서 fsid 및 mon 에 관련된 설정만 남기고 제거합니다.
            - Ceph client 는 config 에 설정된 mon 의 ip 정보를 통해 ceph mon 에 접근합니다.
            - Ceph client 를 통해 ceph command 를 실행하기 전에 해당 host 에서 ceph config 에 설정된 `mon ip`로 접근할 수 있는지 확인하시는 것을 권장합니다.
        - `client.admin` section 에있는 keyring 의 path 를 자신이 복사한 keyring 파일의 위치로 수정합니다.
        ```shell
        # 본 예시는 rook ceph 관련 pod 이 deploy 된 host 에서 진행되는 것을 가정합니다. 
        $ sudo ls /var/lib/rook/rook-ceph/
        client.admin.keyring  log  rook-ceph.config
       
        # client.admin.keyring 파일과 rook-ceph.config 를 해당 host의 /etc/ceph dir 에 복사합니다.
        $ sudo cp /var/lib/rook/rook-ceph/client.admin.keyring /etc/ceph
        
        # rook-ceph.config 의 경우, 해당 파일의 이름을 ceph.conf 로 변경해서 복사합니다.
        $ sudo cp /var/lib/rook/rook-ceph/rook-ceph.config /etc/ceph/ceph.conf
        
        $ sudo ls /etc/ceph
        ceph.conf  client.admin.keyring  rbdmap
        
        # ceph.conf 를 수정합니다.
        $ sudo cat /etc/ceph/ceph.conf
        [global]
        fsid                          = 87f4740c-83f2-4ccd-8a0f-f3187cc36906
        mon initial members     = a 
        mon host                  = 10.110.11.66:6789

        [client.admin]
        keyring = /etc/ceph/client.admin.keyring
        
        # Ceph client 의 동작을 확인합니다.
        $ ceph -s
        ```

3. <strong>Ceph volume mount</strong>
    - RBD(block storage), cephfs(file storage) 를 사용하기 위해서는 mds, pool 생성등 추가적으로 필요한 작업들이 존재하지만 이미 완료되었다고 가정합니다.
    - <strong>RBD image mount</strong>
        1. RBD image 를 생성합니다.
            - Command
                - rbd create `${poolName}`/`${imageName}` --size `${sizeInMegabytes}` --image-feature layering
                - rbd create `${imageName}` -p `${poolName}` --size `${sizeInMegabytes}` --image-feature layering
                - rbd create --size 1000 replicapool/test --image-feature layering
			- RBD image feature
			    - RBD kernel module 의 경우, ceph 에서 제공하는 모든 rbd image feature 를 지원하지 않습니다. 따라서, rbd image 를 mount 하여 사용하기 위해서는 rbd image 의 feature 를 layering 으로 설정해야 합니다.  
        2. RBD image 를 해당 host 에 device 로 map 합니다.
            - Map 이란 rbd kernel module 를 통해서 지정된 rbd image 를 host 에 device 로 매핑하는 것을 말합니다.
            - Command
                - sudo rbd map `${poolName}`/`${imageName}`
                - sudo rbd map replicapool/test
			```shell
			$ rbd map replicapool/test
			/dev/rbd2
			
			# lsblk command 를 통해 해당 host 에 rbd image 가 map 된 것을 확인할 수 있습니다.
			$ lsblk
			```
        3. RBD image 을 filesystem 로 포맷합니다.
            - Command
                - sudo mkfs.ext4 -m0 /dev/rbd/`${poolName}`/`${imageName}`
                    - mkfs 는 block storage device 를 특정 filesystem 으로 포맷할 때 사용됩니다.
                    - `-m` 을 사용하여, filesystem 에 super user 을 위해 reserve 하는 공간의 퍼센트를 설정할 수 있습니다.
                - sudo mkfs.ext4 -m0 /dev/rbd/replicapool/test
        4. Host 에 rbd image 기반의 filesystem 을 mount 합니다.
            - Command
                - sudo mount /dev/rbd/`${poolName}`/`${imageName}`  `${mountPoint}`
                - sudo mount /dev/rbd/replicapool/test /home/k8s/testdir
        
    - <strong>Cephfs volume mount</strong>
        - Cephfs volume 를 그대로 mount 하거나 cephfs volume 에 대해서 subvolume 생성을 통해서 여러 개의 partition 으로 나누어 사용할 수 있습니다.
        - Ceph cluster 에 존재하는 cephfs volume 에 대해서 cephfs subvolume를 생성하는 방법
            1. Ceph cluster 에 존재하는 cephfs volume 을 확인합니다.
                - Command
                    - ceph fs volume ls
                - Ceph cluster 에 존재하는 cpehfs volume 이 없다면 cephfs volume 를 생성해야 됩니다.
                ```shell
                $ ceph fs volume ls
                [
                   {
                      "name": "myfs"
                   }
                ]
                ```
                
            2. Cephfs volume 의 subvolume 를 생성합니다.
                - Command
                    - ceph fs subvolume create `${volumeName}` `${subVolumeName}` `${sizeInByte}`
                    - ceph fs subvolume create myfs testsubvol 107374182400
                - Cephfs subvolumegroup 사용하여, subvolume 를 관리할 수 있지만 본 문서에서는 subvolumegroup 에 대해서는 다루지 않습니다.
                ```shell
                $ ceph fs subvolume create myfs testsubvol 107374182400
                ```
                
            3. Cephfs subvolume 의 생성을 확인합니다.
                - Command
                    - ceph fs subvolume ls `${volumeName}`
					    - `ceph fs subvolume ls`의 경우, ceph v14.2.5 부터 지원됩니다.
                    - ceph fs subvolume ls myfs
					
                ```shell
                $ ceph fs subvolume ls myfs
                [
                   {
                      "name": "testsubvol"
                   }
                ]
                ```
                
        -  Host 에 cephfs volume 를 mount 합니다.
            - Command
                - sudo mount -t ceph `${monIp}`:`${subDir}` `${mountPoint}` -o name=`${user}`,secret=`${keyValue}`
                ```shell
                # Ceph 의 ceph.conf, keyring 파일을 확인합니다.
                $ cat /etc/ceph/ceph.conf
                [global]
                fsid                      = b817b0f9-1579-4b81-a543-c6955da811d3
                mon initial members       = a
                mon host                  = 10.106.114.161:6789

                [client.admin]
                keyring = /etc/ceph/client.admin.keyring
                
                $ cat /etc/ceph/client.admin.keyring
                [client.admin]
                    key = AQCkOZle8LljFhAAdCOO+l3JKRVvHazRc0zyHw==
                    caps mds = "allow *"
                    caps mon = "allow *"
                    caps osd = "allow *"
                    caps mgr = "allow *"
                    
                # Cephfs volume 를 host 에 mount 합니다.
                $ sudo mount -t ceph 10.106.114.161:6789:/ /home/k8s/cephfsdir -o name=admin,secret=AQCkOZle8LljFhAAdCOO+l3JKRVvHazRc0zyHw==
                
                # Cephfs subvolume 를 host 에 mount 합니다.
				
                # Cephfs subvolume 의 path 를 확인합니다. 
                # ceph fs subvolume getpath ${volumeName} ${subVolumeName}
                $ ceph fs subvolume getpath myfs testsub
                /volumes/_nogroup/testsub
                
                # Cephfs subvolume mount
                $ sudo mount -t ceph 10.106.114.161:6789:/volumes/_nogroup/testsub /home/k8s/cephfsdir -o name=admin,secret=AQCkOZle8LljFhAAdCOO+l3JKRVvHazRc0zyHw==
                ```

### 참고
- https://docs.ceph.com/docs/nautilus/man/8/rbd/
- https://docs.ceph.com/docs/nautilus/start/quick-rbd/
- https://docs.ceph.com/docs/nautilus/cephfs/fs-volumes/
- https://linux.die.net/man/8/mkfs.ext4
- https://docs.ceph.com/docs/nautilus/man/8/mount.ceph/
- https://docs.ceph.com/docs/master/releases/nautilus/
