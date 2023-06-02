package util

import (
	"encoding"
	"encoding/binary"
	"errors"
	"io"

	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
)

//////////////////// byte \\\\\\\\\\\\\\\\\\\\

func ReadByte(r io.Reader) (byte, error) {
	var b [1]byte
	_, err := r.Read(b[:])
	return b[0], err
}

func WriteByte(w io.Writer, val byte) error {
	b := []byte{val}
	_, err := w.Write(b)
	return err
}

//////////////////// uint8 \\\\\\\\\\\\\\\\\\\\

func Uint8To1Bytes(val uint8) []byte {
	return []byte{val}
}

func Uint8From1Bytes(b []byte) (uint8, error) {
	if len(b) != 1 {
		return 0, errors.New("len(b) != 1")
	}
	return b[0], nil
}

func MustUint8From1Bytes(b []byte) uint8 {
	ret, err := Uint8From1Bytes(b)
	if err != nil {
		panic(err)
	}
	return ret
}

func ReadUint8(r io.Reader, pval *uint8) error {
	var b [1]byte
	if _, err := r.Read(b[:]); err != nil {
		return err
	}
	*pval = b[0]
	return nil
}

func WriteUint8(w io.Writer, val uint8) error {
	_, err := w.Write(Uint8To1Bytes(val))
	return err
}

//////////////////// uint16 \\\\\\\\\\\\\\\\\\\\

func Uint16To2Bytes(val uint16) []byte {
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], val)
	return b[:]
}

func Uint16From2Bytes(b []byte) (uint16, error) {
	if len(b) != 2 {
		return 0, errors.New("len(b) != 2")
	}
	return binary.LittleEndian.Uint16(b), nil
}

func MustUint16From2Bytes(b []byte) uint16 {
	ret, err := Uint16From2Bytes(b)
	if err != nil {
		panic(err)
	}
	return ret
}

func ReadUint16(r io.Reader, pval *uint16) error {
	var b [2]byte
	if _, err := r.Read(b[:]); err != nil {
		return err
	}
	*pval = binary.LittleEndian.Uint16(b[:])
	return nil
}

func WriteUint16(w io.Writer, val uint16) error {
	_, err := w.Write(Uint16To2Bytes(val))
	return err
}

//////////////////// int32 \\\\\\\\\\\\\\\\\\\\

func Int32To4Bytes(val int32) []byte {
	return Uint32To4Bytes(uint32(val))
}

func ReadInt32(r io.Reader, pval *int32) error {
	var b [4]byte
	if _, err := r.Read(b[:]); err != nil {
		return err
	}
	*pval = int32(binary.LittleEndian.Uint32(b[:]))
	return nil
}

//////////////////// uint32 \\\\\\\\\\\\\\\\\\\\

func Uint32To4Bytes(val uint32) []byte {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], val)
	return b[:]
}

func Uint32From4Bytes(b []byte) (uint32, error) {
	if len(b) != 4 {
		return 0, errors.New("len(b) != 4")
	}
	return binary.LittleEndian.Uint32(b), nil
}

func MustUint32From4Bytes(b []byte) uint32 {
	ret, err := Uint32From4Bytes(b)
	if err != nil {
		panic(err)
	}
	return ret
}

func ReadUint32(r io.Reader, pval *uint32) error {
	var b [4]byte
	if _, err := r.Read(b[:]); err != nil {
		return err
	}
	*pval = MustUint32From4Bytes(b[:])
	return nil
}

func WriteUint32(w io.Writer, val uint32) error {
	_, err := w.Write(Uint32To4Bytes(val))
	return err
}

//////////////////// int64 \\\\\\\\\\\\\\\\\\\\

func Int64To8Bytes(val int64) []byte {
	return Uint64To8Bytes(uint64(val))
}

func Int64From8Bytes(b []byte) (int64, error) {
	ret, err := Uint64From8Bytes(b)
	return int64(ret), err
}

func ReadInt64(r io.Reader, pval *int64) error {
	var b [8]byte
	if _, err := r.Read(b[:]); err != nil {
		return err
	}
	*pval = int64(binary.LittleEndian.Uint64(b[:]))
	return nil
}

func WriteInt64(w io.Writer, val int64) error {
	_, err := w.Write(Uint64To8Bytes(uint64(val)))
	return err
}

//////////////////// uint64 \\\\\\\\\\\\\\\\\\\\

func Uint64To8Bytes(val uint64) []byte {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], val)
	return b[:]
}

func Uint64From8Bytes(b []byte) (uint64, error) {
	if len(b) != 8 {
		return 0, errors.New("len(b) != 8")
	}
	return binary.LittleEndian.Uint64(b), nil
}

func MustUint64From8Bytes(b []byte) uint64 {
	ret, err := Uint64From8Bytes(b)
	if err != nil {
		panic(err)
	}
	return ret
}

func ReadUint64(r io.Reader, pval *uint64) error {
	var b [8]byte
	if _, err := r.Read(b[:]); err != nil {
		return err
	}
	*pval = binary.LittleEndian.Uint64(b[:])
	return nil
}

func WriteUint64(w io.Writer, val uint64) error {
	_, err := w.Write(Uint64To8Bytes(val))
	return err
}

//////////////////// bytes \\\\\\\\\\\\\\\\\\\\

