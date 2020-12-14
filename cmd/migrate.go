package cmd

import (
	"github.com/spf13/cobra"
	"redis-migrator/config"
	"redis-migrator/migrator"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Runs redis-migrator to run migration",
	Long:  `Runs redis-migrator to run migration`,
	Run: func(cmd *cobra.Command, args []string) {
		runMigration()
	},
}

func init() {
	migrateCmd.PersistentFlags().StringVarP(&configFilePath, "config.file", "c", "config.yaml", "Location of configuration file to run migration.")
	rootCmd.AddCommand(migrateCmd)
}

func runMigration() {
	data := config.ParseConfig(configFilePath)
	migrator.MigrateRedisData(data)
}
