package uaa

type Name struct {
	FamilyName string `json:"familyName"`
	GivenName  string `json:"givenName"`
}
type PhoneNumber struct {
	Value string `json:"value"`
}
type Email struct {
	Value   string `json:"value"`
	Primary bool   `json:"primary"`
}
type User struct {
	Id           string        `json:"id"`
	Username     string        `json:"userName"`
	Password     string        `json:"password"`
	Name         Name          `json:"name"`
	PhoneNumbers []PhoneNumber `json:"phoneNumbers"`
	Emails       []Email       `json:"emails"`
	Active       bool          `json:"active"`
}
