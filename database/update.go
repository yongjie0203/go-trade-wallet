package database

import (
	"github.com/yongjie0203/go-trade-wallet/model"
	"github.com/yongjie0203/go-universal/db"
)

func TableUpdateRegister() {
	db.TableUpdateRegister("wallet", &model.Wallet{})
	db.TableUpdateRegister("walletToken", &model.WalletToken{})
	db.TableUpdateRegister("walletKey", &model.WalletKey{})
	db.TableUpdateRegister("mainNet", &model.MainNet{})
	db.TableUpdateRegister("token", &model.Token{})
}
