package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

const ext = "yml"

var (
	cfgCredHome, Version string
	debug                bool
)

var rootCmd = &cobra.Command{
	Use:   "dip",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		if debug {
			log.SetReportCaller(true)
			log.SetLevel(log.DebugLevel)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func viperBase(filename string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	viper.SetConfigName(filename)
	viper.SetConfigType(ext)
	viper.AddConfigPath(filepath.Join(home, ".dip"))

	if cfgCredHome != "" {
		viper.AddConfigPath(cfgCredHome)
	}

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("fatal error config file: %v", err)
	}
	return nil
}

func credsValue(key string) (string, error) {
	if err := viperBase("creds"); err != nil {
		return "", err
	}
	value := viper.GetString(key)
	if value == "" {
		return "", fmt.Errorf("no "+key+" found. Check whether the '"+key+"' variable is populated in '%s'", viper.ConfigFileUsed())
	}
	return value, nil
}

func slackToken() (string, error) {
	return credsValue("slack_token")
}

func slackChannelID() (string, error) {
	return credsValue("slack_channel_id")
}

func imagesToBeValidated() (map[string]interface{}, error) {
	if err := viperBase("config"); err != nil {
		return nil, err
	}
	images := viper.GetStringMap("dip_images")
	if len(images) == 0 {
		return nil, fmt.Errorf("no images found. Check whether the 'dip_images' variable is populated in '%s'", viper.ConfigFileUsed())
	}
	return images, nil
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgCredHome, "configCredHome", "", "config and cred file home directory (default is $HOME/.dip)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debugging mode")
}

func initConfig() {
	enableDebug()
}

func enableDebug() {
	log.SetReportCaller(true)
	if debug {
		log.SetLevel(log.DebugLevel)

		// Added to be able to debug viper (used to read the config file)
		// Viper is using a different logger
		jww.SetLogThreshold(jww.LevelTrace)
		jww.SetStdoutThreshold(jww.LevelDebug)
	}
}
