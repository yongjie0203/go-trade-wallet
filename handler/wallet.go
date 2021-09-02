package handler

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/yongjie0203/go-universal/db"
	"log"
	"sync"
	"time"

	"github.com/yongjie0203/go-trade-wallet/model"
	"github.com/yongjie0203/go-trade-wallet/request"
)

func Create(ctx *gin.Context) {
	req := new(request.CreatWallet)
	ctx.Bind(req)

}

func CreateWalletKey(mainNet model.MainNet) *model.WalletKey {
	return CreateETHWalletKey()
}

func CreateETHWalletKey() *model.WalletKey {
	walletKey := new(model.WalletKey)

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	walletKey.Address = address
	walletKey.PrivateKey = hexutil.Encode(privateKeyBytes)[2:]
	walletKey.PublicKey = hexutil.Encode(publicKeyBytes)[4:]
	return walletKey
}

func List(ctx *gin.Context) {
	var wg sync.WaitGroup
	req := new(request.BaseRequest)
	tokens := make([]model.Token, 0)
	mainNets := make([]model.MainNet, 0)
	mainNetMap := make(map[int64]model.MainNet)
	mainNetWalletMap := make(map[int64]*model.Wallet)
	//查询支持的主网信息
	db.Coon.Find(&mainNets)
	for i := range mainNets {
		mainNetMap[mainNets[i].MainNetId] = mainNets[i]
		wallet := model.Wallet{}
		var hasWallet int64 = 0
		db.Coon.Where("main_net_id = ? and uid = ? ", mainNets[i].MainNetId, req.UID).First(&wallet).Count(&hasWallet)
		if hasWallet > 0 {
			mainNetWalletMap[mainNets[i].MainNetId] = &wallet
		}
	}
	//查询支持的所有token
	db.Coon.Find(&tokens)
	walletTokens := make([]model.WalletToken, len(tokens))
	for i := range tokens {
		fmt.Println(tokens[i].Name)
		tk := tokens[i]

		var idx int = i
		mainNet := mainNetMap[tk.MainNetId]
		if tk.IsMainNet == 1 {
			wallet, ok := mainNetWalletMap[tk.MainNetId]
			if !ok {
				//创建钱包秘钥
				walletKey := CreateWalletKey(mainNet)
				walletKey.KeyId = db.NextId("walletKey")
				walletKey.WalletID = db.NextId("wallet")
				walletKey.UID = req.UID
				db.Coon.Create(walletKey)
				//创建主网钱包
				wallet.WalletID = walletKey.WalletID
				wallet.MainNetId = tk.MainNetId
				wallet.UID = req.UID
				wallet.Time = time.Now().UnixNano()
				wallet.Symbol = mainNet.Symbol
				wallet.Address = walletKey.Address
				db.Coon.Create(&wallet)
			}
		}
		wg.Add(1)
		go func() {
			walletToken := model.WalletToken{}

			var count int64 = 0

			db.Coon.Where("token_id = ? and uid = ?", tk.ID, req.UID).First(&walletToken).Count(&count)

			if count == 0 {
				fmt.Println("无记录")
				wallet := model.Wallet{}
				var hasWallet int64 = 0
				db.Coon.Where("main_net_id = ? and uid = ? ", tk.MainNetId, req.UID).First(&wallet).Count(&hasWallet)

				walletToken.ID = db.NextId("walletToken")
				walletToken.TokenId = tk.ID
				walletToken.Symbol = tk.Symbol
				walletToken.Decimals = tk.Decimals
				walletToken.Price = tk.Price
				walletToken.Time = time.Now().UnixNano()
				walletToken.UID = req.UID
				walletToken.Address = wallet.Address
				db.Coon.Create(&walletToken)
			} else {
				wallet, ok := mainNetWalletMap[tk.MainNetId]
				if ok {
					walletToken.Address = wallet.Address
				}
			}

			walletTokens[idx] = walletToken
			wg.Done()
		}()
	}
	wg.Wait()
	ctx.JSON(200, walletTokens)
}
