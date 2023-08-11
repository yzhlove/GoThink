package task

import (
	"context"
	"fmt"
	"github.com/downloader/internal/helper"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
)

type Req struct {
	FilePath    string
	DownloadUrl string
	Total       int
	Current     int
}

type Work struct {
	size        int
	retryTimes  int
	retryDetail sync.Map
	queue       chan *Req
	errCh       chan error
	err         error
	status      atomic.Bool
	ctx         context.Context
	cancel      context.CancelFunc
}

func New() *Work {
	ctx, cancel := context.WithCancel(context.Background())
	w := &Work{
		size:       runtime.NumCPU(),
		retryTimes: 3,
		queue:      make(chan *Req, 16),
		errCh:      make(chan error, 1),
		ctx:        ctx,
		cancel:     cancel,
	}
	w.status.Store(true)
	w.healthCheck()
	w.running()
	return w
}

func (w *Work) healthCheck() {
	fmt.Println("starting health check ...")
	go func() {
		select {
		case <-w.ctx.Done():
			return
		case err := <-w.errCh:
			w.err = err
			w.Destroy()
			return
		}
	}()
}

func (w *Work) running() {
	fmt.Println("starting running ...")
	go func() {
		var wg sync.WaitGroup
		for i := 0; i < w.size; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for {
					select {
					case <-w.ctx.Done():
						return
					case req := <-w.queue:
						fmt.Printf("[%d-%.2d] => %s : %s \n", req.Total, req.Current, req.DownloadUrl, req.FilePath)
						if err := w.handle(req); err != nil {
							w.errCh <- err
							return
						}
					}
				}

			}()
		}
		wg.Wait()
	}()
}

func (w *Work) Destroy() {
	w.status.Store(false)
	w.cancel()
}

func (w *Work) Submit(req *Req) error {
	if w.status.Load() {
		w.queue <- req
		return nil
	}
	return w.err
}

func (w *Work) handle(req *Req) error {
	resp, err := http.Get(req.DownloadUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return helper.WriteFile(req.FilePath, resp.Body)
}
