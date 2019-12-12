# Generate matching strings based on regular expressions

[![Build Status](https://travis-ci.org/wzshiming/crun.svg?branch=master)](https://travis-ci.org/wzshiming/crun)
[![Go Report Card](https://goreportcard.com/badge/github.com/wzshiming/crun)](https://goreportcard.com/report/github.com/wzshiming/crun)
[![GoDoc](https://godoc.org/github.com/wzshiming/crun?status.svg)](https://godoc.org/github.com/wzshiming/crun)
[![GitHub license](https://img.shields.io/github/license/wzshiming/crun.svg)](https://github.com/wzshiming/crun/blob/master/LICENSE)
[![gocover.io](https://gocover.io/_badge/github.com/wzshiming/crun)](https://gocover.io/github.com/wzshiming/crun)

- [English](https://github.com/wzshiming/crun/blob/master/README.md)
- [简体中文](https://github.com/wzshiming/crun/blob/master/README_cn.md)

## Example

``` bash
# Generates a number of all possible combinations of 1 to 6 digits in length
> crun "\d{1,6}"

# Generate random 5 possibilities
> crun -r -l 5 "(root|admin):[0-9]{4,5}"
```

## Download & Install

``` bash
go get -u -v github.com/wzshiming/crun/cmd/crun
```

## License

Pouch is licensed under the MIT License. See [LICENSE](https://github.com/wzshiming/crun/blob/master/LICENSE) for the full license text.
