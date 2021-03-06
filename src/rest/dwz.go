package rest

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"eos-dwz/src/option"
	"eos-dwz/src/pkg/decimal"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/token"
	"github.com/gin-gonic/gin"
)

const DWZ_DECIMAL = 62

var cache = NewCache()
var eosCli = eos.New(option.NodeHost)

type generateRequest struct {
	URL string `json:"url"`
}

type strMap map[string]string

func generate(ctx *gin.Context) {
	p := new(generateRequest)
	if err := ctx.BindJSON(p); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if option.SendAccount == "" || option.PrivateKey == "" {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("no settings for generate api."))
		return
	}

	signer := eos.NewKeyBag()
	err := signer.ImportPrivateKey(option.PrivateKey)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	eosCli.SetSigner(signer)

	// URL写入区块
	// eosCli.Debug = true
	action := token.NewTransfer(option.SendAccount, option.ReceiveAccount, eos.NewEOSAsset(1), p.URL)
	ret, err := eosCli.SignPushActions(action)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	path := decimal.Decimal2Any(int(ret.Processed.BlockNum+1), DWZ_DECIMAL)
	ctx.JSON(http.StatusOK, strMap{"dwz": fmt.Sprintf("%s/%s", option.ServerHost, path)})
}

func dwzHandler(ctx *gin.Context) {
	u := ctx.Request.URL
	if u.Path == "/" {
		ctx.String(http.StatusOK, "%s", "welcome visit this site.")
		return
	}

	// 查询本地缓存，如果有则直接返回
	alias := strings.Split(u.Path, "/")[1]
	if url, ok := cache.Get(alias); ok {
		ctx.Redirect(http.StatusMovedPermanently, url)
		return
	}

	// 从链上查询区块信息
	blockNum := decimal.Any2Decimal(alias, DWZ_DECIMAL)
	ret, err := eosCli.GetBlockByNum(uint32(blockNum))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	rawURL := urlFromBlock(ret.Transactions)
	cache.Set(alias, rawURL)
	ctx.Redirect(http.StatusMovedPermanently, rawURL)
}

func urlFromBlock(transactions []eos.TransactionReceipt) string {
	for _, transaction := range transactions {
		if transaction.Transaction.Packed == nil {
			continue
		}

		a, _ := transaction.Transaction.Packed.Unpack()
		for _, action := range a.Actions {
			if action.Name == "transfer" && action.Account == "eosio.token" {
				transfer := action.Data.(*token.Transfer)
				if transfer.To == option.ReceiveAccount {
					return transfer.Memo
				}
			}
		}
	}

	return ""
}

type Cache struct {
	sync.RWMutex
	store map[string]string
}

func NewCache() *Cache {
	return &Cache{store: make(map[string]string)}
}

func (c *Cache) Set(key, val string) {
	c.Lock()
	defer c.Unlock()

	c.store[key] = val
}

func (c *Cache) Get(key string) (string, bool) {
	c.RLock()
	defer c.RUnlock()

	if val, ok := c.store[key]; ok {
		return val, true
	}

	return "", false
}
