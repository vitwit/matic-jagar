#!/bin/bash

set -e

cd $HOME

echo "--------- Cloning matic-validator-mission-control -----------"

git clone https://github.com/vitwit/matic-jagar.git

cd matic-jagar

git pull origin master

mkdir -p  ~/.matic-jagar/config/

cp example.config.toml ~/.matic-jagar/config/config.toml

cd $HOME

echo "------------ Creating database matic in influxdb-------------"

curl "http://localhost:8086/query" --data-urlencode "q=CREATE DATABASE matic"

echo "------ Updatig config fields with exported values -------"

sed -i '/eth_rpc_endpoint =/c\eth_rpc_endpoint = "'"$ETH_RPC_ENDPOINT"'"'  ~/.matic-jagar/config/config.toml

sed -i '/bor_external_rpc =/c\bor_external_rpc = "'"$BOR_EXTERNAL_RPC"'"'  ~/.matic-jagar/config/config.toml

sed -i '/heimdall_external_rpc =/c\heimdall_external_rpc = "'"$HEIMDALL_EXTERNAL_RPC"'"'  ~/.matic-jagar/config/config.toml

sed -i '/validator_hex_addr =/c\validator_hex_addr = "'"$VAL_HEX_ADDRESS"'"'  ~/.matic-jagar/config/config.toml

sed -i '/signer_address =/c\signer_address = "'"$SIGNER_ADDRESS"'"'  ~/.matic-jagar/config/config.toml

sed -i '/validator_name =/c\validator_name = "'"$VALIDATOR_NAME"'"'  ~/.matic-jagar/config/config.toml

if [ ! -z "${TELEGRAM_CHAT_ID}" ] && [ ! -z "${TELEGRAM_BOT_TOKEN}" ];
then 
    sed -i '/tg_chat_id =/c\tg_chat_id = '"$TELEGRAM_CHAT_ID"''  ~/.matic-jagar/config/config.toml

    sed -i '/tg_bot_token =/c\tg_bot_token = "'"$TELEGRAM_BOT_TOKEN"'"'  ~/.matic-jagar/config/config.toml

    sed -i '/enable_telegram_alerts =/c\enable_telegram_alerts = 'true''  ~/.matic-jagar/config/config.toml
else
    echo "---- Telgram chat id and/or bot token are empty --------"
fi

echo "------ Building and running the code --------"

cd matic-jagar

go build -o matic-jagar
mv matic-jagar $HOME/go/bin

echo "---------- Setup Matic-Jagar service -----------"

echo "[Unit]
Description=Matic-Jagar
After=network-online.target

[Service]
Type=simple
ExecStart=$HOME/go/bin/matic-jagar
Restart=always
RestartSec=3
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target" | sudo tee "/lib/systemd/system/matic_jagar.service"

echo "---------- Start Matic-Jagar service -----------"

sudo systemctl daemon-reload

sudo systemctl enable matic_jagar.service

sudo systemctl start matic_jagar.service

echo "** Done **"
