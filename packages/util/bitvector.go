// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package util

import (
	"io"

	"github.com/iotaledger/wasp/packages/util/rwutil"
)

type BitVector interface {
	SetBits(positions []int) BitVector
	AsInts() []int
	Bytes() []byte
	Read(r io.Reader) error
	Write(w io.Writer) error
}

type fixBitVector struct {
	size int
	data []byte
}

func NewFixedSizeBitVector(size int) BitVector {
	return &fixBitVector{size: size, data: make([]byte, (size-1)/8+1)}
}

func NewFixedSizeBitVectorFromBytes(data []byte) (BitVector, error) {
	return rwutil.ReaderFromBytes(data, new(fixBitVector))
}

func (b *fixBitVector) Bytes() []byte {
	return rwutil.WriterToBytes(b)
}

func (b *fixBitVector) SetBits(positions []int) BitVector {
	for _, p := range positions {
		bytePos, bitMask := b.bitMask(p)
		b.data[bytePos] |= bitMask
	}
	return b
}

func (b *fixBitVector) AsInts() []int {
	ints := []int{}
	for i := 0; i < b.size; i++ {
		bytePos, bitMask := b.bitMask(i)
		if b.data[bytePos]&bitMask != 0 {
			ints = append(ints, i)
		}
	}
	return ints
}

func (b *fixBitVector) bitMask(position int) (int, byte) {
	return position >> 3, 1 << (position & 0x07)
}

func (b *fixBitVector) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	b.size = rr.ReadSize()
	b.data = make([]byte, (b.size-1)/8+1)
	rr.ReadN(b.data)
	return rr.Err
}

func (b *fixBitVector) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteSize(b.size)
	ww.WriteN(b.data)
	return ww.Err
}
