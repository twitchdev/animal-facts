package main

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewTLSServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return cli, s.Close
}

const (
	okResponse = `{
    "broadcaster:999999999": {
        "segment": {
            "segment_type": "broadcaster",
            "channel_id": "999999999"
        },
        "record": {
            "content": "cat"
        }
    },
    "developer:999999999": {
        "segment": {
            "segment_type": "developer",
            "channel_id": "999999999"
        },
        "record": {
            "content": "In the 1750s, Europeans introduced cats into the Americas to control pests."
        }
    }
	}`
)

func TestGetBroadcasterSegment(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(okResponse))
	})
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	cli := &ConfigurationServiceClient{
		client:   httpClient,
		clientID: "some-client-id",
	}

	animalType := cli.GetBroadcasterSegment("999999999")

	assert.Equal(t, "cat", animalType)
}
