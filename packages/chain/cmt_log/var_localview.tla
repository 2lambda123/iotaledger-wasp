---- MODULE var_localview ------------------------------------------------------
(*
Notes:
  - If a block is confirmed already, we don't need to use it as a parent
    block for the next block.
  - Alias/Account outputs are moved in a single TX by the chain.
    If moved outside of the chain, we only correlate them in the scope
    of a milestone/time-slot.

Output:
    - Anchor and Alias outputs to use in the next consensus step.
    - Optionally a block to use as a tip.
    - Optionally a list of TXes to re-publish (in new blocks).

*)
EXTENDS Naturals, Sequences
CONSTANTS Anchors, Accounts, Blocks \* Domains.
VARIABLE cnfAnchor  \* Latest known confirmed anchor output.
VARIABLE cnfAccount \* Latest known confirmed account output.
VARIABLE pending    \* A set of pending TXes.
VARIABLE anchorSI   \* Only assigned in the initial state.
const == <<anchorSI>>
vars == <<cnfAnchor, cnfAccount, pending, const>>

(* Defined explicitly to be able to override. *)
StateIndexes == Nat

(* Defined explicitly to be able to override.
   They are needed here only to model blocks "duplicates"
   by the consensus because of the uncertainty from the L1. *)
LogIndexes == Nat

NIL == CHOOSE NIL :
    /\ NIL \notin Anchors
    /\ NIL \notin Accounts
    /\ NIL \notin Blocks

Entries == [
    si          : StateIndexes,
    li          : LogIndexes,
    anchor      : Anchors,
    account     : Accounts,
    block       : Blocks \cup {NIL},
    consAnchor  : Anchors,  \* Consumed Anchor output.
    consAccount : Accounts, \* Consumed Account output.
    rejected    : BOOLEAN
]

TypeOK ==
    /\ cnfAnchor  \in Anchors  \cup {NIL}
    /\ cnfAccount \in Accounts \cup {NIL}
    /\ pending    \in SUBSET Entries

pendingAfterBySI(e) ==
    { p \in pending : p.si > e.si }

(*
depends(e, d) is true, if there is a chain of pending entries
 through which the entry d depends transitively on the entry e.
*)
RECURSIVE depends(_, _)
depends(anc, d) ==
    \/ d.consAnchor = anc
    \/ \E d2 \in pending : d2.anchor = anc /\ depends(anc, d2)

pendingWithRejected(e) ==
    LET upd(p) == IF p = e \/ depends(e.anchor, p)
                  THEN [p EXCEPT !.rejected = FALSE]
                  ELSE p
    IN { upd(p) : p \in pending }

pendingWithExpired(e) ==
    LET upd(p) == IF p = e
                  THEN [p EXCEPT !.block = NIL]
                  ELSE p
    IN { upd(p) : p \in pending }

\* Replace the existing block, if it was expired already.
pendingWithNew(e) ==
    LET woExpired == { p \in pending : ~(p.anchor = e.anchor /\ p.block = NIL) }
    IN woExpired \cup {e}

(*
The general idea -- clear the rejected entries if all of them are either
confirmed (removed from the pending list) or rejected. This will allow
to proceed with building the chain. Chain can be build when the re is no
rejections in the pending chain.

But the situation is more complicated. Some of the entries can be
with blocks expired, and all this can be forked into several chains..
So we refine the above condition to the following:
  - Only cleanup the pending list, if all entries are either confirmed,
    rejected or expired.
  - When cleaning-up the pending chain, we leave only the entries, that
    are expired, non-forked and depends on the last confirmed output.
*)
pendingCleaned(ps) ==
    IF \A p \in ps : p.rejected \/ p.block = NIL
    THEN
        LET notRejected == { p \in ps : ~p.rejected }
            noForks == { p \in notRejected : \A p2 \in notRejected : p2.si = p.si => p2 = p }
            noGaps == { p \in noForks : depends(cnfAnchor, p) }
        IN noGaps
    ELSE ps

--------------------------------------------------------------------------------
\* Actions.

ConsensusOutputDone ==
    \E anc, cAnc \in Anchors,
       acc, cAcc \in Accounts,
       b \in Blocks,
       li \in LogIndexes
       :
        /\ cnfAnchor # NIL
        /\ cnfAccount # NIL
        /\ anchorSI[anc] > anchorSI[cnfAnchor]
        /\ pending' = pendingWithNew([
                si          |-> anchorSI[anc],
                li          |-> li,
                anchor      |-> anc,
                account     |-> acc,
                block       |-> b,
                consAnchor  |-> cAnc,
                consAccount |-> cAcc,
                rejected    |-> FALSE
            ])
        /\ UNCHANGED <<cnfAnchor, cnfAccount, const>>


(*
We can have multiple entries with the received ao.
That's because multiple blocks can publish the same TX with the same AO.
But all of them will have the same SI, but different LIs.

TODO: If we receive AOs confirmed instead of blocks, we don't
      know which block it was, thus cannot use it for pipelining.
      Check, maybe we receive blocks, not outputs.
*)
AnchorOutputConfirmed == \E anc \in Anchors:
    /\ cnfAnchor' = anc
    /\ IF \E e \in pending : e.anchor = anc THEN
        \E e \in pending : e.anchor = anc \* Should be a singe.
            /\ pending'    = pendingCleaned(pendingAfterBySI(e))
            /\ cnfAccount' = e.account
       ELSE
        \* In this case we don't know the account output anymore,
        \* because that's a change from outside.
        /\ pending'    = {}
        /\ cnfAccount' = NIL
    /\ UNCHANGED <<const>>

