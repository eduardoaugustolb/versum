package biblecorpus

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVerifyAcceptsAGeneratedCanonicalCorpus(t *testing.T) {
	root := writeCanonicalCorpus(t)

	if err := Verify(root); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
}

func TestVerifyRejectsABookWhoseContentDoesNotMatchItsHash(t *testing.T) {
	root := writeCanonicalCorpus(t)
	path := filepath.Join(root, "books", "01-b01.json")
	if err := os.WriteFile(path, []byte(`{"tampered":true}\n`), 0o644); err != nil {
		t.Fatalf("tamper book: %v", err)
	}

	if err := Verify(root); err == nil {
		t.Fatal("Verify() error = nil, want error for a tampered book")
	}
}

func writeCanonicalCorpus(t *testing.T) string {
	t.Helper()

	root := filepath.Join(t.TempDir(), "corpus", "v1")
	if _, err := Generate(writeRawCorpus(t, expectedBookCount), root); err != nil {
		t.Fatalf("generate canonical corpus: %v", err)
	}
	return root
}
