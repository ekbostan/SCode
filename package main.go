// cmd/main.go
package main

import (
   "github.com/ekbostan/SCode/network"
)

func main() {
   trlocal := network.NewLocalTransport("LOCAL")

   opts := network.ServerOptions{
       Transports: []network.Transport{trlocal},
   }

   server := network.NewServer(opts)
   server.Start()
   

}