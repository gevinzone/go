package main

import (
	"flag"
	"fmt"
)

// go run main.go -name=gevin -x 111
func main() {
	var name, x string
	flag.StringVar(&name, "name", "Go 语言命令行", "帮助信息")
	flag.StringVar(&name, "n", "Go 语言命令行", "帮助信息")
	flag.StringVar(&x, "x", "Go 语言命令行", "帮助信息")
	flag.Parse()
	fmt.Printf("name: %s, x: %s\n", name, x)
}
