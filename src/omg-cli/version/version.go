package version

//go:generate go run gen.go
func Semver() string {
	return semver
}
