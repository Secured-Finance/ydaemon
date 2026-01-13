package env

import (
	"math"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/yearn/ydaemon/internal/models"
)

var FILECOIN_CALIBRATION = TChain{
	ID:              314159,
	RpcURI:          `https://api.calibration.node.glif.io/rpc/v1`,
	SubgraphURI:     ``,
	EtherscanURI:    `https://filecoin-testnet.blockscout.com/api`,
	MaxBlockRange:   10_000,
	MaxBatchSize:    math.MaxInt64,
	// AvgBlocksPerDay: 28_800, // ~30 second blocks
	AvgBlocksPerDay: 1_440, // ~60 second blocks
	CanUseWebsocket: false,
	LensContract: TContractData{
		Address: common.Address{},
		Block:   0,
	},
	MulticallContract: TContractData{
		Address: common.HexToAddress(`0xcA11bde05977b3631167028862bE2a173976CA11`),
		Block:   0,
	},
	Coin: models.TERC20Token{
		Address:                   DEFAULT_COIN_ADDRESS,
		UnderlyingTokensAddresses: []common.Address{},
		Type:                      models.TokenTypeNative,
		Name:                      `Testnet Filecoin`,
		Symbol:                    `tFIL`,
		DisplayName:               `Testnet Filecoin`,
		DisplaySymbol:             `tFIL`,
		Description:               `Filecoin Calibration is the primary testing network for Filecoin.`,
		Icon:                      BASE_ASSET_URL + strconv.FormatUint(314159, 10) + `/` + DEFAULT_COIN_ADDRESS.Hex() + `/logo-128.png`,
		Decimals:                  18,
		ChainID:                   314159,
	},
	Registries:            []TContractData{
		{
			Address: common.HexToAddress("0x0377b4daDDA86C89A0091772B79ba67d0E5F7198"),
			Version: 4,
			Block:   3_085_546,
			Label:   `YEARN`,
		},
	},
	APROracleContract:     TContractData{Address: common.Address{}, Block: 0},
	ExtraVaults:       []models.TVaultsFromRegistry{
		{
			// USDFC
			ChainID:         314159,
			Address:         common.HexToAddress(`0x641e8d9634a069d3fe32542207ffc233df4eed19`),
			RegistryAddress: common.HexToAddress(`0x0000000000000000000000000000000000000000`),
			TokenAddress:    common.HexToAddress(`0xb3042734b608a1b16e9e86b374a3f3e389b4cdf0`),
			APIVersion:      `3.0.4`,
			BlockNumber:     3088760,
			Type:            models.TokenTypeStandardVault,
			Kind:            models.VaultKindMultiple,
		},
	},
	BlacklistedVaults:     []common.Address{},
	ExtraTokens:           []common.Address{},
	IgnoredTokens:         []common.Address{},
	Curve: TChainCurve{
		RegistryAddress: common.Address{},
		FactoryAddress:  common.Address{},
		PoolsURIs:       []string{},
		GaugesURI:       ``,
	},
	ExtraURI: TChainExtraURI{
		GammaMerklURI: ``,
		PendleCoreURI: ``,
	},
}
