package negotiator

import (
	"sort"
	"strconv"
	"strings"
)

type mediaType struct {
	Type    string
	Subtype string
	Params  map[string]string
	Q       float64
	I       int
	S       int
}

// PreferredMediaTypes returns the preferred media types from an Accept header.
// If no Accept header is given, the default is to return all provided media types ordered by client preference.
//
// The provided media types are ordered by priority from highest to lowest preference.
// The first media type with a quality value > 0 will be the preferred media type.
// If no media types in the provided list have a quality value > 0, then the preferred media type is the first media type in the provided list.
//
// Example:
//
//	PreferredMediaTypes("text/html, application/json", "application/json", "text/html")
//	// -> []string{"application/json", "text/html"}
//
//	PreferredMediaTypes("text/html;q=0.2, application/json;q=0.8", "application/json", "text/html")
//	// -> []string{"application/json", "text/html"}
//
//	PreferredMediaTypes("text/html, text/plain, */*", "application/json")
//	// -> []string{"application/json"}
func PreferredMediaTypes(accept string, provided ...string) []string {
	if accept == "" {
		accept = "*/*"
	}
	accepts := parseAccept(accept)

	sort.SliceStable(accepts, func(i, j int) bool {
		if accepts[i].Q != accepts[j].Q {
			return accepts[i].Q > accepts[j].Q
		}
		if accepts[i].S != accepts[j].S {
			return accepts[i].S > accepts[j].S
		}
		return accepts[i].I < accepts[j].I
	})

	if len(provided) == 0 {
		// Sorted list of all types
		types := make([]string, 0, len(accepts))
		for _, mediaType := range accepts {
			if isQuality(&mediaType) {
				types = append(types, getFullType(&mediaType))
			}
		}
		return types
	}

	priorities := []mediaType{}
	for i, typ := range provided {
		if priority := getMediaTypePriority(typ, accepts, i); priority != nil {
			priorities = append(priorities, *priority)
		}
	}

	// Sorted list of accepted types
	types := make([]string, 0, len(priorities))
	for _, priority := range priorities {
		if isQuality(&priority) {
			types = append(types, provided[priority.I])
		}
	}
	return types
}

// parseAccept parses the Accept header and returns a list of media types.
func parseAccept(accept string) []mediaType {
	accepts := splitMediaTypes(accept)
	parsedAccepts := make([]mediaType, 0, len(accepts))

	for i := 0; i < len(accepts); i++ {
		mediaType := parseMediaType(strings.TrimSpace(accepts[i]), i)
		if mediaType != nil {
			parsedAccepts = append(parsedAccepts, *mediaType)
		}
	}

	return parsedAccepts
}

// parseMediaType parses a media type from the Accept header.
func parseMediaType(str string, i int) *mediaType {
	parts := strings.Split(str, ";")

	typeAndSubtype := strings.SplitN(parts[0], "/", 2)
	if len(typeAndSubtype) != 2 {
		return nil
	}

	mediaType := &mediaType{
		Type:    strings.TrimSpace(typeAndSubtype[0]),
		Subtype: strings.TrimSpace(typeAndSubtype[1]),
		Params:  make(map[string]string),
		Q:       1.0,
		I:       i,
	}

	for j := 1; j < len(parts); j++ {
		param := strings.SplitN(parts[j], "=", 2)
		if len(param) != 2 {
			continue
		}

		key := strings.TrimSpace(param[0])
		value := strings.TrimSpace(param[1])

		if key == "q" {
			q, err := strconv.ParseFloat(value, 64)
			if err == nil {
				mediaType.Q = q
			}
		} else {
			mediaType.Params[key] = value
		}
	}

	return mediaType
}

// getMediaTypePriority returns the priority of a media type.
func getMediaTypePriority(typ string, accepted []mediaType, index int) *mediaType {
	var priority *mediaType

	for i := 0; i < len(accepted); i++ {
		spec := specify(typ, &accepted[i], index)

		if spec != nil && (priority == nil ||
			(spec.S > priority.S) ||
			(spec.S == priority.S && spec.Q > priority.Q)) {
			priority = spec
		}
	}

	return priority
}

// specify returns the specificity of the media type.
func specify(typ string, spec *mediaType, index int) *mediaType {
	p := parseMediaType(typ, 0)

	if p == nil {
		return nil
	}

	s := 0

	if strings.EqualFold(spec.Type, p.Type) {
		s |= 4
	} else if spec.Type != "*" {
		return nil
	}

	if strings.EqualFold(spec.Subtype, p.Subtype) {
		s |= 2
	} else if spec.Subtype != "*" {
		return nil
	}

	for key, val := range spec.Params {
		if val == "*" || strings.EqualFold(val, p.Params[key]) {
			s |= 1
		} else {
			return nil
		}
	}

	return &mediaType{
		Type:    spec.Type,
		Subtype: spec.Subtype,
		Params:  spec.Params,
		Q:       spec.Q,
		I:       index,
		S:       s,
	}
}

// isQuality checks if a media type has any quality.
func isQuality(mediaType *mediaType) bool {
	return mediaType.Q > 0
}

// getFullType returns the full type string.
func getFullType(mediaType *mediaType) string {
	return mediaType.Type + "/" + mediaType.Subtype
}

// splitMediaTypes splits an Accept header into media types.
func splitMediaTypes(accept string) []string {
	parts := make([]string, 0)

	for len(accept) > 0 {
		accept = strings.TrimSpace(accept)
		idx := strings.IndexByte(accept, ',')
		if idx == -1 {
			parts = append(parts, accept)
			break
		}
		inQuotes := quoteCount(accept[:idx])%2 == 1
		if inQuotes {
			continue
		}
		parts = append(parts, accept[:idx])
		accept = accept[idx+1:]
	}

	return parts
}

// quoteCount counts the number of quotes in a string.
func quoteCount(str string) int {
	return strings.Count(str, "\"")
}
