
## Performance of Forward / Reverse MiMC

[MiMC](https://eprint.iacr.org/2016/492) is a exceptionally easy to implement [_Verifiable Delay Function_](https://eprint.iacr.org/2018/601.pdf).  Below we compare 4 language implementations -- C, Go, Python, and Node.js performance on a short 8192 step Forward and Reverse MiMC problem, the same as used in the STARK code.  These are the results:

| Implementation      | Forward MiMC | Reverse MiMC |
| --------------      | ------------ | ------------ |
| C - `mimc.c`        |       1.85ms |      122.7ms |
| Go - `mimc_test.go` |       6.11ms |      278.8ms |
| Python - `mimc.py`  |      13.40ms |    1,291.4ms |
| Node.js - `mimc.js` |     108.29ms |   13,910.1ms |

### C

After installing [GMP](https://gmplib.org/):
```
$ gcc mimc.c -o mimc -lgmp; ./mimc
forward-mimc: 2432 microseconds
reverse-mimc: 106364 microseconds
PASS
```

### Go

```
$ go test -run MiMC
forward-mimc: 6.110952ms
reverse-mimc: 278.815254ms
PASS
ok  	github.com/wolkdb/plasma/vdf	0.291s
```

### Python

```
$ python3.7 mimc.py  
forward-mimc 0.0134 sec
reverse-mimc: 1.2914 sec
PASS
```

### Node.js

```
$ npm install big-integer
$ node mimc.js
forward-mimc: 108.286ms
reverse-mimc: 13910.061ms
PASS
```
