package lib

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type TransactionType struct {
	t string
}

// Each transaction is a record of an action that is taken on an order
// or an action by a user such as a deposit or withdraw.
var (
	DepositUSD   = TransactionType{"DepositUSD"}
	DepositGold  = TransactionType{"DepositGold"}
	WithdrawUSD  = TransactionType{"WithdrawUSD"}
	WithdrawGold = TransactionType{"WithdrawGold"}
	BuyGold      = TransactionType{"BuyGold"}
	SellGold     = TransactionType{"SellGold"}
)

// String satisifies fmt.Stringer
func (t TransactionType) String() string {
	return t.t
}

// (TransactionType).Scan receives a series of bytes from the database
// and converts it to a TransactionType{} type with the value of the received string.
// This satisifies the sql.Scanner interface
func (t *TransactionType) Scan(s interface{}) (err error) {
	if t == nil {
		return errors.New("TransactionType: Scan on nil pointer")
	}
	*t = TransactionType{string(s.([]uint8))}
	return nil
}

// (TransactionType).MarshalJSON satisfies the json.Marshaller interface
func (t *TransactionType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.t)), nil
}

func (t *TransactionType) validate(s string) error {
	return nil
}

func (t *TransactionType) assign(s string) {
	t.t = s
}

func (t TransactionType) Value() (driver.Value, error) {
	return t.t, nil
}

// (TransactionType).UnmarshalJSON satisfies the json.Unmarshaller interface
func (t *TransactionType) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err = json.Unmarshal(b, &s); err == nil {
		if err = t.validate(s); err == nil {
			t.assign(s)
		}
	}
	return err
}
