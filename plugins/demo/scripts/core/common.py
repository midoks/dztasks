import re
import os
import sys
import base64
import json

from random import Random
from bson import ObjectId

import binascii
import subprocess
import time

class JSONEncoder(json.JSONEncoder):
    def default(self, o):
        if isinstance(o, ObjectId):
            return str(o)
        return json.JSONEncoder.default(self, o)

def getJson(data):
    return json.dumps(data,cls=JSONEncoder)

def returnJson(code, msg, data=None, count=0):
    if data == None:
        return getJson({'code': code, 'msg': msg})
    return getJson({'code': code, 'msg': msg, 'data': data,'count':count})

def getValDef(data, k = 'a1',def_val = ''):
    if k in data:
        return data[k]
    else:
        return def_val

def getUpdatePos(file = ''):
    root_dir = '/tmp'

    if file == '':
        return 0
    dst_file = root_dir +'/'+file+'.txt'
    if not os.path.exists(dst_file):
        return 0
    content = readFile(dst_file).strip()
    if content == '':
        return 0
    return content

def setUpdatePos(file = '', pos = 0):
    # print(file,pos)
    root_dir = '/tmp'
    dst_file = root_dir +'/'+file+'.txt'
    writeFile(dst_file, str(pos))

def filterWord(search_key):
    search_key = search_key.split('(')[0]
    search_key = search_key.split('/')[0]
    search_key = search_key.split('・')[0]
    search_key = search_key.split('&')[0]
    
    return search_key

def execShell(cmdstring, cwd=None, timeout=None, shell=True):

    if shell:
        cmdstring_list = cmdstring
    else:
        cmdstring_list = shlex.split(cmdstring)
    if timeout:
        end_time = datetime.datetime.now() + datetime.timedelta(seconds=timeout)

    sub = subprocess.Popen(cmdstring_list, cwd=cwd, stdin=subprocess.PIPE,
                           shell=shell, bufsize=4096, stdout=subprocess.PIPE, stderr=subprocess.PIPE)

    while sub.poll() is None:
        time.sleep(0.1)
        if timeout:
            if end_time <= datetime.datetime.now():
                raise Exception("Timeout：%s" % cmdstring)

    if sys.version_info[0] == 2:
        return sub.communicate()

    data = sub.communicate()
    # python3 fix 返回byte数据
    if isinstance(data[0], bytes):
        t1 = str(data[0], encoding='utf-8')

    if isinstance(data[1], bytes):
        t2 = str(data[1], encoding='utf-8')
    return (t1, t2)

def getRunDir():
    return os.getcwd()


def getRootDir():
    return os.path.dirname(os.path.dirname(getRunDir()))


def getAppDir():
    return os.path.dirname(getRootDir())


def getTracebackInfo():
    import traceback
    errorMsg = traceback.format_exc()
    return errorMsg

def to_base64(s):
    return base64.b64encode(s.encode('utf-8'))
 
def image_to_base64(file_path):
    with open(file_path, "rb") as image_file:
        encoded_string = base64.b64encode(image_file.read())
    return encoded_string.decode("utf-8")

# base64转换成图片
def base64_to_image(abs_path, base64_encod_str):
    res = base64_encod_str.split(',')[1]
    img_b64decode = base64.b64decode(res)
    # open_image(img_b64decode)  # 打开图片验证下
    # 保存图片
    with open(abs_path, 'wb') as png:
        png.write(img_b64decode)

def base64_to_image_md5(base64_encod_str):
    res = base64_encod_str.split(',')[1]
    img_b64decode = base64.b64decode(res)
    # open_image(img_b64decode)  # 打开图片验证下
    with open('/tmp/t.jpg', 'wb') as png:
        png.write(img_b64decode)
    return file_md5('/tmp/t.jpg')


def open_image(img_b64decode):
    # 打开图片
    image = io.BytesIO(img_b64decode)
    # print(image)
    img = Image.open(image)
    img.show()

def toPinyin(txt):
    if txt == '':
        return ''
    import pohan
    from pohan.pinyin.pinyin import Style
    pinyin_list = pohan.pinyin.han2pinyin(txt, style=Style.NORMAL)

    s = ''
    for x in pinyin_list:
        s += x[0]
    return s


def get_date_int():
    # 取格式时间
    import time
    return int(time.time())

# print(get_date_int())
# sys.exit()


def md5(content):
    # 生成MD5
    try:
        import hashlib
        m = hashlib.md5()
        m.update(content.encode("utf-8"))
        return m.hexdigest()
    except Exception as ex:
        return False


