// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	"github.com/iotaledger/hive.go/ierrors"
	"github.com/iotaledger/wasp/tools/schema/model"
)

func Unmarshal(in []byte, def *model.SchemaDef) error {
	root := Parse(in)
	if root == nil {
		return ierrors.New("failed to parse input yaml file")
	}
	return Convert(root, def)
}
