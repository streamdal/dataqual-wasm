package dataqual_wasm

import (
	"bytes"
	"compress/gzip"
	"context"
	"net"
	"os"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/batchcorp/plumber-schemas/build/go/protos"
	"github.com/batchcorp/plumber-schemas/build/go/protos/common"
)

type MockPlumberServer struct {
	server *grpc.Server
	lis    net.Listener
}

// compress data using gzip. Used before uploading WASM files to plumber server
func compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	if _, err := gz.Write(data); err != nil {
		return nil, errors.Wrap(err, "unable to write to gzip writer")
	}

	if err := gz.Flush(); err != nil {
		return nil, errors.Wrap(err, "unable to flush gzip writer")
	}

	if err := gz.Close(); err != nil {
		return nil, errors.Wrap(err, "unable to close gzip writer")
	}

	return b.Bytes(), nil
}

func (m MockPlumberServer) GetAllConnections(ctx context.Context, request *protos.GetAllConnectionsRequest) (*protos.GetAllConnectionsResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) GetConnection(ctx context.Context, request *protos.GetConnectionRequest) (*protos.GetConnectionResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) CreateConnection(ctx context.Context, request *protos.CreateConnectionRequest) (*protos.CreateConnectionResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) TestConnection(ctx context.Context, request *protos.TestConnectionRequest) (*protos.TestConnectionResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) UpdateConnection(ctx context.Context, request *protos.UpdateConnectionRequest) (*protos.UpdateConnectionResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) DeleteConnection(ctx context.Context, request *protos.DeleteConnectionRequest) (*protos.DeleteConnectionResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) GetAllRelays(ctx context.Context, request *protos.GetAllRelaysRequest) (*protos.GetAllRelaysResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) GetRelay(ctx context.Context, request *protos.GetRelayRequest) (*protos.GetRelayResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) CreateRelay(ctx context.Context, request *protos.CreateRelayRequest) (*protos.CreateRelayResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) UpdateRelay(ctx context.Context, request *protos.UpdateRelayRequest) (*protos.UpdateRelayResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) ResumeRelay(ctx context.Context, request *protos.ResumeRelayRequest) (*protos.ResumeRelayResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) StopRelay(ctx context.Context, request *protos.StopRelayRequest) (*protos.StopRelayResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) DeleteRelay(ctx context.Context, request *protos.DeleteRelayRequest) (*protos.DeleteRelayResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) GetTunnel(ctx context.Context, request *protos.GetTunnelRequest) (*protos.GetTunnelResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) GetAllTunnels(ctx context.Context, request *protos.GetAllTunnelsRequest) (*protos.GetAllTunnelsResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) CreateTunnel(ctx context.Context, request *protos.CreateTunnelRequest) (*protos.CreateTunnelResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) StopTunnel(ctx context.Context, request *protos.StopTunnelRequest) (*protos.StopTunnelResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) ResumeTunnel(ctx context.Context, request *protos.ResumeTunnelRequest) (*protos.ResumeTunnelResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) UpdateTunnel(ctx context.Context, request *protos.UpdateTunnelRequest) (*protos.UpdateTunnelResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) DeleteTunnel(ctx context.Context, request *protos.DeleteTunnelRequest) (*protos.DeleteTunnelResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) ListWasmFiles(ctx context.Context, request *protos.ListWasmFilesRequest) (*protos.ListWasmFilesResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) UploadWasmFile(ctx context.Context, request *protos.UploadWasmFileRequest) (*protos.UploadWasmFileResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) DownloadWasmFile(ctx context.Context, request *protos.DownloadWasmFileRequest) (*protos.DownloadWasmFileResponse, error) {
	var data []byte
	var err error

	switch request.Name {
	case "transform.wasm":
		data, err = os.ReadFile("build/transform.wasm")
		if err != nil {
			return nil, err
		}
	case "match.wasm":
		data, err = os.ReadFile("build/match.wasm")
		if err != nil {
			return nil, err
		}
	}

	compressed, err := compress(data)
	if err != nil {
		return nil, err
	}

	return &protos.DownloadWasmFileResponse{
		Status: &common.Status{
			Code:      common.Code_OK,
			RequestId: uuid.New().String(),
		},
		Data: compressed,
	}, nil
}

