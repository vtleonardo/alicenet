import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import {
  BigNumber,
  ContractReceipt,
  ContractTransaction,
  Wallet,
} from "ethers/lib/ethers";
import hre, { ethers, expect } from "hardhat";
import { deployCreateAndRegister } from "../../scripts/lib/alicenetFactory";
import {
  CONTRACT_ADDR,
  EVENT_DEPLOYED_RAW,
  LOCK_UP,
} from "../../scripts/lib/constants";
import { BonusPool, Lockup, RewardPool } from "../../typechain-types";
import { getEventVar } from "../factory/Setup";
import {
  BaseTokensFixture,
  deployFactoryAndBaseTokens,
  deployUpgradeableWithFactory,
  mineBlocks,
  posFixtureSetup,
  preFixtureSetup,
} from "../setup";
import { Distribution1 } from "./test.data";

export const numberOfLockingUsers = 5;
export const stakedAmount = ethers.utils.parseEther("100000000").toBigInt();
export const totalBonusAmount = ethers.utils.parseEther("2000000");
export const originalLockedAmount = ethers.utils
  .parseEther("20000000")
  .toBigInt();
export const ENROLLMENT_PERIOD = 100;
export const lockDuration = 100;
export const LockupStates = {
  PreLock: 0,
  InLock: 1,
  PostLock: 2,
};
export const example = {
  totalStakedALCA: 100000000,
  distribution: Distribution1,
};

export const profitETH = ethers.utils.parseEther(
  example.distribution.profitETH
);
export const profitALCA = ethers.utils.parseEther(
  example.distribution.profitALCA
);

export interface Fixture extends BaseTokensFixture {
  lockup: Lockup;
  rewardPool: RewardPool;
  bonusPool: BonusPool;
}
interface UserDistribution {
  owner: string;
  tokenID: number;
  shares: string;
  percentageFromTotal: number;
  profitETH: string;
  profitALCA: string;
}

interface UsersDistribution {
  [key: string]: UserDistribution;
}
export interface Distribution {
  profitETH: string;
  profitALCA: string;
  users: UsersDistribution;
}

export interface UnlockState {
  [key: string]: {
    tokenID: number;
    shares: number;
    percentageFromTotalLockup: number;
    bonusSharesALCA: number;
    bonusProfitEth: number;
    bonusProfitALCA: number;
    rewardHeldEth: number;
    rewardHeldALCA: number;
    lockupAggregatedEth: number;
    lockupAggregatedALCA: number;
    totalEarnedEth: string;
    totalEarnedALCA: string;
  };
}

interface UsersState {
  [key: string]: {
    address: string;
    alca: bigint;
    eth: bigint;
    tokenId: bigint;
    tokenOwner: string;
    position: bigint;
    rewardEth?: bigint;
    rewardToken?: bigint;
  };
}

interface LockupPositionsState {
  [key: string]: {
    index: bigint;
    owner: string;
    tokenId: bigint;
  };
}

interface StakingPositionsState {
  [key: string]: {
    shares: bigint;
    tokenId: bigint;
  };
}

interface ContractsState {
  [key: string]: {
    alca: bigint;
    eth: bigint;
    lockedPositions?: bigint;
  };
}

export interface State {
  contracts: ContractsState;
  users: UsersState;
  lockupPositions: LockupPositionsState;
  stakingPositions: StakingPositionsState;
}

