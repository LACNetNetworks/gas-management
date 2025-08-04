package service

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/LACNetNetworks/gas-relay-signer/model"
	"github.com/LACNetNetworks/gas-relay-signer/rpc"
	"github.com/ethereum/go-ethereum/common"
)

var sequence uint8 = 0

func TestInit(t *testing.T) {
	srv := serverMock()
	defer srv.Close()

	dir, _ := os.Getwd()

	createKeyMock(dir + "/keyMock")
	setKeyMock()

	applicationConfig := model.ApplicationConfig{NodeURL: srv.URL + "/getRelayHubContract", NodeKeyPath: dir + "/keyMock"}
	applicationConfig.ContractAddress = "0xdD37c69fF29C4b93A346Ed6dF184f48A71800b7E"
	config := model.Config{Application: applicationConfig}
	relaySignerService := new(RelaySignerService)

	relaySignerService.Init(&config)

	err := os.RemoveAll(dir + "/log")
	if err != nil {
		log.Fatal(err)
	}

	err = os.Remove("keyMock")
	if err != nil {
		log.Fatal(err)
	}

	if relaySignerService.Config.Application.Key != "b3e7374dca5ca90c3899dbb2c978051437fb15534c945bf59df16d6c80be27c0" {
		t.Errorf("Private Key wasn't loaded from file")
	}

	if relaySignerService.Config.Application.RelayHubContractAddress.Hex() != "0xfF6D55d01FB12695EA00c071aD8aF3CE44cF3A91" {
		t.Errorf("RelayHub Smart Contract can't be loaded from Proxy")
	}
}

func TestFailInit(t *testing.T) {
	applicationConfig := model.ApplicationConfig{NodeKeyPath: "./keyMock"}
	config := model.Config{Application: applicationConfig}
	relaySignerService := new(RelaySignerService)

	os.Unsetenv("WRITER_KEY")

	err := relaySignerService.Init(&config)

	if err.Error() != "Environment variable WRITER_KEY not set" {
		t.Errorf("A environment variable shouldn't be set")
	}
}

func TestGetTransactionCount(t *testing.T) {
	srv := serverMock()
	defer srv.Close()

	contents := []byte(`{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["0x92c9885663f6e84127c857d3137936c424b7e07555d2bc7d8bd781b3f0847ac8"],"id":53}`)

	var rpcMessage rpc.JsonrpcMessage
	_ = json.Unmarshal(contents, &rpcMessage)

	var params []string
	_ = json.Unmarshal(rpcMessage.Params, &params)

	applicationConfig := model.ApplicationConfig{NodeURL: srv.URL + "/getTransactionCount", Key: "b3e7374dca5ca90c3899dbb2c978051437fb15534c945bf59df16d6c80be27c0"}
	config := model.Config{Application: applicationConfig}

	relaySignerService := new(RelaySignerService)
	_ = relaySignerService.Init(&config)
	relayHubAddress := common.HexToAddress("0xdD37c69fF29C4b93A346Ed6dF184f48A71800b7E")
	relaySignerService.Config.Application.RelayHubContractAddress = &relayHubAddress
	relaySignerService.Config.Application.ContractAddress = "0xdD37c69fF29C4b93A346Ed6dF184f48A71800b7E"
	jsonResponse := relaySignerService.GetTransactionCount(rpcMessage.ID, params[0], false)

	fmt.Println(jsonResponse)

	if jsonResponse.String() != `{"jsonrpc":"2.0","id":53,"result":"0x159"}` {
		t.Errorf("Incorrect nonce was gotten")
	}
}

func TestGetTransactionCountLatest(t *testing.T) {
	srv := serverMock()
	defer srv.Close()

	contents := []byte(`{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["0x92c9885663f6e84127c857d3137936c424b7e07555d2bc7d8bd781b3f0847ac8","latest"],"id":53}`)

	var rpcMessage rpc.JsonrpcMessage
	_ = json.Unmarshal(contents, &rpcMessage)

	var params []string
	_ = json.Unmarshal(rpcMessage.Params, &params)

	applicationConfig := model.ApplicationConfig{NodeURL: srv.URL + "/getTransactionCount", Key: "b3e7374dca5ca90c3899dbb2c978051437fb15534c945bf59df16d6c80be27c0"}
	config := model.Config{Application: applicationConfig}

	relaySignerService := new(RelaySignerService)
	_ = relaySignerService.Init(&config)
	relayHubAddress := common.HexToAddress("0xdD37c69fF29C4b93A346Ed6dF184f48A71800b7E")
	relaySignerService.Config.Application.RelayHubContractAddress = &relayHubAddress
	relaySignerService.Config.Application.ContractAddress = "0xdD37c69fF29C4b93A346Ed6dF184f48A71800b7E"
	jsonResponse := relaySignerService.GetTransactionCount(rpcMessage.ID, params[0], false)

	fmt.Println(jsonResponse)

	if jsonResponse.String() != `{"jsonrpc":"2.0","id":53,"result":"0x159"}` {
		t.Errorf("Incorrect nonce was gotten")
	}
}

