package placetypes

import (
	"testing"
)

func TestIsAncestor(t *testing.T) {

	a, err := GetPlacetypeByName("locality")

	if err != nil {
		t.Fatalf("Failed to get locality placetype, %v", err)
	}

	b, err := GetPlacetypeByName("country")

	if err != nil {
		t.Fatalf("Failed to get country placetype, %v", err)
	}

	if !IsAncestor(a, b) {
		t.Fatalf("Expected b to be ancestor of a")
	}
}

func TestIsDescendant(t *testing.T) {

	a, err := GetPlacetypeByName("country")

	if err != nil {
		t.Fatalf("Failed to get country placetype, %v", err)
	}

	b, err := GetPlacetypeByName("microhood")

	if err != nil {
		t.Fatalf("Failed to get microhood placetype, %v", err)
	}

	if !IsDescendant(a, b) {
		t.Fatalf("Expected b to be descendant of a")
	}
}
