package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FileName(t *testing.T) {
	assert.Equal(t, "000001.sst", Table(1))
	assert.Equal(t, "1234567891011.sst", Table(1234567891011))

	assert.Equal(t, "MANIFEST-000012", ManifestFileName(12))
	assert.Equal(t, "MANIFEST-123456789", ManifestFileName(123456789))
	assert.Equal(t, "CURRENT", current())
}

func Test_ParseFileName(t *testing.T) {
	assert.Nil(t, ParseFileName("xxx.tt"))
	assert.Nil(t, ParseFileName("aaa.sst"))
	fileDesc := ParseFileName("000001.sst")
	assert.Equal(t, TypeTable, fileDesc.FileType)
	assert.Equal(t, int64(1), fileDesc.FileNumber.Int64())
}
