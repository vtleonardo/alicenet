import {
  BigNumber,
  BigNumberish,
  BytesLike,
  ContractFactory,
  ContractReceipt,
  ContractTransaction,
  Overrides,
} from "ethers";
import { ethers } from "hardhat";
import {
  AliceNetFactory,
  AliceNetFactory__factory as aliceNetFactoryBase,
} from "../../typechain-types";
import { PromiseOrValue } from "../../typechain-types/common";
import {
  ALICENET_FACTORY,
  CONTRACT_ADDR,
  EVENT_DEPLOYED_RAW,
  MULTICALL_GAS_LIMIT,
} from "./constants";
type Ethers = typeof ethers;

export type MultiCallArgsStruct = {
  target: string;
  value: BigNumberish;
  data: BytesLike;
};
export class MultiCallGasError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "MultiCallGasError";
  }
}
export async function deployFactory(
  legacyTokenAddress: PromiseOrValue<string>,
  ethers: Ethers,
  factoryBase?: aliceNetFactoryBase,
  overrides?: Overrides & { from?: PromiseOrValue<string> }
): Promise<AliceNetFactory> {
  factoryBase =
    factoryBase === undefined
      ? await ethers.getContractFactory(ALICENET_FACTORY)
      : factoryBase;
  if (overrides === undefined) {
    return factoryBase.deploy(legacyTokenAddress);
  } else {
    return factoryBase.deploy(legacyTokenAddress, overrides);
  }
}

// multicall deploy logic, proxy, and upgrade proxy
/**
 * @description uses multicall to deploy logic contract with deployCreate, deploys proxy with deployProxy, and upgrades proxy with upgradeProxy
 * @dev since upgradeable contracts go through proxies, constructor args can only be used to set immutable variables
 * this function will fail if gas cost exceeds 10 million gas units
 * @param implementationBase ethers contract factory for the implementation contract
 * @param factory an instance of a deployed and connected factory
 * @param ethers ethers object
 * @param initCallData encoded initialization call data for contracts with a initialize function
 * @param constructorArgs a list of arguements to pass to the constructor of the implementation contract, only for immutable variables
 * @param salt bytes32 formatted salt used for deploycreate2 and to reference the contract in lookup
 * @param overrides
 * @returns a promise that resolves to the deployed contracts
 */
export async function multiCallDeployUpgradeable(
  implementationBase: ContractFactory,
  factory: AliceNetFactory,
  ethers: Ethers,
  initCallData: string,
  constructorArgs: any[] = [],
  salt: string,
  overrides?: Overrides & { from?: PromiseOrValue<string> }
): Promise<ContractTransaction> {
  const multiCallArgs = await encodeMultiCallDeployUpgradeableArgs(
    implementationBase,
    factory,
    ethers,
    initCallData,
    constructorArgs,
    salt
  );
  const estimatedMultiCallGas = await factory.estimateGas.multiCall(
    multiCallArgs
  );
  if (estimatedMultiCallGas.gt(BigNumber.from(MULTICALL_GAS_LIMIT))) {
    throw new MultiCallGasError(
      `estimatedGasCost ${estimatedMultiCallGas.toString()} exceeds MULTICALL_GAS_LIMIT ${MULTICALL_GAS_LIMIT}`
    );
  }
  if (overrides === undefined) {
    return factory.multiCall(multiCallArgs);
  } else {
    return factory.multiCall(multiCallArgs, overrides);
  }
}

// upgradeProxy

/**
 * @description multicall deployCreate and upgradeProxy, throws if gas exceeds 10 million
 * @param implementationBase instance of the logic contract base
 * @param factory ethers connected instance of alicenet factory
 * @param ethers ethers js object
 * @param constructorArgs array of constructor arguments
 * @param initCallData encoded init calldata, 0x if no initializer function
 * @param salt bytes32 formatted salt used for deployProxy and to reference the contract in lookup
 * @param overrides transaction overrides
 * @returns a promise that resolves to the ContractTransaction
 */
