package tstemplates

var paramsTs = map[string]string{
	// *******************************
	"params.ts": `
$#emit importWasmLib
$#emit importSc
$#each func paramsFunc
`,
	// *******************************
	"paramsFunc": `
$#if params paramsFuncParams
`,
	// *******************************
	"paramsFuncParams": `
$#set Kind Param
$#set mut Immutable
$#if param paramsProxyStruct
$#set mut Mutable
$#if param paramsProxyStruct
`,
	// *******************************
	"paramsProxyStruct": `
$#set TypeName $mut$FuncName$+Params
$#each param proxyContainers

export class $TypeName extends wasmlib.ScMapID {
$#each param proxyMethods
}
`,
}
