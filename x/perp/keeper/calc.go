package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/NibiruChain/nibiru/x/common"
	"github.com/NibiruChain/nibiru/x/perp/types"
)

// NOTE hardcoded for now. Need to discuss whether this should be part of the
// Params of x/perp
var initMarginRatio = sdk.MustNewDecFromStr("0.1")

type RemainingMarginWithFundingPayment struct {
	// Margin: amount of quote token (y) backing the position.
	Margin sdk.Dec

	/* BadDebt: Bad debt (margin units) cleared by the PerpEF during the tx.
	   Bad debt is negative net margin past the liquidation point of a position. */
	BadDebt sdk.Dec

	/* FundingPayment: A funding payment (margin units) made or received by the trader on
	    the current position. 'fundingPayment' is positive if 'owner' is the sender
		and negative if 'owner' is the receiver of the payment. Its magnitude is
		abs(vSize * fundingRate). Funding payments act to converge the mark price
		(vPrice) and index price (average price on major exchanges).
	*/
	FundingPayment sdk.Dec

	/* LatestCumulativePremiumFraction: latest cumulative premium fraction. Units are (margin units)/position size. */
	LatestCumulativePremiumFraction sdk.Dec
}

func (k Keeper) CalcRemainMarginWithFundingPayment(
	ctx sdk.Context,
	currentPosition types.Position,
	marginDelta sdk.Dec,
) (remaining RemainingMarginWithFundingPayment, err error) {
	remaining.LatestCumulativePremiumFraction, err = k.
		getLatestCumulativePremiumFraction(ctx, common.TokenPair(currentPosition.Pair))
	if err != nil {
		return remaining, err
	}

	if currentPosition.Size_.IsZero() {
		remaining.FundingPayment = sdk.ZeroDec()
	} else {
		remaining.FundingPayment = remaining.LatestCumulativePremiumFraction.
			Sub(currentPosition.LastUpdateCumulativePremiumFraction).
			Mul(currentPosition.Size_)
	}

	remainingMargin := currentPosition.Margin.Add(marginDelta).Sub(remaining.FundingPayment)

	if remainingMargin.IsNegative() {
		// the remaining margin is negative, liquidators didn't do their job
		// and we have negative margin that must come out of the ecosystem fund
		remaining.BadDebt = remainingMargin.Abs()
		remaining.Margin = sdk.ZeroDec()
	} else {
		remaining.Margin = remainingMargin.Abs()
		remaining.BadDebt = sdk.ZeroDec()
	}

	return remaining, err
}

func (k Keeper) calcFreeCollateral(ctx sdk.Context, pos types.Position, fundingPayment sdk.Dec,
) (sdk.Int, error) {
	pair, err := common.NewTokenPairFromStr(pos.Pair)
	if err != nil {
		return sdk.Int{}, err
	}
	err = k.requireVpool(ctx, pair)
	if err != nil {
		return sdk.Int{}, err
	}

	unrealizedPnL, positionNotional, err := k.
		getPreferencePositionNotionalAndUnrealizedPnL(
			ctx,
			pos,
			types.PnLPreferenceOption_MIN,
		)
	if err != nil {
		return sdk.Int{}, err
	}
	freeMargin := pos.Margin.Add(fundingPayment)
	accountValue := unrealizedPnL.Add(freeMargin)
	minCollateral := sdk.MinDec(accountValue, freeMargin)

	// Get margin requirement. This rounds up, so 16.5 margin required -> 17
	var marginRequirement sdk.Int
	if pos.Size_.IsPositive() {
		// if long position, use open notional
		marginRequirement = initMarginRatio.Mul(pos.OpenNotional).RoundInt()
	} else {
		// if short, use current notional
		marginRequirement = initMarginRatio.Mul(positionNotional).RoundInt()
	}
	freeCollateral := minCollateral.Sub(marginRequirement.ToDec()).TruncateInt()
	return freeCollateral, nil
}
