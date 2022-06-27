package cmd

import (
	"os"
	"path/filepath"
)

func RunSetup(workingDir string) error {

	rootPath := GetProjectRootPath()

	// Generate working dir structure
	err := os.MkdirAll(filepath.Join(workingDir, "keystores", "keys"), os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(workingDir, "config"), os.ModePerm)
	if err != nil {
		return err
	}

	// Copy files
	_, err = CopyFileToFolder(filepath.Join(rootPath, "scripts", "base-files", "deploymentList"), filepath.Join(workingDir, "deploymentList"))
	if err != nil {
		return err
	}

	_, err = CopyFileToFolder(filepath.Join(rootPath, "scripts", "base-files", "deploymentArgsTemplate"), filepath.Join(workingDir, "deploymentArgsTemplate"))
	if err != nil {
		return err
	}

	_, err = CopyFileToFolder(filepath.Join(rootPath, "scripts", "base-files", "genesis.json"), filepath.Join(workingDir, "genesis.json"))
	if err != nil {
		return err
	}

	_, err = CopyFileToFolder(filepath.Join(rootPath, "scripts", "base-files", "baseConfig"), filepath.Join(workingDir, "baseConfig"))
	if err != nil {
		return err
	}

	// TODO - autogenerate this file
	_, err = CopyFileToFolder(filepath.Join(rootPath, "scripts", "base-files", "0x546f99f244b7b58b855330ae0e2bc1b30b41302f"), filepath.Join(workingDir, "keystores", "keys", "0x546f99f244b7b58b855330ae0e2bc1b30b41302f"))
	if err != nil {
		return err
	}

	return nil
}
