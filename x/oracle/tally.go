package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/NibiruChain/nibiru/x/oracle/keeper"
	"github.com/NibiruChain/nibiru/x/oracle/types"
)

// Tally calculates the median and returns it. Sets the set of voters to be rewarded, i.e. voted within
// a reasonable spread from the weighted median to the store
// CONTRACT: pb must be sorted
func Tally(_ sdk.Context, pb types.ExchangeRateBallot, rewardBand sdk.Dec, validatorClaimMap map[string]types.ValidatorPerformance) (weightedMedian sdk.Dec) {
	weightedMedian = pb.WeightedMedianWithAssertion()

	standardDeviation := pb.StandardDeviation(weightedMedian)
	rewardSpread := weightedMedian.Mul(rewardBand.QuoInt64(2))

	if standardDeviation.GT(rewardSpread) {
		rewardSpread = standardDeviation
	}

	for _, vote := range pb {
		// Filter ballot winners & abstain voters
		if (vote.ExchangeRate.GTE(weightedMedian.Sub(rewardSpread)) &&
			vote.ExchangeRate.LTE(weightedMedian.Add(rewardSpread))) ||
			!vote.ExchangeRate.IsPositive() {
			key := vote.Voter.String()
			claim := validatorClaimMap[key]
			claim.Weight += vote.Power
			claim.WinCount++
			validatorClaimMap[key] = claim
		}
	}

	return
}

// ballot for the asset is passing the threshold amount of voting power
func ballotIsPassing(ballot types.ExchangeRateBallot, thresholdVotes sdk.Int) (sdk.Int, bool) {
	ballotPower := sdk.NewInt(ballot.Power())
	return ballotPower, !ballotPower.IsZero() && ballotPower.GTE(thresholdVotes)
}

// PickReferencePair choose Reference pair with the highest voter turnout
// If the voting power of the two denominations is the same,
// select reference pair in alphabetical order.
func PickReferencePair(ctx sdk.Context, k keeper.Keeper, voteTargets map[string]struct{}, voteMap map[string]types.ExchangeRateBallot) string {
	largestBallotPower := int64(0)
	referencePair := ""

	totalBondedPower := sdk.TokensToConsensusPower(k.StakingKeeper.TotalBondedTokens(ctx), k.StakingKeeper.PowerReduction(ctx))
	voteThreshold := k.VoteThreshold(ctx)
	thresholdVotes := voteThreshold.MulInt64(totalBondedPower).RoundInt()

	for pair, ballot := range voteMap {
		// If pair is not in the voteTargets, or the ballot for it has failed, then skip
		// and remove it from voteMap for iteration efficiency
		if _, exists := voteTargets[pair]; !exists {
			delete(voteMap, pair)
			continue
		}

		ballotPower := int64(0)

		// If the ballot is not passed, remove it from the voteTargets array
		// to prevent slashing validators who did valid vote.
		if power, ok := ballotIsPassing(ballot, thresholdVotes); ok {
			ballotPower = power.Int64()
		} else {
			delete(voteTargets, pair)
			delete(voteMap, pair)
			continue
		}

		if ballotPower > largestBallotPower || largestBallotPower == 0 {
			referencePair = pair
			largestBallotPower = ballotPower
		} else if largestBallotPower == ballotPower && referencePair > pair {
			referencePair = pair
		}
	}

	return referencePair
}
