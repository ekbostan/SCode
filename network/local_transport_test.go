package network

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
    tra := NewLocalTransport("A")
    trb := NewLocalTransport("B")
    
    tra.Connect(trb)
    trb.Connect(tra)

    assert.Equal(t, trb, tra.peers[trb.GetAdress()])
    assert.Equal(t, tra, trb.peers[tra.GetAdress()])
}

func TestSendMessage(t *testing.T) {
    tra := NewLocalTransport("A")
    trb := NewLocalTransport("B")
    
    tra.Connect(trb)
    trb.Connect(tra)

    message := make([]byte, 1000)
    assert.Nil(t, tra.SendMessage(trb.GetAdress(), message))

    rpc := <-trb.Consume()
    assert.Equal(t, message, rpc.Payload)
    assert.Equal(t, tra.GetAdress(), rpc.sender)
}

func TestServer(t *testing.T) {
    transport1 := NewLocalTransport("server")
    transport2 := NewLocalTransport("client")
    
    opts := ServerOptions{
        Transports: []Transport{transport1, transport2},
        Workers:    runtime.NumCPU(),
    }

    server := NewServer(opts)

    assert.NotNil(t, server)
    assert.Equal(t, opts.Transports, server.Options.Transports)
    assert.Equal(t, opts.Workers, server.Options.Workers)

    messageReceived := make(chan bool)
    server.RegisterHandler(func(rpc *RPC) error {
        assert.Equal(t, []byte("test message"), rpc.Payload)
        assert.Equal(t, transport2.GetAdress(), rpc.sender)
        messageReceived <- true
        return nil
    })

    server.Start()

    
    transport1.Connect(transport2)
    transport2.Connect(transport1)

   
    err := transport2.SendMessage(transport1.GetAdress(), []byte("test message"))
    assert.Nil(t, err)


    select {
    case <-messageReceived:
    case <-time.After(time.Second):
        t.Fatal("Message was not processed within timeout")
    }

    server.Shutdown()
}

func TestServerWorkerCount(t *testing.T) {
    cpuCount := runtime.NumCPU()
    
    tests := []struct {
        name          string
        workerCount   int
        expectedCount int
    }{
        {
            name:          "Default worker count",
            workerCount:   0,
            expectedCount: cpuCount,
        },
        {
            name:          "Custom worker count",
            workerCount:   4,
            expectedCount: 4,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            opts := ServerOptions{
                Transports: []Transport{NewLocalTransport("test")},
                Workers:    tt.workerCount,
            }
            
            server := NewServer(opts)
            assert.Equal(t, tt.expectedCount, server.Options.Workers)
        })
    }
}