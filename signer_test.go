package signature_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	signature "github.com/lopz82/aws-signature"
)

func TestSignRequest(t *testing.T) {

	testCases := []struct {
		region, service string
		body            io.Reader
	}{
		{
			"eu-central-1",
			"es",
			bytes.NewBufferString("SIMPLE BODY TO BE ENCRYPTED"),
		},
		{
			"eu-central-1",
			"es",
			http.NoBody,
		},
	}

	for _, test := range testCases {

		cfg := signature.CreateConfig()
		ctx := context.Background()
		next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

		cfg.Region = test.region
		cfg.Service = test.service
		body := test.body

		handler, err := signature.New(ctx, next, cfg, "aws-signer")
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", body)
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(recorder, req)

		assertContainsSignedHeaders(t, req)
	}
}

func assertContainsSignedHeaders(t *testing.T, req *http.Request) {
	assertHeaderExists(t, req, "Authorization")
	assertHeaderExists(t, req, "X-Amz-Date")
}

func assertHeaderExists(t *testing.T, req *http.Request, key string) {
	t.Helper()

	if req.Header.Get(key) == "" {
		t.Errorf("header inexistent: %s", key)
	}
}

func TestIsExistingServiceInRegion(t *testing.T) {
	testCases := []struct {
		region, service string
		expected        error
	}{
		{
			region:   "eu-central-1",
			service:  "es",
			expected: nil,
		},
		{
			"eu-central",
			"es",
			errors.New("region eu-central does not exist"),
		}, {
			"eu-central-1",
			"ese",
			errors.New("service ese does not exist"),
		},
	}
	for _, test := range testCases {
		res := signature.IsExistingServiceInRegion(test.region, test.service)
		if res != nil {
			if res.Error() != test.expected.Error() {
				t.Errorf("invalid test: got %s expecting %s", res, test.expected)
			}
		}
	}
}
