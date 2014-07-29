package sidekick

import (
	"github.com/coreos/go-etcd/etcd"
	"log"
	"strings"
	"time"
)

// TTL is the time-to-live in seconds for updates to etcd
const TTL = 10
// UpdateInterval is how frequently in seconds the key will be updated until Stop() is called
const UpdateInterval = 8

// Sidekick periodically updates a key in etcd until Stop() is called
type Sidekick struct {
	client *etcd.Client
	key    string
	value  string
	logger *log.Logger
	quitCh chan bool
	closed bool
}

// New returns a Sidekick pointer if there was no error initially setting the value,
// and there will be a goroutine updating the value every UpdateInterval seconds until
// Stop() is called
// If a non-nil error is returned there will be no goroutine performing updates
func New(servers string, key string, value string) (*Sidekick, error) {
	etcdNodes := strings.Split(servers, ",")
	sk := &Sidekick{}
	sk.client = etcd.NewClient(etcdNodes)
	sk.quitCh = make(chan bool)
	sk.key = key
	_, err := sk.client.Set(sk.key, value, TTL)
	if err != nil {
		return nil, err
	}
	go sk.updateLoop()
	return sk, nil
}

// SetLogger sets a logger, by default no logs are written
func (sk *Sidekick) SetLogger(logger *log.Logger) {
	sk.logger = logger
}
func (sk *Sidekick) updateLoop() {
	ticker := time.NewTicker(UpdateInterval * time.Second)
	for {
		select {
		case <-ticker.C:
			_, err := sk.client.Set(sk.key, sk.value, TTL)
			if err != nil && sk.logger != nil {
				sk.logger.Printf("error updating %s %s\n", sk.key, err.Error())
			}

		case <-sk.quitCh:
			ticker.Stop()
			return
		}
	}
}

// Stop stops the goroutine performing updates
func (sk *Sidekick) Stop() {
	if sk.closed {
		return
	}
	sk.closed = true
	close(sk.quitCh)
	_, err := sk.client.Delete(sk.key, false)
	if err != nil && sk.logger != nil {
		sk.logger.Printf("error deleting %s %s\n", sk.key, err.Error())
	}
	if sk.logger != nil {
		sk.logger.Printf("stopped")
	}
}
