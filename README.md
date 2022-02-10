## Usage

```console
$ ./kconfig help cert
Create kubeconfig file with a specified certificate resources.

Usage:
  kconfig cert [flags]

Flags:
  -g, --group stringArray   group name
  -h, --help                help for cert
      --kubeconfig string   (optional) absolute path to the kubeconfig file (default /home/x/.kube/config)
  -o, --output string       output file - default stdout
  -u, --username string     user name

$ ./kconfig cert -u hello -g hello -o hello.config

$ kubectl get po --kubeconfig hello.config 
NAME                          READY   STATUS    RESTARTS       AGE
echoserver-78cc7857c5-d7bfb   1/1     Running   10 (17h ago)   43d
nginx-765b5f545d-4rn74        1/1     Running   10 (17h ago)   43d
nginx-765b5f545d-kv45x        1/1     Running   10 (17h ago)   43d
```
