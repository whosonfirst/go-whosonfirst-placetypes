package draw

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
)

// DrawPlacetypesGraphToWriter will derive a Graphviz "dot" representation of 'spec' and render it as an image
// written to 'path'. Valid image formats (derived from 'path' 's file extension) are: jpeg (jpg) and png.
func DrawPlacetypesGraphToFile(spec *placetypes.WOFPlacetypeSpecification, path string) error {

	ext := filepath.Ext(path)
	format := strings.TrimLeft(ext, ".")

	switch format {
	case "jpg", "jpeg", "png":
		// pass
	default:
		return fmt.Errorf("Invalid or unsupported format, %s", format)
	}

	wr, err := os.Create(path)

	if err != nil {
		return fmt.Errorf("Failed to open %s for writing, %w", path, err)
	}

	err = DrawPlacetypesGraphToWriter(spec, format, wr)

	if err != nil {
		return fmt.Errorf("Failed to write placetypes graph,%w", err)
	}

	err = wr.Close()

	if err != nil {
		return fmt.Errorf("Failed to close placetypes graph,%w", err)
	}

	return nil
}

// DrawPlacetypesGraphToWriter will derive a Graphviz "dot" representation of 'spec' and render it as an image
// written to 'wr'. Valid image formats (defined by the 'format' argument) are: jpeg (jpg) and png.
func DrawPlacetypesGraphToWriter(spec *placetypes.WOFPlacetypeSpecification, format string, wr io.Writer) error {

	switch format {
	case "jpg", "jpeg", "png":
		// pass
	default:
		return fmt.Errorf("Invalid or unsupported format, %s", format)
	}

	im, err := DrawPlacetypesGraph(spec)

	if err != nil {
		return fmt.Errorf("Failed to draw placetypes graph, %w", err)
	}

	switch format {
	case "jpg", "jpeg":

		opts := &jpeg.Options{
			Quality: 90,
		}

		err = jpeg.Encode(wr, im, opts)

	case "png":
		err = png.Encode(wr, im)
	default:
		// pass
	}

	if err != nil {
		return fmt.Errorf("Failed to encode image, %w", err)
	}

	return nil
}

// DrawPlacetypesGraph will derive a Graphviz "dot" representation of 'spec' and render it as an `image.Image` instance.
func DrawPlacetypesGraph(spec *placetypes.WOFPlacetypeSpecification) (image.Image, error) {

	var buf bytes.Buffer
	buf_wr := bufio.NewWriter(&buf)

	err := spec.PlacetypesToGraphviz(buf_wr)

	if err != nil {
		return nil, fmt.Errorf("Failed to generate graphviz, %v", err)
	}

	buf_wr.Flush()

	gv := graphviz.New()

	graph, err := graphviz.ParseBytes(buf.Bytes())

	if err != nil {
		return nil, fmt.Errorf("Failed to parse graphviz data, %v", err)
	}

	im, err := gv.RenderImage(graph)

	if err != nil {
		return nil, fmt.Errorf("Failed to render graphviz data, %v", err)
	}

	return im, nil
}