func TestGetTransactionCountPending(t *testing.T) {
	srv := serverMock()
	defer srv.Close()

	contents := []byte(`{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["0x92c9885663f6e84127c857d3137936c424b7e07555d2bc7d8bd781b3f0847ac8","pending"],"id":53}`)

	var rpcMessage rpc.JsonrpcMessage
	_ = json.Unmarshal(contents, &rpcMessage)

	var params []string
	_ = json.Unmarshal(rpcMessage.Params, &params)

	applicationConfig := model.ApplicationConfig{NodeURL: srv.URL + "/getTransactionCount", ContractAddress: "0xdD37c69fF29C4b93A346Ed6dF184f48A71800b7E", Key: "b3e7374dca5ca90c3899dbb2c978051437fb15534c945bf59df16d6c80be27c0"}
	config := model.Config{Application: applicationConfig}
	relaySignerService := new(RelaySignerService)
	//_ = relaySignerService.Init(&config)
	if err := relaySignerService.Init(&config); err != nil {
		t.Fatal(err)
	}
	// pon el RelayHub
	hub := common.HexToAddress("0xdD37c69fF29C4b93A346Ed6dF184f48A71800b7E")
	relaySignerService.Config.Application.RelayHubContractAddress = &hub

	// 3) Inicializa la cola con la dirección *checksummed*
	addr := common.HexToAddress(params[0]).Hex()

	relaySignerService.senders = make(map[string]*BigIntQueue)
	relaySignerService.senders[addr] = &BigIntQueue{new(big.Int).SetUint64(200)}
	jsonResponse := relaySignerService.GetTransactionCount(rpcMessage.ID, params[0], true)

	//jsonResponse := relaySignerService.GetTransactionCount(rpcMessage.ID, params[0], true)
	want := `{"jsonrpc":"2.0","id":53,"result":"0xc8"}`
	if jsonResponse.String() != want {
		t.Errorf("Incorrect nonce: got %s, want %s", jsonResponse.String(), want)
	}
}

func TestGetTransactionCountPendingNoValue(t *testing.T) {
	srv := serverMock()
	defer srv.Close()

	contents := []byte(`{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["0x92c9885663f6e84127c857d3137936c424b7e07555d2bc7d8bd781b3f0847ac8","pending"],"id":53}`)

	var rpcMessage rpc.JsonrpcMessage
	_ = json.Unmarshal(contents, &rpcMessage)

	var params []string
	_ = json.Unmarshal(rpcMessage.Params, &params)

	applicationConfig := model.ApplicationConfig{NodeURL: srv.URL + "/getTransactionCount", ContractAddress: "0xdD37c69fF29C4b93A346Ed6dF184f48A71800b7E", Key: "b3e7374dca5ca90c3899dbb2c978051437fb15534c945bf59df16d6c80be27c0"}
	config := model.Config{Application: applicationConfig}
	relaySignerService := new(RelaySignerService)
	_ = relaySignerService.Init(&config)
	relayHubAddress := common.HexToAddress("0xdD37c69fF29C4b93A346Ed6dF184f48A71800b7E")
	relaySignerService.Config.Application.RelayHubContractAddress = &relayHubAddress
	relaySignerService.Config.Application.ContractAddress = "0xdD37c69fF29C4b93A346Ed6dF184f48A71800b7E"
	relaySignerService.senders = make(map[string]*BigIntQueue)
	jsonResponse := relaySignerService.GetTransactionCount(rpcMessage.ID, params[0], true)

	if jsonResponse.String() != `{"jsonrpc":"2.0","id":53,"result":"0x159"}` {
		t.Errorf("Incorrect nonce was gotten")
	}
}

func TestGetTransactionReceipt(t *testing.T) {
	srv := serverMock()
	defer srv.Close()

	contents := []byte(`{"jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["0x504ce587a65bdbdb6414a0c6c16d86a04dd79bfcc4f2950eec9634b30ce5370f"],"id":53}`)

	var rpcMessage rpc.JsonrpcMessage
	_ = json.Unmarshal(contents, &rpcMessage)

	var params []string
	_ = json.Unmarshal(rpcMessage.Params, &params)

	applicationConfig := model.ApplicationConfig{NodeURL: srv.URL + "/getReceipt"}
	config := model.Config{Application: applicationConfig}
	relaySignerService := new(RelaySignerService)
	_ = relaySignerService.Init(&config)
	jsonResponse := relaySignerService.GetTransactionReceipt(rpcMessage.ID, params[0])

	var result map[string]interface{}
	json.Unmarshal([]byte(jsonResponse.String()), &result)

	blockHash := result["result"].(map[string]interface{})

	if blockHash["blockHash"] != "0x6e3aa24e261e61832624749b64049104c6105ba870d3375484548ffdb133eeea" {
		t.Errorf("Incorrect blockHash was gotten")
	}
}

