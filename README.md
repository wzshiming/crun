# Generate all possibilities based on regexp

* [English](./README.md)
* [简体中文](./README_cn.md)

## Usage

```
Usage of crun:
       crun [Options] [regexp]
    or crun "\d{3}"
    or crun "[0-9a-z]{2}"
    or crun "(root|admin) [0-9]{1}"

Options:
	-e # Execute the generated text
```

### Example

``` bash
# Generates a number of all possible combinations of 1 to 6 digits in length
crun "\d{1,6}"

# Violence aesthetics
crun "(root|admin):[0-9]{4,5}"

# Note: If the number is too large, super slow
```

### Download & Install
``` bash
go get -u -v github.com/wzshiming/crun/cmd/crun
```

## MIT License

Copyright © 2017-2018 wzshiming<[https://github.com/wzshiming](https://github.com/wzshiming)>.

MIT is open-sourced software licensed under the [MIT License](https://opensource.org/licenses/MIT).
