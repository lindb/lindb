package kv

import "github.com/lindb/lindb/pkg/logger"

const dummy = ""
const RollupContext = "RollupContext"
const defaultMaxFileSize = int32(256 * 1024 * 1024)
const defaultCompactThreshold = 4
const defaultRollupThreshold = 3

var defaultCompactCheckInterval = 60
var kvLogger = logger.GetLogger("kv", "store")
