# 九七式印字機

こっそり艦娘の機嫌を確認します。

## ビルド方法

```
$ sudo apt install libpcap-dev
$ go get github.com/google/gopacket
$ go build qnana.go
```

## 実行方法

```
$ ./qnana -h
Usage of ./qnana:
  -assembly_debug_log
        If true, the github.com/google/gopacket/tcpassembly library will
        log verbose debugging information (at least one line per packet)
  -assembly_memuse_log
        If true, the github.com/google/gopacket/tcpassembly library will
        log information regarding its memory use every once in a while.
  -f string
        Selects which packets will be processed.
  -i string
        Listen on interface. (default "eth0")
  -r string
        Read packets from file.
  -s int
        Snarf bytes of data from each packet. (default 1600)
$ sudo ./qnana -i eth0 2>/dev/null
```
