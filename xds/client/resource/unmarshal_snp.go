package resource

import (
	"fmt"
	dubbogoLogger "github.com/dubbogo/gost/log/logger"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"istio.io/api/dubbo/v1alpha1"
)

// UnmarshalEndpoints processes resources received in an EDS response,
// validates them, and transforms them into a native struct which contains only
// fields we are interested in.
func UnmarshalDubboServiceNameMapping(opts *UnmarshalOptions) (map[string]DubboServiceNameMappingTypeErrTuple, UpdateMetadata, error) {
	update := make(map[string]DubboServiceNameMappingTypeErrTuple)
	md, err := processAllResources(opts, update)
	return update, md, err
}

func unmarshalDubboServiceNameMapping(r *anypb.Any, logger dubbogoLogger.Logger) (string, DubboServiceNameMappingUpdate, error) {
	if !IsDubboServiceNameMappingResource(r.GetTypeUrl()) {
		return "", DubboServiceNameMappingUpdate{}, fmt.Errorf("unexpected resource type: %q ", r.GetTypeUrl())
	}

	snp := &v1alpha1.ServiceMappingXdsResponse{}
	if err := proto.Unmarshal(r.GetValue(), snp); err != nil {
		return "", DubboServiceNameMappingUpdate{}, fmt.Errorf("failed to unmarshal resource: %v", err)
	}

	return snp.InterfaceName, DubboServiceNameMappingUpdate{
		Namespace:        snp.GetNamespace(),
		InterfaceName:    snp.GetInterfaceName(),
		ApplicationNames: snp.GetApplicationNames(),
		Raw:              r,
	}, nil
}
