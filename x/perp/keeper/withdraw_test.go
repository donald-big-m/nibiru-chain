package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/NibiruChain/nibiru/x/perp/types"
	"github.com/NibiruChain/nibiru/x/testutil/sample"
)

func TestWithdraw(t *testing.T) {
	tests := []struct {
		name                  string
		initialVaultBalance   int64
		initialPrepaidBadDebt int64
		amountToWithdraw      int64

		expectedPerpEFWithdrawal    int64
		expectedFinalPrepaidBadDebt int64
	}{
		{
			name:                  "no bad debt",
			initialVaultBalance:   100,
			initialPrepaidBadDebt: 0,

			amountToWithdraw: 10,

			expectedPerpEFWithdrawal:    0,
			expectedFinalPrepaidBadDebt: 0,
		},
		{
			name:                  "creates prepaid bad debt",
			initialVaultBalance:   9,
			initialPrepaidBadDebt: 0,

			amountToWithdraw: 10,

			expectedPerpEFWithdrawal:    1,
			expectedFinalPrepaidBadDebt: 1,
		},
		{
			name:                  "increases existing prepaid bad debt",
			initialVaultBalance:   9,
			initialPrepaidBadDebt: 1,

			amountToWithdraw: 10,

			expectedPerpEFWithdrawal:    1,
			expectedFinalPrepaidBadDebt: 2,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Log("initialize variables")
			perpKeeper, mocks, ctx := getKeeper(t)
			receiver := sample.AccAddress()
			denom := "NUSD"

			t.Log("mock account keeper")
			vaultAddr := authtypes.NewModuleAddress(types.VaultModuleAccount)
			mocks.mockAccountKeeper.EXPECT().GetModuleAddress(
				types.VaultModuleAccount).
				Return(vaultAddr)

			t.Log("mock bank keeper")
			mocks.mockBankKeeper.EXPECT().GetBalance(ctx, vaultAddr, denom).
				Return(sdk.NewInt64Coin(denom, tc.initialVaultBalance))
			mocks.mockBankKeeper.EXPECT().SendCoinsFromModuleToAccount(
				ctx, types.VaultModuleAccount, receiver,
				sdk.NewCoins(sdk.NewInt64Coin(denom, tc.amountToWithdraw)),
			).Return(nil)
			if tc.expectedPerpEFWithdrawal > 0 {
				mocks.mockBankKeeper.EXPECT().SendCoinsFromModuleToModule(
					ctx, types.PerpEFModuleAccount, types.VaultModuleAccount,
					sdk.NewCoins(sdk.NewInt64Coin(denom, tc.expectedPerpEFWithdrawal)),
				).Return(nil)
			}

			t.Log("initial prepaid bad debt")
			perpKeeper.PrepaidBadDebtState().Set(ctx, denom, sdk.NewInt(tc.initialPrepaidBadDebt))

			t.Log("execute withdrawal")
			err := perpKeeper.Withdraw(ctx, denom, receiver, sdk.NewInt(tc.amountToWithdraw))
			require.NoError(t, err)

			t.Log("assert new prepaid bad debt")
			prepaidBadDebt := perpKeeper.PrepaidBadDebtState().Get(ctx, denom)
			assert.EqualValues(t, tc.expectedFinalPrepaidBadDebt, prepaidBadDebt.Int64())
		})
	}
}
