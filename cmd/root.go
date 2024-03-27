/*
Copyright © 2024 Jirayu Kaewsing strixz.self@gmail.com, kernel137 kostamecev@pm.me
*/
package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var size_table string = `
N 	SHA256   	SHA1 (Git)
1 	8x8      	8x5
2 	16x16    	16x10
3 	32x32    	32x20
4 	64x64    	64x40
5 	128x128  	128x80
6 	256x256  	256x160
7 	512x512  	512x320
8 	1024x1024	1024x640
9 	2048x2048	2048x1280
10	4096x4096	4096x2560`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "shavis-go [SHA256 hash]",
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
	rootCmd.PersistentFlags().IntP("size", "s", 7, "Specified a size for output image"+size_table)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile != "" {

		current_dir, _ := os.Getwd()

		viper.SetConfigFile(cfgFile)
		fmt.Printf("Using config file: %s\n", path.Join(current_dir, viper.ConfigFileUsed()))

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
	theme := viper.GetString("config.theme")
	themes := viper.GetStringMapStringSlice("theme")

	theme_name, _ := cmd.Flags().GetString("theme")
	git_hash, _ := cmd.Flags().GetString("git")
	use_latest, _ := cmd.Flags().GetBool("git-latest")
	output_name, _ := cmd.Flags().GetString("output")
	file_name, _ := cmd.Flags().GetString("file")
	user_size, _ := cmd.Flags().GetInt("size")

	if theme_name != "" {
		theme = theme_name
	}

	if (output_name != "") && (!strings.HasSuffix(output_name, ".png")) {
		fmt.Println("Error: output name must ended with \".png\"")
		os.Exit(1)
	}

	if user_size < 1 || user_size > 10 {
		fmt.Println("Error: output size must be between 1 and 10", size_table)
		os.Exit(1)
	}

	if user_size != 7 {
		size_ptr := &size
		*size_ptr = user_size
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

			err = image_from_hash(input_hash, fmt.Sprintf("%s.png", input_hash), 8, 5, size, &themes, theme)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			return

		}

		err = image_from_hash(input_hash, output_name, 8, 5, size, &themes, theme)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

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

			err = image_from_hash(input_hash, fmt.Sprintf("%s.png", input_hash), 8, 5, size, &themes, theme)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			return

		}

		err = image_from_hash(input_hash, output_name, 8, 5, size, &themes, theme)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return

	}

	if file_name != "" {

		input_hash, err := filedata_to_hash(file_name)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if output_name == "" {

			err = image_from_hash(input_hash, fmt.Sprintf("%s.png", input_hash), 8, 8, size, &themes, theme)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			return

		}

		err = image_from_hash(input_hash, output_name, 8, 8, size, &themes, theme)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

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

		err = image_from_hash(input_hash, fmt.Sprintf("%s.png", input_hash), 8, 8, size, &themes, theme)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return

	}

	err = image_from_hash(input_hash, output_name, 8, 8, size, &themes, theme)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
