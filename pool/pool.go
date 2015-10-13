package pool

import (
	"errors"
	"sync"
)

var (
	NameInUseErr            = errors.New("name in use")
	NameNotFoundErr         = errors.New("name not found")
	ResourceNotAvailableErr = errors.New("resource not available")
)

type Pool interface {
	Idle() int
	Leased() int

	Register(name string) error
	Unregister(name string) error

	Lease() (string, error)
	Return(string) error
}

func New() Pool {
	return &pool{
		idle:   make(map[string]bool),
		leased: make(map[string]bool),
	}
}

type pool struct {
	m sync.Mutex

	idle   map[string]bool
	leased map[string]bool
}

func (p *pool) Idle() int {
	p.m.Lock()
	defer p.m.Unlock()

	return len(p.idle)
}

func (p *pool) Leased() int {
	p.m.Lock()
	defer p.m.Unlock()

	return len(p.leased)
}

func (p *pool) Register(name string) error {
	println("Register")
	p.m.Lock()
	defer p.m.Unlock()

	_, idleExists := p.idle[name]
	_, leasedExists := p.leased[name]
	if idleExists || leasedExists {
		return NameInUseErr
	}

	p.idle[name] = true
	return nil
}

func (p *pool) Unregister(name string) error {
	println("Unregister")
	p.m.Lock()
	defer p.m.Unlock()

	_, idleExists := p.idle[name]
	_, leasedExists := p.leased[name]
	if !idleExists && !leasedExists {
		return NameNotFoundErr
	}

	delete(p.idle, name)
	delete(p.leased, name)
	return nil
}

func (p *pool) Lease() (name string, err error) {
	println("Lease")
	p.m.Lock()
	defer p.m.Unlock()

	for name, _ = range p.idle {
		delete(p.idle, name)
		p.leased[name] = true
		return name, nil
	}

	return "", ResourceNotAvailableErr
}

func (p *pool) Return(name string) error {
	println("Return")
	p.m.Lock()
	defer p.m.Unlock()

	if !p.leased[name] {
		return NameNotFoundErr
	}

	delete(p.leased, name)
	p.idle[name] = true
	return nil
}
