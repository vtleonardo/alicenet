import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { ethers } from "hardhat";
import {
  deployUpgradeable,
  getEventVar,
  multiCallUpgradeProxy,
} from "../../scripts/lib/alicenetFactory";
import { CONTRACT_ADDR, MOCK, PROXY, UTILS } from "../../scripts/lib/constants";
import { getGasPrices } from "../../scripts/lib/deployment/utils";
import { AliceNetFactory, Utils } from "../../typechain-types";
import { expect } from "../chai-setup";
import { deployFactory } from "./Setup";
process.env.silencer = "true";

describe("AliceNetfactory API test", async () => {
  let utilsContract: Utils;
  let factory: AliceNetFactory;

  async function deployFixture() {
    const utilsBase = await ethers.getContractFactory(UTILS);
    const utilsContract = await utilsBase.deploy();
    const factory = await deployFactory();
    return { utilsContract, factory };
  }

  beforeEach(async () => {
    ({ utilsContract, factory } = await loadFixture(deployFixture));

    const cSize = await utilsContract.getCodeSize(factory.address);
    expect(cSize.toNumber()).to.be.greaterThan(0);
  });

  it("deploy Upgradeable", async () => {
    const salt = ethers.utils.formatBytes32String(MOCK);
    const txResponse = await deployUpgradeable(
      MOCK,
      factory,
      ethers,
      "0x",
      ["2", "s"],
      salt
    );
    const receipt = await txResponse.wait();
    const proxyAddress = getEventVar(receipt, "DeployedProxy", CONTRACT_ADDR);

    const proxy = await ethers.getContractAt(PROXY, proxyAddress);
    const implementationAddress = await factory.getProxyImplementation(
      proxy.address
    );
    expect(implementationAddress).to.not.equal(ethers.constants.AddressZero);
    let cSize = await utilsContract.getCodeSize(implementationAddress);
    expect(cSize.toNumber()).to.be.greaterThan(0);
    cSize = await utilsContract.getCodeSize(proxyAddress);
    expect(cSize.toNumber()).to.be.greaterThan(0);
  });

  it("upgrade deployment", async () => {
    const salt = ethers.utils.formatBytes32String(MOCK);
    const logicContractBase = await ethers.getContractFactory(MOCK);
    let txResponse = await deployUpgradeable(
      MOCK,
      factory,
      ethers,
      "0x",
      ["2", "s"],
      salt,
      1,
      await getGasPrices(ethers)
    );
    let receipt = await txResponse.wait();
    const proxyAddress = getEventVar(receipt, "DeployedProxy", CONTRACT_ADDR);
    const proxy = await ethers.getContractAt(PROXY, proxyAddress);
    const implementationAddress = await factory.getProxyImplementation(
      proxy.address
    );
    txResponse = await multiCallUpgradeProxy(
      logicContractBase,
      factory,
      ethers,
      "0x",
      ["2", "s"],
      salt,
      await getGasPrices(ethers)
    );
    receipt = await txResponse.wait();
    const expectedImplementationAddress = getEventVar(
      receipt,
      "DeployedRaw",
      CONTRACT_ADDR
    );
    const newImplementationAddress = await factory.getProxyImplementation(
      proxy.address
    );
    expect(newImplementationAddress).to.not.equal(implementationAddress);
    expect(newImplementationAddress).to.equal(expectedImplementationAddress);
  });
});
