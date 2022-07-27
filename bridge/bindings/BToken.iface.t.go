// Generated by ifacemaker. DO NOT EDIT.

package bindings

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// IBTokenTransactor ...
type IBTokenTransactor interface {
	// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
	//
	// Solidity: function approve(address spender, uint256 amount) returns(bool)
	Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error)
	// Burn is a paid mutator transaction binding the contract method 0xb390c0ab.
	//
	// Solidity: function burn(uint256 amount_, uint256 minEth_) returns(uint256 numEth)
	Burn(opts *bind.TransactOpts, amount_ *big.Int, minEth_ *big.Int) (*types.Transaction, error)
	// BurnTo is a paid mutator transaction binding the contract method 0x9b057203.
	//
	// Solidity: function burnTo(address to_, uint256 amount_, uint256 minEth_) returns(uint256 numEth)
	BurnTo(opts *bind.TransactOpts, to_ common.Address, amount_ *big.Int, minEth_ *big.Int) (*types.Transaction, error)
	// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
	//
	// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
	DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error)
	// Deposit is a paid mutator transaction binding the contract method 0x00838172.
	//
	// Solidity: function deposit(uint8 accountType_, address to_, uint256 amount_) returns(uint256)
	Deposit(opts *bind.TransactOpts, accountType_ uint8, to_ common.Address, amount_ *big.Int) (*types.Transaction, error)
	// DepositTokensOnBridges is a paid mutator transaction binding the contract method 0x1a116e27.
	//
	// Solidity: function depositTokensOnBridges(uint256 maxEth, uint256 maxTokens, uint16 bridgeVersion, bytes data) payable returns()
	DepositTokensOnBridges(opts *bind.TransactOpts, maxEth *big.Int, maxTokens *big.Int, bridgeVersion uint16, data []byte) (*types.Transaction, error)
	// DestroyBTokens is a paid mutator transaction binding the contract method 0x2dc6b024.
	//
	// Solidity: function destroyBTokens(uint256 numBTK_) returns(bool)
	DestroyBTokens(opts *bind.TransactOpts, numBTK_ *big.Int) (*types.Transaction, error)
	// DestroyPreApprovedBTokens is a paid mutator transaction binding the contract method 0xdfda26fd.
	//
	// Solidity: function destroyPreApprovedBTokens(address account, uint256 numBTK_) returns(bool)
	DestroyPreApprovedBTokens(opts *bind.TransactOpts, account common.Address, numBTK_ *big.Int) (*types.Transaction, error)
	// Distribute is a paid mutator transaction binding the contract method 0xe4fc6b6d.
	//
	// Solidity: function distribute() returns(uint256 minerAmount, uint256 stakingAmount, uint256 lpStakingAmount, uint256 foundationAmount)
	Distribute(opts *bind.TransactOpts) (*types.Transaction, error)
	// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
	//
	// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
	IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error)
	// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
	//
	// Solidity: function initialize() returns()
	Initialize(opts *bind.TransactOpts) (*types.Transaction, error)
	// Mint is a paid mutator transaction binding the contract method 0xa0712d68.
	//
	// Solidity: function mint(uint256 minBTK_) payable returns(uint256 numBTK)
	Mint(opts *bind.TransactOpts, minBTK_ *big.Int) (*types.Transaction, error)
	// MintDeposit is a paid mutator transaction binding the contract method 0x4f232628.
	//
	// Solidity: function mintDeposit(uint8 accountType_, address to_, uint256 minBTK_) payable returns(uint256)
	MintDeposit(opts *bind.TransactOpts, accountType_ uint8, to_ common.Address, minBTK_ *big.Int) (*types.Transaction, error)
	// MintTo is a paid mutator transaction binding the contract method 0x449a52f8.
	//
	// Solidity: function mintTo(address to_, uint256 minBTK_) payable returns(uint256 numBTK)
	MintTo(opts *bind.TransactOpts, to_ common.Address, minBTK_ *big.Int) (*types.Transaction, error)
	// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
	//
	// Solidity: function setAdmin(address admin_) returns()
	SetAdmin(opts *bind.TransactOpts, admin_ common.Address) (*types.Transaction, error)
	// SetSplits is a paid mutator transaction binding the contract method 0x767bc1bf.
	//
	// Solidity: function setSplits(uint256 validatorStakingSplit_, uint256 publicStakingSplit_, uint256 liquidityProviderStakingSplit_, uint256 protocolFee_) returns()
	SetSplits(opts *bind.TransactOpts, validatorStakingSplit_ *big.Int, publicStakingSplit_ *big.Int, liquidityProviderStakingSplit_ *big.Int, protocolFee_ *big.Int) (*types.Transaction, error)
	// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
	//
	// Solidity: function transfer(address to, uint256 amount) returns(bool)
	Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error)
	// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
	//
	// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
	TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error)
	// VirtualMintDeposit is a paid mutator transaction binding the contract method 0x92178278.
	//
	// Solidity: function virtualMintDeposit(uint8 accountType_, address to_, uint256 amount_) returns(uint256)
	VirtualMintDeposit(opts *bind.TransactOpts, accountType_ uint8, to_ common.Address, amount_ *big.Int) (*types.Transaction, error)
}
