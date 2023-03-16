package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	authnv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	var kubeconfig *string
	var serverURL *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	serverURL = flag.String("server-url", "http://localhost:8080", "url to server")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("Failed to find kubeconfig: %v\n", err)
		os.Exit(1)
	}

	token, err := getAuthToken(config)
	if err != nil {
		fmt.Printf("Failed to get auth token: %v\n", err)
		os.Exit(1)
	}
	response, err := sendHello(token, *serverURL)
	if err != nil {
		fmt.Printf("Failed to send hello: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Received response: %s\n", response)
	fmt.Println("Done!")
	os.Exit(0)
}

func getAuthToken(config *rest.Config) (string, error) {
	if config.BearerToken != "" {
		return config.BearerToken, nil
	}
	if config.BearerTokenFile != "" {
		token, err := os.ReadFile(config.BearerTokenFile)
		if err != nil {
			return "", fmt.Errorf("failed to read bearer token file %s: %v", config.BearerTokenFile, err)
		}
		return string(token), nil
	}
	kclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", fmt.Errorf("failed to create kubernetes client: %v", kclient)
	}
	// KLUGE: If we have no token, ask for one from the default service account in the default namespace
	// This is not recommended for anything in production.
	token, err := requestDefaultSAToken(kclient)
	if err != nil {
		return "", err
	}

	return token, nil
}

func requestDefaultSAToken(client kubernetes.Interface) (string, error) {
	tokenRequest := &authnv1.TokenRequest{
		Spec: authnv1.TokenRequestSpec{
			Audiences: []string{},
		},
	}
	tokenRequest, err := client.CoreV1().ServiceAccounts("default").CreateToken(context.TODO(), "default", tokenRequest, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create serviceaccount token: %v", err)
	}
	return tokenRequest.Status.Token, nil
}

func sendHello(bearerToken, server string) (string, error) {
	url, err := url.Parse(server)
	if err != nil {
		return "", err
	}
	// Handle fun edge cases around URL parsing
	if url.Scheme == "" {
		url.Scheme = "http"
	}
	if url.Host == "" && url.Path != "" {
		// URLs like "example.com" do not get parsed as a host, but as a path.
		// This is to adhere to RFC3986 Sec. 3 https://www.rfc-editor.org/rfc/rfc3986#section-3
		// Switch host and path in this case.
		url.Host = url.Path
		url.Path = ""
	}
	// Set the path to "hello" for the request
	url.Path = "hello"
	req, err := http.NewRequest(http.MethodGet, url.String(), &bytes.Buffer{})
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer: %s", bearerToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("hello request failed: %v", err)
	}
	buf := &bytes.Buffer{}
	defer resp.Body.Close()
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	body := buf.String()
	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("response returned status code %d", resp.StatusCode)
	}
	return body, nil
}
