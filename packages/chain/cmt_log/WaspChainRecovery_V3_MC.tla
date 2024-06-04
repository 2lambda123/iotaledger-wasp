---- MODULE WaspChainRecovery_V3_MC ----
EXTENDS WaspChainRecovery_V3

MC_ActionConstraint ==
    \A n \in CN :
        /\ qcConsOut[n].sent < MaxLI
        /\ qcL1RepAO[n].sent < MaxLI
        /\ qcRecover[n].sent < MaxLI
        /\ qcStarted[n].sent < MaxLI

====
