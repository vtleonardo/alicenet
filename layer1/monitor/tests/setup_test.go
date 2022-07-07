package tests

import (
	"github.com/alicenet/alicenet/layer1/ethereum"
	"github.com/alicenet/alicenet/logging"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupEthereum(t *testing.T, n int) *tests.ClientFixture {
	logger := logging.GetLogger("test").WithField("test", t.Name())
	fixture := tests.NewClientFixture(HardHat, 0, n, logger, true, true, true)
	assert.NotNil(t, fixture)

	eth := fixture.Client
	assert.NotNil(t, eth)
	assert.Equal(t, n, len(eth.GetKnownAccounts()))

	t.Cleanup(func() {
		fixture.Close()
		ethereum.CleanGlobalVariables(t)
	})

	return fixture
}
