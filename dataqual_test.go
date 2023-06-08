package dataqual_wasm

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	kafka "github.com/streamdal/segmentio-kafka-go"
	sarama "github.com/streamdal/shopify-sarama"

	"github.com/pkg/errors"
	"github.com/streamdal/dataqual"
	"google.golang.org/grpc"

	"github.com/batchcorp/plumber-schemas/build/go/protos"
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
}

func BenchmarkRabbit_Produce(b *testing.B) {
	b.Run("MatchRuleAndDiscard", func(b *testing.B) {
		os.Setenv("PLUMBER_URL", "localhost:9090")
		os.Setenv("PLUMBER_TOKEN", "test")
		r, err := setupRabbit()
		if err != nil {
			b.Error(err.Error())
		}
		b.Cleanup(func() { r.Cleanup() })

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			r.Produce([]byte(`{"type": "hello world@gmail.com"}`))
		}
	})

	b.Run("NoRuleMatch", func(b *testing.B) {
		time.Sleep(time.Second * 3) // Give time for k.Cleanup() to finish
		os.Setenv("PLUMBER_URL", "localhost:9090")
		os.Setenv("PLUMBER_TOKEN", "test")

		r, err := setupRabbit()
		if err != nil {
			panic(err)
		}
		b.Cleanup(func() { r.Cleanup() })

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			r.Produce([]byte(`{"type": "hello world@gmail"}`))
		}
	})

	b.Run("NoDataqual", func(b *testing.B) {
		time.Sleep(time.Second * 3) // Give time for k.Cleanup() to finish
		os.Setenv("PLUMBER_URL", "")
		os.Setenv("PLUMBER_TOKEN", "")

		r, err := setupRabbit()
		if err != nil {
			panic(err)
		}
		b.Cleanup(func() { r.Cleanup() })

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			r.Produce([]byte(`{"type": "hello world@gmail.com"}`))
		}
	})
}

func BenchmarkSegmentio_Produce(b *testing.B) {
	b.SetParallelism(0)

	b.Run("MatchRuleAndDiscard", func(b *testing.B) {
		os.Setenv("PLUMBER_URL", "localhost:9090")
		os.Setenv("PLUMBER_TOKEN", "test")

		k, err := setupSegmentio(KafkaTopic)
		if err != nil {
			b.Error(err)
		}
		b.Cleanup(func() { k.Cleanup() })

		time.Sleep(time.Second * 3)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			if err := k.Writer.WriteMessages(context.Background(), kafka.Message{
				Topic: KafkaTopic,
				Value: []byte(`{"type": "hello world@gmail.com"}`),
			}); err != nil {
				if !strings.Contains(err.Error(), dataqual.ErrMessageDropped.Error()) {
					b.Error(err.Error())
				}
			}
		}
	})

	b.Run("NoRuleMatch", func(b *testing.B) {
		os.Setenv("PLUMBER_URL", "localhost:9090")
		os.Setenv("PLUMBER_TOKEN", "test")

		k, err := setupSegmentio(KafkaTopic)
		if err != nil {
			b.Error(err)
		}
		b.Cleanup(func() { k.Cleanup() })

		time.Sleep(time.Second * 3)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			if err := k.Writer.WriteMessages(context.Background(), kafka.Message{
				Topic: KafkaTopic,
				Value: []byte(`{"type": "hello world@gmail"}`),
			}); err != nil {
				if !strings.Contains(err.Error(), dataqual.ErrMessageDropped.Error()) {
					b.Error(err.Error())
				}
			}
		}
	})

	b.Run("NoDataqual", func(b *testing.B) {
		os.Setenv("PLUMBER_URL", "")
		os.Setenv("PLUMBER_TOKEN", "")

		k, err := setupSegmentio(KafkaTopic)
		if err != nil {
			b.Error(err)
		}
		b.Cleanup(func() { k.Cleanup() })

		time.Sleep(time.Second * 3)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			if err := k.Writer.WriteMessages(context.Background(), kafka.Message{
				Topic: KafkaTopic,
				Value: []byte(`{"type": "hello world@gmail.com"}`),
			}); err != nil {
				if !strings.Contains(err.Error(), dataqual.ErrMessageDropped.Error()) {
					b.Error(err.Error())
				}
			}
		}
	})
}

func BenchmarkProduce_Sarama(b *testing.B) {
	b.Run("MatchRuleAndDiscard", func(b *testing.B) {
		os.Setenv("PLUMBER_URL", "localhost:9090")
		os.Setenv("PLUMBER_TOKEN", "test")

		k, err := setupSarama(KafkaTopic)
		if err != nil {
			b.Error(err)
		}
		b.Cleanup(func() { k.Cleanup() })

		time.Sleep(time.Second * 3)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			if _, _, err := k.Producer.SendMessage(&sarama.ProducerMessage{
				Topic: KafkaTopic,
				Value: sarama.StringEncoder(`{"type": "hello world@gmail.com"}`),
			}); err != nil {
				if err != dataqual.ErrMessageDropped {
					b.Error(err.Error())
				}
			}
		}
	})

	b.Run("NoRuleMatch", func(b *testing.B) {
		os.Setenv("PLUMBER_URL", "localhost:9090")
		os.Setenv("PLUMBER_TOKEN", "test")

		k, err := setupSarama(KafkaTopic)
		if err != nil {
			b.Error(err)
		}
		b.Cleanup(func() { k.Cleanup() })

		time.Sleep(time.Second * 3)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			if _, _, err := k.Producer.SendMessage(&sarama.ProducerMessage{
				Topic: KafkaTopic,
				Value: sarama.StringEncoder(`{"type": "hello world@gmail"}`),
			}); err != nil {
				if err != dataqual.ErrMessageDropped {
					b.Error(err.Error())
				}
			}
		}
	})

	b.Run("NoDataqual", func(b *testing.B) {
		os.Setenv("PLUMBER_URL", "")
		os.Setenv("PLUMBER_TOKEN", "")

		k, err := setupSarama(KafkaTopic)
		if err != nil {
			b.Error(err)
		}
		b.Cleanup(func() { k.Cleanup() })

		time.Sleep(time.Second * 3)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			if _, _, err := k.Producer.SendMessage(&sarama.ProducerMessage{
				Topic: KafkaTopic,
				Value: sarama.StringEncoder(`{"type": "hello world@gmail.com"}`),
			}); err != nil {
				if err != dataqual.ErrMessageDropped {
					b.Error(err.Error())
				}
			}
		}
	})
}
