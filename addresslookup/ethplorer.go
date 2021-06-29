package addresslookup

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type EthplorerServiceResponse struct {
	PublicTags []string `json:"publicTags"`
	IsContract bool     `json:"isContract"`
}

func EthplorerServiceAddressLookup(addr string) (res EthplorerServiceResponse, err error) {
	url := fmt.Sprintf("https://ethplorer.io/service/service.php?data=%s&showTx=none", addr)
	resp, err := http.Get(url)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}
