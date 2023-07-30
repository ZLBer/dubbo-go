package resource

import "google.golang.org/protobuf/types/known/anypb"

// EndpointsUpdate contains an EDS update.
type DubboServiceNameMappingUpdate struct {
	Namespace        string
	InterfaceName    string
	ApplicationNames []string
	// Raw is the resource from the xds response.
	Raw *anypb.Any
}

type DubboServiceNameMappingTypeErrTuple struct {
	Update DubboServiceNameMappingUpdate
	Err    error
}
