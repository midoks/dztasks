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


core_dir = os.getcwd()
sys.path.append(core_dir + "/scripts/core")

import mgdb
import common

collection = mgdb.getTable('test','dztask')

def getArgs():
    args = sys.argv[2:]
    data = json.loads(args[0])
    return data

def mongodbList():
    args = getArgs()

    page = int(args['page'])
    limit = int(args['limit'])

    mgdb_where = {}
    if 'args' in args:
        post_args = args['args']
        post_args = json.loads(post_args)
        if 'status' in post_args and post_args['status'] != '':
            mgdb_where['status'] = int(post_args['status'])
        if 'zd' in post_args and post_args['zd'] != '' and 'key' in post_args and post_args['key'] != '':
            zd = post_args['zd']
            key = post_args['key']
            if zd == 'content':
                mgdb_where[zd] =  {'$regex':key}
            elif zd == '_id':
                mgdb_where[zd] = ObjectId(key)
            else:
                mgdb_where[zd] = key

        if 'kstime' in post_args and post_args['kstime'] != '' and 'jstime' in post_args and post_args['jstime'] != '':
            kstime = post_args['kstime']
            jstime = post_args['jstime']

            kstime = int(time.mktime(time.strptime(kstime, "%Y-%m-%d %H:%M:%S")))
            jstime = int(time.mktime(time.strptime(jstime, "%Y-%m-%d %H:%M:%S")))

            # print(kstime,jstime)
            mgdb_where['time'] = {
                '$gt': kstime,
                '$lt': jstime
            }

    # print(mgdb_where)
    start_index = (page - 1) * limit
    result = collection.find(mgdb_where).skip(start_index).limit(limit).sort({'time':-1})

    mlist = []
    for row in result:
        mlist.append(row)

    count = collection.count_documents(mgdb_where)
    # print(count,mgdb_where)
    return common.returnJson(0, 'ok', mlist, count)

def mongodbPush():
    args = getArgs()
    # print(args)
    extra = args['extra']
    return common.returnJson(0, 'ok:'+extra)

def mongodbDelete():
    args = getArgs()
    extra = args['extra']

    # print({"_id":ObjectId(extra)})
    r = collection.delete_one({"_id":ObjectId(extra)})
    # print(r)
    if r:
        return common.returnJson(0, '删除成功!')
    return common.returnJson(0, 'ok:'+extra)

def pythonError():
    return common.returnJson(-1, '缺少命令')


if __name__ == "__main__":
    func = sys.argv[1]
    if func == 'list':
        print(mongodbList())
    elif func == 'push':
        print(mongodbPush())
    elif func == 'delete':
        print(mongodbDelete())
    else:
        print(pythonError())