package addresslookup

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/metachris/go-ethutils/addressdetail"
	"github.com/metachris/go-ethutils/smartcontracts"
	"github.com/metachris/go-ethutils/utils"
)

type AddressLookupService struct {
	Client *ethclient.Client

	// Initialize address cache with data from JSON
	Cache map[string]addressdetail.AddressDetail
}

func NewAddressLookupService(client *ethclient.Client) *AddressLookupService {
	return &AddressLookupService{
		Client: client,
		Cache:  make(map[string]addressdetail.AddressDetail),
	}
}

func (ads *AddressLookupService) EnsureIsLoaded(a *addressdetail.AddressDetail) {
	if !a.IsInitial() {
		return
	}

	b, _ := ads.GetAddressDetail(a.Address)
	a.Address = b.Address
	a.Type = b.Type
	a.Name = b.Name
	a.Symbol = b.Symbol
	a.Decimals = b.Decimals
}

// GetAddressDetail returns the addressdetail.AddressDetail from JSON. If not exists then query the Blockchain and caches it for future use
func (ads *AddressLookupService) GetAddressDetail(address string) (detail addressdetail.AddressDetail, found bool) {
	// Check in Cache + JSON dataset
	if detail, found := ads.Cache[strings.ToLower(address)]; found {
		return detail, true
	}

	// Without connection, return Detail with just address
	detail = addressdetail.NewAddressDetail(address) // default
	if ads.Client == nil {
		return detail, false
	}

	// Look up in Blockchain and cache (no matter if found or not, to avoid unnecessary repeated calls)
	detail, found = smartcontracts.GetAddressDetailFromBlockchain(address, ads.Client)
	ads.AddAddressDetailToCache(detail)
	return detail, found
}

func (ads *AddressLookupService) GetAddressDetailFromBlockchain(address string) (detail addressdetail.AddressDetail, found bool) {
	return smartcontracts.GetAddressDetailFromBlockchain(address, ads.Client)
}

func (ads *AddressLookupService) AddAddressDetailToCache(detail addressdetail.AddressDetail) {
	ads.Cache[strings.ToLower(detail.Address)] = detail
}

func (ads *AddressLookupService) AddAddressDetailsToCache(details []addressdetail.AddressDetail) {
	for _, detail := range details {
		ads.Cache[strings.ToLower(detail.Address)] = detail
	}
}

func (ads *AddressLookupService) ClearCache() {
	ads.Cache = make(map[string]addressdetail.AddressDetail)
}

func (ads *AddressLookupService) AddAddressesFromJsonUrl(url string) error {
	details, err := GetAddressesFromJsonUrl(url)
	if err != nil {
		return err
	}

	if utils.DebugEnabled {
		fmt.Printf("adding %d entries from %s\n", len(details), url)
	}
	ads.AddAddressDetailsToCache(details)
	return nil
}

func (ads *AddressLookupService) AddAllAddresses() error {
	jsonUrls := []string{
		JsonUrlAddresses,
		JsonUrlEtherscanTopminers,
		JsonUrlEthplorerExchangeAddresses,
	}

	for _, url := range jsonUrls {
		err := ads.AddAddressesFromJsonUrl(url)
		if err != nil {
			return errors.New(err.Error() + " - " + url)
		}
	}

	return nil
}