func TestGetTransactionReceiptRevertReason(t *testing.T) {
	srv := serverMock()
	defer srv.Close()

	contents := []byte(`{"jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["0x504ce587a65bdbdb6414a0c6c16d86a04dd79bfcc4f2950eec9634b30ce5370f"],"id":53}`)

	var rpcMessage rpc.JsonrpcMessage
	_ = json.Unmarshal(contents, &rpcMessage)

	var params []string
	_ = json.Unmarshal(rpcMessage.Params, &params)

	applicationConfig := model.ApplicationConfig{NodeURL: srv.URL + "/getReceiptRevertReason"}
	config := model.Config{Application: applicationConfig}
	relaySignerService := new(RelaySignerService)
	_ = relaySignerService.Init(&config)
	jsonResponse := relaySignerService.GetTransactionReceipt(rpcMessage.ID, params[0])

	var result map[string]interface{}
	json.Unmarshal([]byte(jsonResponse.String()), &result)

	blockHash := result["result"].(map[string]interface{})

	if blockHash["revertReason"] != "0x08c379a0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000096e616275636f646f730000000000000000000000000000000000000000000000" {
		t.Errorf("Incorrect revert reason was gotten")
	}
}

func TestSendMetatransaction(t *testing.T) {
	srv := serverMock()
	defer srv.Close()

	contents := []byte(`{"id":2914410858336929,"jsonrpc":"2.0","params":["0xf8840180831e8480946e6bbf31aa45042d53128339383fcd1c377b42c780a46057361d00000000000000000000000000000000000000000000000000000000000001591ba028934b543809922b277e85f6bcf7b1f25e937de05c5138e17fdfa480ba74e84ba055a2a611763ffcb748547408093551928c9549f95a0a9cabd3b1f1f2e166cc16"],"method":"eth_sendRawTransaction"}`)

	var rpcMessage rpc.JsonrpcMessage
	_ = json.Unmarshal(contents, &rpcMessage)

	var params []string
	_ = json.Unmarshal(rpcMessage.Params, &params)

	dir, _ := os.Getwd()

	createKeyMock(dir + "/keyMock")
	setKeyMock()

	applicationConfig := model.ApplicationConfig{NodeURL: srv.URL + "/getRelayHubContract", ContractAddress: "0x0ae2Da68515Ef8DC4bBCa1fA1bcE00C508b2Af4B", NodeKeyPath: dir + "/keyMock"}
	config := model.Config{Application: applicationConfig}
	relaySignerService := new(RelaySignerService)
	_ = relaySignerService.Init(&config)

	relaySignerService.Config.Application.NodeURL = srv.URL + "/sendMetatransaction"

	to := common.HexToAddress("0x82a978b3f5962a5b0957d9ee9eef472ee55b42f1")

	encodedFunction, _ := hex.DecodeString("0xf861808082ea6094fd32cfc2e71611626d6368a41f915d0077a306a180b8446057361d000000000000000000000000000000000000000000000000000000000000003c000000000000000000000000173cf75f0905338597fcd38f5ce13e6840b230e9")

	var gasLimit uint64
	var r [32]byte
	var s [32]byte

	gasLimit = 200000

	sender := "0x92c9885663f6e84127c857d3137936c424b7e07555d2bc7d8bd781b3f0847ac8"

	jsonResponse := relaySignerService.SendMetatransaction(rpcMessage.ID, &to, gasLimit, encodedFunction, 27, r, s, sender, 34)

	err := os.Remove("keyMock")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(jsonResponse.String())

	if jsonResponse.String() != `{"jsonrpc":"2.0","id":2914410858336929,"result":"0x9c2fb4956ce18491021a534106fe50e7cfe86bcc373b1626623fa0366f4cc3bc"}` {
		t.Errorf("Incorrect transactionHash was gotten")
	}
}

