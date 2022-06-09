package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/NibiruChain/nibiru/x/common"
	"github.com/NibiruChain/nibiru/x/perp/events"
	"github.com/NibiruChain/nibiru/x/perp/types"
)

/* AddMargin deleverages an existing position by adding margin (collateral)
to it. Adding margin increases the margin ratio of the corresponding position.
*/
func (k Keeper) AddMargin(
	goCtx context.Context, msg *types.MsgAddMargin,
) (res *types.MsgAddMarginResponse, err error) {
	// ------------- Message Setup -------------
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate trader
	msgSender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// validate margin amount
	addedMargin := msg.Margin.Amount
	if !addedMargin.IsPositive() {
		err = fmt.Errorf("margin must be positive, not: %v", addedMargin.String())
		k.Logger(ctx).Debug(
			err.Error(),
			"margin_amount",
			msg.Margin.Amount.String(),
		)
		return nil, err
	}

	// validate token pair
	pair, err := common.NewAssetPairFromStr(msg.TokenPair)
	if err != nil {
		k.Logger(ctx).Debug(
			err.Error(),
			"token_pair",
			msg.TokenPair,
		)
		return nil, err
	}
	// validate vpool exists
	if err = k.requireVpool(ctx, pair); err != nil {
		return nil, err
	}

	// validate margin denom
	if msg.Margin.Denom != pair.GetQuoteTokenDenom() {
		err = fmt.Errorf("invalid margin denom")
		k.Logger(ctx).Debug(
			err.Error(),
			"margin_denom",
			msg.Margin.Denom,
			"quote_token_denom",
			pair.GetQuoteTokenDenom(),
		)
		return nil, err
	}

	// ------------- AddMargin -------------
	position, err := k.Positions().Get(ctx, pair, msgSender)
	if err != nil {
		k.Logger(ctx).Debug(
			err.Error(),
			"pair",
			pair.String(),
			"trader",
			msg.Sender,
		)
		return nil, err
	}

	position.Margin = position.Margin.Add(addedMargin.ToDec())
	coinToSend := sdk.NewCoin(pair.GetQuoteTokenDenom(), addedMargin)
	if err = k.BankKeeper.SendCoinsFromAccountToModule(
		ctx, msgSender, types.VaultModuleAccount, sdk.NewCoins(coinToSend),
	); err != nil {
		k.Logger(ctx).Debug(
			err.Error(),
			"trader",
			msg.Sender,
			"coin",
			coinToSend.String(),
		)
		return nil, err
	}

	events.EmitTransfer(ctx,
		/* coin */ coinToSend,
		/* from */ k.AccountKeeper.GetModuleAddress(types.VaultModuleAccount),
		/* to */ msgSender,
	)

	k.Positions().Set(ctx, pair, msgSender, position)

	// TODO(https://github.com/NibiruChain/nibiru/issues/323): calculate the funding payment
	fPayment := sdk.ZeroDec()
	events.EmitMarginChange(ctx, msgSender, pair.String(), addedMargin, fPayment)

	fmt.Println("STEVENDEBUG add margin done")
	return &types.MsgAddMarginResponse{}, nil
}

/* RemoveMargin further leverages an existing position by directly removing
the margin (collateral) that backs it from the vault. This also decreases the
margin ratio of the position.
*/
// STEVENDEBUG remove margin
func (k Keeper) RemoveMargin(
	goCtx context.Context, msg *types.MsgRemoveMargin,
) (res *types.MsgRemoveMarginResponse, err error) {
	// ------------- Message Setup -------------
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate trader
	msgSender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	fmt.Println("STEVENDEBUG keeper remove margin - a")

	// validate margin amount
	if !msg.Margin.Amount.IsPositive() {
		err = fmt.Errorf("margin must be positive, not: %v", msg.Margin.Amount.String())
		k.Logger(ctx).Debug(
			err.Error(),
			"margin_amount",
			msg.Margin.Amount.String(),
		)
		return nil, err
	}

	fmt.Println("STEVENDEBUG keeper remove margin - b")

	// validate token pair
	pair, err := common.NewAssetPairFromStr(msg.TokenPair)
	if err != nil {
		k.Logger(ctx).Debug(
			err.Error(),
			"token_pair",
			msg.TokenPair,
		)
		return nil, err
	}

	fmt.Println("STEVENDEBUG keeper remove margin - c")

	// validate vpool exists
	if err = k.requireVpool(ctx, pair); err != nil {
		return nil, err
	}

	fmt.Println("STEVENDEBUG keeper remove margin - d")

	// validate margin denom
	if msg.Margin.Denom != pair.GetQuoteTokenDenom() {
		err = fmt.Errorf("invalid margin denom")
		k.Logger(ctx).Debug(
			err.Error(),
			"margin_denom",
			msg.Margin.Denom,
			"quote_token_denom",
			pair.GetQuoteTokenDenom(),
		)
		return nil, err
	}

	fmt.Println("STEVENDEBUG keeper remove margin - e")

	// ------------- RemoveMargin -------------
	position, err := k.Positions().Get(ctx, pair, msgSender)

	fmt.Println("STEVENDEBUG keeper remove margin - e - 1 err = ", err)

	if err != nil {
		k.Logger(ctx).Debug(
			err.Error(),
			"pair",
			pair.String(),
			"trader",
			msg.Sender,
		)
		return nil, err
	}

	fmt.Println("STEVENDEBUG keeper remove margin - f")

	marginDelta := msg.Margin.Amount.Neg()
	remaining, err := k.CalcRemainMarginWithFundingPayment(
		ctx, *position, marginDelta.ToDec())

	fmt.Println("STEVENDEBUG keeper remove margin - f err - 1 ", err)

	fmt.Println("STEVENDEBUG keeper remove margin - f err - bad debt ", remaining.BadDebt)

	if err != nil {
		return nil, err
	}
	if !remaining.BadDebt.IsZero() {
		err = types.ErrFailedToRemoveDueToBadDebt
		k.Logger(ctx).Debug(
			err.Error(),
			"remaining_bad_debt",
			remaining.BadDebt.String(),
		)

		fmt.Println("STEVENDEBUG keeper remove margin - f err - 2 ", err)

		newsdkerr := sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())

		fmt.Println("STEVENDEBUG keeper remove margin - f newsdkerr - 2 ", newsdkerr)

		return nil, err
	}

	fmt.Println("STEVENDEBUG keeper remove margin - g")

	position.Margin = remaining.Margin
	position.LastUpdateCumulativePremiumFraction = remaining.LatestCumulativePremiumFraction
	freeCollateral, err := k.calcFreeCollateral(
		ctx, *position, remaining.FundingPayment)
	if err != nil {
		return res, err
	} else if !freeCollateral.IsPositive() {
		return res, fmt.Errorf("not enough free collateral")
	}

	fmt.Println("STEVENDEBUG keeper remove margin - h")

	k.Positions().Set(ctx, pair, msgSender, position)

	coinToSend := sdk.NewCoin(pair.GetQuoteTokenDenom(), msg.Margin.Amount)
	err = k.BankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.VaultModuleAccount, msgSender, sdk.NewCoins(coinToSend))
	if err != nil {
		k.Logger(ctx).Debug(
			err.Error(),
			"to",
			msg.Sender,
			"coin",
			coinToSend.String(),
		)
		return nil, err
	}

	fmt.Println("STEVENDEBUG keeper remove margin - i")

	events.EmitTransfer(ctx,
		/* coin */ coinToSend,
		/* from */ k.AccountKeeper.GetModuleAddress(types.VaultModuleAccount),
		/* to */ msgSender,
	)

	fmt.Println("STEVENDEBUG keeper remove margin - j")

	events.EmitMarginChange(
		ctx,
		msgSender,
		pair.String(),
		msg.Margin.Amount,
		remaining.FundingPayment,
	)

	fmt.Println("STEVENDEBUG keeper remove margin - k")

	return &types.MsgRemoveMarginResponse{
		MarginOut:      coinToSend,
		FundingPayment: remaining.FundingPayment,
	}, nil
}

