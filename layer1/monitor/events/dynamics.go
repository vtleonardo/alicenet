package events

import (
	"fmt"

	"github.com/MadBase/MadNet/consensus/db"
	"github.com/MadBase/MadNet/layer1"
	"github.com/MadBase/MadNet/layer1/ethereum"
	"github.com/MadBase/MadNet/layer1/executor/tasks/dkg/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

// ProcessValueUpdated handles a dynamic value updating coming from our smart contract
func ProcessValueUpdated(eth layer1.Client, logger *logrus.Entry, log types.Log,
	monDB *db.Database) error {

	logger.Info("ProcessValueUpdated() ...")

	dkgState, err := state.GetDkgState(monDB)
	if err != nil {
		return err
	}

	// todo: ask Hunter only validators allowed?
	if !dkgState.IsValidator {
		return nil
	}

	event, err := ethereum.GetContracts().Governance().ParseValueUpdated(log)
	if err != nil {
		return err
	}

	logger = logger.WithFields(logrus.Fields{
		"Epoch": event.Epoch.Uint64(),
		"Key":   event.Key.String(),
		"Value": fmt.Sprintf("0x%x", event.Value),
	})

	logger.Infof("Value updated")

	logger.Warnf("Dropping dynamic value on the floor")
	return nil
}