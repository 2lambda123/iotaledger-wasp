// Copyright 2022 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"bytes"
	"io"

	"golang.org/x/crypto/blake2b"

	// <<<<<<< HEAD
	// =======
	// 	"github.com/iotaledger/hive.go/core/kvstore/mapdb"
	// 	"github.com/iotaledger/hive.go/serializer/v2"
	// 	iotago "github.com/iotaledger/iota.go/v3"
	// 	"github.com/iotaledger/wasp/packages/isc"
	// 	"github.com/iotaledger/wasp/packages/kv"
	// >>>>>>> gpa-mp
	"github.com/iotaledger/hive.go/core/kvstore/mapdb"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/buffered"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/trie"
)

type block struct {
	trieRoot             trie.VCommitment
	mutations            *buffered.Mutations
	previousL1Commitment *L1Commitment
}

var _ Block = &block{}

func BlockFromBytes(blockBytes []byte) (*block, error) {
	buf := bytes.NewBuffer(blockBytes)

	trieRoot, err := trie.ReadVectorCommitment(buf)
	if err != nil {
		return nil, err
	}

	muts := buffered.NewMutations()
	err = muts.Read(buf)
	if err != nil {
		return nil, err
	}

	var hasPrevL1Commitment bool
	if hasPrevL1Commitment, err = codec.DecodeBool(buf.Next(1)); err != nil {
		return nil, err
	}
	var prevL1Commitment *L1Commitment
	if hasPrevL1Commitment {
		prevL1Commitment = new(L1Commitment)
		err = prevL1Commitment.Read(buf)
		if err != nil {
			return nil, err
		}
	}

	return &block{
		trieRoot:             trieRoot,
		mutations:            muts,
		previousL1Commitment: prevL1Commitment,
	}, nil
}

func (b *block) Mutations() *buffered.Mutations {
	return b.mutations
}

func (b *block) MutationsReader() kv.KVStoreReader {
	return buffered.NewBufferedKVStoreForMutations(
		kv.NewHiveKVStoreReader(mapdb.NewMapDB()),
		b.mutations,
	)
}

func (b *block) TrieRoot() trie.VCommitment {
	return b.trieRoot
}

func (b *block) PreviousL1Commitment() *L1Commitment {
	return b.previousL1Commitment
}

func (b *block) essenceBytes() []byte {
	var w bytes.Buffer
	b.writeEssence(&w)
	return w.Bytes()
}

func (b *block) writeEssence(w io.Writer) {
	w.Write(b.Mutations().Bytes())

	w.Write(codec.EncodeBool(b.PreviousL1Commitment() != nil))
	if b.PreviousL1Commitment() != nil {
		w.Write(b.PreviousL1Commitment().Bytes())
	}
}

func (b *block) Bytes() []byte {
	var w bytes.Buffer
	w.Write(b.TrieRoot().Bytes())
	b.writeEssence(&w)
	return w.Bytes()
}

func (b *block) Hash() BlockHash {
	return BlockHashFromData(b.essenceBytes())
}

func (b *block) L1Commitment() *L1Commitment {
	return newL1Commitment(b.TrieRoot(), b.Hash())
}

func (b *block) GetHash() (ret BlockHash) {
	r := blake2b.Sum256(b.essenceBytes())
	copy(ret[:BlockHashSize], r[:BlockHashSize])
	return
}

func (b *block) Equals(other Block) bool {
	return b.GetHash().Equals(other.Hash())
}
