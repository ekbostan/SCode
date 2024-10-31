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

    // Fixed assertions
    assert.Equal(t, trb, tra.peers[trb.sender])
    assert.Equal(t, tra, trb.peers[tra.sender])
}


func TestSendMessage(t *testing.T){
	tra := NewLocalTransport("A")
    trb := NewLocalTransport("B")
    tra.Connect(trb)
    trb.Connect(tra)
	message := make([]byte,1000)
	assert.Nil(t,tra.SendMessage(trb.sender,message))

	rpc := <- trb.Consume()
	assert.Equal(t,rpc.Payload,message)
	assert.Equal(t,rpc.sender,tra)

}