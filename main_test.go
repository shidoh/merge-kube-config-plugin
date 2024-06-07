package main

import (
	"io/ioutil"
	"os"
	"testing"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func TestReadKubeconfig(t *testing.T) {
	// Create a temporary kubeconfig file for testing
	kubeconfigContent := `apiVersion: v1
clusters:
- cluster:
    server: https://example.com
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
  name: test-context
current-context: test-context
kind: Config
users:
- name: test-user
  user:
    token: test-token`

	tmpFile, err := ioutil.TempFile("", "kubeconfig")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(kubeconfigContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	tmpFile.Close()

	// Test readKubeconfig function
	config, err := readKubeconfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("Error reading kubeconfig: %v", err)
	}

	if config.Contexts["test-context"] == nil {
		t.Fatalf("Expected context 'test-context' to be present")
	}
}

func TestWriteKubeconfig(t *testing.T) {
	// Create a temporary file to write the kubeconfig
	tmpFile, err := ioutil.TempFile("", "kubeconfig")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Create a simple kubeconfig struct for testing
	config := &api.Config{
		Contexts: map[string]*api.Context{
			"test-context": {
				Cluster:  "test-cluster",
				AuthInfo: "test-user",
			},
		},
		Clusters: map[string]*api.Cluster{
			"test-cluster": {
				Server: "https://example.com",
			},
		},
		AuthInfos: map[string]*api.AuthInfo{
			"test-user": {
				Token: "test-token",
			},
		},
		CurrentContext: "test-context",
	}

	// Test writeKubeconfig function
	if err := writeKubeconfig(tmpFile.Name(), config); err != nil {
		t.Fatalf("Error writing kubeconfig: %v", err)
	}

	// Read the written file and verify its contents
	data, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Error reading written file: %v", err)
	}

	loadedConfig, err := clientcmd.Load(data)
	if err != nil {
		t.Fatalf("Error loading written kubeconfig: %v", err)
	}

	// Instead of directly comparing the entire config structs, let's compare individual fields
	if config.CurrentContext != loadedConfig.CurrentContext {
		t.Fatalf("Current context does not match. Expected %s, got %s", config.CurrentContext, loadedConfig.CurrentContext)
	}

	if len(config.Contexts) != len(loadedConfig.Contexts) {
		t.Fatalf("Number of contexts does not match. Expected %d, got %d", len(config.Contexts), len(loadedConfig.Contexts))
	}

	for name, context := range config.Contexts {
		loadedContext, ok := loadedConfig.Contexts[name]
		if !ok {
			t.Fatalf("Context %s not found in loaded config", name)
		}
		if context.Cluster != loadedContext.Cluster || context.AuthInfo != loadedContext.AuthInfo {
			t.Fatalf("Context %s does not match. Expected %v, got %v", name, context, loadedContext)
		}
	}

	if len(config.Clusters) != len(loadedConfig.Clusters) {
		t.Fatalf("Number of clusters does not match. Expected %d, got %d", len(config.Clusters), len(loadedConfig.Clusters))
	}

	for name, cluster := range config.Clusters {
		loadedCluster, ok := loadedConfig.Clusters[name]
		if !ok {
			t.Fatalf("Cluster %s not found in loaded config", name)
		}
		if cluster.Server != loadedCluster.Server {
			t.Fatalf("Cluster %s does not match. Expected %v, got %v", name, cluster, loadedCluster)
		}
	}

	if len(config.AuthInfos) != len(loadedConfig.AuthInfos) {
		t.Fatalf("Number of authInfos does not match. Expected %d, got %d", len(config.AuthInfos), len(loadedConfig.AuthInfos))
	}

	for name, authInfo := range config.AuthInfos {
		loadedAuthInfo, ok := loadedConfig.AuthInfos[name]
		if !ok {
			t.Fatalf("AuthInfo %s not found in loaded config", name)
		}
		if authInfo.Token != loadedAuthInfo.Token {
			t.Fatalf("AuthInfo %s does not match. Expected %v, got %v", name, authInfo, loadedAuthInfo)
		}
	}
}

func TestMergeKubeconfigs(t *testing.T) {
	// Create two kubeconfig structs for testing
	config1 := &api.Config{
		Contexts: map[string]*api.Context{
			"context1": {
				Cluster:  "cluster1",
				AuthInfo: "user1",
			},
		},
		Clusters: map[string]*api.Cluster{
			"cluster1": {
				Server: "https://cluster1.example.com",
			},
		},
		AuthInfos: map[string]*api.AuthInfo{
			"user1": {
				Token: "token1",
			},
		},
	}

	config2 := &api.Config{
		Contexts: map[string]*api.Context{
			"context2": {
				Cluster:  "cluster2",
				AuthInfo: "user2",
			},
		},
		Clusters: map[string]*api.Cluster{
			"cluster2": {
				Server: "https://cluster2.example.com",
			},
		},
		AuthInfos: map[string]*api.AuthInfo{
			"user2": {
				Token: "token2",
			},
		},
	}

	// Test mergeKubeconfigs function
	mergedConfig := mergeKubeconfigs(config1, config2)

	// Check the number of contexts, clusters, and authInfos
	if len(mergedConfig.Contexts) != 2 || len(mergedConfig.Clusters) != 2 || len(mergedConfig.AuthInfos) != 2 {
		t.Fatalf("Expected merged config to contain 2 contexts, clusters, and authInfos")
	}

	// Check the presence of contexts
	if mergedConfig.Contexts["context1"] == nil || mergedConfig.Contexts["context2"] == nil {
		t.Fatalf("Expected contexts 'context1' and 'context2' to be present")
	}

	// Check the presence of clusters
	if mergedConfig.Clusters["cluster1"] == nil || mergedConfig.Clusters["cluster2"] == nil {
		t.Fatalf("Expected clusters 'cluster1' and 'cluster2' to be present")
	}

	// Check the presence of authInfos
	if mergedConfig.AuthInfos["user1"] == nil || mergedConfig.AuthInfos["user2"] == nil {
		t.Fatalf("Expected authInfos 'user1' and 'user2' to be present")
	}

	// Verify context details
	if mergedConfig.Contexts["context1"].Cluster != "cluster1" {
		t.Fatalf("Context 'context1' has wrong cluster value")
	}
	if mergedConfig.Contexts["context2"].Cluster != "cluster2" {
		t.Fatalf("Context 'context2' has wrong cluster value")
	}

	// Verify cluster details
	if mergedConfig.Clusters["cluster1"].Server != "https://cluster1.example.com" {
		t.Fatalf("Cluster 'cluster1' has wrong server value")
	}
	if mergedConfig.Clusters["cluster2"].Server != "https://cluster2.example.com" {
		t.Fatalf("Cluster 'cluster2' has wrong server value")
	}

	// Verify authInfo details
	if mergedConfig.AuthInfos["user1"].Token != "token1" {
		t.Fatalf("AuthInfo 'user1' has wrong token value")
	}
	if mergedConfig.AuthInfos["user2"].Token != "token2" {
		t.Fatalf("AuthInfo 'user2' has wrong token value")
	}
}
