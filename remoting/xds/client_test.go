package xds

import (
	"context"
	"fmt"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	_struct "github.com/golang/protobuf/ptypes/struct"
	"time"

	envoy_service_discovery_v3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
	"istio.io/api/dubbo/v1alpha1"
	"log"
	"testing"
)

const (
	xdsServerAddress = "localhost:15010" // 将此地址替换为实际的xDS服务器地址
)

func sendXdsRequest(stream envoy_service_discovery_v3.AggregatedDiscoveryService_StreamAggregatedResourcesClient, resourceType, resourceName string, node *envoy_config_core_v3.Node) error {
	req := &envoy_service_discovery_v3.DiscoveryRequest{
		TypeUrl:       resourceType,
		ResourceNames: []string{resourceName},
		ResponseNonce: time.Now().String(),
		Node:          node,
	}

	return stream.Send(req)
}
func TestXDS(t *testing.T) {
	conn, err := grpc.Dial(xdsServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to xDS server: %v", err)
	}
	defer conn.Close()

	adsClient := envoy_service_discovery_v3.NewAggregatedDiscoveryServiceClient(conn)
	stream, err := adsClient.StreamAggregatedResources(context.Background())
	if err != nil {
		log.Fatalf("Failed to open ADS stream: %v", err)
	}
	//sidecar~ip~{POD_NAME}~{NAMESPACE_NAME}.svc.cluster.local
	node := &envoy_config_core_v3.Node{
		Id:      "sidecar~127.0.0.1~xds_client~default.svc.cluster.local",
		Cluster: "default",
		Metadata: &_struct.Struct{
			Fields: map[string]*structpb.Value{
				"env": {
					Kind: &structpb.Value_StringValue{
						StringValue: "test",
					},
				},
			},
		},
	}

	//发送请求
	err = sendXdsRequest(stream, "dubbo.networking.v1alpha1.v1.servicenamemapping", "a|default", node)
	if err != nil {
		log.Fatalf("Failed to send Listener request: %v", err)
	}

	//err = sendXdsRequest(stream, resource.ClusterType, "cluster", node)
	//if err != nil {
	//	log.Fatalf("Failed to send Cluster request: %v", err)
	//}

	// 接收响应
	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Fatalf("Failed to receive xDS response: %v", err)
		}
		log.Printf("Received xDS response: %v", resp)
	}
}

func TestClient(t *testing.T) {

	dial, err := grpc.Dial("127.0.0.1:15010", grpc.WithTransportCredentials(insecure.NewCredentials()))

	fmt.Println(dial)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := v1alpha1.NewServiceNameMappingServiceClient(dial)

	mapping, err := client.RegisterServiceAppMapping(context.Background(), &v1alpha1.ServiceMappingRequest{
		Namespace:       "default",
		ApplicationName: "application-05",
		InterfaceNames:  []string{"a"},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("res:", mapping)

}