func TestNonceAfterTransactions(t *testing.T) {
	srv := serverMock()
	defer srv.Close()

	contents := []byte(`{"id":2914410858336929,"jsonrpc":"2.0","params":["0xf8840180831e8480946e6bbf31aa45042d53128339383fcd1c377b42c780a46057361d00000000000000000000000000000000000000000000000000000000000001591ba028934b543809922b277e85f6bcf7b1f25e937de05c5138e17fdfa480ba74e84ba055a2a611763ffcb748547408093551928c9549f95a0a9cabd3b1f1f2e166cc16"],"method":"eth_sendRawTransaction"}`)

	var rpcMessage rpc.JsonrpcMessage
	_ = json.Unmarshal(contents, &rpcMessage)

	var params []string
	_ = json.Unmarshal(rpcMessage.Params, &params)

	dir, _ := os.Getwd()

	createKeyMock(dir + "/keyMock")
	setKeyMock()

	applicationConfig := model.ApplicationConfig{NodeURL: srv.URL + "/getRelayHubContract", ContractAddress: "0x0ae2Da68515Ef8DC4bBCa1fA1bcE00C508b2Af4B", NodeKeyPath: dir + "/keyMock"}
	config := model.Config{Application: applicationConfig}
	relaySignerService := new(RelaySignerService)
	_ = relaySignerService.Init(&config)

	relaySignerService.Config.Application.NodeURL = srv.URL + "/sendMetatransaction"

	to := common.HexToAddress("0x82a978b3f5962a5b0957d9ee9eef472ee55b42f1")

	encodedFunction, _ := hex.DecodeString("0xf861808082ea6094fd32cfc2e71611626d6368a41f915d0077a306a180b8446057361d000000000000000000000000000000000000000000000000000000000000003c000000000000000000000000173cf75f0905338597fcd38f5ce13e6840b230e9")

	var gasLimit uint64
	var r [32]byte
	var s [32]byte

	gasLimit = 200000

	sender := "0x92c9885663f6e84127c857d3137936c424b7e07555d2bc7d8bd781b3f0847ac8"

	for i := 34; i < 45; i++ {
		jsonResponse := relaySignerService.SendMetatransaction(rpcMessage.ID, &to, gasLimit, encodedFunction, 27, r, s, sender, uint64(i))
		if jsonResponse.String() != `{"jsonrpc":"2.0","id":2914410858336929,"result":"0x9c2fb4956ce18491021a534106fe50e7cfe86bcc373b1626623fa0366f4cc3bc"}` {
			t.Errorf("Incorrect transactionHash was gotten")
		}

		jsonResponseNonce := relaySignerService.GetTransactionCount(rpcMessage.ID, sender, true)

		fmt.Println(jsonResponseNonce)

		nonce := fmt.Sprintf("%x", i)

		fmt.Println(nonce)

		responseNonce := fmt.Sprintf(`{"jsonrpc":"2.0","id":2914410858336929,"result":"0x%s"}`, nonce)

		if jsonResponseNonce.String() != responseNonce {
			t.Errorf("Incorrect nonce was gotten")
		}
	}

	err := os.Remove("keyMock")
	if err != nil {
		log.Fatal(err)
	}
}

func TestNonceCollisionSendMetatransactionSameSender(t *testing.T) {
    
	const workers = 3
    
    // 1) Stand up your same serverMock
    srv := serverMock()
    defer srv.Close()

    // 2) Prepare the service exactly as in TestSendMetatransaction
    dir, _ := os.Getwd()
    createKeyMock(dir + "/keyMock")
    setKeyMock()

    cfg := &model.Config{
        Application: model.ApplicationConfig{
            NodeURL:         srv.URL + "/getRelayHubContract",
            ContractAddress: "0x0ae2Da68515Ef8DC4bBCa1fA1bcE00C508b2Af4B",
            NodeKeyPath:     dir + "/keyMock",
        },
    }

	contents := []byte(`{"id":2914410858336929,"jsonrpc":"2.0","params":["0xf8840180831e8480946e6bbf31aa45042d53128339383fcd1c377b42c780a46057361d00000000000000000000000000000000000000000000000000000000000001591ba028934b543809922b277e85f6bcf7b1f25e937de05c5138e17fdfa480ba74e84ba055a2a611763ffcb748547408093551928c9549f95a0a9cabd3b1f1f2e166cc16"],"method":"eth_sendRawTransaction"}`)

	var rpcMessage rpc.JsonrpcMessage
	_ = json.Unmarshal(contents, &rpcMessage)

	var params []string
	_ = json.Unmarshal(rpcMessage.Params, &params)

    svc := new(RelaySignerService)
    if err := svc.Init(cfg); err != nil {
        t.Fatalf("Init failed: %v", err)
    }
    // ensure we hit the mock endpoint
    svc.Config.Application.NodeURL = srv.URL + "/sendMetatransaction"

    // 3) Seed the map so every call goes through incrementTransactionCount
    //svc.senders = map[string]*big.Int{senderAddr: big.NewInt(0)}
	svc.senders = make(map[string]*big.Int)
    // 4) Spawn N workers (goroutines) all calling SendMetatransaction
    to := common.HexToAddress("0x82a978b3f5962a5b0957d9ee9eef472ee55b42f1")
    data, _ := hex.DecodeString("0xf861808082ea6094fd32cfc2e71611626d6368a41f915d0077a306a180b8446057361d000000000000000000000000000000000000000000000000000000000000003c000000000000000000000000173cf75f0905338597fcd38f5ce13e6840b230e9")         // empty payload is fine
	var gasLimit uint64
	var nonce uint64
	var r [32]byte
	var s [32]byte

	gasLimit = 200000
	nonce=1
	senderAddr:="0x0000000000000000000000000000000000000000"

    var wg sync.WaitGroup
    wg.Add(workers)
    for i := 0; i < workers; i++ {
        go func() {
            defer wg.Done()
			svc.SendMetatransaction(rpcMessage.ID, &to, gasLimit, data, 27, r, s, senderAddr, nonce)
        }()
    }
    wg.Wait()

	err := os.Remove("keyMock")
	if err != nil {
		log.Fatal(err)
	}

    //Check COLLISION: only one key should exist
    if got := len(svc.senders); got != 1 {
        t.Fatalf("expected 1 entry on collision, got %d", got)
    }

    // And if you increment per call(worker), it should be 3, sender = nonce(3)
    if cnt := svc.senders[senderAddr].Uint64(); cnt != uint64(workers) {
        t.Errorf("expected count=%d, got %d", workers, cnt)
    }	
}

