import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { BigNumber } from "ethers";
import { ethers, expect } from "hardhat";
import { deployCreate } from "../../scripts/lib/alicenetFactory";
import { Redistribution } from "../../typechain-types";
import {
  Fixture,
  getContractAddressFromDeployedRawEvent,
  getFixture,
  mineBlocks,
} from "../setup";

const DEFAULT_WITHDRAWAL_BLOCK_WINDOW = 172800; // 1 month block window
const DEFAULT_MAX_DISTRIBUTION_AMOUNT = ethers.utils.parseEther("10000000");

/**
 * Constructor:
 * -Should not allow sum of allowedAmounts to be greater than the maxRedistributionAmount
 * -Should not allow addresses and amounts with different length
 * -Should not allow length zero in any of the parameters
 * -Should not allow duplicate addresses and address zero
 * -Should not allow distribution amount 0
 * -Check totalAllowances must be equal to the sum of allowedAmounts
 * -Check expireBlock must be equal to block.number + withdrawalBlockWindow
 *
 * SetOperator:
 * -Should not allow call from address other than factory
 * -Only factory can call
 * -Should be able to set a operator as factory (operator is public we can check with the getter)
 *
 * CreateRedistributionStakedPosition:
 * -Should not allow call from address other than factory
 * -Only factory can call
 * -Cannot be called after expiration
 * -Should fail if there was not alca approval from the factory
 * -Should create a position if all the conditions are met, check tokenID and the position properties match the maxRedistributionAmount
 * -Should not allow to create a position if it already exists
 *
 * registerAddressForDistribution
 * Should not allow call from address other than operator (not even factory)
 * Only operator can call
 * Cannot be called after expiration
 * Should not allow to register an address that is already registered
 * Should not allow
 */

interface FixtureWithRedistribution extends Fixture {
  redistribution: Redistribution;
  accounts: SignerWithAddress[];
}

async function deployFunc(): Promise<FixtureWithRedistribution> {
  const baseFixture = await getFixture();
  const allAccounts = await ethers.getSigners();
  const distributionAccounts = allAccounts.slice(0, 5).map((a) => a.address);
  const accountAmounts = distributionAccounts.map(() =>
    // ethers.utils.parseEther(`${i + 1}`).mul(500_000)
    ethers.utils.parseEther(`500000`)
  );
  return await deployRedistribution(
    baseFixture,
    allAccounts,
    distributionAccounts,
    accountAmounts,
    DEFAULT_WITHDRAWAL_BLOCK_WINDOW,
    DEFAULT_MAX_DISTRIBUTION_AMOUNT
  );
}

async function deployRedistribution(
  fixture: Fixture,
  accounts: SignerWithAddress[],
  distributionAccounts: string[],
  accountAmounts: BigNumber[],
  withdrawalBlockWindow: number,
  maxDistributionAmount: BigNumber
): Promise<FixtureWithRedistribution> {
  const tx = await deployCreate("Redistribution", fixture.factory, ethers, [
    withdrawalBlockWindow,
    maxDistributionAmount,
    distributionAccounts,
    accountAmounts,
  ]);
  const redistributionAddress = await getContractAddressFromDeployedRawEvent(
    tx
  );
  const redistribution = await ethers.getContractAt(
    "Redistribution",
    redistributionAddress
  );
  return {
    ...fixture,
    accounts,
    redistribution,
  };
}

