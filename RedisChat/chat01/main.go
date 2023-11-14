package main

import (
	"fmt"
	"log"
	"sync"
)

func init() {
	if err := Init(); err != nil {
		log.Fatal(err)
	}

}

func main() {

	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			msg := &Msg{
				Data:    "hello world",
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

//go:generate msgp -tests=false -io=false

type Msg struct {
	Data    string
	Version int64
}

func UpdateMsg(msg *Msg) error {

	bytes, err := msg.MarshalMsg(nil)
	if err != nil {
		return err
	}

	conn := Conn()
	defer conn.Close()

	_, err = conn.Do("SET", "MsgToken", bytes)
	return err
}
