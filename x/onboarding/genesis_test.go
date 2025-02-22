package onboarding_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/version"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/evmos/ethermint/tests"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	althea "github.com/AltheaFoundation/althea-L1/app"
	altheaconfig "github.com/AltheaFoundation/althea-L1/config"
	"github.com/AltheaFoundation/althea-L1/x/onboarding"
	"github.com/AltheaFoundation/althea-L1/x/onboarding/types"
)

type GenesisTestSuite struct {
	suite.Suite

	ctx sdk.Context

	app     *althea.AltheaApp
	genesis types.GenesisState
}

func (suite *GenesisTestSuite) SetupTest() {
	// consensus key
	consAddress := sdk.ConsAddress(tests.GenerateAddress().Bytes())

	suite.app = althea.NewSetup(false, func(app *althea.AltheaApp, gs simapp.GenesisState) simapp.GenesisState {
		evmGenesis := evmtypes.DefaultGenesisState()
		evmGenesis.Params.EvmDenom = altheaconfig.BaseDenom
		evmGenesis.Params.AllowUnprotectedTxs = false

		gs[evmtypes.ModuleName] = app.AppCodec().MustMarshalJSON(evmGenesis)

		return gs
	})

	// nolint: exhaustruct
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{
		Height:          1,
		ChainID:         altheaconfig.DefaultChainID(),
		Time:            time.Now().UTC(),
		ProposerAddress: consAddress.Bytes(),

		Version: tmversion.Consensus{
			Block: version.BlockProtocol,
		},
		LastBlockId: tmproto.BlockID{
			Hash: tmhash.Sum([]byte("block_id")),
			PartSetHeader: tmproto.PartSetHeader{
				Total: 11,
				Hash:  tmhash.Sum([]byte("partset_header")),
			},
		},
		AppHash:            tmhash.Sum([]byte("app")),
		DataHash:           tmhash.Sum([]byte("data")),
		EvidenceHash:       tmhash.Sum([]byte("evidence")),
		ValidatorsHash:     tmhash.Sum([]byte("validators")),
		NextValidatorsHash: tmhash.Sum([]byte("next_validators")),
		ConsensusHash:      tmhash.Sum([]byte("consensus")),
		LastResultsHash:    tmhash.Sum([]byte("last_result")),
	})

	suite.genesis = *types.DefaultGenesisState()
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) TestOnboardingInitGenesis() {
	testCases := []struct {
		name     string
		genesis  types.GenesisState
		expPanic bool
	}{
		{
			"default genesis",
			suite.genesis,
			false,
		},
		{
			"custom genesis - onboarding disabled",
			types.GenesisState{
				// nolint: exhaustruct
				Params: types.Params{
					EnableOnboarding: false,
				},
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			if tc.expPanic {
				suite.Require().Panics(func() {
					onboarding.InitGenesis(suite.ctx, *suite.app.OnboardingKeeper, tc.genesis)
				})
			} else {
				suite.Require().NotPanics(func() {
					onboarding.InitGenesis(suite.ctx, *suite.app.OnboardingKeeper, tc.genesis)
				})

				params := suite.app.OnboardingKeeper.GetParams(suite.ctx)
				suite.Require().Equal(tc.genesis.Params, params)
			}
		})
	}
}

func (suite *GenesisTestSuite) TestOnboardingExportGenesis() {
	onboarding.InitGenesis(suite.ctx, *suite.app.OnboardingKeeper, suite.genesis)

	genesisExported := onboarding.ExportGenesis(suite.ctx, *suite.app.OnboardingKeeper)
	suite.Require().Equal(genesisExported.Params, suite.genesis.Params)
}
