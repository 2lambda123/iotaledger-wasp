package util

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"io"
	"math"

	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
)

var errInsufficientBytes = errors.New("insufficient bytes")

func Read(r io.Reader, data []byte) error {
	n, err := r.Read(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return errInsufficientBytes
	}
	return nil
}

func Write(w io.Writer, data []byte) error {
	n, err := w.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return errInsufficientBytes
	}
	return nil
}

//////////////////// byte \\\\\\\\\\\\\\\\\\\\

func ReadByte(r io.Reader) (byte, error) {
	var b [1]byte
	err := Read(r, b[:])
	return b[0], err
}

func WriteByte(w io.Writer, val byte) error {
	return Write(w, []byte{val})
}

//////////////////// uint8 \\\\\\\\\\\\\\\\\\\\

func Uint8ToBytes(val uint8) []byte {
	return []byte{val}
}

func Uint8FromBytes(b []byte) (uint8, error) {
	if len(b) != 1 {
		return 0, errors.New("len(b) != 1")
	}
	return b[0], nil
}

func MustUint8FromBytes(b []byte) uint8 {
	val, err := Uint8FromBytes(b)
	if err != nil {
		panic(err)
	}
	return val
}

func ReadUint8(r io.Reader) (uint8, error) {
	var b [1]byte
	err := Read(r, b[:])
	return b[0], err
}

func WriteUint8(w io.Writer, val uint8) error {
	return Write(w, []byte{val})
}

//////////////////// uint16 \\\\\\\\\\\\\\\\\\\\

func Uint16ToBytes(val uint16) []byte {
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], val)
	return b[:]
}

func Uint16FromBytes(b []byte) (uint16, error) {
	if len(b) != 2 {
		return 0, errors.New("len(b) != 2")
	}
	return binary.LittleEndian.Uint16(b), nil
}

func MustUint16FromBytes(b []byte) uint16 {
	val, err := Uint16FromBytes(b)
	if err != nil {
		panic(err)
	}
	return val
}

func ReadUint16(r io.Reader) (uint16, error) {
	var b [2]byte
	err := Read(r, b[:])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(b[:]), nil
}

func WriteUint16(w io.Writer, val uint16) error {
	return Write(w, Uint16ToBytes(val))
}

//////////////////// int32 \\\\\\\\\\\\\\\\\\\\

func Int32ToBytes(val int32) []byte {
	return Uint32ToBytes(uint32(val))
}

func Int32FromBytes(b []byte) (int32, error) {
	ret, err := Uint32FromBytes(b)
	return int32(ret), err
}

func MustInt32FromBytes(b []byte) int32 {
	return int32(MustUint32FromBytes(b))
}

func ReadInt32(r io.Reader) (int32, error) {
	val, err := ReadUint32(r)
	return int32(val), err
}

func WriteInt32(w io.Writer, val int32) error {
	return WriteUint32(w, uint32(val))
}

//////////////////// uint32 \\\\\\\\\\\\\\\\\\\\

func Uint32ToBytes(val uint32) []byte {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], val)
	return b[:]
}

func Uint32FromBytes(b []byte) (uint32, error) {
	if len(b) != 4 {
		return 0, errors.New("len(b) != 4")
	}
	return binary.LittleEndian.Uint32(b), nil
}

func MustUint32FromBytes(b []byte) uint32 {
	val, err := Uint32FromBytes(b)
	if err != nil {
		panic(err)
	}
	return val
}

func ReadUint32(r io.Reader) (uint32, error) {
	var b [4]byte
	err := Read(r, b[:])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(b[:]), nil
}

func WriteUint32(w io.Writer, val uint32) error {
	return Write(w, Uint32ToBytes(val))
}

//////////////////// int64 \\\\\\\\\\\\\\\\\\\\

func Int64ToBytes(val int64) []byte {
	return Uint64ToBytes(uint64(val))
}

func Int64FromBytes(b []byte) (int64, error) {
	ret, err := Uint64FromBytes(b)
	return int64(ret), err
}

