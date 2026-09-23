package httpapi_test

import (
	"context"
	"net/http"
	"testing"

	miniredis "github.com/alicebob/miniredis/v2"
	redisadapter "github.com/eduardoaugustolb/versum/api/internal/cache/redis"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/queries"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
	"github.com/eduardoaugustolb/versum/api/internal/health"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi"
	redisclient "github.com/redis/go-redis/v9"
)

type noopCatalogRepository struct{}

func (noopCatalogRepository) ReplaceBook(context.Context, domain.Book, []domain.Verse) error {
	return nil
}
func (noopCatalogRepository) ListBooks(context.Context) ([]domain.Book, error) {
	return nil, nil
}
func (noopCatalogRepository) FindChapter(context.Context, string, int) (domain.Chapter, error) {
	return domain.Chapter{}, domain.ErrChapterNotFound
}

func router(t *testing.T, br booksRepository, cr chaptersRepository) http.Handler {
	t.Helper()

	redisServer := miniredis.RunT(t)
	redisClient := redisclient.NewClient(&redisclient.Options{Addr: redisServer.Addr()})
	t.Cleanup(func() { _ = redisClient.Close() })

	return httpapi.NewHandler(httpapi.Dependencies{
		Health: health.CheckHealth{},
		Catalog: httpapi.CatalogDependencies{
			ListBooks:  queries.NewListBooks(br),
			GetChapter: queries.NewGetChapter(cr),
		},
		Cache: redisadapter.NewRedisCache(redisClient),
	})
}
