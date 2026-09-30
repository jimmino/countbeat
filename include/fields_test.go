package include

import (
	"bytes"
	"os"
	"testing"

	"github.com/elastic/beats/v7/libbeat/asset"
)

// The compressed bytes in fields.go vary with the Go version that generated
// them, so compare the decoded content rather than the file itself.
func TestFieldsGoMatchesFieldsYml(t *testing.T) {
	want, err := os.ReadFile("../fields.yml")
	if err != nil {
		t.Fatal(err)
	}
	got, err := asset.DecodeData(AssetFieldsYml())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Error("include/fields.go does not match fields.yml; run 'make update'")
	}
}
