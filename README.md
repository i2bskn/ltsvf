ltsvf (LTSV filter)
===================

Install
-------

Download from [releases](https://github.com/i2bskn/ltsvf/releases) and stored in the `$PATH`.

or With go lang:

```
go get github.com/i2bskn/ltsvf
```

Usage
-----

```
$ cat test.txt | ltsvf
a:0 b:1
a:2 b:3
$ cat test.txt | ltsvf -f a:0,b:1
a:0 b:1
```

