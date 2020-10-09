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
	"os"

	"github.com/merkio/dev-tools/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var accountID string
var bucket string
var from string
var to string
var force bool

// s3Cmd represents the s3 command
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Operations with s3 storage",
}

var listBuckets = &cobra.Command{
	Use:   "list-buckets",
	Short: "List of buckets",
	Run: func(cmd *cobra.Command, args []string) {
		for _, bucket := range utils.ListBuckets() {
			fmt.Println(bucket)
		}
	},
}

var listObjects = &cobra.Command{
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

var deleteBucket = &cobra.Command{
	Use:   "rm-bucket",
	Short: "Remove bucket",
	Run: func(cmd *cobra.Command, args []string) {
		if bucket == "" {
			fmt.Println("Bucket name is required!")
			os.Exit(1)
		}
		utils.DeleteBucket(bucket, force)
	},
}

var backupToDisk = &cobra.Command{
	Use:   "backup-to-disk",
	Short: "Backup files to the local disk",
	Run: func(cmd *cobra.Command, args []string) {
		if to == "" {
			fmt.Println("Bucket name was not provided, using '~/bkms-backup'")
		} else {
			viper.Set("BackupDir", to)
		}

		for _, bucket := range utils.ListBuckets() {
			utils.DownloadBucket(bucket, viper.GetString("BackupDir"))
		}
	},
}

var restoreFromDisk = &cobra.Command{
	Use:   "restore-from-disk",
	Short: "Restore files from the local disk",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	deleteBucket.Flags().StringVarP(&bucket, "bucket", "b", "", "Bucket name")
	force = *deleteBucket.Flags().BoolP("force", "f", false, "Force delete, if bucket contains data")

	listObjects.Flags().StringVarP(&accountID, "account-id", "a", "", "Account Id")
	listObjects.Flags().StringVarP(&bucket, "bucket", "b", "", "Bucket name")

	backupToDisk.Flags().StringVarP(&accountID, "account-id", "a", "", "Account Id")
	backupToDisk.Flags().StringVarP(&to, "to", "t", "", "Absolute local path")

	restoreFromDisk.Flags().StringVarP(&accountID, "account-id", "a", "", "Account Id")
	restoreFromDisk.Flags().StringVarP(&from, "from", "f", "", "Absolute local path")

	s3Cmd.AddCommand(listBuckets)
	s3Cmd.AddCommand(listObjects)
	s3Cmd.AddCommand(backupToDisk)
	s3Cmd.AddCommand(restoreFromDisk)
	rootCmd.AddCommand(s3Cmd)
}
