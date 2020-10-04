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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// s3Cmd represents the s3 command
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Operations with s3 storage",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("s3 called")
	},
}

var s3ListBuckets = &cobra.Command{
	Use:   "list-buckets",
	Short: "List of buckets",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("s3 list buckets")
	},
}

var s3ListObjects = &cobra.Command{
	Use:   "list-objects",
	Short: "List of objects",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("s3 list buckets")
	},
}

func init() {
	viper.BindPFlag("account-id", s3Cmd.PersistentFlags().Lookup("accountId"))

	s3Cmd.AddCommand(s3ListBuckets)
	rootCmd.AddCommand(s3Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// s3Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// s3Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