func ReadBytes(r io.Reader) ([]byte, error) {
	length, err := ReadSize32(func() (byte, error) { return ReadByte(r) })
	if err != nil {
		return nil, err
	}
	if length == 0 {
		return []byte{}, nil
	}
	ret := make([]byte, length)
	_, err = r.Read(ret)
	return ret, err
}

func WriteBytes(w io.Writer, data []byte) error {
	size := uint32(len(data))
	_, err := w.Write(Size32ToBytes(size))
	if err != nil {
		return err
	}
	if size != 0 {
		_, err = w.Write(data)
	}
	return err
}

//////////////////// bool \\\\\\\\\\\\\\\\\\\\

func ReadBool(r io.Reader, cond *bool) error {
	var b [1]byte
	if _, err := r.Read(b[:]); err != nil {
		return err
	}
	if (b[0] & 0xfe) != 0x00 {
		return errors.New("ReadBool: unexpected value")
	}
	*cond = b[0] != 0
	return nil
}

func WriteBool(w io.Writer, cond bool) error {
	var b [1]byte
	if cond {
		b[0] = 1
	}
	_, err := w.Write(b[:])
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

func ReadSize32(readByte func() (byte, error)) (uint32, error) {
	b, err := readByte()
	if err != nil {
		return 0, err
	}

	// simplest case, continuation bit not set
	if b < 0x80 {
		return uint32(b), nil
	}

	// first group of 7 bits
	value := uint32(b & 0x7f)
	shift := 7

	for {
		b, err = readByte()
		if err != nil {
			return 0, err
		}

		// continuation bit not set
		if b < 0x80 {
			return value | (uint32(b) << shift), nil
		}

		// next group of 7 bits
		value |= uint32(b&0x7f) << shift
		shift += 7
		if shift >= 32 {
			return 0, errors.New("size32 overflow")
		}
	}
}

func WriteSize32(w io.Writer, value uint32) error {
	_, err := w.Write(Size32ToBytes(value))
	return err
}

func BytesToSize32(buf []byte) uint32 {
	b := uint32(buf[0])
	if (b & 0x80) == 0 {
		return b
	}
	ret := b & 0x7f
	b = uint32(buf[1])
	if (b & 0x80) == 0 {
		return ret | (b << 7)
	}
	ret |= (b & 0x7f) << 7
	b = uint32(buf[2])
	if (b & 0x80) == 0 {
		return ret | (b << 14)
	}
	ret |= (b & 0x7f) << 14
	b = uint32(buf[3])
	if (b & 0x80) == 0 {
		return ret | (b << 21)
	}
	ret |= (b & 0x7f) << 21
	b = uint32(buf[4])
	if (b & 0xf0) == 0 {
		return ret | (b << 28)
	}
	panic("invalid ULEB32")
}

func Size32ToBytes(value uint32) []byte {
	if value < 0x80 {
		return []byte{
			byte(value),
		}
	}
	if value < 0x4000 {
		return []byte{
			byte(value | 0x80),
			byte(value >> 7),
		}
	}
	if value < 0x200000 {
		return []byte{
			byte(value | 0x80),
			byte((value >> 7) | 0x80),
			byte(value >> 14),
		}
	}
	if value < 0x10000000 {
		return []byte{
			byte(value | 0x80),
			byte((value >> 7) | 0x80),
			byte((value >> 14) | 0x80),
			byte(value >> 21),
		}
	}
	return []byte{
		byte(value | 0x80),
		byte((value >> 7) | 0x80),
		byte((value >> 14) | 0x80),
		byte((value >> 21) | 0x80),
		byte(value >> 28),
	}
}

//////////////////// string, uint16 length \\\\\\\\\\\\\\\\\\\\

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

func ReadStringMu(mu *marshalutil.MarshalUtil) (string, error) {
	ret, err := ReadBytesMu(mu)
	return string(ret), err
}

func WriteStringMu(mu *marshalutil.MarshalUtil, str string) {
	WriteBytesMu(mu, []byte(str))
}

//////////////////// marshaled \\\\\\\\\\\\\\\\\\\\

// ReadMarshaled supports kyber.Point, kyber.Scalar and similar.
func ReadMarshaled(r io.Reader, val encoding.BinaryUnmarshaler) error {
	var err error
	var bin []byte
	if bin, err = ReadBytes(r); err != nil {
		return err
	}
	return val.UnmarshalBinary(bin)
}

// WriteMarshaled supports kyber.Point, kyber.Scalar and similar.
func WriteMarshaled(w io.Writer, val encoding.BinaryMarshaler) error {
	var err error
	var bin []byte
	if bin, err = val.MarshalBinary(); err != nil {
		return err
	}
	return WriteBytes(w, bin)
}

func WriteBytesMu(mu *marshalutil.MarshalUtil, data []byte) {
	size := uint32(len(data))
	mu.WriteBytes(Size32ToBytes(size))
	if size != 0 {
		mu.WriteBytes(data)
	}
}

func ReadBytesMu(mu *marshalutil.MarshalUtil) ([]byte, error) {
	size, err := ReadSize32(mu.ReadByte)
	if err != nil {
		return nil, err
	}
	if size == 0 {
		return []byte{}, nil
	}
	return mu.ReadBytes(int(size))
}
