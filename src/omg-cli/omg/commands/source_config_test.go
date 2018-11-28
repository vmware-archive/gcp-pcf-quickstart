package commands

import "testing"

func TestNameToEnv(t *testing.T) {

	cases := []struct {
		name    string
		envName string
	} {
		{
			name:    "DNSZoneName",
			envName: "DNS_ZONE_NAME",
		},
		{
			name:    "HappyCloud",
			envName: "HAPPY_CLOUD",
		},
	}

	for _, tc := range cases {
		actual := nameToEnv(tc.name)
		if tc.envName != actual {
			t.Fatalf("unexpected env name. wanted %s but got %s", tc.envName, actual)
		}
	}
}