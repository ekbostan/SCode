package network

import (
   "context"
   "sync"
)

type ServerOptions struct {
    Transports []Transport
}

type Server struct {
   Options   ServerOptions
   rpcPool   *sync.Pool
   rpcChan   chan *RPC
   workers   int
   cancel    context.CancelFunc
   ctx       context.Context
}

func NewServer(opts ServerOptions) *Server {
   ctx, cancel := context.WithCancel(context.Background())
   
   s := &Server{
       Options: opts,
       rpcPool: &sync.Pool{
           New: func() interface{} {
               return &RPC{}
           },
       },
       rpcChan: make(chan *RPC, 1024),
       workers: 10,
       ctx:     ctx,
       cancel:  cancel,
   }
   return s
}

func (s *Server) initTransports() {
   for i := 0; i < s.workers; i++ {
       go s.processRPC()
   }

   for _, tr := range s.Options.Transports {
       go func(tr Transport) {
           for {
               select {
               case <-s.ctx.Done():
                   return
               default:
                   rpc := s.rpcPool.Get().(*RPC)
                   
                   select {
                   case msg := <-tr.Consume():
                       rpc.sender = msg.sender
                       rpc.Payload = msg.Payload
                       
                       select {
                       case s.rpcChan <- rpc:
                       case <-s.ctx.Done():
                           s.rpcPool.Put(rpc)
                           return
                       }
                   case <-s.ctx.Done():
                       s.rpcPool.Put(rpc)
                       return
                   }
               }
           }
       }(tr)
   }
}

func (s *Server) processRPC() {
   for {
       select {
       case <-s.ctx.Done():
           return
       case rpc := <-s.rpcChan:
           s.rpcPool.Put(rpc)
       }
   }
}

func (s *Server) Start() {
   s.initTransports()
}

func (s *Server) Shutdown() {
   s.cancel()
   close(s.rpcChan)
}