package kv

import (
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/logger"
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
	// GetTargetFamily returns the target family based on source family name
	GetTargetFamily(sourceFamilyName string) Family
}

// needRollup returns if need rollup source family data
func (f *family) needRollup() bool {
	if f.rolluping.Load() {
		// has background rollup job running
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
	if rollupFilesLen >= threshold {
		kvLogger.Info("need to rollup level0 files", logger.String("family", f.familyInfo()),
			logger.Any("numOfFiles", rollupFilesLen), logger.Any("threshold", f.option.RollupThreshold))
		return true
	}
	//FIXME need add time threshold????
	return false
}

// rollup does rollup in source family, need trigger target family does rollup compact job
func (f *family) rollup() {
	// if has background rollup job running, return it.
	if f.rolluping.CAS(false, true) {
		defer func() {
			// clean up unused files, maybe some file not used
			f.deleteObsoleteFiles()
			f.rolluping.Store(false)
		}()

		rollupFiles := f.familyVersion.GetLiveRollupFiles()
		if len(rollupFiles) == 0 {
			return
		}
		var interval timeutil.Interval
		var sourceFiles []table.FileNumber
		for file, i := range rollupFiles {
			// only allow ont target rollup interval
			interval = i
			sourceFiles = append(sourceFiles, file)
		}

		// do rollup job in target family
		rollup, ok := f.store.getRollup(interval)
		if !ok {
			kvLogger.Warn("skip rollup because cannot get target rollup",
				logger.String("family", f.familyInfo()),
				logger.Int64("interval", interval.Int64()))
			return
		}
		editLog := version.NewEditLog(f.ID())
		targetFamily := rollup.GetTargetFamily(f.name)

		if err := targetFamily.doRollupWork(f, rollup, sourceFiles); err != nil {
			kvLogger.Error("do rollup work fail",
				logger.String("family", f.familyInfo()),
				logger.Int64("interval", interval.Int64()),
				logger.Any("files", sourceFiles))
			return
		}

		// after rollup job successfully, need add delete rollup file edit log
		for _, file := range sourceFiles {
			editLog.Add(version.CreateDeleteRollupFile(file))
		}

		// finally need commit edit log
		f.commitEditLog(editLog)
	}
}

// doRollupWork does rollup work in target family,
// 1. reads data from source family
// 2. merge these
// 3. finally builds new sst files in target family
func (f *family) doRollupWork(sourceFamily Family, rollup Rollup, sourceFiles []table.FileNumber) (err error) {
	if len(sourceFiles) == 0 {
		return
	}
	targetFiles := make(map[table.FileNumber]struct{})
	for _, file := range sourceFiles {
		targetFiles[file] = struct{}{}
	}
	referenceFiles := f.familyVersion.GetLiveReferenceFiles()
	files, ok := referenceFiles[sourceFamily.ID()]
	if ok {
		for _, file := range files {
			_, exist := targetFiles[file]
			if exist {
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
		return
	}

	snapshot := sourceFamily.GetSnapshot()
	defer func() {
		snapshot.Close()
	}()
	compaction := version.NewCompaction(f.ID(), -1, nil, nil)

	compactionState := newCompactionState(f.maxFileSize, snapshot, compaction)
	compactJob := newCompactJobFunc(f, compactionState, rollup)
	if err := compactJob.Run(); err != nil {
		return err
	}
	return nil
}