export async function multiCallUpgradeProxy(
  implementationBase: ContractFactory,
  factory: AliceNetFactory,
  ethers: Ethers,
  initCallData: string,
  constructorArgs: any[] = [],
  salt: string,
  overrides?: Overrides & { from?: PromiseOrValue<string> }
) {
  const multiCallArgs = await encodeMultiCallUpgradeProxyArgs(
    implementationBase,
    factory,
    ethers,
    initCallData,
    constructorArgs,
    salt
  );
  const estimatedMultiCallGas = await factory.estimateGas.multiCall(
    multiCallArgs
  );
  if (estimatedMultiCallGas.gt(MULTICALL_GAS_LIMIT)) {
    throw new MultiCallGasError(
      `estimatedGasCost ${estimatedMultiCallGas.toString()} exceeds MULTICALL_GAS_LIMIT ${MULTICALL_GAS_LIMIT}`
    );
  }
  if (overrides === undefined) {
    return factory.multiCall(multiCallArgs);
  } else {
    return factory.multiCall(multiCallArgs, overrides);
  }
}

/**
 * @description uses factory multicall to deploy a proxy contract with deployProxy, then upgrades the proxy with upgradeProxy
 * @param logicAddress address of the logic contract already deployed
 * @param factory instance of deployed and connected alicenetFactory
 * @param salt bytes32 formatted salt used for deployCreate2 and to reference the contract in lookup
 * @param initCallData encoded init calldata, 0x if no initializer function
 * @returns
 */
export async function multiCallDeployProxyAndUpgradeProxy(
  logicAddress: string,
  factory: AliceNetFactory,
  salt: string,
  initCallData: string
) {
  const multiCallArgs = await encodeMultiCallDeployProxyAndUpgradeProxyArgs(
    logicAddress,
    factory,
    initCallData,
    salt
  );
  return factory.multiCall(multiCallArgs);
}

export async function factoryMultiCall(
  factory: AliceNetFactory,
  multiCallArgs: MultiCallArgsStruct[],
  overrides?: Overrides & { from?: PromiseOrValue<string> }
) {
  if (overrides === undefined) {
    return factory.multiCall(multiCallArgs);
  } else {
    return factory.multiCall(multiCallArgs, overrides);
  }
}

/**
 * @description encodes multicall for deployProxy and upgradeProxy
 * @dev this function is used if the logic contract is too big to be deployed with a full multicall
 * this just deploys a proxy and upgrade the proxy to point to the implementationAddress
 * @param implementationAddress address of the implementation contract
 * @param factory instance of deployed and connected alicenet factory
 * @param initCallData encoded init calldata, 0x if no initializer function
 * @param salt bytes32 format of the salt that references the proxy contract
 * @returns
 */
export async function encodeMultiCallDeployProxyAndUpgradeProxyArgs(
  implementationAddress: string,
  factory: AliceNetFactory,
  initCallData: string,
  salt: string
) {
  const deployProxyCallData: BytesLike = factory.interface.encodeFunctionData(
    "deployProxy",
    [salt]
  );
  const deployProxy = encodeMultiCallArgs(
    factory.address,
    0,
    deployProxyCallData
  );
  const upgradeProxyCallData = factory.interface.encodeFunctionData(
    "upgradeProxy",
    [salt, implementationAddress, initCallData]
  );
  const upgradeProxy = encodeMultiCallArgs(
    factory.address,
    0,
    upgradeProxyCallData
  );
  return [deployProxy, upgradeProxy];
}

/**
 * @description encodes the arguments for alicenet factory multicall to
 * deploy a logic contract with deploycreate,
 * deploy a proxy with deployProxy,
 * and upgrade the proxy with upgradeProxy
 * @param implementationBase ethers contract factory for the implementation contract
 * @param factory instance of deployed and connected alicenetFactory
 * @param ethers instance of ethers
 * @param initCallData encoded call data for the initialize function of the implementation contract
 * @param constructorArgs string array of constructor arguments, only used to set immutable variables
 * @param salt bytes32 formatted salt used for deploycreate2 and to reference the contract in lookup
 * @returns an array of encoded multicall data for deployCreate, deployProxy, and upgradeProxy
 */
