package block

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/yongjie0203/go-trade-wallet/erc20"
	"github.com/yongjie0203/go-trade-wallet/model"
	"github.com/yongjie0203/go-trade-wallet/token"
	"github.com/yongjie0203/go-universal/db"
	"log"
	"math/big"
	"strings"
)

var Hub ETHHub = ETHHub{
	Headers: make(chan *types.Header),
	Blocks:  make(chan *types.Block),
	ETHTx:   make(chan *Transaction),
	ERC20Tx: make(chan *Transaction),
	Logs:    make(chan types.Log),
}

type ETHHub struct {
	Client         *ethclient.Client
	Headers        chan *types.Header
	Blocks         chan *types.Block
	ETHTx          chan *Transaction
	ERC20Tx        chan *Transaction
	Logs           chan types.Log
	SubBlockHeader ethereum.Subscription
	SubLog         ethereum.Subscription
}

type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}

type Transaction struct {
	Block *types.Block
	Tx    *types.Transaction
}

func InitClient() {
	url := "https://mainnet.infura.io/v3/4880a9fbd7de4f8a86c198aabc0fedc5"
	url = "wss://rinkeby.infura.io/ws/v3/4880a9fbd7de4f8a86c198aabc0fedc5"
	//add := "0x86699E7b0f09299C75334F81C9CC7108eC59E695"
	//url := "/home/yongjie/.ethereum/geth.ipc"
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	Hub.Client = client
}

func StartETHHub() {
	InitClient()
	var client = Hub.Client
	var err error

	Hub.SubBlockHeader, err = client.SubscribeNewHead(context.Background(), Hub.Headers)
	if err != nil {
		log.Fatal(err)
	}
	//0xdac17f958d2ee523a2206206994597c13d831ec7
	// 0x Protocol (ZRX) token address

	OnNewBlock()

}

func SubLog(block *types.Block) {

	//contractAddress := common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7")
	query := ethereum.FilterQuery{
		FromBlock: block.Number(),
		ToBlock:   nil,
		/*Addresses: []common.Address{
		contractAddress,
		},*/

	}
	var err error
	Hub.SubLog, err = Hub.Client.SubscribeFilterLogs(context.Background(), query, Hub.Logs)
	if err != nil {
		log.Fatal(err)
	}
}

func OnNewBlock() {
	for {
		select {
		case err := <-Hub.SubBlockHeader.Err():
			log.Fatal(err)
		case header := <-Hub.Headers:
			go OnHeader(header)
		case block := <-Hub.Blocks:
			go OnBlock(block)
		case log := <-Hub.Logs:
			go OnLog(log)
		case tx := <-Hub.ETHTx:
			go OnETHTx(tx)
		case tx := <-Hub.ERC20Tx:
			go OnErc20Tx(tx)
		}
	}
}

func OnHeader(header *types.Header) {
	log.Println("OnHeader")
	fmt.Println(header.Hash().Hex())
	block, err := Hub.Client.BlockByHash(context.Background(), header.Hash())
	if err != nil {
		log.Fatal(err)
	}

	go func() { Hub.Blocks <- block }()
}

func OnBlock(b *types.Block) {
	log.Println("OnBlock")
	//go SubLog(b)
	OnTxs(b.Transactions(), b)

}

func OnLog(l types.Log) {
	log.Println("OnLog")
	contractAbi, err := abi.JSON(strings.NewReader(string(token.TokenABI)))
	if err != nil {
		log.Fatal(err)
	}

	logTransferSig := []byte("Transfer(address,address,uint256)")

	logTransferSigHash := crypto.Keccak256Hash(logTransferSig).Hex()

	fmt.Printf("Log Block Number: %d\n", l.BlockNumber)
	fmt.Printf("Log Index: %d\n", l.Index)

	switch l.Topics[0].Hex() {
	case logTransferSigHash:
		fmt.Printf("Log Name: Transfer\n")
		fmt.Println(logTransferSigHash, l.Topics[0].Hex())
		var transferEvent LogTransfer
		var o []interface{}
		var err error

		o, err = contractAbi.Unpack("Transfer", l.Data)
		if err != nil {
			log.Println(err)
		}

		transferEvent.From = common.HexToAddress(l.Topics[1].Hex())
		transferEvent.To = common.HexToAddress(l.Topics[2].Hex())
		//transferEvent.Tokens = vLog.Topics[3].Hex()

		fmt.Printf("txHash:%s \n", l.TxHash)
		fmt.Printf("o:%s \n", o)
		fmt.Printf("address:%s \n", l.Address)

		fmt.Printf("From: %s\n", transferEvent.From.Hex())
		fmt.Printf("To: %s\n", transferEvent.To.Hex())
		fmt.Printf("Tokens: %s\n", transferEvent.Tokens.String())

	}

	fmt.Printf("\n\n")

}

