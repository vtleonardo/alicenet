package dynamics

import (
	"context"
	"github.com/alicenet/alicenet/bridge/bindings"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"

	"github.com/alicenet/alicenet/layer1/executor/tasks"
)

// CanonicalVersionCheckTask contains required state for the task.
type CanonicalVersionCheckTask struct {
	*tasks.BaseTask
	//Version info
	MaxUpdateEpoch *big.Int
	Version        bindings.CanonicalVersion
}

// asserting that CanonicalVersionCheckTask struct implements interface tasks.Task.
var _ tasks.Task = &CanonicalVersionCheckTask{}

// NewVersionCheckTask creates a background task that attempts to verify the version check.
func NewVersionCheckTask(maxUpdateEpoch *big.Int, version bindings.CanonicalVersion) *CanonicalVersionCheckTask {
	return &CanonicalVersionCheckTask{
		BaseTask:       tasks.NewBaseTask(0, 0, false, nil),
		MaxUpdateEpoch: maxUpdateEpoch,
		Version:        version,
	}
}

// Prepare prepares for work to be done in the CanonicalVersionCheckTask.
func (t *CanonicalVersionCheckTask) Prepare(ctx context.Context) *tasks.TaskErr {
	logger := t.GetLogger().WithField("method", "Prepare()")
	logger.Debug("preparing task")

	return nil
}

// Execute executes the task business logic.
func (t *CanonicalVersionCheckTask) Execute(ctx context.Context) (*types.Transaction, *tasks.TaskErr) {
	logger := t.GetLogger().WithField("method", "Execute()")
	logger.Debug("initiate execution")

	//todo: get local version, log a message depending on the patch severity

	return nil, nil
}

// ShouldExecute checks if it makes sense to execute the task.
func (t *CanonicalVersionCheckTask) ShouldExecute(ctx context.Context) (bool, *tasks.TaskErr) {
	logger := t.GetLogger().WithField("method", "ShouldExecute()")
	logger.Debug("should execute task")

	//todo: get local version, compare to the event version and respond

	return true, nil
}
