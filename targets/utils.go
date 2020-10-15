package targets

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func ConvertToMatic(amount string) string {
	f, _ := strconv.ParseFloat(amount, 64)
	d := f * math.Pow(10, -18)
	bal := fmt.Sprintf("%.6f", d)

	log.Println("heimdall bal : ", bal)

	return bal
}

func convertToCommaSeparated(amt string) string {
	a, err := strconv.Atoi(amt)
	if err != nil {
		return amt
	}
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", a)
}

func ConvertWeiToEth(num *big.Int) string {
	wei := num.String()

	f, _ := strconv.ParseFloat(wei, 64)
	eth := f * math.Pow(10, -18)
	ether := fmt.Sprintf("%.8f", eth)

	log.Println("eth..", ether)

	return ether
}

func HexToBigInt(hex string) (*big.Int, bool) {
	n := new(big.Int)
	n2, err := n.SetString(hex[2:], 16)

	return n2, err
}

func HexToIntConversion(hex string) int {
	val := hex[2:]

	n, err := strconv.ParseUint(val, 16, 32)
	if err != nil {
		panic(err)
	}
	n2 := int(n)

	return n2
}

func ConvertNanoSecToMinutes(nanoSec int64) int64 {
	sec := nanoSec / 1e9
	minutes := sec / 60

	return minutes
}
