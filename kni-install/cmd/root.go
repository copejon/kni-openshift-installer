/*
Copyright © 2020 Jonathan Cope jcope@redhat.com

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
	"os"
	"path"
	"path/filepath"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kni-install",
	RunE: rootCmdFunc,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// rootCmdFunc pre-configures system dependency locations prior to create() and destroy() calls
func rootCmdFunc(cmd *cobra.Command, _ []string) error {
	site = path.Base(siteRepo)
	siteBuildDir = filepath.Join(kniRoot, site, "final_manifests")
	ocpInstaller = filepath.Join(kniRoot, site, "requirements", "openshift-install")
	return cmd.Help()
}

var (
	isDryRun      bool
	isBareCluster bool
	kniRoot       string
	logLvl        string
	site          string
	siteRepo      string
	siteBuildDir  string
	ocpInstaller  string
)

func init() {
	cobra.OnInitialize(initConfig)

	userHome, _ := os.UserHomeDir()

	rootCmd.PersistentFlags().StringVar(&kniRoot, "kni-dir", filepath.Join(userHome, ".kni"), `(optional) Sets path to non-standard .kni path, useful for running the app outside of a containerized env.`)
	rootCmd.PersistentFlags().StringVar(&siteRepo, "repo", "", `git repo path containing site config files`)
	rootCmd.PersistentFlags().BoolVar(&isDryRun, "dry-run", false, `(optional) If true, prints, but does not execute OS commands.`)
	rootCmd.PersistentFlags().StringVar(&logLvl, "log-level", "info", `Set log level of detail. Accepted input is one of: ["info", "debug"]`)
	rootCmd.PersistentFlags().BoolVar(&isBareCluster, "bare-cluster", false, "when true, complete cluster deployment and stop, do no deploy workload.")
	_ = rootCmd.PersistentFlags().Parse(os.Args[1:])
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".kni-openshift-installer" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".kni-openshift-kni-install")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
