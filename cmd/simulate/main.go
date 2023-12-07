package main

import (
	"context"
	"fmt"

	"github.com/heroticket/internal/config"
	"github.com/heroticket/internal/service/ticket"
	"github.com/heroticket/internal/web3"
	"github.com/heroticket/pkg/contracts/heroticket"
)

func main() {
	cfg, err := config.NewServerConfig("./configs/server/config.json")
	if err != nil {
		panic(err)
	}

	ethclient, err := web3.NewClient(context.Background(), cfg.RpcUrl)
	if err != nil {
		panic(err)
	}

	hero, err := heroticket.NewHeroticket(web3.HexToAddress(cfg.Ticket.ContractAddress), ethclient)
	if err != nil {
		panic(err)
	}

	pvk, err := web3.ParsePrivateKey(cfg.Ticket.PrivateKey)
	if err != nil {
		panic(err)
	}

	service := ticket.New(ethclient, hero, pvk, nil, cfg.Ticket.MoralisApiKey)

	nfts, err := service.GetOwnedNFT(context.Background(), web3.HexToAddress("0x15a88243b4c61ef0071e3527b88873caf4a334dd"))
	if err != nil {
		panic(err)
	}

	for _, nft := range nfts.NFTs {
		fmt.Println(nft.MetaData)
	}
}
