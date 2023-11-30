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

type TfvarVirtualNetworkType struct {
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
	VirtualNetworks       []TfvarVirtualNetworkType       `json:"virtualNetworks"`
	Subnets               []TfvarSubnetType               `json:"subnets"`
	Jumpservers           []TfvarJumpserverType           `json:"jumpservers"`
	NetworkSecurityGroups []TfvarNetworkSecurityGroupType `json:"networkSecurityGroups"`
	KubernetesClusters    []TfvarKubernetesClusterType    `json:"kubernetesClusters"`
	Firewalls             []TfvarFirewallType             `json:"firewalls"`
	ContainerRegistries   []ContainerRegistryType         `json:"containerRegistries"`
	AppGateways           []AppGatewayType                `json:"appGateways"`
}

type Blob struct {
	Name             string `xml:"Name" json:"name"`
	VersionId        string `xml:"VersionId" json:"versionId"`
	IsCurrentVersion bool   `xml:"IsCurrentVersion" json:"isCurrentVersion"`
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

var (
	PrivateLab    = []string{"privatelab", "challengelab"}
	PublicLab     = []string{"publiclab"}
	ProtectedLabs = []string{"readinesslab", "mockcase"}
)

type LabType struct {
	Id               string          `json:"id"`
	Name             string          `json:"name"`
	Description      string          `json:"description"`
	Tags             []string        `json:"tags"`
	Template         TfvarConfigType `json:"template"`
	ExtendScript     string          `json:"extendScript"`
	Message          string          `json:"message"`
	Category         string          `json:"category"`
	Type             string          `json:"type"`
	CreatedBy        string          `json:"createdBy"`
	CreatedOn        string          `json:"createdOn"`
	UpdatedBy        string          `json:"updatedBy"`
	UpdatedOn        string          `json:"updatedOn"`
	Owners           []string        `json:"owners"`
	Editors          []string        `json:"editors"`
	Viewers          []string        `json:"viewers"`
	VersionId        string          `json:"versionId"`
	IsCurrentVersion bool            `json:"isCurrentVersion"`
}

type BlobType struct {
	Name      string `json:"name"`
	Url       string `json:"url"`
	VersionId string `json:"versionId"`
}

type LabService interface {
	// Private Labs
	// Role: user
	// Types: privatelab, challengelab
	GetPrivateLab(typeOfLab string, userId string) ([]LabType, error)
	GetPrivateLabVersions(typeOfLab string, labId string, userId string) ([]LabType, error)
	UpsertPrivateLab(LabType) error
	DeletePrivateLab(typeOfLab string, labId string, userId string) error

	// Public Labs
	// Role: user
	// Types: publiclab
	GetPublicLab(typeOfLab string) ([]LabType, error)
	GetPublicLabVersions(typeOfLab string, labId string) ([]LabType, error)
	UpsertPublicLab(LabType) error
	DeletePublicLab(typeOfLab string, labId string, userId string) error

	// Protected Labs
	// Role: mentor
	// Types: readinesslab, mockcase
	GetProtectedLabs(typeOfLab string) ([]LabType, error)
	GetProtectedLabVersions(typeOfLab string, labId string) ([]LabType, error)
	UpsertProtectedLab(LabType) error
	DeleteProtectedLab(typeOfLab string, labId string) error

	// Shared functions
	GetLabs(typeOfLab string) ([]LabType, error)
	GetLabVersions(typeOfLab string, labId string) ([]LabType, error)
	UpsertLab(LabType) error
	DeleteLab(typeOfLab string, labId string) error
}

type LabRepository interface {

	// Public labs
	GetEnumerationResults(typeOfLab string, includeVersions bool) (EnumerationResults, error)
	GetLab(typeOfLab string, labId string, versionId string) (LabType, error) // send empty versionId ("") to get current version.
	UpsertLab(labId string, lab string, typeOfLab string) error
	DeleteLab(typeOfLab string, labId string) error
}
