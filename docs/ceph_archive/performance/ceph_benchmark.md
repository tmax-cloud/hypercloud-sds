# Ceph Benchmark
> 본 문서는 hypercloud에서 제공하는 rook ceph cluster의 성능분석 방법을 정리한 문서입니다.

## 환경 구성
- rook ceph: `1.1.6`
- kube: `1.15`

## Ceph 성능분석 방법
- Ceph command를 수행하기 위해서는 ceph client의 설치가 필수적입니다.
    - Ceph command를 host를 수행하기 위해서는 ceph common 설치를 통해 ceph client를 해당 host에 설치해야 합니다.  
    - Rook ceph toolbox pod에는 ceph client가 설치되어있기 때문에 toolbox pod에 접속하는 방식을 통해서 ceph command를 사용하실 수 있습니다.
- 성능분석을 하기전에 다음의 command를 통해 시스템에 존재하는 cache를 비웁니다.
    ```shell
    $ sudo echo 3 | sudo tee /proc/sys/vm/drop_caches && sudo sync
    ```
- benchmark tool로 `fio`를 통한 방법을 설명하는 부분이 있습니다. `fio`의 경우, 설정 방법이 아닌 사용 방법에 초점을 맞춰서 기술하였습니다. 

### Pool에 대한 성능분석
> Pool에 대한 성능분석의 경우, `rados`라는 ceph command를 통해서 수행되어 집니다.

- <strong>주의 사항</strong>
    - 기존에 사용하던 pool에 대해 성능분석을 수행할 경우, 해당 pool에 저장되어 있는 object에 영향을 줄 수 있으므로 benchmark 용 pool를 따로 생성하여 진행하는 것을 권장합니다.

- <strong>Pool의 write benchmark</strong>
    - rados bench -p `${poolName}` `${time}` write -t `${concurrentOperations}` -b `${writeObjectSize}` --no-cleanup
        - 특정 pool에 대해서 지정한 시간 동안 write benchmark가 수행됩니다.
        - `--no-cleanup`의 경우, write benchmark에 대한 data를 pool에 남겨 놓는 것을 말합니다.
            - read benchmark를 위해서는 반드시 먼저 write benchmark를 --no-cleanup 옵션과 같이 수행시켜야 합니다.
    ```shell
	# replicapool에 대해서 10초 동안 write benchmark를 수행합니다.
	$ rados bench -p replicapool 10 write --no-cleanup
	hints = 1
    Maintaining 16 concurrent writes of 4194304 bytes to objects of size 4194304 for up to 10 seconds or 0 objects
    Object prefix: benchmark_data_ubuntu-bionic_629
      sec Cur ops   started  finished  avg MB/s  cur MB/s last lat(s)  avg lat(s)
        0      16        16         0         0         0           -           0
        1      16        89        73   291.154       292    0.229995    0.187743
		                ...
       10      16       454       438   175.039        88    0.704013    0.355826
    Total time run:         10.4372
    Total writes made:      454
	Write size:             4194304
	Object size:            4194304
	Bandwidth (MB/sec):     173.994
	Stddev Bandwidth:       64.5201
	Max bandwidth (MB/sec): 292
	Min bandwidth (MB/sec): 88
	Average IOPS:           43
	Stddev IOPS:            16.13
	Max IOPS:               73
	Min IOPS:               22
	Average Latency(s):     0.366409
	Stddev Latency(s):      0.195026
	Max latency(s):         1.10837
	Min latency(s):         0.0281251
	```
- <strong>Pool의 read benchmark</strong>
    - rados bench -p `${poolName}` `${time}` `${readType}` -t `${concurrentOperations}`
        - 특정 pool에 대해서 지정한 시간 동안 read benchmark를 수행합니다.
            - `${readType}`의 경우, rand와 seq가 있습니다.
	```shell
	# replicapool에 대해서 seq read benchmark를 수행합니다.
	$ rados bench -p replicapool 5 seq
	hints = 1
	sec Cur ops   started  finished  avg MB/s  cur MB/s last lat(s)  avg lat(s)
		0      16        16         0         0         0           -           0
		1      16       100        84   335.881       336    0.194277    0.172069
		                        ...
		5      16       409       393   314.243       356   0.0671083    0.198888
	Total time run:       5.18534
	Total reads made:     410
	Read size:            4194304
	Object size:          4194304
	Bandwidth (MB/sec):   316.276
	Average IOPS:         79
	Stddev IOPS:          21.5012
	Max IOPS:             95
	Min IOPS:             41
	Average Latency(s):   0.201392
	Max latency(s):       0.850704
	Min latency(s):       0.0114133
	
	# replicapool에 대해서 rand read benchmark를 수행합니다.
	$ rados bench -p replicapool 5 rand
	# 위의 command와 결과가 출력되는 형식이 동일하므로 해당 command의 결과는 생략합니다.
	```
