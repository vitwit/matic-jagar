# Validator Mission Control


## Install Prerequisites
- **Go 13.x+**
- **Grafana 6.7+**
- **InfluxDB 1.7+**
- **Telegraf 1.14+**


### Install Grafana for Ubuntu
Download the latest .deb file and extract it by using the following commands

```sh
$ cd $HOME
$ sudo -S apt-get install -y adduser libfontconfig1
$ wget https://dl.grafana.com/oss/release/grafana_7.3.1_amd64.deb
$ sudo -S dpkg -i grafana_7.3.1_amd64.deb
```

Start the grafana server
```sh
$ sudo -S systemctl daemon-reload

$ sudo -S systemctl start grafana-server

Grafana will be running on port :3000 (ex:: https://localhost:3000)
```

### Install InfluxDB and Telegraf

```sh
$ cd $HOME
$ wget -qO- https://repos.influxdata.com/influxdb.key | sudo apt-key add -
$ source /etc/lsb-release
$ echo "deb https://repos.influxdata.com/${DISTRIB_ID,,} ${DISTRIB_CODENAME} stable" | sudo tee /etc/apt/sources.list.d/influxdb.list
```

Start influxDB

```sh
$ sudo -S apt-get update && sudo apt-get install influxdb
$ sudo -S service influxdb start

The default port that runs the InfluxDB HTTP service is :8086
```

**Note :** If you want cusomize the configuration, edit `influxdb.conf` at `/etc/influxdb/influxdb.conf` and don't forget to restart the server after the changes. You can find a sample 'influxdb.conf' [file here](https://github.com/jheyman/influxdb/blob/master/influxdb.conf).


Start telegraf

```sh
$ sudo -S apt-get update && sudo apt-get install telegraf
$ sudo -S service telegraf start
```

## Install and configure the Validator Mission Control

### Get the code

```bash
$ git clone https://github.com/vitwit/matic-jagar.git
$ cd matic-jagar
$ git fetch && git checkout refactor
$ cp example.config.toml config.toml
```

### Configure the following variables in `config.toml`

