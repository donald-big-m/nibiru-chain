package epochs

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/NibiruChain/nibiru/x/epochs/keeper"
	"github.com/NibiruChain/nibiru/x/epochs/types"
)

// BeginBlocker of epochs module.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	k.IterateEpochInfo(ctx, func(index int64, epochInfo types.EpochInfo) (stop bool) {
		// If blocktimme < initial epoch start time, return
		if ctx.BlockTime().Before(epochInfo.StartTime) {
			return false
		}

		logger := k.Logger(ctx)

		// Epoch has not started yet
		shouldInitialEpochStart := !epochInfo.EpochCountingStarted

		epochEndTime := epochInfo.CurrentEpochStartTime.Add(epochInfo.Duration)
		// StartTime is set to a pinpointed timestamp that is after the default value of CurrentEpochStartTime (=0).
		shouldEpochStart := shouldInitialEpochStart || ctx.BlockTime().After(epochEndTime)

		if !shouldEpochStart {
			return false
		}

		// we deduced that a new epoch tick should happen
		epochInfo.CurrentEpochStartHeight = ctx.BlockHeight()

		if shouldInitialEpochStart {
			epochInfo.EpochCountingStarted = true
			epochInfo.CurrentEpoch = 1
			epochInfo.CurrentEpochStartTime = epochInfo.StartTime
			logger.Info(fmt.Sprintf("Starting new epoch with identifier %s epoch number %d", epochInfo.Identifier, epochInfo.CurrentEpoch))
		} else {
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeEpochEnd,
					sdk.NewAttribute(types.AttributeEpochNumber, fmt.Sprintf("%d", epochInfo.CurrentEpoch)),
				),
			)
			k.AfterEpochEnd(ctx, epochInfo.Identifier, epochInfo.CurrentEpoch)
			epochInfo.CurrentEpoch += 1
			epochInfo.CurrentEpochStartTime = epochInfo.CurrentEpochStartTime.Add(epochInfo.Duration)
			logger.Info(fmt.Sprintf("Starting epoch with identifier %s epoch number %d", epochInfo.Identifier, epochInfo.CurrentEpoch))
		}

		// emit new epoch start event, set epoch info, and run BeforeEpochStart hook
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeEpochStart,
				sdk.NewAttribute(types.AttributeEpochNumber, fmt.Sprintf("%d", epochInfo.CurrentEpoch)),
				sdk.NewAttribute(types.AttributeEpochStartTime, fmt.Sprintf("%d", epochInfo.CurrentEpochStartTime.Unix())),
			),
		)
		k.SetEpochInfo(ctx, epochInfo)
		k.BeforeEpochStart(ctx, epochInfo.Identifier, epochInfo.CurrentEpoch)

		return false
	})
}
