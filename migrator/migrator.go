package migrator

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"redis-migrator/client"
	"redis-migrator/config"
)

// MigrateRedisData is the function to migrate keys from old to new redis
func MigrateRedisData(redConfig config.Configuration) error {
	concurrentWorkers := max(1, redConfig.ConcurrentWorkers)
	if concurrentWorkers > len(redConfig.Databases) {
		concurrentWorkers = len(redConfig.Databases)
	}
	logrus.Infof("Migrating with %d concurrent processes", concurrentWorkers)
	errCh := make(chan error, len(redConfig.Databases))
	defer close(errCh)
	databaseCh := make(chan int, concurrentWorkers)
	// Schedule task for workers
	go func() {
		for _, database := range redConfig.Databases {
			databaseCh <- database
		}
		close(databaseCh)
	}()

	for i := 0; i < concurrentWorkers; i++ {
		go func() {
			var threadErr error
			defer func() { errCh <- threadErr }()

			for db := range databaseCh {
				logrus.Infof("Migrating database: %d", db)
				if err := migrateDB(redConfig, db); err != nil {
					threadErr = err
					return
				}
				logrus.Infof("Migrated database: %d", db)
			}
		}()
	}

	for i := 0; i < len(redConfig.Databases); i++ {
		if err := <-errCh; err != nil {
			return err
		}
	}

	return nil
}

func migrateDB(redConfig config.Configuration, db int) error {
	oldRedisClient, err := client.OldRedisClient(redConfig, db)
	if err != nil {
		return fmt.Errorf("[DB %d] Error while connecting with redis %v", db, err)
	}
	newRedisClient, err := client.NewRedisClient(redConfig, db)
	if err != nil {
		return fmt.Errorf("[DB %d] Error while connecting with redis %v", db, err)
	}
	keys, err := redis.Strings(oldRedisClient.Do("KEYS", "*"))
	if err != nil {
		return fmt.Errorf("[DB %d] Error while listing redis keys %v", db, err)
	}
	logrus.Infof("[DB %d] Migrating %d keys", db, len(keys))
	for i, key := range keys {
		keyType, err := redis.String(oldRedisClient.Do("TYPE", key))
		if err != nil {
			return fmt.Errorf("[DB %d] Not able to get the key type %s: %v", db, key, err)
		}
		switch keyType {
		case "string":
			migrateStringKeys(oldRedisClient, newRedisClient, key)
		case "hash":
			migrateHashKeys(oldRedisClient, newRedisClient, key)
		case "list":
			migrateListKeys(oldRedisClient, newRedisClient, key)
		default:
			return errors.New(fmt.Sprintf("[DB %d] key type is not supported: %s", db, keyType))
		}
		logrus.Debugf("[DB %d] Migrated %d/%d", db, i, len(keys))
	}
	return nil
}

func migrateListKeys(oldClient redis.Conn, newClient redis.Conn, key string) {
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
