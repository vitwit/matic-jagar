package targets

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/alerter"
	"github.com/vitwit/matic-jagar/config"
	db "github.com/vitwit/matic-jagar/influxdb"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
	"github.com/vitwit/matic-jagar/utils"
)

// Proposals is to get all the proposals and stores it in db
// Alerter will send alerts if there is any new proposal or any change in proposal status
func Proposals(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := db.CreateBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	p, err := scraper.GetProposals(ops)
	if err != nil {
		log.Printf("Error in proposals: %v", err)
		return
	}

	for _, proposal := range p.Result {
		validatorVoted := CheckValidatorVoted(proposal.ID, cfg, c)
		validatorDeposited := CheckValidatorDeposited(proposal.ID, cfg, c)
		err = SendVotingPeriodProposalAlerts(cfg, c)
		if err != nil {
			log.Printf("Error while sending voting period proposal alert: %v", err)
		}

		tag := map[string]string{"id": proposal.ID}
		fields := map[string]interface{}{
			"proposal_id":               proposal.ID,
			"content_type":              proposal.Content.Type,
			"content_value_title":       proposal.Content.Value.Title,
			"content_value_description": proposal.Content.Value.Description,
			"proposal_status":           proposal.ProposalStatus,
			"final_tally_result":        proposal.FinalTallyResult,
			"submit_time":               utils.GetUserDateFormat(proposal.SubmitTime),
			"deposit_end_time":          utils.GetUserDateFormat(proposal.DepositEndTime),
			"total_deposit":             proposal.TotalDeposit,
			"voting_start_time":         utils.GetUserDateFormat(proposal.VotingStartTime),
			"voting_end_time":           utils.GetUserDateFormat(proposal.VotingEndTime),
			"validator_voted":           validatorVoted,
			"validator_deposited":       validatorDeposited,
		}
		newProposal := false
		proposalStatus := ""
		q := client.NewQuery(fmt.Sprintf("SELECT * FROM heimdall_proposals WHERE proposal_id = '%s'", proposal.ID), cfg.InfluxDB.Database, "")
		if response, err := c.Query(q); err == nil && response.Error() == nil {
			for _, r := range response.Results {
				if len(r.Series) == 0 {
					newProposal = true
					break
				} else {
					for idx, col := range r.Series[0].Columns {
						if col == "proposal_status" {
							v := r.Series[0].Values[0][idx]
							proposalStatus = fmt.Sprintf("%v", v)
						}
					}
				}
			}

			if newProposal {
				log.Printf("New Proposal Came: %s", proposal.ID)
				_ = db.WriteToInfluxDb(c, bp, "heimdall_proposals", tag, fields)

				if proposal.ProposalStatus == "Rejected" || proposal.ProposalStatus == "Passed" {
					_ = alerter.SendTelegramAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been %s", proposal.ID, proposal.ProposalStatus), cfg)
					_ = alerter.SendEmailAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been = %s", proposal.ID, proposal.ProposalStatus), cfg)
				} else if proposal.ProposalStatus == "VotingPeriod" {
					_ = alerter.SendTelegramAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been moved to %s", proposal.ID, proposal.ProposalStatus), cfg)
					_ = alerter.SendEmailAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been moved to %s", proposal.ID, proposal.ProposalStatus), cfg)
				} else {
					_ = alerter.SendTelegramAlert(fmt.Sprintf("A new proposal "+proposal.Content.Type+" has been added to "+proposal.ProposalStatus+" with proposal id = %s", proposal.ID), cfg)
					_ = alerter.SendEmailAlert(fmt.Sprintf("A new proposal "+proposal.Content.Type+" has been added to "+proposal.ProposalStatus+" with proposal id = %s", proposal.ID), cfg)
				}
			} else {
				q := client.NewQuery(fmt.Sprintf("DELETE FROM heimdall_proposals WHERE id = '%s'", proposal.ID), cfg.InfluxDB.Database, "")
				if response, err := c.Query(q); err == nil && response.Error() == nil {
					log.Printf("Delete proposal %s from heimdall_proposals", proposal.ID)
				} else {
					log.Printf("Failed to delete proposal %s from heimdall_proposals", proposal.ID)
					log.Printf("Reason for proposal deletion failure %v", response)
				}
				log.Printf("Writing the proposal: %s", proposal.ID)
				_ = db.WriteToInfluxDb(c, bp, "heimdall_proposals", tag, fields)
				if proposal.ProposalStatus != proposalStatus {
					if proposal.ProposalStatus == "Rejected" || proposal.ProposalStatus == "Passed" {
						_ = alerter.SendTelegramAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been %s", proposal.ID, proposal.ProposalStatus), cfg)
						_ = alerter.SendEmailAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been = %s", proposal.ID, proposal.ProposalStatus), cfg)
					} else {
						_ = alerter.SendTelegramAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been moved to %s", proposal.ID, proposal.ProposalStatus), cfg)
						_ = alerter.SendEmailAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been moved to %s", proposal.ID, proposal.ProposalStatus), cfg)
					}
				}
			}
		}
	}

	// Calling fucntion to delete deposit proposals
	// which are ended
	err = DeleteDepoitEndProposals(cfg, c, p)
	if err != nil {
		log.Printf("Error while deleting proposals")
	}
}

