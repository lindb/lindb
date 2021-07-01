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

package kv

import (
	"errors"
	"fmt"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/logger"
)

//go:generate mockgen -source ./compact_job.go -destination=./compact_job_mock.go -package kv

// CompactJob represents the compact job which does merge sst files
type CompactJob interface {
	// Run runs compact logic
	Run() error
}

// compactJob represents the compaction job, merges input files
type compactJob struct {
	family Family
	state  *compactionState
	merger NewMerger
	rollup Rollup
}

// newCompactJob creates a compaction job
func newCompactJob(family Family, state *compactionState, rollup Rollup) CompactJob {
	return &compactJob{
		family: family,
		merger: family.getNewMerger(),
		state:  state,
		rollup: rollup,
	}
}

// Run runs compact job
func (c *compactJob) Run() error {
	compaction := c.state.compaction
	switch {
	case compaction.IsTrivialMove():
		c.moveCompaction()
	default:
		if err := c.mergeCompaction(); err != nil {
			return err
		}
	}
	return nil
}

// moveCompaction moves low level file to  up level, just does metadata change
func (c *compactJob) moveCompaction() {
	compaction := c.state.compaction
	kvLogger.Info("starting compaction job, just move file to next level", logger.String("family", c.family.familyInfo()))
	//move file to next level
	fileMeta := compaction.GetLevelFiles()[0]
	level := compaction.GetLevel()
	compaction.DeleteFile(level, fileMeta.GetFileNumber())
	compaction.AddFile(level+1, fileMeta)
	c.family.commitEditLog(compaction.GetEditLog())
	// TODO add cost ?????
	kvLogger.Info("finish move file compaction", logger.String("family", c.family.familyInfo()))
}

// mergeCompaction merges input files to up level
func (c *compactJob) mergeCompaction() (err error) {
	kvLogger.Info("starting compaction job, do merge compaction",
		logger.String("family", c.family.familyInfo()))
	defer func() {
		// cleanup compaction context, include temp pending output files
		c.cleanupCompaction()
	}()

	// do merge logic
	if err = c.doMerge(); err != nil {
		return err
	}
	// if merge success install compaction results into manifest
	c.installCompactionResults()
	return nil
}

// doMerge merges the input files based on merger interface which need use implements
func (c *compactJob) doMerge() error {
	it, err := c.makeInputIterator()
	if err != nil {
		return err
	}
	merger := c.merger()
	if c.rollup != nil {
		merger.Init(map[string]interface{}{RollupContext: c.rollup})
	}

	var needMerge [][]byte
	var previousKey uint32
	start := true
	for it.HasNext() {
		key := it.Key()
		value := it.Value()
		switch {
		case start || key == previousKey:
			// if start or same keys, append to need merge slice
			needMerge = append(needMerge, value)
			start = false
		case key != previousKey:
			//FIXME stone1100 merge data maybe is one block

			// 1. if new key != previous key do merge logic based on user define
			mergedValue, err := merger.Merge(previousKey, needMerge)
			if err != nil {
				return err
			}
			// 2. add new k/v pair into new store build
			if err := c.add(previousKey, mergedValue); err != nil {
				return err
			}
			// 3. prepare next merge loop
			// init value for next loop
			needMerge = needMerge[:0]
			// add value to need merge slice
			needMerge = append(needMerge, value)
		}
		// set previous merge key
		previousKey = key
	}

	// if has pending merge values after iterator, need do merge
	if len(needMerge) > 0 {
		mergedValue, err := merger.Merge(previousKey, needMerge)
		if err != nil {
			return err
		}
		if err := c.add(previousKey, mergedValue); err != nil {
			return err
		}
	}
	// if has store builder opened, need close it
	if c.state.builder != nil {
		if err := c.finishCompactionOutputFile(); err != nil {
			return err
		}
	}
	return nil
}

// installCompactionResults installs compactions results.
// 1. mark input files is deletion which compaction job picked.
// 2. add output files to up level.
// 3. commit edit log for manifest.
func (c *compactJob) installCompactionResults() {
	// marks compaction input files for deletion
	c.state.compaction.MarkInputDeletes()
	// adds compaction outputs
	level := c.state.compaction.GetLevel()
	for _, output := range c.state.outputs {
		c.state.compaction.AddFile(level+1, output)
	}
	c.family.commitEditLog(c.state.compaction.GetEditLog())
}

// add adds new k/v pair into new store build,
// if store builder is nil need create a new store builder,
// if file size > max file limit, closes current builder.
func (c *compactJob) add(key uint32, value []byte) error {
	if len(value) == 0 {
		return nil
	}
	// generates output file number and creates store build if necessary
	if c.state.builder == nil {
		if err := c.openCompactionOutputFile(); err != nil {
			return err
		}
	}
	// add key/value into store builder
	if err := c.state.builder.Add(key, value); err != nil {
		return err
	}
	// close current store build's file if it is big enough
	if c.state.builder.Size() >= c.state.maxFileSize {
		if err := c.finishCompactionOutputFile(); err != nil {
			return err
		}
	}
	return nil
}

// makeInputIterator makes a merged iterator by compaction pick input files
func (c *compactJob) makeInputIterator() (table.Iterator, error) {
	var its []table.Iterator
	for which := 0; which < 2; which++ {
		files := c.state.compaction.GetInputs()[which]
		if len(files) > 0 {
			for _, fileMeta := range files {
				reader, err := c.state.snapshot.GetReader(fileMeta.GetFileNumber())
				if err != nil {
					return nil, err
				}
				its = append(its, reader.Iterator())
			}
		}
	}
	return table.NewMergedIterator(its), nil
}

// openCompactionOutputFile opens a new compaction store build, and adds the file number into pending output
func (c *compactJob) openCompactionOutputFile() error {
	//TODO add lock
	builder, err := c.family.newTableBuilder()
	if err != nil {
		return err
	}
	fileNumber := builder.FileNumber()
	c.family.addPendingOutput(fileNumber)
	c.state.currentFileNumber = fileNumber
	c.state.builder = builder
	return nil
}

// finishCompactionOutputFile closes current store builder, then generates a new file into edit log
func (c *compactJob) finishCompactionOutputFile() (err error) {
	builder := c.state.builder
	// finally need cleanup store build if no error
	defer func() {
		if err == nil {
			c.state.builder = nil
		}
	}()
	if builder == nil {
		return errors.New("store build is nil")
	}
	if builder.Count() == 0 {
		// if no data after compact
		return err
	}
	if err = builder.Close(); err != nil {
		return fmt.Errorf("close table builder error when compaction job, error:%w", err)
	}
	fileMeta := version.NewFileMeta(builder.FileNumber(), builder.MinKey(), builder.MaxKey(), builder.Size())
	c.state.addOutputFile(fileMeta)
	return err
}

// cleanupCompaction cleanups the compaction context, such as remove pending output files etc.
func (c *compactJob) cleanupCompaction() {
	if c.state.builder != nil {
		if err := c.state.builder.Abandon(); err != nil {
			kvLogger.Warn("abandon store build error when do compact job",
				logger.String("family", c.family.familyInfo()),
				logger.Int64("file", c.state.currentFileNumber.Int64()))
		}
		c.family.removePendingOutput(c.state.currentFileNumber)
	}
	for _, output := range c.state.outputs {
		c.family.removePendingOutput(output.GetFileNumber())
	}
}
