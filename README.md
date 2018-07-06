# Deep Blockchains

In our research, we describe how a deep blockchain architecture can address bottlenecks of a modern single-layer blockchain without sacrificing their core benefits of immutability, security, or trustlessness. Fundamentally, this is achievable with higher layer blockchains submitting block transactions containing summaries of the higher layer’s block for inclusion in blocks of the lower layer blockchain, with all higher layer blockchains supervenient on lower level blockchain features.

We describe in particular a 3-layer blockchain for provable storage and bandwidth in detail:
* Layer 1 is MainNet, which stores both (a) registered Layer 3 blockchain roots and (b) Layer 2 Block Merkle roots;

* Layer 2 is a Plasma Cash chain, storing (a) Plasma tokens redeemable for bandwidth and (b) Layer 3 Block hashes;

* Layer 3 blockchains are any number of blockchains using storage and bandwidth of Layer 2, e.g those that package NoSQL / SQL transactions for typical database operations;

Layers 2 and 3 utilize a Cloudstore abstraction to store and retrieve blocks and the chunks created by these blocks in Ethereum SWARM and multiple cloud computing providers. We demonstrate repeated use of Sparse Merkle Trees and show how this construct can be used in our core deep blockchain to provide provable data storage with Deep Merkle Proofs. We aim to demonstrate implementation results of high-throughput, low-latency layer 3 blockchains resting on  economically secure Layer 2 Plasma Cash blockchains, taken together which are fundamentally capable of scaling for modern web applications.

[Download DeepBlockchains paper](https://github.com/wolkdb/deepblockchains/blob/master/Deep_Blockchains.pdf)

## Layer 1: Root Chain Plasma Contract

A feature complete smart contract is available at [PlasmaCash](https://github.com/wolkdb/deepblockchains/tree/master/Plasmacash) directory

## Layer 2: Plasma Cash - POA

Coming soon.  Email services@wolk.com for early access.

## Layer 3: SQL + NoSQL Chain

Coming soon. Email services@wolk.com for early access.
