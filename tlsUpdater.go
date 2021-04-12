package main

import (
	"bytes"
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"strings"
	"time"
)

func initTlsUpdater() {

	optionsModifier := func(options *metav1.ListOptions) {
		options.LabelSelector = "tls-updater=true"
		options.TypeMeta = metav1.TypeMeta{
			Kind: "kubernetes.io/tls",
		}
	}

	watchList := cache.NewFilteredListWatchFromClient(
		k8s.CoreV1().RESTClient(),
		string(v1.ResourceSecrets),
		v1.NamespaceAll,
		optionsModifier,
	)

	_, controller := cache.NewInformer(
		watchList,
		&v1.Secret{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				secret := obj.(*v1.Secret)
				dests := secret.Annotations["tls-updater-dests"]
				destSlice := strings.Split(dests, ",")
				fmt.Printf("Secret created: %s. Checking TLS Cert...\n", secret.Name)
				updateCerts(secret, destSlice)
			},

			DeleteFunc: func(obj interface{}) {
				secret := obj.(*v1.Secret)
				fmt.Printf("Secret deleted: %s. Nothing to do...\n", secret.Name)
			},

			UpdateFunc: func(oldObj, newObj interface{}) {
				oldSecret := oldObj.(*v1.Secret)
				newSecret := newObj.(*v1.Secret)

				fmt.Printf("Secret changed %s, %s. Checking TLS Cert...\n", oldSecret.Name, newSecret.Name)
				dests := newSecret.Annotations["tls-updater-dests"]
				destSlice := strings.Split(dests, ",")
				updateCerts(newSecret, destSlice)
			},
		},
	)

	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}

}

func updateCert(secret *v1.Secret, destinationSecretName string) {

	destinationSecret, err := k8s.CoreV1().Secrets(secret.Namespace).Get(context.TODO(), destinationSecretName, metav1.GetOptions{})
	if err != nil {
		return
	}

	updated := false

	if ! bytes.Equal(secret.Data["tls.key"], destinationSecret.Data["tls.key"]) {
		destinationSecret.Data["tls.key"] = secret.Data["tls.key"]
		updated = true
	}

	if ! bytes.Equal(secret.Data["tls.crt"], destinationSecret.Data["tls.crt"]) {
		destinationSecret.Data["tls.crt"] = secret.Data["tls.crt"]
		updated = true
	}

	if updated {
		fmt.Printf("Secret changed %s. Updating secret %s!!\n", secret.Name, destinationSecret.Name)
		_, _ = k8s.CoreV1().Secrets(destinationSecret.Namespace).Update(context.TODO(), destinationSecret, metav1.UpdateOptions{})
	} else {
		fmt.Printf("Secret %s is synched with secret %s. Not updating keypair\n", destinationSecret.Name, secret.Name)
	}

}

func updateCerts(secret *v1.Secret, destinationCerts []string) {

	for _, e := range destinationCerts {
		updateCert(secret, e)
	}

}
