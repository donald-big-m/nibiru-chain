package keeper

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/NibiruChain/nibiru/x/common"
	epochtypes "github.com/NibiruChain/nibiru/x/epochs/types"
	"github.com/NibiruChain/nibiru/x/perp/types"
	testutilevents "github.com/NibiruChain/nibiru/x/testutil/events"
)

func TestEndOfEpochTwapCalculation(t *testing.T) {
	tests := []struct {
		name                            string
		indexPrice                      sdk.Dec
		markPrice                       sdk.Dec
		expectedCumulativeFundingRates  []sdk.Dec
		expectedFundingRateChangedEvent *types.FundingRateChangedEvent
	}{
		{
			name:                            "check empty prices",
			indexPrice:                      sdk.ZeroDec(),
			markPrice:                       sdk.ZeroDec(),
			expectedCumulativeFundingRates:  []sdk.Dec{sdk.ZeroDec()},
			expectedFundingRateChangedEvent: nil,
		},
		{
			name:                            "empty index price",
			indexPrice:                      sdk.ZeroDec(),
			markPrice:                       sdk.NewDec(10),
			expectedCumulativeFundingRates:  []sdk.Dec{sdk.ZeroDec()},
			expectedFundingRateChangedEvent: nil,
		},
		{
			name:                            "empty mark price",
			indexPrice:                      sdk.NewDec(10),
			markPrice:                       sdk.ZeroDec(),
			expectedCumulativeFundingRates:  []sdk.Dec{sdk.ZeroDec()},
			expectedFundingRateChangedEvent: nil,
		},
		{
			name:                           "equal prices",
			indexPrice:                     sdk.NewDec(10),
			markPrice:                      sdk.NewDec(10),
			expectedCumulativeFundingRates: []sdk.Dec{sdk.ZeroDec(), sdk.ZeroDec()},
			expectedFundingRateChangedEvent: &types.FundingRateChangedEvent{
				Pair:                  common.Pair_BTC_NUSD.String(),
				MarkPrice:             sdk.NewDec(10),
				IndexPrice:            sdk.NewDec(10),
				LatestFundingRate:     sdk.ZeroDec(),
				CumulativeFundingRate: sdk.ZeroDec(),
				BlockHeight:           1,
				BlockTimeMs:           1,
			},
		},
		{
			name:                           "calculate funding rate with higher index price",
			markPrice:                      sdk.NewDec(19),
			indexPrice:                     sdk.NewDec(462),
			expectedCumulativeFundingRates: []sdk.Dec{sdk.ZeroDec(), sdk.MustNewDecFromStr("-0.039953102453102453")},
			expectedFundingRateChangedEvent: &types.FundingRateChangedEvent{
				Pair:                  common.Pair_BTC_NUSD.String(),
				MarkPrice:             sdk.NewDec(19),
				IndexPrice:            sdk.NewDec(462),
				LatestFundingRate:     sdk.MustNewDecFromStr("-0.039953102453102453"),
				CumulativeFundingRate: sdk.MustNewDecFromStr("-0.039953102453102453"),
				BlockHeight:           1,
				BlockTimeMs:           1,
			},
		},
		{
			name:                           "calculate funding rate with higher mark price",
			markPrice:                      sdk.NewDec(745),
			indexPrice:                     sdk.NewDec(64),
			expectedCumulativeFundingRates: []sdk.Dec{sdk.ZeroDec(), sdk.MustNewDecFromStr("0.443359375000000000")},
			expectedFundingRateChangedEvent: &types.FundingRateChangedEvent{
				Pair:                  common.Pair_BTC_NUSD.String(),
				MarkPrice:             sdk.NewDec(745),
				IndexPrice:            sdk.NewDec(64),
				LatestFundingRate:     sdk.MustNewDecFromStr("0.443359375000000000"),
				CumulativeFundingRate: sdk.MustNewDecFromStr("0.443359375000000000"),
				BlockHeight:           1,
				BlockTimeMs:           1,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			perpKeeper, mocks, ctx := getKeeper(t)
			ctx = ctx.WithBlockHeight(1).WithBlockTime(time.UnixMilli(1))

			t.Log("initialize params")
			initParams(ctx, perpKeeper)

			t.Log("set mocks")
			setMockPrices(ctx, mocks, tc.indexPrice, tc.markPrice)

			perpKeeper.AfterEpochEnd(ctx, "30 min", 1)

			t.Log("assert PairMetadataState")
			pair, err := perpKeeper.PairsMetadata.Get(ctx, common.Pair_BTC_NUSD)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedCumulativeFundingRates, pair.CumulativeFundingRates)

			if tc.expectedFundingRateChangedEvent != nil {
				t.Log("assert FundingRateChangedEvent")
				testutilevents.RequireContainsTypedEvent(t, ctx, tc.expectedFundingRateChangedEvent)
			}
		})
	}
}

func initParams(ctx sdk.Context, k Keeper) {
	k.SetParams(ctx, types.Params{
		Stopped:                 false,
		FeePoolFeeRatio:         sdk.MustNewDecFromStr("0.00001"),
		EcosystemFundFeeRatio:   sdk.MustNewDecFromStr("0.000005"),
		LiquidationFeeRatio:     sdk.MustNewDecFromStr("0.000007"),
		PartialLiquidationRatio: sdk.MustNewDecFromStr("0.00001"),
		FundingRateInterval:     "30 min",
		TwapLookbackWindow:      15 * time.Minute,
	})
	setPairMetadata(k, ctx, types.PairMetadata{
		Pair: common.Pair_BTC_NUSD,
		// start with one entry to ensure we append
		CumulativeFundingRates: []sdk.Dec{sdk.ZeroDec()},
	})
}

func setMockPrices(ctx sdk.Context, mocks mockedDependencies, indexPrice sdk.Dec, markPrice sdk.Dec) {
	mocks.mockVpoolKeeper.EXPECT().ExistsPool(ctx, common.Pair_BTC_NUSD).Return(true)

	mocks.mockEpochKeeper.EXPECT().GetEpochInfo(ctx, "30 min").Return(
		epochtypes.EpochInfo{Duration: time.Hour},
	).MaxTimes(1)

	mocks.mockPricefeedKeeper.EXPECT().
		GetCurrentTWAP(ctx, common.Pair_BTC_NUSD.Token0, common.Pair_BTC_NUSD.Token1).Return(indexPrice, nil).MaxTimes(1)

	mocks.mockVpoolKeeper.EXPECT().
		GetSpotTWAP(ctx, common.Pair_BTC_NUSD, 15*time.Minute).
		Return(markPrice, nil).MaxTimes(1)
}
