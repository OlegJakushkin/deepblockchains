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
Converted round constants into a polynomial and low-degree extended it [6.602174ms]
Computational trace output (8192 steps) [11.172972ms]
Converted computational steps into a polynomial [41.382899ms]
Extended it [258.199748ms]
Computed C(P, K) polynomial [93.196225ms]
Computed D polynomial [127.13986ms]
Computed d_evaluations [12.423474ms]
Computed inv_z2_evaluations [158.935997ms]
Computed b_evaluations, MerkleTree input [29.315647ms]
Computed hash root [41.135149ms] mtree[1] d494ec987dc1dfb66522849ed16f0a3b0bf57bc943e855ee6222ccc51436bb60
Computed random linear combination [23.022553ms]
Proving 65536 values are degree <= 16384 [341.929872ms] 2e3a9f04e27b5eadd61aee10e8c4523ca5af563962178195cb949ba23f3d2732
Proving 16384 values are degree <= 4096 [87.520293ms] ea9e14a3f0bf6b6da3ede789cf3436c2d02aa18d21e9e9a01e1bc17ea711a122
Proving 4096 values are degree <= 1024 [23.644541ms] 8d5ebb56cbe7a657638ac3801e63d94d627b3340fdbf135f61656c228f57ac0a
Proving 1024 values are degree <= 256 [8.809678ms] 81a40695570df55f6a19fdbad697d07b8ef7fdbc064afd3fa8751aa0ac6974f4
Proving 256 values are degree <= 64 [2.63234ms] 8ac23b191bb8ed3442a7e1d33fdae1028dc21af5c584550d6ae738e413a65309
Produced FRI proof
STARK computed in 1.290008954s

STARK Proof size:  layer 0: 110322 bytes |  layer 1: 97122 bytes |  layer 2: 83922 bytes |  layer 3: 70722 bytes |  layer 4: 57480 bytes |  layer 5: 2120 bytes |
Approx proof length: 156238 bytes (branches), 421688 bytes (FRI proof), 577926 bytes (total)
------
Verifying degree <= 16 [3.001955ms]
Verifying degree (0) <= 16384 [4.436105ms]
Verifying degree (1) <= 4096 [5.71256ms]
Verifying degree (2) <= 1024 [5.877444ms]
Verifying degree (4) <= 64 [6.438763ms]
Verifying degree (3) <= 256 [6.505624ms]
FRI proof verified [6.518201ms]
MIMC computed in 12.752269ms

Verified 80 consistency checks [10.308633ms]
STARK verified in 19.121431ms
PASS
ok	github.com/wolkdb/plasma/stark	1.334s
```

## Performance

Basic parallelization (go routines, WaitGroups) of the proof generation and verification have been added, with a first pass review of big.Int usage.  

### 2014 MacBook Pro 2.2 GHz Intel Core i7 (8 "logical" cores)

| `NUM_CORES` | STARK Proof Generation | STARK Verification |
| --- | --- | --- |
| 1 | 3.12s | 46.11ms |
| 2 | 2.01s | 24.70ms |
| 4 | 1.49s | 15.64ms |
| 8 | 1.42s | 16.58ms |
| 16 | 1.49s | 18.30ms |
| 32 | 1.48s | 20.18ms |
| [ethereum/research python](https://github.com/ethereum/research/tree/master/mimc_stark) | 3.78s | 52.10ms |

## Contributions

If you would like to contribute your updates to this (and run a workshop), please feel free to submit a pull request and email me.

## License info

GPLv3 [License](stark/LICENSE.md)
