import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { expect } from "chai";
import { BigNumber, BigNumberish, Contract } from "ethers";
import {
  callFunctionAndGetReturnValues,
  getTokenIdFromTx,
  mineBlocks,
} from "../setup";

export interface Position {
  shares: bigint;
  freeAfter: bigint;
  withdrawFreeAfter: bigint;
  accumulatorEth: bigint;
  accumulatorToken: bigint;
}

export const newPosition = (
  shares: bigint,
  freeAfter: bigint,
  withdrawFreeAfter: bigint,
  accumulatorEth: bigint,
  accumulatorToken: bigint
): Position => {
  return {
    shares,
    freeAfter,
    withdrawFreeAfter,
    accumulatorEth,
    accumulatorToken,
  };
};

export const getPosition = async (
  contract: Contract,
  tokenID: BigNumberish
): Promise<Position> => {
  const [
    _shares,
    _freeAfter,
    _withdrawFreeAfter,
    _accumulatorEth,
    _accumulatorToken,
  ] = await contract.getPosition(tokenID);
  return {
    shares: _shares.toBigInt(),
    freeAfter: _freeAfter.toBigInt(),
    withdrawFreeAfter: _withdrawFreeAfter.toBigInt(),
    accumulatorEth: _accumulatorEth.toBigInt(),
    accumulatorToken: _accumulatorToken.toBigInt(),
  };
};

export const assertTotalReserveAndZeroExcess = async (
  contract: Contract,
  reserveAmountToken: bigint,
  reserveAmountEth: bigint
) => {
  expect((await contract.getTotalReserveMadToken()).toBigInt()).to.be.equals(
    reserveAmountToken,
    "Total reserve tokens don't match expected value!"
  );
  expect((await contract.getTotalReserveEth()).toBigInt()).to.be.equals(
    reserveAmountEth,
    "Total reserve eth don't match expected value!"
  );
  expect((await contract.estimateExcessToken()).toBigInt()).to.be.equals(
    BigInt(0),
    "Excess token is not 0!"
  );
  expect((await contract.estimateExcessEth()).toBigInt()).to.be.equals(
    BigInt(0),
    "Excess eth is not 0!"
  );
};

export const assertPositions = async (
  contract: Contract,
  tokenID: number,
  expectedPosition: Position,
  ownerAddress?: string,
  expectedBalance?: bigint,
  reserveAmountToken?: bigint,
  reserveAmountEth?: bigint
) => {
  expect(await getPosition(contract, tokenID)).to.be.deep.equals(
    expectedPosition,
    "Actual position didn't match expected position!"
  );
  if (ownerAddress !== undefined) {
    expect((await contract.ownerOf(tokenID)).toLowerCase()).to.be.equals(
      ownerAddress.toLowerCase(),
      "Owner address didn't match expected address!"
    );
    if (expectedBalance !== undefined) {
      expect((await contract.balanceOf(ownerAddress)).toBigInt()).to.be.equals(
        expectedBalance,
        "Balance didn't match the expected amount!"
      );
    }
  }
  if (reserveAmountToken !== undefined && reserveAmountEth !== undefined) {
    await assertTotalReserveAndZeroExcess(
      contract,
      reserveAmountToken,
      reserveAmountEth
    );
  }
};

export const assertERC20Balance = async (
  contract: Contract,
  ownerAddress: string,
  expectedAmount: bigint
) => {
  expect((await contract.balanceOf(ownerAddress)).toBigInt()).to.be.equals(
    expectedAmount,
    "ERC20 Balance didn't match the expected amount!"
  );
};

export const collectEth = async (
  contract: Contract,
  user: SignerWithAddress,
  tokenID: number
): Promise<BigNumber> => {
  const [collectedEth] = await callFunctionAndGetReturnValues(
    contract,
    "collectEth",
    user,
    [tokenID]
  );
  return collectedEth;
};

export const collectToken = async (
  contract: Contract,
  user: SignerWithAddress,
  tokenID: number
): Promise<BigNumber> => {
  const [collectedEth] = await callFunctionAndGetReturnValues(
    contract,
    "collectToken",
    user,
    [tokenID]
  );
  return collectedEth;
};

export const mintPosition = async (
  contract: Contract,
  user: SignerWithAddress,
  amountToMint: bigint,
  expectedAccumulatorEth: bigint = 0n,
  expectedAccumulatorToken: bigint = 0n
): Promise<[number, Position]> => {
  const tx = await contract.connect(user).mint(amountToMint);
  await mineBlocks(1n);
  const tokenID = await getTokenIdFromTx(tx);
  const blockNumber = BigInt((await tx.wait()).blockNumber);
  const expectedPosition = newPosition(
    amountToMint,
    blockNumber + 1n,
    blockNumber + 1n,
    expectedAccumulatorEth,
    expectedAccumulatorToken
  );
  await assertPositions(contract, tokenID, expectedPosition, user.address);

  return [tokenID, expectedPosition];
};

export const assertAccumulatorAndSlushEth = async (
  contract: Contract,
  accumulatorExpectedAmount: bigint,
  slushExpectedAmount: bigint
) => {
  const [accumulator, slush] = await contract.getEthAccumulator();
  expect(accumulator.toBigInt()).to.be.equals(
    accumulatorExpectedAmount,
    "Accumulator didn't match the expected amount!"
  );
  expect(slush.toBigInt()).to.be.equals(
    slushExpectedAmount,
    "Slush didn't match the expected amount!"
  );
};

export const assertAccumulatorAndSlushToken = async (
  contract: Contract,
  accumulatorExpectedAmount: bigint,
  slushExpectedAmount: bigint
) => {
  const [accumulator, slush] = await contract.getTokenAccumulator();
  expect(accumulator.toBigInt()).to.be.equals(
    accumulatorExpectedAmount,
    "Accumulator didn't match the expected amount!"
  );
  expect(slush.toBigInt()).to.be.equals(
    slushExpectedAmount,
    "Slush didn't match the expected amount!"
  );
};

export const estimateAndCollectTokens = async (
  contract: Contract,
  tokenID: number,
  user: SignerWithAddress,
  expectedCollectAmount: bigint
) => {
  expect(
    (await contract.estimateTokenCollection(tokenID)).toBigInt()
  ).to.be.equals(
    expectedCollectAmount,
    "Token collection don't match expected value!"
  );

  expect((await collectToken(contract, user, tokenID)).toBigInt()).to.be.equals(
    expectedCollectAmount,
    "Collected amount to match expected value!"
  );
};
