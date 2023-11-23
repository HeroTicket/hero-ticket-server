package ipfs

import "fmt"

var ErrUnknown = fmt.Errorf("unknown error")

const PinFileToIPFSUrl = "https://api.pinata.cloud/pinning/pinFileToIPFS"

type PinFileResponse struct {
	IpfsHash    string `json:"IpfsHash"`
	PinSize     int    `json:"PinSize"`
	Timestamp   string `json:"Timestamp"`
	IsDuplicate bool   `json:"IsDuplicate"`
}
