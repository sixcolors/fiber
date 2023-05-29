package negotiator

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2/utils"
)

func Test_PreferredMediaTypes(t *testing.T) {
	t.Parallel()

	utils.AssertEqual(t, []string{"application/json"}, PreferredMediaTypes("", "application/json"))
	utils.AssertEqual(t, []string{"application/json", "text/html"}, PreferredMediaTypes("text/html, application/json", "application/json", "text/html"))

	utils.AssertEqual(t, []string{"*/*"}, PreferredMediaTypes("*/*"))
	utils.AssertEqual(t, []string{}, PreferredMediaTypes("application/json;q=0", "application/json"), "q=0 should be ignored")
	utils.AssertEqual(t, []string{"text/plain"}, PreferredMediaTypes("text/plain, application/json;q=0.5, text/html, text/xml, text/yaml, text/javascript, text/csv, text/css, text/rtf, text/markdown, application/octet-stream;q=0.2, */*;q=0.1", "text/plain"))

	accept := "text/html, application/*;q=0.2, image/jpeg;q=0.8"
	utils.AssertEqual(t, []string{"text/html", "image/jpeg", "application/*"}, PreferredMediaTypes(accept))

	provided := []string{"text/html", "text/plain", "application/json"}
	utils.AssertEqual(t, []string{"text/html", "application/json"}, PreferredMediaTypes(accept, provided...))

	// check wildcard
	utils.AssertEqual(t, []string{"application/json"}, PreferredMediaTypes("*/*", "application/json"))
	utils.AssertEqual(t, []string{"text/plain"}, PreferredMediaTypes("*/*", "text/plain"))
	utils.AssertEqual(t, []string{"application/xml"}, PreferredMediaTypes("text/html, */*", "application/xml"))

}

func Benchmark_PerferedMediaTypes(b *testing.B) {
	accepts := []string{
		"text/html,application/xhtml+xml,application/xml;q=0.9",
		"text/html, application/*;q=0.2, image/jpeg;q=0.8, text/plain, application/json;q=0, application/octet-stream;q=0.2, */*;q=0.1",
		"text/html, application/*;q=0.2, image/jpeg;q=0.8, text/plain, application/json, text/xhtml, text/xml, text/yaml, text/javascript, text/csv, text/css, text/rtf, text/markdown, application/octet-stream;q=0.2, */*;q=0.1",
	}
	for i := 0; i < len(accepts); i++ {
		b.Run(fmt.Sprintf("run-%#v", accepts[i]), func(bb *testing.B) {
			bb.ReportAllocs()
			bb.ResetTimer()

			for n := 0; n < bb.N; n++ {
				PreferredMediaTypes(accepts[n%3])
			}
		})
	}
}

func Test_parseMediaType(t *testing.T) {
	t.Parallel()

	accept := ""
	mediatype := parseMediaType(accept, 0)

	if mediatype != nil {
		t.Fatalf("Expected nil, got %v", mediatype)
	}

	accept = "text/html"
	mediatype = parseMediaType(accept, 0)

	utils.AssertEqual(t, "text", mediatype.Type)
	utils.AssertEqual(t, "html", mediatype.Subtype)
	utils.AssertEqual(t, 1.0, mediatype.Q)
	utils.AssertEqual(t, 0, len(mediatype.Params))

	accept = "text/html;q=0.8"
	mediatype = parseMediaType(accept, 0)

	utils.AssertEqual(t, "text", mediatype.Type)
	utils.AssertEqual(t, "html", mediatype.Subtype)
	utils.AssertEqual(t, 0.8, mediatype.Q)
	utils.AssertEqual(t, 0, len(mediatype.Params))

	accept = "text/html;foo=bar"
	mediatype = parseMediaType(accept, 0)

	utils.AssertEqual(t, "text", mediatype.Type)
	utils.AssertEqual(t, "html", mediatype.Subtype)
	utils.AssertEqual(t, 1.0, mediatype.Q)
	utils.AssertEqual(t, 1, len(mediatype.Params))
	utils.AssertEqual(t, "bar", mediatype.Params["foo"])
}

func Test_getMediaTypePriority(t *testing.T) {
	t.Parallel()

	accept := []mediaType{
		{
			Type:    "text",
			Subtype: "html",
			Q:       1.0,
			Params:  map[string]string{},
		},
		{
			Type:    "text",
			Subtype: "*",
			Q:       0.8,
			Params:  map[string]string{},
		},
		{
			Type:    "*",
			Subtype: "*",
			Q:       0.1,
			Params:  map[string]string{},
		},
	}

	mediatype := getMediaTypePriority("text/html", accept, 0)

	utils.AssertEqual(t, "text", mediatype.Type)
	utils.AssertEqual(t, "html", mediatype.Subtype)
	utils.AssertEqual(t, 1.0, mediatype.Q)
	utils.AssertEqual(t, 0, len(mediatype.Params))
}

func Test_specify(t *testing.T) {
	t.Parallel()

	spec := mediaType{
		Type:    "text",
		Subtype: "html",
		Q:       1.0,
		Params:  map[string]string{},
	}

	mediatype := specify("text/html", &spec, 0)

	utils.AssertEqual(t, "text", mediatype.Type)
	utils.AssertEqual(t, "html", mediatype.Subtype)
	utils.AssertEqual(t, 1.0, mediatype.Q)
	utils.AssertEqual(t, 0, len(mediatype.Params))
}
