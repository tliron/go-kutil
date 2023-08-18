package util

import (
	"testing"
)

func TestLetterCase(t *testing.T) {
	s := "IpAddress"
	assertTransform(t, ToDromedaryCase, s, "ipAddress")
	assertTransform(t, ToSnakeCase, s, "ip_address")
	assertTransform(t, ToKebabCase, s, "ip-address")

	s = "IPAddress"
	assertTransform(t, ToDromedaryCase, s, "ipaddress") // hmm...
	assertTransform(t, ToSnakeCase, s, "ip_address")
	assertTransform(t, ToKebabCase, s, "ip-address")

	s = "IP_Address"
	assertTransform(t, ToDromedaryCase, s, "ipAddress")
	assertTransform(t, ToSnakeCase, s, "ip_address")
	assertTransform(t, ToKebabCase, s, "ip-address")

	s = "IP-Address"
	assertTransform(t, ToDromedaryCase, s, "ipAddress")
	assertTransform(t, ToSnakeCase, s, "ip_address")
	assertTransform(t, ToKebabCase, s, "ip-address")

	s = "IP Address"
	assertTransform(t, ToDromedaryCase, s, "ipAddress")
	assertTransform(t, ToSnakeCase, s, "ip_address")
	assertTransform(t, ToKebabCase, s, "ip-address")
}
