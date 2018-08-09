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
$ go test -run Stark
Computational trace output (8192 steps) [12.091336ms]
    fft setup [8.649762ms]
    inv_fft core [31.536506ms]
    invfft final [5.678538ms]
Converted computational steps into a polynomial [47.023251ms]
    fft setup [52.540103ms]
    reg_fft core [200.458932ms]
Extended it [253.051972ms]
Converted round constants into a polynomial and low-degree extended it [5.429781ms]
Computed C(P, K) polynomial [42.363488ms]
Computed D polynomial [129.364969ms]
Computed d_evaluations [26.4655ms]
Computed interpolant [38.006013ms]
Computed inv_z2_evaluations [180.55469ms]
Computed b_evaluations [21.630954ms] 
Computed hash root [65.766893ms] mtree[1] d494ec987dc1dfb66522849ed16f0a3b0bf57bc943e855ee6222ccc51436bb60
Computed random linear combination 2fdcce751344aad36fddf977e7c11cbd71cdd5b7f697b0bf9fce3beefcab2fb7 [109.350859ms]
Proving 65536 values are degree <= 16384 [400.974559ms] 2e3a9f04e27b5eadd61aee10e8c4523ca5af563962178195cb949ba23f3d2732
Proving 16384 values are degree <= 4096 [94.452123ms] ea9e14a3f0bf6b6da3ede789cf3436c2d02aa18d21e9e9a01e1bc17ea711a122
Proving 4096 values are degree <= 1024 [32.865383ms] 8d5ebb56cbe7a657638ac3801e63d94d627b3340fdbf135f61656c228f57ac0a
Proving 1024 values are degree <= 256 [11.112383ms] 81a40695570df55f6a19fdbad697d07b8ef7fdbc064afd3fa8751aa0ac6974f4
Proving 256 values are degree <= 64 [2.861338ms] 8ac23b191bb8ed3442a7e1d33fdae1028dc21af5c584550d6ae738e413a65309
Produced FRI proof
STARK computed in 1.474048858s

STARK Proof size:  layer 0: 110322 bytes |  layer 1: 97122 bytes |  layer 2: 83922 bytes |  layer 3: 70722 bytes |  layer 4: 57480 bytes |  layer 5: 2120 bytes | 
Approx proof length: 156238 bytes (branches), 421688 bytes (FRI proof), 577926 bytes (total)
------
Verifying degree <= 16
   [level 5 3.438503ms]
Verifying degree (3) <= 256 [5.853218ms]
Verifying degree (4) <= 64 [5.9343ms]
Verifying degree (2) <= 1024 [6.440503ms]
Verifying degree (0) <= 16384 [6.617105ms]
Verifying degree (1) <= 4096 [6.696188ms]
FRI proof verified [6.709595ms]
MIMC computed in 11.064569ms

Verified 80 consistency checks [9.905963ms]
STARK verified in 17.104753ms
PASS
ok	github.com/wolkdb/deepblockchains/stark	1.523s
```

## Contributions

If you would like to contribute your updates to this, please feel free to submit a pull request and email me.

## License info

GPLv3 [License](stark/LICENSE.md)
