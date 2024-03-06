package asa

import (
	"fmt"
	"testing"
)

func TestAmountCents(t *testing.T) {
	tests := []struct {
		m    Money
		want int
	}{
		{
			m:    Money{Amount: "0", Currency: "USD"},
			want: 0,
		},
		{
			m:    Money{Amount: "100", Currency: "USD"},
			want: 10000,
		},
		{
			m:    Money{Amount: "100.1", Currency: "USD"},
			want: 10010,
		},
		{
			m:    Money{Amount: "100.12", Currency: "USD"},
			want: 10012,
		},
		{
			m:    Money{Amount: "100.12", Currency: "USD"},
			want: 10012,
		},
	}
	for _, testCase := range tests {
		t.Run(fmt.Sprintf("%s %s", testCase.m.Amount, testCase.m.Currency), func(t *testing.T) {
			if got := testCase.m.AmountCents(); got != testCase.want {
				t.Errorf("Money.AmountCents() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestMoneyWithCents(t *testing.T) {
	tests := []struct {
		amount   int
		currency string
		want     Money
	}{
		{
			amount:   0,
			currency: "USD",
			want:     Money{Amount: "0", Currency: "USD"},
		},
		{
			amount:   5,
			currency: "USD",
			want:     Money{Amount: "0.05", Currency: "USD"},
		},
		{
			amount:   50,
			currency: "USD",
			want:     Money{Amount: "0.50", Currency: "USD"},
		},
		{
			amount:   53,
			currency: "USD",
			want:     Money{Amount: "0.53", Currency: "USD"},
		},
		{
			amount:   153,
			currency: "USD",
			want:     Money{Amount: "1.53", Currency: "USD"},
		},
		{
			amount:   10000,
			currency: "GBP",
			want:     Money{Amount: "100.00", Currency: "GBP"},
		},
		{
			amount:   10010,
			currency: "USD",
			want:     Money{Amount: "100.10", Currency: "USD"},
		},
		{
			amount:   10012,
			currency: "USD",
			want:     Money{Amount: "100.12", Currency: "USD"},
		},
	}
	for _, testCase := range tests {
		t.Run(fmt.Sprintf("%d %s", testCase.amount, testCase.currency), func(t *testing.T) {
			if got := MoneyWithCents(testCase.amount, testCase.currency); got != testCase.want {
				t.Errorf("MoneyWithCents() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestMul(t *testing.T) {
	tests := []struct {
		m      Money
		factor float64
		want   Money
	}{
		{
			m:      Money{Amount: "0", Currency: "USD"},
			factor: 1.0,
			want:   Money{Amount: "0", Currency: "USD"},
		},
		{
			m:      Money{Amount: "1.00", Currency: "USD"},
			factor: 0.5,
			want:   Money{Amount: "0.50", Currency: "USD"},
		},
		{
			m:      Money{Amount: "1.00", Currency: "GBP"},
			factor: 2.0,
			want:   Money{Amount: "2.00", Currency: "GBP"},
		},
		{
			m:      Money{Amount: "1.00", Currency: "GBP"},
			factor: 0.0,
			want:   Money{Amount: "0", Currency: "GBP"},
		},
		{
			m:      Money{Amount: "100", Currency: "GBP"},
			factor: 10,
			want:   Money{Amount: "1000.00", Currency: "GBP"},
		},
	}
	for _, testCase := range tests {
		t.Run(fmt.Sprintf("%s_times_%f", testCase.m.Amount, testCase.factor), func(t *testing.T) {
			if got := testCase.m.Mul(testCase.factor); got != testCase.want {
				t.Errorf("Money.Mul() = %v, want %v", got, testCase.want)
			}
		})
	}

}
