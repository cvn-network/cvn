// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/draft-ERC20Permit.sol";

contract Soul is ERC20, AccessControl, ERC20Permit {
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");

    address public constant ERC20_MODULE_ADDRESS = 0x47EeB2eac350E1923b8CBDfA4396A077b36E62a0;

    bool private _transferEnabled = false;

    constructor() ERC20("SOUL", "SOUL") ERC20Permit("SOUL") {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(MINTER_ROLE, msg.sender);

        _grantRole(DEFAULT_ADMIN_ROLE, ERC20_MODULE_ADDRESS);
        _grantRole(MINTER_ROLE, ERC20_MODULE_ADDRESS);
    }

    function transferEnabled() public view returns (bool) {
        return _transferEnabled;
    }

    function enableTransfer() public onlyRole(DEFAULT_ADMIN_ROLE) {
        require(!_transferEnabled, "Soul: transfer is already enabled");
        _transferEnabled = true;
    }

    function disableTransfer() public onlyRole(DEFAULT_ADMIN_ROLE) {
        require(_transferEnabled, "Soul: transfer is already disabled");
        _transferEnabled = false;
    }

    function mint(address to, uint256 amount) public onlyRole(MINTER_ROLE) {
        _mint(to, amount);
    }

    function transfer(address to, uint256 amount) public virtual override returns (bool) {
        address owner = _msgSender();
        require(_transferEnabled && owner != ERC20_MODULE_ADDRESS, "Soul: transfer is disabled");
        _transfer(owner, to, amount);
        return true;
    }

    function transferFrom(
        address from,
        address to,
        uint256 amount
    ) public virtual override returns (bool) {
        address spender = _msgSender();
        require(_transferEnabled && spender != ERC20_MODULE_ADDRESS, "Soul: transfer is disabled");
        _spendAllowance(from, spender, amount);
        _transfer(from, to, amount);
        return true;
    }
}
