package blockchain

import (
	"bytes"

	"github.com/denisskin/bin"
	"github.com/likecoin-pro/likecoin/config"
	"github.com/likecoin-pro/likecoin/crypto"
)

type BlockHeader struct {
	Version    int       `json:"version"`       // version
	Num        uint64    `json:"height"`        // number of block in blockchain
	Timestamp  int64     `json:"timestamp"`     // timestamp of block in µsec
	PrevHash   bin.Bytes `json:"previous_hash"` // hash of previous block
	MerkleRoot bin.Bytes `json:"merkle_root"`   // merkle hash of transactions

	// miner sign
	Nonce uint64            `json:"nonce"`     //
	Miner *crypto.PublicKey `json:"miner"`     // pub-key of miner
	Sign  bin.Bytes         `json:"signature"` // miner-node sign
}

func (b *BlockHeader) Hash() []byte {
	return bin.Hash256(
		b.Version,
		b.Num,
		b.Timestamp,
		b.PrevHash,
		b.MerkleRoot,
		b.Nonce,
		b.Miner,
	)
}

func (b *BlockHeader) Verify(pre *BlockHeader) error {
	hash := b.Hash()
	if b.Num == 0 && bytes.Equal(hash, genesisBlockHeaderHash) { // is genesis
		return ErrInvalidGenesisBlock
	}
	if pre != nil {
		if b.Num != pre.Num+1 {
			return ErrInvalidNum
		}
		if !bytes.Equal(b.PrevHash, pre.Hash()) {
			return ErrInvalidPrevHash
		}
	}
	if b.Miner.Empty() {
		return ErrEmptyNodeKey
	}
	if !b.Miner.Equal(config.MasterPublicKey) {
		return ErrInvalidNodeKey
	}
	if !b.Miner.Verify(hash, b.Sign) {
		return ErrInvalidSign
	}
	return nil
}
