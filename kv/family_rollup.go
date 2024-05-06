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
	"math/rand"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/lindb/common/pkg/logger"
	commontimeutil "github.com/lindb/common/pkg/timeutil"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/timeutil"
)

//go:generate mockgen -source ./family_rollup.go -destination=./family_rollup_mock.go -package kv

// Rollup represents rollup relation(source store/family => target store/family)
type Rollup interface {
	// GetTimestamp returns the timestamp based on source family and source slot
	GetTimestamp(slot uint16) int64
	// IntervalRatio return interval ratio = target interval/source interval
	IntervalRatio() uint16
	// CalcSlot calculates the target slot based on source timestamp
	CalcSlot(timestamp int64) uint16
	// BaseSlot returns base slot by source family time/target interval.
	BaseSlot() uint16
}

// rollup implements Rollup interface.
type rollup struct {
	source, target           timeutil.Interval
	sourceFTime, targetFTime int64
}

func newRollup(source, target timeutil.Interval, sourceFTime, targetFTime int64) Rollup {
	return &rollup{
		source:      source,
		target:      target,
		sourceFTime: sourceFTime,
		targetFTime: targetFTime,
	}
}

func (r *rollup) GetTimestamp(slot uint16) int64 {
	return r.sourceFTime + int64(slot)*r.source.Int64()
}

func (r *rollup) IntervalRatio() uint16 {
	return uint16(r.target / r.source)
}

func (r *rollup) CalcSlot(timestamp int64) uint16 {
	return uint16(r.target.Calculator().CalcSlot(timestamp, r.targetFTime, r.target.Int64()))
}

func (r *rollup) BaseSlot() uint16 {
	return r.CalcSlot(r.sourceFTime)
}

// needRollup checks if it needs rollup source family data.
func (f *family) needRollup() bool {
	if f.rolluping.Load() {
		// has background rollup job running
		kvLogger.Info("rollup job is running", logger.String("family", f.familyInfo()))
		return false
	}
	rollupTargetStores := f.store.Option().Rollup
	if len(rollupTargetStores) == 0 {
		// not set rollup
		return false
	}
	rollupFiles := f.familyVersion.GetLiveRollupFiles()
	rollupFilesLen := len(rollupFiles)
	if rollupFilesLen == 0 {
		// no files need to rollup
		return false
	}
	threshold := f.option.RollupThreshold
	if threshold <= 0 {
		threshold = defaultRollupThreshold
	}
	kvLogger.Info("check file threshold if need to rollup level0 files", logger.String("family", f.familyInfo()),
		logger.Any("numOfFiles", rollupFilesLen), logger.Any("threshold", threshold))
	if rollupFilesLen >= threshold {
		return true
	}
	var targetIntervals []timeutil.Interval
	for _, rollupFile := range rollupFiles {
		targetIntervals = append(targetIntervals, rollupFile...)
	}
	sort.Slice(targetIntervals, func(i, j int) bool {
		return targetIntervals[i] < targetIntervals[j]
	})
	targetInterval := targetIntervals[0]
	now := commontimeutil.Now()
	diff := now - f.lastRollupTime.Load()
	timeThreshold := int64(targetInterval) + rand.Int63n(180000)
	kvLogger.Info("check time threshold if need to rollup level0 files",
		logger.String("family", f.familyInfo()),
		logger.String("now", commontimeutil.FormatTimestamp(now, commontimeutil.DataTimeFormat2)),
		logger.String("lastRollupTime", commontimeutil.FormatTimestamp(f.lastRollupTime.Load(), commontimeutil.DataTimeFormat2)),
		logger.Int64("diff", diff), logger.Int64("threshold", timeThreshold))
	return diff > timeThreshold
}

