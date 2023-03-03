package placetypes

import (
	"testing"
)

func TestIsCorePlacetype(t *testing.T) {

	is_core := []string{
		"region",
		"custom",
		"disputed",
		"localadmin",
	}

	not_core := []string{
		"airport",
		"gate",
	}

	for _, n := range is_core {

		if !isCorePlacetype(n) {
			t.Fatalf("Expected %s to be core placetype", n)
		}
	}

	for _, n := range not_core {

		if isCorePlacetype(n) {
			t.Fatalf("%s not expected to be core placetype", n)
		}
	}

}

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
