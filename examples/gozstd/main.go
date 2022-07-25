package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sunvim/utils/tools"
	"github.com/valyala/gozstd"
)

// the size is 22
func main() {
	fmt.Println("test gozstd")
	fs, _ := os.OpenFile("test.txt", os.O_RDWR|os.O_CREATE, 0666)
	defer fs.Close()
	compress := gozstd.Compress
	for i := 0; i < 10240; i++ {
		bs := strings.Repeat(strconv.Itoa(i), i+1)
		cs := compress(nil, tools.StringToBytes(bs))
		fmt.Fprintf(fs, "index: %d compress %d bytes to %d bytes\n", i, len(bs), len(cs))
	}
}
