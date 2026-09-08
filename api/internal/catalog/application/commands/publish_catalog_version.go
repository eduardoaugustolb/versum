package commands

import (
	"context"

	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/ports"
)

type PublishCatalogVersionInput struct {
	CorpusSHA256 string
}

type PublishCatalogVersion struct {
	repository ports.CatalogVersionRepository
}

func NewPublishCatalogVersion(repository ports.CatalogVersionRepository) PublishCatalogVersion {
	return PublishCatalogVersion{
		repository: repository,
	}
}

func (uc *PublishCatalogVersion) Execute(ctx context.Context, input PublishCatalogVersionInput) error {
	return uc.repository.Record(ctx, input.CorpusSHA256)
}
