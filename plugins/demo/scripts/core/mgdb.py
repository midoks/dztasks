# coding:utf-8

import pymongo

def conn():
    client = pymongo.MongoClient('127.0.0.1', 27017)
    return client


def getTable(name = 'test', dbname = 'demo'):
    mgdb_client = conn()
    db = mgdb_client[dbname]
    collection = db[name]
    return collection

def getTableCount(dbname='demo',name='test', where={}):
    doc = getTable(name,dbname)
    result = doc.count_documents(where)
    return result


def getTableList(dbname='demo',name='test',page=1, size = 5,where = {}):
    collection = getTable(name,dbname)

    start_index = (page - 1) * size
    end_index = page * size

    result = collection.find(where).skip(start_index).limit(size).sort({'_id':-1})

    d = []
    for document in result:
        d.append(document)
    return d


