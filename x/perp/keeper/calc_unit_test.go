package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/NibiruChain/nibiru/x/common"
	"github.com/NibiruChain/nibiru/x/perp/types"
	"github.com/NibiruChain/nibiru/x/testutil/sample"
)

func Test_calcFreeCollateral(t *testing.T) {
	testCases := []struct {
		name string
		test func()
	}{
		{
			name: "invalid token pair - error",
			test: func() {
				k, _, ctx := getKeeper(t)
				fundingPayment := sdk.ZeroInt()
				the3pool := "dai:usdc:usdt"
				alice := sample.AccAddress()
				pos := types.ZeroPosition(ctx, common.TokenPair(the3pool), alice.String())
				_, err := k.calcFreeCollateral(ctx, *pos, fundingPayment)
				assert.Error(t, err)
				assert.ErrorContains(t, err, common.ErrInvalidTokenPair.Error())
			},
		},
		{
			name: "token pair not found - error",
			test: func() {
				k, mocks, ctx := getKeeper(t)

				fundingPayment := sdk.ZeroInt()
				validPair := common.TokenPair("xxx:yyy")
				alice := sample.AccAddress()
				pos := types.ZeroPosition(ctx, validPair, alice.String())
				mocks.mockVpoolKeeper.EXPECT().ExistsPool(ctx, validPair).
					Return(false)
				_, err := k.calcFreeCollateral(ctx, *pos, fundingPayment)
				assert.Error(t, err)
				assert.ErrorContains(t, err, types.ErrPairNotFound.Error())
			},
		},
		{
			name: "zero free collateral, zero position - happy path",
			test: func() {
				k, mocks, ctx := getKeeper(t)

				fundingPayment := sdk.ZeroInt()
				validPair := common.TokenPair("xxx:yyy")
				alice := sample.AccAddress()
				pos := types.ZeroPosition(ctx, validPair, alice.String())
				mocks.mockVpoolKeeper.EXPECT().ExistsPool(ctx, validPair).
					Return(true)
				freeCollateral, err := k.calcFreeCollateral(ctx, *pos, fundingPayment)
				assert.NoError(t, err)
				assert.EqualValues(t, sdk.ZeroInt(), freeCollateral)
			},
		},
		{
			name: "negative free collateral, zero position - happy path",
			test: func() {
				k, mocks, ctx := getKeeper(t)

				fundingPayment := sdk.NewInt(10)
				validPair := common.TokenPair("xxx:yyy")
				alice := sample.AccAddress()
				pos := types.ZeroPosition(ctx, validPair, alice.String())
				mocks.mockVpoolKeeper.EXPECT().ExistsPool(ctx, validPair).
					Return(true)
				freeCollateral, err := k.calcFreeCollateral(ctx, *pos, fundingPayment)
				assert.NoError(t, err)
				assert.EqualValues(t, sdk.NewInt(-10), freeCollateral)
			},
		},
		{
			name: "positive free collateral, zero position - happy path",
			test: func() {
				k, mocks, ctx := getKeeper(t)

				fundingPayment := sdk.NewInt(-100)
				validPair := common.TokenPair("xxx:yyy")
				alice := sample.AccAddress()
				pos := types.ZeroPosition(ctx, validPair, alice.String())
				mocks.mockVpoolKeeper.EXPECT().ExistsPool(ctx, validPair).
					Return(true)
				freeCollateral, err := k.calcFreeCollateral(ctx, *pos, fundingPayment)
				assert.NoError(t, err)
				assert.EqualValues(t, sdk.NewInt(100), freeCollateral)
			},
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			tc.test()
		})
	}
}