export async function getState(fixture: Fixture | BaseTokensFixture) {
  const signers = await ethers.getSigners();
  const contracts = [
    fixture.lockup,
    fixture.publicStaking,
    fixture.bonusPool,
    fixture.rewardPool,
    fixture.factory,
  ];
  const contractNames = [
    "lockup",
    "publicStaking",
    "bonusPool",
    "rewardPool",
    "factory",
  ];
  const contractsState: ContractsState = {};
  const usersState: UsersState = {};
  for (let i = 0; i < contracts.length; i++) {
    if (contractNames[i] === "lockup")
      contractsState[contractNames[i]] = {
        alca: (await fixture.alca.balanceOf(contracts[i].address)).toBigInt(),
        eth: (
          await ethers.provider.getBalance(contracts[i].address)
        ).toBigInt(),
        lockedPositions: (
          await fixture.lockup.getCurrentNumberOfLockedPositions()
        ).toBigInt(),
      };
    else {
      contractsState[contractNames[i]] = {
        alca: (await fixture.alca.balanceOf(contracts[i].address)).toBigInt(),
        eth: (
          await ethers.provider.getBalance(contracts[i].address)
        ).toBigInt(),
      };
    }
  }
  for (let i = 1; i <= numberOfLockingUsers; i++) {
    const [rewardEth, rewardALCA] =
      await fixture.lockup.getTemporaryRewardBalance(signers[i].address);
    usersState["user" + i] = {
      address: signers[i].address,
      alca: (await fixture.alca.balanceOf(signers[i].address)).toBigInt(),
      eth: (await ethers.provider.getBalance(signers[i].address)).toBigInt(),
      tokenId: (await fixture.lockup.tokenOf(signers[i].address)).toBigInt(),
      tokenOwner: await fixture.lockup.ownerOf(
        await fixture.lockup.tokenOf(signers[i].address)
      ),
      position: (
        await fixture.lockup.getPositionByIndex(
          await fixture.lockup.tokenOf(signers[i].address)
        )
      ).toBigInt(),
      rewardEth: rewardEth.toBigInt(),
      rewardToken: rewardALCA.toBigInt(),
    };
  }

  usersState.bonusPool = {
    address: fixture.bonusPool.address,
    alca: (await fixture.alca.balanceOf(fixture.bonusPool.address)).toBigInt(),
    eth: (
      await ethers.provider.getBalance(fixture.bonusPool.address)
    ).toBigInt(),
    tokenId: (
      await fixture.lockup.tokenOf(fixture.bonusPool.address)
    ).toBigInt(),
    tokenOwner: await fixture.lockup.ownerOf(
      await fixture.lockup.tokenOf(fixture.bonusPool.address)
    ),
    position: (
      await fixture.lockup.getPositionByIndex(
        await fixture.lockup.tokenOf(fixture.bonusPool.address)
      )
    ).toBigInt(),
  };

  const positionsState: LockupPositionsState = {};
  const positions = await fixture.lockup.getCurrentNumberOfLockedPositions();
  for (let i = 1; i <= positions; i++) {
    const position = (await fixture.lockup.getPositionByIndex(i)).toBigInt();
    const owner = await fixture.lockup.ownerOf(position);
    const tokenId = (await fixture.lockup.tokenOf(owner)).toBigInt();
    positionsState[i] = {
      owner: await fixture.lockup.ownerOf(position),
      tokenId: (await fixture.lockup.tokenOf(owner)).toBigInt(),
      index: (await fixture.lockup.getIndexByTokenId(tokenId)).toBigInt(),
    };
  }

  const stakingsState: StakingPositionsState = {};
  for (let i = 1; i <= numberOfLockingUsers; i++) {
    const owner = signers[i].address;
    const tokenId_ = (await fixture.lockup.tokenOf(owner)).toBigInt();
    if (tokenId_ !== 0n) {
      const [positionShares, , ,] = await fixture.publicStaking.getPosition(
        tokenId_
      );
      stakingsState["user" + i] = {
        shares: positionShares.toBigInt(),
        tokenId: tokenId_,
      };
    }
  }

  const state: State = {
    contracts: contractsState,
    users: usersState,
    lockupPositions: positionsState,
    stakingPositions: stakingsState,
  };
  return state;
}