func MustInt64FromBytes(b []byte) int64 {
	return int64(MustUint64FromBytes(b))
}

func ReadInt64(r io.Reader) (int64, error) {
	val, err := ReadUint64(r)
	return int64(val), err
}

func WriteInt64(w io.Writer, val int64) error {
	return WriteUint64(w, uint64(val))
}

//////////////////// uint64 \\\\\\\\\\\\\\\\\\\\

func Uint64ToBytes(val uint64) []byte {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], val)
	return b[:]
}

func Uint64FromBytes(b []byte) (uint64, error) {
	if len(b) != 8 {
		return 0, errors.New("len(b) != 8")
	}
	return binary.LittleEndian.Uint64(b), nil
}

func MustUint64FromBytes(b []byte) uint64 {
	val, err := Uint64FromBytes(b)
	if err != nil {
		panic(err)
	}
	return val
}

func ReadUint64(r io.Reader) (uint64, error) {
	var b [8]byte
	err := Read(r, b[:])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(b[:]), nil
}

func WriteUint64(w io.Writer, val uint64) error {
	return Write(w, Uint64ToBytes(val))
}

//////////////////// bytes \\\\\\\\\\\\\\\\\\\\

func ReadBytes(r io.Reader) ([]byte, error) {
	length, err := ReadSize32(r)
	if err != nil {
		return nil, err
	}
	if length == 0 {
		return []byte{}, nil
	}
	ret := make([]byte, length)
	err = Read(r, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func WriteBytes(w io.Writer, data []byte) error {
	size := len(data)
	if size > math.MaxUint32 {
		panic("data size overflow")
	}
	err := WriteSize32(w, uint32(size))
	if err != nil {
		return err
	}
	if size != 0 {
		return Write(w, data)
	}
	return nil
}

//////////////////// bool \\\\\\\\\\\\\\\\\\\\

func ReadBool(r io.Reader) (bool, error) {
	var b [1]byte
	err := Read(r, b[:])
	if err != nil {
		return false, err
	}
	if (b[0] & 0xfe) != 0x00 {
		return false, errors.New("ReadBool: unexpected value")
	}
	return b[0] != 0, nil
}

func WriteBool(w io.Writer, cond bool) error {
	var b [1]byte
	if cond {
		b[0] = 1
	}
	err := Write(w, b[:])
	return err
}

///////////// []int as a bit vector \\\\\\\\\\\\\

func WriteIntsAsBits(w io.Writer, ints []int) error {
	max := 0
	for _, b := range ints {
		if max < b {
			max = b
		}
	}
	size := max/8 + 1
	data := make([]byte, size)
	for _, b := range ints {
		var bitMask byte = 1
		bitMask <<= b % 8
		bytePos := b / 8
		data[bytePos] |= bitMask
	}
	return WriteBytes(w, data)
}

func ReadIntsAsBits(r io.Reader) ([]int, error) {
	data, err := ReadBytes(r)
	if err != nil {
		return nil, err
	}
	ints := []int{}
	for bytePos := range data {
		var bitMask byte = 1
		for i := 0; i < 8; i++ {
			if data[bytePos]&bitMask != 0 {
				ints = append(ints, bytePos*8+i)
			}
			bitMask <<= 1
		}
	}
	return ints, nil
}

//////////////////// size32 \\\\\\\\\\\\\\\\\\\\

func decodeSize32(readByte func() (byte, error)) (uint32, error) {
	b, err := readByte()
	if err != nil {
		return 0, err
	}
	if b < 0x80 {
		return uint32(b), nil
	}
	value := uint32(b & 0x7f)

	b, err = readByte()
	if err != nil {
		return 0, err
	}
	if b < 0x80 {
		return value | (uint32(b) << 7), nil
	}
	value |= uint32(b&0x7f) << 7

	b, err = readByte()
	if err != nil {
		return 0, err
	}
	if b < 0x80 {
		return value | (uint32(b) << 14), nil
	}
	value |= uint32(b&0x7f) << 14

	b, err = readByte()
	if err != nil {
		return 0, err
	}
	if b < 0x80 {
		return value | (uint32(b) << 21), nil
	}
	value |= uint32(b&0x7f) << 21

	b, err = readByte()
	if err != nil {
		return 0, err
	}
	if b < 0xf0 {
		return value | (uint32(b) << 28), nil
	}
	return 0, errors.New("size32 overflow")
}

func Size32FromBytes(buf []byte) (uint32, error) {
	return ReadSize32(bytes.NewReader(buf))
}

func Size32ToBytes(s uint32) []byte {
	switch {
	case s < 0x80:
		return []byte{byte(s)}
	case s < 0x4000:
		return []byte{byte(s | 0x80), byte(s >> 7)}
	case s < 0x200000:
		return []byte{byte(s | 0x80), byte((s >> 7) | 0x80), byte(s >> 14)}
	case s < 0x10000000:
		return []byte{byte(s | 0x80), byte((s >> 7) | 0x80), byte((s >> 14) | 0x80), byte(s >> 21)}
	default:
		return []byte{byte(s | 0x80), byte((s >> 7) | 0x80), byte((s >> 14) | 0x80), byte((s >> 21) | 0x80), byte(s >> 28)}
	}
}

func MustSize32FromBytes(b []byte) uint32 {
	size, err := Size32FromBytes(b)
	if err != nil {
		panic(err)
	}
	return size
}

func ReadSize32(r io.Reader) (uint32, error) {
	return decodeSize32(func() (byte, error) {
		return ReadByte(r)
	})
}

func WriteSize32(w io.Writer, value uint32) error {
	return Write(w, Size32ToBytes(value))
}

//////////////////// string \\\\\\\\\\\\\\\\\\\\

func ReadString(r io.Reader) (string, error) {
	ret, err := ReadBytes(r)
	if err != nil {
		return "", err
	}
	return string(ret), err
}

func WriteString(w io.Writer, str string) error {
	return WriteBytes(w, []byte(str))
}

//////////////////// marshaled \\\\\\\\\\\\\\\\\\\\

// ReadMarshaled supports kyber.Point, kyber.Scalar and similar.
func ReadMarshaled(r io.Reader, val encoding.BinaryUnmarshaler) error {
	bin, err := ReadBytes(r)
	if err != nil {
		return err
	}
	return val.UnmarshalBinary(bin)
}

// WriteMarshaled supports kyber.Point, kyber.Scalar and similar.
func WriteMarshaled(w io.Writer, val encoding.BinaryMarshaler) error {
	bin, err := val.MarshalBinary()
	if err != nil {
		return err
	}
	return WriteBytes(w, bin)
}

func WriteBytesToMarshalUtil(data []byte, mu *marshalutil.MarshalUtil) {
	size := uint32(len(data))
	mu.WriteBytes(Size32ToBytes(size)).WriteBytes(data)
}

func ReadBytesFromMarshalUtil(mu *marshalutil.MarshalUtil) ([]byte, error) {
	size, err := decodeSize32(mu.ReadByte)
	if err != nil {
		return nil, err
	}
	ret, err := mu.ReadBytes(int(size))
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func ReaderFromBytes[T interface{ Read(r io.Reader) error }](data []byte, reader T) (T, error) {
	r := bytes.NewBuffer(data)
	err := reader.Read(r)
	if err != nil {
		return reader, err
	}
	if r.Len() != 0 {
		return reader, errors.New("excess bytes")
	}
	return reader, nil
}

func WriterToBytes(writer interface{ Write(w io.Writer) error }) []byte {
	w := new(bytes.Buffer)
	err := writer.Write(w)
	// should never happen when writing to bytes.Buffer
	if err != nil {
		panic(err)
	}
	return w.Bytes()
}

func WriteFromBytes(w io.Writer, bytes interface{ Bytes() []byte }) error {
	return Write(w, bytes.Bytes())
}
