## Usage

```console
$ ./konfig cert -u hello -g hello -o hello.config

$ kubectl create clusterrolebinding --clusterrole=view --user hello hello:view
clusterrolebinding.rbac.authorization.k8s.io/hello:view created

$ kubectl get po --kubeconfig hello.config 
NAME                          READY   STATUS    RESTARTS      AGE
echoserver-78cc7857c5-d7bfb   1/1     Running   9 (47h ago)   43d
nginx-765b5f545d-4rn74        1/1     Running   9 (47h ago)   43d
nginx-765b5f545d-kv45x        1/1     Running   9 (47h ago)   43d
```