(BigInt.prototype as any).toJSON = function () {
  const dotString = this.toString().split(/(?=(?:\d{18})+(?:\.|$))/g)[0];
  const commaString = dotString
    .toString()
    .split(/(?=(?:\d{3})+(?:\.|$))/g)
    .join(",");
  return commaString;
};

export function showState(title: string, state: State) {
  if (process.env.npm_config_detailed === "true") {
    // execute "npm --detailed=true test" to see this output
    console.log(title, JSON.parse(JSON.stringify(state, null, 2)));
  }
}
export function showVariable(title: string, data: any) {
  if (process.env.npm_config_detailed === "true") {
    // execute "npm --detailed=true test" to see this output
    console.log(title, data);
  }
}
export const getEthConsumedAsGas = (receipt: ContractReceipt): bigint => {
  return receipt.cumulativeGasUsed.mul(receipt.effectiveGasPrice).toBigInt();
};

export const getImpersonatedSigner = async (
  addressToImpersonate: string
): Promise<any> => {
  const [admin] = await ethers.getSigners();
  const testUtils = await (
    await (await ethers.getContractFactory("TestUtils")).deploy()
  ).deployed();
  await admin.sendTransaction({
    to: testUtils.address,
    value: ethers.utils.parseEther("1"),
  });
  await testUtils.payUnpayable(addressToImpersonate);
  await hre.network.provider.request({
    method: "hardhat_impersonateAccount",
    params: [addressToImpersonate],
  });
  return ethers.getImpersonatedSigner(addressToImpersonate);
};

export async function deployLockupContract(
  baseTokensFixture: BaseTokensFixture,
  enrollmentPeriod: number = ENROLLMENT_PERIOD
) {
  const txResponse = await deployCreateAndRegister(
    LOCK_UP,
    baseTokensFixture.factory,
    ethers,
    [enrollmentPeriod, lockDuration, totalBonusAmount],
    ethers.utils.formatBytes32String(LOCK_UP)
  );

  // get the address from the event
  const lockupAddress = await getEventVar(
    txResponse,
    EVENT_DEPLOYED_RAW,
    CONTRACT_ADDR
  );
  const lockupStartBlock =
    (txResponse.blockNumber as number) + ENROLLMENT_PERIOD;
  await posFixtureSetup(baseTokensFixture.factory, baseTokensFixture.alca);
  return {
    lockup: await ethers.getContractAt(LOCK_UP, lockupAddress),
    lockupStartBlock,
  };
}

export async function getSimulatedStakingPositions(
  fixture: Fixture,
  signers: SignerWithAddress[],
  numberOfUsers: number,
  createBonusPosition: boolean = true
) {
  const tokenIDs = [];
  const asFactory = await getImpersonatedSigner(fixture.factory.address);
  await fixture.alca
    .connect(signers[0])
    .increaseAllowance(fixture.publicStaking.address, stakedAmount);
  await fixture.alca
    .connect(signers[0])
    .transfer(fixture.bonusPool.address, totalBonusAmount);
  for (let i = 1; i <= numberOfUsers * 10; i++) {
    if (i % 10 === 0) {
      // stake test positions only for tokens 10,20,30,40 & 50
      const index = i / 10;
      const user = ("user" + index) as string;
      const stakedAmount = ethers.utils.parseEther(
        example.distribution.users[user].shares
      );
      await fixture.publicStaking
        .connect(signers[0])
        .mintTo(signers[index].address, stakedAmount, 0);
      const tokenID = await fixture.publicStaking.tokenOfOwnerByIndex(
        signers[index].address,
        0
      );
      tokenIDs[index] = tokenID;
    } else {
      if (i % 2 === 0) {
        // for the rest stake 1M if even
        await fixture.publicStaking
          .connect(signers[0])
          .mintTo(signers[0].address, ethers.utils.parseEther("1000000"), 0);
      } else {
        // or 500K if odd
        await fixture.publicStaking
          .connect(signers[0])
          .mintTo(signers[0].address, ethers.utils.parseEther("500000"), 0);
      }
    }
  }
  createBonusPosition &&
    (await fixture.bonusPool.connect(asFactory).createBonusStakedPosition());
  const leftOver =
    stakedAmount - (await fixture.publicStaking.getTotalShares()).toBigInt();
  await fixture.publicStaking
    .connect(signers[0])
    .mintTo(signers[0].address, leftOver, 0);
  if (createBonusPosition) {
    tokenIDs[tokenIDs.length] = await fixture.publicStaking.tokenOfOwnerByIndex(
      fixture.bonusPool.address,
      0
    );
  }

  return tokenIDs;
}

