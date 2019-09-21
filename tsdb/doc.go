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
+-----^-+------+       +-----^-+------+              │ Scanner
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
       │ NameIDIndexReader    │ MetaIndexReader      │ Filter       │ MetaGetter           │ Scanner
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

━━━━━━━━━━━━━━━━━━━━━━━Layout of Series Forward Index Table━━━━━━━━━━━━━━━━━━━━━━━━

                   Level1
                   +---------+---------+---------+---------+---------+---------+
                   │ Metric  │ Metric  │ Metric  │ Metric  │ Metric  │ Footer  │
                   │ Block   │ Block   │ Block   │ Offset  │ Index   │         │
                   +---------+---------+---------+---------+---------+---------+
                  /           \                  │          \         \
                 /             \                 │           \         \
                /               \                +            \         +------+
               /                 \                \            \                \
  +-----------+                   +--------+       \            +--------+       \
 /                 Level2                   \       \                     \       \
v--------+--------+--------+--------+--------v       v--------+---+--------v-------v
│ Version│ Version│ Version│ Version│ Footer │       │ Offset │...│ Offset │ Metric│
│ Entry1 │ Entry2 │ Entry3 │ Offsets│        │       │        │   │        │ Bitmap│
+--------+--------+--------+--------+--------+       +--------+---+--------+-------+
         │        │
         │        │
  +------+        +------------------------------------------+
 /                 Level3                                     \
v--------+--------+--------+--------+--------+--------+--------v
│  Time  │ TagKeys│ Dict   │ Series │Offsets │SeriesID│ Footer │
│  Range │ Block  │ Block  │LUTBlock│ Block  │ BitMap │        │
+--------+--------+--------+--------+--------+--------+--------+

Level1(KV table: MetricID -> MetricBlock)
Level1 is same as MetricDataTable as below


Level2(Version Offsets Block)
┌────────────────────────────────┐┌──────────────────────────────────────────────────────┐┌─────────────────────┐
│          Version Entries       ││                     Version Offsets                  ││        Footer       │
├──────────┬──────────┬──────────┤├──────────┬──────────┬──────────┬──────────┬──────────┤├──────────┬──────────┤
│  Version │  Version │  Version ││ Versions │ Version1 │ Version1 │ Version2 │ Version2 ││VersionOff│ CRC32    │
│  Entry1  │  Entry2  │  Entry3  ││  Count   │   int64  │  Length  │   int64  │  Length  ││ setsPos  │ CheckSum │
├──────────┼──────────┼──────────┤├──────────┼──────────┼──────────┼──────────┼──────────┤├──────────┼──────────┤
│  N Bytes │  N Bytes │  N Bytes ││ uvariant │  8 Bytes │ uvariant │  8 Bytes │ uvariant ││ 4 Bytes  │ 4 Bytes  │
└──────────┴──────────┴──────────┘└──────────┴──────────┴──────────┴──────────┴──────────┘└──────────┴──────────┘


Level3(Version Entry Block)
TagKeysBlock stores all tagKeys of the metric
┌─────────────────────┐┌──────────────────────────────────────────────────────┐┌──────────┐┌─────────────────────┐
│  Time Range Block   ││                      TagKeys Block                   ││Dict Block││      Tags Blocks    │
├──────────┬──────────┤├──────────┬──────────┬──────────┬──────────┬──────────┤├──────────┤├──────────┬──────────┤
│   Start  │   End    ││  TagKey  │  TagKey1 │  TagKey1 │  TagKey2 │  TagKey2 ││          ││TagsBlock1│TagsBlock2│
│TimeDelta │TimeDelta ││  Count   │  Length  │          │  Length  │          ││  .....   ││          │          │
├──────────┼──────────┤├──────────┼──────────┼──────────┼──────────┼──────────┤├──────────┤├──────────┼──────────┤
│  4 Bytes │  4 Bytes ││ uvariant │ uvariant │  N Bytes │ uvariant │  N Bytes ││  N Bytes ││  N Bytes │  N Bytes │
└──────────┴──────────┘└──────────┴──────────┴──────────┴──────────┴──────────┘└──────────┘^──────────^──────────┘
                                                                                           │          │
                                                                                         PosOfTags1 PosOfTags2

Level3(Dict Block)
Dict Block is composed of 2 parts:
1) String Block Offsets
   TagValues of the metric are split into multi string blocks(each block size is up to 400)

2) Snappy Compressed String Blocks
   Theoretically, one compressed string block may cost 1-3 pages(4KB/page)

