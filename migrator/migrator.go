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
		if err != nil {
			logrus.Errorf("Error while connecting with redis %v", err)
		}
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
			case "list":
				migarteListKeys(oldRedisClient, newRedisClient, key)
			}
		}
	}
}

func migarteListKeys(oldClient redis.Conn, newClient redis.Conn, key string) {
	value, err := redis.Strings(oldClient.Do("LPOP", key))
	if err != nil {
		logrus.Errorf("Not able to get the value for key %s: %v", key, err)
	}
	var data = []interface{}{key}
	for _, v := range value {
		data = append(data, v)
	}
	_, err = newClient.Do("LPUSH", data...)
	if err != nil {
		logrus.Errorf("Error while pushing list keys %v", err)
	}
	logrus.Debugf("Migrated %s key with value: %v", key, data)
}

func migrateHashKeys(oldClient redis.Conn, newClient redis.Conn, key string) {
	value, err := redis.StringMap(oldClient.Do("HGETALL", key))
	if err != nil {
		logrus.Errorf("Not able to get the value for key %s: %v", key, err)
	}
	var data = []interface{}{key}
	for k, v := range value {
		data = append(data, k, v)
	}
	_, err = newClient.Do("HMSET", data...)
	if err != nil {
		logrus.Errorf("Error while pushing list keys %v", err)
	}
	logrus.Debugf("Migrated %s key with value: %v", key, data)
}

func migrateStringKeys(oldClient redis.Conn, newClient redis.Conn, key string) {
	value, err := redis.String(oldClient.Do("GET", key))
	if err != nil {
		logrus.Errorf("Not able to get the value for key %s: %v", key, err)
	}
	_, err = newClient.Do("SET", key, value)
	if err != nil {
		logrus.Errorf("Error while pushing list keys %v", err)
	}
	logrus.Debugf("Migrated %s key with value: %v", key, value)
}
