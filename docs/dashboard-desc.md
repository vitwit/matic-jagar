### 1. Validator monitoring metrics (Heimdall)
The following list of metrics are displayed in this dashboard:

- Validator Details :  Displays the details of a validator like moniker, validator signer address and hex address.
- Node Status :  Displays whether the node is running or not in the form of **UP** or **DOWN**.
- Validator Status :  Displays the validator health. Shows **Voting** if the validator is in active state or else **Jailed**.
- Validator Caught Up : Displays whether the validator node is in sync with the network or not.
- Block Time Difference : Displays the time difference between the previous and current block.
- Current Block Height - Validator : Displays the latest block height committed by the validator.
- Latest Block Height - Network :  Displays the latest block height of a network.
- Height Difference : Displays the difference between heights of validator current block height and network latest block height.
- Missed Blocks : A graphical display of missed blocks.
- Last Missed Block Range : Displays the continuous missed blocks range based on the threshold provided in the `config.toml`
- Blocks Missed In last 48h : Displays the count of blocks missed by the validator in last 48 hours.
- Unconfirmed Txns : Displays the number of unconfirmed transactions on that node.
- Latest Checkpoint : Displays the height of the latest check point.
- No.of Peers : Displays the total number of peers connected to the validator.
- Peer Address : Displays the ip addresses of connected peers.
- Latency : Displays the latency of connected peers and the validator.
- Validator Fee : Displays the commission rate of the validator.
- Voting Power : Displays the voting power of the validator.
- Max Tx Gas : Displays the max transaction gas.
- Rewards : Displays the rewards of your validator.
- Last proposed Block Height : Displays height of the last block proposed by the validator.
- Last Proposed Block Time : Displays the time of the last block proposed by the validator.
- Heimdall Current Balance : Displays the account balance of the validator.
- ETH Current Balance : Displays the current ETH balance of your signer address.
- Self Stake : Displays the amount of self stake done by the signer address.
- Voting Period Proposals : Displays the list of the proposals which are currently in voting period.
- Deposit Period Proposals : Displays the list of the proposals which are currently in deposit period.
- Completed Proposals : Displays the list of the proposals which are completed with their status as passed or rejected.


**Note:** The above mentioned metrics will be calculated and displayed according to the validator address which will be configured in config.toml.

### 2. Bor
The following list of metrics are displayed in this dashboard:

- Current Block Height - validator : Displays the current height of bor on validator.
- Current Block Height - network : Displays the  height of bor on network. 
- Block Height Difference : Displays the difference between heights of validator current block height and network latest block height
- Current Span : Displays the current span
- Pending Transactions : Displays the number of unconfirmed transactions on that node
- Current Block Proposer : Displays the signer address of block producer.  
- No.of Blocks Proposed : Displays the count of blocks proposed by validator. 
- No.of Blocks Signed : Displays the count of blocks signed by validator.
- Missed Blocks Range : Displays the continuous missed blocks range based on the threshold provided in the `config.toml`.
- Missed Blocks In Last 48 hours 

### 3. Summary Dashboard
This dashboard displays a quick information summary of validator details and system metrics. It includes the following details.

- Validator identity (Moniker and hex Address)
- Validator summary (Node Status, Validator Status, Voting Power, Height Difference and No.Of peers) are the metrics being displayed from Validator details.
- Server uptime,CPU usage, RAM Usage, Memory usage and information about disk usage are the metrics being displayed from System details.
 
### 4. Heimdall network metrics

This dashboard displays the tendermint prometheus metrics.

 - Total txs : Count of txs on the network
 - Block size : Graphical display of block size.
 - Validators : Gauge display of count of **Active**, **Missing** and **Byzantine** validators.
 - Failed txs : Count of failed txs on the node
 - Total network I/O : Bar graph of total network throughput. 
### 5. System Monitoring Metrics
These metrics are are collected by the node_exporter and displays all the metrics related to 
 
 - CPU
 - Memory
 - Disk
 - Network traffic
 - System processes
 - Systemd
 - Storage
 - Hardware Misc
 
 

