# coding:utf-8

import sys
import io
import os
import time
import re
import string
import subprocess


mgdb_client = pymongo.MongoClient('127.0.0.1', 27017)
db = mgdb_client["dztask"]
collection = db["test"]


def mongodbList():
	return 'test'


if __name__ == "__main__":
    func = sys.argv[1]
    if func == 'list':
        print(mongodbList())
    else:
        print('error')