func TestNonceRaceSendMetatransactionDifferentSenders(t *testing.T) {
    
	const workers = 3
    
    // 1) Stand up your same serverMock
    srv := serverMock()
    defer srv.Close()

    // 2) Prepare the service exactly as in TestSendMetatransaction
    dir, _ := os.Getwd()
    createKeyMock(dir + "/keyMock")
    setKeyMock()

    cfg := &model.Config{
        Application: model.ApplicationConfig{
            NodeURL:         srv.URL + "/getRelayHubContract",
            ContractAddress: "0x0ae2Da68515Ef8DC4bBCa1fA1bcE00C508b2Af4B",
            NodeKeyPath:     dir + "/keyMock",
        },
    }

	contents := []byte(`{"id":2914410858336929,"jsonrpc":"2.0","params":["0xf8840180831e8480946e6bbf31aa45042d53128339383fcd1c377b42c780a46057361d00000000000000000000000000000000000000000000000000000000000001591ba028934b543809922b277e85f6bcf7b1f25e937de05c5138e17fdfa480ba74e84ba055a2a611763ffcb748547408093551928c9549f95a0a9cabd3b1f1f2e166cc16"],"method":"eth_sendRawTransaction"}`)

	var rpcMessage rpc.JsonrpcMessage
	_ = json.Unmarshal(contents, &rpcMessage)

	var params []string
	_ = json.Unmarshal(rpcMessage.Params, &params)

    svc := new(RelaySignerService)
    if err := svc.Init(cfg); err != nil {
        t.Fatalf("Init failed: %v", err)
    }
    // ensure we hit the mock endpoint
    svc.Config.Application.NodeURL = srv.URL + "/sendMetatransaction"

    // 3) Seed the map so every call goes through incrementTransactionCount
    //svc.senders = map[string]*big.Int{senderAddr: big.NewInt(0)}
	svc.senders = make(map[string]*big.Int)
    // 4) Spawn N workers (goroutines) all calling SendMetatransaction
    to := common.HexToAddress("0x82a978b3f5962a5b0957d9ee9eef472ee55b42f1")
    data, _ := hex.DecodeString("0xf861808082ea6094fd32cfc2e71611626d6368a41f915d0077a306a180b8446057361d000000000000000000000000000000000000000000000000000000000000003c000000000000000000000000173cf75f0905338597fcd38f5ce13e6840b230e9")         // empty payload is fine
	var gasLimit uint64
	var nonce uint64
	var r [32]byte
	var s [32]byte

	gasLimit = 200000
	nonce=1

    var wg sync.WaitGroup
    wg.Add(workers)
    for i := 0; i < workers; i++ {
        go func(i int) {
            defer wg.Done()
			senderAddr := fmt.Sprintf("0x%040x", i)
			svc.SendMetatransaction(rpcMessage.ID, &to, gasLimit, data, 27, r, s, senderAddr, nonce)
			//fmt.Println("=Sender=")
			//fmt.Println(senderAddr)
			//fmt.Println("=Nonce=")
			//fmt.Println(i)
        }(i)
    }
    wg.Wait()

	err := os.Remove("keyMock")
	if err != nil {
		log.Fatal(err)
	}

    // 5) Check how many entries we actually got back
	 // RACE: each goroutine gets its own key
	 if got, want := len(svc.senders), workers; got != want {
        t.Fatalf("expected %d distinct entries, got %d", want, got)
    }

    // And each one saw exactly nonce==1
    for i := 0; i < workers; i++ {
        sender := fmt.Sprintf("0x%040x", i)
        if cnt := svc.senders[sender].Uint64(); cnt != 1 {
            t.Errorf("for %s: expected count=1, got %d", sender, cnt)
        }
    }
}

func serverMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/getTransactionCount", mockGetNonce)
	handler.HandleFunc("/getReceipt", mockGetReceipt)
	handler.HandleFunc("/getReceiptRevertReason", mockGetReceiptRevertReason)
	handler.HandleFunc("/sendMetatransaction", mockSendMetatransaction)
	handler.HandleFunc("/getRelayHubContract", mockGetRelayHubContract)

	srv := httptest.NewServer(handler)

	return srv
}

func mockGetNonce(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`{"jsonrpc" : "2.0","id" : 53,"result" : "0x0000000000000000000000000000000000000000000000000000000000000159"}`))
}

func mockSendMetatransaction(w http.ResponseWriter, r *http.Request) {
	switch sequence {
	case 0, 1, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 16, 17, 18, 19, 20, 21, 22, 23:
		_, _ = w.Write([]byte(`{"jsonrpc" : "2.0","id" : 53,"result" : "0x6"}`))
	case 2:
		_, _ = w.Write([]byte(`{"jsonrpc" : "2.0","id" : 53,"result" : "0x6"}`))
	case 15:
		_, _ = w.Write([]byte(`{"jsonrpc" : "2.0","id" : 53,"result" : "0x6"}`))
	}
	sequence++
}

