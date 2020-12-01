package cmd

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	configFilePath string
	outputFormat   string
	logLevel       string
	logFmt         string
	version        string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log.level", "", logrus.InfoLevel.String(), "redis-migrator logging level.")
	rootCmd.PersistentFlags().StringVarP(&logFmt, "log.format", "", "text", "redis-migrator log format.")
	flag.Parse()
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

var rootCmd = &cobra.Command{
	Use:   "redis-migrator",
	Short: "redis-migrator",
	Long:  `A tool for migrating redis database and its keys.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		parsedLevel, err := logrus.ParseLevel(logLevel)
		if err != nil {
			logrus.Errorf("log-level flag has invalid value %s", logLevel)
		} else {
			logrus.SetLevel(parsedLevel)
		}
		if logFmt == "json" {
			logrus.SetFormatter(&logrus.JSONFormatter{})
		} else {
			logrus.SetFormatter(&logrus.TextFormatter{
				DisableColors: true,
				FullTimestamp: true,
			})
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			logrus.Error(err)
		}
		os.Exit(1)
	},
}

// Execute the stuff
func Execute(VERSION string) {
	version = VERSION
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
