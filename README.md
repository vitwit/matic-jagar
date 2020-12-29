# Matic-jagar

![](https://github.com/vitwit/matic-jagar/blob/main/docs/logo.jpg)

**Matic-jagar** is a monitoring and alerting tool for validators operating on Matic Network. It provides separate grafana dashboards to better monitor the health of the validator and server. It has integration with Telegram and Sendgrid which enables it to provide updates via notifications or email. It uses InfluxDb and Prometheus to store the metrics and Grafana to display them. 

The alerting part of the tool has a modular approach which enables the user to decide on which metrics the alerts should be sent. Any and all notifications can be modified by a user to fit ones preference. By default a basic level of notifications is enabled for inexperienced users which can be modified by editing the config file.

## Features

- [Click here](./docs/dashboard-desc.md) to read about the different dashboards which are provided and a short description about them.

- [Click here](./docs/metric-calc.md) to read about how the metrics are calculated and displayed on the dashboards.

## Installation

- [Click here](./INSTRUCTIONS.md) to find the installation instructions.

- [Click here](./docs/upgrade.md) to find upgrade instructions for the tool.

If you would like to see any modifications or additional features to this tool, please feel free to open an issue and we will consider adding it to a future release.
