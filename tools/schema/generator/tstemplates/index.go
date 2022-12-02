// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package tstemplates

var indexTs = map[string]string{
	// *******************************
	"index_impl.ts": `
export * from "./$package";
export * from "./lib";
`,
	// *******************************
	"index.ts": `
export * from "./consts";
export * from "./contract";
$#if events exportEvents
$#if params exportParams
$#if results exportResults
$#if state exportState
$#if structs exportStructs
$#if typedefs exportTypedefs
`,
	// *******************************
	"exportEvents": `
export * from "./events";
`,
	// *******************************
	"exportParams": `
export * from "./params";
`,
	// *******************************
	"exportResults": `
export * from "./results";
`,
	// *******************************
	"exportState": `
export * from "./state";
`,
	// *******************************
	"exportStructs": `
export * from "./structs";
`,
	// *******************************
	"exportTypedefs": `
export * from "./typedefs";
`,
}