export async function deployFixtureWithoutImpersonatingFactory() {
  return deployFixture(undefined, false, true);
}
export async function deployFixtureForAggregateProfits() {
  return deployFixture(1000, false, false, false, false, false);
}
export async function deployFixtureWithoutStaking() {
  return deployFixture(undefined, true, false);
}

export async function deployFixtureWithoutBonusPosition() {
  return deployFixture(undefined, true, true, false);
}

export async function deployFixture(
  enrollementPeriod: number = ENROLLMENT_PERIOD,
  impersonateLockup: boolean = true,
  simulateStakedPosition: boolean = true,
  createBonusPosition: boolean = true,
  impersonateRewardPool: boolean = true,
  impersonatePublicStaking: boolean = true
) {
  await preFixtureSetup();
  const signers = await ethers.getSigners();
  const baseTokensFixture = await deployFactoryAndBaseTokens();
  // deploying foundation so terminate doesn't fail
  await deployUpgradeableWithFactory(
    baseTokensFixture.factory,
    "Foundation",
    undefined
  );
  await ethers.provider.getBlockNumber();
  const { lockup, lockupStartBlock } = await deployLockupContract(
    baseTokensFixture,
    enrollementPeriod
  );
  // get the address of the reward pool from the lockup contract
  const rewardPoolAddress = await lockup.getRewardPoolAddress();
  const rewardPool = await ethers.getContractAt(
    "RewardPool",
    rewardPoolAddress
  );
  // get the address of the bonus pool from the reward pool contract
  const bonusPoolAddress = await rewardPool.getBonusPoolAddress();
  const bonusPool = await ethers.getContractAt("BonusPool", bonusPoolAddress);
  const factorySigner = await getImpersonatedSigner(
    baseTokensFixture.factory.address
  );
  let pblicStakingSigner;
  if (impersonatePublicStaking) {
    pblicStakingSigner = await getImpersonatedSigner(
      baseTokensFixture.publicStaking.address
    );
  }
  let rewardPoolSigner;
  if (impersonateLockup) {
    rewardPoolSigner = await getImpersonatedSigner(rewardPoolAddress);
  }
  const fixture: Fixture = {
    ...baseTokensFixture,
    rewardPool,
    lockup,
    bonusPool,
    lockupStartBlock,
  };
  let tokenIDs: BigNumber[] = [];
  if (impersonateRewardPool) {
    rewardPoolSigner = await getImpersonatedSigner(rewardPoolAddress);
  }
  if (simulateStakedPosition) {
    tokenIDs = await getSimulatedStakingPositions(
      fixture,
      signers,
      5,
      createBonusPosition
    );
    expect(
      (await fixture.publicStaking.getTotalShares()).toBigInt()
    ).to.be.equals(stakedAmount);
  }
  return {
    fixture,
    accounts: signers,
    stakedTokenIDs: tokenIDs,
    asFactory: factorySigner,
    asPublicStaking: pblicStakingSigner,
    asRewardPool: rewardPoolSigner,
  };
}

export async function distributeProfits(
  fixture: BaseTokensFixture,
  fundsSourceAddress: SignerWithAddress,
  profitETH: BigNumber,
  profitALCA: BigNumber
) {
  await fixture.alca
    .connect(fundsSourceAddress)
    .increaseAllowance(fixture.publicStaking.address, profitALCA);
  await fixture.publicStaking.connect(fundsSourceAddress).depositEth(42, {
    value: profitETH,
  });
  await fixture.publicStaking
    .connect(fundsSourceAddress)
    .depositToken(42, profitALCA);
}

