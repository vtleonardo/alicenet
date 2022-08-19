// This file is auto-generated by hardhat generate-immutable-auth-contract task. DO NOT EDIT.
// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

import "contracts/utils/DeterministicAddress.sol";

abstract contract ImmutableFactory is DeterministicAddress {
    address private immutable _factory;
    error OnlyFactory(address sender, address expected);

    modifier onlyFactory() {
        if (msg.sender != _factory) {
            revert OnlyFactory(msg.sender, _factory);
        }
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
    error OnlyAToken(address sender, address expected);

    modifier onlyAToken() {
        if (msg.sender != _aToken) {
            revert OnlyAToken(msg.sender, _aToken);
        }
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

abstract contract ImmutableATokenBurner is ImmutableFactory {
    address private immutable _aTokenBurner;
    error OnlyATokenBurner(address sender, address expected);

    modifier onlyATokenBurner() {
        if (msg.sender != _aTokenBurner) {
            revert OnlyATokenBurner(msg.sender, _aTokenBurner);
        }
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
    error OnlyATokenMinter(address sender, address expected);

    modifier onlyATokenMinter() {
        if (msg.sender != _aTokenMinter) {
            revert OnlyATokenMinter(msg.sender, _aTokenMinter);
        }
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

abstract contract ImmutableBToken is ImmutableFactory {
    address private immutable _bToken;
    error OnlyBToken(address sender, address expected);

    modifier onlyBToken() {
        if (msg.sender != _bToken) {
            revert OnlyBToken(msg.sender, _bToken);
        }
        _;
    }

    constructor() {
        _bToken = getMetamorphicContractAddress(
            0x42546f6b656e0000000000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _bTokenAddress() internal view returns (address) {
        return _bToken;
    }

    function _saltForBToken() internal pure returns (bytes32) {
        return 0x42546f6b656e0000000000000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableDynamics is ImmutableFactory {
    address private immutable _dynamics;
    error OnlyDynamics(address sender, address expected);

    modifier onlyDynamics() {
        if (msg.sender != _dynamics) {
            revert OnlyDynamics(msg.sender, _dynamics);
        }
        _;
    }

    constructor() {
        _dynamics = getMetamorphicContractAddress(
            0x44796e616d696373000000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _dynamicsAddress() internal view returns (address) {
        return _dynamics;
    }

    function _saltForDynamics() internal pure returns (bytes32) {
        return 0x44796e616d696373000000000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableFoundation is ImmutableFactory {
    address private immutable _foundation;
    error OnlyFoundation(address sender, address expected);

    modifier onlyFoundation() {
        if (msg.sender != _foundation) {
            revert OnlyFoundation(msg.sender, _foundation);
        }
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
    error OnlyGovernance(address sender, address expected);

    modifier onlyGovernance() {
        if (msg.sender != _governance) {
            revert OnlyGovernance(msg.sender, _governance);
        }
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

abstract contract ImmutableInvalidTxConsumptionAccusation is ImmutableFactory {
    address private immutable _invalidTxConsumptionAccusation;
    error OnlyInvalidTxConsumptionAccusation(address sender, address expected);

    modifier onlyInvalidTxConsumptionAccusation() {
        if (msg.sender != _invalidTxConsumptionAccusation) {
            revert OnlyInvalidTxConsumptionAccusation(msg.sender, _invalidTxConsumptionAccusation);
        }
        _;
    }

    constructor() {
        _invalidTxConsumptionAccusation = getMetamorphicContractAddress(
            0x92a73f2b6573522d63c8fc84b5d8e5d615fbb685c1b3d7fad2155fe227daf848,
            _factoryAddress()
        );
    }

    function _invalidTxConsumptionAccusationAddress() internal view returns (address) {
        return _invalidTxConsumptionAccusation;
    }

    function _saltForInvalidTxConsumptionAccusation() internal pure returns (bytes32) {
        return 0x92a73f2b6573522d63c8fc84b5d8e5d615fbb685c1b3d7fad2155fe227daf848;
    }
}

abstract contract ImmutableLiquidityProviderStaking is ImmutableFactory {
    address private immutable _liquidityProviderStaking;
    error OnlyLiquidityProviderStaking(address sender, address expected);

    modifier onlyLiquidityProviderStaking() {
        if (msg.sender != _liquidityProviderStaking) {
            revert OnlyLiquidityProviderStaking(msg.sender, _liquidityProviderStaking);
        }
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

abstract contract ImmutableMultipleProposalAccusation is ImmutableFactory {
    address private immutable _multipleProposalAccusation;
    error OnlyMultipleProposalAccusation(address sender, address expected);

    modifier onlyMultipleProposalAccusation() {
        if (msg.sender != _multipleProposalAccusation) {
            revert OnlyMultipleProposalAccusation(msg.sender, _multipleProposalAccusation);
        }
        _;
    }

    constructor() {
        _multipleProposalAccusation = getMetamorphicContractAddress(
            0xcfdffd500b4a956e03976b2afd69712237ffa06e35093df1e05e533688959fdc,
            _factoryAddress()
        );
    }

    function _multipleProposalAccusationAddress() internal view returns (address) {
        return _multipleProposalAccusation;
    }

    function _saltForMultipleProposalAccusation() internal pure returns (bytes32) {
        return 0xcfdffd500b4a956e03976b2afd69712237ffa06e35093df1e05e533688959fdc;
    }
}

abstract contract ImmutablePublicStaking is ImmutableFactory {
    address private immutable _publicStaking;
    error OnlyPublicStaking(address sender, address expected);

    modifier onlyPublicStaking() {
        if (msg.sender != _publicStaking) {
            revert OnlyPublicStaking(msg.sender, _publicStaking);
        }
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
    error OnlySnapshots(address sender, address expected);

    modifier onlySnapshots() {
        if (msg.sender != _snapshots) {
            revert OnlySnapshots(msg.sender, _snapshots);
        }
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

abstract contract ImmutableStakingPositionDescriptor is ImmutableFactory {
    address private immutable _stakingPositionDescriptor;
    error OnlyStakingPositionDescriptor(address sender, address expected);

    modifier onlyStakingPositionDescriptor() {
        if (msg.sender != _stakingPositionDescriptor) {
            revert OnlyStakingPositionDescriptor(msg.sender, _stakingPositionDescriptor);
        }
        _;
    }

    constructor() {
        _stakingPositionDescriptor = getMetamorphicContractAddress(
            0x5374616b696e67506f736974696f6e44657363726970746f7200000000000000,
            _factoryAddress()
        );
    }

    function _stakingPositionDescriptorAddress() internal view returns (address) {
        return _stakingPositionDescriptor;
    }

    function _saltForStakingPositionDescriptor() internal pure returns (bytes32) {
        return 0x5374616b696e67506f736974696f6e44657363726970746f7200000000000000;
    }
}

abstract contract ImmutableValidatorPool is ImmutableFactory {
    address private immutable _validatorPool;
    error OnlyValidatorPool(address sender, address expected);

    modifier onlyValidatorPool() {
        if (msg.sender != _validatorPool) {
            revert OnlyValidatorPool(msg.sender, _validatorPool);
        }
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
    error OnlyValidatorStaking(address sender, address expected);

    modifier onlyValidatorStaking() {
        if (msg.sender != _validatorStaking) {
            revert OnlyValidatorStaking(msg.sender, _validatorStaking);
        }
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

abstract contract ImmutableETHDKGAccusations is ImmutableFactory {
    address private immutable _ethdkgAccusations;
    error OnlyETHDKGAccusations(address sender, address expected);

    modifier onlyETHDKGAccusations() {
        if (msg.sender != _ethdkgAccusations) {
            revert OnlyETHDKGAccusations(msg.sender, _ethdkgAccusations);
        }
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
    error OnlyETHDKGPhases(address sender, address expected);

    modifier onlyETHDKGPhases() {
        if (msg.sender != _ethdkgPhases) {
            revert OnlyETHDKGPhases(msg.sender, _ethdkgPhases);
        }
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
    error OnlyETHDKG(address sender, address expected);

    modifier onlyETHDKG() {
        if (msg.sender != _ethdkg) {
            revert OnlyETHDKG(msg.sender, _ethdkg);
        }
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
