package targets

import (
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// ConvertToMatic converts amount into matic and returns
func ConvertToMatic(amount string) string {
	f, _ := strconv.ParseFloat(amount, 64)
	d := f * math.Pow(10, -18)
	bal := fmt.Sprintf("%.4f", d)

	log.Println("heimdall bal : ", bal)

	return bal
}

// convertToCommaSeparated converts value into comma seperated
//for user friendly
func convertToCommaSeparated(amt string) string {
	a, err := strconv.Atoi(amt)
	if err != nil {
		return amt
	}
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", a)
}

// ConvertWeiToEth converts wei into eth value
func ConvertWeiToEth(num *big.Int) string {
	wei := num.String()

	f, _ := strconv.ParseFloat(wei, 64)
	eth := f * math.Pow(10, -18)
	ether := fmt.Sprintf("%.8f", eth)

	log.Println("eth..", ether)

	return ether
}

// ConvertWeiToEth converts wei into eth value
func ConvertValueToEth(num int64) float64 {
	wei := strconv.FormatInt(num, 10)

	f, _ := strconv.ParseFloat(wei, 64)
	eth := f * math.Pow(10, -18)

	return eth
}

// HexToBigInt convert hex value into big int
func HexToBigInt(hex string) (*big.Int, bool) {
	n := new(big.Int)
	n2, err := n.SetString(hex[2:], 16)

	return n2, err
}

// HexToIntConversion converts hex into int format
func HexToIntConversion(hex string) int {
	val := hex[2:]

	n, err := strconv.ParseUint(val, 16, 64)
	if err != nil {
		panic(err)
	}
	n2 := int(n)

	return n2
}

// ConvertNanoSecToMinutes converts nano seconds into minutes
func ConvertNanoSecToMinutes(nanoSec int64) int64 {
	sec := nanoSec / 1e9
	minutes := sec / 60

	return minutes
}

// EncodeToHex encodes b as a hex string with 0x prefix.
func EncodeToHex(b []byte) string {

	return Encode(b)
}

// Encode encodes b as a hex string with 0x prefix.
func Encode(b []byte) string {
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], b)
	return string(enc)
}

// DecodeEthCallResult decodes the eth_call result and resturns the array
func DecodeEthCallResult(resp string) []string {
	s := resp[2:]
	n := len(resp) / 64
	index := 0

	var SubArray []string

	for i := 0; i < n; i++ {
		startIndex := index
		end := startIndex + 64
		sub := s[startIndex:end]
		SubArray = append(SubArray, sub)
		index = end
	}

	return SubArray
}

// GetUserDateFormat to which returns date in a user friendly
func GetUserDateFormat(timeToConvert string) string {
	time, err := time.Parse(time.RFC3339, timeToConvert)
	if err != nil {
		fmt.Println("Error while converting date ", err)
	}
	date := time.Format("Mon Jan _2 15:04:05 2006")
	fmt.Println("Converted time into date format : ", date)
	return date
}

var (
	// Denom of matic
	MaticDenom = "MATIC"
)
