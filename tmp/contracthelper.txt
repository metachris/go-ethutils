// package contracthelper

// import (
// 	"context"

// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/ethclient"
// 	"github.com/metachris/eth-go-bindings/erc165"
// 	"github.com/metachris/eth-go-bindings/erc20"
// 	"github.com/metachris/eth-go-bindings/erc721"
// )

// func IsContract(address string, client *ethclient.Client) bool {
// 	addr := common.HexToAddress(address)
// 	b, err := client.CodeAt(context.Background(), addr, nil)
// 	Perror(err)
// 	return len(b) > 0
// }

// func SmartContractSupportsInterface(address string, interfaceId [4]byte, client *ethclient.Client) bool {
// 	addr := common.HexToAddress(address)
// 	instance, err := erc165.NewErc165(addr, client) // the SupportsInterface signature is the same for all contract types, so we can just use the ERC721 interface
// 	Perror(err)
// 	isSupported, err := instance.SupportsInterface(nil, erc165.InterfaceIdErc721)
// 	return err == nil && isSupported
// }

// // TODO: Currently returns true for every SC that supports INTERFACEID_ERC165. It should really be INTERFACEID_ERC721,
// // but that doesn't detect some SCs, eg. cryptokitties https://etherscan.io/address/0x06012c8cf97BEaD5deAe237070F9587f8E7A266d#readContract
// // As a quick fix, just checks ERC165 and count it as ERC721 address. Improve with further/better SC method checks.
// func IsErc721(address string, client *ethclient.Client) (isErc721 bool, detail AddressDetail) {
// 	detail.Address = address

// 	addr := common.HexToAddress(address)
// 	instance, err := erc721.NewErc721(addr, client)
// 	Perror(err)

// 	isErc721, err = instance.SupportsInterface(nil, erc165.InterfaceIdErc165)
// 	if err != nil || !isErc721 {
// 		return false, detail
// 	}

// 	// It appears to be ERC721
// 	detail.Type = AddressTypeErc721

// 	// Try to get a name and symbol
// 	detail.Name, _ = instance.Name(nil)
// 	// if err != nil {
// 	// 	// eg. "abi: cannot marshal in to go slice: offset 33 would go over slice boundary (len=32)"
// 	// 	// ignore, since we don't check erc721 metadata extension
// 	// }

// 	detail.Symbol, _ = instance.Symbol(nil)
// 	// if err != nil {
// 	// 	// ignore, since we don't check erc721 metadata extension
// 	// }

// 	return true, detail
// }

// func IsErc20(address string, client *ethclient.Client) (isErc20 bool, detail AddressDetail) {
// 	detail.Address = address
// 	addr := common.HexToAddress(address)
// 	instance, err := erc20.NewErc20(addr, client)
// 	Perror(err)

// 	detail.Name, err = instance.Name(nil)
// 	if err != nil || len(detail.Name) == 0 {
// 		// fmt.Println(1, err)
// 		return false, detail
// 	}

// 	// Needs symbol
// 	detail.Symbol, err = instance.Symbol(nil)
// 	if err != nil || len(detail.Symbol) == 0 {
// 		// fmt.Println(2)
// 		return false, detail
// 	}

// 	// Needs decimals
// 	detail.Decimals, err = instance.Decimals(nil)
// 	if err != nil {
// 		// fmt.Println(3, err)
// 		return false, detail
// 	}

// 	// Needs totalSupply
// 	_, err = instance.TotalSupply(nil)
// 	if err != nil {
// 		// fmt.Println(4)
// 		return false, detail
// 	}

// 	detail.Type = AddressTypeErc20
// 	return true, detail
// }

// func GetAddressDetailFromBlockchain(address string, client *ethclient.Client) (detail AddressDetail, found bool) {
// 	ret := NewAddressDetail(address)

// 	if isErc721, detail := IsErc721(address, client); isErc721 {
// 		return detail, true
// 	}

// 	if isErc20, detail := IsErc20(address, client); isErc20 {
// 		return detail, true
// 	}

// 	if IsContract(address, client) {
// 		ret.Type = AddressTypeOtherContract
// 		return ret, true
// 	}

// 	ret.Type = AddressTypePubkey
// 	return ret, false
// }
