package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
)

func main() {

	iterator_uri := flag.String("iterator-uri", "directory://", "")
	iterator_source := flag.String("iterator-source", "/usr/local/whosonfirst/whosonfirst-placetypes/placetypes", "")

	flag.Parse()

	ctx := context.Background()

	mu := new(sync.RWMutex)
	parent_map := new(sync.Map)

	wof_placetypes := make([]*placetypes.WOFPlacetype, 0)

	iter_cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {

		if filepath.Ext(path) != ".json" {
			return nil
		}

		var pt *placetypes.WOFPlacetype

		dec := json.NewDecoder(r)
		err := dec.Decode(&pt)

		if err != nil {
			return fmt.Errorf("Failed to decode %s, %w", path, err)
		}

		parent_map.Store(pt.Name, pt.Id)

		mu.Lock()
		defer mu.Unlock()

		log.Printf("Add %s (%d)\n", path, pt.Id)

		wof_placetypes = append(wof_placetypes, pt)
		return nil
	}

	iter, err := iterator.NewIterator(ctx, *iterator_uri, iter_cb)

	if err != nil {
		log.Fatalf("Failed to create iterator, %v", err)
	}

	err = iter.IterateURIs(ctx, *iterator_source)

	if err != nil {
		log.Fatalf("Failed to iterate URIs, %v", err)
	}

	// START OF... not sure...

	type spec_pt struct {
		Role   string              `json:"role"`
		Name   string              `json:"name"`
		Parent []int64             `json:"parent"`
		Names  map[string][]string `json:"names"`
	}

	spec := make(map[int64]spec_pt)

	for _, pt := range wof_placetypes {

		id := pt.Id

		parents := pt.Parent
		parent_ids := make([]int64, len(parents))

		for idx, p := range parents {

			p_id, ok := parent_map.Load(p)

			if !ok {
				log.Fatalf("Unable to derive parent ID for %s", p)
			}

			parent_ids[idx] = p_id.(int64)
		}

		spec[id] = spec_pt{
			Role:   pt.Role,
			Name:   pt.Name,
			Parent: parent_ids,
		}
	}

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(spec)

	if err != nil {
		log.Fatalf("Failed to encode spec, %v", err)
	}
}
