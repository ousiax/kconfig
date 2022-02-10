package version

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	apimachineryversion "k8s.io/apimachinery/pkg/version"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/pkg/version"
	"k8s.io/client-go/tools/clientcmd"

	cmdutil "github.com/qqbuby/kconfig/cmd/util"
)

// Version is a struct for version information
type Version struct {
	ClientVersion *apimachineryversion.Info `json:"clientVersion,omitempty" yaml:"clientVersion,omitempty"`
	ServerVersion *apimachineryversion.Info `json:"serverVersion,omitempty" yaml:"serverVersion,omitempty"`
}

var (
	versionExample = `
		# Print the client and server versions for the current context
		kconfig version`
)

// Options is a struct to support version command
type Options struct {
	ClientOnly bool
	Short      bool
	Output     string
	Out        io.Writer

	discoveryClient discovery.CachedDiscoveryInterface
}

// NewCmdVersion returns a cobra command for fetching versions
func NewCmdVersion(configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	o := Options{
		Out: os.Stdout,
	}
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print the client and server version information",
		Long:    "Print the client and server version information for the current context.",
		Example: versionExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(configFlags))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}
	cmd.Flags().BoolVar(&o.ClientOnly, "client", o.ClientOnly, "If true, shows client version only (no server required).")
	cmd.Flags().BoolVar(&o.Short, "short", o.Short, "If true, print just the version number.")
	cmd.Flags().StringVarP(&o.Output, "output", "o", o.Output, "One of 'yaml' or 'json'.")
	return cmd
}

// Complete completes all the required options
func (o *Options) Complete(configFlags *genericclioptions.ConfigFlags) error {
	var err error
	if o.ClientOnly {
		return nil
	}
	o.discoveryClient, err = configFlags.ToDiscoveryClient()
	// if we had an empty rest.Config, continue and just print out client information.
	// if we had an error other than being unable to build a rest.Config, fail.
	if err != nil && !clientcmd.IsEmptyConfig(err) {
		return err
	}
	return nil
}

// Validate validates the provided options
func (o *Options) Validate() error {
	if o.Output != "" && o.Output != "yaml" && o.Output != "json" {
		return errors.New(`--output must be 'yaml' or 'json'`)
	}

	return nil
}

// Run executes version command
func (o *Options) Run() error {
	var (
		serverVersion *apimachineryversion.Info
		serverErr     error
		versionInfo   Version
	)

	clientVersion := version.Get()
	versionInfo.ClientVersion = &clientVersion

	if !o.ClientOnly && o.discoveryClient != nil {
		// Always request fresh data from the server
		o.discoveryClient.Invalidate()
		serverVersion, serverErr = o.discoveryClient.ServerVersion()
		versionInfo.ServerVersion = serverVersion
	}

	switch o.Output {
	case "":
		if o.Short {
			fmt.Fprintf(o.Out, "Client Version: %s\n", clientVersion.GitVersion)
			if serverVersion != nil {
				fmt.Fprintf(o.Out, "Server Version: %s\n", serverVersion.GitVersion)
			}
		} else {
			fmt.Fprintf(o.Out, "Client Version: %s\n", fmt.Sprintf("%#v", clientVersion))
			if serverVersion != nil {
				fmt.Fprintf(o.Out, "Server Version: %s\n", fmt.Sprintf("%#v", *serverVersion))
			}
		}
	case "yaml":
		marshalled, err := yaml.Marshal(&versionInfo)
		if err != nil {
			return err
		}
		fmt.Fprintln(o.Out, string(marshalled))
	case "json":
		marshalled, err := json.MarshalIndent(&versionInfo, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(o.Out, string(marshalled))
	default:
		// There is a bug in the program if we hit this case.
		// However, we follow a policy of never panicking.
		return fmt.Errorf("VersionOptions were not validated: --output=%q should have been rejected", o.Output)
	}

	return serverErr
}