export async function encodeMultiCallDeployUpgradeableArgs(
  implementationBase: ContractFactory,
  factory: AliceNetFactory,
  ethers: Ethers,
  initCallData: string,
  constructorArgs: string[] = [],
  salt: string
) {
  const deployProxyCallData: BytesLike = factory.interface.encodeFunctionData(
    "deployProxy",
    [salt]
  );
  const deployProxy = encodeMultiCallArgs(
    factory.address,
    0,
    deployProxyCallData
  );
  const [deployCreate, upgradeProxy] = await encodeMultiCallUpgradeProxyArgs(
    implementationBase,
    factory,
    ethers,
    initCallData,
    constructorArgs,
    salt
  );
  return [deployCreate, deployProxy, upgradeProxy];
}
/**
 * @decription encodes a multicall for deploying a logic contract with deployCreate, and upgradeProxy to point to the newly deployed implementation contract
 * @param implementationBase ethers contract instance of the implementation contract
 * @param factory connected instance of alicenetFactory
 * @param ethers instance of hardhat ethers
 * @param initCallData encoded call data for the initialize function of the implementation contract
 * @param constructorArgs encoded constructor arguments
 * @param salt bytes32 formatted salt used to deploy the proxy
 * @returns
 */
export async function encodeMultiCallUpgradeProxyArgs(
  implementationBase: ContractFactory,
  factory: AliceNetFactory,
  ethers: Ethers,
  initCallData: string,
  constructorArgs: any[] = [],
  salt: string
) {
  const deployTxData = implementationBase.getDeployTransaction(
    ...constructorArgs
  ).data as BytesLike;
  const deployCreateCallData = factory.interface.encodeFunctionData(
    "deployCreate",
    [deployTxData]
  );
  const implementationContractAddress = await calculateDeployCreateAddress(
    factory.address,
    ethers
  );
  const upgradeProxyCallData = factory.interface.encodeFunctionData(
    "upgradeProxy",
    [salt, implementationContractAddress, initCallData]
  );
  const deployCreate = encodeMultiCallArgs(
    factory.address,
    0,
    deployCreateCallData
  );
  const upgradeProxy = encodeMultiCallArgs(
    factory.address,
    0,
    upgradeProxyCallData
  );
  return [deployCreate, upgradeProxy];
}

export function encodeMultiCallArgs(
  targetAddress: string,
  value: BigNumberish,
  callData: BytesLike
): MultiCallArgsStruct {
  const output: MultiCallArgsStruct = {
    target: targetAddress,
    value,
    data: callData,
  };
  return output;
}

export async function calculateDeployCreateAddress(
  deployerAddress: string,
  ethers: Ethers
) {
  const factoryNonce = await ethers.provider.getTransactionCount(
    deployerAddress
  );
  return ethers.utils.getContractAddress({
    from: deployerAddress,
    nonce: factoryNonce,
  });
}

export async function deployUpgradeableGasSafe(
  contractName: string,
  factory: AliceNetFactory,
  ethers: Ethers,
  initCallData: string,
  constructorArgs: any[],
  salt: string,
  waitConfirmantions: number = 0,
  overrides?: Overrides & { from?: PromiseOrValue<string> }
) {
  const ImplementationBase = await ethers.getContractFactory(contractName);
  try {
    return await multiCallDeployUpgradeable(
      ImplementationBase,
      factory,
      ethers,
      initCallData,
      constructorArgs,
      salt,
      overrides
    );
  } catch (err) {
    if (err instanceof MultiCallGasError) {
      return deployUpgradeable(
        contractName,
        factory,
        ethers,
        initCallData,
        constructorArgs,
        salt,
        waitConfirmantions,
        overrides
      );
    }
    throw err;
  }
}
/**
 * @description attempts to upgrade a proxy using a multicall deploycreate and upgradeProxy,
 * if the gas is too high, it will deploy the implementation contract and upgrade the proxy with 2 separate calls
 * @param contractName name of the contract to deploy
 * @param factory connected instance of AliceNetFactory
 * @param ethers instance of ethers js
 * @param initCallData encoded inititalize call data
 * @param constructorArgs constructor arguments (can only be used for immutable variables)
 * @param salt bytes32 formatted salt used to deploy the proxy
 * @param waitConfirmations
 * @param overrides
 * @returns
 */
