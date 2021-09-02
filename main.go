package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/yongjie0203/go-trade-wallet/database"
	"github.com/yongjie0203/go-trade-wallet/handler"
	"github.com/yongjie0203/go-universal/db"
	"log"
)

func main() {

	//go block.StartETHHub()

	db.InitDB()
	database.TableUpdateRegister()
	db.DBUpdate()
	//go ws.SetUpWebSocket()

	ginRouter := gin.Default()

	api := ginRouter.Group("/v1")
	wallet := api.Group("/wallet")
	{
		wallet.GET("/create", handler.Create)
		wallet.GET("/list", handler.List)
		//trade.GET("/order", handler.Order)
		//trade.POST("/list/:size/:page", controller.GetArticleList)
		//trade.GET("/get/:id", controller.GetArticle)
		//trade.GET("/delete/:id", controller.eteArticle)
	}

	ginRouter.Run()

	/*
			url:="https://mainnet.infura.io/v3/4880a9fbd7de4f8a86c198aabc0fedc5"
			url = "wss://rinkeby.infura.io/ws/v3/4880a9fbd7de4f8a86c198aabc0fedc5"
			add := "0x86699E7b0f09299C75334F81C9CC7108eC59E695"
			//url := "/home/yongjie/.ethereum/geth.ipc"
			client, err := ethclient.Dial(url)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("we have a connection")
			_ = client // we'll use this in the upcoming sections
			account := common.HexToAddress(add)
			balance, err := client.BalanceAt(context.Background(), account, nil)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(balance) // 25893180161173005034

			gasPrice, err := client.SuggestGasPrice(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("建议燃料费用："+gasPrice.String())

			headers := make(chan *types.Header)
			sub, err := client.SubscribeNewHead(context.Background(), headers)
			if err != nil {
				log.Fatal(err)
			}
			for {
				select {
				case err := <-sub.Err():
					log.Fatal(err)
				case header := <-headers:
					fmt.Println(header.Hash().Hex()) // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
					block, err := client.BlockByHash(context.Background(), header.Hash())
					if err != nil {
						log.Fatal(err)
					}

					fmt.Println(block.Hash().Hex())        // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
					fmt.Println(block.Number().Uint64())   // 3477413
					fmt.Println(block.Time())     // 1529525947
					fmt.Println(block.Nonce())             // 130524141876765836
					fmt.Println(len(block.Transactions())) // 7
					transactions := block.Transactions()
					for i := range transactions {
						tx:= transactions[i]
						from,_:=types.Sender(types.NewLondonSigner(tx.ChainId()),tx)
						fmt.Println("from:"+ from.Hex())
						//str:=hexutil.Encode(tx.Data())
						//fmt.Println("str:"+str)
						var isTransfer bool = false
						if tx.Data() != nil {
							if len(tx.Data()) > 0  {
								transferFnSignature := []byte("transfer(address,uint256)")
								transfer:=crypto.Keccak256Hash(transferFnSignature).Hex()
								data := hexutil.Encode(tx.Data())
								if strings.HasPrefix(data,transfer) {
									log.Println("transfer")
									isTransfer = true
								}


							}
							isTransfer = false
						}
						if !isTransfer {
							continue
						}

						hash:=tx.Hash()
						//data:=tx.Data()
						nonce:=tx.Nonce()
						value:=tx.Value()
						accessList:=tx.AccessList()
						chainId:=tx.ChainId()
						cost:=tx.Cost()
						gas:=tx.Gas()
						gasPrice:=tx.GasPrice()
						gasFeeCap:=tx.GasFeeCap()
						gasTipCap:=tx.GasTipCap()
						to:=tx.To()
						t:=tx.Type()
						v,r,s:=tx.RawSignatureValues()
						fmt.Printf(`
		hash:%s

		nonce:%v
		value:%v
		accessList:%v
		chainId:%v
		cost:%v
		gas:%v
		gasPrice:%v
		gasFeeCap:%v
		gasTipCap:%v
		to:%v
		type:%v
		v:%v
		r:%v
		s:%v
							`,
							hash,
						//	string(data),
							nonce,
							value,
							accessList,
							chainId,
							cost,
							gas,
							gasPrice,
							gasFeeCap,
							gasTipCap,
							to,
							t,
							v,
							r,
							s,
							)

					}
				}
			}*/

	/*
		blockNumber := big.NewInt(12928207)
		balanceAt, err := client.BalanceAt(context.Background(), account, blockNumber)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(balanceAt) // 25729324269165216042
	*/
	/*fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	fmt.Println(ethValue) // 25.729324269165216041

	pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
	fmt.Println(pendingBalance) // 25729324269165216042

	for i:=0;i<10;i++ {
		address:=test()
		account := common.HexToAddress(address)
		balance, err := client.BalanceAt(context.Background(), account, nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(balance) // 25893180161173005034

	}*/

}

func NewWallet() *Wallet {
	//创建曲线
	curve := elliptic.P256()
	//生成私钥
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	//生成公钥
	pubKeyOrig := privateKey.PublicKey

	//拼接X, Y
	pubKey := append(pubKeyOrig.X.Bytes(), pubKeyOrig.Y.Bytes()...)

	return &Wallet{Private: privateKey, PubKey: pubKey}
}

func test() string {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println("privateKey:" + hexutil.Encode(privateKeyBytes)[2:]) // 0xfad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("publicKey:" + hexutil.Encode(publicKeyBytes)[4:]) // 0x049a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("address:" + address) // 0x96216849c49358B10257cb55b28eA603c874b05E

	return address

	//return &Wallet{Private: privateKey, PubKey: publicKeyBytes}
}

type Wallet struct {
	//私钥
	Private *ecdsa.PrivateKey
	//约定，这里的PubKey不存储原始的公钥，而是存储X和Y拼接的字符串，在校验端重新拆分（参考r,s传递）
	PubKey []byte
}
