<p align="center">
  <h3 align="center">定制</h3>
</p>


### 安装脚本
```
curl --insecure -fsSL https://raw.githubusercontent.com/midoks/dztasks/master/script/install.sh?$(date +%s) | bash
```


### 卸载脚本
```
curl --insecure -fsSL https://raw.githubusercontent.com/midoks/dztasks/master/script/uninstall.sh?$(date +%s) | bash
```

### 用户密码
```
cat /opt/dztasks/custom/conf/app.conf
```

### 调式
```
curl --insecure -fsSL https://cdn.jsdelivr.net/gh/midoks/dztasks@latest/script/install.sh | bash

wget --no-check-certificate -O /tmp/dztasks.sh https://raw.githubusercontent.com/midoks/dztasks/master/script/install.sh?$(date +%s) && bash /tmp/dztasks.sh
```
