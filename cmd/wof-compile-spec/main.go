package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
)

func main() {

	iterator_uri := flag.String("iterator-uri", "directory://", "")
	iterator_source := flag.String("iterator-source", "/usr/local/whosonfirst/whosonfirst-placetypes/placetypes", "")

	flag.Parse()

	ctx := context.Background()

	mu := new(sync.RWMutex)
	parent_map := new(sync.Map)

	wof_placetypes := make([]*placetypes.WOFPlacetypeRecord, 0)

	iter, err := iterate.NewIterator(ctx, *iterator_uri)

	if err != nil {
		log.Fatalf("Failed to create iterator, %v", err)
	}
	
	for rec, err := range iter.Iterate(ctx, *iterator_source) {

		if err != nil {
			log.Fatalf("Iterator yielded error, %v", err)
		}

		defer rec.Body.Close()

		if filepath.Ext(rec.Path) != ".json" {
			continue
		}

		var pt *placetypes.WOFPlacetypeRecord

		dec := json.NewDecoder(rec.Body)
		err := dec.Decode(&pt)

		if err != nil {
			log.Fatalf("Failed to decode %s, %v", rec.Path, err)
		}

		parent_map.Store(pt.Name, pt.Id)

		mu.Lock()
		defer mu.Unlock()

		wof_placetypes = append(wof_placetypes, pt)
	}

	// START OF... not sure...

	spec := make(map[string]*placetypes.WOFPlacetype)

	for _, pt := range wof_placetypes {

		// Legacy stuff, oh well...
		str_id := strconv.FormatInt(pt.Id, 10)

		parents := pt.Parent
		parent_ids := make([]int64, len(parents))

		for idx, p := range parents {

			p_id, ok := parent_map.Load(p)

			if !ok {
				log.Fatalf("Unable to derive parent ID for %s", p)
			}

			parent_ids[idx] = p_id.(int64)
		}

		spec[str_id] = &placetypes.WOFPlacetype{
			Id:     pt.Id,
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
