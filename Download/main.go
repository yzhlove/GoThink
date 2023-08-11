package main

import (
	"fmt"
	"github.com/downloader/internal/task"
	"time"
)

func main() {

	tk := task.New()

	for i := 0; i < 10; i++ {
		req := &task.Req{
			FilePath:    fmt.Sprintf("baidu.%d.txt", i+1),
			DownloadUrl: "http://www.baidu.com",
			Total:       10,
			Current:     i + 1,
		}
		if err := tk.Submit(req); err != nil {
			fmt.Println("error -> ", err)
			return
		}
	}

	time.Sleep(time.Second)
	tk.Destroy()
	fmt.Println("download OK.")

}
