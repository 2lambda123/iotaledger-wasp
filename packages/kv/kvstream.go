package kv

import (
	"errors"
	"io"
	"os"

	"github.com/iotaledger/wasp/packages/util/rwutil"
)

// Interfaces for writing/reading persistent streams of key/values

// StreamWriter represents an interface specific to write a sequence of key/value pairs
type StreamWriter interface {
	Write(key, value []byte) error
	Stats() (int, int) // num k/v pairs and num bytes so far
}

// StreamIterator is an interface to iterate stream
type StreamIterator interface {
	Iterate(func(k, v []byte) bool) error
}

// BinaryStreamWriter writes stream of k/v pairs in binary format.
// Keys encoding is 'size32' and values is 'size32'
type BinaryStreamWriter struct {
	w         io.Writer
	kvCount   int
	byteCount int
}

func NewBinaryStreamWriter(w io.Writer) *BinaryStreamWriter {
	return &BinaryStreamWriter{w: w}
}

// BinaryStreamWriter implements StreamWriter interface
var _ StreamWriter = &BinaryStreamWriter{}

func (b *BinaryStreamWriter) Write(key, value []byte) error {
	ww := rwutil.NewWriter(b.w)
	ww.WriteUint16(uint16(len(key)))
	ww.WriteN(key)
	ww.WriteUint32(uint32(len(value)))
	ww.WriteN(value)
	b.byteCount += len(key) + len(value) + 6
	b.kvCount++
	return ww.Err
}

func (b *BinaryStreamWriter) Stats() (int, int) {
	return b.kvCount, b.byteCount
}

type BinaryStreamIterator struct {
	r io.Reader
}

func NewBinaryStreamIterator(r io.Reader) *BinaryStreamIterator {
	return &BinaryStreamIterator{r: r}
}

func (b BinaryStreamIterator) Iterate(fun func(k []byte, v []byte) bool) error {
	for {
		rr := rwutil.NewReader(b.r)
		key := make([]byte, rr.ReadUint16())
		if errors.Is(rr.Err, io.EOF) {
			return nil
		}
		rr.ReadN(key)
		value := make([]byte, rr.ReadUint32())
		rr.ReadN(value)
		if rr.Err != nil {
			return rr.Err
		}
		if !fun(key, value) {
			return nil
		}
	}
}

type BinaryStreamFileWriter struct {
	*BinaryStreamWriter
	File *os.File
}

// CreateKVStreamFile create a new k/v file
func CreateKVStreamFile(fname string) (*BinaryStreamFileWriter, error) {
	file, err := os.Create(fname)
	if err != nil {
		return nil, err
	}
	return &BinaryStreamFileWriter{
		BinaryStreamWriter: NewBinaryStreamWriter(file),
		File:               file,
	}, nil
}

type BinaryStreamFileIterator struct {
	*BinaryStreamIterator
	File *os.File
}

// OpenKVStreamFile opens existing file with k/v stream
func OpenKVStreamFile(fname string) (*BinaryStreamFileIterator, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	return &BinaryStreamFileIterator{
		BinaryStreamIterator: NewBinaryStreamIterator(file),
		File:                 file,
	}, nil
}
