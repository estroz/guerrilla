package guerrilla

import (
	"sync"

	log "github.com/Sirupsen/logrus"
)

type Guerrilla struct {
	Config  *AppConfig
	servers []*Server
	pool    *ClientPool
}

// Returns a new instance of Guerrilla with the given config, not yet running.
func New(ac *AppConfig) *Guerrilla {
	g := &Guerrilla{ac, []*Server{}, NewClientPool(10)}

	// Instantiate servers
	for _, sc := range ac.Servers {
		// Add app-wide allowed hosts to each server
		sc.AllowedHosts = ac.AllowedHosts
		server, err := NewServer(sc, g.pool)
		if err != nil {
			log.WithError(err).Error("Failed to create server")
		}
		g.servers = append(g.servers, server)
	}
	return g
}

// Entry point for the application. Starts all servers.
func (g *Guerrilla) Run() {
	for _, s := range g.servers {
		go s.Run()
	}
}

type ClientPool struct {
	m       sync.Mutex
	clients chan *Client
}

func NewClientPool(size int) *ClientPool {
	return &ClientPool{clients: make(chan *Client, size)}
}

func (cp *ClientPool) Put(c *Client) {
	select {
	case cp.clients <- c:
	default:
		// Internal channel is saturated, double size
		cp.resize(cap(cp.clients) * 2)
		cp.Put(c)
	}
}

func (cp *ClientPool) Get() *Client {
	select {
	case c := <-cp.clients:
		return c
	default:
		return NewClient()
	}
}

func (cp *ClientPool) resize(newSize int) {
	cp.m.Lock()
	defer cp.m.Unlock()

	if cap(cp.clients) == newSize {
		return // another sneaky goroutine already resized!
	}

	temp := cp.clients
	cp.clients = make(chan *Client, newSize)
	var c *Client
	for len(temp) > 0 {
		c = <-temp
		cp.clients <- c
	}
}
