# blockchain

[Blockchain.info API](https://blockchain.info/api) Go interface.

[GoDoc](http://godoc.org/github.com/qedus/blockchain) documentation is available here.

Check the test files for example use. However it is along the lines of:

```go

package main

import (
	"fmt"
	"net/http"

	"github.com/qedus/blockchain"
)

func main() {
	bc := blockchain.New(http.DefaultClient)
	address := &blockchain.Address{Address: "an address hash"}
	if err := bc.Request(address); err != nil {
		panic(err)
	}

	// Loop through all the transactions associated with a certain
	// address "an address hash" and print their transaction hashes.
	for {
		tx, err := address.NextTransaction()
		if err == blockchain.TransactionsDone {
			break
		} else if err != nil {
			panic(err)
		}
		fmt.Println(tx.Hash)
	}
}
```

## Blockchain.info API Issues
By creating this client I have noticed a few issues/inconsistencies with the Blockchain.info API. I have taken the stance of not correcting them with this client and instead document them here.

### Coinbase block input has an empty hash
Look at this [second block](https://blockchain.info/block-height/1?format=json) for example. You will notice that the coinbase transaction, i.e. the first transaction, contains an input with an empty hash as so `inputs: [{}]`. It should really be `inputs:[]`. The effect this has on the Go JSON parser is to produce an empty [`Input{}`](https://github.com/qedus/blockchain/blob/master/transaction.go#L15) struct for all coinbase transactions. Note that coinbase transactions should have *no* inputs - not empty ones.

The consequence of this is you cannot use `len(tx.Inputs) == 0` to determine if a transaction is a coinbase transaction. Instead I have created [`tx.IsCoinbase()`](https://github.com/qedus/blockchain/blob/master/transaction.go#L53) which should be used to determine if a transaction is a coinbase one.

### Transaction input and output counts not always correct
Take transaction `2355913fc1a3d71efbc228c69dc5d74340e07b9012377b4b9f6d5522116d0509` in [block 312373](https://blockchain.info/block/312373?format=json) for example. The JSON says there are 14 inputs (`vin_sz=14`) when in fact it only lists 4 in `inputs=[...]`. Therefore `tx.InputCount` and `tx.OutputCount` do not always match `len(tx.Inputs)` and `len(tx.Outputs)`. However if you select the [transaction individually](https://blockchain.info/rawtx/2355913fc1a3d71efbc228c69dc5d74340e07b9012377b4b9f6d5522116d0509) the correct number of inputs are listed in `inputs=[...]`.

It appears that the block API consolidates inputs with the same address where as the individual transaction API does not. This does not appear to be documented. I therefore advise that if you want to iterate through all transaction inputs and outputs you do not rely on `tx.InputCount` or `tx.OutputCount` but instead use `len(tx.Inputs)` and `len(tx.Outputs)` or their `:= range ...` equivalents.

```go
    // Do NOT do this.
    for i := 0; i < tx.InputCount; i++ {
        ...
    }
    
    // Do this instead.
    for i, input := range tx.Inputs {
        ...
    }
```

## Warning
If you run the tests ensure that you use a dummy set of credentials for `wallet_test.go` otherwise the test will spend your money. 
