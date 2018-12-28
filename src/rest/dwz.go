package rest

import (
	"fmt"
	"net/http"
	"strings"

	"eos-dwz/src/pkg/decimal"
	"github.com/gin-gonic/gin"
)

const SYS_HOST = "http://localhost:8128"
const DWZ_DECIMAL = 62

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

	blockId := 34428259
	path := decimal.Decimal2Any(blockId, DWZ_DECIMAL)
	ctx.JSON(http.StatusOK, strMap{"dwz": fmt.Sprintf("%s/%s", SYS_HOST, path)})
}

func dwzHandler(ctx *gin.Context) {
	u := ctx.Request.URL
	alias := strings.Split(u.Path, "/")[1]
	blockId := decimal.Any2Decimal(alias, DWZ_DECIMAL)
	fmt.Println(u, alias, blockId)

	//根据blockId获取区块数据

	//rawURL := ""
	//ctx.Redirect(http.StatusMovedPermanently, rawURL)
}
