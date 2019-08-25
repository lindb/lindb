package version

// Compaction represents the compaction job context
type Compaction struct {
	level int

	inputs        [][]*FileMeta
	levelInputs   []*FileMeta
	levelUpInputs []*FileMeta

	editLog *EditLog
}

// NewCompaction create a compaction job context
func NewCompaction(familyID, level int, levelInputs, levelUpInputs []*FileMeta) *Compaction {
	return &Compaction{
		level:         level,
		inputs:        [][]*FileMeta{levelInputs, levelUpInputs},
		levelInputs:   levelInputs,
		levelUpInputs: levelUpInputs,
		editLog:       NewEditLog(familyID),
	}
}

// IsTrivialMove returns a trivial compaction that can be implemented by just
// moving a single input file to the next level (no merging or splitting).
// returns true: can just moving file to the next level
func (c *Compaction) IsTrivialMove() bool {
	return len(c.levelInputs) == 1 && len(c.levelUpInputs) == 0
}

// GetLevelFiles returns low level files
func (c *Compaction) GetLevelFiles() []*FileMeta {
	return c.levelInputs
}

// DeleteFile deletes a old file which compaction input file
func (c *Compaction) DeleteFile(level int, fileNumber int64) {
	c.editLog.Add(NewDeleteFile(int32(level), fileNumber))
}

// AddFile adds a new file which compaction output file
func (c *Compaction) AddFile(level int, file *FileMeta) {
	c.editLog.Add(CreateNewFile(int32(level), file))
}

// MarkInputDeletes marks all inputs of this compaction as deletion, adds log to version edit.
func (c *Compaction) MarkInputDeletes() {
	for _, input := range c.levelInputs {
		c.editLog.Add(NewDeleteFile(int32(c.level), input.fileNumber))
	}
	for _, upInput := range c.levelUpInputs {
		c.editLog.Add(NewDeleteFile(int32(c.level+1), upInput.fileNumber))
	}
}

// GetLevel returns compaction level
func (c *Compaction) GetLevel() int {
	return c.level
}

// GetEditLog returns edit log
func (c *Compaction) GetEditLog() *EditLog {
	return c.editLog
}

// GetInputs returns all input files for compaction job
func (c *Compaction) GetInputs() [][]*FileMeta {
	return c.inputs
}
