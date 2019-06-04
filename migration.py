import redis
from itertools import izip_longest
import sys
from multiprocessing import Pool
import signal


OLD = "SHARD_ENDPOINT_FOR_OLD_REDIS"
NEW = "SHARD_ENDPOINT_FOR_NEW_REDIS"

src = redis.Redis(host=OLD, port=6379, db=0,decode_responses=False)
dst = redis.Redis(host=NEW, port=6379, db=0,decode_responses=False)

def init_worker():
    signal.signal(signal.SIGINT, signal.SIG_IGN)

def dump_and_restore(key):
    ttl = src.ttl(str(key))
    if ttl == -1: # required as you can't create keys witn ttl -1, for no ttl , you can give 0 value.
        ttl = 0
    print "Dumping key: %s" % key
    value = src.dump(str(key))
    print "Restoring key: %s" % key
    try:
        dst.restore(str(key), ttl, value, replace=True)
    except redis.exceptions.ResponseError as e:
        print "Failed to restore key: %s" % key
        print "Exception: %s" % str(e)
        pass            
    return    

try:
    current = 0
    noOfWorkers = 10
    pool = Pool(noOfWorkers, init_worker)
    is_first = True
    while int(current) != 0 or is_first:
        is_first = False
        current,keys = src.execute_command('scan', int(current), 'count', 1000)
        pool.map(dump_and_restore,keys)

except KeyboardInterrupt:
    print "Caught KeyboardInterrupt, terminating workers"
    pool.terminate()
    pool.join()
