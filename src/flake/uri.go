package flake

import (
	"net/url"
)

type FlakeUrl struct {
	URL       string
	Base      string
	Scheme    string
	Authority *string
	Path      string
	Query     url.Values
	Fragment  string
}

func UrlFromString(rawURL string) (*FlakeUrl, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	scheme := parsed.Scheme
	authority := parsed.Host
	path := parsed.Path
	query := parsed.Query()
	fragment := parsed.Fragment

	base := scheme + "://" + authority

	var authPtr *string
	if authority != "" {
		authPtr = &authority
	}

	return &FlakeUrl{
		URL:       rawURL,
		Base:      base,
		Scheme:    scheme,
		Authority: authPtr,
		Path:      path,
		Query:     query,
		Fragment:  fragment,
	}, nil
}

// ToString returns a string representation of ParsedURL
func (p *FlakeUrl) NixFlakeUrl() string {
	return p.URL
}
