package utils

import (
	"fmt"
	"github.com/merkio/dev-tools/config"
	"github.com/shirou/gopsutil/process"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// LocalEnvSetup preparation steps, for launch services
func LocalEnvSetup() {
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
		_ = ExecuteCommand("kubectl", "patch", "svc", "minio", "-p",
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
	if IsExistCommand("vagrant") {
		configMap := config.GetConfigMap()
		err := os.Chdir(configMap.Repositories["k3s-local"]["directory"])
		if err != nil {
			log.Fatal(err)
		}
		err = ExecuteCommand("vagrant", "up")
		if err != nil {
			log.Fatal(err)
		}
		StartS3()
	}
}

// StopLocalCluster start local kubernetes cluster
func StopLocalCluster() {
	if IsExistCommand("vagrant") {
		configMap := config.GetConfigMap()
		err := os.Chdir(configMap.Repositories["k3s-local"]["directory"])
		if err != nil {
			log.Fatal(err)
		}
		err = ExecuteCommand("vagrant", "halt")
		if err != nil {
			log.Fatal(err)
		}
		StartS3()
	}
}

// UpdateServiceDependency update source code of the dependency services
func UpdateServiceDependency(service string) {
	configMap := config.GetConfigMap()

	dServices := configMap.Dependency[service]
	for _, ds := range dServices {
		if ds == "k3s-local" {
			continue
		}
		svc := configMap.Repositories[ds]
		repository := svc["directory"]
		fmt.Printf("Update repository %s for service %s\n", repository, ds)

		PullChanges(repository)
	}
}

// ListPods show list of pods for the namespace (by default env = env)
func ListPods(namespace string) {
	fmt.Printf("List pods for namespace [%s]", namespace)
	if IsExistCommand("kubectl") {
		err := ExecuteCommand("kubectl", "get", "pod", "-n", namespace)

		if err != nil {
			fmt.Println(err)
		}
	}
}

// StartService start single service
func StartService(service string, mode string, trigger string, namespace string) {
	configMap := config.GetConfigMap()

	repo := configMap.Repositories[service]["directory"]
	profile := configMap.Repositories[service]["profile"]

	if namespace != "" {
		namespace = "env"
	}

	startService(repo, "skaffold", argsDefaultCommand(service, profile, mode, trigger, namespace)...)
}

// StartDependencies start all
func StartDependencies(service string) {
	configMap := config.GetConfigMap()

	dServices := configMap.Dependency[service]
	for _, ds := range dServices {
		StartService(ds, "run", "", "env")
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

	if service == "" {
		for svc, values := range configMap.Repositories {
			if svc == "k3s-local" {
				continue
			}

			repo := values["directory"]
			profile := values["profile"]

			killProcesses(svc, namespace)
			stopService(repo, namespace, profile)
		}
	} else {
		repo := configMap.Repositories[service]["directory"]
		profile := configMap.Repositories[service]["profile"]

		killProcesses(service, namespace)
		stopService(repo, namespace, profile)
	}

}

func startService(path string, command string, args ...string) {

	k8s := filepath.Join(path, "k8s")
	err := os.Chdir(k8s)
	if err != nil {
		log.Fatalf("Unable to change directory to %s, %v", k8s, err)
	}
	err = ExecuteCommand(command, args...)
	if err != nil {
		fmt.Println(err)
	}
}

func stopService(path string, command string, args ...string) {

	k8s := filepath.Join(path, "k8s")
	err := os.Chdir(k8s)
	if err != nil {
		log.Fatalf("Unable to change directory to %s, %v", k8s, err)
	}
	err = ExecuteCommand(command, args...)
	if err != nil {
		fmt.Println(err)
	}
}

// killProcesses kill processes with service name and namespace
func killProcesses(service string, namespace string) {

	ps, err := process.Processes()

	if err != nil {
		log.Fatal(err)
	}

	for _, proc := range ps {
		name, e := proc.Name()
		if e != nil {
			fmt.Println(e)
			return
		}
		if strings.Contains(name, "skaffold") || strings.Contains(name, "kubectl") {
			cmdLine, e := proc.Cmdline()
			if e != nil {
				fmt.Println(e)
				return
			}
			if strings.Contains(cmdLine, service) && strings.Contains(cmdLine, namespace) {
				e := proc.Kill()
				if e != nil {
					fmt.Println(e)
				}
			}
		}
	}
}

func argsDefaultCommand(service string, profile string, mode string, trigger string, namespace string) []string {
	args := []string{
		mode, "--insecure-registry", "registry.192.168.56.11.nip.io:5000", "-n", namespace,
		fmt.Sprintf("--default-repo=registry.192.168.56.11.nip.io:5000/bkag/%s/local", service),
		"--port-forward", "--force=false", "--tail",
	}

	if profile != "" {
		args = append(args, []string{"--profile", profile}...)
	}
	if trigger != "" {
		args = append(args, []string{"--trigger", trigger}...)
	}

	return args
}
