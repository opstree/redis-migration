#! /usr/bin/python3
'''
This utility aims to help anyone Wile migrating redis.
It is simple and convinent and compatible.

Author:- OpsTree Solutions
'''
import logging
from multiprocessing import  Process, Pool
import argparse
import signal
import redis
import tqdm
import config_with_yaml as config

parser = argparse.ArgumentParser(description='Redis Migration')
parser.add_argument('--version', help="Version of this application",\
                     action='version', version='Redis-Migration 1.1')
parser.add_argument('--log_level', help="Default is Warning, you can choose INFO,ERROR",\
                     action="store", type=str, default="WARNING")
parser.add_argument('--log_file', help="If you want to save log in file , just pass file path",\
                     action="store", type=str, default=None)
parser.add_argument('--log_json', help="Use this when you want log in json format",\
                     action="store_true", default=False)


class DatabaseError(Exception):
    '''
    To Generate own Exception for wrong entry in selected_database
    '''
    pass

def init_worker():
    '''
    To capture Keyboard Interrupt
    '''
    signal.signal(signal.SIGINT, signal.SIG_IGN)


def fetch_connection_object(yaml_config, database_index):
    '''
    To make connection object for given database.
    It return source_redis and destination_redis
    '''
    logging.info("Started migration for Selected DB %s", database_index)
    source_redis = redis.Redis(\
                    host=yaml_config.getProperty("old_redis.ip"), \
                    port=yaml_config.getPropertyWithDefault("old_redis.port", 6379), \
                    password=yaml_config.getPropertyWithDefault("NEW_REDIS.password", None), \
                    db=database_index, \
                    decode_responses=False)
    destination_redis = redis.Redis(\
                    host=yaml_config.getProperty("new_redis.ip"), \
                    port=yaml_config.getPropertyWithDefault("new_redis.port", 6379), \
                    password=yaml_config.getPropertyWithDefault("new_redis.password", None), \
                    db=database_index, \
                    decode_responses=False)
    return source_redis, destination_redis

def dump_and_restore_keys(key):
    """
    This function based recieve a key , then dump it and restore it to new redis
    TTL (Time to live )required as you can't create keys witn ttl -1
    for no ttl , you can give 0 value.
    Used str.decode('utf-8').
    It properly converts byte into string so that none type error can be eliminated
    It also makes utility more compatible.
    """
    ttl = SOURCE_REDIS.ttl(key.decode('utf-8'))
    if ttl == -1:
        ttl = 0
    logging.debug("Dumping key: %s", key.decode('utf-8'))
    value = SOURCE_REDIS.dump(key.decode('utf-8'))
    logging.debug("Restoring key: %s", key.decode('utf-8'))
    try:
        DESTINATION_REDIS.restore(key, ttl, value, replace=True)
    except redis.exceptions.ResponseError as error:
        logging.error("Failed to restore key: %s", key)
        logging.error("Exception: %s", str(error))

def start_migration(yaml_config, database_index):
    '''
    This function will call connection and retrieve SOURCE_REDIS, DESTINATION_REDIS,
    Then Based on no. of workers a Multiprocessing Pool will be used to make Program efficient.
    '''
    try:
        global SOURCE_REDIS, DESTINATION_REDIS
        SOURCE_REDIS, DESTINATION_REDIS = fetch_connection_object(yaml_config, database_index)
        current = 0
        no_of_workers = yaml_config.getPropertyWithDefault("no_of_workers", 6)
        pool = Pool(no_of_workers, init_worker)
        is_first = True
        while int(current) != 0 or is_first:
            is_first = False
            try:
                current, keys = SOURCE_REDIS.execute_command('scan', int(current), 'count', 2000)
                logging.info(" scan is executed and keys are fetched")
                logging.info(len(keys))
            except Exception as error:
                logging.error("Connection Error")
                logging.error(error)
                raise
            try:
                for _ in tqdm.tqdm(pool.map(dump_and_restore_keys, keys), total=len(keys)):
                    pass
            except Exception as error:
                logging.error("Error while pool mapping i.e when  workers were assigning keys ")
                logging.error(error)
                raise

        logging.info("Successfully Completed for selected database")

    except KeyboardInterrupt:
        logging.error("Caught KeyboardInterrupt, terminating workers")
        pool.terminate()
        pool.join()

def database_to_migrate():
    '''
    As redis database has logical index from 0-15
    If all is passed in config file then 0-15 is selected
    Or user can enter selected index in form of list
    '''
    database_index = yaml_config.getProperty("old_redis.select_database")
    if database_index == "all":
        return range(0, 16)
    elif isinstance(database_index) is list:
        return database_index
    elif isinstance(database_index) is int:
        single_index = [database_index]
        return single_index
    else:
        logging.error("Select Proper Database")
        raise DatabaseError("Entered Wrong Database")

'''
This function will be imported and called in main.py
Log Format Function it will return format of log
'''
def log_type(log_json):
    """
    If user passed log_json then log format will be json
    Default is simple log
    """
    simple_log = '%(asctime)s - %(filename)s - %(levelname)s - %(message)s'
    json = "{ \"Time\": \"%(asctime)s\" , \"Log Level\": \"%(levelname)s\",\"FILENAME\": \"%(filename)s\",\"Line Number\": \"%(lineno)d\",\"Message\": \"%(message)s\"}"
    format_log = simple_log
    if bool(log_json):
        format_log = json
    return format_log


if __name__ == "__main__":
    args = parser.parse_args()
    yaml_config = config.load("config.yaml")
    logging.basicConfig(filename=args.log_file, \
                        level=args.log_level, \
                        format=log_type(args.log_json))
    for each_database in database_to_migrate():
        print("Selected Db %s" % each_database)
        '''
        Process is used here to improve program efficiency
        Process is used for less task whereas Pool is used for large no.of task.
        '''
        createprocess = Process(target=start_migration, args=(yaml_config, each_database,))
        createprocess.start()
        createprocess.join()
