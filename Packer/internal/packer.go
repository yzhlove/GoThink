package internal

import (
	"bytes"
	"context"
	"sync"
	"sync/atomic"
)

type Packer interface {
	GetName() string
	GetReader() *bytes.Reader
}

type Pack struct {
	size     int
	names    []string
	monitors []*monitor
	waiting  sync.WaitGroup
	errCh    chan error
	err      error
	reqQueue chan Packer

	ctx    context.Context
	cancel context.CancelFunc
}

func (p *Pack) running() {
	for i := 0; i < len(p.monitors); i++ {
		p.waiting.Add(1)
		go func(m *monitor) {
			defer p.waiting.Done()
			p.pack(m)
		}(p.monitors[i])
	}
	p.waiting.Wait()
}

func (p *Pack) pack(m *monitor) {
	defer m.close()
	for {
		select {
		case <-p.ctx.Done():
			return
		case req, ok := <-p.reqQueue:
			if !ok {
				return
			}
			ok, err := m.submit(req)
			if err != nil {
				p.throwErr(err)
				return
			}
			if !ok {
				return
			}
		}
	}
}

func (p *Pack) throwErr(err error) {
	select {
	case p.errCh <- err:
	default:
	}
}

func (p *Pack) Submit(req Packer) error {

	return nil
}

type monitor struct {
	limit    int64
	size     atomic.Int64
	reqQueue chan Packer
	write    *ZipWriter
	errCh    chan error
	err      error
}

func (m *monitor) throwErr(err error) {
	select {
	case m.errCh <- err:
		m.err = err
	default:
	}
}

func (m *monitor) isWrite(size int64) bool {
	if m.size.Load() >= m.limit {
		return false
	}
	m.size.Add(size)
	return true
}

func (m *monitor) newMonitor(name string) (*monitor, error) {
	w, err := NewZipWrite(name)
	if err != nil {
		return nil, err
	}

	return &monitor{
		reqQueue: make(chan Packer, 16),
		write:    w,
	}, err
}

func (m *monitor) running() {
	for {
		if req, ok := <-m.reqQueue; ok {
			reader := req.GetReader()
			if m.isWrite(reader.Size()) {
				if err := m.write.Write(req.GetName(), reader); err != nil {
					m.throwErr(err)
					return
				}
				continue
			}
			return
		}
		return
	}
}

func (m *monitor) submit(req Packer) (bool, error) {

	return false, nil
}

func (m *monitor) close() {

}
