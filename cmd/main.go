package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/yearn/ydaemon/common/ethereum"
	"github.com/yearn/ydaemon/common/logs"
	"github.com/yearn/ydaemon/internal"
	"github.com/yearn/ydaemon/internal/storage"
)

func processServer(chainID uint64, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	setStatusForChainID(chainID, `Loading`)
	defer setStatusForChainID(chainID, `OK`)

	logs.Info(`Initializing chain ` + strconv.FormatUint(chainID, 10) + ` indexing process`)

	logs.Info(`Setting up WebSocket client for chain ` + strconv.FormatUint(chainID, 10))
	ethereum.GetWSClient(chainID, true)

	logs.Info(`Initializing block timestamps for chain ` + strconv.FormatUint(chainID, 10))
	ethereum.InitBlockTimestamp(chainID)

	logs.Info(`Starting main indexer for chain ` + strconv.FormatUint(chainID, 10))
	internal.InitializeV2(chainID, nil)

	logs.Info(`Chain ` + strconv.FormatUint(chainID, 10) + ` initialization completed`)
	TriggerInitializedStatus(chainID)
}

/**************************************************************************************************
** Main entry point for the daemon, handling everything from initialization to running external
** processes.
**************************************************************************************************/
func main() {
	initFlags()
	ethereum.Initialize()
	storage.InitializeStorage()
	// go ListenToSignals()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logs.Info(`Starting indexing processes for ` + strconv.Itoa(len(chains)) + ` chains: ` + fmt.Sprintf("%v", chains))

	// Wait for all chains to complete initial data loading before starting the HTTP server
	var wg sync.WaitGroup
	for _, chainID := range chains {
		wg.Add(1)
		go processServer(chainID, &wg)
	}

	logs.Info(`Waiting for all chains to complete initial data loading...`)
	wg.Wait()
	logs.Success(`All chains initialized successfully`)

	// Now that all data is loaded, start the HTTP server
	logs.Info(`Starting HTTP server on port ` + port)
	go NewRouter().Run(`:` + port)
	// go TriggerTgMessage(`💛 - yDaemon v` + GetVersion() + ` is ready to accept requests: https://ydaemon.yearn.fi/`)

	logs.Success(`Server ready on port ` + port + ` !`)
	select {}
}
