package cmd

// // ContextConfig returns the config of the current context in kubeconfig
// func ContextConfig() (restclient.Config, error) {
// 	var kubeconfig *string
// 	if home := homeDir(); home != "" {
// 		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
// 	} else {
// 		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
// 	}
// 	flag.Parse()

// 	// use the current context in kubeconfig
// 	return clientcmd.BuildConfigFromFlags("", *kubeconfig)
// }
