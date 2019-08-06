package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUserNamespaces(t *testing.T) {
	// given
	createResourcesCmd.SetOutput(os.Stdout)
	createResourcesCmd.SetArgs([]string{"-f", "/Users/xcoulon/code/go/src/github.com/codeready-toolchain/oc-crt/template.yml", "-n", "test-tmpl"})
	// /apis/authorization.openshift.io/v1/namespaces/test-tmpl/rolebindings/user-edit
	// when
	err := createResourcesCmd.Execute()
	// then
	require.NoError(t, err)
}
