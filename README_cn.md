# 根据正则生成所有可能性

[![Build Status](https://travis-ci.org/wzshiming/crun.svg?branch=master)](https://travis-ci.org/wzshiming/crun)
[![Go Report Card](https://goreportcard.com/badge/github.com/wzshiming/crun)](https://goreportcard.com/report/github.com/wzshiming/crun)
[![GoDoc](https://godoc.org/github.com/wzshiming/crun?status.svg)](https://godoc.org/github.com/wzshiming/crun)
[![GitHub license](https://img.shields.io/github/license/wzshiming/crun.svg)](https://github.com/wzshiming/crun/blob/master/LICENSE)
[![cover.run](https://cover.run/go/github.com/wzshiming/crun.svg?style=flat&tag=golang-1.10)](https://cover.run/go?tag=golang-1.10&repo=github.com%2Fwzshiming%2Fcrun)

- [English](./README.md)
- [简体中文](./README_cn.md)

## 用法

``` log
Usage of crun:
       crun [Options] [regexp]
    or crun "\d{3}"
    or crun "[0-9a-z]{2}"
    or crun "(root|admin) [0-9]{1}"

Options:
    -e # 执行产生的文本
```

## 实例

``` bash
# 生成 1到6位长度的数字所有可能性组合
crun "\d{1,6}"

# 暴力美学
crun "(root|admin):[0-9]{4,5}"

# 注意: 如果量太大会超卡的
```

## 下载 & 安装

``` bash
go get -u -v github.com/wzshiming/crun/cmd/crun
```

## 许可证

软包根据MIT License。有关完整的许可证文本，请参阅[LICENSE](https://github.com/wzshiming/crun/blob/master/LICENSE)。
