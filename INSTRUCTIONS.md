# Matic-Jagar setup


## Install Prerequisites
- **Go 13.x+**
- **Grafana 6.7+**
- **InfluxDB 1.7+**


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

You can view the `grafana` logs using this:
```
journalctl -u grafana-server -f
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
$ cp prometheus-2.22.1.linux-amd64/prometheus $GOBIN
$ cp prometheus-2.22.1.linux-amd64/prometheus.yml $HOME
```

- Add the following in prometheus.yml using your editor of choice and replace the values of `<sentry-IP>` with the IP addresses of your sentries.

```sh
 scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: 'validator'

    static_configs:
    - targets: ['localhost:26660']


  - job_name: 'node_exporter'

    static_configs:
    - targets: ['localhost:9100']

  - job_name: 'sentry-1'

    static_configs:
    - targets: ['<sentry-1-IP>:26660']

  - job_name: 'sentry-2'

    static_configs:
    - targets: ['<sentry-2-IP>:26660']
```
Indentations in the `prometheus.yml` file are important. You can find a sample configuration file [here](./docs/prometheus.yml)


- Setup Prometheus System service
 ```
 sudo nano /lib/systemd/system/prometheus.service
 ```
 
 Copy-paste the following.
 
 **Note :** It is assumed for this setup purposes you are running the services as `ubuntu`. If your `user` is different please make the necessary changes in systemd file.
 
 
 ```
 [Unit]
Description=Prometheus
After=network-online.target

[Service]
Type=simple
ExecStart=/home/ubuntu/go/bin/prometheus --config.file=/home/ubuntu/prometheus.yml
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
You can view the `prometheus` logs using this:
```
journalctl -u prometheus -f
``` 

### Install node exporter


```sh
$ cd $HOME
$ curl -LO https://github.com/prometheus/node_exporter/releases/download/v0.18.1/node_exporter-0.18.1.linux-amd64.tar.gz
$ tar -xvf node_exporter-0.18.1.linux-amd64.tar.gz
$ cp node_exporter-0.18.1.linux-amd64/node_exporter $GOBIN
```
- Setup Node exporter service

```
 sudo nano /lib/systemd/system/node_exporter.service
 ```
 
 
 Copy-paste the following.
 
 **Note :** It is assumed for this setup purposes you are running the services as `ubuntu`. If your `user` is different please make the necessary changes in systemd file.
 
 ```
 [Unit]
Description=Node_exporter
After=network-online.target

[Service]
Type=simple
ExecStart=/home/ubuntu/go/bin/node_exporter
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

You can view the `node_exporter` logs using this:
```
journalctl -u node_exporter -f
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
$ go build -o matic && ./matic
```

Installation of the tool is completed lets configure the Grafana dashboards.

## Grafana Dashboards

The repo provides five dashboards

1. Validator Monitoring Metrics - Displays the validator metrics which are calculated and stored in influxdb.
2. Bor - Displays the bor metrics of validator which are calculated and stored in influxdb.
3. System Metrics - Displays the metrics related to your validator server on which this tool is hosted on.
4. Summary - Displays a quick overview of heimdall, bor and system metrics.
5. Setup Overview - Displays current block height, block time, number of connected peers and number of unconfirmed transactions.

Information on all the dashboards can be found [here](./docs/dashboard-desc.md)


## Importing dashboards

### 1. Login to your Grafana dashboard
- Open your web browser and go to http://<your_ip>:3000/. `3000` is the default HTTP port that Grafana runs on if you haven’t configured a different port.
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
- To import the json file of the **validator monitoring metrics** click the *plus* button present on left hand side of the dashboard. Click on import and load the validator_monitoring_metrics.json present in the grafana_template folder. 

- Select the datasources and click on import.

- Follow the same steps to import **system_monitoring_metrics.json**, **setup_overview.json**, **bor.json**, **summary.json**. 


- *For more info about grafana dashboard imports you can refer https://grafana.com/docs/grafana/latest/reference/export_import/*



