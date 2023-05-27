import { ethers } from "hardhat";

async function main() {

  const soul = await ethers.deployContract("Soul");

  await soul.waitForDeployment();

  console.log("Soul deployed to:", soul.target);
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
