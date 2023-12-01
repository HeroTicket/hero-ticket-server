package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/heroticket/internal/config"
	"github.com/heroticket/internal/web3"
	"github.com/heroticket/pkg/contracts/heroticket"
)

func main() {
	cfg, err := config.NewServerConfig("./configs/server/config.json")
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ethclient, err := web3.NewClient(ctx, cfg.RpcUrl)
	if err != nil {
		panic(err)
	}

	hero, err := heroticket.NewHeroticket(web3.HexToAddress(cfg.Ticket.ContractAddress), ethclient)
	if err != nil {
		panic(err)
	}

	// pvk, err := web3.ParsePrivateKey(cfg.Ticket.PrivateKey)
	// if err != nil {
	// 	panic(err)
	// }

	address := common.HexToAddress("0x3557db220dbfdbbb8cf5489495bf02aac9a889ed")

	tba, err := hero.TbaAddress(&bind.CallOpts{}, address)
	if err != nil {
		panic(err)
	}

	fmt.Println(tba.Hex())
	// nonce, err := ethclient.PendingNonceAt(ctx, address)
	// if err != nil {
	// 	panic(err)
	// }

	// gasPrice, err := ethclient.SuggestGasPrice(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	// chainID, err := ethclient.ChainID(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	// auth, err := bind.NewKeyedTransactorWithChainID(pvk, chainID)
	// if err != nil {
	// 	panic(err)
	// }

	// auth.GasPrice = gasPrice
	// auth.Nonce = big.NewInt(int64(nonce))
	// auth.GasLimit = 3000000

	// // 1 day = 86400

	// tx, err := hero.IssueTicket(auth, "Test Ticket", "TT", "https://ipfs.io/ipfs/QmfFbvLH37DebBqmVBm7V8ecfzgjFPnPeHRYiYk1PNoW84/2level.png",
	// 	address, big.NewInt(100), big.NewInt(1000000000), big.NewInt(100), big.NewInt(86401))
	// if err != nil {
	// 	panic(err)
	// }

	// receipt, err := bind.WaitMined(ctx, ethclient, tx)
	// if err != nil {
	// 	panic(err)
	// }

	// var ticketIssued *heroticket.HeroticketTicketIssued

	// for _, log := range receipt.Logs {
	// 	ticketIssued, err = hero.ParseTicketIssued(*log)
	// 	if err == nil {
	// 		break
	// 	}
	// }

	// if ticketIssued == nil {
	// 	panic("TicketIssued not found")
	// }

	// fmt.Println(ticketIssued.TicketAddress.Hex())
}
