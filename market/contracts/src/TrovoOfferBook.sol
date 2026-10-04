// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {Pausable} from "@openzeppelin/contracts/utils/Pausable.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {EIP712} from "@openzeppelin/contracts/utils/cryptography/EIP712.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";

/// @title TrovoOfferBook
/// @notice Fixed-price offers between Trovo's curated regulated tokens.
///
/// A seller escrows `sellToken` in an offer and sets a price in each token
/// it accepts as payment. A buyer fills any part of it: in one transaction
/// the payment goes to the offer's proceeds recipient and the sold tokens
/// to the buyer's chosen recipient, or nothing happens.
///
/// Every fill needs an authorization signed by one of the platform's
/// authorizers (EIP-712 `Fill`), binding the taker, recipient, offer,
/// payment token, amount, maximum payment, a single-use nonce and a
/// deadline. This is how the platform keeps its compliance checks (KYC,
/// purchase caps, sale windows) in front of trades it cannot see on-chain.
/// An authorizer key can only let someone buy at the seller's own price;
/// it can never move a seller's or buyer's tokens.
///
/// Only tokens the owner lists as tradable can be sold or accepted; a
/// delisted token stops trading, while sellers can always cancel and take
/// back what is left.
contract TrovoOfferBook is Ownable, Pausable, ReentrancyGuard, EIP712 {
    using SafeERC20 for IERC20;

    struct Offer {
        address seller;
        address sellToken;
        address proceedsRecipient;
        uint256 remaining;
        bool open;
    }

    /// @dev Paying for `amount` base units of the sell token costs
    /// ceil(amount * num / den) base units of the payment token.
    struct Price {
        uint128 num;
        uint128 den;
    }

    /// @dev A purchase: `amount` base units of offer `offerId`'s sell token
    /// for at most `maxPayment` of `paymentToken`, delivered to `recipient`.
    /// `nonce` (any unused number) and `deadline` bound the authorization.
    struct FillRequest {
        uint256 offerId;
        address paymentToken;
        uint256 amount;
        uint256 maxPayment;
        address recipient;
        uint256 nonce;
        uint256 deadline;
    }

    bytes32 public constant FILL_TYPEHASH = keccak256(
        "Fill(uint256 offerId,address taker,address recipient,address paymentToken,uint256 amount,uint256 maxPayment,uint256 nonce,uint256 deadline)"
    );

    uint256 public nextOfferId = 1;
    mapping(uint256 => Offer) private _offers;
    mapping(uint256 => mapping(address => Price)) private _prices;

    mapping(address => bool) public tradable;
    mapping(address => bool) public authorizers;
    mapping(uint256 => bool) public nonceUsed;

    event TradableSet(address indexed token, bool tradable);
    event AuthorizerSet(address indexed authorizer, bool allowed);
    event OfferCreated(uint256 indexed offerId, address indexed seller, address indexed sellToken, uint256 amount, address proceedsRecipient);
    event PriceSet(uint256 indexed offerId, address indexed paymentToken, uint128 num, uint128 den);
    event OfferCancelled(uint256 indexed offerId, uint256 returned);
    event Filled(
        uint256 indexed offerId,
        address indexed taker,
        address indexed recipient,
        address paymentToken,
        uint256 amount,
        uint256 payment,
        uint256 nonce
    );

    error NotTradable(address token);
    error InvalidPrice();
    error InvalidAmount();
    error InvalidAddress();
    error NotSeller();
    error OfferClosed(uint256 offerId);
    error NotAccepted(uint256 offerId, address paymentToken);
    error InsufficientRemaining(uint256 remaining, uint256 amount);
    error PaymentTooHigh(uint256 payment, uint256 maxPayment);
    error AuthorizationExpired(uint256 deadline);
    error NonceUsed(uint256 nonce);
    error Unauthorized();

    constructor(address initialOwner, address[] memory initialAuthorizers)
        Ownable(initialOwner)
        EIP712("TrovoOfferBook", "1")
    {
        for (uint256 i = 0; i < initialAuthorizers.length; i++) {
            _setAuthorizer(initialAuthorizers[i], true);
        }
    }

    // ---------------------------------------------------------------------
    // Sellers
    // ---------------------------------------------------------------------

    /// @notice Escrows `amount` of `sellToken` from the caller in a new
    /// offer, priced in each of `paymentTokens`. Payments go to
    /// `proceedsRecipient`.
    function createOffer(
        address sellToken,
        uint256 amount,
        address proceedsRecipient,
        address[] calldata paymentTokens,
        Price[] calldata prices
    ) external whenNotPaused nonReentrant returns (uint256 offerId) {
        if (!tradable[sellToken]) revert NotTradable(sellToken);
        if (proceedsRecipient == address(0)) revert InvalidAddress();
        if (amount == 0 || paymentTokens.length == 0 || paymentTokens.length != prices.length) revert InvalidAmount();

        offerId = nextOfferId++;
        _offers[offerId] = Offer({seller: msg.sender, sellToken: sellToken, proceedsRecipient: proceedsRecipient, remaining: 0, open: true});
        for (uint256 i = 0; i < paymentTokens.length; i++) {
            if (prices[i].num == 0) revert InvalidPrice();
            _setPrice(offerId, sellToken, paymentTokens[i], prices[i]);
        }

        // escrow what actually arrives
        uint256 before = IERC20(sellToken).balanceOf(address(this));
        IERC20(sellToken).safeTransferFrom(msg.sender, address(this), amount);
        uint256 received = IERC20(sellToken).balanceOf(address(this)) - before;
        if (received == 0) revert InvalidAmount();
        _offers[offerId].remaining = received;

        emit OfferCreated(offerId, msg.sender, sellToken, received, proceedsRecipient);
    }

    /// @notice Sets (or, with num 0, withdraws) the offer's price in
    /// `paymentToken`.
    function setPrice(uint256 offerId, address paymentToken, Price calldata price) external {
        Offer storage o = _offers[offerId];
        if (o.seller != msg.sender) revert NotSeller();
        if (!o.open) revert OfferClosed(offerId);
        _setPrice(offerId, o.sellToken, paymentToken, price);
    }

    /// @notice Closes the offer and returns what is left to the seller.
    /// Works while paused or after the token is delisted.
    function cancelOffer(uint256 offerId) external nonReentrant {
        Offer storage o = _offers[offerId];
        if (o.seller != msg.sender) revert NotSeller();
        if (!o.open) revert OfferClosed(offerId);
        uint256 left = o.remaining;
        o.open = false;
        o.remaining = 0;
        if (left > 0) {
            IERC20(o.sellToken).safeTransfer(o.seller, left);
        }
        emit OfferCancelled(offerId, left);
    }

    // ---------------------------------------------------------------------
    // Buyers
    // ---------------------------------------------------------------------

    /// @notice Buys `f.amount` base units of the offer's sell token, paying
    /// in `f.paymentToken` (at most `f.maxPayment`), and sends them to
    /// `f.recipient`. The caller must have approved this contract for the
    /// payment and hold a platform authorization for exactly this fill.
    function fill(FillRequest calldata f, bytes calldata authorization)
        external
        whenNotPaused
        nonReentrant
        returns (uint256 payment)
    {
        Offer storage o = _offers[f.offerId];
        if (!o.open) revert OfferClosed(f.offerId);
        if (!tradable[o.sellToken]) revert NotTradable(o.sellToken);
        if (!tradable[f.paymentToken]) revert NotTradable(f.paymentToken);
        if (f.recipient == address(0)) revert InvalidAddress();
        if (f.amount == 0) revert InvalidAmount();
        if (f.amount > o.remaining) revert InsufficientRemaining(o.remaining, f.amount);
        if (block.timestamp > f.deadline) revert AuthorizationExpired(f.deadline);
        if (nonceUsed[f.nonce]) revert NonceUsed(f.nonce);

        Price memory p = _prices[f.offerId][f.paymentToken];
        if (p.num == 0) revert NotAccepted(f.offerId, f.paymentToken);
        payment = Math.mulDiv(f.amount, p.num, p.den, Math.Rounding.Ceil);
        if (payment > f.maxPayment) revert PaymentTooHigh(payment, f.maxPayment);

        (address signer, ECDSA.RecoverError err,) = ECDSA.tryRecover(_fillDigest(f, msg.sender), authorization);
        if (err != ECDSA.RecoverError.NoError || !authorizers[signer]) revert Unauthorized();

        nonceUsed[f.nonce] = true;
        o.remaining -= f.amount;

        IERC20(f.paymentToken).safeTransferFrom(msg.sender, o.proceedsRecipient, payment);
        IERC20(o.sellToken).safeTransfer(f.recipient, f.amount);

        emit Filled(f.offerId, msg.sender, f.recipient, f.paymentToken, f.amount, payment, f.nonce);
    }

    // ---------------------------------------------------------------------
    // Views
    // ---------------------------------------------------------------------

    function getOffer(uint256 offerId) external view returns (Offer memory) {
        return _offers[offerId];
    }

    function priceOf(uint256 offerId, address paymentToken) external view returns (Price memory) {
        return _prices[offerId][paymentToken];
    }

    /// @notice What buying `amount` of the offer costs in `paymentToken`
    /// (0 when the token is not accepted).
    function quote(uint256 offerId, address paymentToken, uint256 amount) external view returns (uint256) {
        Price memory p = _prices[offerId][paymentToken];
        if (p.num == 0) return 0;
        return Math.mulDiv(amount, p.num, p.den, Math.Rounding.Ceil);
    }

    /// @notice The EIP-712 digest an authorizer signs for `taker` to make
    /// fill `f`.
    function fillDigest(FillRequest calldata f, address taker) external view returns (bytes32) {
        return _fillDigest(f, taker);
    }

    // ---------------------------------------------------------------------
    // Owner
    // ---------------------------------------------------------------------

    function setTradable(address token, bool allowed) external onlyOwner {
        if (token == address(0)) revert InvalidAddress();
        tradable[token] = allowed;
        emit TradableSet(token, allowed);
    }

    function setAuthorizer(address authorizer, bool allowed) external onlyOwner {
        _setAuthorizer(authorizer, allowed);
    }

    function pause() external onlyOwner {
        _pause();
    }

    function unpause() external onlyOwner {
        _unpause();
    }

    // ---------------------------------------------------------------------

    function _setPrice(uint256 offerId, address sellToken, address paymentToken, Price memory price) private {
        if (price.num != 0) {
            if (!tradable[paymentToken]) revert NotTradable(paymentToken);
            if (paymentToken == sellToken || price.den == 0) revert InvalidPrice();
        }
        _prices[offerId][paymentToken] = price;
        emit PriceSet(offerId, paymentToken, price.num, price.den);
    }

    function _fillDigest(FillRequest calldata f, address taker) private view returns (bytes32) {
        return _hashTypedDataV4(
            keccak256(
                abi.encode(FILL_TYPEHASH, f.offerId, taker, f.recipient, f.paymentToken, f.amount, f.maxPayment, f.nonce, f.deadline)
            )
        );
    }

    function _setAuthorizer(address authorizer, bool allowed) private {
        if (authorizer == address(0)) revert InvalidAddress();
        authorizers[authorizer] = allowed;
        emit AuthorizerSet(authorizer, allowed);
    }
}