def file_md5(file_path):
    import hashlib
    # 定义函数来计算文件的MD5值
    with open(file_path, 'rb') as f:
        data = f.read()
    return hashlib.md5(data).hexdigest()

def readFile(filename):
    # 读文件内容
    try:
        fp = open(filename, 'r')
        fBody = fp.read()
        fp.close()
        return fBody
    except Exception as e:
        print(e)
        return False


def make_sign(data):
    a2 = sorted(data.items(), key=lambda x: x[0])
    # print(a2)
    data = dict(a2)
    a = []
    for x in data:
        a.append(x + '=' + str(data[x]))

    s = '&'.join(a) + data['apikey']
    # print(s)
    s_md5 = md5(s)
    # print(a)
    # print(s)
    # print(s_md5.upper())
    return s_md5.upper()


def db_make_insert_sql(table, item):
    sql = "insert into " + table
    keyStr = '('
    valueStr = ' values('
    for i in item:
        # print i, item[i]
        keyStr += '`' + str(i) + '`,'
        valueStr += "\"" + str(item[i]) + "\","

    # name = MySQLdb.escape_string(name)
    keyStrLen = len(keyStr)
    keyStr = keyStr[0:keyStrLen - 1]
    keyStr += ') '

    valueStrLen = len(valueStr)
    valueStr = valueStr[0:valueStrLen - 1]
    valueStr += ') '

    sql += keyStr
    sql += valueStr
    return sql

def readFile(filename):
    # 读文件内容
    try:
        fp = open(filename, 'r')
        fBody = fp.read()
        fp.close()
        return fBody
    except Exception as e:
        # print(e)
        return False
        
def writeFile(filename, content, mode='w+'):
    # 写文件内容
    try:
        fp = open(filename, mode)
        fp.write(content)
        fp.close()
        return True
    except Exception as e:
        # print(e)
        return False

def db_make_update_sql(table, item, mid):
    sql = "update " + table + " set "
    keyStr = ''
    for i in item:
        keyStr += '`' + str(i) + '`=' + "\"" + str(item[i]) + "\","

    keyStrLen = len(keyStr)
    keyStr = keyStr[0:keyStrLen - 1]
    sql += keyStr
    sql += " where id = '" + str(mid) + "'"
    return sql
    
def getRandomString(length):
    # 取随机字符串
    rnd_str = ''
    chars = 'AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789'
    chrlen = len(chars) - 1
    random = Random()
    for i in range(length):
        rnd_str += chars[random.randint(0, chrlen)]
    return rnd_str


def encrypt_data(data, key=''):
    from Crypto.Cipher import AES

    import json
    if type(data) != str:
        data = json.dumps(data)

    # pad = lambda s: s + (16 - len(s) % 16) * chr(16 - len(s) % 16)
    # data = pad(data)
    vi = getRandomString(16)
    if key == '':
        key = 'xxxxxxxxxx'

    cipher = AES.new(key[:16].encode('utf8'), AES.MODE_CFB, vi.encode('utf8'))
    encryptedbytes = cipher.encrypt(data.encode('utf8'))
    encryptedbytes = binascii.hexlify(encryptedbytes)
    s = str(encryptedbytes, encoding="utf8")

    vi_en = binascii.hexlify(vi.encode('utf8'))
    vi_en = str(vi_en, encoding="utf8")

    # print('vi data:', vi)
    # print('vi hexlify data:', vi_en)
    # print('vi hexlify data len:', len(vi_en))
    # print('key data:', key)
    # print('encode data:', s)
    ret = s[0:16] + vi_en + s[16:]
    return ret


def decrypt_data(data, key=''):
    from Crypto.Cipher import AES

    if data == '':
        return ''

    if key == '':
        key = 'xxxxxxxxxx'

    slen = len(data)
    # print('data len:', slen)
    if slen < 48:
        content = data[0:slen - 32]
        iv = data[-32:]
    else:
        content = data[0:16] + data[48:]
        iv = data[16:48]

    # content = binascii.unhexlify(content)
    # print('content:', content, 'iv:', iv)
    iv = binascii.unhexlify(iv)
    content = binascii.unhexlify(content)

    # print(content, iv)

    cipher = AES.new(key[:16].encode('utf8'), AES.MODE_CFB, iv)
    text_decrypted = cipher.decrypt(content)

    text_decrypted = str(text_decrypted, encoding="utf8")
    # print(text_decrypted)
    try:
        import json
        return json.loads(text_decrypted)
    except Exception as e:
        pass
    return text_decrypted