package helpers

func RemoveMinterAddressPrefix(address string) string {
	return address[2:42]
}
