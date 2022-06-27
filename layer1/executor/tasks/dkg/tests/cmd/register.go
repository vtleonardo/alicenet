package cmd

import (
	"strings"
)

func RunRegister(factoryAddress string, validators []string) error {

	bridgeDir := GetBridgePath()

	// Register validator
	_, _, err := executeCommand(bridgeDir, "npx", "hardhat --network dev --show-stack-traces registerValidators --factory-address", factoryAddress, strings.Join(validators, " "))
	if err != nil {
		return err
	}

	return nil
}
