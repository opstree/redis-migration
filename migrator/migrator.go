package migrator

import (
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"redis-migrator/client"
	"redis-migrator/config"
)

// MigrateRedisData is the function to migrate keys from old to new redis
func MigrateRedisData(redConfig config.Configuration) {
	for _, database := range redConfig.Databases {
		logrus.Errorf("Executing migrator for database: %v", database)
		oldRedisClient, err := client.OldRedisClient(redConfig, database)
		newRedisClient, err := client.NewRedisClient(redConfig, database)
		if err != nil {
			logrus.Errorf("Error while connecting with redis %v", err)
		}
		keys, err := redis.Strings(oldRedisClient.Do("KEYS", "*"))
		if err != nil {
			logrus.Errorf("Error while listing redis keys %v", err)
		}
		for _, key := range keys {
			value, err := redis.String(oldRedisClient.Do("GET", key))
			if err != nil {
				logrus.Errorf("Not able to get the value for key %s: %v", key, err)
			}
			newRedisClient.Do("SET", key, value)
			logrus.Debugf("Migrated %s key with value: %v", key, value)
		}
	}
}
