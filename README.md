# gRPC Long-lived Streaming

This repository holds a minimalistic example of a gRPC long lived streaming application.

To learn more, visit the blog post: https://dev.bitolog.com/grpc-long-lived-streaming/



## Instructions

To compile the proto file, run the following command from the `protos` folder:

```
$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative longlived.proto
```

Note that this was tested on protoc version: `libprotoc 3.6.1` 



## Running the server

Navigate to `/server` folder and run the following:

```
$ go build server.go
$ ./server
2021/03/04 08:48:09 Starting server on address 127.0.0.1:7070
2021/03/04 08:48:09 Starting data generation
```



## Running the client(s)

The client process emulates several clients (default is 10).

Once the server is up, navigate to `/client` and run the following:

```
$ go build client.go
$ ./client
2021/03/04 09:19:29 Subscribing client ID: 1
 2021/03/04 09:19:29 Client ID 1 got response: "data mock for: 1"
 2021/03/04 09:19:30 Client ID 1 got response: "data mock for: 1"
 2021/03/04 09:19:31 Subscribing client ID: 2
 2021/03/04 09:19:31 Client ID 1 got response: "data mock for: 1"
 2021/03/04 09:19:31 Client ID 2 got response: "data mock for: 2"
 2021/03/04 09:19:32 Client ID 1 got response: "data mock for: 1"
 2021/03/04 09:19:32 Client ID 2 got response: "data mock for: 2"
 2021/03/04 09:19:33 Subscribing client ID: 3
 2021/03/04 09:19:33 Client ID 2 got response: "data mock for: 2"
 2021/03/04 09:19:33 Client ID 1 got response: "data mock for: 1"
 2021/03/04 09:19:33 Client ID 3 got response: "data mock for: 3"
 2021/03/04 09:19:34 Client ID 1 got response: "data mock for: 1"
 2021/03/04 09:19:34 Client ID 3 got response: "data mock for: 3"
 2021/03/04 09:19:34 Client ID 2 got response: "data mock for: 2"
 2021/03/04 09:19:35 Subscribing client ID: 4
 2021/03/04 09:19:35 Client ID 2 got response: "data mock for: 2"
 2021/03/04 09:19:35 Client ID 4 got response: "data mock for: 4"
```

