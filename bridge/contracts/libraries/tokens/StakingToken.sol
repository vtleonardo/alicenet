// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.16;

abstract contract StakingToken {
    address private immutable _stakingToken;
    error OnlyStakingToken(address sender, address expected);

    modifier onlyStakingToken() {
        if (msg.sender != _stakingToken) {
            revert OnlyStakingToken(msg.sender, _stakingToken);
        }
        _;
    }

    constructor(address stakingAddress_) {
        _stakingToken = stakingAddress_;
    }

    function _stakingTokenAddress() internal view returns (address) {
        return _stakingToken;
    }
}