- <strong>Benchmark 관련 pool의 object 정리</strong>
    - rados -p `${poolName}` cleanup
    ```shell
    $ rados -p replicapool cleanup
    Removed 282 objects
    ```
	
### RBD에 대한 성능분석
> RBD는 `rbd` command 또는 `fio`를 사용하여 성능분석을 진행할 수 있습니다

- <strong>rbd command</strong>
    - `bench-write` option을 통해 deprecate되었고 이제는 `bench` option 사용을 권장합니다. 따라서, 본 문서에서는 `bench` option에 대해서 설명합니다.
    - bench option
        - bench --io-type`[read|write]` --io-size `[size-in-B/K/M/G/T]` --io-threads `[num of threads]` --io-total `[size-in-B/K/M/G/T]` --io-pattern `[seq|rand]`
    ```shell
	$ rbd bench replicapool/test3 --io-type write --io-size 8192 --io-threads 10 --io-total 10G --io-pattern seq
	bench  type write io_size 8388608 io_threads 5 bytes 5368709120 pattern sequential
	  SEC       OPS   OPS/SEC   BYTES/SEC
        1        45     44.46  372959697.41
        2        75     39.12  328143449.11
		          ...
       17       566     33.76  283195575.90
       18       601     32.37  271518807.94
       19       631     34.04  285531375.55
    elapsed:    19  ops:      640  ops/sec:    33.06  bytes/sec: 277317603.36
	```
- <strong>fio</strong>
    - fio가 제공하는 rbd ioengine를 사용하면 fio를 통해 rbd image에 대한 성능분석을 할 수 있습니다.
    - rbd image를 host에 mount시킨 후, 해당 mountpoint를 대상으로 libaio ioengine를 사용하여 rbd image에 대한 성능분석을 하는 방법도 있습니다. 그러나 해당 방법을 사용할 경우, rbd image에 mount된 file system이 성능분석에 주는 영향을 고려해야 됩니다.
    ```shell
    # fio의 parameter관련 설명은 해당 문서에서는 생략합니다. 아래의 참고 부분에 fio doc 링크을 참고해주세요.
    $ fio --ioengine=rbd --name rbdtest --randrepeat=0 --rw=write --clientname=admin --pool=replicapool --rbdname=image0 --invalidate=0 --bs=1M --direct=1 --time_based=1 --runtime=30 --numjobs=1 --iodepth=1 --output=rbdtestresult
    $ cat rbdtestresult
    # fio의 result에 대한 모습은 생략하도록 하겠습니다.
    ```
	
### Cephfs에 대한 성능분석
> Cephfs는 `fio`를 사용하여 성능분석을 진행할 수 있습니다.

- Cephfs에 대한 성능분석을 시작하기 위해서는 먼저 ceph commond를 통해 cephfs를 특정 host에 mount해야 됩니다.
- Cephfs에 대한 mount 작업 완료 후, libaio ioengine를 사용하여 성능분석을 시작하면 됩니다.
    ```shell
    # fio의 parameter관련 설명은 해당 문서에서는 생략합니다. 아래의 참고 부분에 fio doc 링크을 참고해주세요.
    $ fio --directory=/root/cephfs --name cephfstest --direct=1 --rw=randwrite --bs=4k --size=50M --numjobs=16 --time_based --runtime=60 --group_reporting --output=cephfsresult
    $ cat cephfsresult
    # fio의 result에 대한 모습은 스킵하도록 하겠습니다.
    ```
	
## 참고
- https://docs.ceph.com/docs/nautilus/man/8/rados/
- https://docs.ceph.com/docs/nautilus/man/8/rbd/
- https://fio.readthedocs.io/en/latest/