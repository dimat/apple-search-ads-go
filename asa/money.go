package asa

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Money is the response to requests for budget amounts in campaigns
//
// https://developer.apple.com/documentation/apple_search_ads/money
type Money struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

// AmountCents returns the whole amount in cents
func (m Money) AmountCents() int {
	// Avoiding any floating point arithmetic
	decimalIndex := strings.Index(m.Amount, ".")

	amount := m.Amount
	decimalPlaces := 0
	if decimalIndex >= 0 {
		amount = m.Amount[:decimalIndex] + m.Amount[decimalIndex+1:]
		decimalPlaces = len(m.Amount) - decimalIndex - 1
	}

	amountInt, err := strconv.Atoi(amount)
	if err != nil {
		return 0 // we trust Apple to give us a valid amount
	}
	return amountInt * int(math.Pow10(2-decimalPlaces))

}

func MoneyWithCents(amount int, currency string) Money {
	amountString := strconv.Itoa(amount)
	// insert decimal point leaving two decimal places, but only if there are enough digits
	if len(amountString) > 2 {
		amountString = amountString[:len(amountString)-2] + "." + amountString[len(amountString)-2:]
	} else if amountString == "0" {
		amountString = "0"
	} else {
		amountString = fmt.Sprintf("0.%02s", amountString)
	}

	return Money{
		Amount:   amountString,
		Currency: currency,
	}
}

func (m Money) Mul(factor float64) Money {
	amountCents := m.AmountCents()
	amountCents = int(float64(amountCents) * factor)
	return MoneyWithCents(amountCents, m.Currency)
}
