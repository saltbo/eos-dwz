package rest

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"eos-dwz/src/pkg/decimal"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/token"
	"github.com/gin-gonic/gin"
)

const SYS_HOST = "http://localhost:8128"
const DWZ_DECIMAL = 62
const TRANSFER_FROM = "igetgetchain"
const TRANSFER_TO = ""

var eosCli = eos.New("https://openapi.eos.ren")
var cache = NewCache()

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

	// URL写入区块
	URLBase64 := base64.URLEncoding.EncodeToString([]byte(p.URL))
	action := token.NewTransfer(TRANSFER_FROM, TRANSFER_TO, eos.NewEOSAsset(1), URLBase64)

	signer := eos.NewKeyBag()
	signer.ImportPrivateKey("")
	eosCli.SetSigner(signer)
	ret, err := eosCli.SignPushActions(action)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	path := decimal.Decimal2Any(int(ret.BlockNum), DWZ_DECIMAL)
	ctx.JSON(http.StatusOK, strMap{"dwz": fmt.Sprintf("%s/%s", SYS_HOST, path)})
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
				if transfer.From == TRANSFER_FROM {
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
