package network

import (
	"runtime"
	"sync"
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


func TestServer(t *testing.T){
    transport1 := NewLocalTransport("server")
    transport2 := NewLocalTransport("client")
    transports := []Transport{transport1,transport2}

    opts := ServerOptions{
        Transports: transports,
        Workers:    runtime.NumCPU(),
    }

    server := NewServer(opts)

    assert.Equal(t,opts.Workers, server.Options.Workers)
    assert.Equal(t, opts.Transports, server.Options.Transports)
    assert.NotNil(t,server.rpcChan)
    assert.NotNil(t,server.done)


    var receivedMessages    []string
    var mu  sync.Mutex

    server.RegisterHandler(func(rpc *RPC) error {
        mu.Lock()
        receivedMessages = append(receivedMessages, string(rpc.Payload))
        mu.Unlock()
        return nil
    })
    

    err := server.Start()
    assert.NoError(t,err)


    transport1.Connect(transport2)
    transport2.Connect(transport1)


    messages := []string{
        "Message 1",
        "Message 2",
        "Message 3",
    }

    var wg sync.WaitGroup

    for _,msg := range messages{
        wg.Add(1)
        go func(message string){
            defer wg.Done()
            err := transport2.SendMessage(transport1.sender,[]byte(msg))
            assert.NoError(t,err)

        }(msg)
    }

    wg.Wait()

    time.Sleep(100 * time.Millisecond)

    mu.Lock()
    assert.Equal(t,len(messages),len(receivedMessages))

    for _, msg := range messages{
        assert.Contains(t,receivedMessages,msg)
    }
    mu.Unlock()

    err = server.Shutdown()
    assert.NoError(t,err)


    select{
    case _,ok := <- server.done:
        assert.False(t,ok,"done channel should be closed")
    default:
    }

    select{
    case _,ok := <-server.rpcChan:
        assert.False(t,ok,"rpcChan should be closed")
    default:
    }
}

func TestServerWorkerCount(t *testing.T) {
        tests := []struct {
            name          string
            workerCount   int
            expectedCount int
        }{
            {
                name:          "Zero workers defaults to NumCPU",
                workerCount:   0,
                expectedCount: runtime.NumCPU(),
            },
            {
                name:          "Specific worker count",
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
    