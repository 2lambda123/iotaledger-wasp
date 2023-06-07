// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package rwutil

import (
	"bytes"
	"errors"
	"io"
)

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
