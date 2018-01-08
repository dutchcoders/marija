package blockchain_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/qedus/blockchain"
)

const (
	blockHash = "0000000000000bae09a7a393a8acded75aa67e46cb81f7acaa5ad94f9eacd103"
)

func TestRequestBlock(t *testing.T) {
	bc := blockchain.New(http.DefaultClient)
	block := &blockchain.Block{Hash: blockHash}
	if err := bc.Request(block); err != nil {
		t.Fatal(err)
	}

	if block.Fee != 200000 {
		t.Fatal("fee not 200000")
	}

	if block.TransactionCount != 22 {
		t.Fatal("transaction count not 22")
	}

	if len(block.Transactions) != 22 {
		t.Fatal("transactions length not 22")
	}

	if block.Transactions[0].Hash != "5b09bbb8d3cb2f8d4edbcf30664419fb7c9deaeeb1f62cb432e7741c80dbe5ba" {
		t.Fatal("first transaction hash incorrect")
	}

	if !block.Transactions[0].IsCoinbase() {
		t.Fatal("not coinbase transaction")
	}
}

func TestRequestLatestBlock(t *testing.T) {

	bc := blockchain.New(http.DefaultClient)
	block := &blockchain.LatestBlock{}
	if err := bc.Request(block); err != nil {
		t.Fatal(err)
	}

	if time.Unix(block.Time, 0).Before(time.Now().Add(-30 * time.Minute)) {
		t.Fatalf("latest block too old at %s minutes and now is %s",
			time.Unix(block.Time, 0), time.Now())
	}

	if len(block.TransactionIndexes) < 1 {
		t.Fatal("no transactions in latest block")
	}
}

func TestRequestBlockHeight(t *testing.T) {
	bc := blockchain.New(http.DefaultClient)
	bh := &blockchain.BlockHeight{Height: 285180}
	if err := bc.Request(bh); err != nil {
		t.Fatal(err)
	}

	/* It appears that Blockchain.info has fogotten which blocks were previously
	       orphan blocks. Therefore this test no longer works.
		if len(bh.Blocks) != 2 {
			t.Fatal("should be two blocks")
		}

		if bh.Blocks[1].MainChain {
			t.Fatal("this block should not on the main chain")
		}
	*/

	if !bh.Blocks[0].MainChain {
		t.Fatal("this block should be on the main chain")
	}
}
