package network

import (
   "sync"
)

type ServerOptions struct {
   Transports []Transport
   Workers    int
}

type Server struct {
   Options  ServerOptions
   rpcChan  chan *RPC
   done     chan struct{}
   wg       sync.WaitGroup
   handler  Handler
}

type Handler func(rpc *RPC) error

func NewServer(opts ServerOptions) *Server {
   if opts.Workers == 0 {
       opts.Workers = 8
   }
   
   return &Server{
       Options: opts,
       rpcChan: make(chan *RPC, 1024),
       done:    make(chan struct{}),
   }
}

func (s *Server) RegisterHandler(handler Handler) {
   s.handler = handler
}

func (s *Server) initTransports() {
   for i := 0; i < s.Options.Workers; i++ {
       s.wg.Add(1)
       go s.processRPC()
   }

   for _, tr := range s.Options.Transports {
       s.wg.Add(1)
       go s.handleTransport(tr)
   }
}

func (s *Server) handleTransport(tr Transport) {
   defer s.wg.Done()
   
   for {
       select {
       case <-s.done:
           return
       case msg := <-tr.Consume():
           rpc := &RPC{
               sender:  msg.sender,
               Payload: msg.Payload,
           }
           
           select {
           case s.rpcChan <- rpc:
           case <-s.done:
               return
           }
       }
   }
}

func (s *Server) processRPC() {
   defer s.wg.Done()
   
   for {
       select {
       case <-s.done:
           return
       case rpc := <-s.rpcChan:
           if s.handler != nil {
               if err := s.handler(rpc); err != nil {
                   // Handle error
                   continue
               }
           }
       }
   }
}

func (s *Server) Start() error {
   s.initTransports()
   return nil
}

func (s *Server) Shutdown() error {
   close(s.done)
   s.wg.Wait()
   close(s.rpcChan)
   return nil
}