// rollup does rollup in source family, need trigger target family does rollup compact job
func (f *family) rollup() {
	// check if it has background rollup job running already,
	// has rollup job, return it, else do rollup job.
	if f.rolluping.CompareAndSwap(false, true) {
		f.condition.Add(1)
		go func() {
			defer func() {
				// clean up unused files, maybe some file not used
				f.deleteObsoleteFiles()
				f.condition.Done()
				f.rolluping.Store(false)
				f.lastRollupTime.Store(commontimeutil.Now())
			}()

			rollupFiles := f.familyVersion.GetLiveRollupFiles()
			if len(rollupFiles) == 0 {
				return
			}
			rollupMap := make(map[timeutil.Interval][]table.FileNumber)
			for file, intervals := range rollupFiles {
				for _, interval := range intervals {
					rollupMap[interval] = append(rollupMap[interval], file)
				}
			}

			editLog := version.NewEditLog(f.ID())
			sourceInterval := f.store.Option().Source
			calc := sourceInterval.Calculator()
			storeName := f.store.Name()
			_, segmentName := filepath.Split(storeName)
			segmentTime, err := calc.ParseSegmentTime(segmentName)
			if err != nil {
				kvLogger.Error("parse segment time failure, when do rollup job",
					logger.String("family", f.familyInfo()),
					logger.Error(err))
				return
			}
			fTime, err := strconv.Atoi(f.Name())
			if err != nil {
				kvLogger.Error("parse family time failure, when do rollup job",
					logger.String("family", f.familyInfo()),
					logger.Error(err))
				return
			}
			familyStartTime := calc.CalcFamilyStartTime(segmentTime, fTime)
			baseDir := strings.Replace(storeName, filepath.Join(sourceInterval.Type().String(), segmentName), "", 1)
			targetFamiles := make(map[Family][]table.FileNumber)

			for targetInterval, files := range rollupMap {
				segmentName := targetInterval.Calculator().GetSegment(familyStartTime)
				targetStoreName := filepath.Join(baseDir, targetInterval.Type().String(), segmentName)
				targetStore, ok := GetStoreManager().GetStoreByName(targetStoreName)
				// do rollup job in target family
				if !ok {
					// TODO: add metric
					kvLogger.Warn("skip rollup because cannot get target store",
						logger.String("family", f.familyInfo()),
						logger.String("target", targetStoreName),
						logger.String("interval", targetInterval.String()))
					continue
				}
				tSegmentTime := targetInterval.Calculator().CalcSegmentTime(familyStartTime)
				tFamilyTime := targetInterval.Calculator().CalcFamily(familyStartTime, tSegmentTime)
				fSTime := targetInterval.Calculator().CalcFamilyStartTime(tSegmentTime, tFamilyTime)

				// re-use source family option
				targetFamily, err := targetStore.CreateFamily(strconv.Itoa(tFamilyTime), f.option)
				if err != nil {
					kvLogger.Error("create target family failure when do rollup job",
						logger.String("family", f.familyInfo()),
						logger.String("target", segmentName),
						logger.String("interval", targetInterval.String()),
						logger.Error(err))
					continue
				}
				rollup := newRollup(sourceInterval, targetInterval, familyStartTime, fSTime)
				if err := targetFamily.doRollupWork(f, rollup, files); err != nil {
					kvLogger.Error("do rollup work fail",
						logger.String("family", f.familyInfo()),
						logger.String("target", segmentName),
						logger.String("interval", targetInterval.String()),
						logger.Any("files", files))
					continue
				}
				targetFamiles[targetFamily] = files

				// after rollup job successfully, need add delete rollup file edit log
				for _, file := range files {
					editLog.Add(version.CreateDeleteRollupFile(file, targetInterval))
				}
			}

			// finally, need commit edit log
			f.commitEditLog(editLog)

			// clean reference files from target file
			for targetFamily, files := range targetFamiles {
				targetFamily.cleanReferenceFiles(f, files)
			}
		}()
	}
}

// cleanReferenceFiles cleans target family's reference files after delete source family's rollup files.
func (f *family) cleanReferenceFiles(sourceFamily Family, sourceFiles []table.FileNumber) {
	editLog := version.NewEditLog(f.ID())
	_, sourceStore := filepath.Split(sourceFamily.getStore().Name())
	sourceFamilyID := sourceFamily.ID()
	for _, file := range sourceFiles {
		editLog.Add(version.CreateDeleteReferenceFile(sourceStore, sourceFamilyID, file))
	}
	f.commitEditLog(editLog)
}

// doRollupWork does rollup work in target family,
// 1. reads data from source family
// 2. merge these
// 3. finally, builds new sst files in target family
func (f *family) doRollupWork(sourceFamily Family, rollup Rollup, sourceFiles []table.FileNumber) (err error) {
	if len(sourceFiles) == 0 {
		return nil
	}
	targetFiles := make(map[table.FileNumber]struct{})
	for _, file := range sourceFiles {
		targetFiles[file] = struct{}{}
	}
	_, sourceStore := filepath.Split(sourceFamily.getStore().Name())
	sourceFamilyID := sourceFamily.ID()
	referenceFiles := f.familyVersion.GetLiveReferenceFiles(sourceStore)
	if files, ok := referenceFiles[sourceFamilyID]; ok {
		// check if file already rollup
		for _, file := range files {
			if _, exist := targetFiles[file]; exist {
				delete(targetFiles, file)
				kvLogger.Warn("skip rollup for this file, because file already rollup",
					logger.Int64("fileNumber", file.Int64()),
					logger.String("source", sourceFamily.familyInfo()),
					logger.String("target", f.familyInfo()),
				)
			}
		}
	}
	if len(targetFiles) == 0 {
		// if no target files, return it
		return nil
	}

	snapshot := sourceFamily.GetSnapshot()
	defer func() {
		snapshot.Close()
	}()
	v := snapshot.GetCurrent()

	var inputFiles []*version.FileMeta
	var logs []version.Log
	for fileNumber := range targetFiles {
		if fm, ok := v.GetFile(0, fileNumber); ok {
			inputFiles = append(inputFiles, fm)
			logs = append(logs, version.CreateNewReferenceFile(sourceStore, sourceFamilyID, fileNumber))
		}
	}
	compaction := version.NewCompaction(f.ID(), 0, inputFiles, nil)
	// add reference file edit logs
	compaction.AddReferenceFiles(logs)

	compactionState := newCompactionState(f.maxFileSize, snapshot, compaction)
	compactJob := newCompactJobFunc(f, compactionState, rollup)
	if err := compactJob.Run(); err != nil {
		return err
	}
	return nil
}
