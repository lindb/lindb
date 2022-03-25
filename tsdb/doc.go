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

package tsdb

/*


━━━━━━━━━━━━━━━━━━━━━━━━━━IO Flow━━━━━━━━━━━━━━━━━━━━━━━━

Each shard contains a MemoryDatabase, the Index Database is a global singleton

a) Write Flow

+-------------------------------------+
│                                     │
│               Engine                │
│                                     │
+------+----------------------+-------+
       │                      │
       │                      │
Shard  │               Shard  │
+------v-------+       +------v-------+
│ Data │ Memory│       │ Data │ Memory+----------------------------------------------------+
│  DB  │   DB  │       │  DB  │   DB  +--------------+                                     │
+-----^-+------+       +-----^-+------+              │                                     │
      │ │                    │ │                     │                                     │
      │ │                    │ │ ID                  │                                     │
      │ │                    │ │ Generator           │                                     │
+-----+-v--------------------+-v------+              │                                     │
│                                     │              │                                     │
│            ID Sequencer             │              │                                     │
│                                     │              +--------------+                      │
+------+----------------------+-------+              │              │                      │
       │                      │                      │              │                      │
       │                      │                      │ SeriesIndex- │ ForwardIndex-        │
       │ NameIDIndexFlusher   │ MetaIndexFlusher     │ Flusher      │ Flusher              │ MetricDataFlusher
+------v-------+       +------v-------+       +------v-------+------v-------+       +------v-------+
│ MetricNameID │       │  MetricMeta  │       │SeriesInverted│ SeriesForward│       │  MetricData  │
│  IndexTable  │       │  IndexTable  │       │  IndexTable  │  IndexTable  │       │    Table     │
+--------------+       +--------------+       +--------------+--------------+       +--------------+


b) Query flow

Shard                  Shard
+------+-------+       +-----+--------+                Suggester
│ Data │ Memory│       │ Data │ Memory<--------------+ MetaGetter
│  DB  │   DB  │       │  DB  │   DB  │              │ Filter
+-----^-+------+       +-----^-+------+              │ Load
      │ │                    │ │                     +----------------------
      │ │                    │ │
      │ │           IDGetter │ │
+-----+-v--------------------+-v------+
│                                     <----------------------------------------------------+
│            ID Sequencer             │                                                    │
│                                     <--------------+--------------+                      │
+------^----------------------^-------+              │              │                      │
       │                      │                      │ Suggest-     │                      │
       ^ SuggestMetrics       ^ SuggestTagKeys       ^ TagValues    ^                      ^
       │ NameIDIndexReader    │ MetaIndexReader      │ Filter       │ MetaGetter           │ Load
+------+-------+       +------+-------+       +------+-------+------+-------+       +------+-------+
│ MetricNameID │       │  MetricMeta  │       │SeriesInverted│ SeriesForward│       │  MetricData  │
│  IndexTable  │       │  IndexTable  │       │  IndexTable  │  IndexTable  │       │    Table     │
+--------------+       +--------------+       +--------------+--------------+       +--------------+



━━━━━━━━━━━━━━━━━━━━━━━━━━Layout of MemoryDatabase━━━━━━━━━━━━━━━━━━━━━━━━

+--------------+       +--------------+
│              │------>│              │
│              │-+     │              │-+
│   Memory     │ │     │  Metric      │ │
│   Database   │ │-+   │  Store       │ │-+
│   RwMutex    │ │ │   │              │ │ │
│              │ │ │   │              │ │ │
│              │ │ │   │              │ │ │
│              │ │ │   │              │ │ │
│              │ │ │   │              │ │ │
│              │ │ │   │              │ │ │
+-+------------+ │ │   +--------------+ │ │
  +--------------+ │     +-----│--------+ │
    +--------------+       +---│----------+
                               │
                               V
+--------------+       +--------------+
│              │<------│              │
│              │-+     │              │-+
│              │ │     │              │ │
│   Field      │ │-+   │  TimeSeries  │ │-+
│   Store      │ │ │   │  Store       │ │ │
│              │ │ │   │              │ │ │
│              │ │ │   │              │ │ │
│              │ │ │   │              │ │ │
│              │ │ │   │              │ │ │
│              │ │ │   │              │ │ │
+--------------+ │ │   +--------------+ │ │
  +--------------+ │     +--------------+ │
    +--------------+       +--------------+


━━━━━━━━━━━━━━━━━━━━━━━Layout of TagKeys Meta Table━━━━━━━━━━━━━━━━━━━━━━━━

                   Level1
                   +---------+---------+---------+---------+---------+---------+
                   │ TagKey  │ TagKey  │ TagKey  │ Offsets │ Bitmap  │ Footer  │
                   │  Meta   │  Meta   │  Meta   │         │         │         │
                   +---------+---------+---------+---------+---------+---------+
                  /           \                  |         |
                 /             \                 |          \
                /                \              /            \
               /                   \           /               \
  +-----------+                     |        /                   \
 /                     Level2       |       |                     |
v--------+--------+--------+--------v       v--------+---+--------v
│  Trie  │TagValue│ Offsets│ Footer │       │ Offset │...│ Offset │
│  Tree  │IDBitmap│        │        │       │        │   │        │
+--------+--------+--------+--------+       +--------+---+--------+


Level1(KV table: TagKeyID -> TagKeyMeta data)
Level1 is same as MetricDataTable as below

Level2(Footer)

┌───────────────────────────────────────────┐
│                 Footer                    │
├──────────┬──────────┬──────────┬──────────┤
│  BitMap  │  Offsets │ TagValue │  CRC32   │
│ Position │ Position │ Sequence │ CheckSum │
├──────────┼──────────┼──────────┼──────────┤
│ 4 Bytes  │ 4 Bytes  │ 4 Bytes  │ 4 Bytes  │
└──────────┴──────────┴──────────┴──────────┘


━━━━━━━━━━━━━━━━━━━━━━━Layout of Metric Inverted Index Table━━━━━━━━━━━━━━━━━━━━━━━━

                   Level1
                   +---------+---------+---------+---------+---------+---------+
                   │ Inverted│ Inverted│ Inverted│ Blocks  │ Blocks  │ Footer  │
                   │ Block   │  Block  │  Block  │ Offsets │ Offsets │         │
                   +---------+---------+---------+---------+---------+---------+
                  /          |
                 /            \
                /               \
               /                  \
  +-----------+                    +-------+
 /                 Level2                   \
v--------+--------+--------+--------+--------v
│TagValue│TagValue│TagValue| HighKey│ Footer │
│ Bucket │ Bucket │ Bitmap | Offsets│        │
+--------+--------+--------+--------+--------+
|        |
|        |         Level3
v--------+--------+--------+--------v
│SeriesID│SeriesID│ LowKey │Offsets │
│ Bitmap │ Bitmap │ Offsets│Length  │
+--------+--------+--------+--------+


Level2(Footer)
┌─────────────────────────────────────────────────────────────────┐
│                       Level2 Footer                            │
├──────────┬──────────┬──────────┬──────────┬──────────┬──────────┤
│  TagKey  │  TagKey  │ TagKeyID │  TagKey  │  TagKey  │  TagID   │
│   Len    │          │          │   Len    │          │          │
├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┤
│  1 Byte  │ N Bytes  │ 4 Bytes  │  1 Byte  │ N Bytes  │ 4 Bytes  │
└──────────┴──────────┴──────────┴──────────┴──────────┴──────────┘


━━━━━━━━━━━━━━━━━━━━━━━━━━Layout of Metric Data Table━━━━━━━━━━━━━━━━━━━━━━

                   Level1
                   +---------+---------+---------+---------+---------+
                   │ Metric  │ Metric  │ Metric  │ Metric  │ Footer  │
                   │ Block   │ Block   │ Offsets │ Bitmap  │         │
                   +---------+---------+---------+---------+---------+
                  /           \
                 /             \
                /               \
               /                 \
  +-----------+                   +-----------------+
 /                 Level2                            \
v--------+--------+--------+--------+--------+--------v
│ Series │ Series │  Field | Series │ HighKey│ Footer │
│ Bucket │ Bucket │  Metas | Bitmap │ Offsets│        │
+--------+--------+--------+--------+--------+--------+
│        │
│        │
│        │         Level3
v--------v--------+--------+--------+--------+
│ Series │ Series │ Series │ LowKey │ PosOf  |
│ Entry  │ Entry  │ Entry  │ Offsets│Offsets |
+--------+--------+--------+--------+--------+
│         \        \        \
│          \        \        \
│           \        \        |
│            \        \       +--------------------------+
│             \        +--------------------------+       \
│              +------------------+                \       \
│                  Level4          \                \        \
v--------+--------+--------+--------+                +--------+
│ Field  │ Field  │ Field  | LenOf  │                │ Field  |
│ Data   │ Data   │ Offsets│ Offsets│                │ Data   |
+--------+--------+--------+--------+                +--------+


Level1 (KV table: MetricBlocks, Offset, Keys)
┌───────────────────────────────────────────┬─────────────────────┐
│               Metric Blocks               │  Offset And Keys    │
├──────────┬──────────┬──────────┬──────────┼──────────┬──────────┤
│  Metric  │  Metric  │  Metric  │  Metric  │  Offset  │  High    │
│  Block1  │  Block2  │  Block3  │  Block4  │          │  Keys    │
├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┤
│  N Bytes │  N Bytes │  N Bytes │ N Bytes  │  N Bytes │  N Bytes │
└──────────┴──────────┴──────────┴──────────^──────────^──────────┘
                                            │          │
                                       posOfOffset  posOfKeys

Level2 (KV table: Series Bucket Footer)
┌──────────────────────────────────────────────────────┐
│                    Footer                            │
├──────────┬──────────┬──────────┬──────────┬──────────┤
│   time   │ position │ position │ position │  CRC32   │
│   range  │ OfMetas  │ OfBitmap │ OfOffsets│ CheckSum │
├──────────┼──────────┼──────────┼──────────┼──────────┤
│  4 Byte  │ 4 Bytes  │ 4 Bytes  │ 4 Bytes  │  4 Bytes │
└──────────┴──────────┴──────────┴──────────┴──────────┘


Level2(SingleField Meta)
┌─────────────────────────────────────────────────────────────────┐
│                      SingleField Meta                                │
├──────────┬──────────┬──────────┬──────────┬──────────┬──────────┤
│   Count  │ FieldID  │  Field   │ FieldID  │  Field   │          │
│          │ (uint16) │  Type    │ (uint16) │  Type    │  ......  │
├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┤
│  1 Byte  │  1 Bytes │ 1 Byte   │  1 Bytes │ 1 Byte   │          │
└──────────┴──────────┴──────────┴──────────┴──────────┴──────────┘


*/
