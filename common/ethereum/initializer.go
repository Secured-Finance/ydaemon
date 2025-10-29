package ethereum

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/yearn/ydaemon/common/env"
	"github.com/yearn/ydaemon/common/logs"
)

/**************************************************************************************************
** authTransport is a custom HTTP transport that adds Authorization header to all requests
***************************************************************************************************/
type authTransport struct {
	authHeader string
	Transport  http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	reqClone := req.Clone(req.Context())
	reqClone.Header.Set("Authorization", t.authHeader)
	return t.Transport.RoundTrip(reqClone)
}

/**************************************************************************************************
** The init function is a special function triggered directly on execution of the package.
** It is used to initialize the package.
** This init is responsible of creating the RPC clients for all the chains supported by yDaemon
** and storing them in the RPC map.
** Then it will create the multicall clients for each chain.
** Then, it will init the blockTimeSyncMap for all the chains.
***************************************************************************************************/
func Initialize() {
	// Check if verbose blocktime logging is enabled
	if os.Getenv("VERBOSE_BLOCKTIME") == "true" {
		EnableVerboseBlocktime()
	}

	// Create the RPC client for all the chains supported by yDaemon
	for _, chain := range env.GetChains() {
		logs.Info(`Dial RPC URI for chain`, chain.ID)
		rpcURI := GetRPCURI(chain.ID)

		// Check if Authorization token is configured for this chain (Bearer authentication)
		authToken := os.Getenv("RPC_AUTH_TOKEN_FOR_" + strconv.FormatUint(chain.ID, 10))

		var client *ethclient.Client
		var err error

		if authToken != "" {
			// Create custom HTTP client with Bearer Authorization header
			httpClient := &http.Client{
				Transport: &authTransport{
					authHeader: "Bearer " + authToken,
					Transport:  http.DefaultTransport,
				},
			}

			// Create RPC client with custom HTTP client
			ctx := context.Background()
			rpcClient, err := rpc.DialOptions(ctx, rpcURI, rpc.WithHTTPClient(httpClient))
			if err != nil {
				logs.Error(err, "Failed to connect to node with auth")
				continue
			}
			client = ethclient.NewClient(rpcClient)
		} else {
			// Use default dial without custom headers
			client, err = ethclient.Dial(rpcURI)
			if err != nil {
				logs.Error(err, "Failed to connect to node")
				continue
			}
		}

		RPC[chain.ID] = client
	}

	// Create the multicall client for all the chains supported by yDaemon
	for _, chain := range env.GetChains() {
		rpcToUse := GetRPCURI(chain.ID)
		multiCallURI, exists := os.LookupEnv("MULTICALL_RPC_URI_FOR_" + strconv.FormatUint(chain.ID, 10))
		if exists {
			rpcToUse = multiCallURI
		}

		MulticallClientForChainID[chain.ID] = NewMulticallWithAuth(
			rpcToUse,
			chain.MulticallContract.Address,
			chain.ID,
		)
	}

	// Initialize the internal block time data storage in background
	logs.Info(`Starting background blocktime data initialization`)
	go func() {
		InitBlockTimeData()
		logs.Info(`Background blocktime data initialization completed`)
	}()

	// Initialize block timestamps for each supported chain  
	logs.Info(`Initializing block timestamps for all chains`)
	for _, chain := range env.GetChains() {
		logs.Info(`Initializing block timestamps for chain ` + strconv.FormatUint(chain.ID, 10))
		InitBlockTimestamp(chain.ID)
	}
	logs.Info(`Completed blockchain initialization for all chains`)
}
