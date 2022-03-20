package textdb

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/iotaledger/hive.go/byteutils"
	"github.com/iotaledger/hive.go/kvstore"
	"github.com/iotaledger/hive.go/logger"
	"github.com/iotaledger/hive.go/types"
	"github.com/iotaledger/wasp/packages/parameters"
	"github.com/mr-tron/base58"
)

const StorePerm = 0o666

// a key/value store implementation that uses text files
type textKV struct {
	sync.RWMutex
	Marshaller
	filename string
	log      *logger.Logger
	realm    []byte
}

type Marshaller interface {
	Marshal(val interface{}) ([]byte, error)
	Unmarshal(buf []byte, v interface{}) error
}

func GetMarshaller() Marshaller {
	regFile := parameters.GetString(parameters.RegistryFile)
	if filepath.Ext(regFile) == "yaml" {
		return YAMLMarshaller()
	}
	return JSONMarshaller()
}

// a key/value store for text storage. Works with both yaml and json.
func NewTextKV(log *logger.Logger, filename string) kvstore.KVStore {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND, StorePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err != nil {
		panic(err)
	}
	fd, err := f.Stat()
	if err != nil {
		panic(err)
	}
	if fd.Size() == 0 {
		err = os.WriteFile(f.Name(), []byte("{}"), StorePerm)
		if err != nil {
			panic(err)
		}
	}
	return &textKV{filename: f.Name(), log: log, Marshaller: GetMarshaller()}
}

// WithRealm is a factory method for using the same underlying storage with a different realm.
func (s *textKV) WithRealm(realm kvstore.Realm) kvstore.KVStore {
	return &textKV{
		filename: s.filename,
		log:      s.log,
		realm:    realm,
	}
}

// Realm returns the configured realm.
func (s *textKV) Realm() kvstore.Realm {
	return byteutils.ConcatBytes(s.realm)
}

// Shutdown marks the store as shutdown.
func (s *textKV) Shutdown() {
}

func (s *textKV) load() (map[string]interface{}, error) {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		return nil, err
	}
	ret := map[string]interface{}{}
	err = s.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// Iterate iterates over all keys and values with the provided prefix. You can pass kvstore.EmptyPrefix to iterate over all keys and values.
func (s *textKV) Iterate(prefix kvstore.KeyPrefix, kvConsumerFunc kvstore.IteratorKeyValueConsumerFunc) error {
	s.RLock()
	rec, err := s.load()
	if err != nil {
		return err
	}
	s.RUnlock()

	copiedElements := make(map[string]interface{})
	keyPrefix := byteutils.ConcatBytesToString(s.realm, prefix)
	for key, value := range rec {
		actualKey, err := base58.Decode(key)
		if err != nil {
			return err
		}
		actualKeyStr := string(actualKey)
		if strings.HasPrefix(actualKeyStr, keyPrefix) {
			copiedElements[actualKeyStr] = value
		}
	}
	for key, value := range copiedElements {
		valB, err := s.Marshal(value)
		if err != nil {
			return err
		}
		if !kvConsumerFunc([]byte(key)[len(s.realm):], valB) {
			break
		}
	}
	return nil
}

// IterateKeys iterates over all keys with the provided prefix. You can pass kvstore.EmptyPrefix to iterate over all keys.
func (s *textKV) IterateKeys(prefix kvstore.KeyPrefix, consumerFunc kvstore.IteratorKeyConsumerFunc) error {
	s.RLock()
	rec, err := s.load()
	if err != nil {
		return err
	}
	s.RUnlock()

	copiedElements := make(map[string]interface{})
	keyPrefix := byteutils.ConcatBytesToString(s.realm, prefix)
	for key := range rec {
		actualKey, err := base58.Decode(key)
		if err != nil {
			return err
		}
		actualKeyStr := string(actualKey)
		if strings.HasPrefix(actualKeyStr, keyPrefix) {
			copiedElements[actualKeyStr] = types.Empty{}
		}
	}

	for key := range copiedElements {
		if !consumerFunc([]byte(key)[len(s.realm):]) {
			break
		}
	}
	return nil
}

// clear the key/value store
func (s *textKV) Clear() error {
	s.Lock()
	defer s.Unlock()
	return os.Truncate(s.filename, 0)
}

