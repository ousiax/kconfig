package util

import (
	"github.com/spf13/cobra"

	"k8s.io/klog/v2"
)

func GetFlagString(cmd *cobra.Command, flag string) string {
	s, err := cmd.Flags().GetString(flag)
	if err != nil {
		klog.Fatalf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return s
}

func CheckErr(err error) {
	if err != nil {
		klog.Fatal(err)
	}
}
