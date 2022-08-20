// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

import "contracts/BToken.sol";
import "hardhat/console.sol";

contract BondingCurveStressMock is BToken {
    error InvalidBurnMintConversion(uint256 input, uint256 output, uint256 poolBalance);

    constructor() BToken() {}

    receive() external payable {}

    function stressMint(
        bool printValues,
        uint256 step,
        uint256 numIterations
    ) public payable {
        for (uint256 i = 1; i <= numIterations; i++) {
            uint256 mintedBToken;
            if (getPoolBalance() == 0) {
                mintedBToken = BToken._mint(msg.sender, BToken._MARKET_SPREAD, 0);
                if (printValues) {
                    console.log("%s,%s,%s", BToken.getPoolBalance(), mintedBToken, 1);
                }
            }
            mintedBToken = BToken._mint(msg.sender, step, 0);
            if (printValues) {
                console.log("%s,%s,%s", BToken.getPoolBalance(), mintedBToken, step);
            }
        }
    }

    function stressBurn(
        bool printValues,
        uint256 step,
        uint256 numIterations
    ) public payable {
        for (uint256 i = 1; i <= numIterations; i++) {
            uint256 receivedEth = BToken._burn(msg.sender, msg.sender, step, 0);
            if (printValues) {
                console.log("%s,%s,%s", BToken.getPoolBalance(), step, receivedEth);
            }
        }
    }

    function stressMintAndBurn(
        bool printValues,
        uint256 step,
        uint256 numIterations
    ) public payable {
        for (uint256 i = 1; i <= numIterations; i++) {
            uint256 mintedBToken = BToken._mint(msg.sender, step, 0);
            uint256 receivedEth = BToken._burn(msg.sender, msg.sender, mintedBToken, 0);
            if (printValues) {
                console.log("%s,%s,%s", BToken.getPoolBalance(), step, receivedEth);
            }
            if (step / receivedEth != _MARKET_SPREAD) {
                revert InvalidBurnMintConversion(step, receivedEth, BToken.getPoolBalance());
            }
        }
    }
}
