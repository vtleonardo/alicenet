// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.16;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "contracts/utils/auth/ImmutableALCA.sol";
import "contracts/utils/auth/ImmutablePublicStaking.sol";

contract Redistribution is ImmutableALCA, ImmutablePublicStaking {
    uint256 internal constant _MAX_MINT_LOCK = 1051200;
    uint256 public immutable maxRedistributionAmount;
    // todo
    address public immutable operator;
    uint256 public immutable expireBlock;
    uint256 internal _tokenID;
    uint256 public totalAllowances;
    mapping(address => uint256) internal _withdrawalBalance;
    mapping(address => bool) internal _takenPositions;

    modifier onlyOperator() {
        require(msg.sender == operator, "Redistribution: not operator");
        _;
    }

    modifier notExpired() {
        require(block.number < expireBlock, "Redistribution: expired");
        _;
    }

    constructor(
        uint256 withdrawalBlockWindow,
        uint256 maxRedistributionAmount_,
        address[] memory allowedAddresses,
        uint256[] memory allowedAmounts
    ) ImmutableFactory(msg.sender) ImmutableALCA() ImmutablePublicStaking() {
        require(allowedAddresses.length == allowedAmounts.length, "Redistribution: invalid input");
        uint256 totalAllowance = 0;
        for (uint256 i = 0; i < allowedAddresses.length; i++) {
            _withdrawalBalance[allowedAddresses[i]] = allowedAmounts[i];
            totalAllowance += allowedAmounts[i];
        }
        require(totalAllowance <= maxRedistributionAmount_, "Redistribution: invalid input");
        maxRedistributionAmount = maxRedistributionAmount_;
        totalAllowances = totalAllowance;
        expireBlock = block.number + withdrawalBlockWindow;
    }

    function setOperator(address operator_) public onlyFactory {
        operator = operator_;
    }

    function createRedistributionStakedPosition() public onlyFactory {
        require(_tokenID == 0, "Redistribution: already created");
        IERC20 alca = IERC20(_alcaAddress());
        alca.transferfrom(msg.sender, address(this), maxRedistributionAmount);
        // approve the staking contract to transfer the ALCA
        alca.approve(_publicStakingAddress(), maxRedistributionAmount);
        uint256 tokenID = IStakingNFT(_publicStakingAddress()).mint(maxRedistributionAmount);
        _tokenID = tokenID;
    }

    function registerAddress(address account, uint256 amount) public onlyOperator notExpired {
        require(
            _withdrawalBalance[account] == 0 && !_takenPositions[account],
            "Redistribution: already registered or already taken position"
        );
        require(
            totalAllowances + amount <= maxRedistributionAmount,
            "Redistribution: not enough funds to distribute"
        );
        _withdrawalBalance[account] = amount;
        totalAllowances += amount;
    }

    function withdrawStakedPosition(address to) public notExpired {
        uint256 withdrawalAmount = _withdrawalBalance[msg.sender];
        require(
            withdrawalAmount > 0 && !_takenPositions[msg.sender],
            "Redistribution: already taken position or not registered"
        );
        _withdrawalBalance[msg.sender] = 0;
        _takenPositions[msg.sender] = true;
        IStakingNFT(_publicStakingAddress()).burn(_tokenID);
        uint256 alcaBalance = IERC20(_alcaAddress()).balanceOf(address(this));
        require(alcaBalance >= withdrawalAmount, "Redistribution: not enough funds to distribute");
        IERC20(_alcaAddress()).approve(_publicStakingAddress(), alcaBalance);
        uint256 tokenID = IStakingNFT(_publicStakingAddress()).mintTo(
            to,
            withdrawalAmount,
            _MAX_MINT_LOCK
        );
        uint256 remainder = alcaBalance - withdrawalAmount;
        if (remainder > 0) {
            _tokenID = IStakingNFT(_publicStakingAddress()).mint(remainder);
        }
    }

    function getTokenID() public view returns (uint256) {
        return _tokenID;
    }

    function getRedistributionAmount(address account) public view returns (uint256) {
        return _withdrawalBalance[account];
    }
}
