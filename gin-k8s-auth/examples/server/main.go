package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
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
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("Failed to create kubernetes client: %v", err)
		os.Exit(1)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Failed to create kubernetes client: %v", err)
		os.Exit(1)
	}

	r.Use(gink8sauth.K8SAuthenticator(client))
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world!",
		})
	})
	fmt.Print("Running hello world server with k8s auth middleware")
	r.Run()
}