// Get gets the given key or nil if it doesn't exist or an error if an error occurred.
func (s *textKV) Get(key kvstore.Key) (value kvstore.Value, err error) {
	s.RLock()
	rec, err := s.load()
	if err != nil {
		return nil, err
	}
	s.RUnlock()

	actualKey := byteutils.ConcatBytes(s.realm, key)
	val, ok := rec[base58.Encode(actualKey)]
	if !ok {
		return nil, kvstore.ErrKeyNotFound
	}
	return s.Marshal(val)
}

// Set sets the given key and value.
func (s *textKV) Set(key kvstore.Key, value kvstore.Value) error {
	s.Lock()
	defer s.Unlock()

	var newVal interface{}
	err := s.Unmarshal(value, &newVal)
	if err != nil {
		return err
	}
	rec, err := s.load()
	if err != nil {
		return err
	}
	actualKey := byteutils.ConcatBytes(s.realm, key)
	rec[base58.Encode(actualKey)] = newVal
	data, err := s.Marshal(rec)
	if err != nil {
		return err
	}
	return os.WriteFile(s.filename, data, StorePerm)
}

// Has checks whether the given key exists.
func (s *textKV) Has(key kvstore.Key) (bool, error) {
	s.RLock()
	defer s.RUnlock()

	rec, err := s.load()
	if err != nil {
		return false, err
	}
	keyStr := base58.Encode(byteutils.ConcatBytes(s.realm, key))
	_, ok := rec[keyStr]
	return ok, nil
}

// Delete deletes the entry for the given key.
func (s *textKV) Delete(key kvstore.Key) error {
	s.Lock()
	defer s.Unlock()

	rec, err := s.load()
	if err != nil {
		return err
	}
	keyStr := base58.Encode(byteutils.ConcatBytes(s.realm, key))
	delete(rec, keyStr)
	data, err := s.Marshal(rec)
	if err != nil {
		return err
	}
	return os.WriteFile(s.filename, data, StorePerm)
}

// DeletePrefix deletes all the entries matching the given key prefix.
func (s *textKV) DeletePrefix(prefix kvstore.KeyPrefix) error {
	s.Lock()
	defer s.Unlock()

	rec, err := s.load()
	if err != nil {
		return err
	}
	for key := range rec {
		keyBytes, err := base58.Decode(key)
		if err != nil {
			return err
		}
		keyPrefix := byteutils.ConcatBytesToString(s.realm, prefix)
		if strings.HasPrefix(string(keyBytes), keyPrefix) {
			delete(rec, key)
		}
	}
	data, err := s.Marshal(rec)
	if err != nil {
		return err
	}
	return os.WriteFile(s.filename, data, StorePerm)
}

// Batched returns a BatchedMutations interface to execute batched mutations.
func (s *textKV) Batched() kvstore.BatchedMutations {
	return &batchedMutations{
		kvStore:          s,
		deleteOperations: make(map[string]types.Empty),
		setOperations:    make(map[string]kvstore.Value),
	}
}

// Flush persists all outstanding write operations to disc.
func (s *textKV) Flush() error {
	return nil
}

// Close closes the database file handles.
func (s *textKV) Close() error {
	return nil
}

type batchedMutations struct {
	sync.Mutex
	kvStore          *textKV
	setOperations    map[string]kvstore.Value
	deleteOperations map[string]types.Empty
}

func (b *batchedMutations) Set(key kvstore.Key, value kvstore.Value) error {
	b.Lock()
	defer b.Unlock()

	strKey := string(key)
	delete(b.deleteOperations, strKey)
	b.setOperations[strKey] = value

	return nil
}

func (b *batchedMutations) Delete(key kvstore.Key) error {
	b.Lock()
	defer b.Unlock()

	strKey := string(key)
	delete(b.setOperations, strKey)
	b.deleteOperations[strKey] = types.Void

	return nil
}

func (b *batchedMutations) Cancel() {
	b.Lock()
	defer b.Unlock()

	b.setOperations = make(map[string]kvstore.Value)
	b.deleteOperations = make(map[string]types.Empty)
}

func (b *batchedMutations) Commit() error {
	b.Lock()
	defer b.Unlock()

	for key, value := range b.setOperations {
		err := b.kvStore.Set([]byte(key), value)
		if err != nil {
			return err
		}
	}

	for key := range b.deleteOperations {
		err := b.kvStore.Delete([]byte(key))
		if err != nil {
			return err
		}
	}
	return nil
}
