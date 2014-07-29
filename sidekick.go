package sidekick

import (
	"github.com/coreos/go-etcd/etcd"
	"log"
	"strings"
	"sync"
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
	timer  *time.Timer
	closed bool
	sync.Mutex
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
	sk.timer = time.NewTimer(UpdateInterval * time.Second)
	for {
		select {
		case <-sk.timer.C:
			sk.Lock()
			sk.timer.Reset(UpdateInterval * time.Second)
			sk.Unlock()
			_, err := sk.client.Set(sk.key, sk.value, TTL)
			if err != nil && sk.logger != nil {
				sk.logger.Printf("error updating %s %s\n", sk.key, err.Error())
			}

		case <-sk.quitCh:
			sk.timer.Stop()
			return
		}
	}
}

// Value changes the value used and performs an update.  This is a no-op if Stop() has been called.
func (sk *Sidekick) Value(value string) {
	if sk.closed {
		return
	}
	sk.Lock()
	defer sk.Unlock()
	sk.value = value
	// reset the timer so that an update will happen immediately
	sk.timer.Reset(0)
}

// Stop stops the goroutine performing updates.  This is a no-op if Stop() has been called.
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
