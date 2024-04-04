package types

import (
	sdkmath "cosmossdk.io/math"
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/NibiruChain/nibiru/x/common/asset"
	"github.com/NibiruChain/nibiru/x/common/denoms"
	epochstypes "github.com/NibiruChain/nibiru/x/epochs/types"

	"github.com/cosmos/cosmos-sdk/codec"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Markets:          []Market{},
		Amms:             []AMM{},
		Positions:        []GenesisPosition{},
		ReserveSnapshots: []ReserveSnapshot{},
		CollateralDenom:  TestingCollateralDenomNUSD,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	for _, m := range gs.Markets {
		if err := m.Validate(); err != nil {
			return err
		}
	}

	for _, m := range gs.Amms {
		if err := m.Validate(); err != nil {
			return err
		}
	}

	// TODO: validate positions
	//for _, pos := range gs.Positions {
	//	if err := pos.Validate(); err != nil {
	//		return err
	//	}
	//}

	return nil
}

func DefaultMarket(pair asset.Pair) Market {
	return Market{
		Pair:                            pair,
		Enabled:                         false,
		Version:                         1,
		LatestCumulativePremiumFraction: sdkmath.LegacyZeroDec(),
		ExchangeFeeRatio:                sdkmath.LegacyMustNewDecFromStr("0.0010"),
		EcosystemFundFeeRatio:           sdkmath.LegacyMustNewDecFromStr("0.0010"),
		LiquidationFeeRatio:             sdkmath.LegacyMustNewDecFromStr("0.0500"),
		PartialLiquidationRatio:         sdkmath.LegacyMustNewDecFromStr("0.5000"),
		FundingRateEpochId:              epochstypes.ThirtyMinuteEpochID,
		MaxFundingRate:                  sdkmath.LegacyNewDec(1),
		TwapLookbackWindow:              time.Minute * 30,
		PrepaidBadDebt:                  sdk.NewCoin(TestingCollateralDenomNUSD, sdkmath.ZeroInt()),
		MaintenanceMarginRatio:          sdkmath.LegacyMustNewDecFromStr("0.0625"),
		MaxLeverage:                     sdkmath.LegacyNewDec(10),
		OraclePair:                      asset.NewPair(pair.BaseDenom(), denoms.USD),
	}
}

func GetGenesisStateFromAppState(cdc codec.JSONCodec, appState map[string]json.RawMessage) *GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return &genesisState
}
