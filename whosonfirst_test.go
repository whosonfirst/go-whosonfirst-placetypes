package placetypes

import (
	"context"
	"testing"
)

func TestWhosOnFirstDefinition(t *testing.T) {

	ctx := context.Background()

	uri := "whosonfirst://"

	d, err := NewDefinition(ctx, uri)

	if err != nil {
		t.Fatalf("Failed to create definition for '%s', %v", uri, err)
	}

	if d.Property() != "wof:placetype" {
		t.Fatalf("Unexpected placetype property, %s", d.Property())
	}

	if d.URI() != uri {
		t.Fatalf("Unexpected URI, %s", d.URI())
	}
}
		
