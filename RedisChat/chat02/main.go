package main

import (
	"fmt"
	"github.com/reids-chat/chat02/db"
	"log"
	"sync"
)

//go:generate msgp -tests=false -io=false

type Msg struct {
	Data    []byte
	Version int64
}

func main() {

	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			msg := &Msg{
				Data:    []byte("hello world"),
				Version: int64(i + 1),
			}
			if err := UpdateMsg(msg); err != nil {
				fmt.Println("write msg error: ", err)
			}

			fmt.Println("write db msg is --> ", msg)
		}(i)
	}

	wg.Wait()
	fmt.Println("write OK.")

}

func UpdateMsg(msg *Msg) error {

	bytes, err := msg.MarshalMsg(nil)
	if err != nil {
		return err
	}

	conn := db.Conn()
	defer conn.Close()

	resp, err := redisAtomic.Do(conn, "MSG_TOKEN", msg.Version, bytes)
	if err != nil {
		return err
	}

	fmt.Println("resp --> ", resp)
	return nil
}
