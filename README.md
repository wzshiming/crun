# 根据正则生成字典


## Test

```
Usage of crun:
           crun <regexp>
        or crun "\d{3}"
        or crun "[0-9a-z]{2}"
        or crun "(root|admin) [0-9]{1}"
```


### 实例

``` bash
# 生成 1到6位长度的数字所有可能性组合
crun "\d{1,6}"

# 输出到 ditc.txt 文件
crun "\d{1,6}" > ditc.txt

# 暴力美学
crun "(root|admin):[0-9]{4,10}"

# !!!!! 注意如果量太大会超卡的
```

### 安装
``` bash
# 依赖 golang 和 git

# 设置环境变量 如果已经设置 请忽略
mkdir -p $HOME/gopath
export GOPATH=$HOME/gopath
export GOBIN=$GOPATH/bin
export PATH=$GOBIN:$PATH

# 下载&安装
go get -u -v github.com/wzshiming/crun/cmd/crun

```
