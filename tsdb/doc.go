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
│   RwMutex    │ │ │   │  RWMutex     │ │ │
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
│              │ │ │   │   SpinLock   │ │ │
+--------------+ │ │   +--------------+ │ │
  +----│---------+ │     +--------------+ │
    +--│-----------+       +--------------+
       │
       V
+--------------+
│              │
│              │-+
│              │ │
│   Segment    │ │-+
│   Store      │ │ │
│              │ │ │
│              │ │ │
│              │ │ │
│              │ │ │
│              │ │ │
+--------------+ │ │
  +--------------+ │
    +--------------+


━━━━━━━━━━━━━━━━━━━━━━━Layout of TagKeys Meta Table━━━━━━━━━━━━━━━━━━━━━━━━

                   Level1
                   +---------+---------+---------+---------+---------+---------+
                   │ TagKey  │ TagKey  │ TagKey  │ Offsets │ Bitmap  │ Footer  │
                   │  Meta   │  Meta   │  Meta   │         │         │         │
                   +---------+---------+---------+---------+---------+---------+
                  /           \                  |         |         |         |
                 /             \                 |          \        |         |
                /                \              /            \       |         |
               /                   \           /               \     |         |
  +-----------+                     |        /                   \    \         \
 /                     Level2       |       |                     |    \         |
v--------+--------+--------+--------v       v--------+---+--------v     v--------v
│  Trie  │TagValue│ Offsets│ Footer │       │ Offset │...│ Offset │     │ TagKV  │
│  Tree  │IDBitmap│        │        │       │        │   │        │     │ Bitmap │
+--------+--------+--------+--------+       +--------+---+--------+     +--------+


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


━━━━━━━━━━━━━━━━━━━━━━━Layout of Metric NameID Index Table━━━━━━━━━━━━━━━━━━━━━━━━
Metric-NameID-Table is a gzip compressed k/v pairs of metricNames and metricIDs on disk.

                   Level1
                   +---------+---------+---------+---------+
                   │ Metric  │  Meta   │ Index   │ Footer  │
                   │ KVPair  │         │         │         │
                   +---------+---------+---------+---------+

Level1(Metric NameID KVPair)
┌─────────────────────────────────────────────────────────────────┬─────────────────────┐
│            Gzip Compressed Metric K/V pairs                     │  SequenceNumber     │
├──────────┬──────────┬──────────┬──────────┬──────────┬──────────┼──────────┬──────────┤
│MetricName│MetricName│ MetricID │MetricName│MetricName│ MetricID │ MetricID │ TagKeyID │
│  Length  │          │          │  Length  │          │          │ Sequence │ Sequence │
├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┼──────────┼──────────┤
│ uvariant │ N Bytes  │ 4 Bytes  │ uvariant │ N Bytes  │ 4 Bytes  │ 4 Bytes  │ 4 Bytes  │
└──────────┴──────────┴──────────┴──────────┴──────────┴──────────┴──────────┴──────────┘


━━━━━━━━━━━━━━━━━━━━━━━Layout of Metric Meta Index Table━━━━━━━━━━━━━━━━━━━━━━━━
Metric-Meta stores meta info for metric,
such as tagKey, tagKeyID, fieldID, fieldName and fieldType etc.

                   Level1
                   +---------+---------+---------+---------+---------+---------+
                   │ Metric  │ Metric  │ Metric  │ Metric  │ Metric  │ Footer  │
                   │ Meta    │  Meta   │  Meta   │  Meta   │ Index   │         │
                   +---------+---------+---------+---------+---------+---------+
                  /         /          │         │\        +---------+
                 /         +           |         │ +----------+       \
                /          |           |         +-------+     \       \
               /           |           |                  \     \       \
  +-----------+            |           |                   \     \       \
 /                 Level2  |           |                    \     \       \
