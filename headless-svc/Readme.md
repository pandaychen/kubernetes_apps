##	Headless Service

生成镜像
```
[root@VM_16_9_centos ~/headless/gserver]# docker build -t headless_svc:v1 .
Sending build context to Docker daemon  12.55MB
Step 1/4 : FROM centos:7
7: Pulling from library/centos
d8d02d457314: Pull complete 
Digest: sha256:307835c385f656ec2e2fec602cf093224173c51119bbebd602c53c3653a3d6eb
Status: Downloaded newer image for centos:7
 ---> 67fa590cfc1c
Step 2/4 : LABEL maintainer="pandaychen<ringbuffer@126.com>"
 ---> Running in 0d756382f7f8
Removing intermediate container 0d756382f7f8
 ---> ce39d525e733
Step 3/4 : COPY server /data/
 ---> bf518a1a2654
Step 4/4 : CMD ["/bin/sh", "-c", "cd /data/ && /data/server"]
 ---> Running in 5139ad1658e3
Removing intermediate container 5139ad1658e3
 ---> 03e1b7aa8882
Successfully built 03e1b7aa8882
Successfully tagged headless_svc:v1
[root@VM_16_9_centos ~/headless/gserver]# 

```

查看下images

```
[root@VM_16_9_centos ~/headless/gserver]# docker images |grep head
headless_svc                                               v1                  03e1b7aa8882        2 minutes ago       214MB

```

创建SVC
```
[root@VM_16_9_centos ~/headless]# kubectl apply -f headless_svc.yaml 
deployment.extensions/headsvc-deployment created

```

查看创建的SVC
```
[root@VM_16_9_centos ~/headless]# kubectl get svc|grep hea
headsvc-service      ClusterIP      None           <none>            50051/TCP           44s

```

查看我们创建的POD
```
[root@VM_16_9_centos ~/headless/gclient]# kubectl get pods
NAME                                    READY   STATUS      RESTARTS   AGE
etcd-0                                  1/1     Running     0          53d
etcd-1                                  1/1     Running     0          53d
etcd-2                                  1/1     Running     0          53d
headsvc-deployment-79b8d55dd5-5ccbh     1/1     Running     0          2m44s
headsvc-deployment-79b8d55dd5-972pq     1/1     Running     0          2m44s
headsvc-deployment-79b8d55dd5-h64g2     1/1     Running     0          2m44s
headsvc-deployment-79b8d55dd5-l278z     1/1     Running     0          2m44s
headsvc-deployment-79b8d55dd5-tsl72     1/1     Running     0          2m44s
headsvc-deployment-79b8d55dd5-zrwl5     1/1     Running     0          2m44s
```

查询HeadlessSvc是否生效，在任意一个K8S集群内的容器执行nslookup
```
[root@headsvc-7c5b5d9cd5-2ltj8 /]# nslookup headsvc
Server:		10.3.255.82
Address:	10.3.255.82#53

Name:	headsvc.default.svc.cluster.local
Address: 10.0.0.243
Name:	headsvc.default.svc.cluster.local
Address: 10.0.0.242
Name:	headsvc.default.svc.cluster.local
Address: 10.0.0.241
Name:	headsvc.default.svc.cluster.local
Address: 10.0.0.246
Name:	headsvc.default.svc.cluster.local
Address: 10.0.0.245
Name:	headsvc.default.svc.cluster.local
Address: 10.0.0.244
```


运行客户端
```
./client  -port=50051 -svc headsvc.default.svc.cluster.local
```

验证HeadlessSvc的Pod负载均衡方式
```
for var in `kubectl get pod|grep headsvc|awk '{print $1}'`; do  echo $var;kubectl logs $var; done

headsvc-7c5b5d9cd5-2ltj8
Recv: gRPC
Recv: gRPC
headsvc-7c5b5d9cd5-958qb
Recv: gRPC
Recv: gRPC
headsvc-7c5b5d9cd5-jbmp6
Recv: gRPC
Recv: gRPC
headsvc-7c5b5d9cd5-lv84r
Recv: gRPC
Recv: gRPC
headsvc-7c5b5d9cd5-pcbvp
Recv: gRPC
Recv: gRPC
headsvc-7c5b5d9cd5-tfbrg
Recv: gRPC
Recv: gRPC
```
