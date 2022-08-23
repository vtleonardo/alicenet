// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.16;

import "contracts/interfaces/IStakingToken.sol";
import "contracts/libraries/tokens/StakingToken.sol";
import "contracts/utils/ImmutableAuth.sol";

/// @custom:salt ATokenMinter
/// @custom:deploy-type deployUpgradeable
contract ATokenMinter is ImmutableFactory, StakingToken, IStakingTokenMinter {
    constructor(address stakingAddress_)
        ImmutableFactory(msg.sender)
        StakingToken(stakingAddress_)
        IStakingTokenMinter()
    {}

    function mint(address to, uint256 amount) public onlyFactory {
        IStakingToken(_stakingTokenAddress()).externalMint(to, amount);
    }
}