export async function upgradeProxyGasSafe(
  contractName: string,
  factory: AliceNetFactory,
  ethers: Ethers,
  initCallData: string,
  constructorArgs: any[],
  salt: string,
  waitConfirmations: number = 0,
  overrides?: Overrides & { from?: PromiseOrValue<string> }
) {
  const ImplementationBase = await ethers.getContractFactory(contractName);
  try {
    return await multiCallUpgradeProxy(
      ImplementationBase,
      factory,
      ethers,
      initCallData,
      constructorArgs,
      salt,
      overrides
    );
  } catch (err) {
    if (err instanceof MultiCallGasError) {
      const txResponse = await deployCreate(
        ImplementationBase,
        factory,
        ethers,
        constructorArgs,
        overrides
      );

      const receipt = await txResponse.wait(waitConfirmations);
      const logicAddress = getEventVar(
        receipt,
        EVENT_DEPLOYED_RAW,
        CONTRACT_ADDR
      );
      return upgradeProxy(logicAddress, factory, initCallData, salt, overrides);
    }
    throw err;
  }
}

/**
 * @param contract name of the contract to deploy, or a contract factory
 * @param factory instance of deployed and connected alicenetFactory
 * @param ethers ethers js object
 * @param constructorArgs constructor arguments for the implementation contract
 * @param overrides
 * @returns a promise that resolves to a transaction response
 */
export async function deployCreate(
  contract: string | ContractFactory,
  factory: AliceNetFactory,
  ethers: Ethers,
  constructorArgs: any[] = [],
  overrides?: Overrides & { from?: PromiseOrValue<string> }
) {
  const implementationBase =
    contract instanceof ContractFactory
      ? contract
      : await ethers.getContractFactory(contract);

  const deployTxData = implementationBase.getDeployTransaction(
    ...constructorArgs
  ).data as BytesLike;
  if (overrides === undefined) {
    return await factory.deployCreate(deployTxData);
  } else {
    return await factory.deployCreate(deployTxData, overrides);
  }
}

export async function deployCreate2(
  contractName: string,
  factory: AliceNetFactory,
  ethers: Ethers,
  constructorArgs: any[] = [],
  salt: string,
  overrides?: Overrides & { from?: PromiseOrValue<string> }
) {
  const implementationBase = await ethers.getContractFactory(contractName);
  const deployTxData = implementationBase.getDeployTransaction(
    ...constructorArgs
  ).data as BytesLike;
  if (overrides === undefined) {
    return await factory.deployCreate2(0, salt, deployTxData);
  } else {
    return await factory.deployCreate2(0, deployTxData, salt, overrides);
  }
}

/**
 * @description deploys logic contract with deployCreate, then multiCalls deployProxy and upgradeProxy
 * @param contract name of the contract to deploy, or a contract factory
 * @param factory instance of deployed and connected alicenetFactory
 * @param ethers ethers js object
 * @param initCallData encoded call data for the initialize function of the implementation contract
 * @param constructorArgs constructor arguments for the implementation contract
 * @param salt bytes32 formatted salt used to deploy the proxy
 * @param waitConfirmation number of confirmations to wait for before returning the transaction
 * @param overrides
 * @returns
 */
