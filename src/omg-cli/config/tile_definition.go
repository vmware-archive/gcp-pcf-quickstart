package config

type PivnetMetadata struct {
	Name      string
	VersionId string
	FileId    string
	Sha256    string
}

type StemcellMetadata struct {
	PivnetMetadata
	StemcellName string
}

type OpsManagerMetadata struct {
	Name    string
	Version string
}

type Tile struct {
	Pivnet   PivnetMetadata
	Product  OpsManagerMetadata
	Stemcell *StemcellMetadata
}
