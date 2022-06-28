package executor

import (
	"context"
	"github.com/alicenet/alicenet/consensus/db"
	"github.com/alicenet/alicenet/layer1/executor/tasks"
	"github.com/alicenet/alicenet/logging"
	"github.com/alicenet/alicenet/test/mocks"
	mockrequire "github.com/derision-test/go-mockgen/testutil/require"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func getTaskManager(t *testing.T) (*TasksManager, *mocks.MockClient, *db.Database, *taskResponseChan, *mocks.MockWatcher) {
	db := mocks.NewTestDB()
	client := mocks.NewMockClient()
	client.ExtractTransactionSenderFunc.SetDefaultReturn(common.Address{}, nil)
	client.GetTxMaxStaleBlocksFunc.SetDefaultReturn(10)
	hdr := &types.Header{
		Number: big.NewInt(1),
	}
	client.GetHeaderByNumberFunc.SetDefaultReturn(hdr, nil)
	client.GetBlockBaseFeeAndSuggestedGasTipFunc.SetDefaultReturn(big.NewInt(100), big.NewInt(1), nil)
	client.GetDefaultAccountFunc.SetDefaultReturn(accounts.Account{Address: common.Address{}})
	client.GetTransactionByHashFunc.SetDefaultReturn(nil, false, nil)
	client.GetFinalityDelayFunc.SetDefaultReturn(10)

	logger := logging.GetLogger("test")

	txWatcher := mocks.NewMockWatcher()
	taskManager, err := NewTaskManager(txWatcher, db, logger.WithField("Component", "schedule"))
	assert.Nil(t, err)

	taskRespChan := &taskResponseChan{trChan: make(chan tasks.TaskResponse, 100)}
	return taskManager, client, db, taskRespChan, txWatcher
}

func Test_TaskManagerHappyPath(t *testing.T) {
	manager, client, db, taskRespChan, txWatcher := getTaskManager(t)
	defer taskRespChan.close()

	receipt := &types.Receipt{
		Status:      types.ReceiptStatusSuccessful,
		BlockNumber: big.NewInt(20),
	}

	receiptResponse := mocks.NewMockReceiptResponse()
	receiptResponse.IsReadyFunc.SetDefaultReturn(true)
	receiptResponse.GetReceiptBlockingFunc.SetDefaultReturn(receipt, nil)

	txWatcher.SubscribeFunc.SetDefaultReturn(receiptResponse, nil)

	client.GetTransactionReceiptFunc.SetDefaultReturn(receipt, nil)

	task := mocks.NewMockTask()
	task.PrepareFunc.SetDefaultReturn(nil)
	txn := types.NewTx(&types.LegacyTx{
		Nonce:    1,
		Value:    big.NewInt(1),
		Gas:      1,
		GasPrice: big.NewInt(1),
		Data:     []byte{52, 66, 175, 92},
	})
	task.ExecuteFunc.SetDefaultReturn(txn, nil)
	task.ShouldExecuteFunc.SetDefaultReturn(true, nil)
	task.GetLoggerFunc.SetDefaultReturn(manager.logger)

	mainCtx := context.Background()
	manager.ManageTask(mainCtx, task, "happyPath", "123", db, manager.logger, client, taskRespChan)

	mockrequire.CalledOnce(t, task.Prepare)
	mockrequire.CalledOnce(t, task.Execute)
	mockrequire.CalledOnceWith(t, task.Finish, mockrequire.Values(nil))
}
