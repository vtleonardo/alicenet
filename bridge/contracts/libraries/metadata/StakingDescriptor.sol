// SPDX-License-Identifier: GPL-2.0-or-later
pragma solidity ^0.8.11;

/* solhint-disable */
import "@openzeppelin/contracts/utils/Strings.sol";
import "@openzeppelin/contracts/utils/math/SafeMath.sol";
import "@openzeppelin/contracts/utils/math/SignedSafeMath.sol";
import "contracts/utils/Base64.sol";
import "contracts/libraries/metadata/StakingSVG.sol";

library StakingDescriptor {
    using Strings for uint256;

    struct ConstructTokenURIParams {
        uint256 tokenId;
        uint256 shares;
        uint256 freeAfter;
        uint256 withdrawFreeAfter;
        uint256 accumulatorEth;
        uint256 accumulatorToken;
    }

    function constructTokenURI(ConstructTokenURIParams memory params)
        internal
        pure
        returns (string memory)
    {
        string memory name = generateName(params);
        string memory description = generateDescription();
        string memory attributes = generateAttributes(
            params.tokenId.toString(),
            params.shares.toString(),
            params.freeAfter.toString(),
            params.withdrawFreeAfter.toString(),
            params.accumulatorEth.toString(),
            params.accumulatorToken.toString()
        );
        string memory image = Base64.encode(bytes(generateSVGImage(params)));

        return
            string(
                abi.encodePacked(
                    "data:application/json;utf8,",
                    bytes(
                        abi.encodePacked(
                            '{"name":"',
                            name,
                            '", "description":"',
                            description,
                            '", "attributes": ',
                            attributes,
                            ', "image_data": "',
                            "data:image/svg+xml;base64,",
                            image,
                            '"}'
                        )
                    )
                )
            );
    }

    function escapeQuotes(string memory symbol) internal pure returns (string memory) {
        bytes memory symbolBytes = bytes(symbol);
        uint8 quotesCount = 0;
        for (uint8 i = 0; i < symbolBytes.length; i++) {
            if (symbolBytes[i] == '"') {
                quotesCount++;
            }
        }
        if (quotesCount > 0) {
            bytes memory escapedBytes = new bytes(symbolBytes.length + (quotesCount));
            uint256 index;
            for (uint8 i = 0; i < symbolBytes.length; i++) {
                if (symbolBytes[i] == '"') {
                    escapedBytes[index++] = "\\";
                }
                escapedBytes[index++] = symbolBytes[i];
            }
            return string(escapedBytes);
        }
        return symbol;
    }

    function generateSVGImage(ConstructTokenURIParams memory params)
        internal
        pure
        returns (string memory svg)
    {
        StakingSVG.StakingSVGParams memory svgParams = StakingSVG.StakingSVGParams({
            shares: params.shares.toString(),
            freeAfter: params.freeAfter.toString(),
            withdrawFreeAfter: params.withdrawFreeAfter.toString(),
            accumulatorEth: params.accumulatorEth.toString(),
            accumulatorToken: params.accumulatorToken.toString()
        });

        return StakingSVG.generateSVG(svgParams);
    }

    function generateDescription() private pure returns (string memory) {
        return
            string(
                abi.encodePacked(
                    "This NFT represents a staked position on AliceNet. The owner of this NFT can modify or redeem the position."
                )
            );
    }

    function generateAttributes(
        string memory tokenId,
        string memory shares,
        string memory freeAfter,
        string memory withdrawFreeAfter,
        string memory accumulatorEth,
        string memory accumulatorToken
    ) private pure returns (string memory) {
        return
            string(
                abi.encodePacked(
                    "[",
                    '{"trait_type": "Shares", "value": "',
                    shares,
                    '"},'
                    '{"trait_type": "Free After", "value": "',
                    freeAfter,
                    '"},'
                    '{"trait_type": "Withdraw Free After", "value": "',
                    withdrawFreeAfter,
                    '"},'
                    '{"trait_type": "Accumulator Eth", "value": "',
                    accumulatorEth,
                    '"},'
                    '{"trait_type": "Accumulator Token", "value": "',
                    accumulatorToken,
                    '"},'
                    '{"trait_type": "Token ID", "value": "',
                    tokenId,
                    '"}'
                    "]"
                )
            );
    }

    function generateName(ConstructTokenURIParams memory params)
        private
        pure
        returns (string memory)
    {
        return
            string(
                abi.encodePacked("AliceNet Staked Token For Position #", params.tokenId.toString())
            );
    }
}
/* solhint-enable */
