package blockchain_test

import (
	"net/http"
	"testing"

	"github.com/qedus/blockchain"
)

func TestUnconfirmedTransactions(t *testing.T) {
	bc := blockchain.New(http.DefaultClient)
	ut := &blockchain.UnconfirmedTransactions{}
	if err := bc.Request(ut); err != nil {
		t.Fatal(err)
	}

	if len(ut.Transactions) == 0 {
		t.Fatal("no transactions")
	}

	count := 0
	for {

		tx, err := ut.NextTransaction()

		if err == blockchain.IterDone {
			break
		} else if err != nil {
			t.Fatal(err)
		}

		if tx.Hash == "" {
			t.Fatal("no transaction hash")
		}
		count++
	}
	t.Logf("%d unconfirmed transactions", count)
}

func TestTransactionFee(t *testing.T) {
	bc := blockchain.New(http.DefaultClient)
	b := &blockchain.Block{Index: 312373}
	if err := bc.Request(b); err != nil {
		t.Fatal(err)
	}

	feeSum := int64(0)
	for _, tx := range b.Transactions {
		feeSum = feeSum + tx.Fee()
	}

	if feeSum != b.Fee {
		t.Fatalf("fees do not tally feeSum (%d) vs b.Fee (%d)",
			feeSum, b.Fee)
	}
}

func TestTransactionHash(t *testing.T) {
	bc := blockchain.New(http.DefaultClient)
	tx := &blockchain.Transaction{
		Hash: "2355913fc1a3d71efbc228c69dc5d74340e07b9012377b4b9f6d5522116d0509"}
	if err := bc.Request(tx); err != nil {
		t.Fatal(err)
	}

	if len(tx.Inputs) != 14 {
		t.Fatal("should be 14 inputs")
	}

	if len(tx.Outputs) != 1 {
		t.Fatal("should be 1 output")
	}

	if tx.Index != 30187542 {
		t.Fatalf("incorrect index for transaction (%d)", tx.Index)
	}
}

func TestTransactionIndex(t *testing.T) {
	bc := blockchain.New(http.DefaultClient)
	tx := &blockchain.Transaction{Index: 30187542}

	if err := bc.Request(tx); err != nil {
		t.Fatal(err)
	}

	if len(tx.Inputs) != 14 {
		t.Fatal("should be 14 inputs")
	}

	if len(tx.Outputs) != 1 {
		t.Fatal("should be 1 output")
	}

	if tx.Index != 30187542 {
		t.Fatalf("incorrect index for transaction (%d)", tx.Index)
	}
}
