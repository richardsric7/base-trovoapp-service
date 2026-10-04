// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.28;

import "@account-abstraction/contracts/core/BasePaymaster.sol";
import "@account-abstraction/contracts/core/Helpers.sol";
import "@account-abstraction/contracts/core/UserOperationLib.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";
import "@openzeppelin/contracts/utils/math/Math.sol";

/**
 * @title TrovoTokenPaymaster
 * @notice ERC-4337 (EntryPoint v0.7) paymaster that lets Trovo wallets pay
 *         their own gas in a curated stablecoin (USDC, USDT, cNGN, ...).
 *
 * The paymaster pays the EntryPoint in ETH from its deposit and charges the
 * wallet the equivalent in the chosen token, at an exchange rate quoted and
 * signed off-chain by the Trovo quote service (which applies Trovo's spread).
 *
 * paymasterAndData layout:
 *
 *   [  0: 20] paymaster address          (EntryPoint v0.7 static fields)
 *   [ 20: 36] paymaster verification gas
 *   [ 36: 52] paymaster postOp gas
 *   [ 52: 72] token                      (the stablecoin paying for gas)
 *   [ 72: 78] validUntil                 (uint48, 0 = no expiry)
 *   [ 78: 84] validAfter                 (uint48)
 *   [ 84:116] exchangeRate               (uint256: token base units per 1e18 wei, spread included)
 *   [116:181] quote signature            (65 bytes, eth_sign over getHash())
 *
 * The quote is bound to the whole UserOperation (except its signature), so a
 * quote can only pay for the exact operation it was issued for, and the
 * wallet owner signs the rate as part of the userOpHash.
 *
 * Charging:
 *   - When the wallet is already deployed and has approved at least the
 *     maximum token cost, the maximum is pulled during validation and the
 *     unused part refunded in postOp (the normal case).
 *   - Otherwise (the activation operation, whose calldata also approves this
 *     paymaster), validation only checks the balance, and the actual cost is
 *     pulled in postOp, after the calldata has run. If that pull fails the
 *     cost is recorded as debt, and the wallet cannot use the paymaster again
 *     until it is settled (settleDebt) or forgiven by the owner.
 *
 * Safety bounds against a compromised quote signer: per-token min/max rates
 * set by the owner (a Safe), several quote signers for rotation, and a
 * pauser that can stop the paymaster immediately.
 */
