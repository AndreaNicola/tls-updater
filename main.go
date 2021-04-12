package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
	// or
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var k8s *kubernetes.Clientset

func main() {
	inCluster := os.Getenv("IN_CLUSTER")

	var config *rest.Config
	var err error

	if inCluster == "true" {
		config, err = rest.InClusterConfig()
	} else {

		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}

		flag.Parse()
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)

	}

	logrus.Debug("Connecting to cluster...")
	k8s, err = kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Debug("Connection failed!")
		panic(err.Error())
	}

	logrus.Debug("Successfully connected!")

	initTlsUpdater()

}
