# gRPC Long-lived Streaming

This repository holds a minimalistic example of a gRPC long lived stream.

To learn more, visit the blog post: 



## Instructions

To compile the proto file, run the following command from the `protos` folder:

```
$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative longlived.proto
```

Note that this was tested on protoc version: `libprotoc 3.6.1` 