contract TrovoTokenPaymaster is BasePaymaster, Pausable {
    using SafeERC20 for IERC20;
    using UserOperationLib for PackedUserOperation;

    uint256 private constant TOKEN_OFFSET = PAYMASTER_DATA_OFFSET; // 52
    uint256 private constant VALID_UNTIL_OFFSET = TOKEN_OFFSET + 20; // 72
    uint256 private constant VALID_AFTER_OFFSET = VALID_UNTIL_OFFSET + 6; // 78
    uint256 private constant RATE_OFFSET = VALID_AFTER_OFFSET + 6; // 84
    uint256 private constant SIGNATURE_OFFSET = RATE_OFFSET + 32; // 116
    uint256 private constant SIGNATURE_LENGTH = 65;

    /// @dev exchangeRate is token base units per RATE_DENOMINATOR wei (1 ETH).
    uint256 public constant RATE_DENOMINATOR = 1e18;

    struct TokenConfig {
        bool enabled;
        uint256 minRate;
        uint256 maxRate;
    }

    /// @notice Gas tokens and the rate bounds a quote must fall within.
    mapping(address => TokenConfig) public tokens;

    /// @notice Addresses whose signatures are accepted on quotes.
    mapping(address => bool) public quoteSigners;

    /// @notice May pause the paymaster (only the owner can unpause).
    address public pauser;

    /// @notice Gas charged on top of actualGasCost to cover postOp itself.
    uint256 public postOpGasOverhead;

    /// @notice Unpaid gas, per wallet and token, from failed postOp charges.
    mapping(address => mapping(address => uint256)) public debt;

    /// @notice Number of tokens each wallet owes debt in.
    mapping(address => uint256) public debtTokenCount;

    event TokenConfigured(address indexed token, bool enabled, uint256 minRate, uint256 maxRate);
    event QuoteSignerSet(address indexed signer, bool allowed);
    event PauserSet(address indexed pauser);
    event PostOpGasOverheadSet(uint256 overhead);
    event GasPaidInToken(
        address indexed sender,
        address indexed token,
        bytes32 indexed userOpHash,
        uint256 exchangeRate,
        uint256 tokenCost,
        uint256 actualGasCost
    );
    event GasChargeFailed(address indexed sender, address indexed token, bytes32 indexed userOpHash, uint256 tokenCost);
    event DebtSettled(address indexed sender, address indexed token, uint256 amount, bool forgiven);
    event TokenWithdrawn(address indexed token, address indexed to, uint256 amount);

    error InvalidPaymasterData();
    error TokenNotEnabled(address token);
    error RateOutOfBounds(address token, uint256 rate);
    error InsufficientTokenBalance(address token, uint256 balance, uint256 required);
    error OutstandingDebt(address sender);
    error InvalidRateBounds();
    error NotPauser();
    error NoDebt();

    constructor(
        IEntryPoint entryPoint_,
        address initialOwner,
        address[] memory initialQuoteSigners,
        address pauser_,
        uint256 postOpGasOverhead_
    ) BasePaymaster(entryPoint_) {
        for (uint256 i = 0; i < initialQuoteSigners.length; i++) {
            _setQuoteSigner(initialQuoteSigners[i], true);
        }
        pauser = pauser_;
        emit PauserSet(pauser_);
        postOpGasOverhead = postOpGasOverhead_;
        emit PostOpGasOverheadSet(postOpGasOverhead_);
        if (initialOwner != msg.sender) {
            _transferOwnership(initialOwner);
        }
    }

    // ---------------------------------------------------------------------
    // Quote hashing
    // ---------------------------------------------------------------------

    /**
     * @notice The hash a quote signer signs (with eth_sign / personal_sign)
     *         to authorise `token` at `exchangeRate` for `userOp`. Covers
     *         every UserOperation field except the signature, and the
     *         paymaster gas limits.
     */
    function getHash(
        PackedUserOperation calldata userOp,
        address token,
        uint48 validUntil,
        uint48 validAfter,
        uint256 exchangeRate
    ) public view returns (bytes32) {
        bytes32 opHash = keccak256(
            abi.encode(
                userOp.getSender(),
                userOp.nonce,
                keccak256(userOp.initCode),
                keccak256(userOp.callData),
                userOp.accountGasLimits,
                uint256(bytes32(userOp.paymasterAndData[PAYMASTER_VALIDATION_GAS_OFFSET:PAYMASTER_DATA_OFFSET])),
                userOp.preVerificationGas,
                userOp.gasFees
            )
        );
        return keccak256(abi.encode(opHash, block.chainid, address(this), token, validUntil, validAfter, exchangeRate));
    }

    /// @notice Splits the paymaster-specific part of paymasterAndData.
    function parsePaymasterAndData(bytes calldata paymasterAndData)
        public
        pure
        returns (address token, uint48 validUntil, uint48 validAfter, uint256 exchangeRate, bytes calldata signature)
    {
        if (paymasterAndData.length != SIGNATURE_OFFSET + SIGNATURE_LENGTH) revert InvalidPaymasterData();
        token = address(bytes20(paymasterAndData[TOKEN_OFFSET:VALID_UNTIL_OFFSET]));
        validUntil = uint48(bytes6(paymasterAndData[VALID_UNTIL_OFFSET:VALID_AFTER_OFFSET]));
        validAfter = uint48(bytes6(paymasterAndData[VALID_AFTER_OFFSET:RATE_OFFSET]));
        exchangeRate = uint256(bytes32(paymasterAndData[RATE_OFFSET:SIGNATURE_OFFSET]));
        signature = paymasterAndData[SIGNATURE_OFFSET:];
    }

    /// @notice Token cost of `weiAmount` at `exchangeRate`, rounded up.
    function tokenCost(uint256 weiAmount, uint256 exchangeRate) public pure returns (uint256) {
        return Math.mulDiv(weiAmount, exchangeRate, RATE_DENOMINATOR, Math.Rounding.Ceil);
    }

    // ---------------------------------------------------------------------
    // ERC-4337
    // ---------------------------------------------------------------------

    struct Quote {
        address token;
        uint48 validUntil;
        uint48 validAfter;
        uint256 exchangeRate;
    }

    function _validatePaymasterUserOp(PackedUserOperation calldata userOp, bytes32 userOpHash, uint256 maxCost)
        internal
        override
        whenNotPaused
        returns (bytes memory context, uint256 validationData)
    {
        Quote memory q;
        bytes calldata signature;
        (q.token, q.validUntil, q.validAfter, q.exchangeRate, signature) = parsePaymasterAndData(userOp.paymasterAndData);

        TokenConfig memory cfg = tokens[q.token];
        if (!cfg.enabled) revert TokenNotEnabled(q.token);
        if (q.exchangeRate < cfg.minRate || q.exchangeRate > cfg.maxRate) revert RateOutOfBounds(q.token, q.exchangeRate);
        if (debtTokenCount[userOp.getSender()] != 0) revert OutstandingDebt(userOp.getSender());

        if (!_isQuoteSigner(getHash(userOp, q.token, q.validUntil, q.validAfter, q.exchangeRate), signature)) {
            return ("", _packValidationData(true, q.validUntil, q.validAfter));
        }
        context = _preCharge(userOp, userOpHash, q, maxCost);
        validationData = _packValidationData(false, q.validUntil, q.validAfter);
    }

    function _isQuoteSigner(bytes32 hash, bytes calldata signature) private view returns (bool) {
        (address recovered, ECDSA.RecoverError err,) =
            ECDSA.tryRecover(MessageHashUtils.toEthSignedMessageHash(hash), signature);
        return err == ECDSA.RecoverError.NoError && quoteSigners[recovered];
    }

    function _preCharge(PackedUserOperation calldata userOp, bytes32 userOpHash, Quote memory q, uint256 maxCost)
        private
        returns (bytes memory)
    {
        address sender = userOp.getSender();
        uint256 maxTokenCost = tokenCost(maxCost, q.exchangeRate);
        uint256 balance = IERC20(q.token).balanceOf(sender);
        if (balance < maxTokenCost) revert InsufficientTokenBalance(q.token, balance, maxTokenCost);

        // Pre-charge only an already-deployed wallet with enough allowance:
        // a wallet being deployed by this operation approves in its calldata.
        bool preCharged =
            userOp.initCode.length == 0 && IERC20(q.token).allowance(sender, address(this)) >= maxTokenCost;
        if (preCharged) {
            IERC20(q.token).safeTransferFrom(sender, address(this), maxTokenCost);
        }
        return abi.encode(sender, q.token, q.exchangeRate, maxTokenCost, preCharged, userOpHash);
    }

    function _postOp(PostOpMode, bytes calldata context, uint256 actualGasCost, uint256 actualUserOpFeePerGas)
        internal
        override
    {
        (address sender, address token, uint256 exchangeRate, uint256 maxTokenCost, bool preCharged, bytes32 userOpHash) =
            abi.decode(context, (address, address, uint256, uint256, bool, bytes32));

        uint256 cost = tokenCost(actualGasCost + postOpGasOverhead * actualUserOpFeePerGas, exchangeRate);
        if (cost > maxTokenCost) cost = maxTokenCost;

        if (preCharged) {
            if (maxTokenCost > cost) {
                IERC20(token).safeTransfer(sender, maxTokenCost - cost);
            }
        } else if (!_tryTransferFrom(token, sender, cost)) {
            // Never revert here: a reverted postOp still costs the deposit
            // and would also undo the wallet's action. Record the debt.
            if (debt[sender][token] == 0) debtTokenCount[sender] += 1;
            debt[sender][token] += cost;
            emit GasChargeFailed(sender, token, userOpHash, cost);
            return;
        }
        emit GasPaidInToken(sender, token, userOpHash, exchangeRate, cost, actualGasCost);
    }

    function _tryTransferFrom(address token, address from, uint256 amount) private returns (bool) {
        if (amount == 0) return true;
        (bool ok, bytes memory ret) =
            token.call(abi.encodeCall(IERC20.transferFrom, (from, address(this), amount)));
        return ok && (ret.length == 0 ? token.code.length > 0 : abi.decode(ret, (bool)));
    }

    // ---------------------------------------------------------------------
    // Debt
    // ---------------------------------------------------------------------

    /// @notice Pays `sender`'s gas debt in `token` (needs its allowance). Anyone may call.
    function settleDebt(address sender, address token) external {
        uint256 amount = debt[sender][token];
        if (amount == 0) revert NoDebt();
        _clearDebt(sender, token);
        IERC20(token).safeTransferFrom(sender, address(this), amount);
        emit DebtSettled(sender, token, amount, false);
    }

    function forgiveDebt(address sender, address token) external onlyOwner {
        uint256 amount = debt[sender][token];
        if (amount == 0) revert NoDebt();
        _clearDebt(sender, token);
        emit DebtSettled(sender, token, amount, true);
    }

    function _clearDebt(address sender, address token) private {
        debt[sender][token] = 0;
        debtTokenCount[sender] -= 1;
    }

    // ---------------------------------------------------------------------
    // Administration
    // ---------------------------------------------------------------------

    /**
     * @notice Enables/disables a gas token and bounds the rates quotes may
     *         use (token base units per 1 ETH, spread included).
     */
    function setToken(address token, bool enabled, uint256 minRate, uint256 maxRate) external onlyOwner {
        if (enabled && (minRate == 0 || minRate > maxRate)) revert InvalidRateBounds();
        tokens[token] = TokenConfig(enabled, minRate, maxRate);
        emit TokenConfigured(token, enabled, minRate, maxRate);
    }

    function setQuoteSigner(address signer, bool allowed) external onlyOwner {
        _setQuoteSigner(signer, allowed);
    }

    function _setQuoteSigner(address signer, bool allowed) private {
        quoteSigners[signer] = allowed;
        emit QuoteSignerSet(signer, allowed);
    }

    function setPauser(address pauser_) external onlyOwner {
        pauser = pauser_;
        emit PauserSet(pauser_);
    }

    function setPostOpGasOverhead(uint256 overhead) external onlyOwner {
        postOpGasOverhead = overhead;
        emit PostOpGasOverheadSet(overhead);
    }

    function pause() external {
        if (msg.sender != pauser && msg.sender != owner()) revert NotPauser();
        _pause();
    }

    function unpause() external onlyOwner {
        _unpause();
    }

    /// @notice Moves collected gas tokens to the treasury (for conversion to ETH).
    function withdrawToken(address token, address to, uint256 amount) external onlyOwner {
        IERC20(token).safeTransfer(to, amount);
        emit TokenWithdrawn(token, to, amount);
    }
}
