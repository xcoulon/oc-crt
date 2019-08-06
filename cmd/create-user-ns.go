package cmd

import (
	"fmt"
	"os"

	"github.com/codeready-toolchain/oc-crt/pkg/templates"

	log "github.com/Sirupsen/logrus"
	"github.com/openshift/api/template"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	kcmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

var createResourcesCmd *cobra.Command
var yamlFile string
var projectName string
var username string

func init() {
	createResourcesCmd = &cobra.Command{
		Use:               "create-user-ns",
		Short:             "creates a user's namespaces on the cluster",
		PersistentPreRun:  func(cmd *cobra.Command, args []string) {},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {},
		RunE: func(c *cobra.Command, args []string) error {
			if yamlFile == "" {
				return errors.Errorf("missing yaml template file to process")
			}
			return processAndApply(c, yamlFile, kubeconfig)
		},
	}
	createResourcesCmd.Flags().StringVarP(&yamlFile, "file", "f", "", "the yaml template file to process and apply")
	createResourcesCmd.Flags().StringVarP(&projectName, "projectName", "p", "", "the name of the project to create)")
	createResourcesCmd.Flags().StringVarP(&username, "username", "u", "", "the name of the user owning the project to create)")
	// also, init logger
	log.SetOutput(createResourcesCmd.OutOrStdout())
}

func processAndApply(cmd *cobra.Command, yamlFile string, kubeconfig string) error {
	fmt.Fprintf(cmd.OutOrStdout(), "processing '%s'", yamlFile)
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return errors.Wrapf(err, "unable to process and apply template")
	}
	r, err := os.Open(yamlFile)
	if err != nil {
		return errors.Wrapf(err, "unable to process and apply template")
	}
	defer func() {
		if closeErr := r.Close(); closeErr != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "unable to close template file after processing")
		}
	}()
	scheme := runtime.NewScheme()
	// see https://github.com/openshift/oc/blob/master/cmd/oc/oc.go#L77
	utilruntime.Must(template.Install(scheme))
	kubeConfigFlags := genericclioptions.NewConfigFlags(true)
	// kubeConfigFlags.AddFlags(cmd.PersistentFlags())
	matchVersionKubeConfigFlags := kcmdutil.NewMatchVersionFlags(kubeConfigFlags)
	// matchVersionKubeConfigFlags.AddFlags(cmd.PersistentFlags())
	f := kcmdutil.NewFactory(matchVersionKubeConfigFlags)
	return templates.Apply(yamlFile, r, projectName, username, scheme, f, config)
}
