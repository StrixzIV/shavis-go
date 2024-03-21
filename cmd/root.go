/*
Copyright © 2024 Jirayu Kaewsing strixz.self@gmail.com, kernel137 kostamecev@pm.me
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "shavis-go",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.shavis-go.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolP("git-latest", "l", false, "Use a latest git commit hash to generate 8x5 image")
	rootCmd.PersistentFlags().StringP("git", "g", "", "Use a specified git commit hash to generate 8x5 image")
}

func run(cmd *cobra.Command, args []string) {

	git_hash, _ := cmd.Flags().GetString("git")
	use_latest, _ := cmd.Flags().GetBool("git-latest")

	if use_latest {

		executable, _ := os.Executable()

		repo, _ := git.PlainOpen(filepath.Dir(executable))
		ref, _ := repo.Head()

		git_hash = strings.Split(ref.String(), " ")[0]

	}

	fmt.Println(git_hash)

}
