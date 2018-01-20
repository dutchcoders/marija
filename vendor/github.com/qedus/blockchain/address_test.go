package blockchain_test

import (
	"net/http"
	"testing"

	"github.com/qedus/blockchain"
)

const (
	largeAddress               = "1dice8EMZmqKvrGE4Qc9bUFf9PX3xaYDp"
	largeAddressTxHashOne      = "598af2a2f8eada73ed7262dfb6d6f99970096a4193fdd8043729d01b03423091"
	largeAddressTxHashFiftyOne = "bc6a0683001e38a3c78e3493c5447ff4a21e26367836034aeb9b9910e3cd437e"

	smallAddress = "1416ArzSr5HGzaTbfQHjkLE5RVBBGw3W13"
)

func TestRequestLargeAddress(t *testing.T) {
	bc := blockchain.New(http.DefaultClient)

	address := &blockchain.Address{Address: largeAddress,
		TxSortDescending: false}
	if err := bc.Request(address); err != nil {
		t.Fatal(err)
	}

	if len(address.Transactions) != 50 {
		t.Fatalf("tx count not 50")
	}

	tx, err := address.NextTransaction()
	if err != nil {
		t.Fatal(err)
	}
	if tx.Hash != largeAddressTxHashOne {
		t.Fatalf("tx hash incorrect %s but should be %s",
			tx.Hash, largeAddressTxHashOne)
	}

	for i := 0; i < 49; i++ {
		tx, err = address.NextTransaction()
		if err != nil {
			t.Fatal(err)
		}
	}

	// Check the address iterator goes to the server again.
	tx, err = address.NextTransaction()
	if err != nil {
		t.Fatal(err)
	}
	if tx.Hash != largeAddressTxHashFiftyOne {
		t.Fatalf("tx hash incorrect %s but should be %s",
			tx.Hash, largeAddressTxHashFiftyOne)
	}
	if len(address.Transactions) != 50 {
		t.Fatalf("tx count not 50")
	}
}

func TestRequestSmallAddress(t *testing.T) {
	bc := blockchain.New(http.DefaultClient)
	address := &blockchain.Address{Address: smallAddress}
	if err := bc.Request(address); err != nil {
		t.Fatal(err)
	}

	count := 0
	for {
		_, err := address.NextTransaction()
		if err == blockchain.IterDone {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		count++
	}

	if count != 6 {
		t.Fatalf("expected 6 iterations but got %d", count)
	}
}
