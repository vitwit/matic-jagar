#!/bin/bash

set -e

cd $HOME

echo "------ checking for go if it's not installed then it will be installing here -----"

command_exists () {
    type "$1" &> /dev/null ;
}

if command_exists go ; then
    echo "Golang is already installed"
else
  echo "Install dependencies"
  sudo apt update
  sudo apt install build-essential jq -y

  wget https://dl.google.com/go/go1.15.2.linux-amd64.tar.gz
  tar -xvf go1.15.2.linux-amd64.tar.gz
  sudo mv go /usr/local

  echo "" >> ~/.bashrc
  echo 'export GOPATH=$HOME/go' >> ~/.bashrc
  echo 'export GOROOT=/usr/local/go' >> ~/.bashrc
  echo 'export GOBIN=$GOPATH/bin' >> ~/.bashrc
  echo 'export PATH=$PATH:/usr/local/go/bin:$GOBIN' >> ~/.bashrc

  #source ~/.bashrc
  . ~/.bashrc

  go version
fi

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

echo "----------- Intsalling prometheus -----------"

wget https://github.com/prometheus/prometheus/releases/download/v2.22.1/prometheus-2.22.1.linux-amd64.tar.gz
$ tar -xvf prometheus-2.22.1.linux-amd64.tar.gz

tar -xvf prometheus-2.22.1.linux-amd64.tar.gz

sudo cp prometheus-2.22.1.linux-amd64/prometheus $GOBIN

sudo cp prometheus-2.22.1.linux-amd64/prometheus.yml $HOME


echo "------- Edit prometheus.yml --------------"

echo "scrape_configs:
 
  - job_name: 'validator'

    static_configs:
    - targets: ['localhost:26660']


  - job_name: 'node_exporter'

    static_configs:
    - targets: ['localhost:9100']" >> "prometheus.yml"

echo "------- Setup prometheus system service -------"

echo "[Unit]
Description=Prometheus
After=network-online.target

[Service]
User=$USER
ExecStart=$HOME/go/bin/prometheus --config.file=$HOME/prometheus.yml
Restart=always
RestartSec=3
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target" | sudo tee "/lib/systemd/system/prometheus.service"

echo "------ Start prometheus -----------"

sudo systemctl daemon-reload
sudo systemctl enable prometheus.service
sudo systemctl start prometheus.service


echo "-------- Installing node exporter -----------"

cd $HOME

curl -LO https://github.com/prometheus/node_exporter/releases/download/v0.18.1/node_exporter-0.18.1.linux-amd64.tar.gz

tar -xvf node_exporter-0.18.1.linux-amd64.tar.gz

sudo cp node_exporter-0.18.1.linux-amd64/node_exporter $GOBIN

echo "---------- Setup Node exporter service -----------"

echo "[Unit]
Description=Node_exporter
After=network-online.target

[Service]
User=$USER
ExecStart=$HOME/go/bin/node_exporter
Restart=always
RestartSec=3
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target" | sudo tee "/lib/systemd/system/node_exporter.service"

echo "----------- Start node exporter ------------"

sudo systemctl daemon-reload

sudo systemctl enable node_exporter.service

sudo systemctl start node_exporter.service

echo "---- Cleaning .dep .tar.gz files of grafana, prometheus, influxdb and node exporter --------"

rm influxdb_1.8.3_amd64.deb grafana_7.3.1_amd64.deb node_exporter-0.18.1.linux-amd64.tar.gz prometheus-2.22.1.linux-amd64.tar.gz

echo "------------Creating databases matic-------------"

curl "http://localhost:8086/query" --data-urlencode "q=CREATE DATABASE matic"


echo "--------- Cloning matic-validator-mission-control -----------"

cd $HOME

git clone https://github.com/vitwit/matic-jagar.git

cd matic-jagar

git fetch && git checkout mumbai-testnet

cp example.config.toml config.toml

echo "------ Building and running the code --------"

go build -o matic-jagar && ./matic-jagar