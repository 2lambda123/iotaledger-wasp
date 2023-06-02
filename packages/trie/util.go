package trie

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"

	"golang.org/x/crypto/blake2b"
)

// mustBytes most common way of serialization
func mustBytes(o interface{ Write(w io.Writer) error }) []byte {
	w := new(bytes.Buffer)
	if err := o.Write(w); err != nil {
		panic(err)
	}
	return w.Bytes()
}

func concat(par ...[]byte) []byte {
	w := new(bytes.Buffer)
	for _, p := range par {
		w.Write(p)
	}
	return w.Bytes()
}

// ---------------------------------------------------------------------------
// r/w utility functions

func ReadBytes16(r io.Reader) ([]byte, error) {
	var length uint16
	err := ReadUint16(r, &length)
	if err != nil {
		return nil, err
	}
	if length == 0 {
		return []byte{}, nil
	}
	ret := make([]byte, length)
	_, err = r.Read(ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func WriteBytes16(w io.Writer, data []byte) error {
	if len(data) > math.MaxUint16 {
		panic(fmt.Sprintf("WriteBytes16: too long data (%v)", len(data)))
	}
	err := WriteUint16(w, uint16(len(data)))
	if err != nil {
		return err
	}
	if len(data) != 0 {
		_, err = w.Write(data)
	}
	return err
}

func ReadUint16(r io.Reader, pval *uint16) error {
	var tmp2 [2]byte
	_, err := r.Read(tmp2[:])
	if err != nil {
		return err
	}
	*pval = binary.LittleEndian.Uint16(tmp2[:])
	return nil
}

func WriteUint16(w io.Writer, val uint16) error {
	_, err := w.Write(uint16To2Bytes(val))
	return err
}

func uint16To2Bytes(val uint16) []byte {
	var tmp2 [2]byte
	binary.LittleEndian.PutUint16(tmp2[:], val)
	return tmp2[:]
}

func readByte(r io.Reader) (byte, error) {
	var b [1]byte
	_, err := r.Read(b[:])
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func writeByte(w io.Writer, val byte) error {
	b := []byte{val}
	_, err := w.Write(b)
	return err
}

func ReadBytes(r io.Reader) ([]byte, error) {
	var length uint32
	err := readUint32(r, &length)
	if err != nil {
		return nil, err
	}
	if length == 0 {
		return []byte{}, nil
	}
	ret := make([]byte, length)
	_, err = r.Read(ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func WriteBytes(w io.Writer, data []byte) error {
	if len(data) > math.MaxUint32 {
		panic(fmt.Sprintf("WriteBytes: too long data (%v)", len(data)))
	}
	err := writeUint32(w, uint32(len(data)))
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func uint32To4Bytes(val uint32) []byte {
	var tmp4 [4]byte
	binary.LittleEndian.PutUint32(tmp4[:], val)
	return tmp4[:]
}

func uint32From4Bytes(b []byte) (uint32, error) {
	if len(b) != 4 {
		return 0, errors.New("len(b) != 4")
	}
	return binary.LittleEndian.Uint32(b), nil
}

func mustUint32From4Bytes(b []byte) uint32 {
	ret, err := uint32From4Bytes(b)
	if err != nil {
		panic(err)
	}
	return ret
}

func readUint32(r io.Reader, pval *uint32) error {
	var tmp4 [4]byte
	_, err := r.Read(tmp4[:])
	if err != nil {
		return err
	}
	*pval = mustUint32From4Bytes(tmp4[:])
	return nil
}

func writeUint32(w io.Writer, val uint32) error {
	_, err := w.Write(uint32To4Bytes(val))
	return err
}

func blake2b160(data []byte) (ret [HashSizeBytes]byte) {
	hash, _ := blake2b.New(HashSizeBytes, nil)
	if _, err := hash.Write(data); err != nil {
		panic(err)
	}
	copy(ret[:], hash.Sum(nil))
	return
}

func assertf(cond bool, format string, args ...interface{}) {
	if !cond {
		panic(fmt.Sprintf("assertion failed:: "+format, args...))
	}
}

func assertNoError(err error) {
	assertf(err == nil, "error: %v", err)
}
