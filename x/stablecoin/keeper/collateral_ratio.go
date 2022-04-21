package keeper

import (
	"context"
	"fmt"

	"github.com/NibiruChain/nibiru/x/common"
	"github.com/NibiruChain/nibiru/x/stablecoin/events"
	"github.com/NibiruChain/nibiru/x/stablecoin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ---------------------------------------------------------------------------
// Collateral Ratio Getters and Setters
// ---------------------------------------------------------------------------

/*
The collateral ratio, or 'collRatio' (sdk.Dec), is a value beteween 0 and 1 that
determines what proportion of collateral and governance token is used during
stablecoin mints and burns.
*/

// GetCollRatio queries the 'collRatio'.
func (k *Keeper) GetCollRatio(ctx sdk.Context) (collRatio sdk.Dec) {
	return sdk.NewDec(k.GetParams(ctx).CollRatio).QuoInt64(1_000_000)
}

/*
SetCollRatio manually sets the 'collRatio'. This method is mainly used for
testing. When the chain is live, the collateral ratio cannot be manually set, only
adjusted by a fixed amount (e.g. 0.25%).
*/
func (k *Keeper) SetCollRatio(ctx sdk.Context, collRatio sdk.Dec) (err error) {
	collRatioTooHigh := collRatio.GT(sdk.OneDec())
	collRatioTooLow := collRatio.IsNegative()
	if collRatioTooHigh {
		return fmt.Errorf("input 'collRatio', %d, is higher than 1", collRatio)
	} else if collRatioTooLow {
		return fmt.Errorf("input 'collRatio', %d, is negative", collRatio)
	}

	params := k.GetParams(ctx)
	// TODO this should be rethought for production
	newParams := types.NewParams(
		collRatio,
		params.GetFeeRatioAsDec(),
		params.GetEfFeeRatioAsDec(),
		params.GetBonusRateRecollAsDec(),
	)
	k.ParamSubspace.SetParamSet(ctx, &newParams)

	return err
}

// ---------------------------------------------------------------------------
// Recollateralize
// ---------------------------------------------------------------------------

/*
GetUSDValForTargetCollRatio is the collateral value in USD needed to reach a target
collateral ratio.
*/
func (k *Keeper) GetUSDValForTargetCollRatio(
	ctx sdk.Context,
) (neededStable sdk.Dec, err error) {
	stableSupply := k.GetSupplyNUSD(ctx)
	targetCollRatio := k.GetCollRatio(ctx)
	moduleAddr := k.AccountKeeper.GetModuleAddress(types.ModuleName)
	moduleCoins := k.BankKeeper.SpendableCoins(ctx, moduleAddr)
	collDenoms := []string{common.CollDenom}

	currentTotalCollUSD := sdk.ZeroDec()

	for _, collDenom := range collDenoms {
		amtColl := moduleCoins.AmountOf(collDenom)
		priceColl, err := k.PriceKeeper.GetCurrentPrice(
			ctx, collDenom, common.StableDenom)
		if err != nil {
			return sdk.ZeroDec(), err
		}
		collUSD := priceColl.Price.MulInt(amtColl)
		currentTotalCollUSD = currentTotalCollUSD.Add(collUSD)
	}

	targetCollUSD := targetCollRatio.MulInt(stableSupply.Amount)
	neededStable = targetCollUSD.Sub(currentTotalCollUSD)
	return neededStable, err
}

func (k *Keeper) RecollateralizeCollAmtForTargetCollRatio(
	ctx sdk.Context,
) (neededCollAmount sdk.Int, err error) {
	neededUSDForRecoll, _ := k.GetUSDValForTargetCollRatio(ctx)
	priceCollStable, err := k.PriceKeeper.GetCurrentPrice(
		ctx, common.CollDenom, common.StableDenom)
	if err != nil {
		return sdk.Int{}, err
	}

	neededCollAmountDec := neededUSDForRecoll.Quo(priceCollStable.Price)
	return neededCollAmountDec.Ceil().TruncateInt(), err
}

/*
Recollateralize is a function that incentivizes the caller to add up to the
amount of collateral needed to reach some target collateral ratio
(`collRatioTarget`). Recollateralize checks if the USD value of collateral in
the protocol is below the required amount defined by the current collateral ratio.
Nibiru's NUSD stablecoin is taken to be the dollar that determines USD value.

Args:
  msg (MsgRecollateralize) {
    Creator (string): Caller of 'Recollateralize'
	Coll (sdk.Coin): Input collateral that will be sold to the protocol.
  }

Returns:
  response (MsgRecollateralizeResponse) {
    Gov (sdk.Coin): Governance received as a reward for recollateralizing Nibiru.
  }
  err: Error condition for if the function succeeds or fails.
*/
func (k Keeper) Recollateralize(
	goCtx context.Context, msg *types.MsgRecollateralize,
) (response *types.MsgRecollateralizeResponse, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	caller, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return response, err
	}

	params := k.GetParams(ctx)
	targetCollRatio := params.GetCollRatioAsDec()

	neededCollAmt, err := k.RecollateralizeCollAmtForTargetCollRatio(ctx)
	if err != nil {
		return response, err
	} else if neededCollAmt.LTE(sdk.ZeroInt()) {
		return response, fmt.Errorf(
			"protocol has sufficient COLL, so 'Recollateralize' is not needed")
	}

	// The caller doesn't need to be put in the full amount,
	// just a positive amount that is at most the 'neededCollAmount'.
	inColl := sdk.NewCoin(msg.Coll.Denom, sdk.ZeroInt())
	if msg.Coll.Amount.GT(neededCollAmt) {
		inColl.Amount = neededCollAmt
	} else if msg.Coll.Amount.LTE(sdk.ZeroInt()) {
		return response, fmt.Errorf(
			"collateral input, %v, must be positive", msg.Coll.String())
	} else {
		inColl.Amount = msg.Coll.Amount
	}

	// Send collateral from the caller to the module
	err = k.checkEnoughBalance(ctx, inColl, caller)
	if err != nil {
		return response, err
	}
	err = k.BankKeeper.SendCoinsFromAccountToModule(
		ctx, caller, types.ModuleName, sdk.NewCoins(inColl),
	)
	if err != nil {
		return response, err
	}
	events.EmitTransfer(
		ctx,
		/* coin */ inColl,
		/* from */ caller.String(),
		/* to   */ k.AccountKeeper.GetModuleAddress(types.ModuleName).String(),
	)

	// Compute GOV rewarded to user
	priceCollStable, err := k.PriceKeeper.GetCurrentPrice(
		ctx, common.CollDenom, common.StableDenom)
	if err != nil {
		return response, err
	}
	inUSD := priceCollStable.Price.MulInt(inColl.Amount)
	outGovAmount, err := k.GovAmtFromRecollateralize(ctx, inUSD)
	if err != nil {
		return response, err
	}
	outGov := sdk.NewCoin(common.GovDenom, outGovAmount)

	// Mint and send GOV reward from the module to the caller
	err = k.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(outGov))
	if err != nil {
		return response, err
	}
	events.EmitMintNIBI(ctx, outGov)

	err = k.BankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, caller, sdk.NewCoins(outGov),
	)
	if err != nil {
		return response, err
	}
	events.EmitTransfer(
		ctx, outGov,
		/* from */ k.AccountKeeper.GetModuleAddress(types.ModuleName).String(),
		/* to   */ caller.String(),
	)

	events.EmitRecollateralize(
		ctx,
		/* inCoin    */ inColl,
		/* outCoin   */ outGov,
		/* caller    */ caller.String(),
		/* collRatio */ targetCollRatio,
	)
	return &types.MsgRecollateralizeResponse{
		Gov: outGov,
	}, err
}

