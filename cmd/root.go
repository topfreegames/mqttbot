package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CfgFile is the path to the configuration file
var CfgFile string

//RootCmd is the root command for mqttbot CLI application
var RootCmd = &cobra.Command{
	Use:   "mqttbot",
	Short: "mqttbot is a bot spies on your mqtt topics",
	Long:  `Use mqttbot to spy on your mqtt topics, persist all its messages and act like a bot`,
}

//Execute runs RootCmd to initialize mqttbot CLI application
func Execute(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(
		&CfgFile, "config", "c", "./config/local.yaml",
		"config file (default is ./config/local.yaml)",
	)
	initConfig()
}

func initConfig() {
	if CfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(CfgFile)
	}
	viper.SetEnvPrefix("mqttbot")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigName(".mqttbot") // name of config file (without extension)
	viper.AddConfigPath("$HOME")    // adding home directory as first search path
	viper.AutomaticEnv()            // read in environment variables that match
	viper.SetDefault("logger.level", "DEBUG")

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		fmt.Println("Using config file: ", viper.ConfigFileUsed())
	} else {
		fmt.Printf("Error loading config: %v", err)
	}
}
