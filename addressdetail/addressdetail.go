package addressdetail

import "fmt"

type AddressType string

const (
	AddressTypeInit AddressType = "" // Init value

	// After detection
	AddressTypeErc20         AddressType = "Erc20"
	AddressTypeErc721        AddressType = "Erc721"
	AddressTypeErcToken      AddressType = "ErcToken"
	AddressTypeOtherContract AddressType = "OtherContract"
	AddressTypeWallet        AddressType = "Wallet" // couldn't detect a smart contract, might be just a wallet
)

type AddressDetail struct {
	Address  string      `json:"address"`
	Type     AddressType `json:"type"`
	Name     string      `json:"name"`
	Symbol   string      `json:"symbol"`
	Decimals uint8       `json:"decimals"`
}

// Returns a new unknown address detail
func NewAddressDetail(address string) AddressDetail {
	return AddressDetail{Address: address, Type: AddressTypeInit}
}

func (a AddressDetail) String() string {
	return fmt.Sprintf("%s [%s] name=%s, symbol=%s, decimals=%d", a.Address, a.Type, a.Name, a.Symbol, a.Decimals)
}

// func (a *AddressDetail) IsLoaded() bool {
func (a *AddressDetail) IsInitial() bool {
	return a.Type == AddressTypeInit
}

func (a *AddressDetail) IsErc20() bool {
	return a.Type == AddressTypeErc20
}

func (a *AddressDetail) IsErc721() bool {
	return a.Type == AddressTypeErc721
}
