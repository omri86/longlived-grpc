package main

import (
	"../protos"
	"context"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

func main() {
	// Create multiple clients and start receiving data
	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		client, err := mkLonglivedClient(int32(i))
		if err != nil {
			// TODO handle
			log.Fatal(err)
		}
		go client.start()
		time.Sleep(time.Second*2)
	}

	// The wait group purpose is to avoid exiting, the clients do not exit
	wg.Wait()
}

type longlivedClient struct {
	client protos.LonglivedClient           // client is the gRPC client
	id     int32                            // id is the client ID used for subscribing
	conn   *grpc.ClientConn                 // conn is the client connection
	stream protos.Longlived_SubscribeClient // stream holds the gRPC stream between the client and the server
}

// mkLonglivedClient returns a new gRPC client instance
func mkLonglivedClient(id int32) (*longlivedClient, error) {
	conn, err := mkConnection()
	if err != nil {
		return nil, err
	}
	return &longlivedClient{
		client: protos.NewLonglivedClient(conn),
		conn:   conn,
		id:     id,
	}, nil
}

// close is not used but is here as an example of how to close the gRPC client connection
func (c *longlivedClient) close() {
	if err := c.conn.Close(); err != nil {
		log.Fatal(err)
	}
}

// subscribe subscribes to messages from the gRPC server
func (c *longlivedClient) subscribe() (protos.Longlived_SubscribeClient, error) {
	return c.client.Subscribe(context.Background(), &protos.Request{Id: c.id})
}

func (c *longlivedClient) start() {
	log.Printf("Subscribing client ID: %d", c.id)
	var err error
	for {
		if c.stream == nil {
			if c.stream, err = c.subscribe(); err != nil {
				log.Printf("Failed to subscribe: %v", err)
				c.sleep()
				continue
			}
		}
		response, err := c.stream.Recv()
		if err != nil {
			log.Printf("Failed to recieve message: %v", err)
			// Clearing the stream will force the client to resubscribe on next iteration
			c.stream = nil
			c.sleep()
			continue
		}
		log.Printf("Client ID %v got response: %v", c.id, response.Data)
	}
}

// sleep is used to give the server time to unsubscribe the client and reset the stream
func (c *longlivedClient) sleep() {
	time.Sleep(time.Second * 5)
}

func mkConnection() (*grpc.ClientConn, error) {
	return grpc.Dial("127.0.0.1:7070", []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}...)
}