func OnETHTx(transaction *Transaction) {
	tx := transaction.Tx
	from, _ := types.Sender(types.NewLondonSigner(tx.ChainId()), tx)
	log.Println("txHash:", tx.Hash())
	log.Println("eth转账金额：", tx.Value(), tx.Hash())
	log.Println("from:", from)
	log.Println("to:", tx.To())
	log.Println("value:", tx.Value())
	log.Println("cost:", tx.Cost())
}

func OnErc20Tx(transaction *Transaction) {
	tx := transaction.Tx
	from, _ := types.Sender(types.NewLondonSigner(tx.ChainId()), tx)
	log.Println("txHash:", tx.Hash())
	log.Println("token 转账：", tx.Hash())
	log.Println("from:", from)
	log.Println("to:", tx.To()) //合约地址
	log.Println("value:", tx.Value())
	log.Println("cost:", tx.Cost())

	txData := hexutil.Encode(tx.Data())
	//log.Println("转token到：",txData[10:74])
	//log.Println("转token数量：：",txData[74:])
	to := common.HexToAddress(txData[10:74])
	//num:="0x"+string(trimLeftZeroes(txData[74:]))
	log.Println("numstring:", "0x"+string(trimLeftZeroes(txData[74:])))
	num := hexutil.MustDecodeBig("0x" + string(trimLeftZeroes(txData[74:])))

	log.Println("转token到：", to)
	log.Println("转token数量：：", num, num.Uint64())

	token := erc20.GetERC20TokenInfo(tx.To().Hex(), Hub.Client)
	log.Println("token.Name：", token.Name)
	log.Println("token.Symbol:", token.Symbol)
	log.Println("token.Decimals：", token.Decimals)

	var tk = model.Token{}
	tk.Address = tx.To().Hex()
	tk.Decimals = token.Decimals
	tk.Name = token.Name
	tk.Symbol = token.Symbol
	tk.TotalSupply = token.TotalSupply

	tk.ID = db.NextId("")
	db.Coon.FirstOrCreate(&tk, "address = ?", tx.To().Hex())

}

func trimLeftZeroes(s string) string {
	idx := 0
	for ; idx < len(s); idx++ {
		if s[idx] != '0' {
			break
		}
	}
	return s[idx:]
}

func OnTx(tx *types.Transaction, block *types.Block) {
	//log.Println("OnTx")
	if tx.Value() != nil {
		if tx.Value().Int64() > 0 {

			t := new(Transaction)
			t.Tx = tx
			t.Block = block
			go func() { Hub.ETHTx <- t }()
			return
		}
	}
	if tx.Data() != nil {
		if len(tx.Data()) > 0 {
			//transferFnSignature := []byte("transfer(address,uint256)")
			//transfer:=crypto.Keccak256Hash(transferFnSignature).Hex()
			data := hexutil.Encode(tx.Data())
			if strings.HasPrefix(data, "0xa9059cbb") {

				t := new(Transaction)
				t.Tx = tx
				t.Block = block
				go func() { Hub.ERC20Tx <- t }()
				return
			}

		}
	}
	return
}

func OnTxs(txs types.Transactions, block *types.Block) {
	log.Println("OnTxs")
	for i := range txs {
		tx := txs[i]
		go OnTx(tx, block)
	}
}
