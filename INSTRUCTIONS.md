# Matic-Jagar setup


## Prerequisites
- **Go 13.x+**
- **Grafana 6.7+**
- **InfluxDB 1.7+**
- **Prometheus 2.x+**

### Prerequisite installation using script
#### 1) You can run the installation script to install prerequisites

- Script downloads and installs grafana, prometheus, influxdb and node exporter.
- It also downloads go if it's not already installed.
- This script takes `sentry-1` and/or `sentry-2` env variables and writes them to `prometheus.yml` file for gathering prometheus metrics emitted from the respective nodes. 
- Export the env variables using the following commands:
```bash
cd $HOME
export SENTRY1="<sentry-1-IP>" # ex:- export SENTRY1="143.125.36.5" 
export SENTRY2="<sentry-2-IP>" # ex:- export SENTRY2="143.185.336.95"
```
- If you don't have any sentries or have one, you can skip this or export only one IP address.
- **Note**: By default prometheus metrics are enabled on your nodes on port 26660. If you have changed the prometheus port on your node please edit the `~/prometheus.yml` and enter your custom port. 
- You can find the script [here](./scripts/install_prerequisites.sh)
- Execute the script using the following command:
```bash
curl -s -L https://raw.githubusercontent.com/vitwit/matic-jagar/main/scripts/install_prerequisites.sh | bash
```
Source your `.bashrc` after executing the script for the env variables to take effect.
```
source ~/.bashrc
```

**Note**: This script installs the prerequistes and enables them to run on their default ports ie. `Grafana` by default runs on port 3000, `InfluxDb` by default runs on port 8086, `Prometheus` by default runs on port 9090 and `Node Exporter` by default runs on port 9100. If you want to change the default ports please follow these [instructions](./docs/custom-port.md).



### 2) Manual installation of prerequisites

To manually install the prerequistes please follow this [guide](./docs/prereq-manual.md).



## Install and configure the tool

### 1) You can run the tool installation script to build and deploy

- It clones and sets up the monitoring tool as a system service.
- Please export the following env variables first as they will be used to initialize the `config.toml` file for the tool.
```bash
cd $HOME
export ETH_RPC_ENDPOINT="<infura-endpoint>" # Ex- export ETH_RPC_ENDPOINT= "https://goerli.prylabs.net"
export BOR_EXTERNAL_RPC="http://<sentry-ip>:8545" # Ex - export BOR_EXTERNAL_RPC="http://156.23.25.21:8545"
export HEIMDALL_EXTERNAL_RPC="http://<sentry-ip>:26657" # Ex - export HEIMDALL_EXTERNAL_RPC="http://156.23.25.21:26657"
export VAL_HEX_ADDRESS="<hex-address>" # Ex - export VAL_HEX_ADDRESS="E4B8E9225842401AD16D4D826732953DAF07C7E2". You can get it by running this cmd on validator- heimdallcli status | jq .validator_info.address
export SIGNER_ADDRESS="0xE4b8e9222705621aD16d4d826732953DAf07C7E2" # Ex- export SIGNER_ADDRESS="0xE4b8e9222705621aD16d4d826732953DAf07C7E2"
export VALIDATOR_NAME="moniker" # Your validator moniker
export TELEGRAM_CHAT_ID=<id> # Ex - export TELEGRAM_CHAT_ID=22828812
export TELEGRAM_BOT_TOKEN="<token>" # Ex - TELEGRAM_BOT_TOKEN="1117273891:AAEtr3ZU5x4JRj5YSF5LBeu1fPF0T4xj-UI"
```
- **Note**: If you don't want telegram notifications you can skip exporting `TELEGRAM_CHAT_ID` and `TELEGRAM_BOT_TOKEN` but the rest are mandatory.
- You can find the tool installation script [here](./scripts/tool_installation.sh).
- Run the script using the following command 
```bash
curl -s -L https://raw.githubusercontent.com/vitwit/matic-jagar/main/scripts/tool_installation.sh | bash
```

### 2) Manual installation of tool

```bash
git clone https://github.com/vitwit/matic-jagar.git

cd matic-jagar

git pull origin main 

mkdir -p  ~/.matic-jagar/config/

cp example.config.toml ~/.matic-jagar/config/config.toml
```

Edit the `config.toml` with your changes. Information on all the fields in `config.toml` can be found [here](./docs/config-desc.md)


-  Build and run the monitoring binary

```sh
$ go build -o matic-jagar && ./matic-jagar
```

Installation of the tool is completed lets configure the Grafana dashboards.

## Grafana Dashboards

The repo provides five dashboards

1. Validator Monitoring Metrics - Displays the validator metrics which are calculated and stored in influxdb.
2. Bor - Displays the bor metrics of validator which are calculated and stored in influxdb.
3. System Metrics - Displays the metrics related to your validator server on which this tool is hosted on.
4. Summary - Displays a quick overview of heimdall, bor and system metrics.
5. Setup Overview - Displays current block height, block time, number of connected peers and number of unconfirmed transactions of validator and two sentries.

Information on all the dashboards can be found [here](./docs/dashboard-desc.md).


## Importing dashboards

### 1. Login to your Grafana dashboard
- Open your web browser and go to http://<your_validator_ip>:3000/. `3000` is the default HTTP port that Grafana runs on if you haven’t configured a different port. Please make sure your firewall allows it to be accesed from the browser.
- If you are a first time user type `admin` for the username and password in the login page.
- You can change the password after login.

### 2. Create Datasources

- Before importing the dashboards you have to create a datasource of **InfluxDBMatic**.
- To create datasources go to configuration and select Data Sources.
- Click on `Add data source` and select InfluxDB from Time series databases section.
- Name the datasource as`InfluxDBMatic`. Replace the URL with `http://localhost:8086`. In InfluxDB Details section replace Database name as `matic`.
- Click on **Save & Test** . Now you have a working Datasource of InfluxDBMatic.

- For a **Prometheus** data source, click on `Add data source` and select `Prometheus`. Replace the URL with `http//localhost:9090`. Click on **Save & Test** . Now you have a working Datasource of Prometheus.


### 3. Import the dashboards
- To import the dashboards click the **+** button present on left hand side of the dashboard. Click on import and paste the UID of the dashboards on the text field below **Import via grafana.com** and click on load. 

- Select the datasources and click on import.

UID of dashboards are as follows:

- **13441**: Validator monitoring dashboard
- **13442**: Bor metrics dashboard
- **13443**: Summary dashboard
- **13444**: Setup overview dashboard
- **13445**: System monitoring metrics dashboard

*For more info about grafana dashboard imports you can refer https://grafana.com/docs/grafana/latest/reference/export_import/*




