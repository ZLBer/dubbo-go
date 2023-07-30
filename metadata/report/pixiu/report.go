package pixiu

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/metadata/identifier"
	"dubbo.apache.org/dubbo-go/v3/metadata/report"
	"dubbo.apache.org/dubbo-go/v3/registry"
	"encoding/json"
	gxset "github.com/dubbogo/gost/container/set"
	"google.golang.org/grpc"
	"istio.io/api/dubbo/v1alpha1"
)

type pixiuMetadataReportFactory struct {
}

// CreateMetadataReport create a new metadata report
func (mf *pixiuMetadataReportFactory) CreateMetadataReport(url *common.URL) report.MetadataReport {

	conn, err := grpc.Dial(url.Location)
	if err != nil {
		panic(err)
	}

	snpClient := v1alpha1.NewServiceNameMappingServiceClient(conn)

	metaClient := v1alpha1.NewServiceMetadataServiceClient(conn)

	return &pixiuMetadataReport{snpClient: snpClient, metaClient: metaClient}
}

type pixiuMetadataReport struct {
	snpClient  v1alpha1.ServiceNameMappingServiceClient
	metaClient v1alpha1.ServiceMetadataServiceClient
}

func (p pixiuMetadataReport) StoreProviderMetadata(providerIdentifier *identifier.MetadataIdentifier, serviceDefinitions string) error {
	panic("implement me")
}

func (p pixiuMetadataReport) StoreConsumerMetadata(metadataIdentifier *identifier.MetadataIdentifier, s string) error {
	panic("implement me")
}

func (p pixiuMetadataReport) SaveServiceMetadata(metadataIdentifier *identifier.ServiceMetadataIdentifier, url *common.URL) error {
	panic("implement me")
}

func (p pixiuMetadataReport) RemoveServiceMetadata(metadataIdentifier *identifier.ServiceMetadataIdentifier) error {
	panic("implement me")
}

func (p pixiuMetadataReport) GetExportedURLs(metadataIdentifier *identifier.ServiceMetadataIdentifier) ([]string, error) {
	panic("implement me")
}

func (p pixiuMetadataReport) SaveSubscribedData(metadataIdentifier *identifier.SubscriberMetadataIdentifier, s string) error {
	panic("implement me")
}

func (p pixiuMetadataReport) GetSubscribedURLs(metadataIdentifier *identifier.SubscriberMetadataIdentifier) ([]string, error) {
	panic("implement me")
}

func (p pixiuMetadataReport) GetServiceDefinition(metadataIdentifier *identifier.MetadataIdentifier) (string, error) {
	panic("implement me")
}

func (p pixiuMetadataReport) GetAppMetadata(metadataIdentifier *identifier.SubscriberMetadataIdentifier) (*common.MetadataInfo, error) {
	response, err := p.metaClient.Get(context.Background(), &v1alpha1.GetServiceMetadataRequest{
		Namespace:       metadataIdentifier.Group,
		ApplicationName: "",
		Revision:        metadataIdentifier.Revision,
	})
	if err != nil {
		return nil, err
	}
	data := response.GetMetadataInfo()
	var metadataInfo common.MetadataInfo
	err = json.Unmarshal([]byte(data), &metadataInfo)
	if err != nil {
		return nil, err
	}
	return &metadataInfo, nil
}

func (p pixiuMetadataReport) PublishAppMetadata(metadataIdentifier *identifier.SubscriberMetadataIdentifier, info *common.MetadataInfo) error {

	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	_, err = p.metaClient.Publish(context.Background(), &v1alpha1.PublishServiceMetadataRequest{
		Namespace:       metadataIdentifier.Group,
		ApplicationName: metadataIdentifier.Application,
		Revision:        metadataIdentifier.Revision,
		MetadataInfo:    string(data),
	})
	if err != nil {
		return err
	}
	return err
}

func (p pixiuMetadataReport) RegisterServiceAppMapping(serviceInterface string, group string, appName string) error {
	_, err := p.snpClient.RegisterServiceAppMapping(context.Background(), &v1alpha1.ServiceMappingRequest{
		Namespace:       group,
		ApplicationName: appName,
		InterfaceNames:  []string{serviceInterface},
	})
	if err != nil {
		return err
	}
	return nil
}

func (p pixiuMetadataReport) GetServiceAppMapping(key string, group string, listener registry.MappingListener) (*gxset.HashSet, error) {
	panic("implement me")
}

func (p pixiuMetadataReport) RemoveServiceAppMappingListener(key string, group string) error {
	panic("implement me")
}