export async function lockStakedNFT(
  fixture: Fixture,
  account: Wallet | SignerWithAddress,
  tokenID: BigNumber
): Promise<ContractTransaction> {
  return fixture.publicStaking
    .connect(account)
    ["safeTransferFrom(address,address,uint256,bytes)"](
      account.address,
      fixture.lockup.address,
      tokenID,
      "0x"
    );
}

export async function jumpToInlockState(fixture: BaseTokensFixture) {
  const blocksToMine = (await fixture.lockup.getLockupStartBlock())
    .sub(await ethers.provider.getBlockNumber())
    .toBigInt();
  await mineBlocks(blocksToMine + 1n);
  expect(await fixture.lockup.getState()).to.be.equals(LockupStates.InLock);
}

export async function jumpToPostLockState(fixture: BaseTokensFixture) {
  const blocksToMine = (await fixture.lockup.getLockupEndBlock())
    .sub(await ethers.provider.getBlockNumber())
    .toBigInt();
  await mineBlocks(blocksToMine + 1n);
  expect(await fixture.lockup.getState()).to.be.equals(LockupStates.PostLock);
}

export async function ensureBlockIsAtLeast(targetBlock: number): Promise<void> {
  const currentBlock = await ethers.provider.getBlockNumber();
  if (currentBlock < targetBlock) {
    const blockDelta = targetBlock - currentBlock;
    await mineBlocks(BigInt(blockDelta));
  }
}

export async function getUserLockingInfo(fixture: Fixture, userId: number) {
  const totalShares = await fixture.publicStaking.getTotalShares();
  const signers = await ethers.getSigners();
  const owner_ = signers[userId];
  const tokenId_ = await fixture.lockup.tokenOf(owner_.address);
  const index_ = await fixture.lockup.getIndexByTokenId(tokenId_);
  const userInitialShares_ = ethers.utils.parseEther(
    example.distribution.users["user" + userId].shares
  );
  const [userCurrentShares, , , ,] = await fixture.publicStaking.getPosition(
    tokenId_
  );
  const profitALCAUser_ = profitALCA
    .mul(userCurrentShares)
    .div(totalShares)
    .toBigInt();
  const [estimatePayoutEth_, estimatePayoutALCA_] =
    await fixture.lockup.estimateProfits(tokenId_);
  const profitETHUser_ = profitETH
    .mul(userCurrentShares)
    .div(totalShares)
    .toBigInt();
  const reservedProfitALCAUser_ = (
    await fixture.lockup.getReservedAmount(profitALCAUser_)
  ).toBigInt();
  const reservedProfitETHUser_ = (
    await fixture.lockup.getReservedAmount(profitETHUser_)
  ).toBigInt();
  const freeProfitALCAUser_ = profitALCAUser_ - reservedProfitALCAUser_;
  const freeProfitETHUser_ = profitETHUser_ - reservedProfitETHUser_;
  return {
    userInitialShares: userInitialShares_.toBigInt(),
    userCurrentShares: userCurrentShares.toBigInt(),
    profitALCAUser: profitALCAUser_,
    reservedProfitALCAUser: reservedProfitALCAUser_,
    profitETHUser: profitETHUser_,
    reservedProfitETHUser: reservedProfitETHUser_,
    index: index_,
    tokenId: tokenId_.toBigInt(),
    estimatePayoutEth: estimatePayoutEth_.toBigInt(),
    estimatePayoutALCA: estimatePayoutALCA_.toBigInt(),
    freeProfitALCAUser: freeProfitALCAUser_,
    freeProfitETHUser: freeProfitETHUser_,
    owner: owner_,
  };
}
