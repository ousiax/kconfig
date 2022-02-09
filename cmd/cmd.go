package cmd

import (
	"flag"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/qqbuby/konfig/cmd/cert"
)

func NewCmdKonfig() *cobra.Command {
	var cmds = &cobra.Command{
		Use: "konfig",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	flags := cmds.Flags()
	logFlags := &flag.FlagSet{}
	klog.InitFlags(logFlags)
	flags.AddGoFlagSet(logFlags)

	cmds.AddCommand(cert.NewCmdCert())
	return cmds
}
