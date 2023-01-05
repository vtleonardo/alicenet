// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.16;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "contracts/utils/auth/ImmutableALCA.sol";
import "contracts/utils/auth/ImmutablePublicStaking.sol";
import "contracts/utils/auth/ImmutableFoundation.sol";
import "contracts/interfaces/IStakingNFT.sol";
import "contracts/utils/ERC20SafeTransfer.sol";
import "contracts/utils/MagicEthTransfer.sol";

contract Redistribution is
    ImmutableALCA,
    ImmutablePublicStaking,
    ImmutableFoundation,
    ERC20SafeTransfer,
    MagicEthTransfer
{
    struct accountInfo {
        uint248 balance;
        bool isPositionTaken;
    }

    event Withdrawn(address indexed user, uint256 amount);
    event TokenAlreadyTransferred();

    error NotOperator();
    error WithdrawalWindowExpired();
    error WithdrawalWindowNotExpiredYet();
    error IncorrectLength();
    error ZeroAmountNotAllowed();
    error InvalidAllowanceSum(uint256 totalAllowance, uint256 maxRedistributionAmount);
    error DistributionTokenAlreadyCreated();
    error PositionAlreadyRegisteredOrTaken();
    error InvalidDistributionAmount(uint256 amount, uint256 maxAllowed);
    error NotEnoughFundsToRedistribute(uint256 withdrawAmount, uint256 currentAmount);
    error PositionAlreadyTaken();

    /// The amount of blocks that the withdraw position will be locked against burning. This is
    /// approximately 6 months.
    uint256 public constant MAX_MINT_LOCK = 1051200;
    /// The total amount of ALCA that can be redistributed to accounts via this contract.
    uint256 public immutable maxRedistributionAmount;
    /// The block number that the withdrawal window will expire.
    uint256 public immutable expireBlock;
    /// The address of the operator of the contract. The operator will be able to register new
    /// accounts that will have rights to withdraw funds.
    address public operator;
    /// The amount from the `maxRedistributionAmount` already reserved for distribution.
    uint256 public totalAllowances;
    /// The current tokenID of the public staking position that holds the ALCA to be distributed.
    uint256 public tokenID;
    mapping(address => accountInfo) internal _accounts;

    modifier onlyOperator() {
        if (msg.sender != operator) {
            revert NotOperator();
        }
        _;
    }

    modifier notExpired() {
        if (block.number > expireBlock) {
            revert WithdrawalWindowExpired();
        }
        _;
    }

    /**
     * @notice This function is used to receive ETH from the public staking contract.
     */
    receive() external payable onlyPublicStaking {}

    constructor(
        uint256 withdrawalBlockWindow,
        uint256 maxRedistributionAmount_,
        address[] memory allowedAddresses,
        uint248[] memory allowedAmounts
    ) ImmutableFactory(msg.sender) ImmutableALCA() ImmutablePublicStaking() ImmutableFoundation() {
        if (allowedAddresses.length != allowedAmounts.length || allowedAddresses.length == 0) {
            revert IncorrectLength();
        }
        uint256 totalAllowance = 0;
        for (uint256 i = 0; i < allowedAddresses.length; i++) {
            if (allowedAddresses[i] == address(0) || _accounts[allowedAddresses[i]].balance > 0) {
                revert PositionAlreadyRegisteredOrTaken();
            }
            if (allowedAmounts[i] == 0) {
                revert ZeroAmountNotAllowed();
            }
            _accounts[allowedAddresses[i]] = accountInfo(allowedAmounts[i], false);
            totalAllowance += allowedAmounts[i];
        }
        if (totalAllowance > maxRedistributionAmount_) {
            revert InvalidAllowanceSum(totalAllowance, maxRedistributionAmount_);
        }
        maxRedistributionAmount = maxRedistributionAmount_;
        totalAllowances = totalAllowance;
        expireBlock = block.number + withdrawalBlockWindow;
    }

    /**
     * @notice Set a new operator for the contract. This function can only be called by the factory.
     * @param operator_ The new operator address.
     */
    function setOperator(address operator_) public onlyFactory {
        operator = operator_;
    }

    /**
     * @notice Creates the total staked position for the redistribution. This function can only be
     * called by the factory. This function can only be called if the withdrawal window has not expired
     * yet.
     * @dev the maxRedistributionAmount should be approved to this contract before calling this
     * function.
     */
    function createRedistributionStakedPosition() public onlyFactory notExpired {
        if (tokenID != 0) {
            revert DistributionTokenAlreadyCreated();
        }
        _safeTransferFromERC20(
            IERC20Transferable(_alcaAddress()),
            msg.sender,
            maxRedistributionAmount
        );
        // approve the staking contract to transfer the ALCA
        IERC20(_alcaAddress()).approve(_publicStakingAddress(), maxRedistributionAmount);
        tokenID = IStakingNFT(_publicStakingAddress()).mint(maxRedistributionAmount);
    }

    /**
     * @notice register an new address for a distribution amount. This function can only be called
     * by the operator. The distribution amount can not be greater that the total amount left for
     * distribution. Only one amount can be registered per address. Amount for already registered
     * addresses cannot be changed.
     * @dev This function can only be called if the withdrawal window has not expired yet.
     * @param user The address to register for distribution.
     * @param distributionAmount The amount to register for distribution.
     */
    function registerAddressForDistribution(
        address user,
        uint248 distributionAmount
    ) public onlyOperator notExpired {
        if (distributionAmount == 0) {
            revert ZeroAmountNotAllowed();
        }
        accountInfo memory account = _accounts[user];
        if (account.balance > 0 || account.isPositionTaken) {
            revert PositionAlreadyRegisteredOrTaken();
        }
        uint256 distributionLeft = _getDistributionLeft();
        if (distributionAmount > distributionLeft) {
            revert InvalidDistributionAmount(distributionAmount, distributionLeft);
        }
        _accounts[user] = accountInfo(distributionAmount, false);
        totalAllowances += distributionAmount;
    }

    /**
     *  @notice Withdraw the staked position to the user's address. It will burn the Public
     *  Staking position held by this contract and mint a new one to the user's address with the
     *  owned amount and in case there is a remainder, it will mint a new position to this contract.
     *  THE CALLER OF THIS FUNCTION MUST BE AN EOA (EXTERNAL OWNED ACCOUNT) OR PROXY WALLET THAT
     *  ACCEPTS AND HANDLE ERC721 POSITIONS. BEWARE IF THIS REQUIREMENTS ARE NOT FOLLOWED, THE
     *  POSITION CAN BE FOREVER LOST.
     *  @dev This function can only be called by the user that has the right to withdraw a staked
     *  position. This function can only be called if the withdrawal window has not expired yet.
     *  @param to The address to send the staked position to.
     */
    function withdrawStakedPosition(address to) public notExpired {
        accountInfo memory account = _accounts[msg.sender];
        if (account.balance == 0 || account.isPositionTaken) {
            revert PositionAlreadyTaken();
        }
        _accounts[msg.sender] = accountInfo(0, true);
        IStakingNFT staking = IStakingNFT(_publicStakingAddress());
        IERC20 alca = IERC20(_alcaAddress());
        staking.burn(tokenID);
        uint256 alcaBalance = alca.balanceOf(address(this));
        if (alcaBalance < account.balance) {
            revert NotEnoughFundsToRedistribute(alcaBalance, account.balance);
        }
        alca.approve(_publicStakingAddress(), alcaBalance);
        staking.mintTo(to, account.balance, MAX_MINT_LOCK);
        uint256 remainder = alcaBalance - account.balance;
        if (remainder > 0) {
            tokenID = staking.mint(remainder);
        }
        // send any eth balance collected to the foundation
        uint256 ethBalance = address(this).balance;
        if (ethBalance > 0) {
            _safeTransferEthWithMagic(IMagicEthTransfer(_foundationAddress()), ethBalance);
        }
        emit Withdrawn(msg.sender, account.balance);
    }

    /**
     *  @notice Send any remaining funds that were not claimed during the valid time back to the
     *  factory. It will transfer the Public Staking position (in case it exists) and any ALCA back
     *  to the Factory. Ether will be send to the foundation.
     *  @dev This function can only be called by the AliceNet factory. This function never fails and
     *  can act as a skim of ether and ALCA.
     *  function never fails.
     */
    function sendExpiredFundsToFactory() public onlyFactory {
        if (block.number <= expireBlock) {
            revert WithdrawalWindowNotExpiredYet();
        }
        try
            IERC721(_publicStakingAddress()).transferFrom(address(this), _factoryAddress(), tokenID)
        {} catch {
            emit TokenAlreadyTransferred();
        }
        uint256 alcaBalance = IERC20(_alcaAddress()).balanceOf(address(this));
        if (alcaBalance > 0) {
            _safeTransferERC20(IERC20Transferable(_alcaAddress()), _factoryAddress(), alcaBalance);
        }
        uint256 ethBalance = address(this).balance;
        if (ethBalance > 0) {
            _safeTransferEthWithMagic(IMagicEthTransfer(_foundationAddress()), ethBalance);
        }
    }

    /**
     * @notice Returns the account info for a given user
     * @param user The address of the user
     */
    function getRedistributionInfo(address user) public view returns (accountInfo memory account) {
        account = _accounts[user];
    }

    /**
     * @notice Returns the amount of ALCA left to distribute
     */
    function getDistributionLeft() public view returns (uint256) {
        return _getDistributionLeft();
    }

    // internal function to get the amount of ALCA left to distribute
    function _getDistributionLeft() internal view returns (uint256) {
        return maxRedistributionAmount - totalAllowances;
    }
}