export async function deployUpgradeable(
  contract: string | ContractFactory,
  factory: AliceNetFactory,
  ethers: Ethers,
  initCallData: string,
  constructorArgs: Array<string>,
  salt: string,
  waitConfirmation: number = 0,
  overrides?: Overrides & { from?: PromiseOrValue<string> }
) {
  const txResponse = await deployCreate(
    contract,
    factory,
    ethers,
    constructorArgs,
    overrides
  );
  const receipt = await txResponse.wait(waitConfirmation);
  const implementationContractAddress = await getEventVar(
    receipt,
    "DeployedRaw",
    "contractAddr"
  );

  // use mutlticall to deploy proxy and upgrade proxy
  const multiCallArgs = await encodeMultiCallDeployProxyAndUpgradeProxyArgs(
    implementationContractAddress,
    factory,
    initCallData,
    salt
  );
  if (overrides === undefined) {
    return await factory.multiCall(multiCallArgs);
  } else {
    return await factory.multiCall(multiCallArgs, overrides);
  }
}

export async function deployCreateAndRegister(
  contractName: string,
  factory: AliceNetFactory,
  ethers: Ethers,
  constructorArgs: any[],
  salt: string,
  overrides?: Overrides & { from?: PromiseOrValue<string> }
): Promise<ContractTransaction> {
  const logicContract: any = await ethers.getContractFactory(contractName);
  // if not constructor ars is provide and empty array is used to indicate no constructor args
  // encode deployBcode,
  const deployTxData = logicContract.getDeployTransaction(...constructorArgs)
    .data as BytesLike;
  if (overrides === undefined) {
    return await factory.deployCreateAndRegister(deployTxData, salt);
  } else {
    return await factory.deployCreateAndRegister(deployTxData, salt, overrides);
  }
}

/**
 * @description deploys logic contract with deployCreate, then upgradeProxy with the logic contract address
 * @param logicAddress address of the logic contract
 * @param factory connected instance of AlicenetFactory
 * @param initCallData encoded call data for the initialize function of the implementation contract
 * @param salt bytes32 formatted salt used to deploy the proxy
 * @param overrides tx detail overrides
 * @returns
 */
export async function upgradeProxy(
  logicAddress: string,
  factory: AliceNetFactory,
  initCallData: string,
  salt: string,
  overrides?: Overrides & { from?: PromiseOrValue<string> }
) {
  // upgrade the proxy
  if (overrides === undefined) {
    return await factory.upgradeProxy(salt, logicAddress, initCallData);
  } else {
    return await factory.upgradeProxy(
      salt as BytesLike,
      logicAddress,
      initCallData,
      overrides
    );
  }
}

/**
 * @description returns everything on the left side of the :
 * ie: src/proxy/Proxy.sol:Mock => src/proxy/Proxy.sol
 * @param qualifiedName the relative path of the contract file + ":" + name of contract
 * @returns the relative path of the contract
 */
export function extractPath(qualifiedName: string) {
  return qualifiedName.split(":")[0];
}

/**
 * @description goes through the receipt from the
 * transaction and extract the specified event name and variable
 * @param receipt tx object returned from the tran
 * @param eventName
 * @param varName
 * @returns
 */
export function getEventVar(
  receipt: ContractReceipt,
  eventName: string,
  varName: string
) {
  let result = "0x";
  if (receipt.events !== undefined) {
    const events = receipt.events;
    for (let i = 0; i < events.length; i++) {
      // look for the event
      if (events[i].event === eventName) {
        if (events[i].args !== undefined) {
          const args = events[i].args;
          // extract the deployed mock logic contract address from the event
          result = args !== undefined ? args[varName] : undefined;
          if (result !== undefined) {
            return result;
          }
        } else {
          throw new Error(
            `failed to extract ${varName} from event: ${eventName}`
          );
        }
      }
    }
  }
  throw new Error(`failed to find event: ${eventName}`);
}

/**
 *
 * @param factoryAddress address of the factory that deployed the contract
 * @param salt value specified by custom:salt in the contrac
 * @param ethers ethersjs object
 * @returns returns the address of the metamorphic contract deployed with the following metamorphic code "0x6020363636335afa1536363636515af43d36363e3d36f3"
 */
export function getMetamorphicAddress(
  factoryAddress: string,
  salt: string,
  ethers: Ethers
) {
  const initCode = "0x6020363636335afa1536363636515af43d36363e3d36f3";
  return ethers.utils.getCreate2Address(
    factoryAddress,
    salt,
    ethers.utils.keccak256(initCode)
  );
}
