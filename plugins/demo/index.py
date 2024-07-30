# coding:utf-8

import sys
import io
import os
import time
import re
import string
import subprocess
import pymongo
import json
from bson import ObjectId

class JSONEncoder(json.JSONEncoder):
    def default(self, o):
        if isinstance(o, ObjectId):
            return str(o)
        return json.JSONEncoder.default(self, o)


def getJson(data):
    return json.dumps(data,cls=JSONEncoder)


def getArgs():
    args = sys.argv[2:]
    # print(args)
    data = json.loads(args[0])
    # print(data)
    return data

def returnData(status, msg, data=None):
    return {'status': status, 'msg': msg, 'data': data}


def returnJson(code, msg, data=None, count=0):

    # if data == None:
    #     return {'status': status, 'msg': msg}
    # return {'status': status, 'msg': msg, 'data': data}

    if data == None:
        return getJson({'code': code, 'msg': msg,'count':count})
    return getJson({'code': code, 'msg': msg, 'data': data,'count':count})


def mgdbConn():
    client = pymongo.MongoClient('127.0.0.1', 27017)
    return client

def getTable(name = 'test'):
    mgdb_client = mgdbConn()
    db = mgdb_client["dztask"]
    collection = db[name]
    return collection

def getTableCount(name='test'):
    info = getTable(name)
    result = info.count_documents({})
    return result

def getTableList(name='test',where = {},page=1, size = 5):
    collection = getTable(name)

    start_index = (page - 1) * size
    end_index = page * size

    result = collection.find(where).skip(start_index).limit(size).sort({'_id':-1})

    d = []
    for document in result:
        d.append(document)
    return d

def mongodbList():
    args = getArgs()
    page = int(args['page'])
    limit = int(args['limit'])
    mlist = getTableList('test',{}, page, limit)
    count = getTableCount()
    return returnJson(0, 'ok', mlist, count)

def mongodbPush():
    args = getArgs()
    # print(args)
    extra = args['extra']
    return returnJson(0, 'ok:'+extra)

def pythonError():
    return returnJson(-1, '缺少命令')


if __name__ == "__main__":
    func = sys.argv[1]
    if func == 'list':
        print(mongodbList())
    elif func == 'push':
        print(mongodbPush())
    else:
        print(pythonError())