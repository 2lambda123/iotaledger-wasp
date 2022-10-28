//
//
//
//
//
//

package smUtils

import (
	"github.com/iotaledger/hive.go/core/logger"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/util"
)

type NodeRandomiser struct {
	me          gpa.NodeID
	nodeIDs     []gpa.NodeID
	permutation *util.Permutation16
	log         *logger.Logger
}

func NewNodeRandomiser(me gpa.NodeID, nodeIDs []gpa.NodeID, log *logger.Logger) *NodeRandomiser {
	result := NewNodeRandomiserNoInit(me, log)
	result.init(nodeIDs)
	return result
}

// Before using the returned NodeRandomiser, it must be initted: UpdateNodeIDs
// method must be called.
func NewNodeRandomiserNoInit(me gpa.NodeID, log *logger.Logger) *NodeRandomiser {
	return &NodeRandomiser{
		me:          me,
		nodeIDs:     nil, // Will be set in result.UpdateNodeIDs([]gpa.NodeID).
		permutation: nil, // Will be set in result.UpdateNodeIDs([]gpa.NodeID).
		log:         log.Named("nr"),
	}
}

func (nrT *NodeRandomiser) init(allNodeIDs []gpa.NodeID) {
	nrT.nodeIDs = make([]gpa.NodeID, 0, len(allNodeIDs)-1)
	for _, nodeID := range allNodeIDs {
		if nodeID != nrT.me { // Do not include self to the permutation.
			nrT.nodeIDs = append(nrT.nodeIDs, nodeID)
		}
	}
	var err error
	nrT.permutation, err = util.NewPermutation16(uint16(len(nrT.nodeIDs)))
	if err != nil {
		nrT.log.Warnf("Failed to generate cryptographically secure random domains permutation; will use insecure one: %v", err)
		return
	}
}

func (nrT *NodeRandomiser) UpdateNodeIDs(nodeIDs []gpa.NodeID) {
	nrT.init(nodeIDs)
}

func (nrT *NodeRandomiser) IsInitted() bool {
	return nrT.permutation != nil
}

func (nrT *NodeRandomiser) GetRandomOtherNodeIDs(upToNumPeers int) []gpa.NodeID {
	if upToNumPeers > len(nrT.nodeIDs) {
		upToNumPeers = len(nrT.nodeIDs)
	}
	ret := make([]gpa.NodeID, upToNumPeers)
	for i := 0; i < upToNumPeers; {
		ret[i] = nrT.nodeIDs[nrT.permutation.NextNoCycles()]
		distinct := true
		for j := 0; j < i && distinct; j++ {
			if ret[i].Equals(ret[j]) {
				distinct = false
			}
		}
		if distinct {
			i++
		}
	}
	return ret
}
