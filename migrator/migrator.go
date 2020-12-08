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
		logrus.Debugf("Executing migrator for database: %v", database)
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
			keyType, err := redis.String(oldRedisClient.Do("TYPE", key))
			if err != nil {
				logrus.Errorf("Not able to get the key type %s: %v", key, err)
			}
			switch keyType {
			case "string":
				migrateStringKeys(oldRedisClient, newRedisClient, key)
			case "hash":
				migrateHashKeys(oldRedisClient, newRedisClient, key)
			}
		}
	}
}

func migrateHashKeys(oldClient redis.Conn, newClient redis.Conn, key string) {
	value, err := redis.StringMap(oldClient.Do("HGETALL", key))
	var data = []interface{}{key}
	for k, v := range value {
		data = append(data, k, v)
	}
	if err != nil {
		logrus.Errorf("Not able to get the value for key %s: %v", key, err)
	}
	newClient.Do("HMSET", data...)
	logrus.Debugf("Migrated %s key with value: %v", key, data)
}

func migrateStringKeys(oldClient redis.Conn, newClient redis.Conn, key string) {
	value, err := redis.String(oldClient.Do("GET", key))
	if err != nil {
		logrus.Errorf("Not able to get the value for key %s: %v", key, err)
	}
	newClient.Do("SET", key, value)
	logrus.Debugf("Migrated %s key with value: %v", key, value)
}
