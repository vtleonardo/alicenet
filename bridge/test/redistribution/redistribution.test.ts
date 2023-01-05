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
} from "../setup";

const withdrawalBlockWindow = 172800; // 1 month block window
const maxDistributionAmount = ethers.utils.parseEther("10000000");

/**
 * Constructor:
 * Should not allow sum of allowedAmounts to be greater than the maxRedistributionAmount
 * Should not allow addresses and amounts with different length
 * Should not allow length zero in any of the parameters
 * Should not allow duplicate addresses and address zero
 * Should not allow distribution amount 0
 * Check totalAllowances must be equal to the sum of allowedAmounts
 * Check expireBlock must be equal to block.number + withdrawalBlockWindow
 *
 * SetOperator:
 * Should not allow call from address other than factory
 * Only factory can call
 * Should be able to set a operator as factory (operator is public we can check with the getter)
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

async function deployRedistribution(): Promise<FixtureWithRedistribution> {
  const fixture = await getFixture();
  const accounts = await ethers.getSigners();
  const accountAddresses: string[] = [];
  const accountAmounts: BigInt[] = [];
  for (let i = 0; i < 5; i++) {
    accountAddresses.push(accounts[i].address);
    accountAmounts.push(
      ethers.utils.parseEther(`${i + 1}`).toBigInt() * 500_000n
    );
  }
  const tx = await deployCreate("Redistribution", fixture.factory, ethers, [
    withdrawalBlockWindow,
    maxDistributionAmount,
    accountAddresses,
    accountAmounts,
  ]);
  const redistributionAddress = await getContractAddressFromDeployedRawEvent(
    tx
  );
  const redistribution = await ethers.getContractAt(
    "Redistribution",
    redistributionAddress
  );
  await fixture.factory.callAny(
    fixture.alca.address,
    0,
    fixture.alca.interface.encodeFunctionData("approve", [
      redistribution.address,
      maxDistributionAmount,
    ])
  );
  await fixture.factory.callAny(
    redistribution.address,
    0,
    redistribution.interface.encodeFunctionData(
      "createRedistributionStakedPosition"
    )
  );
  return {
    ...fixture,
    accounts,
    redistribution,
  };
}

describe("CT redistribution", async () => {
  let fixture: FixtureWithRedistribution;
  beforeEach(async () => {
    fixture = await loadFixture(deployRedistribution);
  });

  it("should have correct maxRedistributionAmount", async () => {
    const maxRedistributionAmount =
      await fixture.redistribution.maxRedistributionAmount();
    expect(maxRedistributionAmount).to.equal(maxDistributionAmount);
  });

  it("should not allow minting a position without an operator", async () => {
    await expect(
      fixture.redistribution
        .connect(fixture.accounts[5])
        .registerAddressForDistribution(
          fixture.accounts[6].getAddress(),
          100_000
        )
    ).to.revertedWithCustomError(fixture.redistribution, "NotOperator");
  });

  it("should be able to get tokenID for total staked position", async () => {
    expect((await fixture.redistribution.tokenID()).toBigInt()).to.equal(1n);
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
