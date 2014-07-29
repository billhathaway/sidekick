// sideKick_test.go
package sidekick

import (
	"testing"
)

// This requires a local etcd server to work 
// Todo: pick a better random path for testing that definitely will not interfer
func TestSideKick(t *testing.T) {
	servers := "http://127.0.0.1:4001"
	key := "/tmppath/server1"
	value := "testval1"
	sk, err := New(servers, key, value)
	if err != nil {
                t.Logf("An etcd server must be running on http://127.0.0.1:4001 for tests to work")
		t.Errorf(err.Error())
		return
	}
	sk.Stop()
}
