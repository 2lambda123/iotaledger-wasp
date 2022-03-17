package dbkeys

import (
	"bytes"
	"encoding/hex"

	"github.com/iotaledger/wasp/packages/database/textdb"
)

const (
	ObjectTypeDBSchemaVersion byte = iota
	ObjectTypeChainRecord
	ObjectTypeCommitteeRecord
	ObjectTypeDistributedKeyData
	ObjectTypeStateHash
	ObjectTypeBlock
	ObjectTypeStateVariable
	ObjectTypeNodeIdentity
	ObjectTypeBlobCache
	ObjectTypeBlobCacheTTL
	ObjectTypeTrustedPeer
)

// MakeKey makes key within the partition. It consists to one byte for object type
// and arbitrary byte fragments concatenated together
func MakeKey(objType byte, keyBytes ...[]byte) []byte {
	var buf bytes.Buffer
	buf.WriteByte(objType)
	for _, b := range keyBytes {
		buf.Write(b)
	}
	base58 := hex.EncodeToString(buf.Bytes())
	json, _ := textdb.GetMarshaller().Marshal(base58)
	return json
}
