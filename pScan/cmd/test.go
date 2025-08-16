package cmd

import (
	"fmt"
	"os"
)

func testErrcheck() {
	// 这些都应该产生 errcheck 警告
	os.Open("nonexistent.txt")
	fmt.Println("test")
	os.Mkdir("/tmp/test", 0o755)

	// Close 方法通常也应该检查
	if f, err := os.Create("/tmp/test.txt"); err == nil {
		f.Close()
	}
}

func main() {
	const (
		KindNone int = iota
		KindPrint
		KindPrintf
		KindErrorf
	)

	y := 3
	y++

	x := false
	if y > 4 {
		x = true
	}

	fmt.Println(x)
}
