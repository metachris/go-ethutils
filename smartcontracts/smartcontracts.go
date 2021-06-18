// Detect smart contract properties
package smartcontracts

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/metachris/eth-go-bindings/erc165"
	"github.com/metachris/eth-go-bindings/erc20"
	"github.com/metachris/eth-go-bindings/erc721"
	"github.com/metachris/go-ethutils/addressdetail"
)

func IsContract(address string, client *ethclient.Client) (isContract bool, err error) {
	addr := common.HexToAddress(address)
	b, err := client.CodeAt(context.Background(), addr, nil)
	return len(b) > 0, err
}

func SmartContractSupportsInterface(address string, interfaceId [4]byte, client *ethclient.Client) (supportsInterface bool, err error) {
	addr := common.HexToAddress(address)
	instance, err := erc165.NewErc165(addr, client) // the SupportsInterface signature is the same for all contract types, so we can just use the ERC721 interface
	if err != nil {
		return supportsInterface, err
	}
	supportsInterface, err = instance.SupportsInterface(nil, erc165.InterfaceIdErc721)
	return supportsInterface, err
}

// TODO: Currently returns true for every SC that supports INTERFACEID_ERC165. It should really be INTERFACEID_ERC721,
// but that doesn't detect some SCs, eg. cryptokitties https://etherscan.io/address/0x06012c8cf97BEaD5deAe237070F9587f8E7A266d#readContract
// As a quick fix, just checks ERC165 and count it as ERC721 address. Improve with further/better SC method checks.
func IsErc721(address string, client *ethclient.Client) (isErc721 bool, detail addressdetail.AddressDetail, err error) {
	detail.Address = address

	addr := common.HexToAddress(address)
	instance, err := erc721.NewErc721(addr, client)
	if err != nil {
		return false, detail, err
	}

	isErc721, err = instance.SupportsInterface(nil, erc165.InterfaceIdErc165)
	if err != nil || !isErc721 {
		return isErc721, detail, err
	}

	// It appears to be ERC721
	detail.Type = addressdetail.AddressTypeErc721

	// Try to get a name and symbol
	detail.Name, _ = instance.Name(nil)
	// if err != nil {
	// 	// eg. "abi: cannot marshal in to go slice: offset 33 would go over slice boundary (len=32)"
	// 	// ignore, since we don't check erc721 metadata extension
	// }

	detail.Symbol, _ = instance.Symbol(nil)
	// if err != nil {
	// 	// ignore, since we don't check erc721 metadata extension
	// }

	return true, detail, nil
}

func IsErc20(address string, client *ethclient.Client) (isErc20 bool, detail addressdetail.AddressDetail, err error) {
	detail.Address = address
	addr := common.HexToAddress(address)
	instance, err := erc20.NewErc20(addr, client)
	if err != nil {
		return false, detail, err
	}

	detail.Name, err = instance.Name(nil)
	if err != nil || len(detail.Name) == 0 {
		return false, detail, err
	}

	// Needs symbol
	detail.Symbol, err = instance.Symbol(nil)
	if err != nil || len(detail.Symbol) == 0 {
		return false, detail, err
	}

	// Needs decimals
	detail.Decimals, err = instance.Decimals(nil)
	if err != nil {
		return false, detail, err
	}

	// Needs totalSupply
	_, err = instance.TotalSupply(nil)
	if err != nil {
		return false, detail, err
	}

	detail.Type = addressdetail.AddressTypeErc20
	return true, detail, nil
}

// GetAddressDetailFromBlockchain tries to detect an ERC20 / ERC721 token or generic smart contract, and returns an addressdetail.AddressDetail
// with the received details.
func GetAddressDetailFromBlockchain(address string, client *ethclient.Client) (detail addressdetail.AddressDetail, found bool) {
	detail = addressdetail.NewAddressDetail(address)

	// check for erc721
	if isErc721, detail, _ := IsErc721(address, client); isErc721 {
		return detail, true
	}

	// check for erc20
	if isErc20, detail, _ := IsErc20(address, client); isErc20 {
		return detail, true
	}

	// check if any type of smart contract
	if isContract, _ := IsContract(address, client); isContract {
		detail.Type = addressdetail.AddressTypeOtherContract
		return detail, true
	}

	// return just a wallet
	detail.Type = addressdetail.AddressTypeEOA
	return detail, false
}
