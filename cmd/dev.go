/*Package cmd Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/merkio/dev-tools/utils"
	"github.com/spf13/cobra"
)

var service string
var mode string
var namespace string
var trigger string

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Commands to work with local services (e.g. start, stop, start-dependency etc.)",
}

var startService = &cobra.Command{
	Use:   "start",
	Short: "Start only one service in specific mode",
	Run: func(cmd *cobra.Command, args []string) {
		if mode == "" {
			mode = "dev"
		}
		if service != "" {
			utils.StartService(service, mode, trigger, namespace)
		} else {
			fmt.Println("You need to specify what service do you want to start")
		}
	},
}

var listPodsInNamespace = &cobra.Command{
	Use:   "list",
	Short: "List of pods in the namespace",
	Run: func(cmd *cobra.Command, args []string) {
		if namespace == "" {
			namespace = "env"
		}
		utils.ListPods(namespace)
	},
}

var updateDependencies = &cobra.Command{
	Use:   "update-dep",
	Short: "Update dependencies for the service",
	Run: func(cmd *cobra.Command, args []string) {
		if service == "" {
			log.Fatal("Please provide the service name")
		}
		utils.UpdateServiceDependency(service)
	},
}

var startDependencies = &cobra.Command{
	Use:   "start-dep",
	Short: "Start dependencies for the service",
	Run: func(cmd *cobra.Command, args []string) {
		if service != "" {
			utils.LocalEnvSetup()
			utils.StartDependencies(service)
		} else {
			fmt.Println("You need to specify what service do you want to start")
		}
	},
}

var stopNS = &cobra.Command{
	Use:   "stop",
	Short: "Stop services|service in the namespace",
	Run: func(cmd *cobra.Command, args []string) {
		utils.StopBKMSService(service, namespace)
	},
}

var createS3Buckets = &cobra.Command{
	Use: "create-buckets",
	Short: "Create list of necessary buckets on the local S3",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CreateBuckets()
	},
}

var startLocalCluster = &cobra.Command{
	Use:   "start-cluster",
	Short: "Execute preparation steps for local cluster (e.g. create namespaces, create s3 buckets, create databases)",
	Run: func(cmd *cobra.Command, args []string) {
		utils.StartLocalCluster()
	},
}

var stopLocalCluster = &cobra.Command{
	Use:   "stop-cluster",
	Short: "Execute preparation steps for local cluster (e.g. create namespaces, create s3 buckets, create databases)",
	Run: func(cmd *cobra.Command, args []string) {
		utils.StopLocalCluster()
	},
}

func init() {
	startService.Flags().StringVarP(&service, "service", "s", "", "Service name")
	startService.Flags().StringVarP(&mode, "mode", "m", "", "Service mode")
	startService.Flags().StringVarP(&trigger, "trigger", "t", "", "Trigger to recompile service, by default on every saved changes")
	startService.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace where to start service")

	startDependencies.Flags().StringVarP(&service, "service", "s", "", "Service name")

	updateDependencies.Flags().StringVarP(&service, "service", "s", "", "Service name")

	stopNS.Flags().StringVarP(&service, "service", "s", "", "Service name")
	stopNS.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace")

	listPodsInNamespace.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace")

	devCmd.AddCommand(listPodsInNamespace)
	devCmd.AddCommand(updateDependencies)
	devCmd.AddCommand(createS3Buckets)
	devCmd.AddCommand(startLocalCluster)
	devCmd.AddCommand(stopLocalCluster)
	devCmd.AddCommand(startService)
	devCmd.AddCommand(startDependencies)
	devCmd.AddCommand(stopNS)

	rootCmd.AddCommand(devCmd)
}
