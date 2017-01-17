package blockchain

import (
	"fmt"
)

type Block struct {
	// Request params.
	Hash  string
	Index int64 `json:"block_index"`

	Version          int64  `json:"ver"`
	PreviousBlock    string `json:"prev_block"`
	MerkelRoot       string `json:"mrkl_root"`
	Time             int64
	Bits             int64
	Fee              int64
	Nonce            int64
	TransactionCount int64 `json:"n_tx"`
	Size             int64
	MainChain        bool `json:"main_chain"`
	Height           int64
	ReceivedTime     int64         `json:"received_time"`
	RelayedBy        string        `json:"relayed_by"`
	Transactions     []Transaction `json:"tx"`
}

func blockHashURL(hash string) string {
	return fmt.Sprintf("%s/rawblock/%s", rootURL, hash)
}

func blockIndexURL(index int64) string {
	return fmt.Sprintf("%s/rawblock/%d", rootURL, index)
}

func (b *Block) load(bc *BlockChain) error {
	url := ""
	if b.Hash != "" {
		url = blockHashURL(b.Hash)
	} else {
		url = blockIndexURL(b.Index)
	}
	return bc.httpGetJSON(url, b)
}

type LatestBlock struct {
	Hash               string
	Time               int64
	BlockIndex         int64 `json:"block_index"`
	Height             int64
	TransactionIndexes []int64 `json:"txIndexes"`
}

func latestBlockURL() string {
	return fmt.Sprintf("%s/latestblock", rootURL)
}

func (b *LatestBlock) load(bc *BlockChain) error {
	url := latestBlockURL()
	return bc.httpGetJSON(url, b)
}

type BlockHeight struct {
	Height int64 `json:"-"`
	Blocks []*Block
}

func blockHeightURL(height int64) string {
	return fmt.Sprintf("%s/block-height/%d?format=json", rootURL, height)
}

func (bh *BlockHeight) load(bc *BlockChain) error {
	url := blockHeightURL(bh.Height)
	return bc.httpGetJSON(url, bh)
}