// GetMarginRatio calculates the MarginRatio from a Position
func (k Keeper) GetMarginRatio(
	ctx sdk.Context, position types.Position, priceOption types.MarginCalculationPriceOption,
) (marginRatio sdk.Dec, err error) {
	fmt.Println("STEVENDEBUG GetMarginRatio start")

	if position.Size_.IsZero() {
		return sdk.Dec{}, types.ErrPositionZero
	}

	var (
		unrealizedPnL    sdk.Dec
		positionNotional sdk.Dec
	)

	switch priceOption {
	case types.MarginCalculationPriceOption_MAX_PNL:
		positionNotional, unrealizedPnL, err = k.getPreferencePositionNotionalAndUnrealizedPnL(
			ctx,
			position,
			types.PnLPreferenceOption_MAX,
		)
	case types.MarginCalculationPriceOption_INDEX:
		positionNotional, unrealizedPnL, err = k.getPositionNotionalAndUnrealizedPnL(
			ctx,
			position,
			types.PnLCalcOption_ORACLE,
		)
	case types.MarginCalculationPriceOption_SPOT:
		positionNotional, unrealizedPnL, err = k.getPositionNotionalAndUnrealizedPnL(
			ctx,
			position,
			types.PnLCalcOption_SPOT_PRICE,
		)
	}

	if err != nil {
		return sdk.Dec{}, err
	}
	if positionNotional.IsZero() {
		// NOTE causes division by zero in margin ratio calculation
		return sdk.Dec{},
			fmt.Errorf("margin ratio doesn't make sense with zero position notional")
	}

	remaining, err := k.CalcRemainMarginWithFundingPayment(
		ctx,
		/* oldPosition */ position,
		/* marginDelta */ unrealizedPnL,
	)
	if err != nil {
		return sdk.Dec{}, err
	}

	marginRatio = remaining.Margin.Sub(remaining.BadDebt).
		Quo(positionNotional)
	return marginRatio, nil
}

func (k Keeper) requireVpool(ctx sdk.Context, pair common.AssetPair) (err error) {
	if !k.VpoolKeeper.ExistsPool(ctx, pair) {
		err = fmt.Errorf("%v: %v", types.ErrPairNotFound.Error(), pair.String())
		k.Logger(ctx).Error(
			err.Error(),
			"pair",
			pair.String(),
		)
		return err
	}
	return nil
}

/*
requireMoreMarginRatio checks if the marginRatio corresponding to the margin
backing a position is above or below the 'baseMarginRatio'.
If 'largerThanOrEqualTo' is true, 'marginRatio' must be >= 'baseMarginRatio'.

Args:
  marginRatio: Ratio of the value of the margin and corresponding position(s).
    marginRatio is defined as (margin + unrealizedPnL) / notional
  baseMarginRatio: Specifies the threshold value that 'marginRatio' must meet.
  largerThanOrEqualTo: Specifies whether 'marginRatio' should be larger or
    smaller than 'baseMarginRatio'.
*/
func requireMoreMarginRatio(marginRatio, baseMarginRatio sdk.Dec, largerThanOrEqualTo bool) error {
	if largerThanOrEqualTo {
		if !marginRatio.GTE(baseMarginRatio) {
			return fmt.Errorf("margin ratio did not meet criteria")
		}
	} else {
		if !marginRatio.LT(baseMarginRatio) {
			return fmt.Errorf("margin ratio did not meet criteria")
		}
	}
	return nil
}
