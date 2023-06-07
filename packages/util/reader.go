package util

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"time"

	"github.com/iotaledger/hive.go/serializer/v2"
	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	iotago "github.com/iotaledger/iota.go/v3"
)

//////////////////// basic size-checked read/write \\\\\\\\\\\\\\\\\\\\

func ReadN(r io.Reader, data []byte) error {
	n, err := r.Read(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return errors.New("incomplete read")
	}
	return nil
}

func WriteN(w io.Writer, data []byte) error {
	n, err := w.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return errors.New("incomplete write")
	}
	return nil
}

//////////////////// bool \\\\\\\\\\\\\\\\\\\\

func ReadBool(r io.Reader) (bool, error) {
	var b [1]byte
	err := ReadN(r, b[:])
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
	err := WriteN(w, b[:])
	return err
}

//////////////////// byte \\\\\\\\\\\\\\\\\\\\

func ReadByte(r io.Reader) (byte, error) {
	var b [1]byte
	err := ReadN(r, b[:])
	return b[0], err
}

func WriteByte(w io.Writer, val byte) error {
	return WriteN(w, []byte{val})
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
	err = ReadN(r, ret)
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
		return WriteN(w, data)
	}
	return nil
}

//////////////////// int8 \\\\\\\\\\\\\\\\\\\\

func ReadInt8(r io.Reader) (int8, error) {
	val, err := ReadUint8(r)
	return int8(val), err
}

func WriteInt8(w io.Writer, val int8) error {
	return WriteUint8(w, uint8(val))
}

//////////////////// int16 \\\\\\\\\\\\\\\\\\\\

func ReadInt16(r io.Reader) (int16, error) {
	val, err := ReadUint16(r)
	return int16(val), err
}

func WriteInt16(w io.Writer, val int16) error {
	return WriteUint16(w, uint16(val))
}

//////////////////// int32 \\\\\\\\\\\\\\\\\\\\

func ReadInt32(r io.Reader) (int32, error) {
	val, err := ReadUint32(r)
	return int32(val), err
}

func WriteInt32(w io.Writer, val int32) error {
	return WriteUint32(w, uint32(val))
}

//////////////////// int64 \\\\\\\\\\\\\\\\\\\\

func ReadInt64(r io.Reader) (int64, error) {
	val, err := ReadUint64(r)
	return int64(val), err
}

func WriteInt64(w io.Writer, val int64) error {
	return WriteUint64(w, uint64(val))
}

//////////////////// size32 \\\\\\\\\\\\\\\\\\\\

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
	return WriteN(w, Size32ToBytes(value))
}

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

//////////////////// uint8 \\\\\\\\\\\\\\\\\\\\

func ReadUint8(r io.Reader) (uint8, error) {
	var b [1]byte
	err := ReadN(r, b[:])
	return b[0], err
}

func WriteUint8(w io.Writer, val uint8) error {
	return WriteN(w, []byte{val})
}

//////////////////// uint16 \\\\\\\\\\\\\\\\\\\\

func ReadUint16(r io.Reader) (uint16, error) {
	var b [2]byte
	err := ReadN(r, b[:])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(b[:]), nil
}

func WriteUint16(w io.Writer, val uint16) error {
	return WriteN(w, Uint16ToBytes(val))
}

//////////////////// uint32 \\\\\\\\\\\\\\\\\\\\

func ReadUint32(r io.Reader) (uint32, error) {
	var b [4]byte
	err := ReadN(r, b[:])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(b[:]), nil
}

func WriteUint32(w io.Writer, val uint32) error {
	return WriteN(w, Uint32ToBytes(val))
}

//////////////////// uint64 \\\\\\\\\\\\\\\\\\\\

func ReadUint64(r io.Reader) (uint64, error) {
	var b [8]byte
	err := ReadN(r, b[:])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(b[:]), nil
}

func WriteUint64(w io.Writer, val uint64) error {
	return WriteN(w, Uint64ToBytes(val))
}

//////////////////// Reader \\\\\\\\\\\\\\\\\\\\

