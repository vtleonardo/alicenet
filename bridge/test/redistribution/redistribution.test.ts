import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
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

  it.only("should have correct maxRedistributionAmount", async () => {
    const maxRedistributionAmount =
      await fixture.redistribution.maxRedistributionAmount();
    expect(maxRedistributionAmount).to.equal(maxDistributionAmount);
  });
});
