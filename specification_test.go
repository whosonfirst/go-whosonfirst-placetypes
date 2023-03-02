package placetypes

import (
	"fmt"
	"testing"

	"github.com/dominikbraun/graph"
)

func TestGraphPlacetypes(t *testing.T) {

	spec, err := DefaultWOFPlacetypeSpecification()

	if err != nil {
		t.Fatalf("Failed to load default WOF placetype specification, %v", err)
	}

	gr, err := spec.GraphPlacetypes()

	if err != nil {
		t.Fatalf("Failed to graph placetypes, %v", err)
	}

	ordered, err := graph.TopologicalSort(gr)

	if err != nil {
		t.Fatalf("Failed to sort placetypes graph, %v", err)
	}

	fmt.Println(ordered)
}

func TestAppendSpecification(t *testing.T) {

	spec, err := DefaultWOFPlacetypeSpecification()

	if err != nil {
		t.Fatalf("Failed to load default WOF placetype specification, %v", err)
	}

	parents := []int64{
		102312307, // country
	}

	pt := WOFPlacetype{
		Id:     1,
		Name:   "map",
		Role:   "optional",
		Parent: parents,
	}

	err = spec.AppendPlacetype(pt)

	if err != nil {
		t.Fatalf("Failed to append placetype, %v", err)
	}

	_, err = spec.GetPlacetypeById(pt.Id)

	if err != nil {
		t.Fatalf("Failed to get placetype by ID, %v", err)
	}

	_, err = spec.GetPlacetypeByName(pt.Name)

	if err != nil {
		t.Fatalf("Failed to get placetype by name, %v", err)
	}

}
