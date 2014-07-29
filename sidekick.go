package sidekick

import (
	"errors"
	"github.com/coreos/go-etcd/etcd"
	"log"
	"strings"
	"sync"
	"time"
)

// DefaultTTL is the default time-to-live in seconds for updates to etcd
const DefaultTTL = 10

// DefaultUpdateInterval is the default frequently in seconds the key will be updated until Stop() is called
const DefaultUpdateInterval = 8

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
	ttl            uint64
	updateInterval uint64
}

var (
	ErrIntervalTooSmall = errors.New("interval must be at least 1 second")
)

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
	sk.ttl = DefaultTTL
	sk.updateInterval = DefaultUpdateInterval
	_, err := sk.client.Set(sk.key, sk.value, sk.ttl)
	if err != nil {
		return nil, err
	}
	go sk.loop()
	return sk, nil
}

// SetLogger sets a logger, by default no logs are written. This is a no-op if Stop() has been called.
func (sk *Sidekick) SetLogger(logger *log.Logger) {
	if sk.closed {
		return
	}
	sk.Lock()
	defer sk.Unlock()
	sk.logger = logger
}

// TTL sets the time-to-live on every update made to etcd. This is a no-op if Stop() has been called.
// TODO: validation on TTL
func (sk *Sidekick) TTL(ttl uint64) error {
	if sk.closed {
		return nil
	}
	sk.Lock()
	defer sk.Unlock()
	sk.ttl = ttl
	return nil
}

// UpdateInterval sets the update interval to the value in seconds. This is a no-op if Stop() has been called.
func (sk *Sidekick) UpdateInterval(interval uint64) error {
	if sk.closed {
		return nil
	}
	if interval < 1 {
		return ErrIntervalTooSmall
	}

	sk.Lock()
	defer sk.Unlock()
	sk.updateInterval = interval
	return nil
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

// loop keeps looping until the quitCh is closed
func (sk *Sidekick) loop() {
	sk.timer = time.NewTimer(time.Duration(sk.updateInterval) * time.Second)
	for {
		select {
		case <-sk.timer.C:
			sk.Lock()
			sk.timer.Reset(time.Duration(sk.updateInterval) * time.Second)
			sk.Unlock()
			_, err := sk.client.Set(sk.key, sk.value, sk.ttl)
			if err != nil && sk.logger != nil {
				sk.logger.Printf("error updating %s %s\n", sk.key, err.Error())
			}

		case <-sk.quitCh:
			sk.timer.Stop()
			sk.timer = nil
			return
		}
	}
}
