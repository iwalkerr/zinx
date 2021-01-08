package main

import "gozinx/znet"

func main() {
	s := znet.NewServer()
	s.Serve()
}
