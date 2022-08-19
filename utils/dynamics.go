package utils

import "github.com/alicenet/alicenet/bridge/bindings"

func CompareCanonicalVersion(newVersion bindings.CanonicalVersion) (bool, bool, bool, bindings.CanonicalVersion) {
	localVersion := GetLocalVersion()
	return newVersion.Major > localVersion.Major,
		newVersion.Major == localVersion.Major && newVersion.Minor > localVersion.Minor,
		newVersion.Major == localVersion.Major && newVersion.Minor == localVersion.Minor && newVersion.Patch > localVersion.Patch,
		localVersion
}

func GetLocalVersion() bindings.CanonicalVersion {
	return bindings.CanonicalVersion{
		Major:      1,
		Minor:      4,
		Patch:      7,
		BinaryHash: [32]byte{},
	}
}