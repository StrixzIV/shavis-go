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
	Use:   "shavis [SHA256 hash]",
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
	rootCmd.PersistentFlags().StringP("file", "f", "", "Use a file's SHA256 checksum to generate 8x8 image")
	rootCmd.PersistentFlags().StringP("theme", "t", "", "Use a specified a color theme of generated image")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Specified a name for output image (Ex. result.png, output.png, your_name.png)")

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

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "error: config file not found in home directory")
		fmt.Fprintln(os.Stderr, "please put \".shavis-go.yaml\" in your home directory or use \"shavis [args] --config config_path\" to temporary set an image generation config")
		os.Exit(1)
	}

}

func run(cmd *cobra.Command, args []string) {

	var input_hash string

	size := viper.GetInt("config.size")
	theme := viper.GetStringMapStringSlice("theme")[viper.GetString("config.theme")]

	theme_name, _ := cmd.Flags().GetString("theme")
	git_hash, _ := cmd.Flags().GetString("git")
	use_latest, _ := cmd.Flags().GetBool("git-latest")
	output_name, _ := cmd.Flags().GetString("output")
	file_name, _ := cmd.Flags().GetString("file")

	if theme_name != "" {

		themes := viper.GetStringMapStringSlice("theme")

		theme_data, exists := themes[theme_name]

		if !exists {
			fmt.Printf("Error: theme \"%s\" is not defined in .shavis-go.yaml file\n", theme_name)
			os.Exit(1)
		}

		theme = theme_data

	}

	if (output_name != "") && (!strings.HasSuffix(output_name, ".png")) {
		fmt.Println("Error: output name must ended with \".png\"")
		os.Exit(1)
	}

	if use_latest {

		current_dir, _ := os.Getwd()
		repo, err := git.PlainOpen(current_dir)

		if err != nil {
			fmt.Println("Error: git repository not found in current working directory")
			os.Exit(1)
		}

		ref, _ := repo.Head()
		input_hash = strings.Split(ref.String(), " ")[0]

		if output_name == "" {
			image_from_hash(input_hash, fmt.Sprintf("%s.png", input_hash), 8, 5, size, theme)
			return
		}

		image_from_hash(input_hash, output_name, 8, 5, size, theme)
		return

	}

	if git_hash != "" {

		err := hash_check(git_hash, "SHA1")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		input_hash = git_hash

		if output_name == "" {
			image_from_hash(input_hash, fmt.Sprintf("%s.png", input_hash), 8, 5, size, theme)
			return
		}

		image_from_hash(input_hash, output_name, 8, 5, size, theme)
		return

	}

	if file_name != "" {

		input_hash, err := filedata_to_hash(file_name)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if output_name == "" {
			image_from_hash(input_hash, fmt.Sprintf("%s.png", input_hash), 8, 8, size, theme)
			return
		}

		image_from_hash(input_hash, output_name, 8, 8, size, theme)
		return

	}

	if len(args) == 0 {
		cmd.Help()
		os.Exit(0)
	}

	input_hash = args[0]
	err := hash_check(input_hash, "SHA256")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if output_name == "" {
		image_from_hash(input_hash, fmt.Sprintf("%s.png", input_hash), 8, 8, size, theme)
		return
	}

	image_from_hash(input_hash, output_name, 8, 8, size, theme)

}
