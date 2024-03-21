/*
Copyright © 2024 Jirayu Kaewsing strixz.self@gmail.com, kernel137 kostamecev@pm.me
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "shavis-go",
	Short: "A Go implimentation of SHA256 or SHA1 hash sum visualization, either directly, file or git commit hash based on: https://github.com/kernel137/shavis",
	Long: `A Go implimentation of SHA256 or SHA1 hash sum visualization, either directly, file or git commit hash
based on https://github.com/kernel137/shavis original implimentation in Python
	`,
	Run: run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.shavis-go.yaml)")

	rootCmd.PersistentFlags().BoolP("git-latest", "l", false, "Use a latest git commit from current working directory hash to generate 8x5 image")
	rootCmd.PersistentFlags().StringP("git", "g", "", "Use a specified git commit hash to generate 8x5 image")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".shavis-go" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".shavis-go")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func run(cmd *cobra.Command, args []string) {

	var input_hash string

	size := viper.GetInt("config.size")
	git_hash, _ := cmd.Flags().GetString("git")
	use_latest, _ := cmd.Flags().GetBool("git-latest")

	if use_latest {

		current_dir, _ := os.Getwd()
		repo, err := git.PlainOpen(current_dir)

		if err != nil {
			fmt.Println("Error: git repository not found in current working directory")
			os.Exit(1)
		}

		ref, _ := repo.Head()

		input_hash = strings.Split(ref.String(), " ")[0]
		image_from_hash(input_hash, fmt.Sprintf("%s.png", input_hash), 8, 5, size, viper.GetStringSlice("theme.red"))

		return

	}

	if git_hash != "" {

		err := hash_check(git_hash, "SHA1")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		input_hash = git_hash
		image_from_hash(input_hash, fmt.Sprintf("%s.png", input_hash), 8, 5, size, viper.GetStringSlice("theme.red"))

		return

	}

	if len(args) == 0 {
		cmd.Help()
		os.Exit(0)
	}

	input_hash = args[0]
	image_from_hash(input_hash, fmt.Sprintf("%s.png", input_hash), 8, 8, size, viper.GetStringSlice("theme.red"))

}
