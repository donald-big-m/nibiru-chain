package assertion

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/NibiruChain/nibiru/app"
)

func AllBalancesEqual(account sdk.AccAddress, amount sdk.Coins) *allBalancesEqualAction {
	return &allBalancesEqualAction{Account: account, Amount: amount}
}

type allBalancesEqualAction struct {
	Account sdk.AccAddress
	Amount  sdk.Coins
}

func (b allBalancesEqualAction) Do(app *app.NibiruApp, ctx sdk.Context) (sdk.Context, error, bool) {
	coins := app.BankKeeper.GetAllBalances(ctx, b.Account)
	if !coins.IsEqual(b.Amount) {
		return ctx, fmt.Errorf(
			"account %s balance not equal, expected %s, got %s",
			b.Account.String(),
			b.Amount.String(),
			coins.String(),
		), false
	}

	return ctx, nil, false
}

func BalanceEqual(account sdk.AccAddress, denom string, amount sdk.Int) *balanceEqualAction {
	return &balanceEqualAction{Account: account, Denom: denom, Amount: amount}
}

type balanceEqualAction struct {
	Account sdk.AccAddress
	Denom   string
	Amount  sdk.Int
}

func (b balanceEqualAction) Do(app *app.NibiruApp, ctx sdk.Context) (sdk.Context, error, bool) {
	coin := app.BankKeeper.GetBalance(ctx, b.Account, b.Denom)
	if !coin.Amount.Equal(b.Amount) {
		return ctx, fmt.Errorf(
			"account %s balance not equal, expected %s, got %s",
			b.Account.String(),
			b.Amount.String(),
			coin.String(),
		), false
	}

	return ctx, nil, false
}

func ModuleBalanceEqual(moduleName string, denom string, amount sdk.Int) *moduleBalanceAction {
	return &moduleBalanceAction{ModuleName: moduleName, Denom: denom, Amount: amount}
}

type moduleBalanceAction struct {
	ModuleName string
	Denom      string
	Amount     sdk.Int
}

func (b moduleBalanceAction) Do(app *app.NibiruApp, ctx sdk.Context) (sdk.Context, error, bool) {
	coin := app.BankKeeper.GetBalance(ctx, app.AccountKeeper.GetModuleAddress(b.ModuleName), b.Denom)
	if !coin.Amount.Equal(b.Amount) {
		return ctx, fmt.Errorf(
			"module %s balance not equal, expected %s, got %s",
			b.ModuleName,
			b.Amount.String(),
			coin.String(),
		), false
	}

	return ctx, nil, false
}
