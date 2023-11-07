package grpcx

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GrpcConnect establishes a gRPC connection to the specified address.
//
// Parameters:
//   - connectAddress: the address to connect to.
//
// Returns: a pointer to a gRPC client connection.
func GrpcConnect(connectAddress string) *grpc.ClientConn {
	conn, err := grpc.Dial(connectAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(fmt.Printf("grpc connect failed: %s", err.Error()))
	}
	return conn
}
