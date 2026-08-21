package commands_test

import (
	"context"
	"errors"
	"testing"

	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/commands"
)

type FakeWriter struct {
	err       error
	gotSHA256 string
}

func (w *FakeWriter) Record(ctx context.Context, corpusSHA string) error {
	w.gotSHA256 = corpusSHA
	return w.err
}

func TestPublishCatalogVersion(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
		err       error
	}{
		{
			name:      "success on registry the version",
			wantError: false,
		},
		{
			name:      "returns error on registry the version",
			wantError: true,
			err:       errors.New("testing"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			writer := FakeWriter{
				err: tc.err,
			}
			publishCatalogVersion := commands.NewPublishCatalogVersion(&writer)
			corpusSHA := "testing"
			err := publishCatalogVersion.Execute(t.Context(), commands.PublishCatalogVersionInput{CorpusSHA256: corpusSHA})
			if (err != nil) != tc.wantError {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.wantError && !errors.Is(err, tc.err) {
				t.Fatalf("expected error %v, got %v", tc.err, err)
			}

			if writer.gotSHA256 != corpusSHA {
				t.Fatalf("expected gotSHA256 %v, got %v", corpusSHA, writer.gotSHA256)
			}
		})
	}
}
