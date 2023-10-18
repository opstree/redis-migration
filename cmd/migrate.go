package cmd

import (
	"github.com/opstree/redis-migration/config"
	"github.com/opstree/redis-migration/migrator"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Runs redis-migrator to run migration",
	Long:  `Runs redis-migrator to run migration`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runMigration(); err != nil {
			panic(err)
		}
	},
}

func init() {
	migrateCmd.PersistentFlags().StringVarP(&configFilePath, "config.file", "c", "config.yaml", "Location of configuration file to run migration.")
	rootCmd.AddCommand(migrateCmd)
}

func runMigration() error {
	data := config.ParseConfig(configFilePath)
	return migrator.MigrateRedisData(data)
}
