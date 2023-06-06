// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package testpeers

import (
	"bytes"
	"embed"
	"fmt"

	"github.com/iotaledger/wasp/packages/util"
)

//go:embed testkeys_pregenerated-*.bin
var embedded embed.FS

func pregeneratedDksName(n, t uint16) string {
	return fmt.Sprintf("testkeys_pregenerated-%v-%v.bin", n, t)
}

func pregeneratedDksRead(n, t uint16) [][]byte {
	buf, err := embedded.ReadFile(pregeneratedDksName(n, t))
	if err != nil {
		panic(err)
	}
	r := bytes.NewReader(buf)
	rr := util.NewReader(r)

	bufN := rr.ReadUint16()
	if n != bufN {
		panic("wrong_file")
	}
	res := make([][]byte, n)
	for i := range res {
		res[i] = rr.ReadBytes()
	}
	return res
}
