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

var accountID string
var bucket string

// s3Cmd represents the s3 command
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Operations with s3 storage",
}

var s3ListBuckets = &cobra.Command{
	Use:   "list-buckets",
	Short: "List of buckets",
	Run: func(cmd *cobra.Command, args []string) {
		for _, bucket := range utils.ListBuckets() {
			fmt.Println(bucket)
		}
	},
}

var s3ListObjects = &cobra.Command{
	Use:   "list-objects",
	Short: "List of objects",
	Run: func(cmd *cobra.Command, args []string) {
		if bucket != "" {
			for _, obj := range utils.ListObjects(bucket) {
				fmt.Println(obj)
			}
		} else {
			for _, bucket := range utils.ListBuckets() {
				for _, obj := range utils.ListObjects(bucket) {
					fmt.Println(obj)
				}
			}
		}
	},
}

func init() {
	s3Cmd.Flags().StringVarP(&accountID, "account-id", "a", "", "Account Id")
	s3Cmd.Flags().StringVarP(&bucket, "bucket", "b", "", "Bucket name")

	s3Cmd.AddCommand(s3ListBuckets)
	rootCmd.AddCommand(s3Cmd)
}
