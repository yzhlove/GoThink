package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func main() {

	test := Test{
		msgs: make(chan struct{}, 1),
	}
	test.Monitor()
	test.Start()
	select {}
}

type Test struct {
	variable atomic.Int32
	msgs     chan struct{}
	status   atomic.Bool
}

func (t *Test) Action() {
	t.variable.Add(1)
	fmt.Println("this is value: ", t.variable.Load())
	if t.variable.Load()%5 == 0 {
		t.status.Store(true)
	}
	if t.variable.Load()%10 == 0 {
		t.status.Store(false)
	}
}

func (t *Test) Monitor() {
	go func() {
		for range t.msgs {
			t.Action()
		}
	}()
}

func (t *Test) Signal() {
	t.msgs <- struct{}{}
}

func (t *Test) Start() {
	go func() {
		timer1 := time.NewTicker(time.Second)
		timer2 := time.NewTicker(time.Second * 5)

		for {
			select {
			case <-timer1.C:
				fmt.Println("timer1 starting .")
				if !t.status.Load() {
					t.Signal()
				}
			case <-timer2.C:
				fmt.Println("timer2 starting .")
				if t.status.Load() {
					t.Signal()
				}
			}
		}
	}()
}
