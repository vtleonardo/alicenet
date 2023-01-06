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
 * Should not allow call from address other than factory
 * Only factory can call
 * Cannot be called after expiration
 * Should fail if there was not alca approval from the factory
 * Should create a position if all the conditions are met, check tokenID and the position properties match the maxRedistributionAmount
 * Should not allow to create a position if it already exists
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

    /**
     *
     *
     * Should create a position if all the conditions are met, check tokenID and the position properties match the maxRedistributionAmount
     * Should not allow to create a position if it already exists
     */

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
      const tokenID = await fixture.redistribution.tokenID();
      expect(tokenID.toString()).to.be.not.equal("0");
      expect(await fixture.publicStaking.ownerOf(tokenID)).to.be.equal(
        fixture.redistribution.address
      );
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
  });

  // first tests

  describe("First round of tests", async () => {
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
  });
});
