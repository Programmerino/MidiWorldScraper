package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/cavaliercoder/grab"
)

var emptyFiles int
var maxAtOnce = 10

var wg sync.WaitGroup

func main() {
	for i := 0; emptyFiles <= 100; i++ {
		if i%maxAtOnce == 0 {
			fmt.Println("Waiting for completion")
			wg.Wait()
		}
		wg.Add(1)
		go download(i)
	}
	wg.Wait()
	fmt.Println("Done!")
}

func download(i int) {
	defer wg.Done()
	client := grab.NewClient()
	req, _ := grab.NewRequest(".", fmt.Sprintf("http://www.midiworld.com/download/%v", i))
	resp := client.Do(req)
	for {
		select {
		case <-resp.Done:
			// download is complete
			break
		}
		break
	}
	fmt.Println(fmt.Sprintf("Finished downloading %s", resp.Filename))
	if resp.Size == 0 {
		emptyFiles++
		os.Remove(resp.Filename)
		fmt.Println(fmt.Sprintf("File %v is 0 bytes long! Will stop in %v", i, 100-emptyFiles))
	} else {
		emptyFiles = 0
	}
}
