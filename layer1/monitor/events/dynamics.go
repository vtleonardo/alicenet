package events

import (
	"fmt"
	"github.com/alicenet/alicenet/layer1/executor/tasks"
	"github.com/alicenet/alicenet/layer1/executor/tasks/dynamics"
	"github.com/alicenet/alicenet/layer1/monitor/objects"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"

	"github.com/alicenet/alicenet/consensus/db"
	"github.com/alicenet/alicenet/layer1"
)

// ProcessValueUpdated handles a dynamic value updating coming from our smart contract.
func ProcessValueUpdated(eth layer1.Client, contracts layer1.AllSmartContracts, logger *logrus.Entry, log types.Log, monDB *db.Database) error {
	logger.Info("ProcessValueUpdated() ...")

	event, err := contracts.EthereumContracts().Governance().ParseValueUpdated(log)
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

func ProcessNewAliceNetNodeVersionAvailable(eth layer1.Client, contracts layer1.AllSmartContracts, logger *logrus.Entry, log types.Log, monState *objects.MonitorState, taskRequestChan chan<- tasks.TaskRequest, monDB *db.Database) error {
	logger = logger.WithField("method", "ProcessNewAliceNetNodeVersionAvailable")
	logger.Info("Processing new AliceNet node version...")

	event, err := contracts.EthereumContracts().Dynamics().ParseNewAliceNetNodeVersionAvailable(log)
	if err != nil {
		return err
	}

	logger = logger.WithFields(logrus.Fields{
		"MaxUpdateEpoch": event.MaxUpdateEpoch.Int64(),
		"Major":          event.Version.Major,
		"Minor":          event.Version.Minor,
		"Patch":          event.Version.Patch,
	})

	monState.CanonicalVersion = event.Version
	logger.Infof("New AliceNet node version received and updated")

	// Killing previous task
	taskRequestChan <- tasks.NewKillTaskRequest(&dynamics.CanonicalVersionCheckTask{})

	//todo: check the local version and schedule the task if applies

	// Scheduling task with the new Canonical Version
	taskRequestChan <- tasks.NewScheduleTaskRequest(dynamics.NewVersionCheckTask(event.MaxUpdateEpoch, event.Version))

	return nil
}
