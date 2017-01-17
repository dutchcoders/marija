package blockchain

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// NewAddress allows you to create a new address in your blockchain.info
// wallet. The Label parameter is optional and if set will label the new
// address. The struct is passed to Request and Address and/or Error parameters
// are populated when Request returns.
type NewAddress struct {
	// Optional request parameter
	Label string

	// Response
	Address string

	// Error Response
	Error string
}

func (na *NewAddress) newAddressURL(bc *BlockChain) (string, error) {
	if bc.GUID == "" {
		return "", errors.New("BlockChain.GUID not set")
	}
	if bc.Password == "" {
		return "", errors.New("BlockChain.Password not set")
	}

	v := url.Values{}
	v.Set("password", bc.Password)

	if bc.SecondPassword != "" {
		v.Set("second_password", bc.SecondPassword)
	}
	if bc.APICode == "" {
		v.Set("api_code", bc.APICode)
	}
	if na.Label != "" {
		v.Set("label", na.Label)
	}

	return fmt.Sprintf("%s/merchant/%s/new_address?%s",
		rootURL, bc.GUID, v.Encode()), nil
}

func (na *NewAddress) load(bc *BlockChain) error {
	url, err := na.newAddressURL(bc)
	if err != nil {
		return err
	}
	return bc.httpGetJSON(url, na)
}

// SendPayment struct allows you to send a payment from a blockchain.info
// wallet. ToAddress and Amount are required parameters. FromAddress, Shared
// Fee and Note are optional parameters. Check parameter Error for error
// respnoses from the blockchain.info wallet.
type SendPayment struct {
	//  Required request parameters.
	ToAddress string
	Amount    int64

	// Optional request parameters.
	FromAddress string
	Shared      bool
	Fee         int64
	Note        string

	// Response
	Message         string
	TransactionHash string `json:"tx_hash"`
	Notice          string

	// Error response
	Error string
}

func (sp *SendPayment) sendPaymentURL(bc *BlockChain) (string, error) {
	if bc.GUID == "" {
		return "", errors.New("BlockChain.GUID not set")
	}
	if bc.Password == "" {
		return "", errors.New("BlockChain.Password not set")
	}
	if sp.ToAddress == "" {
		return "", errors.New("SendPayment.ToAddress not set")
	}
	if sp.Amount == 0 {
		return "", errors.New("SendPayment.Amount not set")
	}

	v := url.Values{}
	v.Set("password", bc.Password)

	if bc.SecondPassword != "" {
		v.Set("second_password", bc.SecondPassword)
	}
	if bc.APICode == "" {
		v.Set("api_code", bc.APICode)
	}
	v.Set("to", sp.ToAddress)
	v.Set("amount", strconv.FormatInt(sp.Amount, 10))
	if sp.FromAddress != "" {
		v.Set("from", sp.FromAddress)
	}
	if sp.Fee != 0 {
		v.Set("fee", strconv.FormatInt(sp.Fee, 10))
	}
	if sp.Note != "" {
		v.Set("note", sp.Note)
	}

	return fmt.Sprintf("%s/merchant/%s/payment?%s",
		rootURL, bc.GUID, v.Encode()), nil
}

func (sp *SendPayment) load(bc *BlockChain) error {
	url, err := sp.sendPaymentURL(bc)
	if err != nil {
		return err
	}
	return bc.httpGetJSON(url, sp)
}

type AddressList struct {
	// Optional request parameter
	Confirmations int

	Addresses []struct {
		Balance       int64
		Address       string
		Label         string
		TotalReceived int64 `json:"total_received"`
	}
}

func (al *AddressList) addressListURL(bc *BlockChain) (string, error) {
	if bc.GUID == "" {
		return "", errors.New("BlockChain.GUID not set")
	}
	if bc.Password == "" {
		return "", errors.New("BlockChain.Password not set")
	}

	v := url.Values{}
	v.Set("password", bc.Password)

	if bc.APICode == "" {
		v.Set("api_code", bc.APICode)
	}

	v.Set("confirmations", strconv.Itoa(al.Confirmations))

	return fmt.Sprintf("%s/merchant/%s/list?%s",
		rootURL, bc.GUID, v.Encode()), nil
}

func (al *AddressList) load(bc *BlockChain) error {
	url, err := al.addressListURL(bc)
	if err != nil {
		return err
	}
	return bc.httpGetJSON(url, al)
}