┌───────────────────────────────────────────┐┌──────────┐┌───────────────────────────────────────────┐
│       Snappy Compressed String Block      ││ StrBlocks││             String Block Offsets          │
├──────────┬──────────┬──────────┬──────────┤├──────────┤├──────────┬──────────┬──────────┬──────────┤
│ TagValue1│ TagValue1│ TagValue2│ TagValue2││  ....... ││  Strings │ StrBlock1│ StrBlock2│ StrBlock3│
│  Length  │          │  Length  │          ││          ││  Count   │  Length  │  Length  │  Length  │
├──────────┼──────────┼──────────┼──────────┤├──────────┤├──────────┼──────────┼──────────┼──────────┤
│ uvariant │  N Bytes │ uvariant │  N Bytes ││          ││ uvariant │ uvariant │ uvariant │ uvariant │
└──────────┴──────────┴──────────┴──────────┘└──────────┘^──────────┴──────────┴──────────┴──────────┘
 \____________________  ___________________/             │
                      \/                                 │
                StrBlock1Length                      PosOfDictBlockOffsets


Level3(Series TagsKeyValue LOOKUP-TABLE Block)
SeriesTagsBlock is composed of 2 parts:
1) bit-array of tagKeys of this seriesID
   TagKeys Block stores all tagKeys of this seriesID,
   If there are 15 tagKeys of the metric, and this series is composed of the 1st, 3rd,5th, 14th,
   then the bit-array is 0101,0100,0000,0010. Offsets are listed in order after the bit-array,

2) tagValue offsets is used to index for the dict block
   each tagValue-index is uvariant encoded

┌──────────────────────────────────────────────────────┐
│             Series Tags LOOKUP-TABLE Block           │
├──────────┬──────────┬──────────┬──────────┬──────────┤
│ TagKeys  │ StrBlock │ StrBlock │ StrBlock │ StrBlock │
│ BitArray │ Sequence1│ Sequence2│ Sequence3│ Sequence4│
├──────────┼──────────┼──────────┼──────────┼──────────┤
│ N Bytes  │ uvariant │ uvariant │ uvariant │ uvariant │
└──────────┴──────────┴──────────┴──────────┴──────────┘


Level3(Footer)
┌────────────────────────────────┐
│              Footer            │
├──────────┬──────────┬──────────┤
│PosOfDictB│ PosOfOff │  PosOf   │
│lockOffset│ setBlock │  BitMap  │
├──────────┼──────────┼──────────┤
│ 4 Bytes  │ 4 Bytes  │  4 Bytes │
└──────────┴──────────┴──────────┘


━━━━━━━━━━━━━━━━━━━━━━━Layout of Series Inverted Index Table━━━━━━━━━━━━━━━━━━━━━━━━

                   Level1
                   +---------+---------+---------+---------+---------+---------+
                   │  TagKV  │  TagKV  │  TagKV  │  TagKV  │  TagKV  │ Footer  │
                   │ EntrySet│ EntrySet│ EntrySet│ Offset  │  Index  │         │
                   +---------+---------+---------+---------+---------+---------+
                  /           \                   \        │\        +-------------------------------+
                 /             \                   \       │ +--------------------------------+       \
                /               \                   \      +-----------------------------+     \       \
               /                 \                   +--------------+                     \     \       \
  +-----------+                   +-----------------+                \                     \     \       \
 /                 Level2                            \                \                     \     \       \
v--------+--------+--------+--------+--------+--------v                v--------+---+--------v     v-------v
│  Time  │ LOUDS  │TagValue│TagValue│Offsets │ Footer │                │ Offset │...│ Offset │     │ TagKV │
│  Range │TrieTree│ Data1  │ Data2  │        │        │                │        │   │        │     │ Bitmap│
+--------+--------+--------+--------+--------+--------+                +--------+---+--------+     +-------+


Level1(KV table: TagKeyID -> EntrySetBlock)
Level1 is same as MetricDataTable as below
This block is alias as EntrySetBlock

Level2(TimeRange & LOUDS Encoded Trie Tree)
This block is alias as TreeBlock
┌─────────────────────┬────────────────────────────────────────────────────────────────────────────┐
│       TimeRange     │                        LOUDS Encoded Trie Tree                             │
├──────────┬──────────┼──────────┬──────────┬──────────┬──────────┬──────────┬──────────┬──────────┤
│ StartTime│  EndTime │   Trie   │  Labels  │  labels  │ isPrefix │ isPrefix │  LOUDS   │  LOUDS   │
│   int64  │   int64  │  TreeLen │  Length  │  Block   │ Key Len  │Key BitMap│  Length  │  BitMap  │
├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┼──────────┼──────────┼──────────┤
│  8 Bytes │  8 Bytes │ uvariant │ uvariant │ N Bytes  │ uvariant │ N Bytes  │ uvariant │ N Bytes  │
└──────────┴──────────┴──────────┴──────────┴──────────┴──────────┴──────────┴──────────┴──────────┘

