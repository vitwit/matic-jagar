#!/bin/bash

set -e

cd $HOME

echo "----------- Installing grafana -----------"

sudo -S apt-get install -y adduser libfontconfig1

wget https://dl.grafana.com/oss/release/grafana_7.3.1_amd64.deb

sudo -S dpkg -i grafana_7.3.1_amd64.deb

echo "------ Starting grafana server using systemd --------"

sudo -S systemctl daemon-reload

sudo -S systemctl start grafana-server

cd $HOME

echo "----------- Installing Influx -----------"

wget https://dl.influxdata.com/influxdb/releases/influxdb_1.8.3_amd64.deb

sudo dpkg -i influxdb_1.8.3_amd64.deb

echo "----------- Starting Influxdb -----------"

sudo systemctl start influxdb 

cd $HOME

echo "----------- Downloading prometheus -----------"

wget https://github.com/prometheus/prometheus/releases/download/v2.22.1/prometheus-2.22.1.linux-amd64.tar.gz
$ tar -xvf prometheus-2.22.1.linux-amd64.tar.gz

tar -xvf prometheus-2.22.1.linux-amd64.tar.gz

sudo cp prometheus-2.22.1.linux-amd64/prometheus $GOBIN

sudo cp prometheus-2.22.1.linux-amd64/prometheus.yml $HOME


echo "------ Start prometheus -----------"

sudo systemctl daemon-reload
sudo systemctl enable prometheus.service
sudo systemctl start prometheus.service


echo "-------- Installing node exporter -----------"

cd $HOME

curl -LO https://github.com/prometheus/node_exporter/releases/download/v0.18.1/node_exporter-0.18.1.linux-amd64.tar.gz

tar -xvf node_exporter-0.18.1.linux-amd64.tar.gz

sudo cp node_exporter-0.18.1.linux-amd64/node_exporter $GOBIN

echo "----------- Start node exporter ------------"

sudo systemctl daemon-reload

sudo systemctl enable node_exporter.service

sudo systemctl start node_exporter.service

echo "------------Creating databases matic-------------"

curl "http://localhost:8086/query" --data-urlencode "q=CREATE DATABASE matic"


echo "--------- Cloning matic-validator-mission-control -----------"

cd $HOME

git clone https://github.com/vitwit/matic-jagar.git

cd matic-jagar

git fetch && git checkout mumbai-testnet

cp example.config.toml config.toml

echo "------ Building and running the code --------"

go build -o matic && ./matic