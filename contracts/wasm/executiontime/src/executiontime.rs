// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

use wasmlib::*;

use crate::*;

pub fn func_f(_ctx: &ScFuncContext, f: &FContext) {
    let n = f.params.n().value();
    let mut x: u32 = 0;
    let mut _y: u32 = 0;

    for _ in 0..n {
        x += 1;
        _y += 3 * (x % 10);
    }
}
