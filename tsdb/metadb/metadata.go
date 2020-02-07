package metadb

import (
	"context"

	"github.com/lindb/lindb/kv"
)

// metadata implements Metadata interface
type metadata struct {
	metadataDatabase MetadataDatabase
	tagMetadata      TagMetadata
}

// NewMetadata creates a metadata
func NewMetadata(ctx context.Context, name, parent string, tagFamily kv.Family) (Metadata, error) {
	db, err := NewMetadataDatabase(ctx, name, parent)
	if err != nil {
		return nil, err
	}
	return &metadata{
		metadataDatabase: db,
		tagMetadata:      NewTagMetadata(tagFamily),
	}, nil
}

// MetadataDatabase returns the metric level metadata
func (m *metadata) MetadataDatabase() MetadataDatabase {
	return m.metadataDatabase
}

// TagMetadata returns the tag metadata
func (m *metadata) TagMetadata() TagMetadata {
	return m.tagMetadata
}