/*
GovAmtFromRecollateralize computes the GOV token given as a reward for calling
recollateralize.

Args:
  ctx (sdk.Context): Carries information about the current state of the application.
  inUSD (sdk.Dec): Value in NUSD stablecoin to be used for recollateralization.
Returns:
  govOut (sdk.Int): Amount of GOV token rewarded for 'Recollateralize'.
*/
func (k *Keeper) GovAmtFromRecollateralize(
	ctx sdk.Context, inUSD sdk.Dec,
) (govOut sdk.Int, err error) {
	params := k.GetParams(ctx)
	bonusRate := params.GetBonusRateRecollAsDec()

	priceGovStable, err := k.PriceKeeper.GetCurrentPrice(
		ctx, common.GovDenom, common.StableDenom)
	if err != nil {
		return sdk.Int{}, err
	}
	govOut = inUSD.Mul(sdk.OneDec().Add(bonusRate)).
		Quo(priceGovStable.Price).TruncateInt()
	return govOut, err
}

func (k *Keeper) GovAmtFromFullRecollateralize(
	ctx sdk.Context,
) (govOut sdk.Int, err error) {
	neededCollUSD, err := k.GetUSDValForTargetCollRatio(ctx)
	if err != nil {
		return sdk.Int{}, err
	}
	return k.GovAmtFromRecollateralize(ctx, neededCollUSD)
}

// ---------------------------------------------------------------------------
// Buyback
// ---------------------------------------------------------------------------

func (k *Keeper) BuybackGovAmtForTargetCollRatio(
	ctx sdk.Context,
) (neededGovAmt sdk.Int, err error) {
	neededUSDForRecoll, _ := k.GetUSDValForTargetCollRatio(ctx)
	neededUSDForBuyback := neededUSDForRecoll.Neg()
	priceGovStable, err := k.PriceKeeper.GetCurrentPrice(
		ctx, common.GovDenom, common.StableDenom)
	if err != nil {
		return sdk.Int{}, err
	}

	neededGovAmtDec := neededUSDForBuyback.Quo(priceGovStable.Price)
	neededGovAmt = neededGovAmtDec.Ceil().TruncateInt()
	return neededGovAmt, err
}

