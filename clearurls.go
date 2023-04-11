package clearurls

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
)

var ErrBlockedURL = errors.New("URL is blocked")

//go:embed data.min.json
var clearURLsData []byte

type pattern struct {
	pattern  string
	compiled *regexp.Regexp
}

func (p *pattern) re() *regexp.Regexp {
	if p.compiled == nil {
		p.compiled = regexp.MustCompile(p.pattern)
	}

	return p.compiled
}

func (p *pattern) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, &p.pattern)
	if err != nil {
		return err
	}

	p.pattern = "(?i)" + p.pattern

	return nil
}

type Provider struct {
	UrlPattern        pattern   `json:"urlPattern"`
	CompleteProvider  bool      `json:"completeProvider"`
	Rules             []pattern `json:"rules"`
	RawRules          []pattern `json:"rawRules"`
	ReferralMarketing []pattern `json:"referralMarketing"`
	Exceptions        []pattern `json:"exceptions"`
	Redirections      []pattern `json:"redirections"`
	ForceRedirection  bool      `json:"forceRedirection"`
}

type URLCleaner struct {
	providers map[string]Provider
}

func New() *URLCleaner {
	var rules struct {
		Providers map[string]Provider `json:"providers"`
	}

	err := json.Unmarshal(clearURLsData, &rules)
	if err != nil {
		panic(err)
	}

	return &URLCleaner{
		providers: rules.Providers,
	}
}

func (c *URLCleaner) Clean(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	return c.CleanURL(u)
}

func (c *URLCleaner) CleanURL(u *url.URL) (string, error) {
MainLoop:
	for _, provider := range c.providers {
		re := provider.UrlPattern.re()
		if !re.MatchString(fmt.Sprintf("%s://%s", u.Scheme, u.Hostname())) {
			continue
		}

		if provider.CompleteProvider {
			return "", ErrBlockedURL
		}

		for _, exception := range provider.Exceptions {
			if exception.re().MatchString(u.String()) {
				continue MainLoop
			}
		}

		for _, redirection := range provider.Redirections {
			matches := redirection.re().FindStringSubmatch(u.String())
			if len(matches) > 0 {
				match, err := url.QueryUnescape(matches[1])
				if err != nil {
					continue
				}

				u, err = u.Parse(match)
				if err != nil {
					continue
				}

				if u.Scheme == "" {
					u.Scheme = "http"
				}

				return c.CleanURL(u)
			}
		}

		if len(u.Query()) > 0 {
			qs := c.cleanQuery(u.Query(), provider)

			u.RawQuery = qs.Encode()
		} else {
			u.RawQuery = u.Query().Encode()
		}

		for _, rule := range provider.RawRules {
			u.Path = rule.re().ReplaceAllString(u.Path, "")
		}
	}

	return u.String(), nil
}

func (c *URLCleaner) cleanQuery(qs url.Values, provider Provider) url.Values {
	rules := append(provider.Rules, provider.ReferralMarketing...)

	for _, qp := range rules {
		for k := range qs {
			if qp.re().MatchString(k) {
				qs.Del(k)
			}
		}
	}

	return qs
}
