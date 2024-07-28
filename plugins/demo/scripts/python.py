import time

def writeFile(filename, content, mode='w+'):
    # 写文件内容
    try:
        fp = open(filename, mode)
        fp.write(content)
        fp.close()
        return True
    except Exception as e:
        print(e)
        return False


writeFile("/tmp/t.txt", str(time.time())+"\n")
print("hello")


