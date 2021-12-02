// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// func select64(x uint64, k int64) int64
TEXT ·select64(SB),$24-24
    MOVQ    x+0(FP), AX
    MOVQ    k+8(FP), CX
    CMPB    ·hasBMI2(SB), $0
    JEQ     fallback
    DECQ    CX
    MOVQ    $1, BX
    SHLQ    CX, BX
    PDEPQ   AX, BX, BX
    TZCNTQ  BX, BX
    MOVQ    BX, ret+16(FP)
    RET
fallback:
    MOVQ    AX, (SP)
    MOVQ    CX, 8(SP)
    CALL    ·select64Broadword(SB)
    MOVQ    16(SP), AX
    MOVQ    AX, ret+16(FP)
    RET