type Reader struct {
	Err error
	r   io.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

func NewBytesReader(data []byte) *Reader {
	return NewReader(bytes.NewBuffer(data))
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

// PushBack returns a pushback writer that allows you to insert data before the stream.
// The Reader will read this data first, and then resume reading from the stream.
// The pushback Writer is only valid for this Reader until it resumes the stream.
func (rr *Reader) PushBack() *Writer {
	pb := &PushBack{rr: rr, r: rr.r, buf: new(bytes.Buffer)}
	rr.r = pb
	return NewWriter(pb.buf)
}

func (rr *Reader) Read(reader interface{ Read(r io.Reader) error }) {
	if rr.Err == nil {
		rr.Err = reader.Read(rr.r)
	}
}

func (rr *Reader) ReadN(ret []byte) {
	if rr.Err == nil {
		rr.Err = ReadN(rr.r, ret)
	}
}

func (rr *Reader) ReadAddress() (ret iotago.Address) {
	addrType := rr.ReadByte()
	if rr.Err != nil {
		return ret
	}
	ret, rr.Err = iotago.AddressSelector(uint32(addrType))
	if rr.Err != nil {
		return ret
	}
	buf := make([]byte, ret.Size())
	buf[0] = addrType
	rr.ReadN(buf[1:])
	if rr.Err != nil {
		return ret
	}
	_, rr.Err = ret.Deserialize(buf, serializer.DeSeriModeNoValidation, nil)
	return ret
}

func (rr *Reader) ReadBool() (ret bool) {
	if rr.Err == nil {
		ret, rr.Err = ReadBool(rr.r)
	}
	return ret
}

//nolint:govet
func (rr *Reader) ReadByte() (ret byte) {
	if rr.Err == nil {
		ret, rr.Err = ReadByte(rr.r)
	}
	return ret
}

func (rr *Reader) ReadBytes() (ret []byte) {
	if rr.Err == nil {
		ret, rr.Err = ReadBytes(rr.r)
	}
	return ret
}

func (rr *Reader) ReadDuration() (ret time.Duration) {
	return time.Duration(rr.ReadInt64())
}

func (rr *Reader) ReadInt8() (ret int8) {
	if rr.Err == nil {
		ret, rr.Err = ReadInt8(rr.r)
	}
	return ret
}

func (rr *Reader) ReadInt16() (ret int16) {
	if rr.Err == nil {
		ret, rr.Err = ReadInt16(rr.r)
	}
	return ret
}

func (rr *Reader) ReadInt32() (ret int32) {
	if rr.Err == nil {
		ret, rr.Err = ReadInt32(rr.r)
	}
	return ret
}

func (rr *Reader) ReadInt64() (ret int64) {
	if rr.Err == nil {
		ret, rr.Err = ReadInt64(rr.r)
	}
	return ret
}

func (rr *Reader) ReadMarshaled(m encoding.BinaryUnmarshaler) {
	buf := rr.ReadBytes()
	if rr.Err == nil {
		if m == nil {
			rr.Err = errors.New("nil unmarshaler")
			return
		}
		rr.Err = m.UnmarshalBinary(buf)
	}
}

type deserializable interface {
	Deserialize([]byte, serializer.DeSerializationMode, interface{}) (int, error)
}

func (rr *Reader) ReadSerialized(s deserializable) {
	data := rr.ReadBytes()
	if rr.Err == nil {
		if s == nil {
			rr.Err = errors.New("nil deserializer")
			return
		}
		var n int
		n, rr.Err = s.Deserialize(data, serializer.DeSeriModeNoValidation, nil)
		if rr.Err == nil && n != len(data) {
			rr.Err = errors.New("incomplete deserialize")
		}
	}
}

func (rr *Reader) ReadSize() (ret int) {
	return int(rr.ReadSize32())
}

func (rr *Reader) ReadSize32() (ret uint32) {
	if rr.Err == nil {
		ret, rr.Err = ReadSize32(rr.r)
	}
	return ret
}

func (rr *Reader) ReadString() (ret string) {
	if rr.Err == nil {
		ret, rr.Err = ReadString(rr.r)
	}
	return ret
}

func (rr *Reader) ReadUint8() (ret uint8) {
	if rr.Err == nil {
		ret, rr.Err = ReadUint8(rr.r)
	}
	return ret
}

func (rr *Reader) ReadUint16() (ret uint16) {
	if rr.Err == nil {
		ret, rr.Err = ReadUint16(rr.r)
	}
	return ret
}

func (rr *Reader) ReadUint32() (ret uint32) {
	if rr.Err == nil {
		ret, rr.Err = ReadUint32(rr.r)
	}
	return ret
}

func (rr *Reader) ReadUint64() (ret uint64) {
	if rr.Err == nil {
		ret, rr.Err = ReadUint64(rr.r)
	}
	return ret
}

//////////////////// Writer \\\\\\\\\\\\\\\\\\\\

type Writer struct {
	Err error
	w   io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

func NewBytesWriter() *Writer {
	return NewWriter(new(bytes.Buffer))
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

func (ww *Writer) Bytes() []byte {
	buf, ok := ww.w.(*bytes.Buffer)
	if !ok {
		panic("writer expects bytes buffer")
	}
	return buf.Bytes()
}

func (ww *Writer) Write(writer interface{ Write(w io.Writer) error }) *Writer {
	if ww.Err == nil {
		ww.Err = writer.Write(ww.w)
	}
	return ww
}

func (ww *Writer) WriteN(val []byte) *Writer {
	if ww.Err == nil {
		ww.Err = WriteN(ww.w, val)
	}
	return ww
}

func (ww *Writer) WriteAddress(val iotago.Address) *Writer {
	if ww.Err == nil {
		buf, _ := val.Serialize(serializer.DeSeriModeNoValidation, nil)
		ww.WriteN(buf)
	}
	return ww
}

func (ww *Writer) WriteBool(val bool) *Writer {
	if ww.Err == nil {
		ww.Err = WriteBool(ww.w, val)
	}
	return ww
}

//nolint:govet
func (ww *Writer) WriteByte(val byte) *Writer {
	if ww.Err == nil {
		ww.Err = WriteByte(ww.w, val)
	}
	return ww
}

func (ww *Writer) WriteBytes(val []byte) *Writer {
	if ww.Err == nil {
		ww.Err = WriteBytes(ww.w, val)
	}
	return ww
}

func (ww *Writer) WriteDuration(val time.Duration) *Writer {
	return ww.WriteInt64(int64(val))
}

func (ww *Writer) WriteFromBytes(bytes interface{ Bytes() []byte }) *Writer {
	if ww.Err == nil {
		ww.WriteBytes(bytes.Bytes())
	}
	return ww
}

func (ww *Writer) WriteInt8(val int8) *Writer {
	if ww.Err == nil {
		ww.Err = WriteInt8(ww.w, val)
	}
	return ww
}

func (ww *Writer) WriteInt16(val int16) *Writer {
	if ww.Err == nil {
		ww.Err = WriteInt16(ww.w, val)
	}
	return ww
}

func (ww *Writer) WriteInt32(val int32) *Writer {
	if ww.Err == nil {
		ww.Err = WriteInt32(ww.w, val)
	}
	return ww
}

func (ww *Writer) WriteInt64(val int64) *Writer {
	if ww.Err == nil {
		ww.Err = WriteInt64(ww.w, val)
	}
	return ww
}

func (ww *Writer) WriteMarshaled(m encoding.BinaryMarshaler) *Writer {
	if ww.Err == nil {
		if m == nil {
			ww.Err = errors.New("nil marshaler")
			return ww
		}
		var buf []byte
		buf, ww.Err = m.MarshalBinary()
		ww.WriteBytes(buf)
	}
	return ww
}

type serializable interface {
	Serialize(serializer.DeSerializationMode, interface{}) ([]byte, error)
}

func (ww *Writer) WriteSerialized(s serializable) *Writer {
	if ww.Err == nil {
		if s == nil {
			ww.Err = errors.New("nil deserializer")
			return ww
		}
		var buf []byte
		buf, ww.Err = s.Serialize(serializer.DeSeriModeNoValidation, nil)
		ww.WriteBytes(buf)
	}
	return ww
}

func (ww *Writer) WriteSize(val int) *Writer {
	return ww.WriteSize32(uint32(val))
}

func (ww *Writer) WriteSize32(val uint32) *Writer {
	if ww.Err == nil {
		ww.Err = WriteSize32(ww.w, val)
	}
	return ww
}

func (ww *Writer) WriteString(val string) *Writer {
	if ww.Err == nil {
		ww.Err = WriteString(ww.w, val)
	}
	return ww
}

type marshalUtilWriter interface {
	WriteToMarshalUtil(mu *marshalutil.MarshalUtil)
}

func (ww *Writer) WriteToMarshalUtil(m marshalUtilWriter) *Writer {
	if ww.Err == nil {
		mu := marshalutil.New()
		m.WriteToMarshalUtil(mu)
		ww.WriteN(mu.Bytes()[:mu.WriteOffset()])
	}
	return ww
}

func (ww *Writer) WriteUint8(val uint8) *Writer {
	if ww.Err == nil {
		ww.Err = WriteUint8(ww.w, val)
	}
	return ww
}

func (ww *Writer) WriteUint16(val uint16) *Writer {
	if ww.Err == nil {
		ww.Err = WriteUint16(ww.w, val)
	}
	return ww
}

func (ww *Writer) WriteUint32(val uint32) *Writer {
	if ww.Err == nil {
		ww.Err = WriteUint32(ww.w, val)
	}
	return ww
}

func (ww *Writer) WriteUint64(val uint64) *Writer {
	if ww.Err == nil {
		ww.Err = WriteUint64(ww.w, val)
	}
	return ww
}

// ////////////////// pushback \\\\\\\\\\\\\\\\\\\\

type PushBack struct {
	r   io.Reader
	rr  *Reader
	buf *bytes.Buffer
}

var _ io.ReadWriter = new(PushBack)

func (pb *PushBack) Read(data []byte) (n int, err error) {
	n, err = pb.buf.Read(data)
	if err != nil {
		if errors.Is(err, io.EOF) {
			pb.rr.r = pb.r
			return pb.r.Read(data)
		}
		return n, err
	}
	if n != len(data) {
		return n, errors.New("invalid pushback read")
	}
	return n, nil
}

func (pb *PushBack) Write(data []byte) (n int, err error) {
	if pb.rr.r == pb.r {
		return 0, errors.New("invalid pushback write")
	}
	return pb.buf.Write(data)
}

//////////////////// marshaling \\\\\\\\\\\\\\\\\\\\

func ReadMarshaled(r io.Reader, val encoding.BinaryUnmarshaler) error {
	bin, err := ReadBytes(r)
	if err != nil {
		return err
	}
	return val.UnmarshalBinary(bin)
}

func WriteMarshaled(w io.Writer, val encoding.BinaryMarshaler) error {
	bin, err := val.MarshalBinary()
	if err != nil {
		return err
	}
	return WriteBytes(w, bin)
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

func WriteBytesToMarshalUtil(data []byte, mu *marshalutil.MarshalUtil) {
	size := uint32(len(data))
	mu.WriteBytes(Size32ToBytes(size)).WriteBytes(data)
}

func ReadFromBytes[T any](rr *Reader, fromBytes func([]byte) (T, error)) (ret T) {
	data := rr.ReadBytes()
	if rr.Err == nil {
		ret, rr.Err = fromBytes(data)
	}
	return ret
}

func WriteFromBytes(w io.Writer, bytes interface{ Bytes() []byte }) error {
	return WriteN(w, bytes.Bytes())
}

func FromMarshalUtil[T any](rr *Reader, fromMu func(mu *marshalutil.MarshalUtil) (T, error)) (ret T) {
	if rr.Err == nil {
		buf, ok := rr.r.(*bytes.Buffer)
		if !ok {
			panic("reader expects bytes buffer")
		}
		mu := marshalutil.New(buf.Bytes())
		ret, rr.Err = fromMu(mu)
		rr.r = bytes.NewBuffer(mu.Bytes()[mu.ReadOffset():])
	}
	return ret
}
