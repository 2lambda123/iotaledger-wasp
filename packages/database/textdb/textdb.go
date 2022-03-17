package textdb

import (
	"encoding/hex"

	"github.com/iotaledger/goshimmer/packages/database"
	"github.com/iotaledger/hive.go/kvstore"
	"github.com/iotaledger/hive.go/logger"
)

type textDB struct {
	kvstore.KVStore
}

func NewTextDB(log *logger.Logger, filename string) (database.DB, error) {
	return &textDB{KVStore: NewTextKV(log, filename)}, nil
}

func (db *textDB) NewStore() kvstore.KVStore {
	return db.KVStore
}

func (db *textDB) Close() error {
	db.KVStore = nil
	return nil
}

func (db *textDB) RequiresGC() bool {
	return false
}

func (db *textDB) GC() error {
	return nil
}

func Base58Text(in []byte, m Marshaller) []byte {
	hexStr := hex.EncodeToString(in)
	newBytes, _ := m.Marshal(hexStr)
	return newBytes
}

func DecodeBase58Text(in []byte, m Marshaller) ([]byte, error) {
	var hexStr string
	err := m.Unmarshal(in, &hexStr)
	if err != nil {
		return nil, err
	}
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	return data, nil
}
