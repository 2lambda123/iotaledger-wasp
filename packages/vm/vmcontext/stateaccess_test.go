package vmcontext

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/hive.go/kvstore/mapdb"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/buffered"
	"github.com/iotaledger/wasp/packages/origin"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/vm"
)

func TestSetThenGet(t *testing.T) {
	db := mapdb.NewMapDB()
	cs := state.NewStore(db)
	origin.InitChain(cs, nil, 0)
	latest, err := cs.LatestBlock()
	require.NoError(t, err)
	stateDraft, err := cs.NewStateDraft(time.Now(), latest.L1Commitment())
	require.NoError(t, err)

	stateUpdate := buffered.NewMutations()
	hname := isc.Hn("test")

	task := &vm.VMTask{}
	taskResult := task.CreateResult()
	taskResult.StateDraft = stateDraft
	vmctx := &VMContext{
		task:               task,
		taskResult:         taskResult,
		currentStateUpdate: stateUpdate,
		callStack:          []*callContext{{contract: hname}},
	}
	s := vmctx.State()

	subpartitionedKey := kv.Key(hname.Bytes()) + "x"

	// contract sets variable x
	s.Set("x", []byte{42})
	require.Equal(t, map[kv.Key][]byte{subpartitionedKey: {42}}, stateUpdate.Sets)
	require.Equal(t, map[kv.Key]struct{}{}, stateUpdate.Dels)

	// contract gets variable x
	v := s.Get("x")
	require.Equal(t, []byte{42}, v)

	// mutation is in currentStateUpdate, prefixed by the contract id
	require.Equal(t, []byte{42}, stateUpdate.Sets[subpartitionedKey])

	// mutation is in the not committed to the virtual state yet
	v = stateDraft.Get(subpartitionedKey)
	require.Nil(t, v)

	// contract deletes variable x
	s.Del("x")
	require.Equal(t, map[kv.Key][]byte{}, stateUpdate.Sets)
	require.Equal(t, map[kv.Key]struct{}{subpartitionedKey: {}}, stateUpdate.Dels)

	// contract sees variable x does not exist
	v = s.Get("x")
	require.Nil(t, v)

	// contract makes several writes to same variable, gets the latest value
	s.Set("x", []byte{2 * 42})
	require.Equal(t, map[kv.Key][]byte{subpartitionedKey: {2 * 42}}, stateUpdate.Sets)
	require.Equal(t, map[kv.Key]struct{}{}, stateUpdate.Dels)

	s.Set("x", []byte{3 * 42})
	require.Equal(t, map[kv.Key][]byte{subpartitionedKey: {3 * 42}}, stateUpdate.Sets)
	require.Equal(t, map[kv.Key]struct{}{}, stateUpdate.Dels)

	v = s.Get("x")
	require.Equal(t, []byte{3 * 42}, v)
}

func TestIterate(t *testing.T) {
	db := mapdb.NewMapDB()
	cs := state.NewStore(db)
	origin.InitChain(cs, nil, 0)
	latest, err := cs.LatestBlock()
	require.NoError(t, err)
	stateDraft, err := cs.NewStateDraft(time.Now(), latest.L1Commitment())
	require.NoError(t, err)

	stateUpdate := buffered.NewMutations()
	hname := isc.Hn("test")

	task := &vm.VMTask{}
	taskResult := task.CreateResult()
	taskResult.StateDraft = stateDraft
	vmctx := &VMContext{
		task:               task,
		taskResult:         taskResult,
		currentStateUpdate: stateUpdate,
		callStack:          []*callContext{{contract: hname}},
	}
	s := vmctx.State()
	s.Set("xy1", []byte{42})
	s.Set("xy2", []byte{42 * 2})

	arr := make([][]byte, 0)
	s.IterateSorted("xy", func(k kv.Key, v []byte) bool {
		require.True(t, strings.HasPrefix(string(k), "xy"))
		arr = append(arr, v)
		return true
	})
	require.EqualValues(t, 2, len(arr))
	require.Equal(t, []byte{42}, arr[0])
	require.Equal(t, []byte{42 * 2}, arr[1])
}

func TestVmctxStateDeletion(t *testing.T) {
	db := mapdb.NewMapDB()
	cs := state.NewStore(db)
	origin.InitChain(cs, nil, 0)

	foo := kv.Key("foo")
	{
		latest, err := cs.LatestBlock()
		require.NoError(t, err)
		stateDraft, err := cs.NewStateDraft(time.Now(), latest.L1Commitment())
		require.NoError(t, err)
		stateDraft.Set(foo, []byte("bar"))
		block := cs.Commit(stateDraft)
		err = cs.SetLatest(block.TrieRoot())
		require.NoError(t, err)
	}

	latest, err := cs.LatestBlock()
	require.NoError(t, err)
	stateDraft, err := cs.NewStateDraft(time.Now(), latest.L1Commitment())
	require.NoError(t, err)
	stateUpdate := buffered.NewMutations()
	task := &vm.VMTask{}
	taskResult := task.CreateResult()
	taskResult.StateDraft = stateDraft
	vmctx := &VMContext{
		task:               task,
		taskResult:         taskResult,
		currentStateUpdate: stateUpdate,
	}
	vmctxStore := vmctx.chainState()
	require.EqualValues(t, "bar", vmctxStore.Get(foo))
	vmctxStore.Del(foo)
	require.False(t, vmctxStore.Has(foo))
	val := vmctxStore.Get(foo)
	require.Nil(t, val)
}
