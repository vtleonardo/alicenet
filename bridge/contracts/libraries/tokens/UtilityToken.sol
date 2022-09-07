// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.16;

abstract contract UtilityToken {
    address private immutable _utilityToken;
    error OnlyUtilityToken(address sender, address expected);

    modifier onlyUtilityToken() {
        if (msg.sender != _utilityToken) {
            revert OnlyUtilityToken(msg.sender, _utilityToken);
        }
        _;
    }

    constructor(address utilityAddress_) {
        _utilityToken = utilityAddress_;
    }

    function _utilityTokenAddress() internal view returns (address) {
        return _utilityToken;
    }
}
