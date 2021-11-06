package tstemplates

var typedefsTs = map[string]string{
	// *******************************
	"typedefs.ts": `
$#emit importWasmLib
$#emit importSc
$#each typedef typedefProxy
`,
	// *******************************
	"typedefProxy": `
$#set mut Immutable
$#if array typedefProxyArray
$#if array typedefProxyAlias
$#if map typedefProxyMap
$#if map typedefProxyAlias
$#set mut Mutable
$#if array typedefProxyArray
$#if array typedefProxyAlias
$#if map typedefProxyMap
$#if map typedefProxyAlias
`,
	// *******************************
	"typedefProxyAlias": `

export class $mut$FldName extends $proxy {
};
`,
	// *******************************
	"typedefProxyArray": `
$#set proxy ArrayOf$mut$FldType
$#if exist else typedefProxyArrayNew
`,
	// *******************************
	"typedefProxyArrayNew": `

export class $proxy {
	objID: i32;

    constructor(objID: i32) {
        this.objID = objID;
    }
$#if mut typedefProxyArrayClear

    length(): i32 {
        return wasmlib.getLength(this.objID);
    }
$#if basetype typedefProxyArrayNewBaseType typedefProxyArrayNewOtherType
}
$#set exist $proxy
`,
	// *******************************
	"typedefProxyArrayClear": `

    clear(): void {
        wasmlib.clear(this.objID);
    }
`,
	// *******************************
	"typedefProxyArrayNewBaseType": `

    get$FldType(index: i32): wasmlib.Sc$mut$FldType {
        return new wasmlib.Sc$mut$FldType(this.objID, new wasmlib.Key32(index));
    }
`,
	// *******************************
	"typedefProxyArrayNewOtherType": `
$#set OldType $FldType
$#if typedef typedefProxyArrayNewOtherTypeTypeDef typedefProxyArrayNewOtherTypeStruct
`,
	// *******************************
	"typedefProxyArrayNewOtherTypeTypeDef": `
$#set varType wasmlib.TYPE_MAP
$#if array setVarTypeArray

	Get$OldType(index: i32): sc.$mut$OldType {
		let subID = wasmlib.getObjectID(this.objID, new wasmlib.Key32(index), $varType);
		return new sc.$mut$OldType(subID);
	}
`,
	// *******************************
	"typedefProxyArrayNewOtherTypeStruct": `

	get$FldType(index: i32): sc.$mut$FldType {
		return new sc.$mut$FldType(this.objID, new wasmlib.Key32(index));
	}
`,
	// *******************************
	"typedefProxyMap": `
$#set proxy Map$FldMapKey$+To$mut$FldType
$#if exist else typedefProxyMapNew
`,
	// *******************************
	"typedefProxyMapNew": `

export class $proxy {
	objID: i32;

    constructor(objID: i32) {
        this.objID = objID;
    }
$#if mut typedefProxyMapClear
$#if basetype typedefProxyMapNewBaseType typedefProxyMapNewOtherType
}
$#set exist $proxy
`,
	// *******************************
	"typedefProxyMapClear": `

    clear(): void {
        wasmlib.clear(this.objID);
    }
`,
	// *******************************
	"typedefProxyMapNewBaseType": `

    get$FldType(key: $FldMapKeyLangType): wasmlib.Sc$mut$FldType {
        return new wasmlib.Sc$mut$FldType(this.objID, $FldMapKeyKey.getKeyID());
    }
`,
	// *******************************
	"typedefProxyMapNewOtherType": `
$#set oldType $fldType
$#set OldType $FldType
$#set OldMapKeyLangType $FldMapKeyLangType
$#set OldMapKeyKey $FldMapKeyKey
$#if typedef typedefProxyMapNewOtherTypeTypeDef typedefProxyMapNewOtherTypeStruct
`,
	// *******************************
	"typedefProxyMapNewOtherTypeTypeDef": `
$#set varType wasmlib.TYPE_MAP
$#if array setVarTypeArray

    get$OldType(key: $OldMapKeyLangType): sc.$mut$OldType {
        let subID = wasmlib.getObjectID(this.objID, $OldMapKeyKey.getKeyID(), $varType);
        return new sc.$mut$OldType(subID);
    }
`,
	// *******************************
	"typedefProxyMapNewOtherTypeStruct": `

    get$OldType(key: $OldMapKeyLangType): sc.$mut$OldType {
        return new sc.$mut$OldType(this.objID, $OldMapKeyKey.getKeyID());
    }
`,
	// *******************************
	"setVarTypeArray": `
$#set varType $arrayTypeID|$FldTypeID
`,
}
