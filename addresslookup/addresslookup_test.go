package addresslookup_test

import (
	"fmt"
	"testing"

	"github.com/metachris/go-ethutils/addresslookup"
)

func TestAddressLookup(t *testing.T) {
	s := addresslookup.NewAddressLookupService(nil)

	addr, found := s.GetAddressDetail("0x3ecef08d0e2dad803847e052249bb4f8bff2d5bb") // MiningPoolHub
	fmt.Println(found, addr)
	if found {
		t.Error("first address shouldn't be found, but was", addr)
	}

	// Add all addresses from web
	err := s.AddAllAddresses()
	if err != nil {
		t.Error("couldn't add all addresses", err)
		return
	}

	_addr := "0x3ecef08d0e2dad803847e052249bb4f8bff2d5bb" // MiningPoolHub, addresses.json
	addr, found = s.GetAddressDetail(_addr)
	if !found {
		t.Error("1. address should have been found", _addr)
		return
	}
	if addr.Name != "MiningPoolHub" {
		t.Error("1. name not MiningPoolHub", addr)
		return
	}

	_addr = "0x0d0707963952f2fba59dd06f2b425ace40b492fe" // Gate.io, ethplorer-exchanges.json
	addr, found = s.GetAddressDetail(_addr)
	if !found {
		t.Error("2. address should have been found", _addr)
		return
	}
	if addr.Name != "Gate.io" {
		t.Error("2. name not Gate.io", addr)
		return
	}

	_addr = "0x00192fb10df37c9fb26829eb2cc623cd1bf599e8" // 2Miners: PPLNS, topminers-etherscan.json
	addr, found = s.GetAddressDetail(_addr)
	if !found {
		t.Error("2. address should have been found", _addr)
		return
	}
	if addr.Name != "2Miners: PPLNS" {
		t.Error("2. name not 2Miners: PPLNS", addr)
		return
	}
}
