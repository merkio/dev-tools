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

	"github.com/merkio/dev-tools/utils"
	"github.com/spf13/cobra"
)

var service string
var mode string
var namespace string

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Commands to work with local services (e.g. start, stop, start-dependency etc.)",
	Long:  ``,
}

var startService = &cobra.Command{
	Use:   "start",
	Short: "Start only one service in specific mode",
	Run: func(cmd *cobra.Command, args []string) {
		if mode == "" {
			mode = "dev"
		}
		if service != "" {
			utils.StartService(service, mode)
		} else {
			fmt.Println("You need to specify what service do you want to start")
		}
	},
}

var startDependencies = &cobra.Command{
	Use:   "start-dep",
	Short: "Start dependencies for the service",
	Run: func(cmd *cobra.Command, args []string) {
		if service != "" {
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
		if service != "" {
			utils.StopBKMSService(service, namespace)
		} else {
			utils.StopBKMSServices(namespace)
		}
	},
}

func init() {
	startService.Flags().StringVarP(&service, "service", "s", "", "Service name")
	startService.Flags().StringVarP(&mode, "mode", "m", "", "Service mode")

	startDependencies.Flags().StringVarP(&service, "service", "s", "", "Service name")

	stopNS.Flags().StringVarP(&service, "service", "s", "", "Service name")
	stopNS.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace")

	devCmd.AddCommand(startService)
	devCmd.AddCommand(startDependencies)
	devCmd.AddCommand(stopNS)

	rootCmd.AddCommand(devCmd)
}
