package utils

import (
	"context"
	"fmt"
	"math/big"

	"github.com/alicenet/alicenet/layer1"
	"github.com/sirupsen/logrus"
)

// AmILeading checks if the current node is a leader for an action
func AmILeading(client layer1.Client, ctx context.Context, logger *logrus.Entry, start int, randHash []byte, numOfValidators int, validatorIndex int, desperationFactor int, desperationDelay int) (bool, error) {
	currentHeight, err := client.GetCurrentHeight(ctx)
	if err != nil {
		return false, err
	}

	blocksSinceDesperation := int(currentHeight) - start - desperationDelay
	amILeading := LeaderElection(numOfValidators, validatorIndex, blocksSinceDesperation, desperationFactor, randHash)

	logger.WithFields(logrus.Fields{
		"currentHeight":          currentHeight,
		"startBlock":             start,
		"desperationDelay":       desperationDelay,
		"blocksSinceDesperation": blocksSinceDesperation,
		"amILeading":             amILeading,
		"randomHash":             fmt.Sprintf("0x%x", randHash),
	}).Info("Checking if I'm leading this action")

	return amILeading, nil
}

// LeaderElection runs the leader election algorithm to check if an index is a leader or not.
func LeaderElection(numValidators int, myIdx int, blocksSinceDesperation int, desperationFactor int, seedHash []byte) bool {
	var numValidatorsAllowed int = 1
	for i := int(blocksSinceDesperation); i > 0; {
		i -= desperationFactor / numValidatorsAllowed
		numValidatorsAllowed++

		if numValidatorsAllowed >= numValidators {
			break
		}
	}

	// use the random nature of seedHash to deterministically define the range of
	// validators that are allowed to take an action
	rand := (&big.Int{}).SetBytes(seedHash)
	start := int((&big.Int{}).Mod(rand, big.NewInt(int64(numValidators))).Int64())
	end := (start + numValidatorsAllowed) % numValidators

	if end > start {
		return myIdx >= start && myIdx < end
	} else {
		return myIdx >= start || myIdx < end
	}
}
