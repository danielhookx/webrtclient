package client

import (
	"flag"
	"net/url"
	"testing"
)

var addr = flag.String("addr", "172.16.101.131:19801", "http service address")
var u = url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}

func TestClient(t *testing.T) {
	flag.Parse()
	c, err := NewClient("test", u)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(c.Read())
}
