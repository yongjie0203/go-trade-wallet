package erc20

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/yongjie0203/go-trade-wallet/token"
	"log"
)

type TokenERC20 struct {
	Address     common.Address
	Token       *token.Token
	Name        string
	Symbol      string
	Decimals    uint8
	TotalSupply uint64
}

func GetERC20TokenInfo(address string, client *ethclient.Client) *TokenERC20 {

	tokenAddress := common.HexToAddress(address)
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	token := new(TokenERC20)
	token.Address = tokenAddress
	token.Token = instance

	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	supply, err := instance.TotalSupply(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	token.Name = name
	token.Symbol = symbol
	token.Decimals = decimals
	token.TotalSupply = supply.Uint64()

	return token
}
