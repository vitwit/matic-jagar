- **[rpc_and_lcd_endpoints]**

    - *eth_rpc_endpoint*

        Ethereum rpc endpoint is used to gather information about validator staking rewards, balance, commission rate and to query valiator share contract address.

    - *bor_rpc_end_point*

        Bor rpc endpoint is used to get metrics of bor node such as block height, current proposer, pending transactions and also precommits of a block to know about missed blocks.

    - *bor_external_rpc*

        Secondary RPC other than your own validator. Useful to get network block height and to calculate block time difference between your validator and network.

    - *heimdall_rpc_endpoint*
        
        Heimdall rpc end point (RPC of your own validator) is used to gather information about network info, validator voting power, unconfirmed txns etc.

    - *heimdall_lcd_endpoint*

        Heimdall rpc end point (ex: http://localhost:1317) is used to gather information like latest block info, balances and staking related metrics.

    - *heimdall_external_rpc*

        Secondary RPC other than your own validator. useful to gather information about network info, validator voting power, unconfirmed txns etc.

- **[validator_details]**

    - *validator_hex_addr*

        Validator hex address is used to verify last proposed block, missed blocks and voting power etc.

    - *signer_address*

        Signer address of your validator is used to get information about staking, balances and voting power.

    - *validator_name*

        Moniker of your validator, to get it displayed in alerts.

    - *stake_manager_contract*

        Address of stake manager contract, this is used to query the methods of stake manager contract and also to get contract address of validator share contract.

- **[enable_alerts]**

    - *enable_telegram_alerts*

        Configure **yes** if you wish to get telegram alerts otherwise make it **no** .
    
    - *enable_email_alerts*
    
        Configure **yes** if you wish to get email alerts otherwise make it **no** .

- **[daily_alert]**

    -   Alert about validator health, i.e. whether it's voting or jailed. You can get alerts twice a day based on the time which can be configured i.e., **alert_time1** and **alert_time2** .

- **[choose_alerts]**

    - *balance_change_alerts*

        If you want to get alerts about heimdall balance change cofigure **yes** otherwise make it **no** .

    - *voting_power_alerts*

        If you want to get alerts about voting power change cofigure **yes** otherwise make it **no** .

    - *proposal_alerts*

        If you want to recieve alerts about new proposal and whenever there is a change in  status like deposit_period to voting_period etc, then configure **yes** otherwise make it **no**.

    - *block_diff_alerts*

        If you want to recieve alerts when there is a gap between your validator block height and network height then make it **yes** otherwise **no** .

    - *missed_block_alerts*

        If you want to get alerts when your validator is missing blocks then configure it **yes** otherwise **no** .

    - *num_peers_alerts*

        If you want to be notified, when there is a drop in number of peers connected to your validator, then configure it **yes** or else **no** .
    
    - *node_sync_alert*

        If you want to be notified about the status of your node syncing then make it **yes** otherwise **no**.

- **[alerting_threholds]**

    - *num_peers_threshold*

        Configure the threshold to get an alert if the no.of connected peers falls below the given threshold.

    - *missed_blocks_threshold*

        Configure the threshold to receive missed block alerts, e.g. a value of 10 would alert you every time you've missed 10 consecutive blocks.
    
    - *block_diff_threshold*

        An integer value to receive block difference alerts, e.g. a value of 2 would alert you if your validator falls 2 or more blocks behind the network's current block height.

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

        Name of your influxdb database in which you want to store the data.

    - *username*

        Provide username if configured.


