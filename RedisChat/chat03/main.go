package main

import (
	"fmt"
	"github.com/reids-chat/chat03/db"
	"github.com/reids-chat/chat03/msg"
	"log"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

func main() {

	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	var msgs = make([]*msg.Msg, 0, 100)
	var mutex sync.Mutex
	var errNumber atomic.Int32

	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func(i int) {
			//time.Sleep(time.Millisecond * time.Duration(rand.Intn(5)))
			defer wg.Done()
			v := &msg.Msg{
				Value:   fmt.Sprintf("this is value {%d}", i),
				Version: time.Now().UnixNano(),
			}

			mutex.Lock()
			msgs = append(msgs, v)
			mutex.Unlock()

			err := msg.Set("msg_test_key", v)
			if err != nil {
				errNumber.Add(1)
			}
		}(i)
	}

	wg.Wait()

	fmt.Println("OK.")

	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].Version > msgs[j].Version
	})

	fmt.Println("msgs[0] ==> ", msgs[0], " err failed number --> ", errNumber.Load())

	message, err := msg.Get("msg_test_key")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("message: %s %v \n", string(message.Value), message.Version)

}
