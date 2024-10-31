package network

import (
	"fmt"
	"sync"
)


type localTransport struct{
	sender 	NetworkAddress
	consumeCh chan RPC
	lock	*sync.RWMutex
	peers	map[NetworkAddress]localTransport

}


func NewLocalTransport(sender NetworkAddress) *localTransport{
	return &localTransport{
		sender:sender,
		consumeCh: make(chan RPC,1024),
		peers: make(map[NetworkAddress]localTransport),
	}

}

func (t *localTransport) Consume() <- chan RPC{
	return t.consumeCh
}

func (t *localTransport) Connect(tr *localTransport) error{
	t.lock.Lock()
	defer t.lock.Unlock()
	t.peers[tr.GetAdress()] = *tr

	return nil
}

func(t *localTransport) SendMessage(to NetworkAddress, payload []byte) error{
	t.lock.RLock()
	defer t.lock.Unlock()
	peer,ok := t.peers[to]
	if !ok{
		return fmt.Errorf("%s: Could not send message to %s",t.sender,to)
	}
	peer.consumeCh <- RPC{
		sender: t.sender,
		Payload: payload ,

	}
	return nil
}

func(t *localTransport) GetAdress() NetworkAddress{
	return t.sender
}