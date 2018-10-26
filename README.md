# Deep Blockchains

In this research, We describe how a deep blockchain architecture can address bottlenecks of a modern single-layer blockchain without sacrificing their core benefits of immutability, security, or trustlessness.

We describe in particular a 3-layer blockchain for provable storage and bandwidth in detail:
* Layer 1: Root Chain - which stores/publishes the SMT (Sparse Merkle Tree) root of {`Plasma Transactions`, and `Anchor Transactions`, `Plasma Accounts`} etc.

* Layer 2: The Plasma-Hybrid Chain - which manages (a) Plasma tokens redeemable for bandwidth and Storage (b) Layer 3 Chains Registrations and Permissions (c) Layer 3 Chains Latest State, all in SMTs.

* Layer 3: Child Chains - Any number of blockchains using storage and bandwidth of Layer 2, e.g those that package NoSQL / SQL transactions for typical database operations;

L2 and L3 utilize a Cloudstore abstraction to store and retrieve blocks and the chunks created by these blocks in Ethereum SWARM and multiple cloud computing providers. We further demonstrate how SMART(Spare Merkle Anchor Root Transactions) Proof can be used to provide provable data storage and on-chain provenance. We aim to demonstrate implementation results of high-throughput, low-latency L3 chains resting on economically secure L2 Plasma-Hybrid Chain, taken together which are capable of scaling for modern web applications.

## Resources [![Open Source Love](https://badges.frapsoft.com/os/v2/open-source.svg?v=103)](https://github.com/wolkdb/go-plasma)

[[Paper](https://github.com/wolkdb/deepblockchains/blob/master/Deep_Blockchains.pdf)] Deep paper describing the scalable multilayer blockchain architecture in details.    

#### Layer 1: RootChain Contract
[[Solidity](https://github.com/wolkdb/deepblockchains/tree/master/Plasmacas)] A feature-complete smart contract, including the Plasma APIs documentation

#### Layer 2: Plasma Chain
[[GoLang](https://github.com/wolkdb/go-plasma)] Wolk's code-complete Plasma-Hybrid implementation

#### Layer 3: SQL + NoSQL Chain
Email services@wolk.com for early access.

#### Cloudstore: Wolk's Decentralized Backend + Erasure Encoding
Email services@wolk.com for early access.

#### SMT Implementation
* [[Solidity](https://github.com/wolkdb/deepblockchains/blob/master/Plasmacash/contracts/RootChain/Libraries/SparseMerkle.sol)] Reference implementation for [ethresearch post](https://ethresear.ch/t/plasma-cash-with-sparse-merkle-trees-bloom-filters-and-probabilistic-transfers/2006) published here
* [[GoLang](https://github.com/wolkdb/deepblockchains/tree/master/smt)] SMT implementation in Go, compatible with Wolk Cloudstore

#### Others

* [[Stark](https://github.com/wolkdb/deepblockchains/tree/master/stark)] ZK-STARK implementation in Go
* [[MiMC](https://github.com/wolkdb/deepblockchains/tree/master/vdf/mimc)] Forward/Reverse MiMC Verifiable Delay Function benchmark in {C,GO, Python, Node.js}
