# Matic-Jagar setup


## Install Prerequisites
- **Go 13.x+**
- **Grafana 6.7+**
- **InfluxDB 1.7+**

## 1.Installation using script
### You can run below installation scripts to install prerequisites and to setup the monitoring tool

#### Script which downloads and runs the prerequisites
- One click installation to download grafana, prometheus, influxdb and node exporter.
- It also downloads go if it's not installed.
- Before running the script, you have to export `sentry-1` and/or `sentry-2` IPs if you have any.
- Follow the below steps to export those IPs just by replacig the values
```bash
cd $HOME
export SENTRY1="http://localhost:8000" # Replace with your sentry-1 IP
export SENTRY2="http://localhost:8000" # Replace with your sentry-2 IP
```
- Then download the [prerequisites script](https://github.com/vitwit/matic-jagar/blob/review-changes/scripts/install_prerequisites.sh) and run.
- To run
```bash
chmod +x install_prerequisites.sh
./install_prerequisites.sh
```

#### Script which downloads and run matic-jagar

- It clones and setup the matic-jagar monitoring tool as a system service.
- Before running the script make sure to follow below steps to export the config fields. The exported values will be refelected in config.toml of matic-jagar.
- Don't forget to change and export the field values.
```bash
cd $HOME
export ETH_RPC_ENDPOINT="https://goerli.prylabs.net" # Replace with infura endpoint
export BOR_EXTERNAL_RPC="http://<sentry-ip>:8545" # Replace the IP address with your sentry IP address
export HEIMDALL_EXTERNAL_RPC="http://<sentry-ip>:26657" # Replace the IP address with your sentry IP address
export VAL_HEX_ADDRESS="E4B8E9225842401AD16D4D826732953DAF07C7E2" # Replace this address with your validator hex address. You can get it by running this cmd on validator- heimdallcli status | jq .validator_info.address
export SIGNER_ADDRESS="0xE4b8e9222705621aD16d4d826732953DAf07C7E2" # Replace this with your valdiator signer address
export VALIDATOR_NAME="moniker" # Your validator moniker
export TELEGRAM_CHAT_ID=22828812 # Replace your chat id here
export TELEGRAM_BOT_TOKEN="1117273891:AAEtr3ZU5x4JRj5YSF5LBeu1fPF0T4xj-UI" # Replace your bot token here
```
- After exporting above fields, you can just download the installation script and run it.
- Here you can find matic-jagar
[installation script](https://github.com/vitwit/matic-jagar/blob/review-changes/scripts/tool_installation.sh)
- To run script 
```bash
chmod +x tool_installation.sh
./tool_installation.sh
```

## 2. Install manually
### Install Grafana for Ubuntu
Download the latest .deb file and extract it:

```sh
$ cd $HOME
$ sudo -S apt-get install -y libfontconfig1
$ wget https://dl.grafana.com/oss/release/grafana_7.3.1_amd64.deb
$ sudo -S dpkg -i grafana_7.3.1_amd64.deb
```

Start the grafana server
```
$ sudo -S systemctl daemon-reload

$ sudo -S systemctl start grafana-server

The default port that Grafana runs on is 3000. 
```

### Install InfluxDB

```sh
$ wget https://dl.influxdata.com/influxdb/releases/influxdb_1.8.3_amd64.deb
$ sudo dpkg -i influxdb_1.8.3_amd64.deb
```

Start influxDB

```sh
$ sudo systemctl start influxdb 

The default port that runs the InfluxDB HTTP service is 8086
```

Create an influxDB database:

```sh
$   cd $HOME
$   influx
>   CREATE DATABASE matic  
$   exit
```

**Note :** If you want to cusomize the configuration, edit `influxdb.conf` at `/etc/influxdb/influxdb.conf` and restart the server after the changes. You can find a sample 'influxdb.conf' [file here](https://github.com/jheyman/influxdb/blob/master/influxdb.conf).


### Install Prometheus 

```sh
$ cd $HOME
$ wget https://github.com/prometheus/prometheus/releases/download/v2.22.1/prometheus-2.22.1.linux-amd64.tar.gz
$ tar -xvf prometheus-2.22.1.linux-amd64.tar.gz
$ sudo cp prometheus-2.22.1.linux-amd64/prometheus $GOBIN
$ sudo cp prometheus-2.22.1.linux-amd64/prometheus.yml $HOME
```

- Add the following in prometheus.yml using your editor of choice

```sh
 scrape_configs:
 
  - job_name: 'validator'

    static_configs:
    - targets: ['localhost:26660']


  - job_name: 'node_exporter'

    static_configs:
    - targets: ['localhost:9100']
```

- Setup Prometheus System service
 ```
 sudo nano /lib/systemd/system/prometheus.service
 ```
 
 Copy-paste the following and replace the <user> variable with your user.
 
 ```
 [Unit]
Description=Prometheus
After=network-online.target

[Service]
User=<user>
ExecStart=/home/<user>/go/bin/prometheus --config.file=/home/<user>/prometheus.yml
Restart=always
RestartSec=3
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target
 ```


```sh
$ sudo systemctl daemon-reload
$ sudo systemctl enable prometheus.service
$ sudo systemctl start prometheus.service
```

### Install node exporter


```sh
$ cd $HOME
$ curl -LO https://github.com/prometheus/node_exporter/releases/download/v0.18.1/node_exporter-0.18.1.linux-amd64.tar.gz
$ tar -xvf node_exporter-0.18.1.linux-amd64.tar.gz
$ sudo cp node_exporter-0.18.1.linux-amd64/node_exporter $GOBIN
```
- Setup Node exporter service

```
 sudo nano /lib/systemd/system/node_exporter.service
 ```
 
 
 Copy-paste the following and replace the <user> variable with your user.
 
 ```
 [Unit]
Description=Node_exporter
After=network-online.target

[Service]
User=<user>
ExecStart=/home/<user>/go/bin/node_exporter
Restart=always
RestartSec=3
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target
 ```

```sh
$ sudo systemctl daemon-reload
$ sudo systemctl enable node_exporter.service
$ sudo systemctl start node_exporter.service
```
#### Clean-up (Optional)

```
$ rm influxdb_1.8.3_amd64.deb grafana_7.3.1_amd64.deb node_exporter-0.18.1.linux-amd64.tar.gz prometheus-2.22.1.linux-amd64.tar.gz
```

## Install and configure the tool

### Get the code

```bash
$ git clone https://github.com/vitwit/matic-jagar.git
$ cd matic-jagar
$ git fetch && git checkout mumbai-testnet
$ cp example.config.toml config.toml
```

Edit the `config.toml` with your changes. Informaion on all the fields in `config.toml` can be found [here](./docs/config-desc.md)


## Build and run the monitoring binary

```sh
$ go build -o matic-jagar && ./matic-jagar
```

Installation of the tool is completed lets configure the Grafana dashboards.

## Grafana Dashboards

The repo provides five dashboards

1. Validator Monitoring Metrics - These are the validator metrics which are calculated and stored in influxdb.
2. Bor - These are the bor metrics which are calculated and stored in influxdb.
3. System Metrics - These are the metrics related to your validator server on which this tool is hosted on.
4. Summary -  gives quick overview of heimall, bor and system metrics.
5. Heimdall network metrics - These are tendermint prometheus metrics emmitted by the node.

Information on all the dashboards can be found [here](./docs/dashboard-desc.md)


## Importing dashboards

### 1. Login to your Grafana dashboard
- Open your web browser and go to http://<your_ip>:3000/. `3000` is the default HTTP port that Grafana listens to if you haven’t configured a different port.
- If you are a first time user type `admin` for the username and password in the login page.
- You can change the password after login.

### 2. Create Datasources

- Before importing the dashboards you have to create a datasource of **InfluxDBMatic**.
- To create datasoruces go to configuration and select Data Sources.
- Click on `Add data source` and select InfluxDB from Time series databases section.
- Name the datasource as`InfluxDBMatic`. Replace the URL with `http://localhost:8086`. In InfluxDB Details section replace Database name as `matic`.
- Click on **Save & Test** . Now you have a working Datasource of InfluxDBMatic.

- For a **Prometheus** data source, click on `Add data source` and select `Prometheus`. Replace the URL with `http//localhost:9090`. Click on **Save & Test** . Now you have a working Datasource of Prometheus.


### 3. Import the dashboards
- To import the json file of the **validator monitoring metrics** click the *plus* button present on left hand side of the dashboard. Click on import and load the validator_monitoring_metrics.json present in the grafana_template folder. 

- Select the datasources and click on import.

- Follow the same steps to import **system_monitoring_metrics.json**, **heimdall_network_metrics.json**, **bor.json**, **summary.json**. 


- *For more info about grafana dashboard imports you can refer https://grafana.com/docs/grafana/latest/reference/export_import/*



