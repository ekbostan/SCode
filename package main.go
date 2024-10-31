package main

import (
    "fmt"
    "sync"
    "time"
    "github.com/ekbostan/SCode/network"
)

func main() {
    trlocal := network.NewLocalTransport("LOCAL")
    trRemote := network.NewLocalTransport("REMOTE")
    
    opts := network.ServerOptions{
        Transports: []network.Transport{trlocal},
    }
    
    server := network.NewServer(opts)
    server.RegisterHandler(func(rpc *network.RPC) error {
        fmt.Printf("Received message: %s\n", string(rpc.Payload))
        return nil
    })
    
    server.Start()
    
    messages := []string{
        "Hello World 1",
        "Hello World 2",
        "Hello World 3",
        "Hello World 4",
        "Hello World 5",
    }
    
    var wg sync.WaitGroup
    
    for _, msg := range messages {
        wg.Add(1)
        go func(message string) {
            defer wg.Done()
            fmt.Printf("Sending message: %s\n", message)
            trRemote.SendMessage(trlocal.GetAdress(), []byte(message))
        }(msg)
    }
    
    wg.Wait()
    time.Sleep(time.Second)
}