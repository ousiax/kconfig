package cmd

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"

	"github.com/qqbuby/kconfig/cmd/cert"
	"github.com/qqbuby/kconfig/cmd/version"
)

func NewCmdKonfig() *cobra.Command {
	var cmds = &cobra.Command{
		Use: "kconfig",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	flags := cmds.PersistentFlags()
	logFlags := &flag.FlagSet{}
	klog.InitFlags(logFlags)
	flags.AddGoFlagSet(logFlags)

	var kubeconfig string
	defaultKubeConfig := ""
	if home := homedir.HomeDir(); home != "" {
		defaultKubeConfig = filepath.Join(home, ".kube", "config")
	}
	flags.StringVar(&kubeconfig, "kubeconfig", "", fmt.Sprintf("(optional) absolute path to the kubeconfig file (default %s)", defaultKubeConfig))
	configFlags := &genericclioptions.ConfigFlags{
		KubeConfig: &kubeconfig,
	}

	cmds.AddCommand(cert.NewCmdCert(configFlags))
	cmds.AddCommand(version.NewCmdVersion(configFlags))

	return cmds
}
