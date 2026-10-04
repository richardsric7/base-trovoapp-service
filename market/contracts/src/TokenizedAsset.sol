// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {ERC20Burnable} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

/// @title TokenizedAsset
/// @notice The token of one tokenized asset: a plain mintable, burnable
/// ERC-20 with its own decimals, owned by the asset's issuing Safe (the
/// only account that can mint). Who may end up holding it is enforced by
/// the platform: purchases go through TrovoOfferBook, which needs a
/// platform authorization for each fill.
contract TokenizedAsset is ERC20Burnable, Ownable {
    uint8 private immutable _customDecimals;

    constructor(string memory name_, string memory symbol_, uint8 decimals_, address initialOwner)
        ERC20(name_, symbol_)
        Ownable(initialOwner)
    {
        _customDecimals = decimals_;
    }

    function decimals() public view override returns (uint8) {
        return _customDecimals;
    }

    function mint(address to, uint256 amount) external onlyOwner {
        _mint(to, amount);
    }
}
