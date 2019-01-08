package option

import (
	"flag"
	"log"

	"github.com/eoscanada/eos-go"
)

var (
	// cluster stuff
	NodeHost       string
	ServerHost     string
	ReceiveAccount eos.AccountName
	SendAccount    eos.AccountName
	PrivateKey     string
)

func init() {
	nodeHost := flag.String("node_host", "", "specify eos node host")
	serverHost := flag.String("server_host", "", "specify eos node host")
	receiveAccount := flag.String("receive_account", "duanwangzhix", "specify receive eos account name")
	sendAccount := flag.String("send_account", "", "specify send eos account name")
	privateKey := flag.String("private_key", "", "specify private key of the send eos account ")
	flag.Parse()

	NodeHost = *nodeHost
	if NodeHost == "" {
		log.Fatalf("empty node_host")
	}

	ServerHost = *serverHost
	if ServerHost == "" {
		log.Fatalf("empty server_host")
	}

	ReceiveAccount = eos.AccountName(*receiveAccount)
	SendAccount = eos.AccountName(*sendAccount)
	PrivateKey = *privateKey
}
