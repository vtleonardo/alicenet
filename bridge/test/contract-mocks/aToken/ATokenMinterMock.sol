// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.16;

import "contracts/interfaces/IStakingToken.sol";

import "contracts/libraries/tokens/StakingToken.sol";
import "contracts/utils/ImmutableAuth.sol";

contract ATokenMinterMock is ImmutableFactory, StakingToken {
    constructor(address stakingAddress_)
        ImmutableFactory(msg.sender)
        StakingToken(stakingAddress_)
    {}

    function mint(address to, uint256 amount) public {
        IStakingToken(_stakingTokenAddress()).externalMint(to, amount);
    }
}