v--------+--------+--------v           v--------+---+--------v     v-------v
│ TagKey │  Field │ PosOf  │           │ Offset │...│ Offset │     │ Metric│
│   Meta │  Meta  │ Field  │           │        │   │        │     │ Bitmap│
+--------+--------+--------+           +--------+---+--------+     +-------+

Level2(TagKey Meta)
┌─────────────────────────────────────────────────────────────────┐
│                       TagKey Meta                               │
├──────────┬──────────┬──────────┬──────────┬──────────┬──────────┤
│  TagKey  │  TagKey  │ TagKeyID │  TagKey  │  TagKey  │  TagID   │
│   Len    │          │          │   Len    │          │          │
├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┤
│  1 Byte  │ N Bytes  │ 4 Bytes  │  1 Byte  │ N Bytes  │ 4 Bytes  │
└──────────┴──────────┴──────────┴──────────┴──────────┴──────────┘

Level2(Field Meta)
┌───────────────────────────────────────────────────────────────────────────────────────┬──────────┐
│                                    Field Meta                                         │          │
├──────────┬──────────┬──────────┬──────────┬──────────┬──────────┬──────────┬──────────┼──────────┤
│  Field   │  Field   │  Field   │  Field   │  Field   │  Field   │  Field   │  Field   │  PosOf   │
│   Len    │  Name    │  Type    │   ID     │   Len    │  Name    │  Type    │   ID     │  Field   │
├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┼──────────┼──────────┼──────────┤
│ uvariant │ N Bytes  │ 1 Byte   │ 2 Bytes  │ uvariant │ N Bytes  │  1 Byte  │ 2 Bytes  │ 4 Bytes  │
└──────────┴──────────┴──────────┴──────────┴──────────┴──────────┴──────────┴──────────┴──────────┘


━━━━━━━━━━━━━━━━━━━━━━━━━━Layout of Metric Data Table━━━━━━━━━━━━━━━━━━━━━━

                   Level1
                   +---------+---------+---------+---------+---------+---------+
                   │ Metric  │ Metric  │ Metric  │ Metric  │ Metric  │ Footer  │
                   │ Block   │ Block   │ Block   │ Offset  │ Index   │         │
                   +---------+---------+---------+---------+---------+---------+
                  /           \                   \        │\        +-------------------------------+
                 /             \                   \       │ +--------------------------------+       \
                /               \                   \      +-----------------------------+     \       \
               /                 \                   +--------------+                     \     \       \
  +-----------+                   +--------+                         \                     \     \       \
 /                 Level2                   \                         \                     \     \       \
v--------+--------+--------+--------+--------v                         v--------+---+--------v     v-------v
│ Version│ Version│ Version│ Version│ Footer │                         │ Offset │...│ Offset │     │ Metric│
│ Entry1 │ Entry2 │ Entry3 │ Offsets│        │                         │        │   │        │     │ Bitmap│
+--------+--------+--------+--------+--------+                         +--------+---+--------+     +-------+
│        │
│        │
│        │         Level3
v--------v--------+--------+--------+--------+--------+--------v
│ Series │ Series │ Series │ Series │ Series │ Fields │ Footer │
│ Entry  │ Entry  │ Entry  │ Offset │ Bitmap │  Meta  │        │
+--------+--------+--------+--------+--------+--------+--------+
│         \                 \       │\        \
│          \                 \      │ \        +-----------------------------------------------+
│           \                 \     │  +----------------------------------------------+         \
│            \                 \    +---------------------------------------------+    \         \
│             \                 +-----------------------------+                    \    \         \
│              +------------------------------------+          \                    \    \         \
│                  Level4                            \          \                    \    \         \
v--------+--------+--------+--------+--------+--------v          v--------+---+-------v    v---------v
│ Fields │ Data   │  Data  │ Data   │ Data   │  Data  │          │ Offset │...│ Offset│    │seriesID │
│ Info   │        │        │        │        │        │          │        │   │       │    │ Bitmap  │
+--------+--------+--------+--------+--------+--------+          +--------+---+-------+    +---------+


