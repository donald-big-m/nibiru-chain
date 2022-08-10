package keeper

import (
	"sort"

	"github.com/NibiruChain/nibiru/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// OrganizeBallotByPair collects all oracle votes for the period, categorized by the votes' pair parameter
func (k Keeper) OrganizeBallotByPair(ctx sdk.Context, validatorClaimMap map[string]types.Claim) (ballots map[string]types.ExchangeRateBallot) {
	ballots = map[string]types.ExchangeRateBallot{}

	// Organize aggregate votes
	aggregateHandler := func(voterAddr sdk.ValAddress, vote types.AggregateExchangeRateVote) (stop bool) {
		// organize ballot only for the active validators
		if claim, ok := validatorClaimMap[vote.Voter]; ok {
			for _, tuple := range vote.ExchangeRateTuples {
				power := claim.Power
				if !tuple.ExchangeRate.IsPositive() {
					// Make the power of abstain vote zero
					power = 0
				}

				ballots[tuple.Pair] = append(ballots[tuple.Pair],
					types.NewBallotVoteForTally(
						tuple.ExchangeRate,
						tuple.Pair,
						voterAddr,
						power,
					),
				)
			}
		}

		return false
	}

	k.IterateAggregateExchangeRateVotes(ctx, aggregateHandler)

	// sort created ballot
	for pair, ballot := range ballots {
		sort.Sort(ballot)
		ballots[pair] = ballot
	}

	return
}

// ClearBallots clears all tallied prevotes and votes from the store
func (k Keeper) ClearBallots(ctx sdk.Context, votePeriod uint64) {
	// Clear all aggregate prevotes
	k.IterateAggregateExchangeRatePrevotes(ctx, func(voterAddr sdk.ValAddress, aggregatePrevote types.AggregateExchangeRatePrevote) (stop bool) {
		if ctx.BlockHeight() > int64(aggregatePrevote.SubmitBlock+votePeriod) {
			k.DeleteAggregateExchangeRatePrevote(ctx, voterAddr)
		}

		return false
	})

	// Clear all aggregate votes
	k.IterateAggregateExchangeRateVotes(ctx, func(voterAddr sdk.ValAddress, aggregateVote types.AggregateExchangeRateVote) (stop bool) {
		k.DeleteAggregateExchangeRateVote(ctx, voterAddr)
		return false
	})
}

// ApplyWhitelist update vote target pair list and set tobin tax with params whitelist
func (k Keeper) ApplyWhitelist(ctx sdk.Context, whitelist types.PairList, voteTargets map[string]struct{}) {

	// check is there any update in whitelist params
	if len(voteTargets) != len(whitelist) {
		k.ClearPairs(ctx)

		for _, item := range whitelist {
			k.SetPair(ctx, item.Name)
		}
	}
}
