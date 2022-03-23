// This file is auto-generated by hardhat generate-immutable-auth-contract task. DO NOT EDIT.
// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

import "./DeterministicAddress.sol";

abstract contract ImmutableFactory is DeterministicAddress {
    address private immutable _factory;

    modifier onlyFactory() {
        require(msg.sender == _factory, "onlyFactory");
        _;
    }

    constructor(address factory_) {
        _factory = factory_;
    }

    function _factoryAddress() internal view returns (address) {
        return _factory;
    }
}

abstract contract ImmutableAToken is ImmutableFactory {
    address private immutable _aToken;

    modifier onlyAToken() {
        require(msg.sender == _aToken, "onlyAToken");
        _;
    }

    constructor() {
        _aToken = getMetamorphicContractAddress(
            0x41546f6b656e0000000000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _aTokenAddress() internal view returns (address) {
        return _aToken;
    }

    function _saltForAToken() internal pure returns (bytes32) {
        return 0x41546f6b656e0000000000000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableFoundation is ImmutableFactory {
    address private immutable _foundation;

    modifier onlyFoundation() {
        require(msg.sender == _foundation, "onlyFoundation");
        _;
    }

    constructor() {
        _foundation = getMetamorphicContractAddress(
            0x466f756e646174696f6e00000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _foundationAddress() internal view returns (address) {
        return _foundation;
    }

    function _saltForFoundation() internal pure returns (bytes32) {
        return 0x466f756e646174696f6e00000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableGovernance is ImmutableFactory {
    address private immutable _governance;

    modifier onlyGovernance() {
        require(msg.sender == _governance, "onlyGovernance");
        _;
    }

    constructor() {
        _governance = getMetamorphicContractAddress(
            0x476f7665726e616e636500000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _governanceAddress() internal view returns (address) {
        return _governance;
    }

    function _saltForGovernance() internal pure returns (bytes32) {
        return 0x476f7665726e616e636500000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableLiquidityProviderStaking is ImmutableFactory {
    address private immutable _liquidityProviderStaking;

    modifier onlyLiquidityProviderStaking() {
        require(msg.sender == _liquidityProviderStaking, "onlyLiquidityProviderStaking");
        _;
    }

    constructor() {
        _liquidityProviderStaking = getMetamorphicContractAddress(
            0x4c697175696469747950726f76696465725374616b696e670000000000000000,
            _factoryAddress()
        );
    }

    function _liquidityProviderStakingAddress() internal view returns (address) {
        return _liquidityProviderStaking;
    }

    function _saltForLiquidityProviderStaking() internal pure returns (bytes32) {
        return 0x4c697175696469747950726f76696465725374616b696e670000000000000000;
    }
}

abstract contract ImmutableMadByte is ImmutableFactory {
    address private immutable _madByte;

    modifier onlyMadByte() {
        require(msg.sender == _madByte, "onlyMadByte");
        _;
    }

    constructor() {
        _madByte = getMetamorphicContractAddress(
            0x4d61644279746500000000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _madByteAddress() internal view returns (address) {
        return _madByte;
    }

    function _saltForMadByte() internal pure returns (bytes32) {
        return 0x4d61644279746500000000000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableMadToken is ImmutableFactory {
    address private immutable _madToken;

    modifier onlyMadToken() {
        require(msg.sender == _madToken, "onlyMadToken");
        _;
    }

    constructor() {
        _madToken = getMetamorphicContractAddress(
            0x4d6164546f6b656e000000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _madTokenAddress() internal view returns (address) {
        return _madToken;
    }

    function _saltForMadToken() internal pure returns (bytes32) {
        return 0x4d6164546f6b656e000000000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutablePublicStaking is ImmutableFactory {
    address private immutable _publicStaking;

    modifier onlyPublicStaking() {
        require(msg.sender == _publicStaking, "onlyPublicStaking");
        _;
    }

    constructor() {
        _publicStaking = getMetamorphicContractAddress(
            0x5075626c69635374616b696e6700000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _publicStakingAddress() internal view returns (address) {
        return _publicStaking;
    }

    function _saltForPublicStaking() internal pure returns (bytes32) {
        return 0x5075626c69635374616b696e6700000000000000000000000000000000000000;
    }
}

abstract contract ImmutableSnapshots is ImmutableFactory {
    address private immutable _snapshots;

    modifier onlySnapshots() {
        require(msg.sender == _snapshots, "onlySnapshots");
        _;
    }

    constructor() {
        _snapshots = getMetamorphicContractAddress(
            0x536e617073686f74730000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _snapshotsAddress() internal view returns (address) {
        return _snapshots;
    }

    function _saltForSnapshots() internal pure returns (bytes32) {
        return 0x536e617073686f74730000000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableStakeNFT is ImmutableFactory {
    address private immutable _stakeNFT;

    modifier onlyStakeNFT() {
        require(msg.sender == _stakeNFT, "onlyStakeNFT");
        _;
    }

    constructor() {
        _stakeNFT = getMetamorphicContractAddress(
            0x5374616b654e4654000000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _stakeNFTAddress() internal view returns (address) {
        return _stakeNFT;
    }

    function _saltForStakeNFT() internal pure returns (bytes32) {
        return 0x5374616b654e4654000000000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableStakeNFTLP is ImmutableFactory {
    address private immutable _stakeNFTLP;

    modifier onlyStakeNFTLP() {
        require(msg.sender == _stakeNFTLP, "onlyStakeNFTLP");
        _;
    }

    constructor() {
        _stakeNFTLP = getMetamorphicContractAddress(
            0x5374616b654e46544c5000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _stakeNFTLPAddress() internal view returns (address) {
        return _stakeNFTLP;
    }

    function _saltForStakeNFTLP() internal pure returns (bytes32) {
        return 0x5374616b654e46544c5000000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableStakeNFTPositionDescriptor is ImmutableFactory {
    address private immutable _stakeNFTPositionDescriptor;

    modifier onlyStakeNFTPositionDescriptor() {
        require(msg.sender == _stakeNFTPositionDescriptor, "onlyStakeNFTPositionDescriptor");
        _;
    }

    constructor() {
        _stakeNFTPositionDescriptor = getMetamorphicContractAddress(
            0x5374616b654e4654506f736974696f6e44657363726970746f72000000000000,
            _factoryAddress()
        );
    }

    function _stakeNFTPositionDescriptorAddress() internal view returns (address) {
        return _stakeNFTPositionDescriptor;
    }

    function _saltForStakeNFTPositionDescriptor() internal pure returns (bytes32) {
        return 0x5374616b654e4654506f736974696f6e44657363726970746f72000000000000;
    }
}

abstract contract ImmutableValidatorNFT is ImmutableFactory {
    address private immutable _validatorNFT;

    modifier onlyValidatorNFT() {
        require(msg.sender == _validatorNFT, "onlyValidatorNFT");
        _;
    }

    constructor() {
        _validatorNFT = getMetamorphicContractAddress(
            0x56616c696461746f724e46540000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _validatorNFTAddress() internal view returns (address) {
        return _validatorNFT;
    }

    function _saltForValidatorNFT() internal pure returns (bytes32) {
        return 0x56616c696461746f724e46540000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableValidatorPool is ImmutableFactory {
    address private immutable _validatorPool;

    modifier onlyValidatorPool() {
        require(msg.sender == _validatorPool, "onlyValidatorPool");
        _;
    }

    constructor() {
        _validatorPool = getMetamorphicContractAddress(
            0x56616c696461746f72506f6f6c00000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _validatorPoolAddress() internal view returns (address) {
        return _validatorPool;
    }

    function _saltForValidatorPool() internal pure returns (bytes32) {
        return 0x56616c696461746f72506f6f6c00000000000000000000000000000000000000;
    }
}

abstract contract ImmutableValidatorStaking is ImmutableFactory {
    address private immutable _validatorStaking;

    modifier onlyValidatorStaking() {
        require(msg.sender == _validatorStaking, "onlyValidatorStaking");
        _;
    }

    constructor() {
        _validatorStaking = getMetamorphicContractAddress(
            0x56616c696461746f725374616b696e6700000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _validatorStakingAddress() internal view returns (address) {
        return _validatorStaking;
    }

    function _saltForValidatorStaking() internal pure returns (bytes32) {
        return 0x56616c696461746f725374616b696e6700000000000000000000000000000000;
    }
}

abstract contract ImmutableATokenBurner is ImmutableFactory {
    address private immutable _aTokenBurner;

    modifier onlyATokenBurner() {
        require(msg.sender == _aTokenBurner, "onlyATokenBurner");
        _;
    }

    constructor() {
        _aTokenBurner = getMetamorphicContractAddress(
            0x41546f6b656e4275726e65720000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _aTokenBurnerAddress() internal view returns (address) {
        return _aTokenBurner;
    }

    function _saltForATokenBurner() internal pure returns (bytes32) {
        return 0x41546f6b656e4275726e65720000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableATokenMinter is ImmutableFactory {
    address private immutable _aTokenMinter;

    modifier onlyATokenMinter() {
        require(msg.sender == _aTokenMinter, "onlyATokenMinter");
        _;
    }

    constructor() {
        _aTokenMinter = getMetamorphicContractAddress(
            0x41546f6b656e4d696e7465720000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _aTokenMinterAddress() internal view returns (address) {
        return _aTokenMinter;
    }

    function _saltForATokenMinter() internal pure returns (bytes32) {
        return 0x41546f6b656e4d696e7465720000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableETHDKGAccusations is ImmutableFactory {
    address private immutable _ethdkgAccusations;

    modifier onlyETHDKGAccusations() {
        require(msg.sender == _ethdkgAccusations, "onlyETHDKGAccusations");
        _;
    }

    constructor() {
        _ethdkgAccusations = getMetamorphicContractAddress(
            0x455448444b4741636375736174696f6e73000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _ethdkgAccusationsAddress() internal view returns (address) {
        return _ethdkgAccusations;
    }

    function _saltForETHDKGAccusations() internal pure returns (bytes32) {
        return 0x455448444b4741636375736174696f6e73000000000000000000000000000000;
    }
}

abstract contract ImmutableETHDKGPhases is ImmutableFactory {
    address private immutable _ethdkgPhases;

    modifier onlyETHDKGPhases() {
        require(msg.sender == _ethdkgPhases, "onlyETHDKGPhases");
        _;
    }

    constructor() {
        _ethdkgPhases = getMetamorphicContractAddress(
            0x455448444b475068617365730000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _ethdkgPhasesAddress() internal view returns (address) {
        return _ethdkgPhases;
    }

    function _saltForETHDKGPhases() internal pure returns (bytes32) {
        return 0x455448444b475068617365730000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableETHDKG is ImmutableFactory {
    address private immutable _ethdkg;

    modifier onlyETHDKG() {
        require(msg.sender == _ethdkg, "onlyETHDKG");
        _;
    }

    constructor() {
        _ethdkg = getMetamorphicContractAddress(
            0x455448444b470000000000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _ethdkgAddress() internal view returns (address) {
        return _ethdkg;
    }

    function _saltForETHDKG() internal pure returns (bytes32) {
        return 0x455448444b470000000000000000000000000000000000000000000000000000;
    }
}
