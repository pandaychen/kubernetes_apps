# K8S服务的健康检查


## 就绪探针(Readiness)

kubelet 使用 readiness probe（就绪探针）来确定容器是否已经就绪可以接受流量。只有当 Pod 中的容器都处于就绪状态时 kubelet 才会认定该 Pod处于就绪状态。该信号的作用是控制哪些 Pod应该作为service的后端。如果 Pod 处于非就绪状态，那么它们将会被从 service 的 load balancer中移除。


## 存活探针(Liveness)


kubelet 使用 liveness probe（存活探针）来确定何时重启容器。例如，当应用程序处于运行状态但无法做进一步操作，liveness 探针将捕获到 deadlock，重启处于该状态下的容器，使应用程序在存在 bug 的情况下依然能够继续运行下去。

 
## YAML配置livenessProbe

```
livenessProbe:
  exec:
    command:
    - /root/client
    - -a
    - 127.0.0.1:19000
  initialDelaySeconds: 2
  periodSeconds: 2

k8s会每2秒执行 /root/rpc_check -a 127.0.0.1:19000, 执行成功的话代表存活,不成功的话k8s会重启pod
```