Level2(Versioned TagValue Data)
alias as TagValueDataBlock
┌──────────┬──────────────────────────────────────────────────────┬─────────────────────┐
│          │                  VersionedTagValue                   │  VersionedTagValues │
├──────────┼──────────┬──────────┬──────────┬──────────┬──────────┼──────────┬──────────┤
│ Version  │ Version1 │StartTime1│ EndTime1 │ BitMap1  │ TagValue1│  Version │ Version  │
│  Count   │   int64  │ (Delta)  │  (Delta) │  Length  │  BitMap  │   Meta2  │  Meta3   │
├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┼──────────┼──────────┤
│ uvariant │ 8 Bytes  │ variant  │ variant  │ uvariant │ N Bytes  │ N Bytes  │  N Bytes │
└──────────┴──────────┴──────────┴──────────┴──────────┴──────────┴──────────┴──────────┘

Level2(Footer)

┌─────────────────────┐
│         Footer      │
├──────────┬──────────┤
│ Offsets  │  CRC32   │
│ Position │ CheckSum │
├──────────┼──────────┤
│ 4 Bytes  │ 4 Bytes  │
└──────────┴──────────┘

Succinct trie tree(Example):
(KEY Value: eleme:1, etcd:2, etrace:3)

Labels: eltecrmdaece
isPrefixKey: 0000000010101
LOUDS: 1011010110101010100100100
Values: [2, 1, 3]


                   +--------+
                   │        │ (pseudo root)
                   │  10    │ (node-0)
                   +--------+
                       │
                   +---v----+
                   │        │ (root)
                   │  10    │ (node-1)
                   +--------+
                       │
                   +---v----+
                   │   e    │
                   │  110   │ (node-2)
                   +---+----+
                      / \
              +------+   +----+
             /                 \
        +---v----+          +---v----+
        │   l    │          │   t    │
        │   10   │ (node-3) │   110  │(node-4)
        +---+----+          +---+----+
            │                   │\_______________
            │                   │                \
        +---v----+          +---v----+        +---v----+
        │   e    │          │   c    │        │   r    │
        │   10   │ (node-5) │   10   │(node-6)│   10   │ (node-7)
        +---+----+          +---+----+        +---+----+
            │                   │                 │
        +---v----+          +---v----+        +---v----+
        │   m    │          │   d    │        │   a    │
        │   10   │ (node-8) │   0    │(node-9)│   10   │ (node-10)
        +---+----+          +--------+        +---+----+
            │                 Value:2             │
        +---v----+                            +---v----+
        │   e    │                            │   c    │
        │   0    │ (node-11)                  │   10   │ (node-12)
        +--------+                            +---+----+
          Value:1                                 │
                                              +---v----+
                                              │   e    │
                                              │   0    │ (node-13)
                                              +--------+
                                               Value:3


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
┌───────────────────────────────────────────────────────────────────────────────────────┐
│                                 Fields Meta                                           │
├──────────┬──────────┬──────────┬──────────┬──────────┬──────────┬──────────┬──────────┤
│StartTime │ EndTime  │ Count    │ FieldID  │  Field   │ FieldName│ FieldName│          │
│ (delta)  │ (delta)  │          │ (uint16) │  Type    │  Length  │          │  ......  │
├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┼──────────┼──────────┤
│ variant  │ variant  │ uvariant │  2 Bytes │ 1 Byte   │ uvariant │ N Bytes  │          │
└──────────┴──────────┴──────────┴──────────┴──────────┴──────────┴──────────┴──────────┘

Level3(Series Footer)
┌────────────────────────────────┐
│           Series Footer        │
├──────────┬──────────┬──────────┤
│ Series   │  Series  │FieldsMeta│
│ OffsetPos│ BitMapPos│   Pos    │
├──────────┼──────────┼──────────┤
│ 4 Bytes  │ 4 Bytes  │ 4 Bytes  │
└──────────┴──────────┴──────────┘


Level4(Fields Info, Fields Data)
┌──────────────────────────────────────────────────────┬─────────────────────┐
│               Fields Info                            │   Fields Data       │
├──────────┬──────────┬──────────┬──────────┬──────────┼──────────┬──────────┤
│StartTime │ EndTime  │ BitArray │  Data1   │  Data2   │  Data1   │ Data2    │
│ (delta)  │ (delta)  │          │  Length  │  Length  │          │          │
├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┼──────────┤
│ variant  │ variant  │ N Bytes  │ uvariant │ uvariant │ N Bytes  │ N Bytes  │
└──────────┴──────────┴──────────┴──────────┴──────────┴──────────┴──────────┘
bit array example(10101001, 1010100110101001)


*/
