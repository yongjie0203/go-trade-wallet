package model

// Wallet 主网钱包信息
type Wallet struct {
	//钱包ID
	WalletID  int64  `json:"wallet_id" gorm:"primary_key"`
	UID       int64  `json:"uid" gorm:"column:uid"`
	MainNetId int64  `json:"main_net_id"`
	Symbol    string `json:"symbol"`
	//钱包地址
	Address string `json:"address"`
	Time    int64  `json:"time"`
	//0:未激活 1：正常 2 冻结
	Status     int    `json:"status"`
	Tags       string `json:"tags"`
	UpdateTime int64  `json:"update_time"`
	IsDelete   int    `json:"is_delete"`
	DeleteTime int64  `json:"delete_time"`
}

// WalletKey 主网钱包privateKey和publicKey
type WalletKey struct {
	KeyId int64 `json:"key_id" gorm:"primary_key"`
	//钱包ID
	WalletID   int64  `json:"wallet_id"`
	UID        int64  `json:"uid"  gorm:"column:uid"`
	MainNetId  int64  `json:"main_net_id" `
	PrivateKey string `json:"private_key" `
	PublicKey  string `json:"public_key"`
	//钱包地址
	Address string `json:"address"`
}

// WalletToken 在主网部署的Token钱包信息
type WalletToken struct {
	Id
	TokenId  int64  `json:"token_id"`
	WalletID int64  `json:"wallet_id"`
	UID      int64  `json:"uid" gorm:"column:uid"`
	Symbol   string `json:"symbol"`
	//token 余额
	Balance int64  `json:"balance"`
	Price   int64  `json:"price"  gorm:"-"`
	Address string `json:"address" gorm:"-"`
	//token 精度
	Decimals uint8 `json:"decimals"`
	Time     int64 `json:"time"`
	//0:未激活 1：激活正常 2 冻结
	Status     int    `json:"status"`
	Tags       string `json:"tags"`
	UpdateTime int64  `json:"update_time"`
	IsDelete   int    `json:"is_delete"`
	DeleteTime int64  `json:"delete_time"`
}

type MainNet struct {
	MainNetId int64  `json:"main_net_id" gorm:"primary_key"`
	Symbol    string `json:"symbol"`

	Time int64 `json:"time"`
	//0:未激活 1：正常 2 下架 3 维护
	Status     int    `json:"status"`
	Tags       string `json:"tags"`
	UpdateTime int64  `json:"update_time"`
	IsDelete   int    `json:"is_delete"`
	DeleteTime int64  `json:"delete_time"`
}

type Token struct {
	Id
	MainNetId int64 `json:"main_net_id"`
	//是否主网token如Eth为以太坊主网 1 主网 0 非主网
	IsMainNet int8 `json:"is_main_net" gorm:"default:0"`
	//1正常
	Status      int    `json:"status"`
	Price       int64  `json:"price"`
	Decimals    uint8  `json:"decimals"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	TotalSupply uint64 `json:"total_supply"`
	//合约地址
	Address string `json:"address"`
	Tail
}
