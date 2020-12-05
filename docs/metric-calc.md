# Metrics calculation

### Validator Monitoring dashboard:

* **Validator Availability header:** 

1) Node Status: Checking if port number 26657 on localhost is in use or not. If the heimdall rpc is active on this port, node status is marked as **UP** else **DOWN**.
2) Validator Status: If the validator has a non zero voting power in the result of http://localhost:26657/status it is marked as **Active** else **Inactive**.
3) Chain-id: Querying the rpc endpoint  http://localhost:26657/status. `result.node_info.network`

* **Valdiator Performance header**


1) Validator caught up: Querying the rpc endpoint http://localhost:26657/status. `result.sync_info.catching_up`. If `catching_up` is `false` it is marked as **Yes** else **No**.
2) Current block height validator : Querying rest endpoint http://localhost/1317/blocks/latest. `block_meta.header.height`
3) Block time difference: Commit time between the current and previous block is calculated and displayed here.
4) Latest block height network: An external heimdall rpc is queried for the block height.
5) Height difference: Difference between latest block height network and current block height validator is calculated and displayed.
6) Missed blocks: Validator's precommit is checked in http://localhost/1317/blocks/latest. If `block_meta.block.last_commit.precommits` does not contain the validator's precommit it's marked as missed.

* **Checkpoint header**

1) Latest checkpoint: Querying rest endpoint http://localhost:1317/checkpoints/count. `result.result`
2)  Checkpoint duration: Querying the rest endpoint http://localhost:1317/checkpoints/params. `result.checkpoint_buffer_time`
3)  Latest checkpoint start-end block: Querying the endpoint http://localhost:1317/checkpoints/latest. `result.start_block` and `result.end_block`.
4)  Last proposed checkpoint and No of checkpoints proposed: Querying the endpoint http://localhost:1317/checkpoints/<latest-checkpoint>. If `result.proposer` is user address then checkpoint number is displayed as proposed. Count of no of checkpoint proposed is incremented by 1.
5) Validator is part of block producer: Querying the rest endpoint http://localhost:/1317/bor/span/<current-span>. If user address is present in `result.selected_producers` it is marked as **Yes** else **No**.
6) Producer count: Querying the rest enpoint http://localhost:/1317/bor/span/<current-span>. Length of `result.selected_producers` array is displayed.
7) Span duration: Querying the rest endpoint http://localhost:1317/bor/params. `result.span_duration`

* **Validator connectivity header**

1) No of peers: querying the rpc endpoint http://localhost:26657/net_info?. `result.n_peers`
2) Peer addresses: querying the rpc endpoint http://localhost:26657/net_info?. `result.peers.remote_ip`

* **Validator details header**

1) Voting power: Querying the rest enpoint http://localhost:1317/staking/validator/<id>. `result.power`
2) Last proposed block height and time: Querying rest endpoint http://localhost/1317/blocks/latest. If `block.proposer_address` is user address then `block.height` and `block.time` are displayed.
3) ETH current balance: User's address is queried on mainnet Ethereum rpc.
4) Heimdall current balance: Querying the rest endpoint http://localhost:1317/bank/balances/<address>. `result.amount`
5) Max tx gas: Querying the rest endpoint http://localhost:1317/auth/params. `result.max_tx_gas`
6) Self stake and rewards: Querying validator share contract for id.

* **Proposals header**

Querying the rest endpoint http://localhost:1317/gov/proposals

### Bor dashboard:

1) Current height validator: Latest block is queried using bor rpc.
2) Current block proposer: Proposer of the latest bor block is being displayed.
3) Current span - Querying the rest endpoint http://localhost:1317/bor/latest-span. `result.span_id`.
4) Pending transactions: Querying the bor rpc for unconfirmed txs.
5) No of blocks proposed: If the user address is the proposer of the bor block then count is incremented by 1.
6) No of blocks signed - If the user address has signed on the latest bor block then count is incremented by 1.

