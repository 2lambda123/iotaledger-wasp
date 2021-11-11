// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

use wasmlib::*;

use crate::*;
use crate::structs::*;

pub fn func_add_test_struct(ctx: &ScFuncContext, f: &AddTestStructContext) {
    let test_structs = f.state.test_structs();
    let length = test_structs.length();
    let description: String = f.params.description().value();
    let test_struct = TestStruct {
        id: length,
        description: description.clone()
    };
    test_structs.get_test_struct(length).set_value(&test_struct);

    ctx.event(&format!(
        "teststruct.added {0} {1}",
        length.to_string(),
        description
    ))
}

pub fn func_clear_all(ctx: &ScFuncContext, f: &ClearAllContext) {
    let mut test_struct: Option<MutableTestStruct> = None;
    test_struct = Some(f.state.test_structs().get_test_struct(0));


    f.state.test_structs().clear();
    ctx.log(&format!(
        "teststructs.cleared {0}",
        f.state.test_structs().length().to_string(),
    ));

    for i in 1..10 {
        let description = f.state.test_structs().get_test_struct(i).value().description;
        let id = f.state.test_structs().get_test_struct(i).value().id;
        ctx.log(&format!(
            "teststructs.cleared.log {0} \n {1} {2} {3}",
            i.to_string(),
            f.state.test_structs().length().to_string(),
            description,
            id.to_string()
        ));
    }

    test_struct.unwrap();
    ctx.log(&format!(
        "teststruct.unwrapped {0} {1}",
        f.state.test_structs().length().to_string(),
        f.state.test_structs().get_test_struct(1).value().description
    ));
}
