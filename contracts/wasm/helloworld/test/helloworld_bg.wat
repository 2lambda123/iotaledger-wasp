(module
  (type (;0;) (func (param i32 i32)))
  (type (;1;) (func (param i32 i32) (result i32)))
  (type (;2;) (func (param i32)))
  (type (;3;) (func (param i32) (result i32)))
  (type (;4;) (func (param i32 i32 i32) (result i32)))
  (type (;5;) (func (param i32 i32 i32)))
  (type (;6;) (func (param i32) (result i64)))
  (type (;7;) (func))
  (type (;8;) (func (param i32 i32 i32 i32)))
  (type (;9;) (func (param i32 i32 i32 i32) (result i32)))
  (type (;10;) (func (param i32 i32 i32 i32 i32)))
  (type (;11;) (func (result i32)))
  (type (;12;) (func (param i64 i32) (result i32)))
  (import "WasmLib" "hostSetBytes" (func (;0;) (type 10)))
  (import "WasmLib" "hostGetObjectID" (func (;1;) (type 4)))
  (import "WasmLib" "hostGetKeyID" (func (;2;) (type 1)))
  (func (;3;) (type 3) (param i32) (result i32)
    (local i32 i32 i32 i32 i32 i32 i32 i32 i32 i32 i32 i32 i32 i32 i32 i32 i64)
    global.get 0
    i32.const 16
    i32.sub
    local.tee 11
    global.set 0
    block  ;; label = @1
      block  ;; label = @2
        local.get 0
        i32.const 245
        i32.ge_u
        if  ;; label = @3
          i32.const 0
          call 74
          local.tee 1
          local.get 1
          i32.const 8
          call 50
          i32.sub
          i32.const 20
          i32.const 8
          call 50
          i32.sub
          i32.const 16
          i32.const 8
          call 50
          i32.sub
          i32.const -65544
          i32.add
          i32.const -9
          i32.and
          i32.const -3
          i32.add
          local.tee 2
          i32.const 0
          i32.const 16
          i32.const 8
          call 50
          i32.const 2
          i32.shl
          i32.sub
          local.tee 1
          local.get 2
          local.get 1
          i32.lt_u
          select
          local.get 0
          i32.le_u
          br_if 2 (;@1;)
          local.get 0
          i32.const 4
          i32.add
          i32.const 8
          call 50
          local.set 4
          i32.const 1049516
          i32.load
          i32.eqz
          br_if 1 (;@2;)
          i32.const 0
          local.get 4
          i32.sub
          local.set 3
          block  ;; label = @4
            block  ;; label = @5
              block (result i32)  ;; label = @6
                i32.const 0
                local.get 4
                i32.const 8
                i32.shr_u
                local.tee 0
                i32.eqz
                br_if 0 (;@6;)
                drop
                i32.const 31
                local.get 4
                i32.const 16777215
                i32.gt_u
                br_if 0 (;@6;)
                drop
                local.get 4
                i32.const 6
                local.get 0
                i32.clz
                local.tee 0
                i32.sub
                i32.const 31
                i32.and
                i32.shr_u
                i32.const 1
                i32.and
                local.get 0
                i32.const 1
                i32.shl
                i32.sub
                i32.const 62
                i32.add
              end
              local.tee 5
              i32.const 2
              i32.shl
              i32.const 1049784
              i32.add
              i32.load
              local.tee 0
              if  ;; label = @6
                local.get 4
                local.get 5
                call 47
                i32.const 31
                i32.and
                i32.shl
                local.set 7
                i32.const 0
                local.set 1
                loop  ;; label = @7
                  block  ;; label = @8
                    local.get 0
                    call 67
                    local.tee 2
                    local.get 4
                    i32.lt_u
                    br_if 0 (;@8;)
                    local.get 2
                    local.get 4
                    i32.sub
                    local.tee 2
                    local.get 3
                    i32.ge_u
                    br_if 0 (;@8;)
                    local.get 0
                    local.set 1
                    local.get 2
                    local.tee 3
                    br_if 0 (;@8;)
                    i32.const 0
                    local.set 3
                    br 3 (;@5;)
                  end
                  local.get 0
                  i32.const 20
                  i32.add
                  i32.load
                  local.tee 2
                  local.get 6
                  local.get 2
                  local.get 0
                  local.get 7
                  i32.const 29
                  i32.shr_u
                  i32.const 4
                  i32.and
                  i32.add
                  i32.const 16
                  i32.add
                  i32.load
                  local.tee 0
                  i32.ne
                  select
                  local.get 6
                  local.get 2
                  select
                  local.set 6
                  local.get 7
                  i32.const 1
                  i32.shl
                  local.set 7
                  local.get 0
                  br_if 0 (;@7;)
                end
                local.get 6
                if  ;; label = @7
                  local.get 6
                  local.set 0
                  br 2 (;@5;)
                end
                local.get 1
                br_if 2 (;@4;)
              end
              i32.const 0
              local.set 1
              i32.const 1
              local.get 5
              i32.const 31
              i32.and
              i32.shl
              call 53
              i32.const 1049516
              i32.load
              i32.and
              local.tee 0
              i32.eqz
              br_if 3 (;@2;)
              local.get 0
              call 60
              i32.ctz
              i32.const 2
              i32.shl
              i32.const 1049784
              i32.add
              i32.load
              local.tee 0
              i32.eqz
              br_if 3 (;@2;)
            end
            loop  ;; label = @5
              local.get 0
              local.get 1
              local.get 0
              call 67
              local.tee 1
              local.get 4
              i32.ge_u
              local.get 1
              local.get 4
              i32.sub
              local.tee 6
              local.get 3
              i32.lt_u
              i32.and
              local.tee 2
              select
              local.set 1
              local.get 6
              local.get 3
              local.get 2
              select
              local.set 3
              local.get 0
              call 46
              local.tee 0
              br_if 0 (;@5;)
            end
            local.get 1
            i32.eqz
            br_if 2 (;@2;)
          end
          i32.const 1049912
          i32.load
          local.tee 0
          local.get 4
          i32.ge_u
          i32.const 0
          local.get 3
          local.get 0
          local.get 4
          i32.sub
          i32.ge_u
          select
          br_if 1 (;@2;)
          local.get 1
          local.tee 0
          local.get 4
          call 72
          local.set 5
          local.get 0
          call 14
          block  ;; label = @4
            local.get 3
            i32.const 16
            i32.const 8
            call 50
            i32.ge_u
            if  ;; label = @5
              local.get 0
              local.get 4
              call 62
              local.get 5
              local.get 3
              call 48
              local.get 3
              i32.const 256
              i32.ge_u
              if  ;; label = @6
                local.get 5
                local.get 3
                call 13
                br 2 (;@4;)
              end
              local.get 3
              i32.const 3
              i32.shr_u
              local.tee 1
              i32.const 3
              i32.shl
              i32.const 1049520
              i32.add
              local.set 6
              block (result i32)  ;; label = @6
                i32.const 1049512
                i32.load
                local.tee 2
                i32.const 1
                local.get 1
                i32.shl
                local.tee 1
                i32.and
                if  ;; label = @7
                  local.get 6
                  i32.load offset=8
                  br 1 (;@6;)
                end
                i32.const 1049512
                local.get 1
                local.get 2
                i32.or
                i32.store
                local.get 6
              end
              local.set 1
              local.get 6
              local.get 5
              i32.store offset=8
              local.get 1
              local.get 5
              i32.store offset=12
              local.get 5
              local.get 6
              i32.store offset=12
              local.get 5
              local.get 1
              i32.store offset=8
              br 1 (;@4;)
            end
            local.get 0
            local.get 3
            local.get 4
            i32.add
            call 45
          end
          local.get 0
          call 74
          local.tee 3
          br_if 2 (;@1;)
          br 1 (;@2;)
        end
        i32.const 16
        local.get 0
        i32.const 4
        i32.add
        i32.const 16
        i32.const 8
        call 50
        i32.const -5
        i32.add
        local.get 0
        i32.gt_u
        select
        i32.const 8
        call 50
        local.set 4
        block  ;; label = @3
          block  ;; label = @4
            block  ;; label = @5
              block (result i32)  ;; label = @6
                block  ;; label = @7
                  block  ;; label = @8
                    i32.const 1049512
                    i32.load
                    local.tee 1
                    local.get 4
                    i32.const 3
                    i32.shr_u
                    local.tee 0
                    i32.const 31
                    i32.and
                    local.tee 2
                    i32.shr_u
                    local.tee 6
                    i32.const 3
                    i32.and
                    i32.eqz
                    if  ;; label = @9
                      local.get 4
                      i32.const 1049912
                      i32.load
                      i32.le_u
                      br_if 7 (;@2;)
                      local.get 6
                      br_if 1 (;@8;)
                      i32.const 1049516
                      i32.load
                      local.tee 0
                      i32.eqz
                      br_if 7 (;@2;)
                      local.get 0
                      call 60
                      i32.ctz
                      i32.const 2
                      i32.shl
                      i32.const 1049784
                      i32.add
                      i32.load
                      local.tee 1
                      call 67
                      local.get 4
                      i32.sub
                      local.set 3
                      local.get 1
                      call 46
                      local.tee 0
                      if  ;; label = @10
                        loop  ;; label = @11
                          local.get 0
                          call 67
                          local.get 4
                          i32.sub
                          local.tee 2
                          local.get 3
                          local.get 2
                          local.get 3
                          i32.lt_u
                          local.tee 2
                          select
                          local.set 3
                          local.get 0
                          local.get 1
                          local.get 2
                          select
                          local.set 1
                          local.get 0
                          call 46
                          local.tee 0
                          br_if 0 (;@11;)
                        end
                      end
                      local.get 1
                      local.tee 2
                      local.get 4
                      call 72
                      local.set 0
                      local.get 1
                      call 14
                      local.get 3
                      i32.const 16
                      i32.const 8
                      call 50
                      i32.lt_u
                      br_if 5 (;@4;)
                      local.get 2
                      local.get 4
                      call 62
                      local.get 0
                      local.get 3
                      call 48
                      i32.const 1049912
                      i32.load
                      local.tee 1
                      i32.eqz
                      br_if 4 (;@5;)
                      local.get 1
                      i32.const 3
                      i32.shr_u
                      local.tee 1
                      i32.const 3
                      i32.shl
                      i32.const 1049520
                      i32.add
                      local.set 7
                      i32.const 1049920
                      i32.load
                      local.set 5
                      i32.const 1049512
                      i32.load
                      local.tee 6
                      i32.const 1
                      local.get 1
                      i32.const 31
                      i32.and
                      i32.shl
                      local.tee 1
                      i32.and
                      i32.eqz
                      br_if 2 (;@7;)
                      local.get 7
                      i32.load offset=8
                      br 3 (;@6;)
                    end
                    block  ;; label = @9
                      local.get 6
                      i32.const -1
                      i32.xor
                      i32.const 1
                      i32.and
                      local.get 0
                      i32.add
                      local.tee 3
                      i32.const 3
                      i32.shl
                      local.tee 0
                      i32.const 1049528
                      i32.add
                      i32.load
                      local.tee 6
                      i32.const 8
                      i32.add
                      i32.load
                      local.tee 2
                      local.get 0
                      i32.const 1049520
                      i32.add
                      local.tee 0
                      i32.ne
                      if  ;; label = @10
                        local.get 2
                        local.get 0
                        i32.store offset=12
                        local.get 0
                        local.get 2
                        i32.store offset=8
                        br 1 (;@9;)
                      end
                      i32.const 1049512
                      local.get 1
                      i32.const -2
                      local.get 3
                      i32.rotl
                      i32.and
                      i32.store
                    end
                    local.get 6
                    local.get 3
                    i32.const 3
                    i32.shl
                    call 45
                    local.get 6
                    call 74
                    local.set 3
                    br 7 (;@1;)
                  end
                  block  ;; label = @8
                    i32.const 1
                    local.get 2
                    i32.shl
                    call 53
                    local.get 6
                    local.get 2
                    i32.shl
                    i32.and
                    call 60
                    i32.ctz
                    local.tee 2
                    i32.const 3
                    i32.shl
                    local.tee 0
                    i32.const 1049528
                    i32.add
                    i32.load
                    local.tee 3
                    i32.const 8
                    i32.add
                    i32.load
                    local.tee 1
                    local.get 0
                    i32.const 1049520
                    i32.add
                    local.tee 0
                    i32.ne
                    if  ;; label = @9
                      local.get 1
                      local.get 0
                      i32.store offset=12
                      local.get 0
                      local.get 1
                      i32.store offset=8
                      br 1 (;@8;)
                    end
                    i32.const 1049512
                    i32.const 1049512
                    i32.load
                    i32.const -2
                    local.get 2
                    i32.rotl
                    i32.and
                    i32.store
                  end
                  local.get 3
                  local.get 4
                  call 62
                  local.get 3
                  local.get 4
                  call 72
                  local.tee 6
                  local.get 2
                  i32.const 3
                  i32.shl
                  local.get 4
                  i32.sub
                  local.tee 2
                  call 48
                  i32.const 1049912
                  i32.load
                  local.tee 0
                  if  ;; label = @8
                    local.get 0
                    i32.const 3
                    i32.shr_u
                    local.tee 0
                    i32.const 3
                    i32.shl
                    i32.const 1049520
                    i32.add
                    local.set 7
                    i32.const 1049920
                    i32.load
                    local.set 5
                    block (result i32)  ;; label = @9
                      i32.const 1049512
                      i32.load
                      local.tee 1
                      i32.const 1
                      local.get 0
                      i32.const 31
                      i32.and
                      i32.shl
                      local.tee 0
                      i32.and
                      if  ;; label = @10
                        local.get 7
                        i32.load offset=8
                        br 1 (;@9;)
                      end
                      i32.const 1049512
                      local.get 0
                      local.get 1
                      i32.or
                      i32.store
                      local.get 7
                    end
                    local.set 0
                    local.get 7
                    local.get 5
                    i32.store offset=8
                    local.get 0
                    local.get 5
                    i32.store offset=12
                    local.get 5
                    local.get 7
                    i32.store offset=12
                    local.get 5
                    local.get 0
                    i32.store offset=8
                  end
                  i32.const 1049920
                  local.get 6
                  i32.store
                  i32.const 1049912
                  local.get 2
                  i32.store
                  local.get 3
                  call 74
                  local.set 3
                  br 6 (;@1;)
                end
                i32.const 1049512
                local.get 1
                local.get 6
                i32.or
                i32.store
                local.get 7
              end
              local.set 1
              local.get 7
              local.get 5
              i32.store offset=8
              local.get 1
              local.get 5
              i32.store offset=12
              local.get 5
              local.get 7
              i32.store offset=12
              local.get 5
              local.get 1
              i32.store offset=8
            end
            i32.const 1049920
            local.get 0
            i32.store
            i32.const 1049912
            local.get 3
            i32.store
            br 1 (;@3;)
          end
          local.get 2
          local.get 3
          local.get 4
          i32.add
          call 45
        end
        local.get 2
        call 74
        local.tee 3
        br_if 1 (;@1;)
      end
      block  ;; label = @2
        block  ;; label = @3
          block  ;; label = @4
            block  ;; label = @5
              block  ;; label = @6
                i32.const 1049912
                i32.load
                local.tee 0
                local.get 4
                i32.lt_u
                if  ;; label = @7
                  i32.const 1049916
                  i32.load
                  local.tee 0
                  local.get 4
                  i32.gt_u
                  br_if 3 (;@4;)
                  i32.const 0
                  local.set 3
                  local.get 11
                  local.get 4
                  i32.const 0
                  call 74
                  local.tee 0
                  i32.sub
                  local.get 0
                  i32.const 8
                  call 50
                  i32.add
                  i32.const 20
                  i32.const 8
                  call 50
                  i32.add
                  i32.const 16
                  i32.const 8
                  call 50
                  i32.add
                  i32.const 8
                  i32.add
                  i32.const 65536
                  call 50
                  call 33
                  local.get 11
                  i32.load
                  local.tee 8
                  i32.eqz
                  br_if 6 (;@1;)
                  local.get 11
                  i32.load offset=8
                  local.set 12
                  i32.const 1049928
                  local.get 11
                  i32.load offset=4
                  local.tee 10
                  i32.const 1049928
                  i32.load
                  i32.add
                  local.tee 1
                  i32.store
                  i32.const 1049932
                  i32.const 1049932
                  i32.load
                  local.tee 0
                  local.get 1
                  local.get 0
                  local.get 1
                  i32.gt_u
                  select
                  i32.store
                  i32.const 1049924
                  i32.load
                  i32.eqz
                  br_if 1 (;@6;)
                  i32.const 1049936
                  local.set 0
                  loop  ;; label = @8
                    local.get 0
                    call 63
                    local.get 8
                    i32.eq
                    br_if 3 (;@5;)
                    local.get 0
                    i32.load offset=8
                    local.tee 0
                    br_if 0 (;@8;)
                  end
                  br 4 (;@3;)
                end
                i32.const 1049920
                i32.load
                local.set 2
                local.get 0
                local.get 4
                i32.sub
                local.tee 1
                i32.const 16
                i32.const 8
                call 50
                i32.lt_u
                if  ;; label = @7
                  i32.const 1049920
                  i32.const 0
                  i32.store
                  i32.const 1049912
                  i32.load
                  local.set 0
                  i32.const 1049912
                  i32.const 0
                  i32.store
                  local.get 2
                  local.get 0
                  call 45
                  local.get 2
                  call 74
                  local.set 3
                  br 6 (;@1;)
                end
                local.get 2
                local.get 4
                call 72
                local.set 0
                i32.const 1049912
                local.get 1
                i32.store
                i32.const 1049920
                local.get 0
                i32.store
                local.get 0
                local.get 1
                call 48
                local.get 2
                local.get 4
                call 62
                local.get 2
                call 74
                local.set 3
                br 5 (;@1;)
              end
              i32.const 1049956
              i32.load
              local.tee 0
              i32.const 0
              local.get 8
              local.get 0
              i32.ge_u
              select
              i32.eqz
              if  ;; label = @6
                i32.const 1049956
                local.get 8
                i32.store
              end
              i32.const 1049960
              i32.const 4095
              i32.store
              i32.const 1049948
              local.get 12
              i32.store
              i32.const 1049940
              local.get 10
              i32.store
              i32.const 1049936
              local.get 8
              i32.store
              i32.const 1049532
              i32.const 1049520
              i32.store
              i32.const 1049540
              i32.const 1049528
              i32.store
              i32.const 1049528
              i32.const 1049520
              i32.store
              i32.const 1049548
              i32.const 1049536
              i32.store
              i32.const 1049536
              i32.const 1049528
              i32.store
              i32.const 1049556
              i32.const 1049544
              i32.store
              i32.const 1049544
              i32.const 1049536
              i32.store
              i32.const 1049564
              i32.const 1049552
              i32.store
              i32.const 1049552
              i32.const 1049544
              i32.store
              i32.const 1049572
              i32.const 1049560
              i32.store
              i32.const 1049560
              i32.const 1049552
              i32.store
              i32.const 1049580
              i32.const 1049568
              i32.store
              i32.const 1049568
              i32.const 1049560
              i32.store
              i32.const 1049588
              i32.const 1049576
              i32.store
              i32.const 1049576
              i32.const 1049568
              i32.store
              i32.const 1049596
              i32.const 1049584
              i32.store
              i32.const 1049584
              i32.const 1049576
              i32.store
              i32.const 1049592
              i32.const 1049584
              i32.store
              i32.const 1049604
              i32.const 1049592
              i32.store
              i32.const 1049600
              i32.const 1049592
              i32.store
              i32.const 1049612
              i32.const 1049600
              i32.store
              i32.const 1049608
              i32.const 1049600
              i32.store
              i32.const 1049620
              i32.const 1049608
              i32.store
              i32.const 1049616
              i32.const 1049608
              i32.store
              i32.const 1049628
              i32.const 1049616
              i32.store
              i32.const 1049624
              i32.const 1049616
              i32.store
              i32.const 1049636
              i32.const 1049624
              i32.store
              i32.const 1049632
              i32.const 1049624
              i32.store
              i32.const 1049644
              i32.const 1049632
              i32.store
              i32.const 1049640
              i32.const 1049632
              i32.store
              i32.const 1049652
              i32.const 1049640
              i32.store
              i32.const 1049648
              i32.const 1049640
              i32.store
              i32.const 1049660
              i32.const 1049648
              i32.store
              i32.const 1049668
              i32.const 1049656
              i32.store
              i32.const 1049656
              i32.const 1049648
              i32.store
              i32.const 1049676
              i32.const 1049664
              i32.store
              i32.const 1049664
              i32.const 1049656
              i32.store
              i32.const 1049684
              i32.const 1049672
              i32.store
              i32.const 1049672
              i32.const 1049664
              i32.store
              i32.const 1049692
              i32.const 1049680
              i32.store
              i32.const 1049680
              i32.const 1049672
              i32.store
              i32.const 1049700
              i32.const 1049688
              i32.store
              i32.const 1049688
              i32.const 1049680
              i32.store
              i32.const 1049708
              i32.const 1049696
              i32.store
              i32.const 1049696
              i32.const 1049688
              i32.store
              i32.const 1049716
              i32.const 1049704
              i32.store
              i32.const 1049704
              i32.const 1049696
              i32.store
              i32.const 1049724
              i32.const 1049712
              i32.store
              i32.const 1049712
              i32.const 1049704
              i32.store
              i32.const 1049732
              i32.const 1049720
              i32.store
              i32.const 1049720
              i32.const 1049712
              i32.store
              i32.const 1049740
              i32.const 1049728
              i32.store
              i32.const 1049728
              i32.const 1049720
              i32.store
              i32.const 1049748
              i32.const 1049736
              i32.store
              i32.const 1049736
              i32.const 1049728
              i32.store
              i32.const 1049756
              i32.const 1049744
              i32.store
              i32.const 1049744
              i32.const 1049736
              i32.store
              i32.const 1049764
              i32.const 1049752
              i32.store
              i32.const 1049752
              i32.const 1049744
              i32.store
              i32.const 1049772
              i32.const 1049760
              i32.store
              i32.const 1049760
              i32.const 1049752
              i32.store
              i32.const 1049780
              i32.const 1049768
              i32.store
              i32.const 1049768
              i32.const 1049760
              i32.store
              i32.const 1049776
              i32.const 1049768
              i32.store
              i32.const 0
              call 74
              local.tee 3
              i32.const 8
              call 50
              local.set 6
              i32.const 20
              i32.const 8
              call 50
              local.set 2
              i32.const 16
              i32.const 8
              call 50
              local.set 1
              local.get 8
              local.get 8
              call 74
              local.tee 0
              i32.const 8
              call 50
              local.get 0
              i32.sub
              local.tee 0
              call 72
              local.set 5
              i32.const 1049916
              local.get 3
              local.get 10
              i32.add
              local.get 6
              i32.sub
              local.get 2
              i32.sub
              local.get 1
              i32.sub
              local.get 0
              i32.sub
              local.tee 3
              i32.store
              i32.const 1049924
              local.get 5
              i32.store
              local.get 5
              local.get 3
              i32.const 1
              i32.or
              i32.store offset=4
              i32.const 0
              call 74
              local.tee 6
              i32.const 8
              call 50
              local.set 2
              i32.const 20
              i32.const 8
              call 50
              local.set 1
              i32.const 16
              i32.const 8
              call 50
              local.set 0
              local.get 5
              local.get 3
              call 72
              local.get 0
              local.get 1
              local.get 2
              local.get 6
              i32.sub
              i32.add
              i32.add
              i32.store offset=4
              i32.const 1049952
              i32.const 2097152
              i32.store
              br 3 (;@2;)
            end
            local.get 0
            call 69
            br_if 1 (;@3;)
            local.get 0
            call 70
            local.get 12
            i32.ne
            br_if 1 (;@3;)
            local.get 0
            i32.const 1049924
            i32.load
            call 43
            i32.eqz
            br_if 1 (;@3;)
            local.get 0
            local.get 0
            i32.load offset=4
            local.get 10
            i32.add
            i32.store offset=4
            i32.const 1049916
            i32.load
            local.set 1
            i32.const 1049924
            i32.load
            local.tee 0
            local.get 0
            call 74
            local.tee 0
            i32.const 8
            call 50
            local.get 0
            i32.sub
            local.tee 0
            call 72
            local.set 5
            i32.const 1049916
            local.get 1
            local.get 10
            i32.add
            local.get 0
            i32.sub
            local.tee 3
            i32.store
            i32.const 1049924
            local.get 5
            i32.store
            local.get 5
            local.get 3
            i32.const 1
            i32.or
            i32.store offset=4
            i32.const 0
            call 74
            local.tee 6
            i32.const 8
            call 50
            local.set 2
            i32.const 20
            i32.const 8
            call 50
            local.set 1
            i32.const 16
            i32.const 8
            call 50
            local.set 0
            local.get 5
            local.get 3
            call 72
            local.get 0
            local.get 1
            local.get 2
            local.get 6
            i32.sub
            i32.add
            i32.add
            i32.store offset=4
            i32.const 1049952
            i32.const 2097152
            i32.store
            br 2 (;@2;)
          end
          i32.const 1049916
          local.get 0
          local.get 4
          i32.sub
          local.tee 1
          i32.store
          i32.const 1049924
          i32.const 1049924
          i32.load
          local.tee 2
          local.get 4
          call 72
          local.tee 0
          i32.store
          local.get 0
          local.get 1
          i32.const 1
          i32.or
          i32.store offset=4
          local.get 2
          local.get 4
          call 62
          local.get 2
          call 74
          local.set 3
          br 2 (;@1;)
        end
        i32.const 1049956
        i32.const 1049956
        i32.load
        local.tee 0
        local.get 8
        local.get 8
        local.get 0
        i32.gt_u
        select
        i32.store
        local.get 8
        local.get 10
        i32.add
        local.set 1
        i32.const 1049936
        local.set 0
        block  ;; label = @3
          loop  ;; label = @4
            local.get 1
            local.get 0
            i32.load
            i32.ne
            if  ;; label = @5
              local.get 0
              i32.load offset=8
              local.tee 0
              br_if 1 (;@4;)
              br 2 (;@3;)
            end
          end
          local.get 0
          call 69
          br_if 0 (;@3;)
          local.get 0
          call 70
          local.get 12
          i32.ne
          br_if 0 (;@3;)
          local.get 0
          i32.load
          local.set 3
          local.get 0
          local.get 8
          i32.store
          local.get 0
          local.get 0
          i32.load offset=4
          local.get 10
          i32.add
          i32.store offset=4
          local.get 8
          call 74
          local.tee 6
          i32.const 8
          call 50
          local.set 2
          local.get 3
          call 74
          local.tee 1
          i32.const 8
          call 50
          local.set 0
          local.get 8
          local.get 2
          local.get 6
          i32.sub
          i32.add
          local.tee 5
          local.get 4
          call 72
          local.set 7
          local.get 5
          local.get 4
          call 62
          local.get 3
          local.get 0
          local.get 1
          i32.sub
          i32.add
          local.tee 0
          local.get 5
          i32.sub
          local.get 4
          i32.sub
          local.set 4
          block  ;; label = @4
            local.get 0
            i32.const 1049924
            i32.load
            i32.ne
            if  ;; label = @5
              i32.const 1049920
              i32.load
              local.get 0
              i32.eq
              br_if 1 (;@4;)
              local.get 0
              call 59
              i32.eqz
              if  ;; label = @6
                block  ;; label = @7
                  local.get 0
                  call 67
                  local.tee 6
                  i32.const 256
                  i32.ge_u
                  if  ;; label = @8
                    local.get 0
                    call 14
                    br 1 (;@7;)
                  end
                  local.get 0
                  i32.const 12
                  i32.add
                  i32.load
                  local.tee 2
                  local.get 0
                  i32.const 8
                  i32.add
                  i32.load
                  local.tee 1
                  i32.ne
                  if  ;; label = @8
                    local.get 1
                    local.get 2
                    i32.store offset=12
                    local.get 2
                    local.get 1
                    i32.store offset=8
                    br 1 (;@7;)
                  end
                  i32.const 1049512
                  i32.const 1049512
                  i32.load
                  i32.const -2
                  local.get 6
                  i32.const 3
                  i32.shr_u
                  i32.rotl
                  i32.and
                  i32.store
                end
                local.get 4
                local.get 6
                i32.add
                local.set 4
                local.get 0
                local.get 6
                call 72
                local.set 0
              end
              local.get 7
              local.get 4
              local.get 0
              call 44
              local.get 4
              i32.const 256
              i32.ge_u
              if  ;; label = @6
                local.get 7
                local.get 4
                call 13
                local.get 5
                call 74
                local.set 3
                br 5 (;@1;)
              end
              local.get 4
              i32.const 3
              i32.shr_u
              local.tee 0
              i32.const 3
              i32.shl
              i32.const 1049520
              i32.add
              local.set 2
              block (result i32)  ;; label = @6
                i32.const 1049512
                i32.load
                local.tee 1
                i32.const 1
                local.get 0
                i32.shl
                local.tee 0
                i32.and
                if  ;; label = @7
                  local.get 2
                  i32.load offset=8
                  br 1 (;@6;)
                end
                i32.const 1049512
                local.get 0
                local.get 1
                i32.or
                i32.store
                local.get 2
              end
              local.set 0
              local.get 2
              local.get 7
              i32.store offset=8
              local.get 0
              local.get 7
              i32.store offset=12
              local.get 7
              local.get 2
              i32.store offset=12
              local.get 7
              local.get 0
              i32.store offset=8
              local.get 5
              call 74
              local.set 3
              br 4 (;@1;)
            end
            i32.const 1049924
            local.get 7
            i32.store
            i32.const 1049916
            i32.const 1049916
            i32.load
            local.get 4
            i32.add
            local.tee 0
            i32.store
            local.get 7
            local.get 0
            i32.const 1
            i32.or
            i32.store offset=4
            local.get 5
            call 74
            local.set 3
            br 3 (;@1;)
          end
          i32.const 1049920
          local.get 7
          i32.store
          i32.const 1049912
          i32.const 1049912
          i32.load
          local.get 4
          i32.add
          local.tee 0
          i32.store
          local.get 7
          local.get 0
          call 48
          local.get 5
          call 74
          local.set 3
          br 2 (;@1;)
        end
        i32.const 1049924
        i32.load
        local.set 9
        i32.const 1049936
        local.set 0
        block  ;; label = @3
          loop  ;; label = @4
            local.get 0
            i32.load
            local.get 9
            i32.le_u
            if  ;; label = @5
              local.get 0
              call 63
              local.get 9
              i32.gt_u
              br_if 2 (;@3;)
            end
            local.get 0
            i32.load offset=8
            local.tee 0
            br_if 0 (;@4;)
          end
          i32.const 0
          local.set 0
        end
        local.get 9
        local.get 0
        call 63
        local.tee 7
        i32.const 20
        i32.const 8
        call 50
        local.tee 16
        i32.sub
        i32.const -23
        i32.add
        local.tee 1
        call 74
        local.tee 0
        i32.const 8
        call 50
        local.get 0
        i32.sub
        local.get 1
        i32.add
        local.tee 0
        local.get 0
        i32.const 16
        i32.const 8
        call 50
        local.get 9
        i32.add
        i32.lt_u
        select
        local.tee 13
        call 74
        local.set 14
        local.get 13
        local.get 16
        call 72
        local.set 0
        i32.const 0
        call 74
        local.tee 5
        i32.const 8
        call 50
        local.set 3
        i32.const 20
        i32.const 8
        call 50
        local.set 6
        i32.const 16
        i32.const 8
        call 50
        local.set 2
        local.get 8
        local.get 8
        call 74
        local.tee 1
        i32.const 8
        call 50
        local.get 1
        i32.sub
        local.tee 1
        call 72
        local.set 15
        i32.const 1049916
        local.get 5
        local.get 10
        i32.add
        local.get 3
        i32.sub
        local.get 6
        i32.sub
        local.get 2
        i32.sub
        local.get 1
        i32.sub
        local.tee 5
        i32.store
        i32.const 1049924
        local.get 15
        i32.store
        local.get 15
        local.get 5
        i32.const 1
        i32.or
        i32.store offset=4
        i32.const 0
        call 74
        local.tee 3
        i32.const 8
        call 50
        local.set 6
        i32.const 20
        i32.const 8
        call 50
        local.set 2
        i32.const 16
        i32.const 8
        call 50
        local.set 1
        local.get 15
        local.get 5
        call 72
        local.get 1
        local.get 2
        local.get 6
        local.get 3
        i32.sub
        i32.add
        i32.add
        i32.store offset=4
        i32.const 1049952
        i32.const 2097152
        i32.store
        local.get 13
        local.get 16
        call 62
        i32.const 1049936
        i64.load align=4
        local.set 17
        local.get 14
        i32.const 8
        i32.add
        i32.const 1049944
        i64.load align=4
        i64.store align=4
        local.get 14
        local.get 17
        i64.store align=4
        i32.const 1049948
        local.get 12
        i32.store
        i32.const 1049940
        local.get 10
        i32.store
        i32.const 1049936
        local.get 8
        i32.store
        i32.const 1049944
        local.get 14
        i32.store
        loop  ;; label = @3
          local.get 0
          i32.const 4
          call 72
          local.set 1
          local.get 0
          i32.const 7
          i32.store offset=4
          local.get 7
          local.get 1
          local.tee 0
          i32.const 4
          i32.add
          i32.gt_u
          br_if 0 (;@3;)
        end
        local.get 9
        local.get 13
        i32.eq
        br_if 0 (;@2;)
        local.get 9
        local.get 13
        local.get 9
        i32.sub
        local.tee 0
        local.get 9
        local.get 0
        call 72
        call 44
        local.get 0
        i32.const 256
        i32.ge_u
        if  ;; label = @3
          local.get 9
          local.get 0
          call 13
          br 1 (;@2;)
        end
        local.get 0
        i32.const 3
        i32.shr_u
        local.tee 0
        i32.const 3
        i32.shl
        i32.const 1049520
        i32.add
        local.set 2
        block (result i32)  ;; label = @3
          i32.const 1049512
          i32.load
          local.tee 1
          i32.const 1
          local.get 0
          i32.shl
          local.tee 0
          i32.and
          if  ;; label = @4
            local.get 2
            i32.load offset=8
            br 1 (;@3;)
          end
          i32.const 1049512
          local.get 0
          local.get 1
          i32.or
          i32.store
          local.get 2
        end
        local.set 0
        local.get 2
        local.get 9
        i32.store offset=8
        local.get 0
        local.get 9
        i32.store offset=12
        local.get 9
        local.get 2
        i32.store offset=12
        local.get 9
        local.get 0
        i32.store offset=8
      end
      i32.const 0
      local.set 3
      i32.const 1049916
      i32.load
      local.tee 0
      local.get 4
      i32.le_u
      br_if 0 (;@1;)
      i32.const 1049916
      local.get 0
      local.get 4
      i32.sub
      local.tee 1
      i32.store
      i32.const 1049924
      i32.const 1049924
      i32.load
      local.tee 2
      local.get 4
      call 72
      local.tee 0
      i32.store
      local.get 0
      local.get 1
      i32.const 1
      i32.or
      i32.store offset=4
      local.get 2
      local.get 4
      call 62
      local.get 2
      call 74
      local.set 3
    end
    local.get 11
    i32.const 16
    i32.add
    global.set 0
    local.get 3)
  (func (;4;) (type 2) (param i32)
    (local i32 i32 i32 i32 i32 i32)
    local.get 0
    call 75
    local.tee 0
    local.get 0
    call 67
    local.tee 2
    call 72
    local.set 1
    block  ;; label = @1
      block  ;; label = @2
        local.get 0
        call 68
        br_if 0 (;@2;)
        local.get 0
        i32.load
        local.set 3
        block  ;; label = @3
          local.get 0
          call 61
          i32.eqz
          if  ;; label = @4
            local.get 2
            local.get 3
            i32.add
            local.set 2
            local.get 0
            local.get 3
            call 73
            local.tee 0
            i32.const 1049920
            i32.load
            i32.ne
            br_if 1 (;@3;)
            local.get 1
            i32.load offset=4
            i32.const 3
            i32.and
            i32.const 3
            i32.ne
            br_if 2 (;@2;)
            i32.const 1049912
            local.get 2
            i32.store
            local.get 0
            local.get 2
            local.get 1
            call 44
            return
          end
          local.get 2
          local.get 3
          i32.add
          i32.const 16
          i32.add
          local.set 0
          br 2 (;@1;)
        end
        local.get 3
        i32.const 256
        i32.ge_u
        if  ;; label = @3
          local.get 0
          call 14
          br 1 (;@2;)
        end
        local.get 0
        i32.const 12
        i32.add
        i32.load
        local.tee 4
        local.get 0
        i32.const 8
        i32.add
        i32.load
        local.tee 5
        i32.ne
        if  ;; label = @3
          local.get 5
          local.get 4
          i32.store offset=12
          local.get 4
          local.get 5
          i32.store offset=8
          br 1 (;@2;)
        end
        i32.const 1049512
        i32.const 1049512
        i32.load
        i32.const -2
        local.get 3
        i32.const 3
        i32.shr_u
        i32.rotl
        i32.and
        i32.store
      end
      block  ;; label = @2
        local.get 1
        call 58
        if  ;; label = @3
          local.get 0
          local.get 2
          local.get 1
          call 44
          br 1 (;@2;)
        end
        block  ;; label = @3
          i32.const 1049924
          i32.load
          local.get 1
          i32.ne
          if  ;; label = @4
            local.get 1
            i32.const 1049920
            i32.load
            i32.ne
            br_if 1 (;@3;)
            i32.const 1049920
            local.get 0
            i32.store
            i32.const 1049912
            i32.const 1049912
            i32.load
            local.get 2
            i32.add
            local.tee 1
            i32.store
            local.get 0
            local.get 1
            call 48
            return
          end
          i32.const 1049924
          local.get 0
          i32.store
          i32.const 1049916
          i32.const 1049916
          i32.load
          local.get 2
          i32.add
          local.tee 1
          i32.store
          local.get 0
          local.get 1
          i32.const 1
          i32.or
          i32.store offset=4
          i32.const 1049920
          i32.load
          local.get 0
          i32.eq
          if  ;; label = @4
            i32.const 1049912
            i32.const 0
            i32.store
            i32.const 1049920
            i32.const 0
            i32.store
          end
          i32.const 1049952
          i32.load
          local.get 1
          i32.ge_u
          br_if 2 (;@1;)
          i32.const 0
          call 74
          local.tee 0
          i32.const 8
          call 50
          local.set 1
          i32.const 20
          i32.const 8
          call 50
          local.set 3
          i32.const 16
          i32.const 8
          call 50
          local.set 2
          i32.const 16
          i32.const 8
          call 50
          local.set 4
          i32.const 1049924
          i32.load
          i32.eqz
          br_if 2 (;@1;)
          local.get 0
          local.get 1
          i32.sub
          local.get 3
          i32.sub
          local.get 2
          i32.sub
          i32.const -65544
          i32.add
          i32.const -9
          i32.and
          i32.const -3
          i32.add
          local.tee 0
          i32.const 0
          local.get 4
          i32.const 2
          i32.shl
          i32.sub
          local.tee 1
          local.get 0
          local.get 1
          i32.lt_u
          select
          i32.eqz
          br_if 2 (;@1;)
          i32.const 0
          call 74
          local.tee 0
          i32.const 8
          call 50
          local.set 1
          i32.const 20
          i32.const 8
          call 50
          local.set 2
          i32.const 16
          i32.const 8
          call 50
          local.set 4
          i32.const 0
          block  ;; label = @4
            i32.const 1049916
            i32.load
            local.tee 5
            local.get 4
            local.get 2
            local.get 1
            local.get 0
            i32.sub
            i32.add
            i32.add
            local.tee 2
            i32.le_u
            br_if 0 (;@4;)
            i32.const 1049924
            i32.load
            local.set 1
            i32.const 1049936
            local.set 0
            block  ;; label = @5
              loop  ;; label = @6
                local.get 0
                i32.load
                local.get 1
                i32.le_u
                if  ;; label = @7
                  local.get 0
                  call 63
                  local.get 1
                  i32.gt_u
                  br_if 2 (;@5;)
                end
                local.get 0
                i32.load offset=8
                local.tee 0
                br_if 0 (;@6;)
              end
              i32.const 0
              local.set 0
            end
            local.get 0
            call 69
            br_if 0 (;@4;)
            local.get 0
            i32.const 12
            i32.add
            i32.load
            drop
            br 0 (;@4;)
          end
          i32.const 0
          call 12
          i32.sub
          i32.ne
          br_if 2 (;@1;)
          i32.const 1049916
          i32.load
          i32.const 1049952
          i32.load
          i32.le_u
          br_if 2 (;@1;)
          i32.const 1049952
          i32.const -1
          i32.store
          return
        end
        local.get 1
        call 67
        local.tee 3
        local.get 2
        i32.add
        local.set 2
        block  ;; label = @3
          local.get 3
          i32.const 256
          i32.ge_u
          if  ;; label = @4
            local.get 1
            call 14
            br 1 (;@3;)
          end
          local.get 1
          i32.const 12
          i32.add
          i32.load
          local.tee 4
          local.get 1
          i32.const 8
          i32.add
          i32.load
          local.tee 1
          i32.ne
          if  ;; label = @4
            local.get 1
            local.get 4
            i32.store offset=12
            local.get 4
            local.get 1
            i32.store offset=8
            br 1 (;@3;)
          end
          i32.const 1049512
          i32.const 1049512
          i32.load
          i32.const -2
          local.get 3
          i32.const 3
          i32.shr_u
          i32.rotl
          i32.and
          i32.store
        end
        local.get 0
        local.get 2
        call 48
        local.get 0
        i32.const 1049920
        i32.load
        i32.ne
        br_if 0 (;@2;)
        i32.const 1049912
        local.get 2
        i32.store
        return
      end
      local.get 2
      i32.const 256
      i32.ge_u
      if  ;; label = @2
        local.get 0
        local.get 2
        call 13
        i32.const 1049960
        i32.const 1049960
        i32.load
        i32.const -1
        i32.add
        local.tee 0
        i32.store
        local.get 0
        br_if 1 (;@1;)
        call 12
        drop
        return
      end
      local.get 2
      i32.const 3
      i32.shr_u
      local.tee 3
      i32.const 3
      i32.shl
      i32.const 1049520
      i32.add
      local.set 1
      block (result i32)  ;; label = @2
        i32.const 1049512
        i32.load
        local.tee 2
        i32.const 1
        local.get 3
        i32.shl
        local.tee 3
        i32.and
        if  ;; label = @3
          local.get 1
          i32.load offset=8
          br 1 (;@2;)
        end
        i32.const 1049512
        local.get 2
        local.get 3
        i32.or
        i32.store
        local.get 1
      end
      local.set 3
      local.get 1
      local.get 0
      i32.store offset=8
      local.get 3
      local.get 0
      i32.store offset=12
      local.get 0
      local.get 1
      i32.store offset=12
      local.get 0
      local.get 3
      i32.store offset=8
    end)
  (func (;5;) (type 9) (param i32 i32 i32 i32) (result i32)
    (local i32 i32 i32 i32 i32 i32)
    block  ;; label = @1
      block  ;; label = @2
        block  ;; label = @3
          local.get 2
          i32.const 9
          i32.ge_u
          if  ;; label = @4
            local.get 3
            local.get 2
            call 10
            local.tee 2
            br_if 1 (;@3;)
            i32.const 0
            return
          end
          i32.const 0
          local.set 2
          i32.const 0
          call 74
          local.tee 1
          local.get 1
          i32.const 8
          call 50
          i32.sub
          i32.const 20
          i32.const 8
          call 50
          i32.sub
          i32.const 16
          i32.const 8
          call 50
          i32.sub
          i32.const -65544
          i32.add
          i32.const -9
          i32.and
          i32.const -3
          i32.add
          local.tee 1
          i32.const 0
          i32.const 16
          i32.const 8
          call 50
          i32.const 2
          i32.shl
          i32.sub
          local.tee 4
          local.get 1
          local.get 4
          i32.lt_u
          select
          local.get 3
          i32.le_u
          br_if 1 (;@2;)
          i32.const 16
          local.get 3
          i32.const 4
          i32.add
          i32.const 16
          i32.const 8
          call 50
          i32.const -5
          i32.add
          local.get 3
          i32.gt_u
          select
          i32.const 8
          call 50
          local.set 6
          local.get 0
          call 75
          local.tee 1
          local.get 1
          call 67
          local.tee 5
          call 72
          local.set 4
          block  ;; label = @4
            block  ;; label = @5
              block  ;; label = @6
                block  ;; label = @7
                  block  ;; label = @8
                    block  ;; label = @9
                      block  ;; label = @10
                        block  ;; label = @11
                          local.get 1
                          call 61
                          i32.eqz
                          if  ;; label = @12
                            local.get 5
                            local.get 6
                            i32.ge_u
                            br_if 1 (;@11;)
                            local.get 4
                            i32.const 1049924
                            i32.load
                            i32.eq
                            br_if 2 (;@10;)
                            local.get 4
                            i32.const 1049920
                            i32.load
                            i32.eq
                            br_if 3 (;@9;)
                            local.get 4
                            call 58
                            br_if 8 (;@4;)
                            local.get 4
                            call 67
                            local.tee 7
                            local.get 5
                            i32.add
                            local.tee 8
                            local.get 6
                            i32.lt_u
                            br_if 8 (;@4;)
                            local.get 8
                            local.get 6
                            i32.sub
                            local.set 5
                            local.get 7
                            i32.const 256
                            i32.lt_u
                            br_if 4 (;@8;)
                            local.get 4
                            call 14
                            br 5 (;@7;)
                          end
                          local.get 1
                          call 67
                          local.set 5
                          local.get 6
                          i32.const 256
                          i32.lt_u
                          br_if 7 (;@4;)
                          local.get 5
                          local.get 6
                          i32.const 4
                          i32.add
                          i32.ge_u
                          if  ;; label = @12
                            local.get 1
                            local.set 4
                            local.get 5
                            local.get 6
                            i32.sub
                            i32.const 131073
                            i32.lt_u
                            br_if 7 (;@5;)
                          end
                          local.get 1
                          i32.load
                          local.tee 7
                          local.get 5
                          i32.add
                          i32.const 16
                          i32.add
                          local.set 8
                          local.get 6
                          i32.const 31
                          i32.add
                          i32.const 65536
                          call 50
                          local.set 5
                          i32.const 0
                          local.tee 6
                          i32.eqz
                          br_if 7 (;@4;)
                          local.get 6
                          local.get 7
                          i32.add
                          local.tee 4
                          local.get 5
                          local.get 7
                          i32.sub
                          local.tee 7
                          i32.const -16
                          i32.add
                          local.tee 9
                          i32.store offset=4
                          local.get 4
                          local.get 9
                          call 72
                          i32.const 7
                          i32.store offset=4
                          local.get 4
                          local.get 7
                          i32.const -12
                          i32.add
                          call 72
                          i32.const 0
                          i32.store offset=4
                          i32.const 1049928
                          i32.const 1049928
                          i32.load
                          local.get 5
                          local.get 8
                          i32.sub
                          i32.add
                          local.tee 5
                          i32.store
                          i32.const 1049956
                          i32.const 1049956
                          i32.load
                          local.tee 7
                          local.get 6
                          local.get 6
                          local.get 7
                          i32.gt_u
                          select
                          i32.store
                          i32.const 1049932
                          i32.const 1049932
                          i32.load
                          local.tee 6
                          local.get 5
                          local.get 6
                          local.get 5
                          i32.gt_u
                          select
                          i32.store
                          br 6 (;@5;)
                        end
                        local.get 1
                        local.set 4
                        local.get 5
                        local.get 6
                        i32.sub
                        local.tee 5
                        i32.const 16
                        i32.const 8
                        call 50
                        i32.lt_u
                        br_if 5 (;@5;)
                        local.get 1
                        local.get 6
                        call 72
                        local.set 4
                        local.get 1
                        local.get 6
                        call 40
                        local.get 4
                        local.get 5
                        call 40
                        local.get 4
                        local.get 5
                        call 8
                        br 4 (;@6;)
                      end
                      i32.const 1049916
                      i32.load
                      local.get 5
                      i32.add
                      local.tee 5
                      local.get 6
                      i32.le_u
                      br_if 5 (;@4;)
                      local.get 1
                      local.get 6
                      call 72
                      local.set 4
                      local.get 1
                      local.get 6
                      call 40
                      local.get 4
                      local.get 5
                      local.get 6
                      i32.sub
                      local.tee 6
                      i32.const 1
                      i32.or
                      i32.store offset=4
                      i32.const 1049916
                      local.get 6
                      i32.store
                      i32.const 1049924
                      local.get 4
                      i32.store
                      br 3 (;@6;)
                    end
                    i32.const 1049912
                    i32.load
                    local.get 5
                    i32.add
                    local.tee 4
                    local.get 6
                    i32.lt_u
                    br_if 4 (;@4;)
                    block  ;; label = @9
                      local.get 4
                      local.get 6
                      i32.sub
                      local.tee 5
                      i32.const 16
                      i32.const 8
                      call 50
                      i32.lt_u
                      if  ;; label = @10
                        local.get 1
                        local.get 4
                        call 40
                        i32.const 0
                        local.set 5
                        i32.const 0
                        local.set 4
                        br 1 (;@9;)
                      end
                      local.get 1
                      local.get 6
                      call 72
                      local.tee 4
                      local.get 5
                      call 72
                      local.set 7
                      local.get 1
                      local.get 6
                      call 40
                      local.get 4
                      local.get 5
                      call 48
                      local.get 7
                      local.get 7
                      i32.load offset=4
                      i32.const -2
                      i32.and
                      i32.store offset=4
                    end
                    i32.const 1049920
                    local.get 4
                    i32.store
                    i32.const 1049912
                    local.get 5
                    i32.store
                    br 2 (;@6;)
                  end
                  local.get 4
                  i32.const 12
                  i32.add
                  i32.load
                  local.tee 9
                  local.get 4
                  i32.const 8
                  i32.add
                  i32.load
                  local.tee 4
                  i32.ne
                  if  ;; label = @8
                    local.get 4
                    local.get 9
                    i32.store offset=12
                    local.get 9
                    local.get 4
                    i32.store offset=8
                    br 1 (;@7;)
                  end
                  i32.const 1049512
                  i32.const 1049512
                  i32.load
                  i32.const -2
                  local.get 7
                  i32.const 3
                  i32.shr_u
                  i32.rotl
                  i32.and
                  i32.store
                end
                local.get 5
                i32.const 16
                i32.const 8
                call 50
                i32.ge_u
                if  ;; label = @7
                  local.get 1
                  local.get 6
                  call 72
                  local.set 4
                  local.get 1
                  local.get 6
                  call 40
                  local.get 4
                  local.get 5
                  call 40
                  local.get 4
                  local.get 5
                  call 8
                  br 1 (;@6;)
                end
                local.get 1
                local.get 8
                call 40
              end
              local.get 1
              local.set 4
            end
            local.get 4
            br_if 3 (;@1;)
          end
          local.get 3
          call 3
          local.tee 4
          i32.eqz
          br_if 1 (;@2;)
          local.get 4
          local.get 0
          local.get 3
          local.get 1
          call 67
          i32.const -8
          i32.const -4
          local.get 1
          call 61
          select
          i32.add
          local.tee 1
          local.get 1
          local.get 3
          i32.gt_u
          select
          call 36
          local.get 0
          call 4
          return
        end
        local.get 2
        local.get 0
        local.get 3
        local.get 1
        local.get 1
        local.get 3
        i32.gt_u
        select
        call 36
        drop
        local.get 0
        call 4
      end
      local.get 2
      return
    end
    local.get 4
    call 61
    drop
    local.get 4
    call 74)
  (func (;6;) (type 4) (param i32 i32 i32) (result i32)
    (local i32 i32 i32 i32 i32 i32 i32)
    i32.const 43
    i32.const 1114112
    local.get 0
    i32.load
    local.tee 3
    i32.const 1
    i32.and
    local.tee 4
    select
    local.set 6
    local.get 2
    local.get 4
    i32.add
    local.set 4
    i32.const 1049144
    i32.const 0
    local.get 3
    i32.const 4
    i32.and
    select
    local.set 7
    block  ;; label = @1
      block  ;; label = @2
        local.get 0
        i32.load offset=8
        i32.const 1
        i32.ne
        if  ;; label = @3
          local.get 0
          local.get 6
          local.get 7
          call 29
          br_if 1 (;@2;)
          br 2 (;@1;)
        end
        local.get 0
        i32.const 12
        i32.add
        i32.load
        local.tee 5
        local.get 4
        i32.le_u
        if  ;; label = @3
          local.get 0
          local.get 6
          local.get 7
          call 29
          br_if 1 (;@2;)
          br 2 (;@1;)
        end
        block  ;; label = @3
          block  ;; label = @4
            block  ;; label = @5
              block  ;; label = @6
                local.get 3
                i32.const 8
                i32.and
                if  ;; label = @7
                  local.get 0
                  i32.load offset=4
                  local.set 8
                  local.get 0
                  i32.const 48
                  i32.store offset=4
                  local.get 0
                  i32.load8_u offset=32
                  local.set 9
                  local.get 0
                  i32.const 1
                  i32.store8 offset=32
                  local.get 0
                  local.get 6
                  local.get 7
                  call 29
                  br_if 5 (;@2;)
                  i32.const 0
                  local.set 3
                  local.get 5
                  local.get 4
                  i32.sub
                  local.tee 4
                  local.set 5
                  i32.const 1
                  local.get 0
                  i32.load8_u offset=32
                  local.tee 6
                  local.get 6
                  i32.const 3
                  i32.eq
                  select
                  i32.const 3
                  i32.and
                  i32.const 1
                  i32.sub
                  br_table 2 (;@5;) 1 (;@6;) 2 (;@5;) 3 (;@4;)
                end
                i32.const 0
                local.set 3
                local.get 5
                local.get 4
                i32.sub
                local.tee 4
                local.set 5
                block  ;; label = @7
                  block  ;; label = @8
                    block  ;; label = @9
                      i32.const 1
                      local.get 0
                      i32.load8_u offset=32
                      local.tee 8
                      local.get 8
                      i32.const 3
                      i32.eq
                      select
                      i32.const 3
                      i32.and
                      i32.const 1
                      i32.sub
                      br_table 1 (;@8;) 0 (;@9;) 1 (;@8;) 2 (;@7;)
                    end
                    local.get 4
                    i32.const 1
                    i32.shr_u
                    local.set 3
                    local.get 4
                    i32.const 1
                    i32.add
                    i32.const 1
                    i32.shr_u
                    local.set 5
                    br 1 (;@7;)
                  end
                  i32.const 0
                  local.set 5
                  local.get 4
                  local.set 3
                end
                local.get 3
                i32.const 1
                i32.add
                local.set 3
                loop  ;; label = @7
                  local.get 3
                  i32.const -1
                  i32.add
                  local.tee 3
                  i32.eqz
                  br_if 4 (;@3;)
                  local.get 0
                  i32.load offset=24
                  local.get 0
                  i32.load offset=4
                  local.get 0
                  i32.load offset=28
                  i32.load offset=16
                  call_indirect (type 1)
                  i32.eqz
                  br_if 0 (;@7;)
                end
                i32.const 1
                return
              end
              local.get 4
              i32.const 1
              i32.shr_u
              local.set 3
              local.get 4
              i32.const 1
              i32.add
              i32.const 1
              i32.shr_u
              local.set 5
              br 1 (;@4;)
            end
            i32.const 0
            local.set 5
            local.get 4
            local.set 3
          end
          local.get 3
          i32.const 1
          i32.add
          local.set 3
          block  ;; label = @4
            loop  ;; label = @5
              local.get 3
              i32.const -1
              i32.add
              local.tee 3
              i32.eqz
              br_if 1 (;@4;)
              local.get 0
              i32.load offset=24
              local.get 0
              i32.load offset=4
              local.get 0
              i32.load offset=28
              i32.load offset=16
              call_indirect (type 1)
              i32.eqz
              br_if 0 (;@5;)
            end
            i32.const 1
            return
          end
          local.get 0
          i32.load offset=4
          local.set 4
          local.get 0
          i32.load offset=24
          local.get 1
          local.get 2
          local.get 0
          i32.load offset=28
          i32.load offset=12
          call_indirect (type 4)
          br_if 1 (;@2;)
          local.get 5
          i32.const 1
          i32.add
          local.set 3
          local.get 0
          i32.load offset=28
          local.set 1
          local.get 0
          i32.load offset=24
          local.set 2
          loop  ;; label = @4
            local.get 3
            i32.const -1
            i32.add
            local.tee 3
            if  ;; label = @5
              local.get 2
              local.get 4
              local.get 1
              i32.load offset=16
              call_indirect (type 1)
              i32.eqz
              br_if 1 (;@4;)
              br 3 (;@2;)
            end
          end
          local.get 0
          local.get 9
          i32.store8 offset=32
          local.get 0
          local.get 8
          i32.store offset=4
          i32.const 0
          return
        end
        local.get 0
        i32.load offset=4
        local.set 4
        local.get 0
        local.get 6
        local.get 7
        call 29
        br_if 0 (;@2;)
        local.get 0
        i32.load offset=24
        local.get 1
        local.get 2
        local.get 0
        i32.load offset=28
        i32.load offset=12
        call_indirect (type 4)
        br_if 0 (;@2;)
        local.get 5
        i32.const 1
        i32.add
        local.set 3
        local.get 0
        i32.load offset=28
        local.set 1
        local.get 0
        i32.load offset=24
        local.set 0
        loop  ;; label = @3
          local.get 3
          i32.const -1
          i32.add
          local.tee 3
          i32.eqz
          if  ;; label = @4
            i32.const 0
            return
          end
          local.get 0
          local.get 4
          local.get 1
          i32.load offset=16
          call_indirect (type 1)
          i32.eqz
          br_if 0 (;@3;)
        end
      end
      i32.const 1
      return
    end
    local.get 0
    i32.load offset=24
    local.get 1
    local.get 2
    local.get 0
    i32.const 28
    i32.add
    i32.load
    i32.load offset=12
    call_indirect (type 4))
  (func (;7;) (type 1) (param i32 i32) (result i32)
    (local i32 i32 i32 i32 i32 i32 i32 i32 i32 i32)
    global.get 0
    i32.const 48
    i32.sub
    local.tee 2
    global.set 0
    local.get 2
    i32.const 36
    i32.add
    i32.const 1048864
    i32.store
    local.get 2
    i32.const 3
    i32.store8 offset=40
    local.get 2
    i64.const 137438953472
    i64.store offset=8
    local.get 2
    local.get 0
    i32.store offset=32
    local.get 2
    i32.const 0
    i32.store offset=24
    local.get 2
    i32.const 0
    i32.store offset=16
    block (result i32)  ;; label = @1
      block  ;; label = @2
        block  ;; label = @3
          block  ;; label = @4
            local.get 1
            i32.load offset=8
            local.tee 3
            if  ;; label = @5
              local.get 1
              i32.load
              local.set 5
              local.get 1
              i32.load offset=4
              local.tee 8
              local.get 1
              i32.const 12
              i32.add
              i32.load
              local.tee 4
              local.get 4
              local.get 8
              i32.gt_u
              select
              local.tee 4
              i32.eqz
              br_if 1 (;@4;)
              local.get 0
              local.get 5
              i32.load
              local.get 5
              i32.load offset=4
              i32.const 1048876
              i32.load
              call_indirect (type 4)
              br_if 3 (;@2;)
              local.get 5
              i32.const 12
              i32.add
              local.set 0
              local.get 1
              i32.load offset=16
              local.set 7
              local.get 4
              local.set 9
              loop  ;; label = @6
                local.get 2
                local.get 3
                i32.const 28
                i32.add
                i32.load8_u
                i32.store8 offset=40
                local.get 2
                local.get 3
                i32.const 4
                i32.add
                i64.load align=4
                i64.const 32
                i64.rotl
                i64.store offset=8
                local.get 3
                i32.const 24
                i32.add
                i32.load
                local.set 6
                i32.const 0
                local.set 10
                i32.const 0
                local.set 1
                block  ;; label = @7
                  block  ;; label = @8
                    block  ;; label = @9
                      local.get 3
                      i32.const 20
                      i32.add
                      i32.load
                      i32.const 1
                      i32.sub
                      br_table 0 (;@9;) 2 (;@7;) 1 (;@8;)
                    end
                    local.get 6
                    i32.const 3
                    i32.shl
                    local.get 7
                    i32.add
                    local.tee 11
                    i32.load offset=4
                    i32.const 18
                    i32.ne
                    br_if 1 (;@7;)
                    local.get 11
                    i32.load
                    i32.load
                    local.set 6
                  end
                  i32.const 1
                  local.set 1
                end
                local.get 2
                local.get 6
                i32.store offset=20
                local.get 2
                local.get 1
                i32.store offset=16
                local.get 3
                i32.const 16
                i32.add
                i32.load
                local.set 1
                block  ;; label = @7
                  block  ;; label = @8
                    block  ;; label = @9
                      local.get 3
                      i32.const 12
                      i32.add
                      i32.load
                      i32.const 1
                      i32.sub
                      br_table 0 (;@9;) 2 (;@7;) 1 (;@8;)
                    end
                    local.get 1
                    i32.const 3
                    i32.shl
                    local.get 7
                    i32.add
                    local.tee 6
                    i32.load offset=4
                    i32.const 18
                    i32.ne
                    br_if 1 (;@7;)
                    local.get 6
                    i32.load
                    i32.load
                    local.set 1
                  end
                  i32.const 1
                  local.set 10
                end
                local.get 2
                local.get 1
                i32.store offset=28
                local.get 2
                local.get 10
                i32.store offset=24
                local.get 7
                local.get 3
                i32.load
                i32.const 3
                i32.shl
                i32.add
                local.tee 1
                i32.load
                local.get 2
                i32.const 8
                i32.add
                local.get 1
                i32.load offset=4
                call_indirect (type 1)
                br_if 4 (;@2;)
                local.get 9
                i32.const -1
                i32.add
                local.tee 9
                i32.eqz
                br_if 3 (;@3;)
                local.get 3
                i32.const 32
                i32.add
                local.set 3
                local.get 0
                i32.const -4
                i32.add
                local.set 1
                local.get 0
                i32.load
                local.set 6
                local.get 0
                i32.const 8
                i32.add
                local.set 0
                local.get 2
                i32.load offset=32
                local.get 1
                i32.load
                local.get 6
                local.get 2
                i32.load offset=36
                i32.load offset=12
                call_indirect (type 4)
                i32.eqz
                br_if 0 (;@6;)
              end
              br 3 (;@2;)
            end
            local.get 1
            i32.load
            local.set 5
            local.get 1
            i32.load offset=4
            local.tee 8
            local.get 1
            i32.const 20
            i32.add
            i32.load
            local.tee 4
            local.get 4
            local.get 8
            i32.gt_u
            select
            local.tee 4
            i32.eqz
            br_if 0 (;@4;)
            local.get 1
            i32.load offset=16
            local.set 3
            local.get 0
            local.get 5
            i32.load
            local.get 5
            i32.load offset=4
            i32.const 1048876
            i32.load
            call_indirect (type 4)
            br_if 2 (;@2;)
            local.get 5
            i32.const 12
            i32.add
            local.set 0
            local.get 4
            local.set 1
            loop  ;; label = @5
              local.get 3
              i32.load
              local.get 2
              i32.const 8
              i32.add
              local.get 3
              i32.const 4
              i32.add
              i32.load
              call_indirect (type 1)
              br_if 3 (;@2;)
              local.get 1
              i32.const -1
              i32.add
              local.tee 1
              i32.eqz
              br_if 2 (;@3;)
              local.get 3
              i32.const 8
              i32.add
              local.set 3
              local.get 0
              i32.const -4
              i32.add
              local.set 9
              local.get 0
              i32.load
              local.set 7
              local.get 0
              i32.const 8
              i32.add
              local.set 0
              local.get 2
              i32.load offset=32
              local.get 9
              i32.load
              local.get 7
              local.get 2
              i32.load offset=36
              i32.load offset=12
              call_indirect (type 4)
              i32.eqz
              br_if 0 (;@5;)
            end
            br 2 (;@2;)
          end
          i32.const 0
          local.set 4
        end
        local.get 8
        local.get 4
        i32.gt_u
        if  ;; label = @3
          local.get 2
          i32.load offset=32
          local.get 5
          local.get 4
          i32.const 3
          i32.shl
          i32.add
          local.tee 0
          i32.load
          local.get 0
          i32.load offset=4
          local.get 2
          i32.load offset=36
          i32.load offset=12
          call_indirect (type 4)
          br_if 1 (;@2;)
        end
        i32.const 0
        br 1 (;@1;)
      end
      i32.const 1
    end
    local.get 2
    i32.const 48
    i32.add
    global.set 0)
  (func (;8;) (type 0) (param i32 i32)
    (local i32 i32 i32 i32)
    local.get 0
    local.get 1
    call 72
    local.set 2
    block  ;; label = @1
      block  ;; label = @2
        local.get 0
        call 68
        br_if 0 (;@2;)
        local.get 0
        i32.load
        local.set 3
        block  ;; label = @3
          local.get 0
          call 61
          i32.eqz
          if  ;; label = @4
            local.get 1
            local.get 3
            i32.add
            local.set 1
            local.get 0
            local.get 3
            call 73
            local.tee 0
            i32.const 1049920
            i32.load
            i32.ne
            br_if 1 (;@3;)
            local.get 2
            i32.load offset=4
            i32.const 3
            i32.and
            i32.const 3
            i32.ne
            br_if 2 (;@2;)
            i32.const 1049912
            local.get 1
            i32.store
            local.get 0
            local.get 1
            local.get 2
            call 44
            return
          end
          local.get 1
          local.get 3
          i32.add
          i32.const 16
          i32.add
          local.set 0
          br 2 (;@1;)
        end
        local.get 3
        i32.const 256
        i32.ge_u
        if  ;; label = @3
          local.get 0
          call 14
          br 1 (;@2;)
        end
        local.get 0
        i32.const 12
        i32.add
        i32.load
        local.tee 4
        local.get 0
        i32.const 8
        i32.add
        i32.load
        local.tee 5
        i32.ne
        if  ;; label = @3
          local.get 5
          local.get 4
          i32.store offset=12
          local.get 4
          local.get 5
          i32.store offset=8
          br 1 (;@2;)
        end
        i32.const 1049512
        i32.const 1049512
        i32.load
        i32.const -2
        local.get 3
        i32.const 3
        i32.shr_u
        i32.rotl
        i32.and
        i32.store
      end
      block  ;; label = @2
        local.get 2
        call 58
        if  ;; label = @3
          local.get 0
          local.get 1
          local.get 2
          call 44
          br 1 (;@2;)
        end
        block  ;; label = @3
          i32.const 1049924
          i32.load
          local.get 2
          i32.ne
          if  ;; label = @4
            local.get 2
            i32.const 1049920
            i32.load
            i32.ne
            br_if 1 (;@3;)
            i32.const 1049920
            local.get 0
            i32.store
            i32.const 1049912
            i32.const 1049912
            i32.load
            local.get 1
            i32.add
            local.tee 1
            i32.store
            local.get 0
            local.get 1
            call 48
            return
          end
          i32.const 1049924
          local.get 0
          i32.store
          i32.const 1049916
          i32.const 1049916
          i32.load
          local.get 1
          i32.add
          local.tee 1
          i32.store
          local.get 0
          local.get 1
          i32.const 1
          i32.or
          i32.store offset=4
          local.get 0
          i32.const 1049920
          i32.load
          i32.ne
          br_if 2 (;@1;)
          i32.const 1049912
          i32.const 0
          i32.store
          i32.const 1049920
          i32.const 0
          i32.store
          return
        end
        local.get 2
        call 67
        local.tee 3
        local.get 1
        i32.add
        local.set 1
        block  ;; label = @3
          local.get 3
          i32.const 256
          i32.ge_u
          if  ;; label = @4
            local.get 2
            call 14
            br 1 (;@3;)
          end
          local.get 2
          i32.const 12
          i32.add
          i32.load
          local.tee 4
          local.get 2
          i32.const 8
          i32.add
          i32.load
          local.tee 2
          i32.ne
          if  ;; label = @4
            local.get 2
            local.get 4
            i32.store offset=12
            local.get 4
            local.get 2
            i32.store offset=8
            br 1 (;@3;)
          end
          i32.const 1049512
          i32.const 1049512
          i32.load
          i32.const -2
          local.get 3
          i32.const 3
          i32.shr_u
          i32.rotl
          i32.and
          i32.store
        end
        local.get 0
        local.get 1
        call 48
        local.get 0
        i32.const 1049920
        i32.load
        i32.ne
        br_if 0 (;@2;)
        i32.const 1049912
        local.get 1
        i32.store
        return
      end
      local.get 1
      i32.const 256
      i32.ge_u
      if  ;; label = @2
        local.get 0
        local.get 1
        call 13
        return
      end
      local.get 1
      i32.const 3
      i32.shr_u
      local.tee 2
      i32.const 3
      i32.shl
      i32.const 1049520
      i32.add
      local.set 1
      block (result i32)  ;; label = @2
        i32.const 1049512
        i32.load
        local.tee 3
        i32.const 1
        local.get 2
        i32.shl
        local.tee 2
        i32.and
        if  ;; label = @3
          local.get 1
          i32.load offset=8
          br 1 (;@2;)
        end
        i32.const 1049512
        local.get 2
        local.get 3
        i32.or
        i32.store
        local.get 1
      end
      local.set 2
      local.get 1
      local.get 0
      i32.store offset=8
      local.get 2
      local.get 0
      i32.store offset=12
      local.get 0
      local.get 1
      i32.store offset=12
      local.get 0
      local.get 2
      i32.store offset=8
    end)
  (func (;9;) (type 1) (param i32 i32) (result i32)
    (local i32 i32 i32 i32)
    global.get 0
    i32.const 32
    i32.sub
    local.tee 2
    global.set 0
    local.get 0
    i32.load
    local.set 4
    block  ;; label = @1
      block  ;; label = @2
        block  ;; label = @3
          block  ;; label = @4
            block  ;; label = @5
              local.get 1
              i32.const 128
              i32.ge_u
              if  ;; label = @6
                local.get 2
                i32.const 0
                i32.store offset=16
                local.get 1
                i32.const 2048
                i32.lt_u
                br_if 1 (;@5;)
                local.get 2
                i32.const 16
                i32.add
                local.set 0
                local.get 1
                i32.const 65536
                i32.lt_u
                if  ;; label = @7
                  local.get 2
                  local.get 1
                  i32.const 63
                  i32.and
                  i32.const 128
                  i32.or
                  i32.store8 offset=18
                  local.get 2
                  local.get 1
                  i32.const 12
                  i32.shr_u
                  i32.const 224
                  i32.or
                  i32.store8 offset=16
                  local.get 2
                  local.get 1
                  i32.const 6
                  i32.shr_u
                  i32.const 63
                  i32.and
                  i32.const 128
                  i32.or
                  i32.store8 offset=17
                  i32.const 3
                  local.set 1
                  br 5 (;@2;)
                end
                local.get 2
                local.get 1
                i32.const 63
                i32.and
                i32.const 128
                i32.or
                i32.store8 offset=19
                local.get 2
                local.get 1
                i32.const 18
                i32.shr_u
                i32.const 240
                i32.or
                i32.store8 offset=16
                local.get 2
                local.get 1
                i32.const 6
                i32.shr_u
                i32.const 63
                i32.and
                i32.const 128
                i32.or
                i32.store8 offset=18
                local.get 2
                local.get 1
                i32.const 12
                i32.shr_u
                i32.const 63
                i32.and
                i32.const 128
                i32.or
                i32.store8 offset=17
                i32.const 4
                local.set 1
                br 4 (;@2;)
              end
              local.get 4
              i32.load offset=8
              local.tee 0
              local.get 4
              i32.const 4
              i32.add
              i32.load
              i32.ne
              if  ;; label = @6
                local.get 4
                i32.load
                local.set 3
                br 3 (;@3;)
              end
              local.get 0
              i32.const 1
              i32.add
              local.tee 3
              local.get 0
              i32.lt_u
              br_if 1 (;@4;)
              local.get 0
              i32.const 1
              i32.shl
              local.tee 5
              local.get 3
              local.get 5
              local.get 3
              i32.gt_u
              select
              local.tee 3
              i32.const 8
              local.get 3
              i32.const 8
              i32.gt_u
              select
              local.set 3
              block  ;; label = @6
                local.get 0
                if  ;; label = @7
                  local.get 2
                  i32.const 24
                  i32.add
                  i32.const 1
                  i32.store
                  local.get 2
                  local.get 0
                  i32.store offset=20
                  local.get 2
                  local.get 4
                  i32.load
                  i32.store offset=16
                  br 1 (;@6;)
                end
                local.get 2
                i32.const 0
                i32.store offset=16
              end
              local.get 2
              local.get 3
              local.get 2
              i32.const 16
              i32.add
              call 21
              local.get 2
              i32.const 8
              i32.add
              i32.load
              local.set 0
              local.get 2
              i32.load offset=4
              local.set 3
              local.get 2
              i32.load
              i32.const 1
              i32.ne
              if  ;; label = @6
                local.get 4
                local.get 3
                i32.store
                local.get 4
                i32.const 4
                i32.add
                local.get 0
                i32.store
                local.get 4
                i32.load offset=8
                local.set 0
                br 3 (;@3;)
              end
              local.get 0
              i32.eqz
              br_if 1 (;@4;)
              local.get 3
              local.get 0
              call 71
              unreachable
            end
            local.get 2
            local.get 1
            i32.const 63
            i32.and
            i32.const 128
            i32.or
            i32.store8 offset=17
            local.get 2
            local.get 1
            i32.const 6
            i32.shr_u
            i32.const 192
            i32.or
            i32.store8 offset=16
            local.get 2
            i32.const 16
            i32.add
            local.set 0
            i32.const 2
            local.set 1
            br 2 (;@2;)
          end
          call 64
          unreachable
        end
        local.get 0
        local.get 3
        i32.add
        local.get 1
        i32.store8
        local.get 4
        local.get 4
        i32.load offset=8
        i32.const 1
        i32.add
        i32.store offset=8
        br 1 (;@1;)
      end
      local.get 4
      local.get 0
      local.get 0
      local.get 1
      i32.add
      call 17
    end
    local.get 2
    i32.const 32
    i32.add
    global.set 0
    i32.const 0)
  (func (;10;) (type 1) (param i32 i32) (result i32)
    (local i32 i32 i32)
    block  ;; label = @1
      local.get 1
      i32.const 9
      i32.ge_u
      if  ;; label = @2
        i32.const 16
        i32.const 8
        call 50
        local.get 1
        i32.gt_u
        if  ;; label = @3
          i32.const 16
          i32.const 8
          call 50
          local.set 1
        end
        i32.const 0
        call 74
        local.tee 3
        local.get 3
        i32.const 8
        call 50
        i32.sub
        i32.const 20
        i32.const 8
        call 50
        i32.sub
        i32.const 16
        i32.const 8
        call 50
        i32.sub
        i32.const -65544
        i32.add
        i32.const -9
        i32.and
        i32.const -3
        i32.add
        local.tee 3
        i32.const 0
        i32.const 16
        i32.const 8
        call 50
        i32.const 2
        i32.shl
        i32.sub
        local.tee 2
        local.get 3
        local.get 2
        i32.lt_u
        select
        local.get 1
        i32.sub
        local.get 0
        i32.le_u
        br_if 1 (;@1;)
        local.get 1
        i32.const 16
        local.get 0
        i32.const 4
        i32.add
        i32.const 16
        i32.const 8
        call 50
        i32.const -5
        i32.add
        local.get 0
        i32.gt_u
        select
        i32.const 8
        call 50
        local.tee 3
        i32.add
        i32.const 16
        i32.const 8
        call 50
        i32.add
        i32.const -4
        i32.add
        call 3
        local.tee 2
        i32.eqz
        br_if 1 (;@1;)
        local.get 2
        call 75
        local.set 0
        block  ;; label = @3
          local.get 1
          i32.const -1
          i32.add
          local.tee 4
          local.get 2
          i32.and
          i32.eqz
          if  ;; label = @4
            local.get 0
            local.set 1
            br 1 (;@3;)
          end
          local.get 2
          local.get 4
          i32.add
          i32.const 0
          local.get 1
          i32.sub
          i32.and
          call 75
          local.set 2
          i32.const 16
          i32.const 8
          call 50
          local.set 4
          local.get 0
          call 67
          local.get 2
          local.get 1
          local.get 2
          i32.add
          local.get 2
          local.get 0
          i32.sub
          local.get 4
          i32.gt_u
          select
          local.tee 1
          local.get 0
          i32.sub
          local.tee 2
          i32.sub
          local.set 4
          local.get 0
          call 61
          i32.eqz
          if  ;; label = @4
            local.get 1
            local.get 4
            call 40
            local.get 0
            local.get 2
            call 40
            local.get 0
            local.get 2
            call 8
            br 1 (;@3;)
          end
          local.get 0
          i32.load
          local.set 0
          local.get 1
          local.get 4
          i32.store offset=4
          local.get 1
          local.get 0
          local.get 2
          i32.add
          i32.store
        end
        block  ;; label = @3
          local.get 1
          call 61
          br_if 0 (;@3;)
          local.get 1
          call 67
          local.tee 2
          i32.const 16
          i32.const 8
          call 50
          local.get 3
          i32.add
          i32.le_u
          br_if 0 (;@3;)
          local.get 1
          local.get 3
          call 72
          local.set 0
          local.get 1
          local.get 3
          call 40
          local.get 0
          local.get 2
          local.get 3
          i32.sub
          local.tee 3
          call 40
          local.get 0
          local.get 3
          call 8
        end
        local.get 1
        call 74
        local.get 1
        call 61
        drop
        return
      end
      local.get 0
      call 3
      local.set 4
    end
    local.get 4)
  (func (;11;) (type 12) (param i64 i32) (result i32)
    (local i32 i32 i32 i32 i32 i64)
    global.get 0
    i32.const 48
    i32.sub
    local.tee 4
    global.set 0
    i32.const 39
    local.set 2
    block  ;; label = @1
      local.get 0
      i64.const 10000
      i64.lt_u
      if  ;; label = @2
        local.get 0
        local.set 7
        br 1 (;@1;)
      end
      loop  ;; label = @2
        local.get 4
        i32.const 9
        i32.add
        local.get 2
        i32.add
        local.tee 3
        i32.const -4
        i32.add
        local.get 0
        local.get 0
        i64.const 10000
        i64.div_u
        local.tee 7
        i64.const 10000
        i64.mul
        i64.sub
        i32.wrap_i64
        local.tee 5
        i32.const 65535
        i32.and
        i32.const 100
        i32.div_u
        local.tee 6
        i32.const 1
        i32.shl
        i32.const 1049228
        i32.add
        i32.load16_u align=1
        i32.store16 align=1
        local.get 3
        i32.const -2
        i32.add
        local.get 5
        local.get 6
        i32.const 100
        i32.mul
        i32.sub
        i32.const 65535
        i32.and
        i32.const 1
        i32.shl
        i32.const 1049228
        i32.add
        i32.load16_u align=1
        i32.store16 align=1
        local.get 2
        i32.const -4
        i32.add
        local.set 2
        local.get 0
        i64.const 99999999
        i64.gt_u
        local.get 7
        local.set 0
        br_if 0 (;@2;)
      end
    end
    local.get 7
    i32.wrap_i64
    local.tee 3
    i32.const 99
    i32.gt_s
    if  ;; label = @1
      local.get 2
      i32.const -2
      i32.add
      local.tee 2
      local.get 4
      i32.const 9
      i32.add
      i32.add
      local.get 7
      i32.wrap_i64
      local.tee 3
      local.get 3
      i32.const 65535
      i32.and
      i32.const 100
      i32.div_u
      local.tee 3
      i32.const 100
      i32.mul
      i32.sub
      i32.const 65535
      i32.and
      i32.const 1
      i32.shl
      i32.const 1049228
      i32.add
      i32.load16_u align=1
      i32.store16 align=1
    end
    block  ;; label = @1
      local.get 3
      i32.const 10
      i32.ge_s
      if  ;; label = @2
        local.get 2
        i32.const -2
        i32.add
        local.tee 2
        local.get 4
        i32.const 9
        i32.add
        i32.add
        local.get 3
        i32.const 1
        i32.shl
        i32.const 1049228
        i32.add
        i32.load16_u align=1
        i32.store16 align=1
        br 1 (;@1;)
      end
      local.get 2
      i32.const -1
      i32.add
      local.tee 2
      local.get 4
      i32.const 9
      i32.add
      i32.add
      local.get 3
      i32.const 48
      i32.add
      i32.store8
    end
    local.get 1
    local.get 4
    i32.const 9
    i32.add
    local.get 2
    i32.add
    i32.const 39
    local.get 2
    i32.sub
    call 6
    local.get 4
    i32.const 48
    i32.add
    global.set 0)
  (func (;12;) (type 11) (result i32)
    (local i32 i32 i32 i32 i32 i32 i32 i32 i32 i32 i32 i32 i32 i32)
    i32.const 1049944
    i32.load
    local.tee 3
    i32.eqz
    if  ;; label = @1
      i32.const 1049960
      i32.const 4095
      i32.store
      i32.const 0
      return
    end
    i32.const 1049936
    local.set 2
    loop  ;; label = @1
      local.get 3
      local.tee 0
      i32.load offset=8
      local.set 3
      local.get 0
      i32.load offset=4
      local.set 4
      local.get 0
      i32.load
      local.set 5
      block  ;; label = @2
        local.get 0
        i32.const 12
        i32.add
        i32.load
        drop
        i32.const 1
        if  ;; label = @3
          local.get 0
          local.set 2
          br 1 (;@2;)
        end
        local.get 0
        call 69
        if  ;; label = @3
          local.get 0
          local.set 2
          br 1 (;@2;)
        end
        local.get 5
        local.get 5
        call 74
        local.tee 1
        i32.const 8
        call 50
        local.get 1
        i32.sub
        i32.add
        local.tee 1
        call 67
        local.set 7
        i32.const 0
        call 74
        local.tee 9
        i32.const 8
        call 50
        local.set 10
        i32.const 20
        i32.const 8
        call 50
        local.set 11
        i32.const 16
        i32.const 8
        call 50
        local.set 12
        local.get 1
        call 59
        local.get 1
        local.get 7
        i32.add
        local.get 5
        local.get 4
        local.get 9
        i32.add
        local.get 10
        i32.sub
        local.get 11
        i32.sub
        local.get 12
        i32.sub
        i32.add
        i32.lt_u
        if  ;; label = @3
          local.get 0
          local.set 2
          br 1 (;@2;)
        end
        if  ;; label = @3
          local.get 0
          local.set 2
          br 1 (;@2;)
        end
        block  ;; label = @3
          local.get 1
          i32.const 1049920
          i32.load
          i32.ne
          if  ;; label = @4
            local.get 1
            call 14
            br 1 (;@3;)
          end
          i32.const 1049912
          i32.const 0
          i32.store
          i32.const 1049920
          i32.const 0
          i32.store
        end
        local.get 1
        local.get 7
        call 13
        local.get 0
        local.set 2
        br 0 (;@2;)
      end
      local.get 6
      i32.const 1
      i32.add
      local.set 6
      local.get 3
      br_if 0 (;@1;)
    end
    i32.const 1049960
    local.get 6
    i32.const 4095
    local.get 6
    i32.const 4095
    i32.gt_u
    select
    i32.store
    local.get 8)
  (func (;13;) (type 0) (param i32 i32)
    (local i32 i32 i32 i32 i32)
    local.get 0
    i64.const 0
    i64.store offset=16 align=4
    local.get 0
    block (result i32)  ;; label = @1
      i32.const 0
      local.get 1
      i32.const 8
      i32.shr_u
      local.tee 2
      i32.eqz
      br_if 0 (;@1;)
      drop
      i32.const 31
      local.get 1
      i32.const 16777215
      i32.gt_u
      br_if 0 (;@1;)
      drop
      local.get 1
      i32.const 6
      local.get 2
      i32.clz
      local.tee 2
      i32.sub
      i32.const 31
      i32.and
      i32.shr_u
      i32.const 1
      i32.and
      local.get 2
      i32.const 1
      i32.shl
      i32.sub
      i32.const 62
      i32.add
    end
    local.tee 2
    i32.store offset=28
    local.get 2
    i32.const 2
    i32.shl
    i32.const 1049784
    i32.add
    local.set 3
    local.get 0
    local.set 4
    block  ;; label = @1
      block  ;; label = @2
        block  ;; label = @3
          block  ;; label = @4
            i32.const 1049516
            i32.load
            local.tee 5
            i32.const 1
            local.get 2
            i32.const 31
            i32.and
            i32.shl
            local.tee 6
            i32.and
            if  ;; label = @5
              local.get 3
              i32.load
              local.set 3
              local.get 2
              call 47
              local.set 2
              local.get 3
              call 67
              local.get 1
              i32.ne
              br_if 1 (;@4;)
              local.get 3
              local.set 2
              br 2 (;@3;)
            end
            i32.const 1049516
            local.get 5
            local.get 6
            i32.or
            i32.store
            local.get 3
            local.get 0
            i32.store
            br 3 (;@1;)
          end
          local.get 1
          local.get 2
          i32.const 31
          i32.and
          i32.shl
          local.set 5
          loop  ;; label = @4
            local.get 3
            local.get 5
            i32.const 29
            i32.shr_u
            i32.const 4
            i32.and
            i32.add
            i32.const 16
            i32.add
            local.tee 6
            i32.load
            local.tee 2
            i32.eqz
            br_if 2 (;@2;)
            local.get 5
            i32.const 1
            i32.shl
            local.set 5
            local.get 2
            local.tee 3
            call 67
            local.get 1
            i32.ne
            br_if 0 (;@4;)
          end
        end
        local.get 2
        i32.load offset=8
        local.tee 1
        local.get 4
        i32.store offset=12
        local.get 2
        local.get 4
        i32.store offset=8
        local.get 4
        local.get 2
        i32.store offset=12
        local.get 4
        local.get 1
        i32.store offset=8
        local.get 0
        i32.const 0
        i32.store offset=24
        return
      end
      local.get 6
      local.get 0
      i32.store
    end
    local.get 0
    local.get 3
    i32.store offset=24
    local.get 4
    local.get 4
    i32.store offset=8
    local.get 4
    local.get 4
    i32.store offset=12)
  (func (;14;) (type 2) (param i32)
    (local i32 i32 i32 i32 i32)
    local.get 0
    i32.load offset=24
    local.set 4
    block  ;; label = @1
      block  ;; label = @2
        local.get 0
        local.get 0
        i32.load offset=12
        i32.eq
        if  ;; label = @3
          local.get 0
          i32.const 20
          i32.const 16
          local.get 0
          i32.const 20
          i32.add
          local.tee 1
          i32.load
          local.tee 3
          select
          i32.add
          i32.load
          local.tee 2
          br_if 1 (;@2;)
          i32.const 0
          local.set 1
          br 2 (;@1;)
        end
        local.get 0
        i32.load offset=8
        local.tee 2
        local.get 0
        i32.load offset=12
        local.tee 1
        i32.store offset=12
        local.get 1
        local.get 2
        i32.store offset=8
        br 1 (;@1;)
      end
      local.get 1
      local.get 0
      i32.const 16
      i32.add
      local.get 3
      select
      local.set 3
      loop  ;; label = @2
        local.get 3
        local.set 5
        local.get 2
        local.tee 1
        i32.const 20
        i32.add
        local.tee 3
        i32.load
        local.tee 2
        i32.eqz
        if  ;; label = @3
          local.get 1
          i32.const 16
          i32.add
          local.set 3
          local.get 1
          i32.load offset=16
          local.set 2
        end
        local.get 2
        br_if 0 (;@2;)
      end
      local.get 5
      i32.const 0
      i32.store
    end
    block  ;; label = @1
      local.get 4
      i32.eqz
      br_if 0 (;@1;)
      block  ;; label = @2
        local.get 0
        local.get 0
        i32.load offset=28
        i32.const 2
        i32.shl
        i32.const 1049784
        i32.add
        local.tee 2
        i32.load
        i32.ne
        if  ;; label = @3
          local.get 4
          i32.const 16
          i32.const 20
          local.get 4
          i32.load offset=16
          local.get 0
          i32.eq
          select
          i32.add
          local.get 1
          i32.store
          local.get 1
          i32.eqz
          br_if 2 (;@1;)
          br 1 (;@2;)
        end
        local.get 2
        local.get 1
        i32.store
        local.get 1
        br_if 0 (;@2;)
        i32.const 1049516
        i32.const 1049516
        i32.load
        i32.const -2
        local.get 0
        i32.load offset=28
        i32.rotl
        i32.and
        i32.store
        br 1 (;@1;)
      end
      local.get 1
      local.get 4
      i32.store offset=24
      local.get 0
      i32.load offset=16
      local.tee 2
      if  ;; label = @2
        local.get 1
        local.get 2
        i32.store offset=16
        local.get 2
        local.get 1
        i32.store offset=24
      end
      local.get 0
      i32.const 20
      i32.add
      i32.load
      local.tee 0
      i32.eqz
      br_if 0 (;@1;)
      local.get 1
      i32.const 20
      i32.add
      local.get 0
      i32.store
      local.get 0
      local.get 1
      i32.store offset=24
    end)
  (func (;15;) (type 0) (param i32 i32)
    (local i32 i32 i32 i32)
    global.get 0
    i32.const -64
    i32.add
    local.tee 2
    global.set 0
    local.get 1
    i32.load offset=4
    local.tee 3
    i32.eqz
    if  ;; label = @1
      local.get 1
      i32.const 4
      i32.add
      local.set 3
      local.get 1
      i32.load
      local.set 4
      local.get 2
      i32.const 0
      i32.store offset=32
      local.get 2
      i64.const 1
      i64.store offset=24
      local.get 2
      local.get 2
      i32.const 24
      i32.add
      i32.store offset=36
      local.get 2
      i32.const 56
      i32.add
      local.get 4
      i32.const 16
      i32.add
      i64.load align=4
      i64.store
      local.get 2
      i32.const 48
      i32.add
      local.get 4
      i32.const 8
      i32.add
      i64.load align=4
      i64.store
      local.get 2
      local.get 4
      i64.load align=4
      i64.store offset=40
      local.get 2
      i32.const 36
      i32.add
      local.get 2
      i32.const 40
      i32.add
      call 7
      drop
      local.get 2
      i32.const 16
      i32.add
      local.tee 4
      local.get 2
      i32.load offset=32
      i32.store
      local.get 2
      local.get 2
      i64.load offset=24
      i64.store offset=8
      block  ;; label = @2
        local.get 1
        i32.load offset=4
        local.tee 5
        i32.eqz
        br_if 0 (;@2;)
        local.get 1
        i32.const 8
        i32.add
        i32.load
        i32.eqz
        br_if 0 (;@2;)
        local.get 5
        call 4
      end
      local.get 3
      local.get 2
      i64.load offset=8
      i64.store align=4
      local.get 3
      i32.const 8
      i32.add
      local.get 4
      i32.load
      i32.store
      local.get 3
      i32.load
      local.set 3
    end
    local.get 1
    i32.const 1
    i32.store offset=4
    local.get 1
    i32.const 12
    i32.add
    i32.load
    local.set 4
    local.get 1
    i32.const 8
    i32.add
    local.tee 1
    i32.load
    local.set 5
    local.get 1
    i64.const 0
    i64.store align=4
    i32.const 12
    i32.const 4
    call 54
    local.tee 1
    i32.eqz
    if  ;; label = @1
      i32.const 12
      i32.const 4
      call 71
      unreachable
    end
    local.get 1
    local.get 4
    i32.store offset=8
    local.get 1
    local.get 5
    i32.store offset=4
    local.get 1
    local.get 3
    i32.store
    local.get 0
    i32.const 1049048
    i32.store offset=4
    local.get 0
    local.get 1
    i32.store
    local.get 2
    i32.const -64
    i32.sub
    global.set 0)
  (func (;16;) (type 5) (param i32 i32 i32)
    (local i32 i32 i32 i32)
    global.get 0
    i32.const 48
    i32.sub
    local.tee 3
    global.set 0
    block  ;; label = @1
      local.get 2
      i32.const 1
      i32.add
      local.tee 4
      local.get 2
      i32.ge_u
      if  ;; label = @2
        i32.const 4
        local.set 5
        block  ;; label = @3
          local.get 1
          i32.load offset=4
          local.tee 2
          i32.const 1
          i32.shl
          local.tee 6
          local.get 4
          local.get 6
          local.get 4
          i32.gt_u
          select
          local.tee 4
          i32.const 4
          local.get 4
          i32.const 4
          i32.gt_u
          select
          local.tee 4
          i32.const 1073741823
          i32.and
          local.get 4
          i32.eq
          if  ;; label = @4
            local.get 4
            i32.const 2
            i32.shl
            local.set 4
            br 1 (;@3;)
          end
          local.get 1
          i32.load offset=4
          local.set 2
          i32.const 0
          local.set 5
        end
        block  ;; label = @3
          local.get 2
          if  ;; label = @4
            local.get 1
            i32.load
            local.set 6
            local.get 3
            i32.const 40
            i32.add
            i32.const 4
            i32.store
            local.get 3
            local.get 2
            i32.const 2
            i32.shl
            i32.store offset=36
            local.get 3
            local.get 6
            i32.store offset=32
            br 1 (;@3;)
          end
          local.get 3
          i32.const 0
          i32.store offset=32
        end
        local.get 3
        i32.const 16
        i32.add
        local.get 4
        local.get 5
        local.get 3
        i32.const 32
        i32.add
        call 22
        i32.const 1
        local.set 2
        local.get 3
        i32.const 24
        i32.add
        i32.load
        local.set 4
        local.get 3
        i32.load offset=20
        local.set 5
        local.get 3
        i32.load offset=16
        i32.const 1
        i32.ne
        if  ;; label = @3
          local.get 1
          local.get 5
          i32.store
          local.get 1
          local.get 4
          i32.const 2
          i32.shr_u
          i32.store offset=4
          i32.const 0
          local.set 2
          br 2 (;@1;)
        end
        local.get 3
        i32.const 8
        i32.add
        local.get 5
        local.get 4
        call 55
        local.get 0
        local.get 3
        i64.load offset=8
        i64.store offset=4 align=4
        br 1 (;@1;)
      end
      local.get 3
      local.get 4
      i32.const 0
      call 55
      local.get 0
      local.get 3
      i64.load
      i64.store offset=4 align=4
      i32.const 1
      local.set 2
    end
    local.get 0
    local.get 2
    i32.store
    local.get 3
    i32.const 48
    i32.add
    global.set 0)
  (func (;17;) (type 5) (param i32 i32 i32)
    (local i32 i32 i32 i32)
    global.get 0
    i32.const 32
    i32.sub
    local.tee 3
    global.set 0
    block  ;; label = @1
      local.get 0
      i32.const 4
      i32.add
      i32.load
      local.tee 5
      local.get 0
      i32.const 8
      i32.add
      i32.load
      local.tee 4
      i32.sub
      local.get 2
      local.get 1
      i32.sub
      local.tee 6
      i32.ge_u
      if  ;; label = @2
        local.get 0
        i32.load
        local.set 2
        br 1 (;@1;)
      end
      block  ;; label = @2
        local.get 4
        local.get 6
        i32.add
        local.tee 2
        local.get 4
        i32.lt_u
        br_if 0 (;@2;)
        local.get 5
        i32.const 1
        i32.shl
        local.tee 4
        local.get 2
        local.get 4
        local.get 2
        i32.gt_u
        select
        local.tee 2
        i32.const 8
        local.get 2
        i32.const 8
        i32.gt_u
        select
        local.set 2
        block  ;; label = @3
          local.get 5
          if  ;; label = @4
            local.get 3
            i32.const 24
            i32.add
            i32.const 1
            i32.store
            local.get 3
            local.get 5
            i32.store offset=20
            local.get 3
            local.get 0
            i32.load
            i32.store offset=16
            br 1 (;@3;)
          end
          local.get 3
          i32.const 0
          i32.store offset=16
        end
        local.get 3
        local.get 2
        local.get 3
        i32.const 16
        i32.add
        call 21
        local.get 3
        i32.const 8
        i32.add
        i32.load
        local.set 4
        local.get 3
        i32.load offset=4
        local.set 2
        local.get 3
        i32.load
        i32.const 1
        i32.ne
        if  ;; label = @3
          local.get 0
          local.get 2
          i32.store
          local.get 0
          i32.const 4
          i32.add
          local.get 4
          i32.store
          local.get 0
          i32.const 8
          i32.add
          i32.load
          local.set 4
          br 2 (;@1;)
        end
        local.get 4
        i32.eqz
        br_if 0 (;@2;)
        local.get 2
        local.get 4
        call 71
        unreachable
      end
      call 64
      unreachable
    end
    local.get 2
    local.get 4
    i32.add
    local.get 1
    local.get 6
    call 36
    drop
    local.get 0
    i32.const 8
    i32.add
    local.tee 0
    local.get 0
    i32.load
    local.get 6
    i32.add
    i32.store
    local.get 3
    i32.const 32
    i32.add
    global.set 0)
  (func (;18;) (type 0) (param i32 i32)
    (local i32 i32 i32 i32)
    global.get 0
    i32.const -64
    i32.add
    local.tee 2
    global.set 0
    local.get 1
    i32.const 4
    i32.add
    local.set 4
    local.get 1
    i32.load offset=4
    i32.eqz
    if  ;; label = @1
      local.get 1
      i32.load
      local.set 3
      local.get 2
      i32.const 0
      i32.store offset=32
      local.get 2
      i64.const 1
      i64.store offset=24
      local.get 2
      local.get 2
      i32.const 24
      i32.add
      i32.store offset=36
      local.get 2
      i32.const 56
      i32.add
      local.get 3
      i32.const 16
      i32.add
      i64.load align=4
      i64.store
      local.get 2
      i32.const 48
      i32.add
      local.get 3
      i32.const 8
      i32.add
      i64.load align=4
      i64.store
      local.get 2
      local.get 3
      i64.load align=4
      i64.store offset=40
      local.get 2
      i32.const 36
      i32.add
      local.get 2
      i32.const 40
      i32.add
      call 7
      drop
      local.get 2
      i32.const 16
      i32.add
      local.tee 3
      local.get 2
      i32.load offset=32
      i32.store
      local.get 2
      local.get 2
      i64.load offset=24
      i64.store offset=8
      block  ;; label = @2
        local.get 1
        i32.load offset=4
        local.tee 5
        i32.eqz
        br_if 0 (;@2;)
        local.get 1
        i32.const 8
        i32.add
        i32.load
        i32.eqz
        br_if 0 (;@2;)
        local.get 5
        call 4
      end
      local.get 4
      local.get 2
      i64.load offset=8
      i64.store align=4
      local.get 4
      i32.const 8
      i32.add
      local.get 3
      i32.load
      i32.store
    end
    local.get 0
    i32.const 1049048
    i32.store offset=4
    local.get 0
    local.get 4
    i32.store
    local.get 2
    i32.const -64
    i32.sub
    global.set 0)
  (func (;19;) (type 2) (param i32)
    (local i32 i32)
    global.get 0
    i32.const 16
    i32.sub
    local.tee 2
    global.set 0
    i32.const 1
    i32.const -36
    i32.const 11
    call 1
    i32.const 2
    i32.ne
    if  ;; label = @1
      i32.const 1
      i32.const -29
      i32.const 13
      i32.const 1048724
      i32.const 18
      call 0
    end
    i32.const 1
    i32.const -30
    i32.const 11
    call 1
    i32.const 3
    i32.ne
    if  ;; label = @1
      i32.const 1
      i32.const -29
      i32.const 13
      i32.const 1048724
      i32.const 18
      call 0
    end
    i32.const 1
    i32.const -34
    i32.const 11
    call 1
    i32.const 4
    i32.ne
    if  ;; label = @1
      i32.const 1
      i32.const -29
      i32.const 13
      i32.const 1048724
      i32.const 18
      call 0
    end
    block  ;; label = @1
      block  ;; label = @2
        block  ;; label = @3
          local.get 0
          i32.const 32768
          i32.and
          i32.eqz
          if  ;; label = @4
            i32.const 1049436
            i32.load
            local.tee 1
            local.get 0
            i32.le_u
            br_if 2 (;@2;)
            local.get 2
            i32.const 8
            i32.add
            i32.const 1049428
            i32.load
            local.get 0
            i32.const 2
            i32.shl
            i32.add
            i32.load
            call_indirect (type 2)
            br 1 (;@3;)
          end
          i32.const 1049448
          i32.load
          local.tee 1
          local.get 0
          i32.const 32767
          i32.and
          local.tee 0
          i32.le_u
          br_if 2 (;@1;)
          i32.const 1048724
          i32.const 1049440
          i32.load
          local.get 0
          i32.const 2
          i32.shl
          i32.add
          i32.load
          call_indirect (type 2)
        end
        local.get 2
        i32.const 16
        i32.add
        global.set 0
        return
      end
      local.get 0
      local.get 1
      i32.const 1048832
      call 24
      unreachable
    end
    local.get 0
    local.get 1
    i32.const 1048816
    call 24
    unreachable)
  (func (;20;) (type 8) (param i32 i32 i32 i32)
    (local i32 i32)
    global.get 0
    i32.const 32
    i32.sub
    local.tee 4
    global.set 0
    i32.const 1
    local.set 5
    i32.const 1049508
    i32.const 1049508
    i32.load
    i32.const 1
    i32.add
    i32.store
    block  ;; label = @1
      block  ;; label = @2
        block  ;; label = @3
          i32.const 1049968
          i32.load
          i32.const 1
          i32.ne
          if  ;; label = @4
            i32.const 1049968
            i64.const 4294967297
            i64.store
            br 1 (;@3;)
          end
          i32.const 1049972
          i32.const 1049972
          i32.load
          i32.const 1
          i32.add
          local.tee 5
          i32.store
          local.get 5
          i32.const 2
          i32.gt_u
          br_if 1 (;@2;)
        end
        local.get 4
        local.get 3
        i32.store offset=28
        local.get 4
        local.get 2
        i32.store offset=24
        local.get 4
        i32.const 1048888
        i32.store offset=20
        local.get 4
        i32.const 1048888
        i32.store offset=16
        i32.const 1049496
        i32.load
        local.tee 2
        i32.const -1
        i32.le_s
        br_if 0 (;@2;)
        i32.const 1049496
        local.get 2
        i32.const 1
        i32.add
        local.tee 2
        i32.store
        i32.const 1049496
        i32.const 1049504
        i32.load
        local.tee 3
        if (result i32)  ;; label = @3
          i32.const 1049500
          i32.load
          local.get 4
          i32.const 8
          i32.add
          local.get 0
          local.get 1
          i32.load offset=16
          call_indirect (type 0)
          local.get 4
          local.get 4
          i64.load offset=8
          i64.store offset=16
          local.get 4
          i32.const 16
          i32.add
          local.get 3
          i32.load offset=12
          call_indirect (type 0)
          i32.const 1049496
          i32.load
        else
          local.get 2
        end
        i32.const -1
        i32.add
        i32.store
        local.get 5
        i32.const 1
        i32.le_u
        br_if 1 (;@1;)
      end
      unreachable
    end
    global.get 0
    i32.const 16
    i32.sub
    local.tee 2
    global.set 0
    local.get 2
    local.get 1
    i32.store offset=12
    local.get 2
    local.get 0
    i32.store offset=8
    unreachable)
  (func (;21;) (type 5) (param i32 i32 i32)
    (local i32 i32 i32)
    i32.const 1
    local.set 3
    i32.const 1
    local.set 4
    block  ;; label = @1
      local.get 1
      i32.const 0
      i32.lt_s
      if  ;; label = @2
        i32.const 0
        local.set 3
        br 1 (;@1;)
      end
      block (result i32)  ;; label = @2
        local.get 2
        i32.load
        local.tee 5
        i32.eqz
        if  ;; label = @3
          i32.const 1
          local.get 1
          i32.eqz
          br_if 1 (;@2;)
          drop
          local.get 1
          i32.const 1
          call 54
          br 1 (;@2;)
        end
        local.get 2
        i32.load offset=4
        local.tee 2
        i32.eqz
        if  ;; label = @3
          i32.const 1
          local.get 1
          i32.eqz
          br_if 1 (;@2;)
          drop
          local.get 1
          i32.const 1
          call 54
          br 1 (;@2;)
        end
        local.get 5
        local.get 2
        i32.const 1
        local.get 1
        call 51
      end
      local.tee 2
      i32.eqz
      if  ;; label = @2
        local.get 0
        local.get 1
        i32.store offset=4
        br 1 (;@1;)
      end
      local.get 0
      local.get 2
      i32.store offset=4
      i32.const 0
      local.set 4
      local.get 1
      local.set 3
    end
    local.get 0
    local.get 4
    i32.store
    local.get 0
    i32.const 8
    i32.add
    local.get 3
    i32.store)
  (func (;22;) (type 8) (param i32 i32 i32 i32)
    (local i32 i32)
    block  ;; label = @1
      block  ;; label = @2
        block  ;; label = @3
          local.get 2
          if  ;; label = @4
            i32.const 1
            local.set 4
            local.get 1
            i32.const 0
            i32.ge_s
            br_if 1 (;@3;)
            br 2 (;@2;)
          end
          local.get 0
          local.get 1
          i32.store offset=4
          i32.const 1
          local.set 4
          br 1 (;@2;)
        end
        block  ;; label = @3
          block  ;; label = @4
            block  ;; label = @5
              block  ;; label = @6
                block  ;; label = @7
                  local.get 3
                  i32.load
                  local.tee 5
                  i32.eqz
                  if  ;; label = @8
                    local.get 1
                    i32.eqz
                    br_if 1 (;@7;)
                    br 3 (;@5;)
                  end
                  local.get 3
                  i32.load offset=4
                  local.tee 3
                  br_if 1 (;@6;)
                  local.get 1
                  br_if 2 (;@5;)
                end
                local.get 2
                local.set 3
                br 3 (;@3;)
              end
              local.get 5
              local.get 3
              local.get 2
              local.get 1
              call 51
              local.tee 3
              i32.eqz
              br_if 1 (;@4;)
              br 2 (;@3;)
            end
            local.get 1
            local.get 2
            call 54
            local.tee 3
            br_if 1 (;@3;)
          end
          local.get 0
          local.get 1
          i32.store offset=4
          local.get 2
          local.set 1
          br 2 (;@1;)
        end
        local.get 0
        local.get 3
        i32.store offset=4
        i32.const 0
        local.set 4
        br 1 (;@1;)
      end
      i32.const 0
      local.set 1
    end
    local.get 0
    local.get 4
    i32.store
    local.get 0
    i32.const 8
    i32.add
    local.get 1
    i32.store)
  (func (;23;) (type 2) (param i32)
    (local i32 i32 i32)
    global.get 0
    i32.const 16
    i32.sub
    local.tee 1
    global.set 0
    local.get 0
    i32.load
    local.tee 2
    i32.const 20
    i32.add
    i32.load
    local.set 3
    block  ;; label = @1
      block (result i32)  ;; label = @2
        block  ;; label = @3
          block  ;; label = @4
            local.get 2
            i32.load offset=4
            br_table 0 (;@4;) 1 (;@3;) 3 (;@1;)
          end
          local.get 3
          br_if 2 (;@1;)
          i32.const 0
          local.set 2
          i32.const 1048888
          br 1 (;@2;)
        end
        local.get 3
        br_if 1 (;@1;)
        local.get 2
        i32.load
        local.tee 3
        i32.load offset=4
        local.set 2
        local.get 3
        i32.load
      end
      local.set 3
      local.get 1
      local.get 2
      i32.store offset=4
      local.get 1
      local.get 3
      i32.store
      local.get 1
      i32.const 1049028
      local.get 0
      i32.load offset=4
      i32.load offset=8
      local.get 0
      i32.load offset=8
      call 20
      unreachable
    end
    local.get 1
    i32.const 0
    i32.store offset=4
    local.get 1
    local.get 2
    i32.store
    local.get 1
    i32.const 1049008
    local.get 0
    i32.load offset=4
    i32.load offset=8
    local.get 0
    i32.load offset=8
    call 20
    unreachable)
  (func (;24;) (type 5) (param i32 i32 i32)
    (local i32)
    global.get 0
    i32.const 48
    i32.sub
    local.tee 3
    global.set 0
    local.get 3
    local.get 1
    i32.store offset=4
    local.get 3
    local.get 0
    i32.store
    local.get 3
    i32.const 28
    i32.add
    i32.const 2
    i32.store
    local.get 3
    i32.const 44
    i32.add
    i32.const 17
    i32.store
    local.get 3
    i64.const 2
    i64.store offset=12 align=4
    local.get 3
    i32.const 1049212
    i32.store offset=8
    local.get 3
    i32.const 17
    i32.store offset=36
    local.get 3
    local.get 3
    i32.const 32
    i32.add
    i32.store offset=24
    local.get 3
    local.get 3
    i32.store offset=40
    local.get 3
    local.get 3
    i32.const 4
    i32.add
    i32.store offset=32
    local.get 3
    i32.const 8
    i32.add
    local.get 2
    call 39
    unreachable)
  (func (;25;) (type 0) (param i32 i32)
    (local i32)
    global.get 0
    i32.const 16
    i32.sub
    local.tee 2
    global.set 0
    block  ;; label = @1
      local.get 0
      i32.load offset=4
      local.get 1
      i32.sub
      i32.const 1
      i32.ge_u
      if  ;; label = @2
        local.get 2
        i32.const 0
        i32.store
        br 1 (;@1;)
      end
      local.get 2
      local.get 0
      local.get 1
      call 16
      local.get 2
      i32.load
      i32.const 1
      i32.ne
      br_if 0 (;@1;)
      local.get 2
      i32.const 8
      i32.add
      i32.load
      local.tee 0
      if  ;; label = @2
        local.get 2
        i32.load offset=4
        local.get 0
        call 71
        unreachable
      end
      call 64
      unreachable
    end
    local.get 2
    i32.const 16
    i32.add
    global.set 0)
  (func (;26;) (type 1) (param i32 i32) (result i32)
    (local i32)
    global.get 0
    i32.const 32
    i32.sub
    local.tee 2
    global.set 0
    local.get 2
    local.get 0
    i32.load
    i32.store offset=4
    local.get 2
    i32.const 24
    i32.add
    local.get 1
    i32.const 16
    i32.add
    i64.load align=4
    i64.store
    local.get 2
    i32.const 16
    i32.add
    local.get 1
    i32.const 8
    i32.add
    i64.load align=4
    i64.store
    local.get 2
    local.get 1
    i64.load align=4
    i64.store offset=8
    local.get 2
    i32.const 4
    i32.add
    local.get 2
    i32.const 8
    i32.add
    call 7
    local.get 2
    i32.const 32
    i32.add
    global.set 0)
  (func (;27;) (type 2) (param i32)
    (local i32 i32)
    i32.const 1049448
    i32.load
    local.tee 1
    local.set 2
    i32.const 1049444
    i32.load
    local.get 1
    i32.eq
    if  ;; label = @1
      i32.const 1049440
      local.get 1
      call 25
      i32.const 1049448
      i32.load
      local.set 2
    end
    i32.const 1049440
    i32.load
    local.get 2
    i32.const 2
    i32.shl
    i32.add
    i32.const 2
    i32.store
    i32.const 1049448
    i32.const 1049448
    i32.load
    i32.const 1
    i32.add
    i32.store
    local.get 0
    i32.load
    local.get 1
    i32.const 32768
    i32.or
    i32.const 13
    i32.const 1048586
    i32.const 13
    call 0)
  (func (;28;) (type 2) (param i32)
    (local i32 i32)
    i32.const 1049436
    i32.load
    local.tee 1
    local.set 2
    i32.const 1049432
    i32.load
    local.get 1
    i32.eq
    if  ;; label = @1
      i32.const 1049428
      local.get 1
      call 25
      i32.const 1049436
      i32.load
      local.set 2
    end
    i32.const 1049428
    i32.load
    local.get 2
    i32.const 2
    i32.shl
    i32.add
    i32.const 1
    i32.store
    i32.const 1049436
    i32.const 1049436
    i32.load
    i32.const 1
    i32.add
    i32.store
    local.get 0
    i32.load
    local.get 1
    i32.const 13
    i32.const 1048576
    i32.const 10
    call 0)
  (func (;29;) (type 4) (param i32 i32 i32) (result i32)
    block (result i32)  ;; label = @1
      local.get 1
      i32.const 1114112
      i32.ne
      if  ;; label = @2
        i32.const 1
        local.get 0
        i32.load offset=24
        local.get 1
        local.get 0
        i32.const 28
        i32.add
        i32.load
        i32.load offset=16
        call_indirect (type 1)
        br_if 1 (;@1;)
        drop
      end
      local.get 2
      i32.eqz
      if  ;; label = @2
        i32.const 0
        return
      end
      local.get 0
      i32.load offset=24
      local.get 2
      i32.const 0
      local.get 0
      i32.const 28
      i32.add
      i32.load
      i32.load offset=12
      call_indirect (type 4)
    end)
  (func (;30;) (type 7)
    (local i32 i32)
    global.get 0
    i32.const 16
    i32.sub
    local.tee 0
    global.set 0
    i32.const 1
    i32.const -20
    i32.const 45
    call 1
    local.tee 1
    i32.const -41
    i32.const 13
    i32.const 1048848
    i32.const 16
    call 0
    local.get 0
    local.get 1
    i32.store offset=12
    local.get 0
    i32.const 12
    i32.add
    call 28
    local.get 0
    i32.const 12
    i32.add
    call 27
    i32.const 1049456
    i32.const 1048576
    i32.const 10
    call 2
    i32.store
    local.get 0
    i32.const 16
    i32.add
    global.set 0)
  (func (;31;) (type 5) (param i32 i32 i32)
    (local i32)
    global.get 0
    i32.const 32
    i32.sub
    local.tee 3
    global.set 0
    local.get 3
    i32.const 20
    i32.add
    i32.const 0
    i32.store
    local.get 3
    i32.const 1049144
    i32.store offset=16
    local.get 3
    i64.const 1
    i64.store offset=4 align=4
    local.get 3
    local.get 1
    i32.store offset=28
    local.get 3
    local.get 0
    i32.store offset=24
    local.get 3
    local.get 3
    i32.const 24
    i32.add
    i32.store
    local.get 3
    local.get 2
    call 39
    unreachable)
  (func (;32;) (type 0) (param i32 i32)
    (local i32 i32)
    local.get 1
    i32.load offset=4
    local.set 2
    local.get 1
    i32.load
    local.set 3
    i32.const 8
    i32.const 4
    call 54
    local.tee 1
    i32.eqz
    if  ;; label = @1
      i32.const 8
      i32.const 4
      call 71
      unreachable
    end
    local.get 1
    local.get 2
    i32.store offset=4
    local.get 1
    local.get 3
    i32.store
    local.get 0
    i32.const 1049064
    i32.store offset=4
    local.get 0
    local.get 1
    i32.store)
  (func (;33;) (type 0) (param i32 i32)
    (local i32)
    local.get 1
    i32.const 16
    i32.shr_u
    memory.grow
    local.set 2
    local.get 0
    i32.const 0
    i32.store offset=8
    local.get 0
    i32.const 0
    local.get 1
    i32.const -65536
    i32.and
    local.get 2
    i32.const -1
    i32.eq
    local.tee 1
    select
    i32.store offset=4
    local.get 0
    i32.const 0
    local.get 2
    i32.const 16
    i32.shl
    local.get 1
    select
    i32.store)
  (func (;34;) (type 2) (param i32)
    (local i32 i32 i32)
    global.get 0
    i32.const 16
    i32.sub
    local.tee 1
    global.set 0
    local.get 0
    i32.load offset=12
    local.tee 2
    i32.eqz
    if  ;; label = @1
      i32.const 1048904
      i32.const 43
      i32.const 1048976
      call 31
      unreachable
    end
    local.get 0
    i32.load offset=8
    local.tee 3
    i32.eqz
    if  ;; label = @1
      i32.const 1048904
      i32.const 43
      i32.const 1048992
      call 31
      unreachable
    end
    local.get 1
    local.get 2
    i32.store offset=8
    local.get 1
    local.get 0
    i32.store offset=4
    local.get 1
    local.get 3
    i32.store
    local.get 1
    call 38
    unreachable)
  (func (;35;) (type 2) (param i32)
    (local i32)
    global.get 0
    i32.const 16
    i32.sub
    local.tee 1
    global.set 0
    local.get 1
    local.get 0
    i32.load
    i32.const 1049456
    i32.load
    call 55
    local.get 1
    local.get 1
    i64.load
    i64.store offset=8
    local.get 1
    i32.const 8
    i32.add
    local.tee 0
    i32.load
    local.get 0
    i32.load offset=4
    i32.const 13
    i32.const 1048711
    i32.const 13
    call 0
    local.get 1
    i32.const 16
    i32.add
    global.set 0)
  (func (;36;) (type 4) (param i32 i32 i32) (result i32)
    (local i32)
    local.get 2
    if  ;; label = @1
      local.get 0
      local.set 3
      loop  ;; label = @2
        local.get 3
        local.get 1
        i32.load8_u
        i32.store8
        local.get 1
        i32.const 1
        i32.add
        local.set 1
        local.get 3
        i32.const 1
        i32.add
        local.set 3
        local.get 2
        i32.const -1
        i32.add
        local.tee 2
        br_if 0 (;@2;)
      end
    end
    local.get 0)
  (func (;37;) (type 2) (param i32)
    (local i32)
    global.get 0
    i32.const 16
    i32.sub
    local.tee 1
    global.set 0
    i32.const 1048652
    i32.const 28
    call 56
    local.get 1
    i64.const 8589934596
    i64.store offset=8
    local.get 1
    i32.const 8
    i32.add
    call 35
    i32.const 1048680
    i32.const 31
    call 56
    local.get 1
    i32.const 16
    i32.add
    global.set 0)
  (func (;38;) (type 2) (param i32)
    (local i32)
    global.get 0
    i32.const 16
    i32.sub
    local.tee 1
    global.set 0
    local.get 1
    i32.const 8
    i32.add
    local.get 0
    i32.const 8
    i32.add
    i32.load
    i32.store
    local.get 1
    local.get 0
    i64.load align=4
    i64.store
    local.get 1
    call 23
    unreachable)
  (func (;39;) (type 0) (param i32 i32)
    (local i32)
    global.get 0
    i32.const 16
    i32.sub
    local.tee 2
    global.set 0
    local.get 2
    local.get 1
    i32.store offset=12
    local.get 2
    local.get 0
    i32.store offset=8
    local.get 2
    i32.const 1049144
    i32.store offset=4
    local.get 2
    i32.const 1049144
    i32.store
    local.get 2
    call 34
    unreachable)
  (func (;40;) (type 0) (param i32 i32)
    local.get 0
    local.get 0
    i32.load offset=4
    i32.const 1
    i32.and
    local.get 1
    i32.or
    i32.const 2
    i32.or
    i32.store offset=4
    local.get 0
    local.get 1
    i32.add
    local.tee 0
    local.get 0
    i32.load offset=4
    i32.const 1
    i32.or
    i32.store offset=4)
  (func (;41;) (type 2) (param i32)
    (local i32)
    block  ;; label = @1
      local.get 0
      i32.load
      local.tee 1
      i32.eqz
      br_if 0 (;@1;)
      local.get 0
      i32.const 4
      i32.add
      i32.load
      i32.eqz
      br_if 0 (;@1;)
      local.get 1
      call 4
    end)
  (func (;42;) (type 2) (param i32)
    (local i32)
    block  ;; label = @1
      local.get 0
      i32.load offset=4
      local.tee 1
      i32.eqz
      br_if 0 (;@1;)
      local.get 0
      i32.const 8
      i32.add
      i32.load
      i32.eqz
      br_if 0 (;@1;)
      local.get 1
      call 4
    end)
  (func (;43;) (type 1) (param i32 i32) (result i32)
    (local i32)
    block  ;; label = @1
      local.get 0
      i32.load
      local.tee 2
      local.get 1
      i32.gt_u
      br_if 0 (;@1;)
      local.get 2
      local.get 0
      i32.load offset=4
      i32.add
      local.get 1
      i32.le_u
      br_if 0 (;@1;)
      i32.const 1
      return
    end
    i32.const 0)
  (func (;44;) (type 5) (param i32 i32 i32)
    local.get 2
    local.get 2
    i32.load offset=4
    i32.const -2
    i32.and
    i32.store offset=4
    local.get 0
    local.get 1
    i32.const 1
    i32.or
    i32.store offset=4
    local.get 0
    local.get 1
    i32.add
    local.get 1
    i32.store)
  (func (;45;) (type 0) (param i32 i32)
    local.get 0
    local.get 1
    i32.const 3
    i32.or
    i32.store offset=4
    local.get 0
    local.get 1
    i32.add
    local.tee 0
    local.get 0
    i32.load offset=4
    i32.const 1
    i32.or
    i32.store offset=4)
  (func (;46;) (type 3) (param i32) (result i32)
    (local i32)
    local.get 0
    i32.load offset=16
    local.tee 1
    if (result i32)  ;; label = @1
      local.get 1
    else
      local.get 0
      i32.const 20
      i32.add
      i32.load
    end)
  (func (;47;) (type 3) (param i32) (result i32)
    i32.const 0
    i32.const 25
    local.get 0
    i32.const 1
    i32.shr_u
    i32.sub
    local.get 0
    i32.const 31
    i32.eq
    select)
  (func (;48;) (type 0) (param i32 i32)
    local.get 0
    local.get 1
    i32.const 1
    i32.or
    i32.store offset=4
    local.get 0
    local.get 1
    i32.add
    local.get 1
    i32.store)
  (func (;49;) (type 2) (param i32)
    i32.const 1048599
    i32.const 25
    call 56
    i32.const 1048711
    i32.const 13
    call 56
    i32.const 1048624
    i32.const 28
    call 56)
  (func (;50;) (type 1) (param i32 i32) (result i32)
    local.get 0
    local.get 1
    i32.add
    i32.const -1
    i32.add
    i32.const 0
    local.get 1
    i32.sub
    i32.and)
  (func (;51;) (type 9) (param i32 i32 i32 i32) (result i32)
    local.get 0
    local.get 1
    local.get 2
    local.get 3
    call 5)
  (func (;52;) (type 4) (param i32 i32 i32) (result i32)
    local.get 0
    i32.load
    local.get 1
    local.get 1
    local.get 2
    i32.add
    call 17
    i32.const 0)
  (func (;53;) (type 3) (param i32) (result i32)
    local.get 0
    i32.const 1
    i32.shl
    local.tee 0
    i32.const 0
    local.get 0
    i32.sub
    i32.or)
  (func (;54;) (type 1) (param i32 i32) (result i32)
    local.get 0
    local.get 1
    call 10)
  (func (;55;) (type 5) (param i32 i32 i32)
    local.get 0
    local.get 2
    i32.store offset=4
    local.get 0
    local.get 1
    i32.store)
  (func (;56;) (type 0) (param i32 i32)
    i32.const 1
    i32.const -26
    i32.const 13
    local.get 0
    local.get 1
    call 0)
  (func (;57;) (type 0) (param i32 i32)
    local.get 0
    i32.const 1049064
    i32.store offset=4
    local.get 0
    local.get 1
    i32.store)
  (func (;58;) (type 3) (param i32) (result i32)
    local.get 0
    i32.load8_u offset=4
    i32.const 2
    i32.and
    i32.const 1
    i32.shr_u)
  (func (;59;) (type 3) (param i32) (result i32)
    local.get 0
    i32.load offset=4
    i32.const 3
    i32.and
    i32.const 1
    i32.ne)
  (func (;60;) (type 3) (param i32) (result i32)
    i32.const 0
    local.get 0
    i32.sub
    local.get 0
    i32.and)
  (func (;61;) (type 3) (param i32) (result i32)
    local.get 0
    i32.load8_u offset=4
    i32.const 3
    i32.and
    i32.eqz)
  (func (;62;) (type 0) (param i32 i32)
    local.get 0
    local.get 1
    i32.const 3
    i32.or
    i32.store offset=4)
  (func (;63;) (type 3) (param i32) (result i32)
    local.get 0
    i32.load
    local.get 0
    i32.load offset=4
    i32.add)
  (func (;64;) (type 7)
    i32.const 1049108
    i32.const 17
    i32.const 1049128
    call 31
    unreachable)
  (func (;65;) (type 1) (param i32 i32) (result i32)
    local.get 0
    i32.load
    drop
    loop  ;; label = @1
      br 0 (;@1;)
    end
    unreachable)
  (func (;66;) (type 1) (param i32 i32) (result i32)
    local.get 0
    i64.load32_u
    local.get 1
    call 11)
  (func (;67;) (type 3) (param i32) (result i32)
    local.get 0
    i32.load offset=4
    i32.const -8
    i32.and)
  (func (;68;) (type 3) (param i32) (result i32)
    local.get 0
    i32.load offset=4
    i32.const 1
    i32.and)
  (func (;69;) (type 3) (param i32) (result i32)
    local.get 0
    i32.load offset=12
    i32.const 1
    i32.and)
  (func (;70;) (type 3) (param i32) (result i32)
    local.get 0
    i32.load offset=12
    i32.const 1
    i32.shr_u)
  (func (;71;) (type 0) (param i32 i32)
    local.get 0
    local.get 1
    i32.const 1049492
    i32.load
    local.tee 0
    i32.const 3
    local.get 0
    select
    call_indirect (type 0)
    unreachable)
  (func (;72;) (type 1) (param i32 i32) (result i32)
    local.get 0
    local.get 1
    i32.add)
  (func (;73;) (type 1) (param i32 i32) (result i32)
    local.get 0
    local.get 1
    i32.sub)
  (func (;74;) (type 3) (param i32) (result i32)
    local.get 0
    i32.const 8
    i32.add)
  (func (;75;) (type 3) (param i32) (result i32)
    local.get 0
    i32.const -8
    i32.add)
  (func (;76;) (type 6) (param i32) (result i64)
    i64.const -5786388802339902217)
  (func (;77;) (type 6) (param i32) (result i64)
    i64.const 9147559743429524724)
  (func (;78;) (type 6) (param i32) (result i64)
    i64.const -9040091204923801142)
  (func (;79;) (type 2) (param i32)
    nop)
  (func (;80;) (type 0) (param i32 i32)
    nop)
  (table (;0;) 21 21 funcref)
  (memory (;0;) 17)
  (global (;0;) (mut i32) (i32.const 1048576))
  (export "memory" (memory 0))
  (export "on_load" (func 30))
  (export "on_call" (func 19))
  (elem (;0;) (i32.const 1) func 49 37 80 79 52 9 26 78 42 15 18 32 57 41 76 77 66 65 79 78)
  (data (;0;) (i32.const 1048576) "helloWorldgetHelloWorldhelloworld.funcHelloWorldhelloworld.funcHelloWorld okhelloworld.viewGetHelloWorldhelloworld.viewGetHelloWorld okHello, world!object id mismatchD:\5cWork\5cgo\5cgithub.com\5ciotaledger\5cwasp\5cpackages\5cvm\5cwasmlib\5csrc\5cexports.rs\00\00\a6\00\10\00H\00\00\00$\00\00\00\0d\00\00\00\a6\00\10\00H\00\00\00)\00\00\00\09\00\00\00Rust:KEY_ZZZZZZZ\04\00\00\00\04\00\00\00\04\00\00\00\05\00\00\00\06\00\00\00\07\00\00\00\04\00\00\00\00\00\00\00\01\00\00\00\08\00\00\00called `Option::unwrap()` on a `None` valuelibrary/std/src/panicking.rs\00s\01\10\00\1c\00\00\00\eb\01\00\00\1f\00\00\00s\01\10\00\1c\00\00\00\ec\01\00\00\1e\00\00\00\09\00\00\00\10\00\00\00\04\00\00\00\0a\00\00\00\0b\00\00\00\04\00\00\00\08\00\00\00\04\00\00\00\0c\00\00\00\0d\00\00\00\0e\00\00\00\0c\00\00\00\04\00\00\00\0f\00\00\00\04\00\00\00\08\00\00\00\04\00\00\00\10\00\00\00library/alloc/src/raw_vec.rscapacity overflow\00\00\00\f8\01\10\00\1c\00\00\00\19\02\00\00\05\00\00\00\13\00\00\00\00\00\00\00\01\00\00\00\14\00\00\00index out of bounds: the len is  but the index is \00\00H\02\10\00 \00\00\00h\02\10\00\12\00\00\0000010203040506070809101112131415161718192021222324252627282930313233343536373839404142434445464748495051525354555657585960616263646566676869707172737475767778798081828384858687888990919293949596979899")
  (data (;1;) (i32.const 1049428) "\04")
  (data (;2;) (i32.const 1049440) "\04"))
