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
(0) Setup: 136.873µs
(3) Converted round constants into a polynomial and low-degree extended it [4.338577ms => 4.513561ms]
(2a) Computational trace output (8192 steps) [10.153501ms => 10.379916ms]
(2b) Converted computational steps into a polynomial [44.85741ms => 55.26943ms]
(1) Powers of the higher-order root of unity [97.431887ms => 97.654749ms]
(2c) Extended it into p_evaluations [271.424934ms => 326.734087ms]
(5) Computed C(P, K) polynomial [110.764331ms => 451.715435ms]
(4a) Computed i_evaluations [150.901125ms]
(6a) Computed z_num_inv [114.10383ms => 566.439306ms]
(6b) Computed d_evaluations [44.098013ms => 610.566241ms]
(4b) Computed inv_z2_evaluations [149.57005ms => 627.294183ms]
(7) Computed b_evaluations, MerkleTree input [35.673054ms => 662.991602ms]
(8) Compute mtree[1] d494ec987dc1dfb66522849ed16f0a3b0bf57bc943e855ee6222ccc51436bb60 [42.509277ms => 705.543331ms]
(9) Computed random linear combination [36.515901ms => 742.083087ms]
(10) Merkelized l_evaluations, Setup Spot check positions [63.531043ms => 769.098099ms]
  (16384-2) Merkelize values [30.592064ms]
  (16384-1a) Calculate the set of x coordinates [41.294518ms]
  (16384-1b) Setup xsets, ysets [10.527359ms]
  (16384-1c) Computed x_polys [221.345699ms]
  (16384-3) Computed column [21.022536ms]
  (4096-2) Merkelize values [5.122202ms]
  (16384-5) Computed branches [5.388873ms]
  (4096-1a) Calculate the set of x coordinates [8.040565ms]
  (4096-1b) Setup xsets, ysets [2.612731ms]
  (4096-1c) Computed x_polys [52.973879ms]
  (4096-3) Computed column [3.023664ms]
  (1024-2) Merkelize values [2.315812ms]
  (4096-5) Computed branches [2.668178ms]
  (1024-1a) Calculate the set of x coordinates [2.617834ms]
  (1024-1b) Setup xsets, ysets [980.663µs]
  (1024-1c) Computed x_polys [15.551014ms]
  (1024-3) Computed column [2.043241ms]
  (256-2) Merkelize values [372.046µs]
  (1024-5) Computed branches [550.909µs]
  (256-1a) Calculate the set of x coordinates [677.658µs]
  (256-1b) Setup xsets, ysets [199.444µs]
  (256-1c) Computed x_polys [7.466321ms]
  (256-3) Computed column [956.16µs]
  (64-2) Merkelize values [91.49µs]
  (64-1a) Calculate the set of x coordinates [156.946µs]
  (64-1b) Setup xsets, ysets [46.534µs]
  (256-5) Computed branches [274.162µs]
  (64-1c) Computed x_polys [1.753876ms]
  (64-3) Computed column [145.251µs]
  (64-4) Computed prove_low_degree [12.319µs]
  (64-5) Computed branches [195.994µs]
  (256-4) Computed prove_low_degree [2.372231ms]
  (1024-4) Computed prove_low_degree [11.755987ms]
  (4096-4) Computed prove_low_degree [33.146026ms]
  (16384-4) Computed prove_low_degree [100.000932ms]
(11) Finished prove_low_degree [394.821483ms => 1.136939027s]
(12) Finalized branches [483.171µs => 1.137430673s]
STARK computed in 1.13744063s

STARK Proof size:  layer 0: 110322 bytes |  layer 1: 97122 bytes |  layer 2: 83922 bytes |  layer 3: 70722 bytes |  layer 4: 57480 bytes |  layer 5: 2120 bytes |
Approx proof length: 156238 bytes (branches), 421688 bytes (FRI proof), 577926 bytes (total)
------
Verifying degree <= 16 [1.100849ms]
(v2) Computed constants_mini_polynomial [302.785µs => 1.623201ms]
Verifying degree (0) <= 16384 [2.263923ms]
Verifying degree (1) <= 4096 [2.137585ms]
(v1) Computed MiMC output 95224774355499767951968048714566316597785297695903697235130434363122555476056 [6.654943ms => 6.721418ms]
Verifying degree (3) <= 256 [1.871259ms]
Verifying degree (4) <= 64 [1.677099ms]
Verifying degree (2) <= 1024 [2.005474ms]
FRI proof verified [9.340058ms]
(v3) verify_low_degree_proof [9.396215ms => 9.549263ms]
(v4) Verified 80 consistency checks [8.875194ms => 15.611056ms]
STARK verified in 15.641346ms
PASS
ok  	github.com/wolkdb/plasma/stark	1.176s
```

## Performance

Basic parallelization (go routines, WaitGroups) of the proof generation and verification have been added, with a first pass review of big.Int usage.  

### STARK Proof Generation
![Proof Parallelization](https://github.com/wolkdb/deepblockchains/blob/master/stark/flows/proof-goroutine-flow.png "Proof Parallelization")


![FRI Parallelization](https://github.com/wolkdb/deepblockchains/blob/master/stark/flows/fri-goroutine-flow.png "FRI Parallelization")

### STARK Verification

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
