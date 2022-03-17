package textdb

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/iotaledger/hive.go/kvstore"
	"github.com/iotaledger/hive.go/logger"
	"github.com/iotaledger/hive.go/types"
	"github.com/iotaledger/wasp/packages/parameters"
)

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
	regFile := parameters.GetString(parameters.DatabaseTextFilename)
	if filepath.Ext(regFile) == "yaml" {
		return NewYAMLMarshaller()
	}
	return NewJSONMarshaller()
}

func NewTextKV(log *logger.Logger, filename string) kvstore.KVStore {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND, 0o666)
	defer f.Close()
	if err != nil {
		// log error
		return nil
	}
	fd, err := f.Stat()
	if err != nil {
		// log err
		return nil
	}
	if fd.Size() == 0 {
		os.WriteFile(f.Name(), []byte("{}"), 0o666)
	}
	return &textKV{filename: f.Name(), log: log, Marshaller: GetMarshaller()}
}

// WithRealm is a factory method for using the same underlying storage with a different realm.
func (s *textKV) WithRealm(realm kvstore.Realm) kvstore.KVStore {
	// return &textKV{s.filename, s.log, realm}

	return &textKV{
		filename: s.filename,
		log:      s.log,
		realm:    realm,
	}
}

// Realm returns the configured realm.
func (s *textKV) Realm() kvstore.Realm {
	return s.realm
}

// Shutdown marks the store as shutdown.
func (s *textKV) Shutdown() {
}

// Iterate iterates over all keys and values with the provided prefix. You can pass kvstore.EmptyPrefix to iterate over all keys and values.
func (s *textKV) Iterate(prefix kvstore.KeyPrefix, kvConsumerFunc kvstore.IteratorKeyValueConsumerFunc) error {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		return err
	}
	rec := map[string]interface{}{}
	err = s.Unmarshal(data, &rec)
	if err != nil {
		return err
	}
	strPrefix := s.unmarshalString(prefix)
	for key, value := range rec {
		if strings.HasPrefix(key, strPrefix) {
			val := s.marshalInterface(value)
			kvConsumerFunc(s.marshalString(key), val)
		}
	}
	return nil
}

// IterateKeys iterates over all keys with the provided prefix. You can pass kvstore.EmptyPrefix to iterate over all keys.
func (s *textKV) IterateKeys(prefix kvstore.KeyPrefix, consumerFunc kvstore.IteratorKeyConsumerFunc) error {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		return err
	}
	rec := map[string]interface{}{}
	err = s.Unmarshal(data, &rec)
	if err != nil {
		return err
	}
	for key := range rec {
		if strings.HasPrefix(key, s.unmarshalString(prefix)) {
			consumerFunc(s.marshalString(key))
		}
	}
	return nil
}

func (s *textKV) Clear() error {
	return os.Truncate(s.filename, 0)
}

// Get gets the given key or nil if it doesn't exist or an error if an error occurred.
func (s *textKV) Get(key kvstore.Key) (value kvstore.Value, err error) {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	err = s.Unmarshal(data, &rec)
	if err != nil {
		return nil, err
	}
	val, ok := rec[s.unmarshalString(key)]
	if !ok {
		return nil, kvstore.ErrKeyNotFound
	}
	return s.Marshal(val)
}

// Set sets the given key and value.
func (s *textKV) Set(key kvstore.Key, value kvstore.Value) error {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		return err
	}
	rec := map[string]interface{}{}
	err = s.Unmarshal(data, &rec)
	if err != nil {
		return err
	}
	var newVal interface{}
	err = s.Unmarshal(value, &newVal)
	if err != nil {
		return err
	}
	keyStr := s.unmarshalString(key)
	rec[keyStr] = newVal
	data, err = s.Marshal(rec)
	if err != nil {
		return err
	}
	return os.WriteFile(s.filename, data, 0o666)
}

// Has checks whether the given key exists.
func (s *textKV) Has(key kvstore.Key) (bool, error) {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		return false, err
	}
	rec := map[string]interface{}{}
	err = s.Unmarshal(data, &rec)
	if err != nil {
		return false, err
	}
	_, ok := rec[s.unmarshalString(key)]
	return ok, nil
}

// Delete deletes the entry for the given key.
func (s *textKV) Delete(key kvstore.Key) error {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		return err
	}
	rec := map[string]interface{}{}
	err = s.Unmarshal(data, &rec)
	if err != nil {
		return err
	}
	delete(rec, s.unmarshalString(key))
	data, err = s.Marshal(rec)
	if err != nil {
		return err
	}
	return os.WriteFile(s.filename, data, 0o666)
}

// DeletePrefix deletes all the entries matching the given key prefix.
func (s *textKV) DeletePrefix(prefix kvstore.KeyPrefix) error {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		return err
	}
	rec := map[string]interface{}{}
	err = s.Unmarshal(data, &rec)
	if err != nil {
		return err
	}
	for key := range rec {
		if strings.HasPrefix(key, s.unmarshalString(prefix)) {
			delete(rec, key)
		}
	}
	data, err = s.Marshal(rec)
	if err != nil {
		return err
	}
	return os.WriteFile(s.filename, data, 0o666)
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

func (s *textKV) unmarshalString(val []byte) string {
	var ret string
	err := s.Unmarshal(val, &ret)
	if err != nil {
		errf := fmt.Errorf("Error unmarshalling string: %w", err)
		panic(errf)
	}
	return ret
}

func (s *textKV) unmarshalInterface(val []byte) interface{} {
	var ret interface{}
	err := s.Unmarshal(val, &ret)
	if err != nil {
		panic(err)
	}
	return ret
}

func (s *textKV) marshalString(val string) []byte {
	data, err := s.Marshal(val)
	if err != nil {
		panic(err)
	}
	return data
}

func (s *textKV) marshalInterface(val interface{}) []byte {
	data, err := s.Marshal(val)
	if err != nil {
		panic(err)
	}
	return data
}

type kvtupple struct {
	key   kvstore.Key
	value kvstore.Value
}

type batchedMutations struct {
	sync.Mutex
	kvStore          *textKV
	setOperations    map[string]kvstore.Value
	deleteOperations map[string]types.Empty

	sets    []kvtupple
	deletes []kvtupple
}

func (b *batchedMutations) Set(key kvstore.Key, value kvstore.Value) error {
	b.Lock()
	defer b.Unlock()

	strKey := b.kvStore.unmarshalString(key)
	delete(b.deleteOperations, strKey)
	b.setOperations[strKey] = value

	return nil
}

func (b *batchedMutations) Delete(key kvstore.Key) error {
	b.Lock()
	defer b.Unlock()

	strKey := b.kvStore.unmarshalString(key)
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
		err := b.kvStore.Set(b.kvStore.marshalString(key), value)
		if err != nil {
			return err
		}
	}

	for key := range b.deleteOperations {
		err := b.kvStore.Delete(b.kvStore.marshalString(key))
		if err != nil {
			return err
		}
	}
	return nil
}
