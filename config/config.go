package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

type (
	//Telegram bot details struct
	Telegram struct {
		BotToken string `mapstructure:"tg_bot_token"`
		ChatID   int64  `mapstructure:"tg_chat_id"`
	}

	//SendGrid tokens
	SendGrid struct {
		Token        string `mapstructure:"sendgrid_token"`
		EmailAddress string `mapstructure:"email_address"`
	}

	//Scraper time interval
	Scraper struct {
		Rate          string `mapstructure:"rate"`
		Port          string `mapstructure:"port"`
		ValidatorRate string `mapstructure:"validator_rate"`
		ContractRate  string `mapstructure:"contract_rate"`
	}

	//InfluxDB details
	InfluxDB struct {
		Port     string `mapstructure:"port"`
		Database string `mapstructure:"database"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	}

	// Endpoints is RPC and LCD endpoints struct
	Endpoints struct {
		EthRPCEndpoint      string `mapstructure:"eth_rpc_endpoint"`
		BorRPCEndpoint      string `mapstructure:"bor_rpc_end_point"`
		BorExternalRPC      string `mapstructure:"bor_external_rpc"`
		HeimdallRPCEndpoint string `mapstructure:"heimdall_rpc_endpoint"`
		HeimdallLCDEndpoint string `mapstructure:"heimdall_lcd_endpoint"`
		HeimdallExternalRPC string `mapstructure:"heimdall_external_rpc"`
	}

	// ValDetails struct
	ValDetails struct {
		ValidatorHexAddress  string `mapstructure:"validator_hex_addr"`
		SignerAddress        string `mapstructure:"signer_address"`
		ValidatorName        string `mapstructure:"validator_name"`
		StakeManagerContract string `mapstructure:"stake_manager_contract"`
	}

	// EnableAlerts struct which holds options to enalbe/disable alerts
	EnableAlerts struct {
		EnableTelegramAlerts string `mapstructure:"enable_telegram_alerts"`
		EnableEmailAlerts    string `mapstructure:"enable_email_alerts"`
	}

	// DailyAlert which holds parameters to send validator statu alerts(twice a day)
	DailyAlert struct {
		AlertTime1 string `mapstructure:"alert_time1"`
		AlertTime2 string `mapstructure:"alert_time2"`
	}

	// ChooseAlerts struct
	ChooseAlerts struct {
		BalanceChangeAlerts string `mapstructure:"balance_change_alerts"`
		VotingPowerAlerts   string `mapstructure:"voting_power_alerts"`
		ProposalAlerts      string `mapstructure:"proposal_alerts"`
		BlockDiffAlerts     string `mapstructure:"block_diff_alerts"`
		MissedBlockAlerts   string `mapstructure:"missed_block_alerts"`
		NumPeersAlerts      string `mapstructure:"num_peers_alerts"`
		NodeSyncAlert       string `mapstructure:"node_sync_alert"`
		NodeStatusAlert     string `mapstructure:"node_status_alert"`
	}

	// AlertingThreshold
	AlertingThreshold struct {
		NumPeersThreshold     int64 `mapstructure:"num_peers_threshold"`
		MissedBlocksThreshold int64 `mapstructure:"missed_blocks_threshold"`
		BlockDiffThreshold    int64 `mapstructure:"block_diff_threshold"`
	}

	//Config
	Config struct {
		Endpoints          Endpoints         `mapstructure:"rpc_and_lcd_endpoints"`
		ValDetails         ValDetails        `mapstructure:"validator_details"`
		EnableAlerts       EnableAlerts      `mapstructure:"enable_alerts"`
		DailyAlert         DailyAlert        `mapstructure:"daily_alert"`
		ChooseAlerts       ChooseAlerts      `mapstructure:"choose_alerts"`
		AlertingThresholds AlertingThreshold `mapstructure:"alerting_threholds"`
		Scraper            Scraper           `mapstructure:"scraper"`
		Telegram           Telegram          `mapstructure:"telegram"`
		SendGrid           SendGrid          `mapstructure:"sendgrid"`
		InfluxDB           InfluxDB          `mapstructure:"influxdb"`
		PagerdutyEmail     string            `mapstructure:"pagerduty_email"`
	}
)

//ReadFromFile to read config details using viper
func ReadFromFile() (*Config, error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("./config/")
	v.SetConfigName("config")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error while reading config.toml: %v", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("error unmarshaling config.toml to application config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("error occurred in config validation: %v", err)
	}

	return &cfg, nil
}

//Validate config struct
func (c *Config) Validate(e ...string) error {
	v := validator.New()
	if len(e) == 0 {
		return v.Struct(c)
	}
	return v.StructExcept(c, e...)
}
