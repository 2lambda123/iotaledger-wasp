// 'root' a core contract on the chain. It is responsible for:
// - initial setup of the chain during chain deployment
// - maintaining of core parameters of the chain
// - maintaining (setting, delegating) chain owner ID
// - maintaining (granting, revoking) smart contract deployment rights
// - deployment of smart contracts on the chain and maintenance of contract registry

package rootimpl

import (
	"github.com/samber/lo"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/isc/coreutil"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/collections"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/vm"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/blob"
	"github.com/iotaledger/wasp/packages/vm/core/blocklog"
	"github.com/iotaledger/wasp/packages/vm/core/errors"
	"github.com/iotaledger/wasp/packages/vm/core/errors/coreerrors"
	"github.com/iotaledger/wasp/packages/vm/core/evm"
	"github.com/iotaledger/wasp/packages/vm/core/governance"
	"github.com/iotaledger/wasp/packages/vm/core/root"
)

var Processor = root.Contract.Processor(nil,
	root.FuncDeployContract.WithHandler(deployContract),
	root.FuncGrantDeployPermission.WithHandler(grantDeployPermission),
	root.FuncRequireDeployPermissions.WithHandler(requireDeployPermissions),
	root.FuncRevokeDeployPermission.WithHandler(revokeDeployPermission),
	root.ViewFindContract.WithHandler(findContract),
	root.ViewGetContractRecords.WithHandler(getContractRecords),
)

func SetInitialState(state kv.KVStore) {
	contractRegistry := collections.NewMap(state, root.VarContractRegistry)
	if contractRegistry.Len() != 0 {
		panic("contract registry must be empty on chain start")
	}

	// forbid deployment of custom contracts by default
	state.Set(root.VarDeployPermissionsEnabled, codec.Bool.Encode(true))

	{
		// register core contracts
		contracts := []*coreutil.ContractInfo{
			root.Contract,
			blob.Contract,
			accounts.Contract,
			blocklog.Contract,
			errors.Contract,
			governance.Contract,
			evm.Contract,
		}

		for _, c := range contracts {
			storeContractRecord(
				state,
				root.ContractRecordFromContractInfo(c),
			)
		}
	}
}

var errInvalidContractName = coreerrors.Register("invalid contract name").Create()

// deployContract deploys contract and calls its 'init' constructor.
// If call to the constructor returns an error or another error occurs,
// removes smart contract form the registry as if it was never attempted to deploy
// Inputs:
//   - ParamName string, the unique name of the contract in the chain. Later used as Hname
//   - ParamProgramHash HashValue is a hash of the blob which represents program binary in the 'blob' contract.
//     In case of hardcoded examples it's an arbitrary unique hash set in the global call examples.AddProcessor
func deployContract(ctx isc.Sandbox) dict.Dict {
	ctx.Log().LogDebugf("root.deployContract.begin")
	if !isAuthorizedToDeploy(ctx) {
		panic(vm.ErrUnauthorized)
	}

	params := ctx.Params()
	progHash := params.MustGetHashValue(root.ParamProgramHash)
	name := params.MustGetString(root.ParamName)
	if name == "" || len(name) > 255 {
		panic(errInvalidContractName)
	}

	// pass to init function all params not consumed so far
	initParams := dict.New()
	params.Dict.Iterate("", func(key kv.Key, value []byte) bool {
		if key != root.ParamProgramHash && key != root.ParamName {
			initParams.Set(key, value)
		}
		return true
	})
	// call to load VM from binary to check if it loads successfully
	err := ctx.Privileged().TryLoadContract(progHash)
	ctx.RequireNoError(err, "root.deployContract.fail 1: ")

	// VM loaded successfully. Storing contract in the registry and calling constructor
	storeContractRecord(ctx.State(), &root.ContractRecord{
		ProgramHash: progHash,
		Name:        name,
	})
	ctx.Call(isc.NewMessage(isc.Hn(name), isc.EntryPointInit, initParams), nil)
	eventDeploy(ctx, progHash, name)
	return nil
}

// grantDeployPermission grants permission to deploy contracts
func grantDeployPermission(ctx isc.Sandbox, deployer isc.AgentID) dict.Dict {
	ctx.RequireCallerIsChainOwner()
	collections.NewMap(ctx.State(), root.VarDeployPermissions).SetAt(deployer.Bytes(), []byte{0x01})
	eventGrant(ctx, deployer)
	return nil
}

// revokeDeployPermission revokes permission to deploy contracts
func revokeDeployPermission(ctx isc.Sandbox, deployer isc.AgentID) dict.Dict {
	ctx.RequireCallerIsChainOwner()
	collections.NewMap(ctx.State(), root.VarDeployPermissions).DelAt(deployer.Bytes())
	eventRevoke(ctx, deployer)
	return nil
}

func requireDeployPermissions(ctx isc.Sandbox, permissionsEnabled bool) dict.Dict {
	ctx.RequireCallerIsChainOwner()
	ctx.State().Set(root.VarDeployPermissionsEnabled, codec.Bool.Encode(permissionsEnabled))
	return nil
}

// findContract view finds and returns encoded record of the contract
func findContract(ctx isc.SandboxView, hname isc.Hname) (bool, *root.ContractRecord) {
	rec := root.FindContract(ctx.StateR(), hname)
	return rec != nil, rec
}

func getContractRecords(ctx isc.SandboxView) map[isc.Hname]*root.ContractRecord {
	ret := make(map[isc.Hname]*root.ContractRecord)
	root.GetContractRegistryR(ctx.StateR()).Iterate(func(elemKey []byte, value []byte) bool {
		ret[lo.Must(codec.Hname.Decode(elemKey))] = lo.Must(root.ContractRecordFromBytes(value))
		return true
	})
	return ret
}
