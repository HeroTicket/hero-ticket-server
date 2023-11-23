package web3

import (
	"context"
	"crypto/ecdsa"
	"regexp"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var AddressRegex = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

func NewClient(ctx context.Context, rpcUrl string) (*ethclient.Client, error) {
	return ethclient.DialContext(ctx, rpcUrl)
}

func ParsePrivateKey(rawPrivateKey string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(rawPrivateKey)
}

func IsAddressValid(address string) bool {
	return AddressRegex.MatchString(address)
}
