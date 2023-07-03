import { TestERC20A } from "../typechain/TestERC20A";
import { TestERC20B } from "../typechain/TestERC20B";
import { TestERC20C } from "../typechain/TestERC20C";
import { ethers } from "hardhat";
import { LiquidInfrastructureNFT } from "../typechain/LiquidInfrastructureNFT";
import { Signer } from "ethers";

type DeployContractsOptions = {
  corruptSig?: boolean;
};

export async function deployContracts(signer?: Signer | undefined) {

  const TestERC20A = await ethers.getContractFactory("TestERC20A", signer);
  const testERC20A = (await TestERC20A.deploy()) as TestERC20A;

  const TestERC20B = await ethers.getContractFactory("TestERC20B", signer);
  const testERC20B = (await TestERC20B.deploy()) as TestERC20B;

  const TestERC20C = await ethers.getContractFactory("TestERC20C", signer);
  const testERC20C = (await TestERC20C.deploy()) as TestERC20C;

  return { testERC20A, testERC20B, testERC20C };
}

export async function deployLiquidAccount(account: string) {
  const LiquidAccount = await ethers.getContractFactory("LiquidInfrastructureNFT");
  return (await LiquidAccount.deploy(account)) as LiquidInfrastructureNFT;
}

export async function liquidAccountAsNewOwner(nftAddress: string, newOwner: Signer) {
  return await ethers.getContractAt("LiquidInfrastructureNFT", nftAddress, newOwner) as LiquidInfrastructureNFT;
}