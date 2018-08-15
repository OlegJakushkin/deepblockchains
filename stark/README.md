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
(0) Setup: 182.571µs
(3) Converted round constants into a polynomial and low-degree extended it [4.4276ms => 4.666351ms]
(2a) Computational trace output (8192 steps) [10.389452ms => 10.700385ms]
(2b) Converted computational steps into a polynomial [46.674131ms => 57.400872ms]
(1) Powers of the higher-order root of unity [95.682927ms => 95.948311ms]
(2c) Extended it into p_evaluations [270.383988ms => 327.829209ms]
(5) Computed C(P, K) polynomial [101.310373ms => 430.617625ms]
(4a) Computed i_evaluations [130.134393ms]
(4b) Computed inv_z2_evaluations [116.376872ms => 574.446241ms]
(6a) Computed z_num_inv [95.80993ms => 670.283144ms]
(6b) Computed d_evaluations [24.213509ms => 694.804971ms]
(7) Computed b_evaluations, MerkleTree input [46.881766ms => 741.716265ms]
(8) Compute mtree[1] d494ec987dc1dfb66522849ed16f0a3b0bf57bc943e855ee6222ccc51436bb60 [28.393291ms => 770.132013ms]
(9) Computed random linear combination [45.518049ms => 815.66801ms]
(10) Merkelized l_evaluations, Setup Spot check positions [72.14833ms => 842.298166ms]
  (16384-2) Merkelize values [26.727857ms]
  (16384-1a) Calculate the set of x coordinates [39.512573ms]
  (16384-1b) Setup xsets, ysets [11.181406ms]
  (16384-1c) Computed x_polys [265.800292ms]
  (16384-3) Computed column [11.267897ms]
  (4096-2) Merkelize values [5.331177ms]
  (16384-5) Computed branches [5.767438ms]
  (4096-1a) Calculate the set of x coordinates [8.182335ms]
  (4096-1b) Setup xsets, ysets [3.21678ms]
  (4096-1c) Computed x_polys [61.048019ms]
  (4096-3) Computed column [2.577056ms]
  (1024-2) Merkelize values [1.313805ms]
  (4096-5) Computed branches [1.670755ms]
  (1024-1a) Calculate the set of x coordinates [2.211452ms]
  (1024-1b) Setup xsets, ysets [741.27µs]
  (1024-1c) Computed x_polys [19.82389ms]
  (1024-3) Computed column [1.377656ms]
  (256-2) Merkelize values [346.107µs]
  (1024-5) Computed branches [570.049µs]
  (256-1a) Calculate the set of x coordinates [685.173µs]
  (256-1b) Setup xsets, ysets [216.53µs]
  (256-1c) Computed x_polys [6.863071ms]
  (256-3) Computed column [780.8µs]
  (64-2) Merkelize values [89.412µs]
  (256-5) Computed branches [246.714µs]
  (64-1a) Calculate the set of x coordinates [143.897µs]
  (64-1b) Setup xsets, ysets [58.229µs]
  (64-1c) Computed x_polys [1.628562ms]
  (64-3) Computed column [183.513µs]
  (64-4) Computed prove_low_degree [10.174µs]
  (64-5) Computed branches [618.981µs]
  (256-4) Computed prove_low_degree [6.3084ms]
  (1024-4) Computed prove_low_degree [14.94238ms]
  (4096-4) Computed prove_low_degree [39.224187ms]
  (16384-4) Computed prove_low_degree [114.427395ms]
(11) Finished prove_low_degree [442.40318ms => 1.258112961s]
(12) Finalized branches [896.618µs => 1.259015913s]
STARK computed in 1.259029107s

STARK Proof size:  layer 0: 110322 bytes |  layer 1: 97122 bytes |  layer 2: 83922 bytes |  layer 3: 70722 bytes |  layer 4: 57480 bytes |  layer 5: 2120 bytes |
Approx proof length: 156238 bytes (branches), 421688 bytes (FRI proof), 577926 bytes (total)
------
Verifying degree <= 16 [1.04213ms]
(v2) Computed constants_mini_polynomial [305.206µs => 1.547463ms]
Verifying degree (2) <= 1024 [2.131157ms]
Verifying degree (3) <= 256 [1.976066ms]
Verifying degree (4) <= 64 [1.899876ms]
Verifying degree (1) <= 4096 [2.397724ms]
Verifying degree (0) <= 16384 [7.027816ms]
FRI proof verified [8.540764ms]
(v3) verify_low_degree_proof [8.590437ms => 8.735205ms]
(v1) Computed MiMC output 95224774355499767951968048714566316597785297695903697235130434363122555476056 [10.786366ms => 10.852428ms]
(v4) Verified 80 consistency checks [14.339083ms]
STARK verified in 14.36828ms
PASS
ok  	github.com/wolkdb/plasma/stark	1.296s
```

## Parallelization Performance

Basic parallelization (go routines, WaitGroups) of the proof generation and verification have been added, with a first pass review of big.Int usage.  

![Proof Parallelization](https://github.com/wolkdb/deepblockchains/blob/master/stark/flows/proof-goroutine-flow.png "Proof Parallelization")

![FRI Parallelization](https://github.com/wolkdb/deepblockchains/blob/master/stark/flows/fri-goroutine-flow.png "FRI Parallelization")

![Verification Parallelization](https://github.com/wolkdb/deepblockchains/blob/master/stark/flows/verify-goroutine-flow.png "Verification Parallelization")



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
