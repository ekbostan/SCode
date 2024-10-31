package network


type NetworkAddress string

type RPC struct{
	sender NetworkAddress
	Payload []byte
}

type Transport interface{
	Consume() <- chan RPC
	Connect(Transport) error
	SendMessage(NetworkAddress,[]byte) error
	GetAdress() NetworkAddress
}