// CheckValidatorVoted is to check validator voted for the proposal or not and returns the status
func CheckValidatorVoted(proposalID string, cfg *config.Config, c client.Client) string {
	var ops types.HTTPOptions
	ops.Endpoint = cfg.Endpoints.HeimdallLCDEndpoint + "/gov/proposals/" + proposalID + "/votes"
	ops.Method = http.MethodGet

	voters, err := scraper.GetProposalVoters(ops)
	if err != nil {
		log.Printf("Error in proposal voters: %v", err)
		return ""
	}

	// Get id using the signer address
	valID := GetValID(cfg, c)

	validatorVoted := "not voted"
	for _, value := range voters.Result {
		if value.Voter == valID {
			validatorVoted = value.Option
		}
	}
	return validatorVoted
}

// SendVotingPeriodProposalAlerts which send alerts of voting period proposals
func SendVotingPeriodProposalAlerts(cfg *config.Config, c client.Client) error {
	bp, err := db.CreateBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return err
	}

	var ops types.HTTPOptions
	ops.Endpoint = cfg.Endpoints.HeimdallLCDEndpoint + "/gov/proposals?status=voting_period"
	ops.Method = http.MethodGet

	p, err := scraper.GetProposals(ops)
	if err != nil {
		log.Printf("Error in voting period proposals: %v", err)
		return err
	}

	for _, proposal := range p.Result {
		proposalVotesURL := cfg.Endpoints.HeimdallLCDEndpoint + "/gov/proposals/" + proposal.ID + "/votes"
		ops.Endpoint = proposalVotesURL
		ops.Method = http.MethodGet

		voters, err := scraper.GetProposalVoters(ops)
		if err != nil {
			log.Printf("Error in proposal voters: %v", err)
			return err
		}

		// Get id using the signer address
		valID := GetValID(cfg, c)

		var validatorVoted string
		for _, value := range voters.Result {
			if value.Voter == valID {
				validatorVoted = value.Option
			}
		}

		if validatorVoted == "No" {
			now := time.Now().UTC()
			votingEndTime, _ := time.Parse(time.RFC3339, proposal.VotingEndTime)
			timeDiff := now.Sub(votingEndTime)
			log.Println("timeDiff...", timeDiff.Hours())

			var proposalAlertCount = 1
			count := GetVotesProposalAlertsCount(cfg, c, proposal.ID)
			if count != "" {
				pac, err := strconv.Atoi(count)
				if err != nil {
					log.Printf("Error while converting proposal alerts count : %v", err)
					return err
				}
				proposalAlertCount = pac
			}

			if timeDiff.Hours() == 24 && proposalAlertCount <= 1 {
				_ = alerter.SendTelegramAlert(fmt.Sprintf("%s validator has not voted on proposal = %s, No.of hours left to vote is : %f h", cfg.ValDetails.ValidatorName, proposal.ID, timeDiff.Hours()), cfg)
				_ = alerter.SendEmailAlert(fmt.Sprintf("%s validator has not voted on proposal = %s, No.of hours left to vote is : %f h", cfg.ValDetails.ValidatorName, proposal.ID, timeDiff.Hours()), cfg)

				proposalAlertCount = proposalAlertCount + 1
				_ = db.WriteToInfluxDb(c, bp, "heimdall_votes_proposal_alert_count", map[string]string{}, map[string]interface{}{"count": proposalAlertCount, "proposal_id": proposal.ID})

				log.Println("Sent alert of voting period proposals")
			}

			if timeDiff.Hours() == 12 && proposalAlertCount <= 2 {
				_ = alerter.SendTelegramAlert(fmt.Sprintf("%s validator has not voted on proposal = %s, No.of hours left to vote is : %f h", cfg.ValDetails.ValidatorName, proposal.ID, timeDiff.Hours()), cfg)
				_ = alerter.SendEmailAlert(fmt.Sprintf("%s validator has not voted on proposal = %s, No.of hours left to vote is : %f h", cfg.ValDetails.ValidatorName, proposal.ID, timeDiff.Hours()), cfg)

				proposalAlertCount = proposalAlertCount + 1
				_ = db.WriteToInfluxDb(c, bp, "heimdall_votes_proposal_alert_count", map[string]string{}, map[string]interface{}{"count": proposalAlertCount, "proposal_id": proposal.ID})

				log.Println("Sent alert of voting period proposals")
			}
		}
	}
	return nil
}

