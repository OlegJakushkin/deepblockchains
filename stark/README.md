# ZK-STARK: MIMC

Author: Sourabh Niyogi ([`sourabh@wolk.com`](mailto:sourabh@wolk.com))

### Background

ZK-STARKs are going to be widely applicable as a trust primitive for providing space/time efficient probabilistic proofs: valid state transitions in multiple layers of blockchains (e.g. including a STARK proof with each block), proving valid erasure coding, among other computational traces.   They do not require a trusted setup and are post-quantum resistant.

This is a Go-based port of Vitalik Buterin's [ZK-STARK/MIMC tutorial code](https://github.com/ethereum/research/tree/master/mimc_stark),
which concretely explores Ben-Sasson et al (2018)'s [Fast Reed-Solomon Interactive Oracle Proofs of Proximity](https://eccc.weizmann.ac.il/report/2017/134/):

* [STARKs, part 1: Proofs with Polynomials](https://vitalik.ca/general/2017/11/09/starks_part_1.html)
* [STARKs, part 2: Thank Goodness it's FRI-day](https://vitalik.ca/general/2017/11/22/starks_part_2.html)
* [STARKs, part 3: Into the Weeds](https://vitalik.ca/general/2018/07/21/starks_part_3.html)
* [Introduction zk SNARKs+STARKs - by Eli Ben-Sasson @ Technion Cyber and Computer Security Summer School](https://youtu.be/VUN35BC11Qw?t=20m7s)

This is not production-level code -- it is meant for blockchain engineers learning how ZK-STARKs (with MIMC VDFs) can be incorporated into next-generation blockchain designs.

### Usage

The following builds a STARK proof and verifies it:

```
$ go test -run Stark
Computational trace output: 95224774355499767951968048714566316597785297695903697235130434363122555476056
Converted computational steps into a polynomial and low-degree extended it 8192
Converted round constants into a polynomial and low-degree extended it 8192 [64]
Computed C(P, K) polynomial c_of_p_evaluations[65534]: 95527946821245685296728071053485177591639220188635432931376167003342370304418
Computed D polynomial
Computed hash root d494ec987dc1dfb66522849ed16f0a3b0bf57bc943e855ee6222ccc51436bb60
Computed random linear combination 2fdcce751344aad36fddf977e7c11cbd71cdd5b7f697b0bf9fce3beefcab2fb7
Computed spot checks [percentage] samples 240
Proving 65536 values are degree <= 16384
Proving 16384 values are degree <= 4096
Proving 4096 values are degree <= 1024
Proving 1024 values are degree <= 256
Proving 256 values are degree <= 64
Proving 64 values are degree <= 16
Produced FRI proof
STARK computed in 2.706264261s

STARK Proof size:  layer 0: 110322 bytes |  layer 1: 97122 bytes |  layer 2: 83922 bytes |  layer 3: 70722 bytes |  layer 4: 57480 bytes |  layer 5: 2120 bytes |
Approx proof length: 156238 bytes (branches), 421688 bytes (FRI proof), 577926 bytes (total)
MIMC computed in 8.924882ms

Verifying degree (0) <= 16384
Verifying degree (1) <= 4096
Verifying degree (2) <= 1024
Verifying degree (3) <= 256
Verifying degree (4) <= 64
Verifying degree <= 16
FRI proof verified (56 pts)
Verified 80 consistency checks
STARK verified in 43.175501ms
PASS
ok  	github.com/wolkdb/deepblockchains/stark	2.780s
```

## Contributions

If you would like to contribute your updates to this, please feel free to submit a pull request and email me.

## License info

GPLv3 [License](stark/LICENSE.md)
