package network

import (
    "testing"
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