import {loadFixture,} from "@nomicfoundation/hardhat-toolbox/network-helpers";
import {expect} from "chai";
import {ethers} from "hardhat";

describe("Soul", function () {
    // We define a fixture to reuse the same setup in every test.
    // We use loadFixture to run this setup once, snapshot that state,
    // and reset Hardhat Network to that snapshot in every test.
    async function deployFixture() {
        // Contracts are deployed using the first signer/account by default
        const [owner, otherAccount] = await ethers.getSigners();

        const Soul = await ethers.getContractFactory("Soul");
        const soul = await Soul.deploy();

        await soul.mint(owner.address, 1);

        return {soul, owner, otherAccount};
    }

    describe("Deployment", function () {
        it("Should set the right token info", async function () {
            const {soul} = await loadFixture(deployFixture);

            expect(await soul.symbol()).to.equal("SOUL");

            expect(await soul.name()).to.equal("SOUL");

            expect(await soul.decimals()).to.equal("18");
        });

        it("Should set the role of owner", async function () {
            const {soul, owner} = await loadFixture(deployFixture);

            expect(await soul.hasRole(await soul.MINTER_ROLE(), owner.address)).to.equal(true);

            expect(await soul.hasRole(await soul.MINTER_ROLE(), await soul.ERC20_MODULE_ADDRESS())).to.equal(true);
        });

        it("Should disable transfer", async function () {
            const {soul} = await loadFixture(deployFixture);

            expect(await soul.transferEnabled()).to.equal(false);
        });
    });

    describe("Transfer", function () {
        it("Should not allow transfer", async function () {
            const {soul, owner, otherAccount} = await loadFixture(deployFixture);

            await expect(soul.transfer(otherAccount.address, 1)).to.be.revertedWith(
                "Soul: transfer is disabled"
            );
        });

        it("Should allow transfer after enable", async function () {
            const {soul, owner, otherAccount} = await loadFixture(deployFixture);

            await soul.enableTransfer();

            expect(await soul.transferEnabled()).to.equal(true);

            await soul.transfer(otherAccount.address, 1);

            expect(await soul.balanceOf(otherAccount.address)).to.equal(1);

            await soul.disableTransfer();

            expect(await soul.transferEnabled()).to.equal(false);
        });
    });
});
