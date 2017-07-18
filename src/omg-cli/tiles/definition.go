package tiles

type PivnetDefinition struct {
	Name      string
	VersionId string
	FileId    string
	Sha256    string
}

type ProductDefinition struct {
	Name    string
	Version string
}

type Definition struct {
	Pivnet  PivnetDefinition
	Product ProductDefinition
}
