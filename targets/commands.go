package targets

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
	"github.com/vitwit/matic-jagar/types"
)

// TelegramAlerting
func TelegramAlerting(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	if strings.ToUpper(strconv.FormatBool(cfg.EnableAlerts.EnableTelegramAlerts)) == "FALSE" {
		return
	}
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.BotToken)
	if err != nil {
		log.Fatalf("Please configure telegram bot token %v:", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	msgToSend := ""

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.Text == "/status" {
			msgToSend = GetStatus(cfg, c)
		} else if update.Message.Text == "/node" {
			msgToSend = NodeStatus(cfg, c)
		} else if update.Message.Text == "/peers" {
			msgToSend = GetPeersCountMsg(cfg, c)
		} else if update.Message.Text == "/balance" {
			msgToSend = GetAccountBal(cfg, c)
		} else if update.Message.Text == "/list" {
			msgToSend = GetHelp()
		} else {
			text := strings.Split(update.Message.Text, "")
			if len(text) != 0 {
				if text[0] == "/" {
					msgToSend = "Command not found do /list to know about available commands"
				} else {
					msgToSend = " "
				}
			}
		}

		log.Printf("[%s] %s", update.Message.From.UserName, msgToSend)

		if msgToSend != " " {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgToSend)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

// GetHelp returns the msg to show for /help
func GetHelp() string {
	msg := "List of available commands\n /status - returns validator status, voting power, current block height " +
		"and network block height\n /peers - returns number of connected peers\n /node - return status of caught-up\n" +
		"/balance - returns the current balance of your account \n /list - list out the available commands"

	return msg
}

// GetPeersCountMsg returns the no of peers for /peers
func GetPeersCountMsg(cfg *config.Config, c client.Client) string {
	var msg string

	count := GetPeersCount(cfg, c) // Get Heimdall No. Of peers
	msg = fmt.Sprintf("No of connected peers on your Heimdall Node is :  %s \n", count)

	return msg
}

// NodeStatus returns the node caught up status /node
func NodeStatus(cfg *config.Config, c client.Client) string {
	var status string

	nodeSync := GetNodeSync(cfg, c) // Getheimdall node sync status
	status = fmt.Sprintf("Your Heimdall validator node is %s \n", nodeSync)

	return status
}

// GetStatus returns the status messages for /status
func GetStatus(cfg *config.Config, c client.Client) string {
	var status string

	valStatus := GetValStatusFromDB(cfg, c)
	if valStatus == "1" {
		valStatus = "voting"
	} else {
		valStatus = "jailed"
	}
	status = fmt.Sprintf("Heimdall Node Status:\nYour validator is currently  %s \n", valStatus)

	valHeight := GetValidatorBlock(cfg, c) // get heimdall validator block height
	status = status + fmt.Sprintf("Validator current block height %s \n", valHeight)

	networkHeight := GetNetworkBlock(cfg, c) // get heimdall network block height
	status = status + fmt.Sprintf("Network current block height %s \n", networkHeight)

	votingPower := GetVotingPowerFromDb(cfg, c) // get heimdall validator voting power
	status = status + fmt.Sprintf("Voting power of your validator is %s \n", votingPower)

	borHeight := GetBorCurrentBlokHeight(cfg, c) // get bor validator block height
	status = status + fmt.Sprintf("\nBor Node :\nValidator current block height %s \n", borHeight)

	spanID := GetBorSpanIDFromDb(cfg, c) // get bor latest span ID
	status = status + fmt.Sprintf("Current span id is %s \n", spanID)

	return status
}

// GetAccountBal returns balance of the corresponding account
func GetAccountBal(cfg *config.Config, c client.Client) string {
	var balanceMsg string

	balance := GetAccountBalWithDenomFromdb(cfg, c) // get heimdall account balance
	balanceMsg = fmt.Sprintf("Heimdall Node : Current balance of your account(%s) is %s \n", cfg.ValDetails.SignerAddress, balance)

	borBalance := GetBorBalanceFromDB(cfg, c) + "ETH" // get bor account balance
	balanceMsg = balanceMsg + fmt.Sprintf("\nBor Node : Current balance of your account(%s) is %s \n", cfg.ValDetails.SignerAddress, borBalance)

	return balanceMsg
}
