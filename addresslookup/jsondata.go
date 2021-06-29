// Tools for dealing with Ethereum addresses: AddressDetail struct, read & write token JSON, get from DB
package addresslookup

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/metachris/go-ethutils/addressdetail"
)

var URL_JSON_ADDRESSES string = "https://metachris.github.io/go-ethutils/addresslookup/json/addresses.json"
var FN_JSON_ADDRESSES string = "addresslookup/json/addresses.json"

func GetAddressesFromJsonUrl(url string) (details []addressdetail.AddressDetail, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return details, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&details)
	if err != nil {
		return details, err
	}
	return details, nil
}

func GetAddressesFromJsonFile(filename string) (details []addressdetail.AddressDetail, err error) {
	fn, _ := filepath.Abs(filename)
	file, err := os.Open(fn)
	if err != nil {
		return details, err
	}

	defer file.Close()

	// Load JSON
	decoder := json.NewDecoder(file)
	var addressDetails []addressdetail.AddressDetail
	err = decoder.Decode(&addressDetails)
	if err != nil {
		return details, err
	}

	// type field is not mandatory In JSON. Use wallet as default.
	for i, v := range addressDetails {
		addressDetails[i].Address = strings.ToLower(addressDetails[i].Address)
		if v.Type == "" {
			addressDetails[i].Type = addressdetail.AddressTypeEOA
		}
	}

	return addressDetails, nil
}

func GetAddressDetailMap(filename *string) (ret map[string]addressdetail.AddressDetail, err error) {
	if filename == nil {
		filename = &FN_JSON_ADDRESSES
	}

	list, err := GetAddressesFromJsonFile(*filename)
	if err != nil {
		return ret, err
	}

	// Convert to map
	AddressDetailMap := make(map[string]addressdetail.AddressDetail)
	for _, v := range list {
		AddressDetailMap[strings.ToLower(v.Address)] = v
	}

	return AddressDetailMap, nil
}