(*
This action is symmetric to the AnchorOutputConfirmed (mod Account/Anchor).
*)
AccountOutputConfirmed == \E acc \in Accounts:
    /\ cnfAccount' = acc
    /\ IF \E e \in pending : e.account = acc THEN
        \E e \in pending : e.account = acc \* Should be a singe.
            /\ pending'   = pendingCleaned(pendingAfterBySI(e))
            /\ cnfAnchor' = e.anchor
       ELSE
        /\ pending'   = {}
        /\ cnfAnchor' = NIL
    /\ UNCHANGED <<const>>

BothOutputsConfirmed == \E anc \in Anchors, acc \in Accounts:
    /\ cnfAnchor' = anc
    /\ cnfAccount' = acc
    /\ IF \E e \in pending : e.anchor = anc /\ e.account = acc THEN
       (* we can use \/ in the following, but either both will be in the TX by us,
          or they are actually externally produced, possibly in separate TXes. *)
        \E e \in pending : e.anchor = anc /\ e.account = acc
            /\ pending' = pendingCleaned(pendingAfterBySI(e))
       ELSE
        /\ pending' = {}
    /\ UNCHANGED <<const>>


AnchorOutputRejected ==
    \E anc \in Anchors:
        \E e \in pending : e.anchor = anc
            /\ pending' = pendingCleaned(pendingWithRejected(e))
            /\ UNCHANGED <<cnfAnchor, cnfAccount, const>>

AccountOutputRejected ==
    \E acc \in Accounts:
        \E e \in pending : e.account = acc
            /\ pending' = pendingCleaned(pendingWithRejected(e))
            /\ UNCHANGED <<cnfAnchor, cnfAccount, const>>

(*
If a block is outdated, then we mark the corresponding entry as not having
the block assigned. The node will wait until all the blocks are either
confirmed, or all of them are outdated (or rejected).

NOTE: We have considered the following alternatives and went with
the case (A), as it is safer, regarding the limited knowledge on
the L1 node behaviour.
  A) Upon reception of the event on the block expiry we mark only a single
     block as outdated. The chain will be built further when all the blocks
     are invalidated individually (or confirmed, etc).
  B) Invalidate the expired block as well as all the the dependent entries.
     Here we can also start building the chain immediately from the last
     non-expired block. That sounds like an optimization, but can cause
     more rejections if new blocks are build on soon-expiring-blocks.
*)
BlockExpired ==
    \E blk \in Blocks:
        \E e \in pending:
            /\ e.block = blk
            /\ pending' = pendingCleaned(pendingWithExpired(e))
            /\ UNCHANGED <<cnfAnchor, cnfAccount, const>>

--------------------------------------------------------------------------------
Init ==
    /\ cnfAnchor = NIL
    /\ cnfAccount = NIL
    /\ pending = {}
    /\ anchorSI \in [Anchors -> StateIndexes]

Next ==
    \/ ConsensusOutputDone
    \/ BothOutputsConfirmed
    \/ AnchorOutputConfirmed \/ AccountOutputConfirmed
    \/ AnchorOutputRejected  \/ AccountOutputRejected
    \/ BlockExpired

Fair == WF_vars(Next)

Spec == Init /\ [][Next]_vars /\ Fair

--------------------------------------------------------------------------------
\* Properties.

(*
We have an output, if we have a confirmed base and have
no pending forks nor unresolved rejections AND we have either
no block expired, or all of the remaining are expired.
*)
HaveOutput ==
    /\ cnfAnchor # NIL /\ cnfAccount # NIL              \* Have a confirmed base.
    /\ \A e1, e2 \in pending: e1.si = e2.si => e1 = e2  \* Have no pending forks.
    /\ \A e \in pending: ~e.rejected                    \* Have no unresolved rejections.
    /\ \E e \in pending: e.block = NIL => \A ee \in pending: ee.block = NIL \* All or none.
    /\ \A e \in pending: depends(cnfAnchor, e)          \* Have chain without gaps.

(*
The output, if exists, is
  - either the last pending anc/acc/blk, or
  - the last confirmed anc/acc, if there is no pending entries.

  TODO: Rejected vs Expired.
*)
Output ==
    IF HaveOutput
    THEN IF \E e \in pending : e.block # NIL
         THEN
            LET last == CHOOSE e \in pending :
                            /\ e.block # NIL
                            /\ \A e2 \in pending: e2.si <= e.si
            IN
            [
                baseAnc |-> last.anchor,
                baseAcc |-> last.account,
                baseBlk |-> last.block,
                reattach |-> {p \in pending : p.si > last.si}
            ]
         ELSE
            [
                baseAnc |-> cnfAnchor,
                baseAcc |-> cnfAccount,
                baseBlk |-> NIL,
                reattach |-> {pending}
            ]
    ELSE
        [
            baseAnc |-> NIL,
            baseAcc |-> NIL,
            baseBlk |-> NIL,
            reattach |-> {}
        ]


================================================================================