- **[rpc_and_lcd_endpoints]**

    - *eth_rpc_endpoint*

        Etherium rpc endpoint useful to gather information about validator staking rewards, balance, commission rate and to query valiator share contract address.

    - *bor_rpc_end_point*

        Bor rpc endpoint is useful to get metrics of bor node such as block height, current proposer, pending transactions and also precommits of a block to know about missed blocks.

    - *bor_external_rpc*

        Bor rxternal open RPC endpoint(secondary RPC other than your own validator). Useful to get network block height and to calculate block time difference of your validator and network.

    - *heimdall_rpc_endpoint*
        
        Heimdall rpc end point (RPC of your own validator) useful to gather information about network info, validator voting power, unconfirmed txns etc.

    - *heimdall_lcd_endpoint*

        Address of your lcd client (ex: http://localhost:1317), Which will be used to gather information like latest block info, blanaces and staking related matrics etc.

    - *heimdall_external_rpc*

        Heimdall rpc end point (RPC of your own validator) useful to gather information about network info, validator voting power, unconfirmed txns etc.

- **[validator_details]**

    - *validator_hex_addr*

        Validator hex address useful to know about last proposed block, missed blocks and voting power, etc.

    - *signer_address*

        Signer address of your validator which will be used to get information about staking, balances an voting power.

    - *validator_name*

        Provide name of your validator, to get it displayed in alerts.

    - *stake_manager_contract*

        Address of stake manager contract, which will be used to query the methods of stake manager contract and also to get contract address of validator share contract.

- **[enable_alerts]**

    - *enable_telegram_alerts*

        Configure **yes** if you wish to get telegram alerts otherwise make it **no** .
    
    - *enable_email_alerts*
    
        Configure **yes** if you wish to get email alerts otherwise make it **no** .

- **[daily_alert]**

    -   Alert about validator health, i.e. whether it's voting or jailed. You can get alerts twice a day based on the time you will configure i.e., **alert_time1** and **alert_time2** .

- **[choose_alerts]**

    - *balance_change_alerts*

        If you want to get alerts about balance change cofigure **yes** otherwise make it **no** .

    - *voting_power_alerts*

        If you want to get alerts about voting power change cofigure **yes** otherwise make it **no** .

    - *proposal_alerts*

        If you want to recieve alerts about new proposal and whenever there is a change in proposal status like deposit_perio to voting_period etc, then configure **yes** otherwise make it **no**.

    - *block_diff_alerts*

        If you want to recieve alerts when there is a change in your validator block height and network height then make it **yes** otherwise **no** .

    - *missed_block_alerts*

        If you want to get alerts when your validator is missing blocks then configure it **yes** otherwise **no** .

    - *num_peers_alerts*

        If you want to be notified, when there is a drop in number of peers connected to your validator, then configure it **yes** or else **no** .
    
    - *node_sync_alert*

        If you want to be notified about your node syncing status then make it **yes** otherwise **no**.

- **[alerting_threholds]**

    - *num_peers_threshold*

        Configure the threshold to get an alert if the no.of connected peers falls below the given threshold.

    - *missed_blocks_threshold*

        Configure the threshold to receive missed block alerts, e.g. a value of 10 would alert you every time you've missed 10 consecutive blocks.
    
    - *block_diff_threshold*

        An integer value to receive block difference alerts, e.g. a value of 2 would alert you if your validator falls 2 or more blocks behind the chain's current block height.

- **[telegram]**

    - *tg_chat_id*

        Telegram chat ID to receive Telegram alerts, required for Telegram alerting.
    
    - *tg_bot_token*

        Telegram bot token, required for Telegram alerting. The bot should be added to the chat and should have send message permission.

- **[sendgrid]**

    - *sendgrid_token*

        Sendgrid mail service api token, required for e-mail alerting.

    - *email_address*

        E-mail address to receive mail notifications, required for e-mail alerting.

- **[influxdb]**

    - *database*

        Name of your influxdb database, to which in which you want to store the data.

    - *username*

        Provice username if have any.

After populating config.toml, check if you have connected to influxdb and created a database which you are going to use.

If not or If your connection throws error "database not found", create a database

```bash
$   cd $HOME
$   influx
>   CREATE DATABASE db_name   (ex: CREATE DATABASE matic)
$   exit
```

After all these steps, build and run the monitoring binary

$ **go build -o matic && ./matic**

We have finished the installation and started the server. Now lets configure the Grafana dashboard.

## Grafana Dashboards

Validator Mission Control provides three dashboards

1. Validator Monitoring Metrics (These are the heimdall metrics which we have calculated and stored in influxdb)
2. Bor (These are the bor metrics which we have calculated and stored in influxdb)
3. System Metrics (These are the metrics related to the system configuration which come from telegraf)
4. Summary (Which gives quick overview of heimall, bor and system metrics)


### 1. Validator monitoring metrics (Heimdall)
The following list of metrics are displayed in this dashboard.

- Validator Details :  Displays the details of a validator like moniker, valiator signer address and hex address.
- Node Status :  Displays whether the node is running or not in the form of UP and DOWN.
- Validator Status :  Displays the validator health. Shows Voting if the validator is in active state or else Jailed.
- Validator Caught Up : Displays whether the validator node is in sync with the network or not.
- Block Time Difference : Displays the time difference between previous block and current block.
- Current Block Height - Validator :  Validator : Displays the current block height committed by the validator.
- Latest Block Height - Network : Network : Displays the latest block height of a network.
- Height Difference : Displays the difference between heights of validator current block height and network latest block height.
- Missed Blocks : Displays a graph about missed blocks.
- Last Missed Block Range : Displays the continuous missed blocks range based on the threshold given in the config.toml
- Blocks Missed In last 48h : Displays the count of blocks missed by the validator in last 48 hours.
- Unconfirmed Txns : Displays the number of unconfirmed transactions on that node.
- Latest Checkpoint : Displays the height of the latest check point.
- No.of Peers : Displays the total number of peers connected to the validator.
- Peer Address : Displays the ip addresses of connected peers.
- Latency : Displays the latency of connected peers with respect to the validator.
- Validator Fee : Displays the commission rate of the validator.
- Voting Power : Displays the voting power of the validator.
- Max Tx Gas : Displays the max transaction gas.
- Rewards : Displays the rewards of your validator.
- Unclaimed Rewards : Displays the current unclaimed rewards amount of the validator.
- Last proposed Block Height : Displays height of the last block proposed by the validator.
- Last Proposed Block Time : Displays the time of the last block proposed by the validator.
- Heimdall Current Balance : Displays the account balance of the validator.
- Bor Current Balance : Displays the current balance of your bor node.
- Self Stake : Displays the self stake of your valiator.
- Voting Period Proposals : Displays the list of the proposals which are currently in voting period.
- Deposit Period Proposals : Displays the list of the proposals which are currently in deposit period.
- Completed Proposals : Displays the list of the proposals which are completed with their status as passed or rejected.


**Note:** The above mentioned metrics will be calculated and displayed according to the validator address which will be configured in config.toml.

For alerts regarding system metrics, a Telegram bot can be set up on the dashboard itself. A new notification channel can be added for the Telegram bot by clicking on the bell icon on the left hand sidebar of the dashboard. 

This will let the user configure the Telegram bot ID and chat ID. **A custom alert** can be set for each graph in a Grafana dashboard by clicking on the edit button and adding alert rules.

### 2. Bor
Displays the metrics of bor node such as,

    - Current Block Height - validator :
    - Current Block Height - network :
    - Block Height Difference :
    - Current Span :
    - Pending Transactions :
    - Current Block Proposer :
    - No.of Blocks Proposed :
    - No.of Blocs Signed :
    - No.of Span Validator is part of :
    - Misse Blocks Range :
    - Last Misse Block Range :
    - Missed Blocks In Last 48 hours :

### 3. System Monitoring Metrics
These are powered by telegraf.

-  For the list of system monitoring metrics, you can refer `telgraf.conf`. You can replace the file with your original telegraf.conf file which will be located at /telegraf/etc/telegraf (installation directory).
 
 ### 4. Summary Dashboard
This dashboard displays a quick information summary of validator details and system metrics. It includes following details.

- Validator identity (Moniker and hex Address)
- Validator summary (Node Status, Validator Status, Voting Power, Height Difference and No.Of peers) are the metrics being displayed from Validator details.
- Server uptime,CPU usage, RAM Usage, Memory usage and information about disk usage are the metrics being displayed from System details.

## How to import these dashboards in your Grafana installation

### 1. Login to your Grafana dashboard
- Open your web browser and go to http://<your_ip>:3000/. `3000` is the default HTTP port that Grafana listens to if you haven’t configured a different port.
- If you are a first time user type `admin` for the username and password in the login page.
- You can change the password after login.

### 2. Create Datasources

- Before importing the dashboards you have to create datasources of InfluxDBTelegraf and InfluxDBMatic.
- To create datasoruces go to configuration and select Data Sources.
- After that you can find Add data source, select InfluxDB from Time series databases section.
- Then to create `InfluxDBMatic` Datasource, follow these configurations. In place of name give InfluxDBMatic, in place of URL give url of influxdb where it is running (ex : http://ip_address:8086). Finaly in InfluxDB Details section give Database name as `matic` (If you haven't created a database with different name). You can give User and Password of influx if you have set anthing, otherwise you can leave it empty.
- After this configuration click on Save & Test. Now you have a working Datasource of InfluxDBMatic.

- Repeat same steps to create `InfluxDBTelegraf` Datasource. In place of name give InfluxDBTelegraf, give URL of telegraf where it is running (ex: http://ip_address:8086). Give Database name as telegraf, user and password (If you have configured any). 

- After this configuration click on Save & Test. Now you have a working Datasource of InfluxDBTelegraf.

### 3. Import the dashboards
- To import the json file of the **validator monitoring metrics** click the *plus* button present on left hand side of the dashboard. Click on import and load the validator_monitoring_metrics.json present in the grafana_template folder. 

- Select the datasources and click on import.

- To import **system monitoring metrics** click the *plus* button present on left hand side of the dashboard. Click on import and load the system_monitoring_metrics.json present in the grafana_template folder.

- While creating this dashboard if you face any issues at valueset, change it to empty and then click on import by selecting the datasources.

- To import **summary**, click the *plus* button present on left hand side of the dashboard. Click on import and load the summary.json present in the grafana_template folder.

- To import **bor**, click the *plus* button present on left hand side of the dashboard. Click on import and load the bor.json present in the grafana_template folder.

- *For more info about grafana dashboard imports you can refer https://grafana.com/docs/grafana/latest/reference/export_import/*