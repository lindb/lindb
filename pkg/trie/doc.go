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

package trie

// This package is heavily reference to the golang implementation of SuRF
// https://github.com/bobotu/myk/tree/master/surf
//
// Only the sparse layer is supported in this package,
// Because what we need is just the full key succinct trie without FPR
//
// In LinDB, the design of KV and other file format is based on the rule,
// where different ids are always compact, so we use the bitmap and fixed-offsets
// to index data in a file.
// And, we will always stores string as key, thus dense layer is not suitable.
//
// Ref: https://en.wikipedia.org/wiki/Trie
//
// A succinct data structure supporting rank/select efficiently is used here for querying and filtering data.
// Definition: https://en.wikipedia.org/wiki/Succinct_data_structure
// [1] SuRF: Practical Range Query Filtering with Fast Succinct Tries:
//     https://db.cs.cmu.edu/papers/2018/mod601-zhangA-hm.pdf
// [2] Fast, Small, Simple Rank/Select on Bitmaps:
//     https://users.dcc.uchile.cl/~gnavarro/ps/sea12.1.pdf
//
// nj11
// nj-2
// nj-3
// sh-4
// sh-5
// sh-6000
// bj-777
// b
// abcdef
// abcdefg
// bj-9
//
// Expanded Tree:
//                    [ ]
//       ______________|___________
//     a/        |b        |n      \s
//    [ ]       [ ]       [ ]      [ ]
//     |b      /$  \j      \j        \ h
//    [ ]    [Ø]   [ ]     [ ]___    [ ]
//     |c           |-      \-   \1    |-
//    [ ]          [ ]      [ ]  [Ø]  [ ]_____
//     |d         /7  \9    /2 \3      |4  \5  \6
//    [ ]       [Ø]   [Ø]  [Ø] [Ø]    [Ø]  [Ø] [Ø]
//     |e
//    [ ]
//     |f
//    [ ]
//    /$ \g
// [Ø]  [$]
//
//
//
// Compact tree:
//                   [ ]
//      ______________|___________
//    a/        |b        |n      \s
//   [ ]       [ ]       [ ]      [ ]
//    |bcdef    |         |j       |h-
//  /$ \g     /$ \j     /- \1     /4 \5  \6
// [Ø]  [$] [Ø]  [ ]  [ ]  [ ]   [Ø] [Ø] [Ø]
//               |-
//             /7 \9   |2 \3
//            [ ] [ ] [ ] [ ]
//
//
// labels: abns$g$j-14567923
//
// prefix offsets: 0, 5, 6, 8
// prefix data: bcdefjh--
// prefix bits: 0101110
//
// hasChild bits: 11110001100000000
//
// LOUDS bits: 10001010101001010
//
// Operations:
// S-ChildNodePos(pos) = select1(S-LOUDS, rank1(SHasChild, pos) + 1)
// to move up, S-ParentNodePos(pos) = select1(SHasChild, rank1(S-LOUDS, pos) - 1)
// to access a value, S-ValuePos(pos) = pos - rank1(S-HasChild, pos)
//
// https://github.com/plar/go-adaptive-radix-tree/tree/master/test/assets
// test data is from go-adaptive-radix-tree