func (m MockPlumberServer) DeleteWasmFile(ctx context.Context, request *protos.DeleteWasmFileRequest) (*protos.DeleteWasmFileResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) GetRuleSets(ctx context.Context, request *protos.GetDataQualityRuleSetsRequest) (*protos.GetDataQualityRuleSetsResponse, error) {
	kafkaRuleID1 := uuid.New().String()
	kafkaRuleID2 := uuid.New().String()
	rabbitRuleID1 := uuid.New().String()
	rabbitRuleID2 := uuid.New().String()

	return &protos.GetDataQualityRuleSetsResponse{
		RuleSets: []*common.RuleSet{
			{
				Name: "kafka publish rules",
				Key:  "dqtest1",
				Bus:  "kafka",
				Mode: common.RuleMode_RULE_MODE_PUBLISH,
				Rules: map[string]*common.Rule{
					kafkaRuleID1: {
						Id:   kafkaRuleID1,
						Type: common.RuleType_RULE_TYPE_MATCH,
						RuleConfig: &common.Rule_MatchConfig{
							MatchConfig: &common.RuleConfigMatch{
								Path: "type",
								Type: "string_contains_any",
								Args: []string{"gmail.com"},
							},
						},
						FailureMode:       common.RuleFailureMode_RULE_FAILURE_MODE_REJECT,
						FailureModeConfig: &common.Rule_Reject{},
					},
				},
			},
			{
				Name: "kafka consume rules",
				Key:  "dqtest1",
				Bus:  "kafka",
				Mode: common.RuleMode_RULE_MODE_CONSUME,
				Rules: map[string]*common.Rule{
					kafkaRuleID2: {
						Id:   kafkaRuleID2,
						Type: common.RuleType_RULE_TYPE_MATCH,
						RuleConfig: &common.Rule_MatchConfig{
							MatchConfig: &common.RuleConfigMatch{
								Path: "type",
								Type: "string_contains_any",
								Args: []string{"gmail.com"},
							},
						},
						FailureMode:       common.RuleFailureMode_RULE_FAILURE_MODE_REJECT,
						FailureModeConfig: &common.Rule_Reject{},
					},
				},
			},
			{
				Name: "rabbit publish rules",
				Key:  "dqtest|dqtest",
				Bus:  "rabbitmq",
				Mode: common.RuleMode_RULE_MODE_PUBLISH,
				Rules: map[string]*common.Rule{
					rabbitRuleID1: {
						Id:   rabbitRuleID1,
						Type: common.RuleType_RULE_TYPE_MATCH,
						RuleConfig: &common.Rule_MatchConfig{
							MatchConfig: &common.RuleConfigMatch{
								Path: "type",
								Type: "string_contains_any",
								Args: []string{"gmail.com"},
							},
						},
						FailureMode:       common.RuleFailureMode_RULE_FAILURE_MODE_REJECT,
						FailureModeConfig: &common.Rule_Reject{},
					},
				},
			},
			{
				Name: "rabbit consume rules",
				Key:  "dqtest",
				Bus:  "rabbitmq",
				Mode: common.RuleMode_RULE_MODE_CONSUME,
				Rules: map[string]*common.Rule{
					rabbitRuleID2: {
						Id:   rabbitRuleID2,
						Type: common.RuleType_RULE_TYPE_MATCH,
						RuleConfig: &common.Rule_MatchConfig{
							MatchConfig: &common.RuleConfigMatch{
								Path: "type",
								Type: "string_contains_any",
								Args: []string{"gmail.com"},
							},
						},
						FailureMode:       common.RuleFailureMode_RULE_FAILURE_MODE_REJECT,
						FailureModeConfig: &common.Rule_Reject{},
					},
				},
			},
		},
	}, nil
}

func (m MockPlumberServer) GetRule(ctx context.Context, request *protos.GetDataQualityRuleRequest) (*protos.GetDataQualityRuleResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) CreateRule(ctx context.Context, request *protos.CreateDataQualityRuleRequest) (*protos.CreateDataQualityRuleResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) UpdateRule(ctx context.Context, request *protos.UpdateDataQualityRuleRequest) (*protos.UpdateDataQualityRuleResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) DeleteRule(ctx context.Context, request *protos.DeleteDataQualityRuleRequest) (*protos.DeleteDataQualityRuleResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) CreateRuleSet(ctx context.Context, request *protos.CreateDataQualityRuleSetRequest) (*protos.CreateDataQualityRuleSetResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) UpdateRuleSet(ctx context.Context, request *protos.UpdateDataQualityRuleSetRequest) (*protos.UpdateDataQualityRuleSetResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) DeleteRuleSet(ctx context.Context, request *protos.DeleteDataQualityRuleSetRequest) (*protos.DeleteDataQualityRuleSetResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) SendRuleNotification(ctx context.Context, request *protos.SendRuleNotificationRequest) (*protos.SendRuleNotificationResponse, error) {
	panic("implement me")
}

func (m MockPlumberServer) GetServerOptions(ctx context.Context, request *protos.GetServerOptionsRequest) (*protos.GetServerOptionsResponse, error) {
	panic("implement me")
}
