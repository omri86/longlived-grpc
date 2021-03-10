package main

import (
	"../protos"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
)

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:7070")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer([]grpc.ServerOption{}...)

	server := &longlivedServer{
		subscribers: make(map[int32]chan<- *protos.Response, 1),
		subsLock: &sync.RWMutex{},
	}

	// Start sending data to subscribers
	go mockDataGenerator(server)

	protos.RegisterLonglivedServer(grpcServer, server)
	log.Printf("Starting server on address %s", lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}

type longlivedServer struct {
	protos.UnimplementedLonglivedServer
	subscribers map[int32]chan<- *protos.Response // subscribers maps a client ID to a channel
	subsLock    *sync.RWMutex                     // subsLock ensures no conflicts when modifying the subscribers
}

// Subscribe handles a subscribe request from a client
// Note that once the scope of this function returns the stream is closed
func (s *longlivedServer) Subscribe(request *protos.Request, stream protos.Longlived_SubscribeServer) error {
	// Handle subscribe request
	log.Printf("Received subscribe request from ID: %d\n", request.Id)

	// Create a channel for this subscriber
	c := make(chan *protos.Response)

	s.subsLock.Lock()
	// Save the subscriber channel according to the ID
	s.subscribers[request.Id] = c
	s.subsLock.Unlock()


	for {
		for msg := range c {
			if err := stream.Send(msg); err != nil {
				// In case of error the client would re-subscribe so close and delete the channel for this subscriber
				s.subsLock.Lock()
				close(s.subscribers[request.Id])
				delete(s.subscribers, request.Id)
				s.subsLock.Unlock()
				log.Printf("Failed to send data to client: %v", err)
				return err
			}
		}
	}
}

func mockDataGenerator(server *longlivedServer) {
	log.Println("Starting data generation")
	for {
		time.Sleep(time.Second)

		server.subsLock.RLock()
		for id, channel := range server.subscribers {
			select {
			case channel <- &protos.Response{Data: fmt.Sprintf("data mock for: %d", id)}:
				log.Printf("Data sent to client %d", id)
			default:
				log.Printf("Channel buffer full - avoid sending data to client: %d", id)
			}
		}
		server.subsLock.RUnlock()
	}
}
