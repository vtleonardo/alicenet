import { BigNumber } from "ethers";
import { ethers, network } from "hardhat";
import { BondingCurveStressMock } from "../../typechain-types";

describe("Stress testing BToken mint and burn", async () => {
  let bTokenMock: BondingCurveStressMock;
  let maxGasLimit: BigNumber;
  beforeEach(async function () {
    bTokenMock = await (
      await ethers.getContractFactory("BondingCurveStressMock")
    ).deploy();
    maxGasLimit = BigNumber.from("9007199254740991");
  });

  it("Stress test minting", async () => {
    const isToPrint = true;
    await network.provider.send("evm_setBlockGasLimit", [
      "0x3000000000000000000",
    ]);
    console.log("ether,mintedBToken,sentEth");
    await bTokenMock.stressMint(isToPrint, 4n * 10n ** 0n, 10, {
      gasLimit: maxGasLimit,
    });
    await bTokenMock.stressMint(isToPrint, 10n ** 3n, 40, {
      gasLimit: maxGasLimit,
    });
    await bTokenMock.stressMint(isToPrint, 10n ** 6n, 100, {
      gasLimit: maxGasLimit,
    });
    await bTokenMock.stressMint(isToPrint, 10n ** 9n, 100, {
      gasLimit: maxGasLimit,
    });
    await bTokenMock.stressMint(isToPrint, 10n ** 12n, 250, {
      gasLimit: maxGasLimit,
    });
    await bTokenMock.stressMint(isToPrint, 10n ** 15n, 500, {
      gasLimit: maxGasLimit,
    });
    await bTokenMock.stressMint(isToPrint, 10n ** 18n, 500, {
      gasLimit: maxGasLimit,
    });
    await bTokenMock.stressMint(isToPrint, 10n ** 19n, 500, {
      gasLimit: maxGasLimit,
    });
    await bTokenMock.stressMint(isToPrint, 10n ** 20n, 500, {
      gasLimit: maxGasLimit,
    });
    await bTokenMock.stressMint(isToPrint, 10n ** 21n, 500, {
      gasLimit: maxGasLimit,
    });
    await bTokenMock.stressMint(isToPrint, 10n ** 22n, 500, {
      gasLimit: maxGasLimit,
    });
    await bTokenMock.stressMint(isToPrint, 10n ** 23n, 500, {
      gasLimit: maxGasLimit,
    });
    await bTokenMock.stressMint(isToPrint, 10n ** 24n, 500, {
      gasLimit: maxGasLimit,
    });
    await bTokenMock.stressMint(isToPrint, 10n ** 25n, 500, {
      gasLimit: maxGasLimit,
    });
  });

  it("Stress test burning", async () => {
    const isToPrint = true;
    const steps = [
      10n ** 21n,
      10n ** 20n,
      10n ** 19n,
      10n ** 18n,
      10n ** 15n,
      10n ** 12n,
      10n ** 9n,
      10n ** 6n,
    ];
    const iterations = [1000n, 500n, 500n, 500n, 500n, 500n, 1000n, 1000n];
    let burnedTotalValue = 0n;
    for (let i = 0; i < steps.length; i++) {
      burnedTotalValue += steps[i] * iterations[i];
    }
    burnedTotalValue = (
      await bTokenMock.getLatestEthToMintBTokens(burnedTotalValue)
    ).toBigInt();
    await network.provider.send("evm_setBlockGasLimit", [
      "0x3000000000000000000",
    ]);
    await bTokenMock.mint(0, {
      value: burnedTotalValue + 1_000_000_000n,
    });
    console.log("ether,burnedBToken,receivedEth");
    for (let i = 0; i < steps.length; i++) {
      await bTokenMock.stressBurn(isToPrint, steps[i], iterations[i], {
        gasLimit: maxGasLimit,
      });
    }
  });

  it.only("Stress test burning and minting", async () => {
    const isToPrint = true;
    const steps = [
      4n * 10n ** 0n,
      10n ** 3n,
      10n ** 6n,
      10n ** 9n,
      10n ** 12n,
      10n ** 15n,
      10n ** 18n,
      10n ** 19n,
      10n ** 20n,
      10n ** 21n,
      10n ** 22n,
      10n ** 23n,
      10n ** 24n,
      10n ** 25n,
    ];
    const iterations = [
      50n,
      50n,
      50n,
      50n,
      50n,
      50n,
      50n,
      50n,
      50n,
      50n,
      25n,
      25n,
      20n,
      10n,
    ];
    const startingPoints = [
      4n * 10n ** 0n,
      10n ** 3n,
      10n ** 6n,
      10n ** 9n,
      10n ** 12n,
      10n ** 15n,
      10n ** 18n,
      10n ** 19n,
      10n ** 20n,
      10n ** 21n,
      10n ** 22n,
      10n ** 23n,
      10n ** 24n,
      10n ** 25n,
    ];
    let totalValue = 0n;
    for (let i = 0; i < steps.length; i++) {
      totalValue += steps[i] * iterations[i];
    }
    await network.provider.send("evm_setBlockGasLimit", [
      "0x3000000000000000000",
    ]);
    console.log("ether,sentEth,receivedEth");
    for (let j = 0; j < startingPoints.length; j++) {
      await bTokenMock.mint(0, {
        value: startingPoints[j],
      });
      for (let i = 0; i < steps.length; i++) {
        await bTokenMock.stressMintAndBurn(isToPrint, steps[i], iterations[i], {
          gasLimit: maxGasLimit,
          value: totalValue,
        });
      }
      await bTokenMock.burn(await bTokenMock.totalSupply(), 0);
    }
  });
});
