package client

import "testing"

func TestPut(t *testing.T) {
	err := Put("http://172.16.101.131:19801/pub", "demo", "testsdphhh")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success")
	return
}
