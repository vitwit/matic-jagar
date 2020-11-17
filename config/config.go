package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

type (
	// Telegram bot details struct
	Telegram struct {
		BotToken string `mapstructure:"tg_bot_token"`
		ChatID   int64  `mapstructure:"tg_chat_id"`
	}

	// SendGrid stores sendgrid API credentials
	SendGrid struct {
		Token        string `mapstructure:"sendgrid_token"`
		EmailAddress string `mapstructure:"email_address"`
	}

	// Scraper defines the time intervals for multiple scrapers to fetch the data
	Scraper struct {
		Rate          string `mapstructure:"rate"`
		Port          string `mapstructure:"port"`
		ValidatorRate string `mapstructure:"validator_rate"`
		ContractRate  string `mapstructure:"contract_rate"`
	}

	// InfluxDB stores influxDB credntials
	InfluxDB struct {
		Port     string `mapstructure:"port"`
		IP       string `mapstructure:"ip"`
		Database string `mapstructure:"database"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	}

	// Endpoints defines multiple API base-urls to fetch the data
	Endpoints struct {
		EthRPCEndpoint      string `mapstructure:"eth_rpc_endpoint"`
		BorRPCEndpoint      string `mapstructure:"bor_rpc_end_point"`
		BorExternalRPC      string `mapstructure:"bor_external_rpc"`
		HeimdallRPCEndpoint string `mapstructure:"heimdall_rpc_endpoint"`
		HeimdallLCDEndpoint string `mapstructure:"heimdall_lcd_endpoint"`
		HeimdallExternalRPC string `mapstructure:"heimdall_external_rpc"`
	}

	// ValDetails stores the validator meta details
	ValDetails struct {
		ValidatorHexAddress  string `mapstructure:"validator_hex_addr"`
		SignerAddress        string `mapstructure:"signer_address"`
		ValidatorName        string `mapstructure:"validator_name"`
		StakeManagerContract string `mapstructure:"stake_manager_contract"`
	}

	// EnableAlerts struct which holds options to enalbe/disable alerts
	EnableAlerts struct {
		EnableTelegramAlerts bool `mapstructure:"enable_telegram_alerts"`
		EnableEmailAlerts    bool `mapstructure:"enable_email_alerts"`
	}

	// RegularStatusAlerts defines time-slots to receive validator status alerts
	RegularStatusAlerts struct {
		AlertTimings []string `mapstructure:"alert_timings"`
	}

	// AlerterPreferences stores individual alert settings to enable/disable particular alert
	AlerterPreferences struct {
		BalanceChangeAlerts string `mapstructure:"balance_change_alerts"`
		VotingPowerAlerts   string `mapstructure:"voting_power_alerts"`
		ProposalAlerts      string `mapstructure:"proposal_alerts"`
		BlockDiffAlerts     string `mapstructure:"block_diff_alerts"`
		MissedBlockAlerts   string `mapstructure:"missed_block_alerts"`
		NumPeersAlerts      string `mapstructure:"num_peers_alerts"`
		NodeSyncAlert       string `mapstructure:"node_sync_alert"`
		NodeStatusAlert     string `mapstructure:"node_status_alert"`
		EthLowBalanceAlert  string `mapstructure:"eth_low_balance_alert"`
	}

	//  AlertingThreshold defines threshold condition for different alert-cases.
	// `Alerter` will send alerts if the condition reaches the threshold
	AlertingThreshold struct {
		NumPeersThreshold     int64   `mapstructure:"num_peers_threshold"`
		MissedBlocksThreshold int64   `mapstructure:"missed_blocks_threshold"`
		BlockDiffThreshold    int64   `mapstructure:"block_diff_threshold"`
		EthBalanceThreshold   float64 `mapstructure:"eth_balance_threshold"`
	}

	// Config defines all the configurations required for the app
	Config struct {
		Endpoints           Endpoints           `mapstructure:"rpc_and_lcd_endpoints"`
		ValDetails          ValDetails          `mapstructure:"validator_details"`
		EnableAlerts        EnableAlerts        `mapstructure:"enable_alerts"`
		RegularStatusAlerts RegularStatusAlerts `mapstructure:"regular_status_alerts"`
		AlerterPreferences  AlerterPreferences  `mapstructure:"alerter_preferences"`
		AlertingThresholds  AlertingThreshold   `mapstructure:"alerting_threholds"`
		Scraper             Scraper             `mapstructure:"scraper"`
		Telegram            Telegram            `mapstructure:"telegram"`
		SendGrid            SendGrid            `mapstructure:"sendgrid"`
		InfluxDB            InfluxDB            `mapstructure:"influxdb"`
		PagerdutyEmail      string              `mapstructure:"pagerduty_email"`
	}
)

// ReadFromFile to read config details using viper
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

// Validate config struct
func (c *Config) Validate(e ...string) error {
	v := validator.New()
	if len(e) == 0 {
		return v.Struct(c)
	}
	return v.StructExcept(c, e...)
}
