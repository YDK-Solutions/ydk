package providers

import "testing"

func TestNetconfServiceProvider_Connect(t *testing.T) {
	provider := NetconfServiceProvider{Address: "127.0.0.1", Username: "admin", Password: "admin", Port: 12022}
	provider.Connect()
}

func TestNetconfServiceProvider_Disconnect(t *testing.T) {
	provider := NetconfServiceProvider{Address: "127.0.0.1", Username: "admin", Password: "admin", Port: 12022}
	provider.Disconnect()
}

func TestNetconfServiceProvider_GetPrivate(t *testing.T) {
	provider := NetconfServiceProvider{Address: "127.0.0.1", Username: "admin", Password: "admin", Port: 12022}
	provider.GetPrivate()
}
