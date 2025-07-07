package cmd

import (
	"fmt"
	"os"
	// "os"
	// "github.com/aidahputri/training/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ConfigFile string
var rootCmd = &cobra.Command{
	Use: "main",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		viper.SetConfigFile(ConfigFile)
		err := viper.ReadInConfig()

		if err != nil {
			fmt.Println("Unable to read config file:", err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Server address is:", viper.GetString("server.listen_addr"))
	},
}

func Execute() error {
	rootCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "c", "config.yaml", "Config file to use")
	return rootCmd.Execute()
}