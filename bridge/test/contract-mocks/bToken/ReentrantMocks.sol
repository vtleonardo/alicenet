// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.16;

import "contracts/utils/ImmutableAuth.sol";
import "contracts/interfaces/IUtilityToken.sol";
import "contracts/libraries/tokens/UtilityToken.sol";
import "contracts/utils/MagicEthTransfer.sol";

contract ReentrantLoopDistributionMock is MagicEthTransfer, ImmutableFactory, UtilityToken {
    uint256 internal _counter;

    constructor(address utilityAddress_)
        ImmutableFactory(msg.sender)
        UtilityToken(utilityAddress_)
    {}

    receive() external payable {
        _internalLoop();
    }

    function depositEth(uint8 magic_) public payable checkMagic(magic_) {
        _internalLoop();
    }

    function _internalLoop() internal {
        _counter++;
        if (_counter <= 3) {
            IUtilityToken(_utilityTokenAddress()).distribute();
        }
    }
}
