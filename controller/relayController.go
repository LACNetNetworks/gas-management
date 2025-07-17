package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	log "github.com/LACNetNetworks/gas-relay-signer/audit"
	"github.com/LACNetNetworks/gas-relay-signer/model"
	"github.com/LACNetNetworks/gas-relay-signer/rpc"
	"github.com/LACNetNetworks/gas-relay-signer/service"
)

var lock sync.Mutex

// RelayController is the main controller
type RelayController struct {
	// The controller's configuration
	Config             *model.Config
	RelaySignerService *service.RelaySignerService
}

// Init controller
func (controller *RelayController) Init(config *model.Config, relaySignerService *service.RelaySignerService) {
	controller.Config = config
	controller.RelaySignerService = relaySignerService
}

// SignTransaction ...
func (controller *RelayController) SignTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//log.GeneralLogger.Println("Body:", r.Body)

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(nil, err)
	}
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

	//	log.GeneralLogger.Println("Request body : ", rdr1)

	var rpcMessage rpc.JsonrpcMessage

	err = json.NewDecoder(rdr1).Decode(&rpcMessage)
	if err != nil {
		log.GeneralLogger.Println("Invalid params")
		log.GeneralLogger.Println(err)
		return
	}

	log.GeneralLogger.Println("JSON-RPC Method:", rpcMessage.Method)

	if rpcMessage.IsPrivTransaction() {
		r.Body = rdr2
		log.GeneralLogger.Println("Is a private Transaction, forward to Besu->Orion")

		serveReverseProxy(controller.Config.Application.NodeURL, w, r)
	} else if rpcMessage.IsPrivRawTransaction() {
		r.Body = rdr2
		log.GeneralLogger.Println("Is a private send Transaction, decrease gas used")

		controller.RelaySignerService.DecreaseGasUsed(rpcMessage.ID)

		log.GeneralLogger.Println("forward to Besu->Orion")
		serveReverseProxy(controller.Config.Application.NodeURL, w, r)
	} else if rpcMessage.IsRawTransaction() {
		processRawTransaction(controller.RelaySignerService, rpcMessage, w)
		return
	} else if rpcMessage.IsGetTransactionReceipt() {
		processGetTransactionReceipt(controller.RelaySignerService, rpcMessage, w)
		return
	} else if rpcMessage.IsGetTransactionCount() {
		processTransactionCount(controller.RelaySignerService, rpcMessage, w)
		return
	} else {
		//	r.Body=rdr2
		err := errors.New("method is not supported")
		data := handleError(rpcMessage.ID, err)
		w.Write(data)
		return
		//	serveReverseProxy(controller.Config.Application.NodeURL,w,r)
	}
}

func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, err := url.Parse(target)
	if err != nil {
		handleError(nil, err)
	}
	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

func handleError(messageID json.RawMessage, err error) []byte {
	//	log.GeneralLogger.Println(err)
	data, err := json.Marshal(service.HandleError(messageID, err))
	if err != nil {
		log.GeneralLogger.Println("Error trying to marshall a response to client")
	}

	return data
}

// GetAllPendingNonceHandler atiende GET /getAllPendingNonce?address=…
func (c *RelayController) GetAllPendingNonceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	addr := r.URL.Query().Get("address")
	if addr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "falta parámetro address"})
		return
	}

	list := c.RelaySignerService.PendingNonces(addr)
	// 200 OK con el array de int
	json.NewEncoder(w).Encode(list)
}
