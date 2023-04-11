package clearurls_test

import (
	"testing"

	"github.com/hoshsadiq/go-clearurls"
)

func TestCleanURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{
			name: "does not modify url",
			url:  "https://natura.com.br/p/2458?consultoria=promotop",
			want: "https://natura.com.br/p/2458?consultoria=promotop",
		},

		{
			name: "removes rawRules",
			url:  "https://amazon.com/ref=blabla",
			want: "https://amazon.com",
		},

		{
			name: "removes utm_source",
			url:  "https://example.org/?utm_source=google",
			want: "https://example.org/",
		},

		{
			name: "removes google utm_source",
			url:  "https://myaccount.google.com/?utm_source=google",
			want: "https://myaccount.google.com/",
		},

		{
			name: "redirects",
			url:  "http://googleadservices.com/link/click?adurl=http://g.co/",
			want: "http://g.co/",
		},

		{
			name: "removes unnecessary characters",
			url:  "http://example.com/?&&&&",
			want: "http://example.com/",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			cu := clearurls.New()

			got, err := cu.Clean(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("CleanURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CleanURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
