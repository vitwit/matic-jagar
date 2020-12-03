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
    
  - job_name: 'sentry-1'

    static_configs:
    - targets: ['<SENTRY1-IP>:26660']
    
  - job_name: 'sentry-2'

    static_configs:
    - targets: ['<SENTRY2-IP>:26660']
```
**Note**: If you don't have any sentries or have one please skip or specify only one `job_name`


- Setup Prometheus System service
 ```
 sudo nano /lib/systemd/system/prometheus.service
 ```
 
 Copy-paste the following:
 
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
For the purpose of this guide it is assumed the `user` is `ubuntu`. If your user is different please make the required changes above.

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
 
 
 Copy-paste the following:
 
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
 For the purpose of this guide it is assumed the `user` is `ubuntu`. If your user is different please make the required changes above.

```sh
$ sudo systemctl daemon-reload
$ sudo systemctl enable node_exporter.service
$ sudo systemctl start node_exporter.service
```

To use custom bind ports for the prerequisites please follow these [instructions](./custom-port.md) 

#### Clean-up (Optional)

```
$ rm influxdb_1.8.3_amd64.deb grafana_7.3.1_amd64.deb node_exporter-0.18.1.linux-amd64.tar.gz prometheus-2.22.1.linux-amd64.tar.gz
```