func mockGetReceipt(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`{
		"jsonrpc" : "2.0",
		"id" : 53,
		"result" : {
		  "blockHash" : "0x6e3aa24e261e61832624749b64049104c6105ba870d3375484548ffdb133eeea",
		  "blockNumber" : "0xaae545",
		  "contractAddress" : null,
		  "cumulativeGasUsed" : "0x309f0",
		  "from" : "0xd00e6624a73f88b39f82ab34e8bf2b4d226fd768",
		  "gasUsed" : "0x309f0",
		  "logs" : [ {
			"address" : "0xff6d55d01fb12695ea00c071ad8af3ce44cf3a91",
			"topics" : [ "0xa37b1b27143f61d990cfcf145e7f5d21c4419700613094ab29654b7ac6c08724" ],
			"data" : "0x0000000000000000000000000000000000000000000000000000000000000001",
			"blockNumber" : "0xaae545",
			"transactionHash" : "0x41167872ab8e13bf7ea5ea366786da656b3f32181410523b97ffecf0ee9cd945",
			"transactionIndex" : "0x0",
			"blockHash" : "0x6e3aa24e261e61832624749b64049104c6105ba870d3375484548ffdb133eeea",
			"logIndex" : "0x0",
			"removed" : false
		  }, {
			"address" : "0xff6d55d01fb12695ea00c071ad8af3ce44cf3a91",
			"topics" : [ "0x1ecdaca0ae98a95eed765c0622982b0f7691f9a345988f8fca91c1c016ce5ee7" ],
			"data" : "0x0000000000000000000000000000000000000000000000000000000000aae545000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000009896800",
			"blockNumber" : "0xaae545",
			"transactionHash" : "0x41167872ab8e13bf7ea5ea366786da656b3f32181410523b97ffecf0ee9cd945",
			"transactionIndex" : "0x0",
			"blockHash" : "0x6e3aa24e261e61832624749b64049104c6105ba870d3375484548ffdb133eeea",
			"logIndex" : "0x1",
			"removed" : false
		  }, {
			"address" : "0xff6d55d01fb12695ea00c071ad8af3ce44cf3a91",
			"topics" : [ "0x79f72f9dacecfa9af3cfe946364971d0ef4826ffd35451658b283d58a382c20f", "0x000000000000000000000000a20aa371a9d05bba5d087bfee8fdf47ffe1088da", "0x000000000000000000000000d00e6624a73f88b39f82ab34e8bf2b4d226fd768" ],
			"data" : "0x",
			"blockNumber" : "0xaae545",
			"transactionHash" : "0x41167872ab8e13bf7ea5ea366786da656b3f32181410523b97ffecf0ee9cd945",
			"transactionIndex" : "0x0",
			"blockHash" : "0x6e3aa24e261e61832624749b64049104c6105ba870d3375484548ffdb133eeea",
			"logIndex" : "0x2",
			"removed" : false
		  }, {
			"address" : "0x91402a50b130cb6ee76b1c85704faf94361cc233",
			"topics" : [ "0xeaf540d6ee39a98c4ab8d5d07d678c306272e18b51a3c93b026c4a2ce84a7afd" ],
			"data" : "0x000000000000000000000000ff6d55d01fb12695ea00c071ad8af3ce44cf3a9100000000000000000000000000000000000000000000000000000000000000430000000000000000000000000000000000000000000000000000000000000043",
			"blockNumber" : "0xaae545",
			"transactionHash" : "0x41167872ab8e13bf7ea5ea366786da656b3f32181410523b97ffecf0ee9cd945",
			"transactionIndex" : "0x0",
			"blockHash" : "0x6e3aa24e261e61832624749b64049104c6105ba870d3375484548ffdb133eeea",
			"logIndex" : "0x3",
			"removed" : false
		  }, {
			"address" : "0xff6d55d01fb12695ea00c071ad8af3ce44cf3a91",
			"topics" : [ "0xfed6f0abc4f5e1923377ee51313db072532b591ea23ea4b4c44a4457e7e5f417", "0x000000000000000000000000d00e6624a73f88b39f82ab34e8bf2b4d226fd768", "0x000000000000000000000000a20aa371a9d05bba5d087bfee8fdf47ffe1088da", "0x00000000000000000000000091402a50b130cb6ee76b1c85704faf94361cc233" ],
			"data" : "0x0000000000000000000000000000000000000000000000000000000000000001",
			"blockNumber" : "0xaae545",
			"transactionHash" : "0x41167872ab8e13bf7ea5ea366786da656b3f32181410523b97ffecf0ee9cd945",
			"transactionIndex" : "0x0",
			"blockHash" : "0x6e3aa24e261e61832624749b64049104c6105ba870d3375484548ffdb133eeea",
			"logIndex" : "0x4",
			"removed" : false
		  }, {
			"address" : "0xff6d55d01fb12695ea00c071ad8af3ce44cf3a91",
			"topics" : [ "0x260359eeed8459102359245337088f93b15364b134b4be9092d508e741bbdee1" ],
			"data" : "0x000000000000000000000000d00e6624a73f88b39f82ab34e8bf2b4d226fd7680000000000000000000000000000000000000000000000000000000000aae5450000000000000000000000000000000000000000000000000000000000004085000000000000000000000000000000000000000000000000000000000989277b0000000000000000000000000000000000000000000000000000000000004085",
			"blockNumber" : "0xaae545",
			"transactionHash" : "0x41167872ab8e13bf7ea5ea366786da656b3f32181410523b97ffecf0ee9cd945",
			"transactionIndex" : "0x0",
			"blockHash" : "0x6e3aa24e261e61832624749b64049104c6105ba870d3375484548ffdb133eeea",
			"logIndex" : "0x5",
			"removed" : false
		  } ],
		  "logsBloom" : "0x0000000000400000000000000000000000000000080000000000000000000000000000000000000000020000010000000000000000000000400000000000000002000000000001000000008000000000100000000020010000100000000000000800000000000000000000000000000000000000000000000000000000000000000000000010000000000000080000000000000000000000000000010000000000004800040000000000000001800000001001400000000000000010000000000000000a000000004008000000000000000000000000000000000000000000000000020000000000000000000000100000000020000000200000000000000000",
		  "status" : "0x1",
		  "to" : "0xff6d55d01fb12695ea00c071ad8af3ce44cf3a91",
		  "transactionHash" : "0x41167872ab8e13bf7ea5ea366786da656b3f32181410523b97ffecf0ee9cd945",
		  "transactionIndex" : "0x0"
		}
	  }`))
}

