package commands

import (
	"context"
)

type PublishCatalogVersionInput struct {
	CorpusSHA256 string
}

type CatalogVersionWriter interface {
	Record(ctx context.Context, corpusSHA256 string) error
}

type PublishCatalogVersion struct {
	writer CatalogVersionWriter
}

func NewPublishCatalogVersion(writer CatalogVersionWriter) PublishCatalogVersion {
	return PublishCatalogVersion{
		writer: writer,
	}
}

func (uc *PublishCatalogVersion) Execute(ctx context.Context, input PublishCatalogVersionInput) error {
	return uc.writer.Record(ctx, input.CorpusSHA256)
}
