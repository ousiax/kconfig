## Usage

```console
$ ./build.sh 

$ ./kconfig version
Client Version: version.Info{Major:"", Minor:"", GitVersion:"v0.3.0-0.23.3", GitCommit:"496dce3db7c59861ee99e347fa7f4f0e366a895e", GitTreeState:"", BuildDate:"2022-02-14T09:38:31Z", GoVersion:"go1.17.2", Compiler:"gc", Platform:"linux/amd64"}
Server Version: version.Info{Major:"1", Minor:"23", GitVersion:"v1.23.0", GitCommit:"ab69524f795c42094a6630298ff53f3c3ebab7f4", GitTreeState:"clean", BuildDate:"2021-12-07T18:09:57Z", GoVersion:"go1.17.3", Compiler:"gc", Platform:"linux/amd64"}

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
