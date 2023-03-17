package main

import (
	"flag"
	"fmt"
	"net/http"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	gink8sauth "github.com/adambkaplan/go-k8s-middleware/gin-k8s-auth"
	"github.com/gin-gonic/gin"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	r := gin.Default()

	client, err := getKubeClient(*kubeconfig)
	if err != nil {
		fmt.Printf("WARN: %v\n", client)
	}
	if client != nil {
		fmt.Printf("INFO: Running hello server with k8s authentication middleware")
		r.Use(gink8sauth.K8SAuthenticator(client))
	}

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world!",
		})
	})
	r.Run()
}

func getKubeClient(kubeconfig string) (kubernetes.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %v", err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %v", err)
	}
	return client, nil
}