describe("CT redistribution", async () => {
  /**
   * Constructor testing
   */
  describe("Constructor testing", async () => {
    it("Should have correct maxRedistributionAmount", async () => {
      const deployFunc = async (): Promise<FixtureWithRedistribution> => {
        const baseFixture = await getFixture();
        const allAccounts = await ethers.getSigners();
        const distributionAccounts = allAccounts
          .slice(0, 5)
          .map((a) => a.address);
        const accountAmounts = distributionAccounts.map(() =>
          ethers.utils.parseEther(`500000`)
        );
        return await deployRedistribution(
          baseFixture,
          allAccounts,
          distributionAccounts,
          accountAmounts,
          DEFAULT_WITHDRAWAL_BLOCK_WINDOW,
          DEFAULT_MAX_DISTRIBUTION_AMOUNT
        );
      };
      const fixture = await loadFixture(deployFunc);
      expect(
        (await fixture.redistribution.maxRedistributionAmount()).toString()
      ).to.be.equal(DEFAULT_MAX_DISTRIBUTION_AMOUNT.toString());
    });

    it("Should not allow sum of allowedAmounts to be greater than the maxRedistributionAmount", async () => {
      const baseFixture = await getFixture();
      const deployFunc = async (): Promise<FixtureWithRedistribution> => {
        const allAccounts = await ethers.getSigners();
        const distributionAccounts = allAccounts
          .slice(0, 5)
          .map((a) => a.address);
        const accountAmounts = distributionAccounts.map(() =>
          // this is 5_000_000 for each of the 5 accounts, which is higher than the 10_000_000 max
          ethers.utils.parseEther(`5000000`)
        );
        return await deployRedistribution(
          baseFixture,
          allAccounts,
          distributionAccounts,
          accountAmounts,
          DEFAULT_WITHDRAWAL_BLOCK_WINDOW,
          DEFAULT_MAX_DISTRIBUTION_AMOUNT
        );
      };

      await expect(loadFixture(deployFunc)).to.be.revertedWithCustomError(
        baseFixture.factory,
        "CodeSizeZero"
      );
    });

    it("Should not allow addresses and amounts with different length", async () => {
      const baseFixture = await getFixture();
      const deployFunc = async (): Promise<FixtureWithRedistribution> => {
        const allAccounts = await ethers.getSigners();
        const distributionAccounts = allAccounts
          .slice(0, 5)
          .map((a) => a.address);
        const accountAmounts = distributionAccounts.map(() =>
          ethers.utils.parseEther(`500000`)
        );
        return await deployRedistribution(
          baseFixture,
          allAccounts,
          distributionAccounts, // passing 5
          accountAmounts.slice(0, 4), // passing 4
          DEFAULT_WITHDRAWAL_BLOCK_WINDOW,
          DEFAULT_MAX_DISTRIBUTION_AMOUNT
        );
      };
      await expect(loadFixture(deployFunc)).to.be.revertedWithCustomError(
        baseFixture.factory,
        "CodeSizeZero"
      );
    });

    it("Should not allow length zero in any of the parameters", async () => {
      const baseFixture = await getFixture();
      const deployFunc = async (): Promise<FixtureWithRedistribution> => {
        const allAccounts = await ethers.getSigners();
        return await deployRedistribution(
          baseFixture,
          allAccounts,
          [], // empty
          [], // empty
          DEFAULT_WITHDRAWAL_BLOCK_WINDOW,
          DEFAULT_MAX_DISTRIBUTION_AMOUNT
        );
      };
      await expect(loadFixture(deployFunc)).to.be.revertedWithCustomError(
        baseFixture.factory,
        "CodeSizeZero"
      );
    });

    it("Should not allow duplicate addresses", async () => {
      const baseFixture = await getFixture();
      const deployFunc = async (): Promise<FixtureWithRedistribution> => {
        const allAccounts = await ethers.getSigners();
        const distributionAccounts = allAccounts
          .slice(0, 4)
          .map((a) => a.address);
        distributionAccounts.push(distributionAccounts[0]); // duplicated
        const accountAmounts = distributionAccounts.map(() =>
          ethers.utils.parseEther(`500000`)
        );
        return await deployRedistribution(
          baseFixture,
          allAccounts,
          distributionAccounts,
          accountAmounts,
          DEFAULT_WITHDRAWAL_BLOCK_WINDOW,
          DEFAULT_MAX_DISTRIBUTION_AMOUNT
        );
      };
      await expect(loadFixture(deployFunc)).to.be.revertedWithCustomError(
        baseFixture.factory,
        "CodeSizeZero"
      );
    });

    it("Should not allow address zero", async () => {
      const baseFixture = await getFixture();
      const deployFunc = async (): Promise<FixtureWithRedistribution> => {
        const allAccounts = await ethers.getSigners();
        const distributionAccounts = allAccounts
          .slice(0, 5)
          .map((a) => a.address);
        // address 0
        distributionAccounts[0] = "0x0000000000000000000000000000000000000000";
        const accountAmounts = distributionAccounts.map(() =>
          ethers.utils.parseEther(`500000`)
        );
        return await deployRedistribution(
          baseFixture,
          allAccounts,
          distributionAccounts,
          accountAmounts,
          DEFAULT_WITHDRAWAL_BLOCK_WINDOW,
          DEFAULT_MAX_DISTRIBUTION_AMOUNT
        );
      };
      await expect(loadFixture(deployFunc)).to.be.revertedWithCustomError(
        baseFixture.factory,
        "CodeSizeZero"
      );
    });

    it("Should not allow distribution amount 0", async () => {
      const baseFixture = await getFixture();
      const deployFunc = async (): Promise<FixtureWithRedistribution> => {
        const allAccounts = await ethers.getSigners();
        const distributionAccounts = allAccounts
          .slice(0, 5)
          .map((a) => a.address);
        const accountAmounts = distributionAccounts.map(
          () => ethers.utils.parseEther(`0`) // amount 0
        );
        return await deployRedistribution(
          baseFixture,
          allAccounts,
          distributionAccounts,
          accountAmounts,
          DEFAULT_WITHDRAWAL_BLOCK_WINDOW,
          DEFAULT_MAX_DISTRIBUTION_AMOUNT
        );
      };
      await expect(loadFixture(deployFunc)).to.be.revertedWithCustomError(
        baseFixture.factory,
        "CodeSizeZero"
      );
    });

    it("totalAllowances must be equal to the sum of allowedAmounts", async () => {
      const baseFixture = await getFixture();
      const deployFunc = async (): Promise<FixtureWithRedistribution> => {
        const allAccounts = await ethers.getSigners();
        const distributionAccounts = allAccounts
          .slice(0, 5)
          .map((a) => a.address);
        const accountAmounts = distributionAccounts.map(() =>
          ethers.utils.parseEther(`500000`)
        );
        return await deployRedistribution(
          baseFixture,
          allAccounts,
          distributionAccounts,
          accountAmounts,
          DEFAULT_WITHDRAWAL_BLOCK_WINDOW,
          DEFAULT_MAX_DISTRIBUTION_AMOUNT
        );
      };
      const fixture = await loadFixture(deployFunc);
      const totalAmounts = ethers.utils.parseEther(`500000`).mul(5);
      expect(
        (await fixture.redistribution.totalAllowances()).toString()
      ).to.be.equal(totalAmounts.toString());
    });

    it("expireBlock must be equal to block.number + withdrawalBlockWindow", async () => {
      const baseFixture = await getFixture();
      const deployFunc = async (): Promise<FixtureWithRedistribution> => {
        const allAccounts = await ethers.getSigners();
        const distributionAccounts = allAccounts
          .slice(0, 5)
          .map((a) => a.address);
        const accountAmounts = distributionAccounts.map(() =>
          ethers.utils.parseEther(`500000`)
        );
        return await deployRedistribution(
          baseFixture,
          allAccounts,
          distributionAccounts,
          accountAmounts,
          DEFAULT_WITHDRAWAL_BLOCK_WINDOW,
          DEFAULT_MAX_DISTRIBUTION_AMOUNT
        );
      };
      const fixture = await loadFixture(deployFunc);
      const latestBlock = await ethers.provider.getBlock("latest");
      const expectedExpireBlock =
        latestBlock.number + DEFAULT_WITHDRAWAL_BLOCK_WINDOW;
      expect(
        (await fixture.redistribution.expireBlock()).toString()
      ).to.be.equal(expectedExpireBlock.toString());
    });
  });

  /**
   * SetOperator testing
   */

  describe("SetOperator testing", async () => {
    let fixture: FixtureWithRedistribution;
    beforeEach(async () => {
      const deployFunc = async (): Promise<FixtureWithRedistribution> => {
        const baseFixture = await getFixture();
        const allAccounts = await ethers.getSigners();
        const distributionAccounts = allAccounts
          .slice(0, 5)
          .map((a) => a.address);
        const accountAmounts = distributionAccounts.map(() =>
          // ethers.utils.parseEther(`${i + 1}`).mul(500_000)
          ethers.utils.parseEther(`500000`)
        );
        return await deployRedistribution(
          baseFixture,
          allAccounts,
          distributionAccounts,
          accountAmounts,
          DEFAULT_WITHDRAWAL_BLOCK_WINDOW,
          DEFAULT_MAX_DISTRIBUTION_AMOUNT
        );
      };
      fixture = await loadFixture(deployFunc);
    });

    it("Should not allow call from address other than factory", async () => {
      await expect(
        fixture.redistribution.setOperator(fixture.accounts[0].address)
      ).to.be.revertedWithCustomError(fixture.redistribution, "OnlyFactory");
    });

    it("Should be able to set a operator as factory", async () => {
      await fixture.factory.callAny(
        fixture.redistribution.address,
        0,
        fixture.redistribution.interface.encodeFunctionData("setOperator", [
          fixture.accounts[0].address,
        ])
      );
      expect(await fixture.redistribution.operator()).to.be.equal(
        fixture.accounts[0].address
      );
    });
  });

  /**
   * CreateRedistributionStakedPosition testing
   */

  describe("CreateRedistributionStakedPosition testing", async () => {
    let fixture: FixtureWithRedistribution;
    beforeEach(async () => {
      fixture = await loadFixture(deployFunc);
    });

    it("Should not allow call from address other than factory", async () => {
      await expect(
        fixture.redistribution.createRedistributionStakedPosition()
      ).to.be.revertedWithCustomError(fixture.redistribution, "OnlyFactory");
    });

    it("Only factory can call", async () => {
      await expect(
        fixture.factory.callAny(
          fixture.alca.address,
          0,
          fixture.alca.interface.encodeFunctionData("approve", [
            fixture.redistribution.address,
            DEFAULT_MAX_DISTRIBUTION_AMOUNT,
          ])
        )
      ).to.be.fulfilled;
      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData(
            "createRedistributionStakedPosition"
          )
        )
      ).to.be.fulfilled;
    });

    it("Cannot be called after expiration", async () => {
      await expect(
        fixture.factory.callAny(
          fixture.alca.address,
          0,
          fixture.alca.interface.encodeFunctionData("approve", [
            fixture.redistribution.address,
            DEFAULT_MAX_DISTRIBUTION_AMOUNT,
          ])
        )
      ).to.be.fulfilled;

      // lets make it expire
      await mineBlocks(
        BigNumber.from(DEFAULT_WITHDRAWAL_BLOCK_WINDOW).toBigInt()
      );

      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData(
            "createRedistributionStakedPosition"
          )
        )
      ).to.be.revertedWithCustomError(
        fixture.redistribution,
        "WithdrawalWindowExpired"
      );
    });

    it("Should fail if there was not alca approval from the factory", async () => {
      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData(
            "createRedistributionStakedPosition"
          )
        )
      ).to.be.revertedWith("ERC20: insufficient allowance");
    });

    it("Should create a position if all the conditions are met", async () => {
      await expect(
        fixture.factory.callAny(
          fixture.alca.address,
          0,
          fixture.alca.interface.encodeFunctionData("approve", [
            fixture.redistribution.address,
            DEFAULT_MAX_DISTRIBUTION_AMOUNT,
          ])
        )
      ).to.be.fulfilled;
      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData(
            "createRedistributionStakedPosition"
          )
        )
      ).to.be.fulfilled;

      // check tokenID
      const tokenID = await fixture.redistribution.tokenID();
      expect(tokenID.toString()).to.be.not.equal("0");
      expect(await fixture.publicStaking.ownerOf(tokenID)).to.be.equal(
        fixture.redistribution.address
      );

      // position matches the maxRedistributionAmount
      const position = await fixture.publicStaking.getPosition(tokenID);
      expect(position.shares.toString()).to.be.equal(
        DEFAULT_MAX_DISTRIBUTION_AMOUNT.toString()
      );
    });

    it("Should not allow to create a position if it already exists", async () => {
      await expect(
        fixture.factory.callAny(
          fixture.alca.address,
          0,
          fixture.alca.interface.encodeFunctionData("approve", [
            fixture.redistribution.address,
            DEFAULT_MAX_DISTRIBUTION_AMOUNT,
          ])
        )
      ).to.be.fulfilled;
      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData(
            "createRedistributionStakedPosition"
          )
        )
      ).to.be.fulfilled;

      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData(
            "createRedistributionStakedPosition"
          )
        )
      ).to.be.revertedWithCustomError(
        fixture.redistribution,
        "DistributionTokenAlreadyCreated"
      );
    });
  });

  /**
   * RegisterAddressForDistribution testing
   */

  describe("RegisterAddressForDistribution testing", async () => {
    let fixture: FixtureWithRedistribution;
    beforeEach(async () => {
      const deployFunc = async (): Promise<FixtureWithRedistribution> => {
        const baseFixture = await getFixture();
        const allAccounts = await ethers.getSigners();
        const distributionAccounts = allAccounts
          .slice(1, 3)
          .map((a) => a.address);
        const accountAmounts = distributionAccounts.map(() =>
          ethers.utils.parseEther(`500000`)
        );
        return await deployRedistribution(
          baseFixture,
          allAccounts,
          distributionAccounts,
          accountAmounts,
          DEFAULT_WITHDRAWAL_BLOCK_WINDOW,
          DEFAULT_MAX_DISTRIBUTION_AMOUNT
        );
      };
      fixture = await loadFixture(deployFunc);

      await expect(
        fixture.factory.callAny(
          fixture.alca.address,
          0,
          fixture.alca.interface.encodeFunctionData("approve", [
            fixture.redistribution.address,
            DEFAULT_MAX_DISTRIBUTION_AMOUNT,
          ])
        )
      ).to.be.fulfilled;
      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData(
            "createRedistributionStakedPosition"
          )
        )
      ).to.be.fulfilled;
    });

    it("Should not allow call from address other than operator", async () => {
      await expect(
        fixture.redistribution.registerAddressForDistribution(
          fixture.accounts[0].address,
          100_000
        )
      ).to.be.revertedWithCustomError(fixture.redistribution, "NotOperator");
    });

    it("Should not allow call from factory", async () => {
      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData(
            "registerAddressForDistribution",
            [fixture.accounts[0].address, 100_000]
          )
        )
      ).to.be.revertedWithCustomError(fixture.redistribution, "NotOperator");
    });

    it("Only operator can call", async () => {
      const operator = fixture.accounts[0];
      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData("setOperator", [
            operator.address,
          ])
        )
      ).to.be.fulfilled;

      await expect(
        fixture.redistribution
          .connect(operator)
          .registerAddressForDistribution(fixture.accounts[3].address, 100_000)
      ).to.be.fulfilled;
    });

    it("Cannot be called after expiration", async () => {
      const operator = fixture.accounts[0];
      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData("setOperator", [
            operator.address,
          ])
        )
      ).to.be.fulfilled;

      // lets make it expire
      await mineBlocks(
        BigNumber.from(DEFAULT_WITHDRAWAL_BLOCK_WINDOW).toBigInt()
      );

      await expect(
        fixture.redistribution
          .connect(operator)
          .registerAddressForDistribution(fixture.accounts[3].address, 100_000)
      ).to.be.revertedWithCustomError(
        fixture.redistribution,
        "WithdrawalWindowExpired"
      );
    });

    it("Should not allow to register an address that is already registered", async () => {
      const operator = fixture.accounts[0];
      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData("setOperator", [
            operator.address,
          ])
        )
      ).to.be.fulfilled;

      await expect(
        fixture.redistribution
          .connect(operator)
          .registerAddressForDistribution(fixture.accounts[3].address, 100_000)
      ).to.be.fulfilled;

      await expect(
        fixture.redistribution
          .connect(operator)
          .registerAddressForDistribution(fixture.accounts[3].address, 100_000)
      ).to.be.revertedWithCustomError(
        fixture.redistribution,
        "PositionAlreadyRegisteredOrTaken"
      );
    });

    it("Should not allow registering amount 0", async () => {
      const operator = fixture.accounts[0];
      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData("setOperator", [
            operator.address,
          ])
        )
      ).to.be.fulfilled;

      await expect(
        fixture.redistribution
          .connect(operator)
          .registerAddressForDistribution(fixture.accounts[3].address, 0)
      ).to.be.revertedWithCustomError(
        fixture.redistribution,
        "ZeroAmountNotAllowed"
      );
    });

    it("Should not allow registering more than available funds", async () => {
      const operator = fixture.accounts[0];
      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData("setOperator", [
            operator.address,
          ])
        )
      ).to.be.fulfilled;

      await expect(
        fixture.redistribution
          .connect(operator)
          .registerAddressForDistribution(fixture.accounts[3].address, 100_000)
      ).to.be.fulfilled;

      await expect(
        fixture.redistribution
          .connect(operator)
          .registerAddressForDistribution(
            fixture.accounts[4].address,
            DEFAULT_MAX_DISTRIBUTION_AMOUNT
          )
      ).to.be.revertedWithCustomError(
        fixture.redistribution,
        "InvalidDistributionAmount"
      );
    });

    it("Should not allow to register an address that already withdrew the position", async () => {
      const operator = fixture.accounts[0];
      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData("setOperator", [
            operator.address,
          ])
        )
      ).to.be.fulfilled;

      await expect(
        fixture.redistribution
          .connect(operator)
          .registerAddressForDistribution(fixture.accounts[3].address, 100_000)
      ).to.be.fulfilled;

      // withdraw
      await expect(
        fixture.redistribution
          .connect(fixture.accounts[3])
          .withdrawStakedPosition(fixture.accounts[3].address)
      ).to.be.fulfilled;

      await expect(
        fixture.redistribution
          .connect(operator)
          .registerAddressForDistribution(fixture.accounts[3].address, 100_000)
      ).to.be.revertedWithCustomError(
        fixture.redistribution,
        "PositionAlreadyRegisteredOrTaken"
      );
    });
  });

  /**
   * SendExpiredFundsToFactory testing
   */

  describe("SendExpiredFundsToFactory testing", async () => {
    let fixture: FixtureWithRedistribution;
    beforeEach(async () => {
      const deployFunc = async (): Promise<FixtureWithRedistribution> => {
        const baseFixture = await getFixture();
        const allAccounts = await ethers.getSigners();
        const distributionAccounts = allAccounts
          .slice(1, 3)
          .map((a) => a.address);
        const accountAmounts = distributionAccounts.map(() =>
          ethers.utils.parseEther(`500000`)
        );
        return await deployRedistribution(
          baseFixture,
          allAccounts,
          distributionAccounts,
          accountAmounts,
          DEFAULT_WITHDRAWAL_BLOCK_WINDOW,
          DEFAULT_MAX_DISTRIBUTION_AMOUNT
        );
      };
      fixture = await loadFixture(deployFunc);
    });

    it("Should not allow call from non-factory address", async () => {
      await expect(
        fixture.redistribution.sendExpiredFundsToFactory()
      ).to.be.revertedWithCustomError(fixture.redistribution, "OnlyFactory");
    });
    it("Should allow call from factory address", async () => {
      // lets make it expire
      await mineBlocks(
        BigNumber.from(DEFAULT_WITHDRAWAL_BLOCK_WINDOW).toBigInt()
      );

      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData(
            "sendExpiredFundsToFactory"
          )
        )
      ).to.be.fulfilled;
    });
    it("Should not allow call when not expired", async () => {
      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData(
            "sendExpiredFundsToFactory"
          )
        )
      ).to.be.revertedWithCustomError(
        fixture.redistribution,
        "WithdrawalWindowNotExpiredYet"
      );
    });
    it("Should emit TokenAlreadyTransferred event if contract is there's no token", async () => {
      // lets make it expire
      await mineBlocks(
        BigNumber.from(DEFAULT_WITHDRAWAL_BLOCK_WINDOW).toBigInt()
      );

      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData(
            "sendExpiredFundsToFactory"
          )
        )
      ).to.emit(fixture.redistribution, "TokenAlreadyTransferred");
    });
    it("Should transfer all ALCA and ETH funds to factory", async () => {
      // lets make it expire
      await mineBlocks(
        BigNumber.from(DEFAULT_WITHDRAWAL_BLOCK_WINDOW).toBigInt()
      );

      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData(
            "sendExpiredFundsToFactory"
          )
        )
      ).to.be.fulfilled;

      expect(
        (
          await fixture.alca.balanceOf(fixture.redistribution.address)
        ).toString()
      ).to.be.equal("0");
      expect(
        (
          await ethers.provider.getBalance(fixture.redistribution.address)
        ).toString()
      ).to.be.equal("0");
    });
    it("Should transfer all ALCA and ETH funds to factory with distribution of yield", async () => {
      await expect(
        fixture.factory.callAny(
          fixture.alca.address,
          0,
          fixture.alca.interface.encodeFunctionData("approve", [
            fixture.redistribution.address,
            DEFAULT_MAX_DISTRIBUTION_AMOUNT,
          ])
        )
      ).to.be.fulfilled;
      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData(
            "createRedistributionStakedPosition"
          )
        )
      ).to.be.fulfilled;

      // transfer ALCA
      const tokenAmount = ethers.utils.parseEther("100");
      await fixture.factory.callAny(
        fixture.alca.address,
        0,
        fixture.alca.interface.encodeFunctionData("approve", [
          fixture.publicStaking.address,
          tokenAmount,
        ])
      );
      await fixture.factory.callAny(
        fixture.publicStaking.address,
        0,
        fixture.publicStaking.interface.encodeFunctionData("depositToken", [
          42,
          tokenAmount,
        ])
      );

      // transfer ETH
      const ethAmount = ethers.utils.parseEther("10").toBigInt();
      await fixture.publicStaking.depositEth(42, { value: ethAmount });

      const factoryALCABalanceBefore = await fixture.alca.balanceOf(
        fixture.factory.address
      );
      const foundationETHBalanceBefore = await ethers.provider.getBalance(
        fixture.foundation.address
      );

      // lets make it expire
      await mineBlocks(
        BigNumber.from(DEFAULT_WITHDRAWAL_BLOCK_WINDOW).toBigInt()
      );

      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData(
            "sendExpiredFundsToFactory"
          )
        )
      ).to.be.fulfilled;

      expect(
        (await fixture.alca.balanceOf(fixture.factory.address)).toString()
      ).to.be.equal(factoryALCABalanceBefore.toString());
      expect(
        (
          await ethers.provider.getBalance(fixture.foundation.address)
        ).toString()
      ).to.be.equal(foundationETHBalanceBefore.toString());

      expect(
        (
          await fixture.alca.balanceOf(fixture.redistribution.address)
        ).toString()
      ).to.be.equal("0");
      expect(
        (
          await ethers.provider.getBalance(fixture.redistribution.address)
        ).toString()
      ).to.be.equal("0");
    });
    it("Should transfer any dangling ether or alca back to factory even if the position was transferred", async () => {
      // lets make it expire
      await mineBlocks(
        BigNumber.from(DEFAULT_WITHDRAWAL_BLOCK_WINDOW).toBigInt()
      );
      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData(
            "sendExpiredFundsToFactory"
          )
        )
      ).to.be.fulfilled;
      // trick to send eth to contract that have receive() protected
      const [admin] = await ethers.getSigners();
      const testUtils = await (
        await (await ethers.getContractFactory("TestUtils")).deploy()
      ).deployed();
      await admin.sendTransaction({
        to: testUtils.address,
        value: ethers.utils.parseEther("1"),
      });
      await testUtils.payUnpayable(fixture.redistribution.address);
      // end of the trick
      // sending alca
      await fixture.factory.callAny(
        fixture.alca.address,
        0,
        fixture.alca.interface.encodeFunctionData("transfer", [
          fixture.redistribution.address,
          ethers.utils.parseEther("2"),
        ])
      );
      const balanceEthBeforeRedistribution = await ethers.provider.getBalance(
        fixture.redistribution.address
      );
      const balanceALCABeforeRedistribution = await fixture.alca.balanceOf(
        fixture.redistribution.address
      );
      expect(balanceEthBeforeRedistribution).to.equal(
        ethers.utils.parseEther("1")
      );
      expect(balanceALCABeforeRedistribution).to.equal(
        ethers.utils.parseEther("2")
      );
      const foundationETHBalanceBefore = await ethers.provider.getBalance(
        fixture.foundation.address
      );
      const factoryALCABalanceBefore = await fixture.alca.balanceOf(
        fixture.factory.address
      );
      // making sure that can use this function to skim funds
      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData(
            "sendExpiredFundsToFactory"
          )
        )
      ).to.be.fulfilled;
      expect(
        (
          await ethers.provider.getBalance(fixture.redistribution.address)
        ).toString()
      ).to.be.equal("0");
      expect(
        (
          await fixture.alca.balanceOf(fixture.redistribution.address)
        ).toString()
      ).to.be.equal("0");
      expect(
        await ethers.provider.getBalance(fixture.foundation.address)
      ).to.be.equal(
        foundationETHBalanceBefore.add(balanceEthBeforeRedistribution)
      );
      expect(await fixture.alca.balanceOf(fixture.factory.address)).to.be.equal(
        factoryALCABalanceBefore.add(balanceALCABeforeRedistribution)
      );
    });
    it("Should transfer any dangling ether or alca back to factory", async () => {
      // lets make it expire
      await mineBlocks(
        BigNumber.from(DEFAULT_WITHDRAWAL_BLOCK_WINDOW).toBigInt()
      );
      // trick to send eth to contract that have receive() protected
      const [admin] = await ethers.getSigners();
      const testUtils = await (
        await (await ethers.getContractFactory("TestUtils")).deploy()
      ).deployed();
      await admin.sendTransaction({
        to: testUtils.address,
        value: ethers.utils.parseEther("1"),
      });
      await testUtils.payUnpayable(fixture.redistribution.address);
      // end of the trick
      // sending alca
      await fixture.factory.callAny(
        fixture.alca.address,
        0,
        fixture.alca.interface.encodeFunctionData("transfer", [
          fixture.redistribution.address,
          ethers.utils.parseEther("2"),
        ])
      );
      const balanceEthBeforeRedistribution = await ethers.provider.getBalance(
        fixture.redistribution.address
      );
      const balanceALCABeforeRedistribution = await fixture.alca.balanceOf(
        fixture.redistribution.address
      );
      expect(balanceEthBeforeRedistribution).to.equal(
        ethers.utils.parseEther("1")
      );
      expect(balanceALCABeforeRedistribution).to.equal(
        ethers.utils.parseEther("2")
      );
      const foundationETHBalanceBefore = await ethers.provider.getBalance(
        fixture.foundation.address
      );
      const factoryALCABalanceBefore = await fixture.alca.balanceOf(
        fixture.factory.address
      );
      // making sure that can use this function to skim funds
      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData(
            "sendExpiredFundsToFactory"
          )
        )
      ).to.be.fulfilled;
      expect(
        (
          await ethers.provider.getBalance(fixture.redistribution.address)
        ).toString()
      ).to.be.equal("0");
      expect(
        (
          await fixture.alca.balanceOf(fixture.redistribution.address)
        ).toString()
      ).to.be.equal("0");
      expect(
        await ethers.provider.getBalance(fixture.foundation.address)
      ).to.be.equal(
        foundationETHBalanceBefore.add(balanceEthBeforeRedistribution)
      );
      expect(await fixture.alca.balanceOf(fixture.factory.address)).to.be.equal(
        factoryALCABalanceBefore.add(balanceALCABeforeRedistribution)
      );
    });
    it("Should be able to burn the position with the factory after the expiration", async () => {
      await fixture.factory.callAny(
        fixture.alca.address,
        0,
        fixture.alca.interface.encodeFunctionData("approve", [
          fixture.redistribution.address,
          DEFAULT_MAX_DISTRIBUTION_AMOUNT,
        ])
      );
      await fixture.factory.callAny(
        fixture.redistribution.address,
        0,
        fixture.redistribution.interface.encodeFunctionData(
          "createRedistributionStakedPosition"
        )
      );
      await mineBlocks(1n);
      // lets make it expire
      await mineBlocks(
        BigNumber.from(DEFAULT_WITHDRAWAL_BLOCK_WINDOW).toBigInt()
      );
      expect(
        await fixture.publicStaking.balanceOf(fixture.factory.address)
      ).to.be.equal(BigNumber.from(0));
      const tokenID = await fixture.redistribution.tokenID();
      const [reserve] = await fixture.publicStaking.getPosition(tokenID);
      await expect(
        fixture.factory.callAny(
          fixture.redistribution.address,
          0,
          fixture.redistribution.interface.encodeFunctionData(
            "sendExpiredFundsToFactory"
          )
        )
      ).to.be.fulfilled;
      expect(
        await fixture.publicStaking.balanceOf(fixture.factory.address)
      ).to.be.equal(BigNumber.from(1));
      const balanceALCABefore = await fixture.alca.balanceOf(
        fixture.factory.address
      );
      await fixture.factory.callAny(
        fixture.publicStaking.address,
        0,
        fixture.publicStaking.interface.encodeFunctionData("burn", [tokenID])
      );
      expect(await fixture.alca.balanceOf(fixture.factory.address)).to.be.equal(
        balanceALCABefore.add(reserve)
      );
    });
  });
  // first tests
  describe("First round of tests", async () => {
    let fixture: FixtureWithRedistribution;
    beforeEach(async () => {
      fixture = await loadFixture(deployFunc);
    });

    it("should have correct maxRedistributionAmount", async () => {
      const maxRedistributionAmount =
        await fixture.redistribution.maxRedistributionAmount();
      expect(maxRedistributionAmount).to.equal(DEFAULT_MAX_DISTRIBUTION_AMOUNT);
    });

    it("should not allow minting a position without an operator", async () => {
      await expect(
        fixture.redistribution
          .connect(fixture.accounts[5])
          .registerAddressForDistribution(fixture.accounts[6].address, 100_000)
      ).to.revertedWithCustomError(fixture.redistribution, "NotOperator");
    });

    it("should register a position as operator", async () => {
      const operator = fixture.accounts[5];
      await fixture.factory.callAny(
        fixture.redistribution.address,
        0,
        fixture.redistribution.interface.encodeFunctionData("setOperator", [
          operator.address,
        ])
      );

      await fixture.redistribution
        .connect(operator)
        .registerAddressForDistribution(fixture.accounts[6].address, 100_000);

      expect(
        await fixture.redistribution.getRedistributionInfo(
          fixture.accounts[6].address
        )
      ).to.be.deep.equal([BigNumber.from(100_000), false]);
    });

    it("should not be able to withdraw if total stake position was not created", async () => {
      await expect(
        fixture.redistribution
          .connect(fixture.accounts[4])
          .withdrawStakedPosition(fixture.accounts[5].address)
      ).to.be.rejectedWith("ERC721: invalid token ID");
    });
  });
  describe("Withdraw tests", async () => {
    let fixture: FixtureWithRedistribution;
    beforeEach(async () => {
      fixture = await loadFixture(deployFunc);
      await fixture.factory.callAny(
        fixture.alca.address,
        0,
        fixture.alca.interface.encodeFunctionData("approve", [
          fixture.redistribution.address,
          DEFAULT_MAX_DISTRIBUTION_AMOUNT,
        ])
      );
      await fixture.factory.callAny(
        fixture.redistribution.address,
        0,
        fixture.redistribution.interface.encodeFunctionData(
          "createRedistributionStakedPosition"
        )
      );
      await mineBlocks(1n);
    });

    it("accounts registered in the constructor should be able to withdraw", async () => {
      for (let i = 0; i < 5; i++) {
        expect(
          await fixture.redistribution.getRedistributionInfo(
            fixture.accounts[i].address
          )
        ).to.be.deep.equal([ethers.utils.parseEther(`500000`), false]);
        expect(
          await fixture.publicStaking.balanceOf(fixture.accounts[i].address)
        ).to.be.equal(BigNumber.from(0));
        await fixture.redistribution
          .connect(fixture.accounts[i])
          .withdrawStakedPosition(fixture.accounts[i].address);
        await mineBlocks(1n);
        expect(
          await fixture.publicStaking.balanceOf(fixture.accounts[i].address)
        ).to.be.equal(BigNumber.from(1));
        expect(
          await fixture.redistribution.getRedistributionInfo(
            fixture.accounts[i].address
          )
        ).to.be.deep.equal([BigNumber.from(0), true]);
      }
    });

    it("accounts registered after constructor should be able to withdraw", async () => {
      const operator = fixture.accounts[5];
      await fixture.factory.callAny(
        fixture.redistribution.address,
        0,
        fixture.redistribution.interface.encodeFunctionData("setOperator", [
          operator.address,
        ])
      );

      await fixture.redistribution
        .connect(operator)
        .registerAddressForDistribution(
          fixture.accounts[6].address,
          ethers.utils.parseEther("100000")
        );

      expect(
        await fixture.redistribution.getRedistributionInfo(
          fixture.accounts[6].address
        )
      ).to.be.deep.equal([ethers.utils.parseEther(`100000`), false]);
      expect(
        await fixture.publicStaking.balanceOf(fixture.accounts[6].address)
      ).to.be.equal(BigNumber.from(0));
      await expect(
        fixture.redistribution
          .connect(fixture.accounts[6])
          .withdrawStakedPosition(fixture.accounts[6].address)
      )
        .to.emit(fixture.redistribution, "Withdrawn")
        .withArgs(
          fixture.accounts[6].address,
          ethers.utils.parseEther(`100000`)
        );
      await mineBlocks(1n);
      expect(
        await fixture.publicStaking.balanceOf(fixture.accounts[6].address)
      ).to.be.equal(BigNumber.from(1));
      expect(
        await fixture.redistribution.getRedistributionInfo(
          fixture.accounts[6].address
        )
      ).to.be.deep.equal([BigNumber.from(0), true]);
      // accounts registered in the contract should also be able to withdraw
      await expect(
        fixture.redistribution
          .connect(fixture.accounts[0])
          .withdrawStakedPosition(fixture.accounts[0].address)
      )
        .to.emit(fixture.redistribution, "Withdrawn")
        .withArgs(
          fixture.accounts[0].address,
          ethers.utils.parseEther(`500000`)
        );
    });

    it("should not allow withdrawing twice", async () => {
      await fixture.redistribution
        .connect(fixture.accounts[0])
        .withdrawStakedPosition(fixture.accounts[0].address);
      await mineBlocks(1n);
      await expect(
        fixture.redistribution
          .connect(fixture.accounts[0])
          .withdrawStakedPosition(fixture.accounts[0].address)
      ).to.be.revertedWithCustomError(
        fixture.redistribution,
        "PositionAlreadyTakenOrInexistent"
      );
    });

    it("should not allow withdrawing if not registered", async () => {
      await expect(
        fixture.redistribution
          .connect(fixture.accounts[5])
          .withdrawStakedPosition(fixture.accounts[5].address)
      ).to.be.revertedWithCustomError(
        fixture.redistribution,
        "PositionAlreadyTakenOrInexistent"
      );
    });

    it("should not allow withdrawing if system expired", async () => {
      await mineBlocks((await fixture.redistribution.expireBlock()).toBigInt());
      await expect(
        fixture.redistribution.withdrawStakedPosition(
          fixture.accounts[0].address
        )
      ).to.be.revertedWithCustomError(
        fixture.redistribution,
        "WithdrawalWindowExpired"
      );
    });

    it("should withdraw with eth public staking yields going to foundation", async () => {
      // since there's only 1 position in the system, all profit should go to it
      const ethAmount = ethers.utils.parseEther("10").toBigInt();
      await fixture.publicStaking.depositEth(42, { value: ethAmount });
      expect(
        await ethers.provider.getBalance(fixture.foundation.address)
      ).to.be.equal(BigNumber.from(0));
      const balanceBefore = await ethers.provider.getBalance(
        fixture.accounts[0].address
      );
      const rcpt = await (
        await fixture.redistribution
          .connect(fixture.accounts[0])
          .withdrawStakedPosition(fixture.accounts[0].address)
      ).wait();
      // balance of the account should be the same as before, but with the gas used
      expect(
        await ethers.provider.getBalance(fixture.accounts[0].address)
      ).to.be.equal(
        balanceBefore.sub(rcpt.gasUsed.mul(rcpt.effectiveGasPrice))
      );
      await mineBlocks(1n);
      // all eth profit should be forward to the foundation
      expect(
        await ethers.provider.getBalance(fixture.foundation.address)
      ).to.be.equal(ethers.utils.parseEther("10"));
    });

    it("should withdraw with alca public staking yields being re-staked", async () => {
      const tokenAmount = ethers.utils.parseEther("100");
      await fixture.factory.callAny(
        fixture.alca.address,
        0,
        fixture.alca.interface.encodeFunctionData("approve", [
          fixture.publicStaking.address,
          tokenAmount,
        ])
      );
      await fixture.factory.callAny(
        fixture.publicStaking.address,
        0,
        fixture.publicStaking.interface.encodeFunctionData("depositToken", [
          42,
          tokenAmount,
        ])
      );
      // since there's only 1 position in the system, all profit should go to it
      const balanceBefore = await fixture.alca.balanceOf(
        fixture.accounts[0].address
      );
      const maxAmountBefore =
        await fixture.redistribution.maxRedistributionAmount();
      const totalAllowancesBefore =
        await fixture.redistribution.totalAllowances();
      const [reserveBefore] = await fixture.publicStaking.getPosition(
        await fixture.redistribution.tokenID()
      );
      await fixture.redistribution
        .connect(fixture.accounts[0])
        .withdrawStakedPosition(fixture.accounts[0].address);
      // balance of the account should be the same as before, but with the gas used
      expect(
        await fixture.alca.balanceOf(fixture.accounts[0].address)
      ).to.be.equal(balanceBefore);
      await mineBlocks(1n);
      // all eth profit should be forward to the foundation
      const [reserveAfter] = await fixture.publicStaking.getPosition(
        await fixture.redistribution.tokenID()
      );
      // the new position should be the same as the old one, but with the 100 alca (profit) more and less the
      // amount sent to the user
      expect(reserveAfter).to.be.equal(
        reserveBefore.add(tokenAmount).sub(ethers.utils.parseEther("500000"))
      );
      expect(
        await fixture.redistribution.maxRedistributionAmount()
      ).to.be.equal(maxAmountBefore);
      expect(await fixture.redistribution.totalAllowances()).to.be.equal(
        totalAllowancesBefore
      );
    });

    it("should not send eth to foundation in case there was no distribution", async () => {
      await fixture.redistribution
        .connect(fixture.accounts[0])
        .withdrawStakedPosition(fixture.accounts[0].address);
      await mineBlocks(1n);
      // all eth profit should be forward to the foundation
      expect(
        await ethers.provider.getBalance(fixture.foundation.address)
      ).to.be.equal(BigNumber.from(0));
    });

    it("should not re-stake if all reserved position were withdrawn and there was not alca yield", async () => {
      const operator = fixture.accounts[5];
      await fixture.factory.callAny(
        fixture.redistribution.address,
        0,
        fixture.redistribution.interface.encodeFunctionData("setOperator", [
          operator.address,
        ])
      );

      await fixture.redistribution
        .connect(operator)
        .registerAddressForDistribution(
          fixture.accounts[6].address,
          await fixture.redistribution.getDistributionLeft()
        );

      for (let i = 0; i < 5; i++) {
        const tokenBefore = await fixture.redistribution.tokenID();
        await fixture.redistribution
          .connect(fixture.accounts[i])
          .withdrawStakedPosition(fixture.accounts[i].address);
        await mineBlocks(1n);
        expect(await fixture.redistribution.tokenID()).to.be.not.equal(
          tokenBefore
        );
      }
      const tokenBefore = await fixture.redistribution.tokenID();
      await fixture.redistribution
        .connect(fixture.accounts[6])
        .withdrawStakedPosition(fixture.accounts[6].address);
      // since all alca was withdrawn, no new position should be created to the redistribution
      // contract
      expect(await fixture.redistribution.tokenID()).to.be.equal(tokenBefore);
    });

    it("should re-stake if all reserved position were withdrawn and there was alca yield", async () => {
      const tokenAmount = ethers.utils.parseEther("100");
      await fixture.factory.callAny(
        fixture.alca.address,
        0,
        fixture.alca.interface.encodeFunctionData("approve", [
          fixture.publicStaking.address,
          tokenAmount,
        ])
      );
      await fixture.factory.callAny(
        fixture.publicStaking.address,
        0,
        fixture.publicStaking.interface.encodeFunctionData("depositToken", [
          42,
          tokenAmount,
        ])
      );
      const operator = fixture.accounts[5];
      await fixture.factory.callAny(
        fixture.redistribution.address,
        0,
        fixture.redistribution.interface.encodeFunctionData("setOperator", [
          operator.address,
        ])
      );

      await fixture.redistribution
        .connect(operator)
        .registerAddressForDistribution(
          fixture.accounts[6].address,
          await fixture.redistribution.getDistributionLeft()
        );

      for (let i = 0; i < 5; i++) {
        const tokenBefore = await fixture.redistribution.tokenID();
        await fixture.redistribution
          .connect(fixture.accounts[i])
          .withdrawStakedPosition(fixture.accounts[i].address);
        await mineBlocks(1n);
        expect(await fixture.redistribution.tokenID()).to.be.not.equal(
          tokenBefore
        );
      }
      const tokenBefore = await fixture.redistribution.tokenID();
      await fixture.redistribution
        .connect(fixture.accounts[6])
        .withdrawStakedPosition(fixture.accounts[6].address);
      // since all alca was withdrawn, no new position should be created to the redistribution
      // contract
      expect(await fixture.redistribution.tokenID()).to.be.not.equal(
        tokenBefore
      );
      const [reserve] = await fixture.publicStaking.getPosition(
        await fixture.redistribution.tokenID()
      );
      expect(reserve).to.be.equal(tokenAmount);
    });
  });
});
