//go:build integration_comp

package dkg_test

import (
	"context"
	"github.com/alicenet/alicenet/layer1/ethereum"
	"github.com/alicenet/alicenet/layer1/executor/tasks/dkg/tests/cmd"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"log"
	"path/filepath"
	"testing"
)

func setupTest(t *testing.T, n int, cleanStart bool) {
	workingDir := cmd.CreateTestWorkingFolder()
	err := cmd.RunSetup(workingDir)
	if err != nil {
		log.Fatalf("Error setup %v", err)
	}

	validators, err := cmd.RunInit(workingDir, n)
	if err != nil {
		log.Fatalf("Error init %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	isRunning, _ := cmd.IsHardHatRunning()
	var details *ethereum.Client
	if !isRunning {
		log.Printf("Hardhat is not running. Start new HardHat")
		details = startHardHat(t, ctx, workingDir, validators[0])
		assert.NotNilf(t, details, "Expected details to be not nil")
	}

	if cleanStart {
		err := cmd.StopHardHat()
		assert.Nilf(t, err, "Failed to stopHardHat")

		details = startHardHat(t, ctx, workingDir, validators[0])
		assert.NotNilf(t, details, "Expected details to be not nil")
	}

	network, err := getEthereumClient(workingDir, validators[0])
	assert.Nilf(t, err, "Failed to build Ethereum endpoint")
	assert.NotNilf(t, network, "Ethereum network should not be Nil")

	log.Printf("miasmo")
}

func getEthereumClient(workingDir string, defaultAccount string) (*ethereum.Client, error) {

	//root := cmd.GetProjectRootPath()
	//assetKey := filepath.Join(root, "assets", "test", "keys")
	//assetPasscode := filepath.Join(root, "assets", "test", "passcodes.txt")
	//assetKey := filepath.Join(workingDir, "scripts", "generated", "keystores", "keys")
	//assetPasscode := filepath.Join(workingDir, "scripts", "generated", "keystores", "passcodes.txt")
	client, err := ethereum.NewClient(
		"http://localhost:8545",
		filepath.Join(workingDir, "keystores", "keys"),
		filepath.Join(workingDir, "keystores", "passcodes.txt"),
		defaultAccount,
		false,
		0,
		300,
		1,
	)
	return client, err
}

func startHardHat(t *testing.T, ctx context.Context, workingDir string, defaultAccount string) *ethereum.Client {

	log.Printf("Starting HardHat ...")
	err := cmd.RunHardHatNode()
	assert.Nilf(t, err, "Error starting hardhat node")

	err = cmd.WaitForHardHatNode(ctx)
	assert.Nilf(t, err, "Failed to wait for hardhat to be up and running")

	client, err := getEthereumClient(workingDir, defaultAccount)
	assert.Nilf(t, err, "Failed to build Ethereum endpoint")
	assert.NotNilf(t, client, "Ethereum network should not be Nil")

	log.Printf("Deploying contracts ...")
	factoryAddress, err := cmd.RunDeploy(workingDir)
	if err != nil {
		client.Close()
		assert.Nilf(t, err, "Error deploying contracts: %v", err)
		return nil
	}
	addr := common.Address{}
	copy(addr[:], common.FromHex(factoryAddress))

	return client
}

// We complete everything correctly, happy path
func TestCompletion_Group_1_AllGood(t *testing.T) {
	n := 4
	setupTest(t, n, false)

	//err := testutils.InitializeValidatorFiles(5)
	//assert.Nil(t, err)
	//
	//suite := dkgTestUtils.StartFromMPKSubmissionPhase(t, n, 100)
	//defer suite.Eth.Close()
	//ctx := context.Background()
	//eth := suite.Eth
	//dkgStates := suite.DKGStates
	//logger := logging.GetLogger("test").WithField("Validator", "")
	//
	//// Do GPKj Submission task
	//for idx := 0; idx < n; idx++ {
	//
	//	err := suite.GpkjSubmissionTasks[idx].Initialize(ctx, logger, eth)
	//	assert.Nil(t, err)
	//	err = suite.GpkjSubmissionTasks[idx].DoWork(ctx, logger, eth)
	//	assert.Nil(t, err)
	//
	//	eth.Commit()
	//	assert.True(t, suite.GpkjSubmissionTasks[idx].Success)
	//}
	//
	//height, err := suite.Eth.GetCurrentHeight(ctx)
	//assert.Nil(t, err)
	//
	//disputeGPKjTasks := make([]*dkg.DisputeGPKjTask, n)
	//completionTasks := make([]*dkg.CompletionTask, n)
	//var completionStart uint64
	//for idx := 0; idx < n; idx++ {
	//	state := dkgStates[idx]
	//	disputeGPKjTask, completionTask := events.UpdateStateOnGPKJSubmissionComplete(state, height)
	//	disputeGPKjTasks[idx] = disputeGPKjTask
	//	completionTasks[idx] = completionTask
	//	completionStart = completionTask.GetStart()
	//}
	//
	//// Advance to Completion phase
	//testutils.AdvanceTo(t, eth, completionStart)
	//
	//for idx := 0; idx < n; idx++ {
	//	state := dkgStates[idx]
	//
	//	err := completionTasks[idx].Initialize(ctx, logger, eth)
	//	assert.Nil(t, err)
	//	amILeading := completionTasks[idx].AmILeading(ctx, eth, logger, state)
	//	err = completionTasks[idx].DoWork(ctx, logger, eth)
	//	if amILeading {
	//		assert.Nil(t, err)
	//		assert.True(t, completionTasks[idx].Success)
	//	} else {
	//		if completionTasks[idx].ShouldExecute(ctx, logger, eth) {
	//			assert.NotNil(t, err)
	//			assert.False(t, completionTasks[idx].Success)
	//		} else {
	//			assert.Nil(t, err)
	//			assert.True(t, completionTasks[idx].Success)
	//		}
	//
	//	}
	//}
}

//
//func TestCompletion_Group_1_StartFromCompletion(t *testing.T) {
//	n := 4
//	suite := dkgTestUtils.StartFromCompletion(t, n, 100)
//	defer suite.Eth.Close()
//	ctx := context.Background()
//	eth := suite.Eth
//	dkgStates := suite.DKGStates
//	logger := logging.GetLogger("test").WithField("Validator", "k")
//
//	// Do Completion task
//	var hasLeader bool
//	for idx := 0; idx < n; idx++ {
//		state := dkgStates[idx]
//		task := suite.CompletionTasks[idx]
//
//		err := task.Initialize(ctx, logger, eth)
//		assert.Nil(t, err)
//		amILeading := task.AmILeading(ctx, eth, logger, state)
//
//		if amILeading {
//			hasLeader = true
//			err = task.DoWork(ctx, logger, eth)
//			eth.Commit()
//			assert.Nil(t, err)
//			assert.False(t, task.ShouldRetry(ctx, logger, eth))
//			assert.True(t, task.Success)
//		}
//	}
//
//	assert.True(t, hasLeader)
//	assert.False(t, suite.CompletionTasks[0].ShouldRetry(ctx, logger, eth))
//
//	callOpts, err := eth.GetCallOpts(ctx, eth.GetDefaultAccount())
//	assert.Nil(t, err)
//
//	phase, err := suite.Eth.Contracts().Ethdkg().GetETHDKGPhase(callOpts)
//	assert.Nil(t, err)
//	assert.Equal(t, uint8(state.Completion), phase)
//
//	// event
//	for j := 0; j < n; j++ {
//		// simulate receiving ValidatorSetCompleted event for all participants
//		suite.DKGStates[j].OnCompletion()
//		assert.Equal(t, suite.DKGStates[j].Phase, state.Completion)
//	}
//}
//
//// We begin by submitting invalid information.
//// This test is meant to raise an error resulting from an invalid argument
//// for the Ethereum interface.
//func TestCompletion_Group_2_Bad1(t *testing.T) {
//	n := 4
//	ecdsaPrivateKeys, _ := testutils.InitializePrivateKeysAndAccounts(n)
//	logger := logging.GetLogger("ethereum")
//	logger.SetLevel(logrus.DebugLevel)
//	eth := testutils.ConnectSimulatorEndpoint(t, ecdsaPrivateKeys, 333*time.Millisecond)
//	defer eth.Close()
//
//	acct := eth.GetKnownAccounts()[0]
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	// Create a task to share distribution and make sure it succeeds
//	state := state.NewDkgState(acct)
//	task := dkg.NewCompletionTask(state, 1, 100)
//	log := logger.WithField("TaskID", "foo")
//
//	err := task.Initialize(ctx, log, eth)
//	assert.NotNil(t, err)
//}
//
//// We test to ensure that everything behaves correctly.
//func TestCompletion_Group_2_Bad2(t *testing.T) {
//	n := 4
//	ecdsaPrivateKeys, _ := testutils.InitializePrivateKeysAndAccounts(n)
//	logger := logging.GetLogger("ethereum")
//	logger.SetLevel(logrus.DebugLevel)
//	eth := testutils.ConnectSimulatorEndpoint(t, ecdsaPrivateKeys, 333*time.Millisecond)
//	defer eth.Close()
//
//	acct := eth.GetKnownAccounts()[0]
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	// Do bad Completion task
//	state := state.NewDkgState(acct)
//	log := logger.WithField("TaskID", "foo")
//	task := dkg.NewCompletionTask(state, 1, 100)
//
//	err := task.Initialize(ctx, log, eth)
//	if err == nil {
//		t.Fatal("Should have raised error")
//	}
//}
//
//// We complete everything correctly, but we do not complete in time
//func TestCompletion_Group_2_Bad3(t *testing.T) {
//	n := 4
//	suite := dkgTestUtils.StartFromMPKSubmissionPhase(t, n, 100)
//	defer suite.Eth.Close()
//	ctx := context.Background()
//	eth := suite.Eth
//	dkgStates := suite.DKGStates
//	logger := logging.GetLogger("test").WithField("Validator", "")
//
//	// Do GPKj Submission task
//	tasksVec := suite.GpkjSubmissionTasks
//	for idx := 0; idx < n; idx++ {
//
//		err := tasksVec[idx].Initialize(ctx, logger, eth)
//		assert.Nil(t, err)
//		err = tasksVec[idx].DoWork(ctx, logger, eth)
//		assert.Nil(t, err)
//
//		eth.Commit()
//		assert.True(t, tasksVec[idx].Success)
//	}
//
//	height, err := suite.Eth.GetCurrentHeight(ctx)
//	assert.Nil(t, err)
//	completionTasks := make([]*dkg.CompletionTask, n)
//	var completionStart, completionEnd uint64
//	for idx := 0; idx < n; idx++ {
//		state := dkgStates[idx]
//		_, completionTask := events.UpdateStateOnGPKJSubmissionComplete(state, height)
//		completionTasks[idx] = completionTask
//		completionStart = completionTask.GetStart()
//		completionEnd = completionTask.GetEnd()
//	}
//
//	// Advance to Completion phase
//	testutils.AdvanceTo(t, eth, completionStart)
//
//	// Advance to end of Completion phase
//	testutils.AdvanceTo(t, eth, completionEnd)
//	eth.Commit()
//
//	err = completionTasks[0].Initialize(ctx, logger, eth)
//	if err != nil {
//		t.Fatal(err)
//	}
//	err = completionTasks[0].DoWork(ctx, logger, eth)
//	if err == nil {
//		t.Fatal("Should have raised error")
//	}
//}
//
//func TestCompletion_Group_3_ShouldRetry_returnsFalse(t *testing.T) {
//	n := 4
//	suite := dkgTestUtils.StartFromCompletion(t, n, 40)
//	defer suite.Eth.Close()
//	ctx := context.Background()
//	eth := suite.Eth
//	dkgStates := suite.DKGStates
//	logger := logging.GetLogger("test").WithField("Validator", "")
//
//	// Do Completion task
//	tasksVec := suite.CompletionTasks
//	var hadLeaders bool
//	for idx := 0; idx < n; idx++ {
//		state := dkgStates[idx]
//
//		err := tasksVec[idx].Initialize(ctx, logger, eth)
//		assert.Nil(t, err)
//		amILeading := tasksVec[idx].AmILeading(ctx, eth, logger, state)
//
//		if amILeading {
//			hadLeaders = true
//			// only perform ETHDKG completion if validator is leading
//			assert.True(t, tasksVec[idx].ShouldRetry(ctx, logger, eth))
//			err = tasksVec[idx].DoWork(ctx, logger, eth)
//			assert.Nil(t, err)
//			assert.False(t, tasksVec[idx].ShouldRetry(ctx, logger, eth))
//		}
//	}
//
//	assert.True(t, hadLeaders)
//
//	// any task is able to tell if ETHDKG still needs completion
//	// if for any reason no validator lead the process,
//	// then all tasks will have ShouldExecute() returning true
//	assert.False(t, tasksVec[0].ShouldRetry(ctx, logger, eth))
//}
//
//func TestCompletion_Group_3_ShouldRetry_returnsTrue(t *testing.T) {
//	n := 4
//	suite := dkgTestUtils.StartFromMPKSubmissionPhase(t, n, 100)
//	defer suite.Eth.Close()
//	ctx := context.Background()
//	eth := suite.Eth
//	dkgStates := suite.DKGStates
//	logger := logging.GetLogger("test").WithField("Validator", "")
//
//	// Do GPKj Submission task
//	for idx := 0; idx < n; idx++ {
//		err := suite.GpkjSubmissionTasks[idx].Initialize(ctx, logger, eth)
//		assert.Nil(t, err)
//		err = suite.GpkjSubmissionTasks[idx].DoWork(ctx, logger, eth)
//		assert.Nil(t, err)
//
//		eth.Commit()
//		assert.True(t, suite.GpkjSubmissionTasks[idx].Success)
//	}
//
//	height, err := suite.Eth.GetCurrentHeight(ctx)
//	assert.Nil(t, err)
//
//	disputeGPKjTasks := make([]*dkg.DisputeGPKjTask, n)
//	completionTasks := make([]*dkg.CompletionTask, n)
//	var completionStart uint64
//	for idx := 0; idx < n; idx++ {
//		state := dkgStates[idx]
//		disputeGPKjTask, completionTask := events.UpdateStateOnGPKJSubmissionComplete(state, height)
//		disputeGPKjTasks[idx] = disputeGPKjTask
//		completionTasks[idx] = completionTask
//		completionStart = completionTask.GetStart()
//	}
//
//	// Advance to Completion phase
//	testutils.AdvanceTo(t, eth, completionStart)
//	eth.Commit()
//
//	err = completionTasks[0].Initialize(ctx, logger, eth)
//	assert.Nil(t, err)
//
//	shouldRetry := completionTasks[0].ShouldExecute(ctx, logger, eth)
//	assert.True(t, shouldRetry)
//}
