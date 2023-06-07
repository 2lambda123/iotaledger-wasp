package util

import (
	"encoding/binary"
	"errors"
)

//////////////////// int8 \\\\\\\\\\\\\\\\\\\\

func Int8ToBytes(val int8) []byte {
	return Uint8ToBytes(uint8(val))
}

func Int8FromBytes(b []byte) (int8, error) {
	ret, err := Uint8FromBytes(b)
	return int8(ret), err
}

func MustInt8FromBytes(b []byte) int8 {
	return int8(MustUint8FromBytes(b))
}

//////////////////// int16 \\\\\\\\\\\\\\\\\\\\

func Int16ToBytes(val int16) []byte {
	return Uint16ToBytes(uint16(val))
}

func Int16FromBytes(b []byte) (int16, error) {
	ret, err := Uint16FromBytes(b)
	return int16(ret), err
}

func MustInt16FromBytes(b []byte) int16 {
	return int16(MustUint16FromBytes(b))
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
