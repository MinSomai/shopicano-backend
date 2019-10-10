package utils

import (
	"fmt"
	"github.com/braintree-go/braintree-go"
)

func IntToDecimal(v int, div int) (*braintree.Decimal, error) {
	var x = float64(v) / float64(div)
	dec := braintree.NewDecimal(0, 0)
	if err := dec.UnmarshalText([]byte(fmt.Sprintf("%.2f", x))); err != nil {
		return nil, err
	}
	return dec, nil
}