Level1(KV table: MetricBlocks, Offset, Keys)
┌───────────────────────────────────────────┬───────────────────────────────────────────┐
│               Metric Blocks               │           Offset And Keys                 │
├──────────┬──────────┬──────────┬──────────┼──────────┬──────────┬──────────┬──────────┤
│  length  │  Metric  │  length  │  Metric  │  length  │  Offset  │  length  │  Keys    │
│          │  Block1  │          │  Block2  │          │          │          │          │
├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┼──────────┼──────────┤
│ uvariant │  N Bytes │ uvariant │ N Bytes  │ uvariant │  N Bytes │ uvariant │  N Bytes │
└──────────┴──────────┴──────────┴──────────^──────────┴──────────^──────────┴──────────┘
                                            │                     │
                                       posOfOffset             posOfKeys

Level1(KV table: Footer)
┌──────────────────────────────────────────────────────┐
│                    Footer                            │
├──────────┬──────────┬──────────┬──────────┬──────────┤
│  length  │ position │ position │  Table   │  Magic   │
│          │ OfOffset │ OfKeys   │ Version  │  Number  │
├──────────┼──────────┼──────────┼──────────┼──────────┤
│  1 Byte  │ 4 Bytes  │ 4 Bytes  │ 1 Bytes  │  8 Bytes │
└──────────┴──────────┴──────────┴──────────┴──────────┘


Level2(Version Offsets Block)
Same as Level2 in ForwardIndexTable
┌────────────────────────────────┐┌──────────────────────────────────────────────────────┐┌─────────────────────┐
│          Version Entries       ││                     Version Offsets                  ││        Footer       │
├──────────┬──────────┬──────────┤├──────────┬──────────┬──────────┬──────────┬──────────┤├──────────┬──────────┤
│  Version │  Version │  Version ││ Versions │ Version1 │ Version1 │ Version2 │ Version2 ││VersionOff│ CRC32    │
│  Entry1  │  Entry2  │  Entry3  ││  Count   │   int64  │  Length  │   int64  │  Length  ││ setsPos  │ CheckSum │
├──────────┼──────────┼──────────┤├──────────┼──────────┼──────────┼──────────┼──────────┤├──────────┼──────────┤
│  N Bytes │  N Bytes │  N Bytes ││ uvariant │  8 Bytes │ uvariant │  8 Bytes │ uvariant ││ 4 Bytes  │ 4 Bytes  │
└──────────┴──────────┴──────────┘└──────────┴──────────┴──────────┴──────────┴──────────┘└──────────┴──────────┘

Level3(Fields Meta)
┌─────────────────────────────────────────────────────────────────┐
│                      Fields Meta                                │
├──────────┬──────────┬──────────┬──────────┬──────────┬──────────┤
│   Count  │ FieldID  │  Field   │ FieldName│ FieldName│          │
│          │ (uint16) │  Type    │  Length  │          │  ......  │
├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┤
│ uvariant │  2 Bytes │ 1 Byte   │ uvariant │ N Bytes  │          │
└──────────┴──────────┴──────────┴──────────┴──────────┴──────────┘

Level3(Series Footer)
┌─────────────────────┐
│   Series Footer     │
├──────────┬──────────┤
│  Series  │FieldsMeta│
│ BitMapPos│   Pos    │
├──────────┼──────────┤
│ 4 Bytes  │ 4 Bytes  │
└──────────┴──────────┘


Level4(Fields Info, Fields Data)
┌────────────────────────────────┬─────────────────────┐
│          Fields Info           │   Fields Data       │
├──────────┬──────────┬──────────┼──────────┬──────────┤
│ BitArray │  Data1   │  Data2   │  Data1   │ Data2    │
│          │  Length  │  Length  │          │          │
├──────────┼──────────┼──────────┼──────────┼──────────┤
│ N Bytes  │ uvariant │ uvariant │ N Bytes  │ N Bytes  │
└──────────┴──────────┴──────────┴──────────┴──────────┘
bit array example(10101001, 1010100110101001)


*/
