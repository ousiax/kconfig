package cert

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	certificatesv1 "k8s.io/api/certificates/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"

	cmduitl "github.com/qqbuby/konfig/cmd/util"
	cmduitlpkix "github.com/qqbuby/konfig/cmd/util/pkix"
)

const (
	flagKubeconfig = "kubeconfig"
	flagUserName   = "username"
	flagGroups     = "group"
	flagExpiration = "expiration"
	flagOutput     = "output"

	expirationSeconds = 60 * 60 * 24 * 365 // one year in seconds
)

type CertOptions struct {
	kubeconfig   string
	clientSet    clientset.Interface
	configAccess clientcmd.ConfigAccess
	csrName      string
	userName     string
	groups       []string
	output       string
}

func NewCmdCert() *cobra.Command {
	o := CertOptions{
		configAccess: clientcmd.NewDefaultPathOptions(),
	}

	cmd := &cobra.Command{
		Use:   "cert",
		Short: "Create kubeconfig file with a specified certificate resources.",
		Run: func(cmd *cobra.Command, args []string) {
			cmduitl.CheckErr(o.Complete(cmd, args))
			cmduitl.CheckErr(o.Validate())
			cmduitl.CheckErr(o.Run())
		},
	}

	kubeconfig := ""
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	cmd.Flags().StringVar(&o.kubeconfig, flagKubeconfig, "", fmt.Sprintf("(optional) absolute path to the kubeconfig file (default %s)", kubeconfig))

	cmd.Flags().StringVarP(&o.userName, flagUserName, "u", "", "user name")
	cmd.MarkFlagRequired(flagUserName)
	cmd.Flags().StringArrayVarP(&o.groups, flagGroups, "g", nil, "group name")
	cmd.MarkFlagRequired(flagGroups)
	cmd.Flags().StringVarP(&o.output, flagOutput, "o", "", "output file - default stdout")

	return cmd
}

func (o *CertOptions) Complete(cmd *cobra.Command, args []string) error {
	o.csrName = o.userName + ":" + strings.Join(o.groups, ":")

	configFlags := &genericclioptions.ConfigFlags{
		KubeConfig: &o.kubeconfig,
	}
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return err
	}
	o.clientSet, err = clientset.NewForConfig(config)
	if err != nil {
		return err
	}
	return nil
}

func (o *CertOptions) Validate() error {
	return nil
}

func (o *CertOptions) Run() error {
	_, err := o.getCertificateSigningRequest()
	if err == nil {
		err := o.deleteCertificatesV1CertificateSigningRequest(err)
		if err != nil {
			return err
		}
	}

	key, request, err := o.createCertificateRequest()
	if err != nil {
		return err
	}
	csr, err := o.createCertificatesV1CertificateSigningRequest(request)
	if err != nil {
		return err
	}

	csr.Status.Conditions = []certificatesv1.CertificateSigningRequestCondition{
		{
			Type:    certificatesv1.CertificateApproved,
			Status:  corev1.ConditionTrue,
			Message: "This CSR was approved by konfig cert approve.",
			Reason:  "KonfigCertApprove",
		},
	}

	csr, err = o.clientSet.CertificatesV1().
		CertificateSigningRequests().
		UpdateApproval(context.TODO(), o.csrName, csr, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	startingConfig, err := o.configAccess.GetStartingConfig()
	if err != nil {
		return err
	}

	ctx := startingConfig.Contexts[startingConfig.CurrentContext]
	kubeconfig := clientcmdapi.Config{
		Clusters: map[string]*clientcmdapi.Cluster{
			ctx.Cluster: startingConfig.Clusters[ctx.Cluster],
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			o.userName: {
				ClientKeyData:         key,
				ClientCertificateData: csr.Status.Certificate,
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			o.userName + "@" + ctx.Cluster: {
				Cluster:   ctx.Cluster,
				AuthInfo:  o.userName,
				Namespace: "default",
			},
		},
		CurrentContext: o.userName + "@" + ctx.Cluster,
	}

	content, err := clientcmd.Write(kubeconfig)
	if err != nil {
		return err
	}

	if len(o.output) != 0 {
		err := os.WriteFile(o.output, content, 0644)
		if err != nil {
			return err
		}
	} else {
		fmt.Fprint(os.Stdout, string(content))
	}

	klog.V(2).Infof("delete csr `%s`.", o.csrName)
	err = o.deleteCertificatesV1CertificateSigningRequest(err)
	if err != nil {
		return err
	}

	return nil
}

func (o *CertOptions) deleteCertificatesV1CertificateSigningRequest(err error) error {
	gracePeriodSeconds := int64(0)
	err = o.clientSet.CertificatesV1().
		CertificateSigningRequests().
		Delete(context.TODO(), o.csrName, metav1.DeleteOptions{
			GracePeriodSeconds: &gracePeriodSeconds,
		})

	return err
}

func (o *CertOptions) createCertificatesV1CertificateSigningRequest(request []byte) (*certificatesv1.CertificateSigningRequest, error) {
	csr, err := o.clientSet.
		CertificatesV1().
		CertificateSigningRequests().
		Create(context.TODO(), &certificatesv1.CertificateSigningRequest{
			ObjectMeta: metav1.ObjectMeta{
				Name: o.csrName,
				Annotations: map[string]string{
					"creator": "konfig.local.io",
				},
			},
			Spec: certificatesv1.CertificateSigningRequestSpec{
				Username: o.userName,
				Groups:   o.groups,
				Usages: []certificatesv1.KeyUsage{
					certificatesv1.UsageClientAuth,
				},
				Request: request,

				SignerName: "kubernetes.io/kube-apiserver-client",
			},
		}, metav1.CreateOptions{})

	return csr, err
}

func (o *CertOptions) getCertificateSigningRequest() (*certificatesv1.CertificateSigningRequest, error) {
	csr, err := o.clientSet.CertificatesV1().
		CertificateSigningRequests().
		Get(context.TODO(), o.csrName, metav1.GetOptions{})
	return csr, err
}

func (o *CertOptions) createCertificateRequest() (keyPem []byte, csrPem []byte, err error) {
	key, csr, err := cmduitlpkix.CreateDefaultCertificateRequest(o.userName, o.groups, nil)
	if err != nil {
		return nil, nil, err
	}

	keyPem, err = cmduitlpkix.PemPkcs8PKey(key)
	if err != nil {
		return nil, nil, err
	}

	csrPem, err = cmduitlpkix.PemCertificateRequest(csr)
	if err != nil {
		return nil, nil, err
	}

	return keyPem, csrPem, nil
}
