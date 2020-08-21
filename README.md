# SimpleYggGen-Go

This program reimplements [SimpleYggGen](https://notabug.org/acetone/SimpleYggGen-Bash) in Go, importing the original Yggdrasil
code for generating keys and utilizing multiple CPU threads for mining.

### Installation

`% go get -u -v git.tdem.in/tdemin/syg_go`

### History

[SimpleYggGen](https://notabug.org/acetone/SimpleYggGen-Bash) is originally a
project by [@acetone](https://notabug.org/acetone), who wrote a Bash miner for
getting "magic" Yggdrasil addresses following a pattern. The main problem with
his implementation was that it ran grep and yggdrasil as separate processes,
making mining very slow. Even though @acetone later made a C++ implementation,
it still relied on running Yggdrasil as a separate process.

As of now (2020-08-12) @acetone reworked his C++ miner implementation, and
[SYG-C++](https://notabug.org/acetone/SimpleYggGen-CPP) is even more performant
than this program (making, like, 15% more iterations within the same time).

### Performance

Obviously far superior to the original SimpleYggGen.

With multiple threads it takes SimpleYggGen **a month** to run through a few
million cycles and find keys for `200::c84:77b0:f66d:b47e:64c7` (targeting
`::`). syg_go has found keys for `206:bcdb::ac47:4e3b:b97e:df4e` with the same
target in **27 minutes**, utilizing 8 threads on AMD Ryzen 1700X.

With 8 threads on Ryzen 1700X while searching for `::` this program reaches:

* 10 000 000 iterations in 2 minutes, 36 seconds
* 100 000 000 iterations in 25 minutes, 58 seconds
* 500 000 000 iterations in 2 hours, 10 minutes

This program contains some modded code from Yggdrasil that aims to improve
performance. If you prefer to use original Yggdrasil code, set `-original`
flag.

### Usage

```
% syg_go -help
Usage of syg_go:
  -highaddr
        high address mining mode (2xx::), excludes regex
  -iter uint
        per how many iterations to output status (default 100000)
  -original
        use original Yggdrasil code
  -regex string
        regex to match addresses against (default "::")
  -threads int
        how many threads to use for mining (default 16)
  -version
        display version
```

### License

See [LICENSE](LICENSE).