func mockGetReceiptRevertReason(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`{
		"jsonrpc" : "2.0",
		"id" : 53,
		"result" : {
		  "blockHash" : "0x6e3aa24e261e61832624749b64049104c6105ba870d3375484548ffdb133eeea",
		  "blockNumber" : "0xaae545",
		  "contractAddress" : null,
		  "cumulativeGasUsed" : "0x309f0",
		  "from" : "0xd00e6624a73f88b39f82ab34e8bf2b4d226fd768",
		  "gasUsed" : "0x309f0",
		  "logs" : [ {
			"address" : "0xff6d55d01fb12695ea00c071ad8af3ce44cf3a91",
			"topics" : [ "0xa37b1b27143f61d990cfcf145e7f5d21c4419700613094ab29654b7ac6c08724" ],
			"data" : "0x0000000000000000000000000000000000000000000000000000000000000001",
			"blockNumber" : "0xaae545",
			"transactionHash" : "0x41167872ab8e13bf7ea5ea366786da656b3f32181410523b97ffecf0ee9cd945",
			"transactionIndex" : "0x0",
			"blockHash" : "0x6e3aa24e261e61832624749b64049104c6105ba870d3375484548ffdb133eeea",
			"logIndex" : "0x0",
			"removed" : false
		  }, {
			"address" : "0xff6d55d01fb12695ea00c071ad8af3ce44cf3a91",
			"topics" : [ "0x1ecdaca0ae98a95eed765c0622982b0f7691f9a345988f8fca91c1c016ce5ee7" ],
			"data" : "0x0000000000000000000000000000000000000000000000000000000000aae545000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000009896800",
			"blockNumber" : "0xaae545",
			"transactionHash" : "0x41167872ab8e13bf7ea5ea366786da656b3f32181410523b97ffecf0ee9cd945",
			"transactionIndex" : "0x0",
			"blockHash" : "0x6e3aa24e261e61832624749b64049104c6105ba870d3375484548ffdb133eeea",
			"logIndex" : "0x1",
			"removed" : false
		  }, {
			"address" : "0xff6d55d01fb12695ea00c071ad8af3ce44cf3a91",
			"topics" : [ "0x79f72f9dacecfa9af3cfe946364971d0ef4826ffd35451658b283d58a382c20f", "0x000000000000000000000000a20aa371a9d05bba5d087bfee8fdf47ffe1088da", "0x000000000000000000000000d00e6624a73f88b39f82ab34e8bf2b4d226fd768" ],
			"data" : "0x",
			"blockNumber" : "0xaae545",
			"transactionHash" : "0x41167872ab8e13bf7ea5ea366786da656b3f32181410523b97ffecf0ee9cd945",
			"transactionIndex" : "0x0",
			"blockHash" : "0x6e3aa24e261e61832624749b64049104c6105ba870d3375484548ffdb133eeea",
			"logIndex" : "0x2",
			"removed" : false
		  }, {
			"address" : "0x91402a50b130cb6ee76b1c85704faf94361cc233",
			"topics" : [ "0xeaf540d6ee39a98c4ab8d5d07d678c306272e18b51a3c93b026c4a2ce84a7afd" ],
			"data" : "0x000000000000000000000000ff6d55d01fb12695ea00c071ad8af3ce44cf3a9100000000000000000000000000000000000000000000000000000000000000430000000000000000000000000000000000000000000000000000000000000043",
			"blockNumber" : "0xaae545",
			"transactionHash" : "0x41167872ab8e13bf7ea5ea366786da656b3f32181410523b97ffecf0ee9cd945",
			"transactionIndex" : "0x0",
			"blockHash" : "0x6e3aa24e261e61832624749b64049104c6105ba870d3375484548ffdb133eeea",
			"logIndex" : "0x3",
			"removed" : false
		  }, {
			"address" : "0xff6d55d01fb12695ea00c071ad8af3ce44cf3a91",
			"topics" : [ "0xfed6f0abc4f5e1923377ee51313db072532b591ea23ea4b4c44a4457e7e5f417", "0x000000000000000000000000d00e6624a73f88b39f82ab34e8bf2b4d226fd768", "0x000000000000000000000000a20aa371a9d05bba5d087bfee8fdf47ffe1088da", "0x00000000000000000000000091402a50b130cb6ee76b1c85704faf94361cc233" ],
			"data" : "0x0000000000000000000000000000000000000000000000000000000000000001",
			"blockNumber" : "0xaae545",
			"transactionHash" : "0x41167872ab8e13bf7ea5ea366786da656b3f32181410523b97ffecf0ee9cd945",
			"transactionIndex" : "0x0",
			"blockHash" : "0x6e3aa24e261e61832624749b64049104c6105ba870d3375484548ffdb133eeea",
			"logIndex" : "0x4",
			"removed" : false
		  }, {
			"address" : "0xff6d55d01fb12695ea00c071ad8af3ce44cf3a91",   
		    "topics": [
			"0x548af85d7bc344f47cbfacdfce1ffea1ecd862e5e235ca9ec919e767c14049a8",
			"0x00000000000000000000000063949701cd0e1cc04dfea0afbf410968f10ff4b6",
			"0x000000000000000000000000bceda2ba9af65c18c7992849c312d1db77cf008e",
			"0x000000000000000000000000938144efd1b3943c3b6658f4f7b72fcd980c55a1"
		    ],
			"data": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000006408c379a0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000096e616275636f646f73000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",		    
			"blockNumber" : "0xaae545",
			"transactionHash" : "0x41167872ab8e13bf7ea5ea366786da656b3f32181410523b97ffecf0ee9cd945",
			"transactionIndex" : "0x0",
			"blockHash" : "0x6e3aa24e261e61832624749b64049104c6105ba870d3375484548ffdb133eeea",
			"logIndex" : "0x5",
			"removed" : false
		  }, {
			"address" : "0xff6d55d01fb12695ea00c071ad8af3ce44cf3a91",
			"topics" : [ "0x260359eeed8459102359245337088f93b15364b134b4be9092d508e741bbdee1" ],
			"data" : "0x000000000000000000000000d00e6624a73f88b39f82ab34e8bf2b4d226fd7680000000000000000000000000000000000000000000000000000000000aae5450000000000000000000000000000000000000000000000000000000000004085000000000000000000000000000000000000000000000000000000000989277b0000000000000000000000000000000000000000000000000000000000004085",
			"blockNumber" : "0xaae545",
			"transactionHash" : "0x41167872ab8e13bf7ea5ea366786da656b3f32181410523b97ffecf0ee9cd945",
			"transactionIndex" : "0x0",
			"blockHash" : "0x6e3aa24e261e61832624749b64049104c6105ba870d3375484548ffdb133eeea",
			"logIndex" : "0x5",
			"removed" : false
		  } ],
		  "logsBloom" : "0x0000000000400000000000000000000000000000080000000000000000000000000000000000000000020000010000000000000000000000400000000000000002000000000001000000008000000000100000000020010000100000000000000800000000000000000000000000000000000000000000000000000000000000000000000010000000000000080000000000000000000000000000010000000000004800040000000000000001800000001001400000000000000010000000000000000a000000004008000000000000000000000000000000000000000000000000020000000000000000000000100000000020000000200000000000000000",
		  "status" : "0x1",
		  "to" : "0xff6d55d01fb12695ea00c071ad8af3ce44cf3a91",
		  "transactionHash" : "0x41167872ab8e13bf7ea5ea366786da656b3f32181410523b97ffecf0ee9cd945",
		  "transactionIndex" : "0x0",
		  "revertReason" : "0x08c379a0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000096e616275636f646f730000000000000000000000000000000000000000000000"	
		}
	  }`))
}

func mockGetRelayHubContract(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`{"jsonrpc" : "2.0","id" : 53,"result" : "0x000000000000000000000000ff6d55d01fb12695ea00c071ad8af3ce44cf3a91"}`))
}

func createKeyMock(path string) {
	d1 := []byte("0xb3e7374dca5ca90c3899dbb2c978051437fb15534c945bf59df16d6c80be27c0")
	err := ioutil.WriteFile(path, d1, 0644)
	if err != nil {
		fmt.Println("err:", err)
	}
}

func setKeyMock() {
	os.Setenv("WRITER_KEY", "0xb3e7374dca5ca90c3899dbb2c978051437fb15534c945bf59df16d6c80be27c0")
}
