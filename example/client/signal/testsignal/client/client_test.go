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
	c := NewClient("test", u)
	t.Log(c.Read())
}