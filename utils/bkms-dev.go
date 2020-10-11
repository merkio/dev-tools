package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/merkio/dev-tools/config"
)

// PrepareLocalCluster preparation steps, execute env-
func PrepareLocalCluster() {
	setupRepo := "https://gitlab.business-keeper.local/dev/v2/kubernetes/compliance-system/" +
		"env-config/local.git/bkms/env-setup/"
	ExecuteCommandsWithPipe(
		"kustomize", []string{"build", setupRepo},
		"kubectl", []string{"apply", "-f", "-"},
	)
}

// StartS3 start S3 in the local cluster
func StartS3() {
	if IsExistCommand("helm") {
		err := ExecuteCommand("helm", "repo", "add", "minio", "https://helm.min.io/")
		if err != nil {
			log.Fatal("Can't add to the helm repo minio", err)
		}
		err = CreateNameSpace("minio")
		if err != nil {
			fmt.Println("Namespace already exist")
		}

		cfgMap := config.Config()
		awsSecret := cfgMap.GetString("aws_secret_access_key")
		awsAccessKey := cfgMap.GetString("aws_access_key_id")

		err = ExecuteCommand("helm", "install", "--namespace", "minio", "minio", "minio/minio",
			"--set", "persistence.storageClass=local-path",
			"--set", fmt.Sprintf("accessKey=%s,secretKey=%s", awsAccessKey, awsSecret))

		if err != nil {
			log.Fatal("Error during deploy minio", err)
		}
		err = ExecuteCommand("kubectl", "patch", "svc", "minio", "-p",
			"{\"spec\":{\"externalIPs\":[\"192.168.56.10\"]}}", "-n", "minio")
	}
}

// CreateNameSpace create new namespace
func CreateNameSpace(namespace string) error {
	if IsExistCommand("kubectl") {
		return ExecuteCommand("kubectl", "create", "namespace", namespace)
	}
	return nil
}

// CreateBuckets create default buckets
func CreateBuckets() {
	cfgMap := config.Config()
	s3Buckets := cfgMap.GetStringSlice("s3-buckets")

	for _, bucket := range s3Buckets {
		CreateBucket(bucket)
	}
}

// StartLocalCluster start local kubernetes cluster
func StartLocalCluster() {

}

// UpdateServiceDependency update source code of the dependency services
func UpdateServiceDependency(service string) {

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
