package entity

type TfvarResourceGroupType struct {
	Location string `json:"location"`
}

type TfvarDefaultNodePoolType struct {
	EnableAutoScaling bool `json:"enableAutoScaling"`
	MinCount          int  `json:"minCount"`
	MaxCount          int  `json:"maxCount"`
}

type TfvarAddonsType struct {
	AppGateway        bool `json:"appGateway"`
	MicrosoftDefender bool `json:"microsoftDefender"`
}

type TfvarKubernetesClusterType struct {
	KubernetesVersion     string                   `json:"kubernetesVersion"`
	NetworkPlugin         string                   `json:"networkPlugin"`
	NetworkPolicy         string                   `json:"networkPolicy"`
	NetworkPluginMode     string                   `json:"networkPluginMode"`
	OutboundType          string                   `json:"outboundType"`
	PrivateClusterEnabled string                   `json:"privateClusterEnabled"`
	Addons                TfvarAddonsType          `json:"addons"`
	DefaultNodePool       TfvarDefaultNodePoolType `json:"defaultNodePool"`
}

type TfvarVirtualNeworkType struct {
	AddressSpace []string
}

type TfvarSubnetType struct {
	Name            string
	AddressPrefixes []string
}

type TfvarNetworkSecurityGroupType struct {
}

type TfvarJumpserverType struct {
	AdminPassword string `json:"adminPassword"`
	AdminUserName string `json:"adminUsername"`
}

type TfvarFirewallType struct {
	SkuName string `json:"skuName"`
	SkuTier string `json:"skuTier"`
}

type ContainerRegistryType struct {
}

type AppGatewayType struct{}

type TfvarConfigType struct {
	ResourceGroup         TfvarResourceGroupType          `json:"resourceGroup"`
	VirtualNetworks       []TfvarVirtualNeworkType        `json:"virtualNetworks"`
	Subnets               []TfvarSubnetType               `json:"subnets"`
	Jumpservers           []TfvarJumpserverType           `json:"jumpservers"`
	NetworkSecurityGroups []TfvarNetworkSecurityGroupType `json:"networkSecurityGroups"`
	KubernetesClusters    []TfvarKubernetesClusterType    `json:"kubernetesClusters"`
	Firewalls             []TfvarFirewallType             `json:"firewalls"`
	ContainerRegistries   []ContainerRegistryType         `json:"containerRegistries"`
	AppGateways           []AppGatewayType                `json:"appGateways"`
}

type Blob struct {
	Name string `xml:"Name" json:"name"`
	//Url  string `xml:"Url" json:"url"`
}

// Ok. if you noted that the its named blob and should be Blobs. I've no idea whose fault is this.
// Read more about the API https://learn.microsoft.com/en-us/rest/api/storageservices/list-blobs?tabs=azure-ad#request
type Blobs struct {
	Blob []Blob `xml:"Blob" json:"blob"`
}

type EnumerationResults struct {
	Blobs Blobs `xml:"Blobs" json:"blobs"`
}

type LabType struct {
	Id           string          `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Tags         []string        `json:"tags"`
	Template     TfvarConfigType `json:"template"`
	ExtendScript string          `json:"extendScript"`
	Message      string          `json:"message"`
	Type         string          `json:"type"`
	CreatedBy    string          `json:"createdBy"`
	CreatedOn    string          `json:"createdOn"`
	UpdatedBy    string          `json:"updatedBy"`
	UpdatedOn    string          `json:"updatedOn"`
}

type BlobType struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LabService interface {
	// Public Labs
	// Includes = sharedlabs, labexercises, mockcases.
	GetPublicLabs(typeOfLab string) ([]LabType, error)
	AddPublicLab(LabType) error
	DeletePublicLab(LabType) error
}

type LabRepository interface {

	// Public labs
	GetEnumerationResults(typeOfLab string) (EnumerationResults, error)
	GetLab(name string, typeOfLab string) (LabType, error)
	AddLab(labId string, lab string, typeOfLab string) error
	DeleteLab(labId string, typeOfLab string) error
}