// CheckValidatorDeposited is to check validator deposited for the proposal or not and returns the status
func CheckValidatorDeposited(proposalID string, cfg *config.Config, c client.Client) string {
	var ops types.HTTPOptions
	proposalURL := cfg.Endpoints.HeimdallLCDEndpoint + "/gov/proposals/" + proposalID + "/deposits"
	ops.Endpoint = proposalURL
	ops.Method = http.MethodGet

	depositors, err := scraper.GetProposalDepositors(ops)
	if err != nil {
		log.Printf("Error in proposal depositors: %v", err)
		return ""
	}

	// Get id using the signer address
	valID := GetValID(cfg, c)

	validateDeposit := "no"
	for _, value := range depositors.Result {
		if value.Depositor == valID && len(value.Amount) != 0 {
			validateDeposit = "yes"
		}
	}
	return validateDeposit
}

// DeleteDepoitEndProposals is to delete proposals from db which are not present in response anymore
func DeleteDepoitEndProposals(cfg *config.Config, c client.Client, p types.Proposals) error {
	var ID string
	found := false
	q := client.NewQuery("SELECT * FROM heimdall_proposals where proposal_status='DepositPeriod'", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx := range r.Series[0].Values {
					proposalID := r.Series[0].Values[idx][7]
					ID = fmt.Sprintf("%v", proposalID)

					for _, proposal := range p.Result {
						if proposal.ID == ID {
							found = true
							break
						} else {
							found = false
						}
					}
					if !found {
						q := client.NewQuery(fmt.Sprintf("DELETE FROM heimdall_proposals WHERE id = '%s'", ID), cfg.InfluxDB.Database, "")
						if response, err := c.Query(q); err == nil && response.Error() == nil {
							log.Printf("Delete proposal %s from heimdall_proposals", ID)
							return err
						}
						log.Printf("Failed to delete proposal %s from heimdall_proposals", ID)
						log.Printf("Reason for proposal deletion failure %v", response)
					}
				}
			}
		}
	}
	return nil
}

// GetVotesProposalAlertsCount returns the count of voting period alerts for particular proposal
func GetVotesProposalAlertsCount(cfg *config.Config, c client.Client, proposalID string) string {
	var count string
	q := client.NewQuery(fmt.Sprintf("SELECT last(count) FROM heimdall_votes_proposal_alert_count WHERE proposal_id = '%s'", proposalID), cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						pc := r.Series[0].Values[0][idx]
						count = fmt.Sprintf("%v", pc)
						break
					}
				}
			}
		}
	}

	return count
}