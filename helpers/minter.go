package helpers

func RemoveMinterWalletPrefix(address string) string {
	return address[2:42]
}
