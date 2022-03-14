# Schema Tool

## Commands

## types

The followings are concrete types. Remember **map** and **array** are not concrete types.

* base types
    * Bool
    * Bytes
    * Int8
    * Int16
    * Int32
    * Int64
    * String
    * Uint8
    * Uint16
    * Uint32
    * Uint64
* smart contract specific types
    * Address
    * AgentID
    * ChainID
    * Color
    * Hash
    * Hname
    * RequestID

### map

Map key should always be concrete type

### array


## typedef

typedefs are originally used to be able to convey multi-dimensional containers.

## struct

A struct is stored as a serialized byte array. So all struct fields need to be serializable.

Not able to contain field with non-concrete datatype (e.g. array, map). At the same time, typedef and struct are not allowed to be treated as field type, too.

## state

## func/view

the function parameters and return results in different datatype can't share the same name. The following is an example.

### parameter naming

#### correct

```yaml
mapOfArraysLength:
    params:
      name: String
    results:
      length: Uint32

arrayOfArraysLength:
    results:
      length: Uint32
```

#### wrong

```yaml
mapOfArraysLength:
    params:
      name: String
    results:
      length: Uint32

arrayOfArraysLength:
    results:
      length: Int8
```

Here, in the wrong case, the same returned results `length`, are in different datatype. One is `Uint32`, and the other one is `Int8`. Thereupon, an error message `redefined result type: length` will be returned.

### function access

Developers can limit the access of caller with following keyword.

* self: Only the smart contract itself can call this function.
* chain: Only the chain owner can call this function.
* creator: Only the contract creator can call this function

### optional function parameter


### func (mutable function)

### view (immutable function)

## contract generation stages

1. yaml/json file parsing
1. SchemaDef object


## template syntax

Any line starting with a special `$#` directive will recursively be processed

### emit

### each

`$#each <key> <template>`

key includes

* event
* events
* func
* mandatory
* param
* params
* result
* results
* state
* struct
* structs
* typedef

See `GenBase.emitEach` for more information

### func

### if

`$#if <condition> [else] <template>`

`<condition>` are predefined conditions, you can go to `GenBase.emitIf` for seeing the list of supporting conditions.

if keyword `else` is provided, then the expression becomes the negate expression of the expression without `else`.

### set

`"$#set <key> <value>"`

`set` directive specifies `<key>` to `<value>`, which can be any string

A special key `exist` is used to add a newly generated type, It can be used to prevent duplicate types from being generated. An example of `exist` directive is `$#set exist proxyExample`.
For `if` directive, if it comes with `exist` key (in this format `$#if exist <template>`), then it schema tool will check whether the type `proxyExample` has existed of not.




TBD
