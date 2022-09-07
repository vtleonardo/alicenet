// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.16;

import "contracts/interfaces/IStakingToken.sol";
import "contracts/libraries/tokens/StakingToken.sol";
import "contracts/utils/ImmutableAuth.sol";

/// @custom:salt ATokenBurner
/// @custom:deploy-type deployUpgradeable
contract ATokenBurner is ImmutableFactory, StakingToken, IStakingTokenBurner {
    constructor(address stakingAddress_)
        ImmutableFactory(msg.sender)
        StakingToken(stakingAddress_)
        IStakingTokenBurner()
    {}

    function burn(address to, uint256 amount) public onlyFactory {
        IStakingToken(_stakingTokenAddress()).externalBurn(to, amount);
    }
}
