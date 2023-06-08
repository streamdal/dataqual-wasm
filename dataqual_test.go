package dataqual_wasm

import (
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/batchcorp/plumber-schemas/build/go/protos"
)

const (
	AmqpExchange   = "dqtest"
	AmqpRoutingKey = "dqtest"
	AmqpQueue      = "dqtest"
)

func startGRPCServer() (*MockPlumberServer, error) {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		return nil, fmt.Errorf("unable to listen on '%s': %s", ":9090", err)
	}
	grpcServer := grpc.NewServer()

	plumberServer := &MockPlumberServer{}
	plumberServer.lis = lis
	plumberServer.server = grpcServer

	protos.RegisterPlumberServerServer(grpcServer, plumberServer)

	errCh := make(chan error, 1)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			errCh <- errors.Wrap(err, "unable to start gRPC server")
		}
	}()

	afterCh := time.After(5 * time.Second)

	select {
	case <-afterCh:
		return plumberServer, nil
	case err := <-errCh:
		return nil, err
	}
}

var srv *MockPlumberServer

func init() {
	grpcSrv, err := startGRPCServer()
	if err != nil {
		panic(err)
	}

	srv = grpcSrv
	println("ran init")
}

func BenchmarkProduce_MatchRuleAndDiscard(b *testing.B) {
	os.Setenv("PLUMBER_URL", "localhost:9090")
	os.Setenv("PLUMBER_TOKEN", "test")

	r, err := setupRabbit()
	if err != nil {
		panic(err)
	}
	defer r.Cleanup()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.Produce([]byte(`{"type": "hello world@gmail.com"}`))
	}
}

func BenchmarkProduce_NoRuleMatch(b *testing.B) {
	os.Setenv("PLUMBER_URL", "localhost:9090")
	os.Setenv("PLUMBER_TOKEN", "test")

	r, err := setupRabbit()
	if err != nil {
		panic(err)
	}
	defer r.Cleanup()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.Produce([]byte(`{"type": "hello world@gmail"}`))
	}
}

func BenchmarkProduce_NoDataqual(b *testing.B) {
	os.Setenv("PLUMBER_URL", "")
	os.Setenv("PLUMBER_TOKEN", "")

	r, err := setupRabbit()
	if err != nil {
		panic(err)
	}
	defer r.Cleanup()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.Produce([]byte(`{"type": "hello world@gmail"}`))
	}
}
