package ports

import "context"

// CatalogVersionRepository persists the version of the published corpus.
type CatalogVersionRepository interface {
	Record(ctx context.Context, corpusSHA256 string) error
}
