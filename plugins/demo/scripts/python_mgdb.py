import time
import pymongo

mgdb_client = pymongo.MongoClient('127.0.0.1', 27017)
db = mgdb_client["dztask"]
collection = db["test"]


now_t = int(time.time())
insert_data = {
    "time": now_t,
    "content":"测试"+str(now_t),
    "status":1
}

collection.insert_one(insert_data)
