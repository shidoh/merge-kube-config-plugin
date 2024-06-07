package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"os"
)

func main() {
	var kubeconfig1, kubeconfig2, output string

	// Set usage message for the flags
	flag.Usage = func() {
		_, err := fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <kubeconfig1> <kubeconfig2> <output>\n", os.Args[0])
		if err != nil {
			return
		}
		flag.PrintDefaults()
	}

	// Define flags
	flag.StringVar(&kubeconfig1, "kubeconfig1", "", "Path to the first kubeconfig file")
	flag.StringVar(&kubeconfig2, "kubeconfig2", "", "Path to the second kubeconfig file")
	flag.StringVar(&output, "output", "", "Path to the merged kubeconfig file")
	flag.Parse()

	// Check non-flag arguments
	args := flag.Args()
	if len(args) >= 3 {
		kubeconfig1 = args[0]
		kubeconfig2 = args[1]
		output = args[2]
	}

	// Check that all arguments are specified
	if kubeconfig1 == "" || kubeconfig2 == "" || output == "" {
		fmt.Println("You must specify the paths to both kubeconfig files and the output file.")
		flag.Usage()
		return
	}

	// Reading the first kubeconfig file
	config1, err := readKubeconfig(kubeconfig1)
	if err != nil {
		fmt.Printf("Error reading kubeconfig1: %v\n", err)
		return
	}

	// Reading the second kubeconfig file
	config2, err := readKubeconfig(kubeconfig2)
	if err != nil {
		fmt.Printf("Error reading kubeconfig2: %v\n", err)
		return
	}

	// Merging the configurations
	mergedConfig := mergeKubeconfigs(config1, config2)

	// Saving the merged kubeconfig file
	if err := writeKubeconfig(output, mergedConfig); err != nil {
		fmt.Printf("Error writing merged kubeconfig: %v\n", err)
		return
	}

	fmt.Println("Kubeconfig files successfully merged!")
}

// Function to read a kubeconfig file
func readKubeconfig(filepath string) (*api.Config, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.Load(data)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// Function to write a kubeconfig file
func writeKubeconfig(filepath string, config *api.Config) error {
	data, err := clientcmd.Write(*config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath, data, 0644)
}

// Function to merge two kubeconfig files
func mergeKubeconfigs(config1, config2 *api.Config) *api.Config {
	// Creating a new empty configuration object
	mergedConfig := api.NewConfig()

	// Adding contexts, clusters, and users from the first file
	for name, context := range config1.Contexts {
		mergedConfig.Contexts[name] = context
	}
	for name, cluster := range config1.Clusters {
		mergedConfig.Clusters[name] = cluster
	}
	for name, authInfo := range config1.AuthInfos {
		mergedConfig.AuthInfos[name] = authInfo
	}

	// Adding contexts, clusters, and users from the second file
	for name, context := range config2.Contexts {
		mergedConfig.Contexts[name] = context
	}
	for name, cluster := range config2.Clusters {
		mergedConfig.Clusters[name] = cluster
	}
	for name, authInfo := range config2.AuthInfos {
		mergedConfig.AuthInfos[name] = authInfo
	}

	return mergedConfig
}