func (k Keeper) Buyback(
	goCtx context.Context, msg *types.MsgBuyback,
) (response *types.MsgBuybackResponse, err error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	caller, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return response, err
	}

	params := k.GetParams(ctx)
	targetCollRatio := params.GetCollRatioAsDec()

	neededGovAmt, err := k.BuybackGovAmtForTargetCollRatio(ctx)
	if err != nil {
		return response, err
	} else if neededGovAmt.LTE(sdk.ZeroInt()) {
		return response, fmt.Errorf(
			"protocol has insufficient COLL, so 'Buyback' is not needed")
	}

	// The caller doesn't need to be put in the full amount,
	// just a positive amount that is at most the 'neededCollAmount'.
	inGov := sdk.NewCoin(msg.Gov.Denom, sdk.ZeroInt())
	if msg.Gov.Amount.GT(neededGovAmt) {
		inGov.Amount = neededGovAmt
	} else if msg.Gov.Amount.LTE(sdk.ZeroInt()) {
		return response, fmt.Errorf(
			"collateral input, %v, must be positive", msg.Gov.String())
	} else {
		inGov.Amount = msg.Gov.Amount
	}

	// Send NIBI from the caller to the module
	err = k.checkEnoughBalance(ctx, inGov, caller)
	if err != nil {
		return response, err
	}
	err = k.BankKeeper.SendCoinsFromAccountToModule(
		ctx, caller, types.ModuleName, sdk.NewCoins(inGov),
	)
	if err != nil {
		return response, err
	}
	events.EmitTransfer(
		ctx,
		/* coin */ inGov,
		/* from */ caller.String(),
		/* to   */ k.AccountKeeper.GetModuleAddress(types.ModuleName).String(),
	)

	// Burn the NIBI that was sent by the caller.
	err = k.BankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(inGov))
	if err != nil {
		return response, err
	}
	events.EmitBurnNIBI(ctx, inGov)

	// Compute USD (stable) value of the GOV sent by the caller: 'inUSD'
	priceGovStable, err := k.PriceKeeper.GetCurrentPrice(
		ctx, common.GovDenom, common.StableDenom)
	if err != nil {
		return response, err
	}
	inUSD := priceGovStable.Price.MulInt(inGov.Amount)

	// Compute collateral amount sent to caller: 'outColl'
	outCollAmount, err := k.CollAmtFromBuyback(ctx, inUSD)
	if err != nil {
		return response, err
	}
	outColl := sdk.NewCoin(common.CollDenom, outCollAmount)

	// Send COLL from the module to the caller
	err = k.BankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, caller, sdk.NewCoins(outColl),
	)
	if err != nil {
		return response, err
	}
	events.EmitTransfer(
		ctx, outColl,
		/* from */ k.AccountKeeper.GetModuleAddress(types.ModuleName).String(),
		/* to   */ caller.String(),
	)

	events.EmitBuyback(
		ctx,
		/* inCoin    */ inGov,
		/* outCoin   */ outColl,
		/* caller    */ caller.String(),
		/* collRatio */ targetCollRatio,
	)
	return &types.MsgBuybackResponse{
		Coll: outColl,
	}, err
}

/*
CollAmtFromBuyback computes the COLL (collateral) given as a reward for calling
buyback.

Args:
  ctx (sdk.Context): Carries information about the current state of the application.
  valUSD (sdk.Dec): Value in NUSD stablecoin to be used for buyback.
Returns:
  collAmt (sdk.Int): Amount of COLL token rewarded for 'Buyback'.
*/
func (k *Keeper) CollAmtFromBuyback(
	ctx sdk.Context, valUSD sdk.Dec,
) (collAmt sdk.Int, err error) {

	priceCollStable, err := k.PriceKeeper.GetCurrentPrice(
		ctx, common.CollDenom, common.StableDenom)
	if err != nil {
		return sdk.Int{}, err
	}
	collAmt = valUSD.
		Quo(priceCollStable.Price).TruncateInt()
	return collAmt, err
}

// TODO hygiene: cover with test cases
func (k *Keeper) CollAmtFromFullBuyback(
	ctx sdk.Context,
) (collAmt sdk.Int, err error) {

	neededUSDForRecoll, err := k.GetUSDValForTargetCollRatio(ctx)
	if err != nil {
		return sdk.Int{}, err
	}
	neededUSDForBuyback := neededUSDForRecoll.Neg()
	return k.CollAmtFromBuyback(ctx, neededUSDForBuyback)
}
