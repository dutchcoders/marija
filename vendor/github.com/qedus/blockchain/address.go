package blockchain

import (
	"fmt"
	"net/url"
	"strconv"
)

type Address struct {
	Hash160          string
	Address          string
	TransactionCount int64         `json:"n_tx"`
	TotalReceived    int64         `json:"total_received"`
	TotalSent        int64         `json:"total_sent"`
	FinalBalance     int64         `json:"final_balance"`
	Transactions     []Transaction `json:"txs"`

	// These are used for the NextTransaction iterator.
	bc               *BlockChain
	txOffset         int
	txPosition       int
	txLimit          int
	TxSortDescending bool
}

func (a *Address) NextTransaction() (Transaction, error) {
	if a.txPosition < len(a.Transactions) {
		a.txPosition = a.txPosition + 1
		return a.Transactions[a.txPosition-1], nil
	}

	if len(a.Transactions) < a.txLimit {
		return Transaction{}, IterDone
	}
	a.Transactions = nil
	if err := a.load(a.bc); err != nil {
		return Transaction{}, err
	}
	return a.NextTransaction()
}

func (a *Address) addressURL() string {
	v := url.Values{}
	v.Set("format", "json")
	if a.TxSortDescending {
		v.Set("sort", "0")
	} else {
		v.Set("sort", "1")
	}
	v.Set("offset", strconv.Itoa(a.txOffset))
	v.Set("limit", strconv.Itoa(a.txLimit))
	return fmt.Sprintf("%s/address/%s?%s", rootURL, a.Address, v.Encode())
}

func (a *Address) load(bc *BlockChain) error {
	a.bc = bc
	if a.txLimit == 0 {
		a.txLimit = maxTransactionLimit
	}
	url := a.addressURL()
	if err := bc.httpGetJSON(url, a); err != nil {
		return err
	}
	a.txOffset = a.txOffset + a.txLimit
	a.txPosition = 0
	return nil
}
