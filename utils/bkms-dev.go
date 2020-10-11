package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/merkio/dev-tools/config"
)

// CreateDatabase preparation step, creates dbs for the services
func CreateDatabase(service string, dbName string) {

}

// StartS3 start S3 in the local cluster
func StartS3() {

}

// CreateBuckets create default buckets
func CreateBuckets() {

}

// CreateNamespaces create namespaces
func CreateNamespaces() {

}

// StartLocalCluster start local kubernetes cluster
func StartLocalCluster() {

}

// UpdateServiceDependency update source code of the dependency services
func UpdateServiceDependency() {

}

// StartService start single service
func StartService(service string, mode string) {
	configMap := config.GetConfigMap()

	repo := configMap.Repositories[service]["directory"]
	profile := configMap.Repositories[service]["profile"]

	startService(repo, profile, mode)
}

// StartDependencies start all
func StartDependencies(service string) {
	configMap := config.GetConfigMap()

	dServices := configMap.Dependency[service]
	for _, ds := range dServices {
		repo := configMap.Repositories[ds]["directory"]
		profile := configMap.Repositories[ds]["profile"]

		startService(repo, profile, "run")
	}
}

// StopBKMSService stop the service in the namespace with profile
func StopBKMSService(service string, namespace string) {
	fmt.Printf("Stop service %s in namespace %s", service, namespace)

	if namespace == "" {
		fmt.Println("Namespace is not specified, using 'env' namespace")
		namespace = "env"
	}

	configMap := config.GetConfigMap()

	repo := configMap.Repositories[service]["directory"]
	profile := configMap.Repositories[service]["profile"]

	stopService(repo, namespace, profile)
}

// StopBKMSServices stop all services from the config file in the namespace
func StopBKMSServices(namespace string) {
	configMap := config.GetConfigMap()

	if namespace == "" {
		fmt.Println("Namespace is not specified, using 'env' namespace")
		namespace = "env"
	}

	for service, values := range configMap.Repositories {
		if service == "k3s-local" {
			continue
		}
		stopService(values["directory"], namespace, values["profile"])
	}
}

func startService(path string, profile string, mode string) {

	k8s := filepath.Join(path, "k8s")
	err := os.Chdir(k8s)
	if err != nil {
		log.Fatalf("Unable to change directory to %s, %v", k8s, err)
	}
	RunService(k8s, profile, mode)
}

func stopService(path string, namespace string, profile string) {

	k8s := filepath.Join(path, "k8s")
	err := os.Chdir(k8s)
	if err != nil {
		log.Fatalf("Unable to change directory to %s, %v", k8s, err)
	}
	StopService(k8s, namespace, profile)
}
