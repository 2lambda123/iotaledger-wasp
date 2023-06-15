package state

import (
	"fmt"
	"io"

	"golang.org/x/crypto/blake2b"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/trie"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/util/rwutil"
)

const BlockHashSize = 20

type BlockHash [BlockHashSize]byte

// L1Commitment represents the data stored as metadata in the anchor output
type L1Commitment struct {
	// root commitment to the state
	trieRoot trie.Hash
	// hash of the essence of the block
	blockHash BlockHash
}

const L1CommitmentSize = trie.HashSizeBytes + BlockHashSize

func BlockHashFromData(data []byte) (ret BlockHash) {
	r := blake2b.Sum256(data)
	copy(ret[:BlockHashSize], r[:BlockHashSize])
	return
}

func BlockHashFromString(blockHashString string) (BlockHash, error) {
	result := BlockHash{}
	slice, err := iotago.DecodeHex(blockHashString)
	if err != nil {
		return result, fmt.Errorf("Error decoding block hash from string %s: %w", blockHashString, err)
	}
	if len(slice) != BlockHashSize {
		return result, fmt.Errorf("Error decoding block hash from string %s: %v bytes obtained; expected %v bytes", blockHashString, len(slice), BlockHashSize)
	}
	copy(result[:], slice)
	return result, nil
}

func newL1Commitment(c trie.Hash, blockHash BlockHash) *L1Commitment {
	return &L1Commitment{
		trieRoot:  c,
		blockHash: blockHash,
	}
}

func (bh BlockHash) String() string {
	return iotago.EncodeHex(bh[:])
}

func (bh BlockHash) Equals(other BlockHash) bool {
	return bh == other
}

func L1CommitmentFromBytes(data []byte) (*L1Commitment, error) {
	return rwutil.ReaderFromBytes(data, new(L1Commitment))
}

func (s *L1Commitment) TrieRoot() trie.Hash {
	return s.trieRoot
}

func (s *L1Commitment) BlockHash() BlockHash {
	return s.blockHash
}

func (s *L1Commitment) Equals(other *L1Commitment) bool {
	return s.blockHash.Equals(other.blockHash) && s.trieRoot.Equals(other.trieRoot)
}

func (s *L1Commitment) Bytes() []byte {
	return rwutil.WriterToBytes(s)
}

func (s *L1Commitment) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	rr.ReadN(s.trieRoot[:])
	rr.ReadN(s.blockHash[:])
	return rr.Err
}

func (s *L1Commitment) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteN(s.trieRoot[:])
	ww.WriteN(s.blockHash[:])
	return nil
}

func (s *L1Commitment) String() string {
	return fmt.Sprintf("<%s;%s>", s.TrieRoot(), s.BlockHash())
}

var L1CommitmentNil = &L1Commitment{}

func init() {
	zs, err := L1CommitmentFromBytes(make([]byte, L1CommitmentSize))
	if err != nil {
		panic(err)
	}
	L1CommitmentNil = zs
}

// PseudoRandL1Commitment is for testing only
func PseudoRandL1Commitment() *L1Commitment {
	d := make([]byte, L1CommitmentSize)
	_, _ = util.NewPseudoRand().Read(d)
	ret, err := L1CommitmentFromBytes(d)
	if err != nil {
		panic(err)
	}
	return ret
}
