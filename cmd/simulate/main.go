package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/heroticket/internal/config"
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

	contractAddress := web3.HexToAddress("0x24f457ac6c4cc4f0ac456c9696d68b809361405b")
	tbaAddress := web3.HexToAddress("0xc50b7cb1651af7d13a89494d58fcacc595a8e905")

	ok, err := hero.HasTicket(&bind.CallOpts{}, tbaAddress, contractAddress)
	if err != nil {
		panic(err)
	}

	fmt.Println(ok)
}
