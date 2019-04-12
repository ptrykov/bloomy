package main

import (
  "runtime"
  server "github.com/ptrykov/bloomy/internal"
)

func main() {
  runtime.GOMAXPROCS(runtime.NumCPU())

  server := server.NewServer();
  server.Run();
}
