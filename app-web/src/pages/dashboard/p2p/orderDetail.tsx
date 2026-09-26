import { useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useSelector } from 'react-redux';
import Header from '../../../components/header';
import {
  useGetP2POrderQuery,
  useGetOpenP2PDisputeForOrderQuery,
  useAcceptP2POrderMutation,
  useRejectP2POrderMutation,
  useCancelP2POrderMutation,
  useMerchantCancelP2POrderMutation,
  useEscrowDepositBuildMutation,
  useEscrowDepositCommitMutation,
  useRegenerateEscrowShortlinkMutation,
  useMarkP2PPaymentSentMutation,
  useConfirmP2PPaymentReceivedMutation,
  useMerchantConfirmsP2PPaymentMutation,
  useBuyerConfirmsP2PNotPaidMutation,
  useGetCustomerP2PPerformanceQuery,
  P2PCreds,
} from '../../../store/api/p2pApis';
import { useP2PIdentity } from '../../../hooks/useP2PIdentity';
import { isCustomer, isMerchant, isAssetDepositor, isFiatPayer, isFiatRecipient, P2POrder } from '../../../types/p2p';
import { P2PListCard, P2PEmptyState, P2PStatusPill, P2PBanner, p2p, statusLabel } from '../../../components/p2p/P2PTheme';
import { signBase64Txn } from '../../../utils/trovoSDK';
import { getExplorerBaseUrl } from '../../../utils/utilities';
import { RootState } from '../../../store/reduxStore';

const STEPS: P2POrder['orderStatus'][] = [
  'AWAITING_APPROVAL',
  'AWAITING_ESCROW_DEPOSIT',
  'AWAITING_PAYMENT',
  'AWAITING_PAYMENT_CONFIRMATION',
  'COMPLETED',
];

// P2POrderDetail is the hub page (Plan Sections 96.4-96.9, consolidated,
// same as app-mobile): the order timeline is the feature's emotional
// centerpiece (Plan Section 97.2), with the one relevant action for the
// order's current state surfaced directly below it.
export default function P2POrderDetail() {
  const { orderId } = useParams();
  const navigate = useNavigate();
  const { creds, username, ready } = useP2PIdentity();
  const [busy, setBusy] = useState(false);
  const [actionError, setActionError] = useState('');

  const { data: order, isLoading, refetch } = useGetP2POrderQuery({ creds, orderId: orderId! }, { skip: !ready || !orderId });
  const { data: dispute } = useGetOpenP2PDisputeForOrderQuery(
    { creds, orderId: orderId! },
    { skip: !ready || !orderId || !order?.isDisputed },
  );

  const walletMode = useSelector((state: RootState) => state.appState!.walletMode);
  const [accept] = useAcceptP2POrderMutation();
  const [reject] = useRejectP2POrderMutation();
  const [cancel] = useCancelP2POrderMutation();
  const [merchantCancel] = useMerchantCancelP2POrderMutation();
  const [buildDeposit] = useEscrowDepositBuildMutation();
  const [commitDeposit] = useEscrowDepositCommitMutation();
  const [regenerateShortlink] = useRegenerateEscrowShortlinkMutation();
  const [markPaymentSent] = useMarkP2PPaymentSentMutation();
  const [confirmPayment] = useConfirmP2PPaymentReceivedMutation();
  const [merchantConfirms] = useMerchantConfirmsP2PPaymentMutation();
  const [buyerConfirmsNotPaid] = useBuyerConfirmsP2PNotPaidMutation();

  if (isLoading) return <div className="p-8 text-center text-gray-400">Loading...</div>;
  if (!order) return <P2PEmptyState message="This order could not be found." />;

  const run = async (action: () => Promise<any>) => {
    setBusy(true);
    setActionError('');
    try {
      await action();
      refetch();
    } catch (e: any) {
      setActionError(e?.data?.message || 'Please try again.');
    } finally {
      setBusy(false);
    }
  };

  const depositFromOwnWallet = async () => {
    if (!window.confirm('Confirm you want to sign this escrow deposit?')) return;
    setBusy(true);
    setActionError('');
    try {
      const built = await buildDeposit({ creds, orderId: order.id }).unwrap();
      if (!built.transaction) {
        setActionError('Could not build the deposit transaction.');
        return;
      }
      const signature = signBase64Txn(creds.secretKey, built.transaction, built.networkPassPhrase || '');
      await commitDeposit({ creds, orderId: order.id, transaction: built.transaction, transactionSignature: signature }).unwrap();
      refetch();
    } catch (e: any) {
      setActionError(e?.data?.message || 'Deposit failed. Please try again.');
    } finally {
      setBusy(false);
    }
  };

  const currentIndex = STEPS.indexOf(order.orderStatus);
  const terminalNonSuccess = ['REJECTED', 'CANCELLED', 'EXPIRED'].includes(order.orderStatus);

  return (
    <div className="flex flex-col space-y-4 p-3 max-w-2xl">
      <Header />
      <P2PListCard>
        <div className="flex items-center justify-between">
          <p className="text-lg font-bold">
            {order.specifiedAssetAmount} {order.asset}
          </p>
          <P2PStatusPill status={order.orderStatus} disputed={order.isDisputed} />
        </div>
        <p className="text-gray-500">
          {order.paymentAmount} {order.currency}
        </p>
        <p className="text-xs text-gray-400">
          {isCustomer(order, username) ? `Merchant: ${order.merchantUsername}` : `Customer: ${order.customerUsername}`}
        </p>
      </P2PListCard>

      {order.isDisputed && <P2PBanner tone="danger">A dispute is open on this order.</P2PBanner>}

      <P2PListCard>
        {terminalNonSuccess ? (
          <p className="font-bold text-gray-500">{statusLabel(order.orderStatus)}</p>
        ) : (
          <div className="space-y-3">
            {STEPS.map((step, i) => {
              const done = currentIndex >= 0 && i < currentIndex;
              const current = i === currentIndex;
              const color = done ? p2p.success : current ? p2p.brandDark : '#D1D5DB';
              return (
                <div key={step} className="flex items-center gap-3">
                  <span className="w-3 h-3 rounded-full" style={{ backgroundColor: color }} />
                  <span style={{ color: current ? p2p.brandDark : '#374151', fontWeight: current ? 700 : 400 }}>
                    {statusLabel(step!)}
                  </span>
                </div>
              );
            })}
          </div>
        )}
      </P2PListCard>

      {actionError && <div className="text-red-600">{actionError}</div>}

      {order.orderStatus === 'AWAITING_APPROVAL' && (
        isMerchant(order, username) ? (
          <>
            <CustomerPerformanceBadge creds={creds} customerId={order.customerUserId} ready={ready} />
            <div className="flex gap-3">
              <ActionButton label="Accept order" onClick={() => run(() => accept({ creds, orderId: order.id }).unwrap())} busy={busy} />
              <ActionButton label="Reject" secondary onClick={() => run(() => reject({ creds, orderId: order.id }).unwrap())} busy={busy} />
            </div>
          </>
        ) : (
          <>
            <p className="text-gray-500">Waiting for the merchant to accept your order.</p>
            <ActionButton label="Cancel order" secondary onClick={() => run(() => cancel({ creds, orderId: order.id }).unwrap())} busy={busy} />
          </>
        )
      )}

      {order.orderStatus === 'AWAITING_ESCROW_DEPOSIT' && (
        isAssetDepositor(order, username) ? (
          <>
            <p className="text-gray-500">
              Deposit {order.sellerEscrowAssetAmount} {order.asset} into escrow to continue.
            </p>
            <div className="flex gap-3">
              <ActionButton label="Deposit from my wallet" onClick={depositFromOwnWallet} busy={busy} />
              {order.escrowDepositShortlink ? (
                <ActionButton
                  label="Share deposit link"
                  secondary
                  onClick={() => navigate(`/dashboard/p2p/order/${order.id}/escrow-share`)}
                  busy={busy}
                />
              ) : (
                <ActionButton
                  label="Generate deposit link"
                  secondary
                  onClick={() => run(() => regenerateShortlink({ creds, orderId: order.id }).unwrap())}
                  busy={busy}
                />
              )}
            </div>
          </>
        ) : (
          <p className="text-gray-500">Waiting for the escrow deposit.</p>
        )
      )}

      {order.orderStatus === 'AWAITING_ESCROW_DEPOSIT' && isMerchant(order, username) && (
        <ActionButton
          label="Cancel order"
          secondary
          onClick={() => run(() => merchantCancel({ creds, orderId: order.id }).unwrap())}
          busy={busy}
        />
      )}

      {order.orderStatus === 'AWAITING_PAYMENT' && (
        <>
          {isFiatPayer(order, username) ? (
            <>
              <P2PListCard>
                <p className="font-bold mb-2">Pay the merchant using</p>
                <p>
                  {order.paymentMethodSnapshot?.paymentChannel} - {order.paymentMethodSnapshot?.provider}
                </p>
                <p className="font-semibold">{order.paymentMethodSnapshot?.account}</p>
              </P2PListCard>
              <ActionButton label="I have sent payment" onClick={() => run(() => markPaymentSent({ creds, orderId: order.id }).unwrap())} busy={busy} />
            </>
          ) : (
            <p className="text-gray-500">Escrow confirmed. Waiting for the buyer's payment.</p>
          )}
          {!order.isDisputed && (
            <DisputeLink onClick={() => navigate(`/dashboard/p2p/order/${order.id}/dispute`)} />
          )}
        </>
      )}

      {order.orderStatus === 'AWAITING_PAYMENT_CONFIRMATION' && (
        <>
          {isFiatRecipient(order, username) ? (
            <>
              <p className="text-gray-500">The buyer marked payment as sent. Confirm receipt to release the asset.</p>
              <ActionButton label="Confirm payment received" onClick={() => run(() => confirmPayment({ creds, orderId: order.id }).unwrap())} busy={busy} />
            </>
          ) : (
            <p className="text-gray-500">Waiting for the seller to confirm your payment.</p>
          )}
          {!order.isDisputed && (
            <DisputeLink onClick={() => navigate(`/dashboard/p2p/order/${order.id}/dispute`)} />
          )}
        </>
      )}

      {order.orderStatus === 'COMPLETED' && (
        <>
          <p className="text-gray-500">This order is complete.</p>
          {order.assetReleaseTransactionHash && (
            <TxHash label="Asset release" hash={order.assetReleaseTransactionHash} walletMode={walletMode} />
          )}
          {order.escrowDepositTransactionHash && (
            <TxHash label="Escrow deposit" hash={order.escrowDepositTransactionHash} walletMode={walletMode} />
          )}
        </>
      )}

      {order.isDisputed && dispute && (
        <P2PListCard>
          <p className="font-bold mb-2">Resolve this dispute</p>
          {isMerchant(order, username) && (
            <ActionButton
              label="I actually received the payment"
              onClick={() => run(() => merchantConfirms({ creds, disputeId: dispute.id }).unwrap())}
              busy={busy}
            />
          )}
          {isCustomer(order, username) && (
            <ActionButton
              label="I have not actually paid yet"
              secondary
              onClick={() => run(() => buyerConfirmsNotPaid({ creds, disputeId: dispute.id }).unwrap())}
              busy={busy}
            />
          )}
        </P2PListCard>
      )}
    </div>
  );
}

// CustomerPerformanceBadge lets a merchant see the customer's trading
// history before deciding to accept their order (Plan Section 8/26).
function CustomerPerformanceBadge({ creds, customerId, ready }: { creds: P2PCreds; customerId: string; ready: boolean }) {
  const { data: perf } = useGetCustomerP2PPerformanceQuery({ creds, customerId }, { skip: !ready || !customerId });
  if (!perf || perf.completedTrades === 0) return null;
  return (
    <p className="text-xs text-gray-400">
      This customer has completed {perf.completedTrades} trade{perf.completedTrades === 1 ? '' : 's'} ({perf.completionRate}%
      completion rate).
    </p>
  );
}

function ActionButton({ label, onClick, busy, secondary }: { label: string; onClick: () => void; busy: boolean; secondary?: boolean }) {
  return (
    <button
      disabled={busy}
      onClick={onClick}
      className="flex-1 py-3 rounded-2xl font-semibold disabled:opacity-50"
      style={
        secondary
          ? { backgroundColor: 'white', color: p2p.brandDark, border: `1px solid ${p2p.brandDark}` }
          : { backgroundColor: p2p.brandDark, color: 'white' }
      }
    >
      {label}
    </button>
  );
}

function DisputeLink({ onClick }: { onClick: () => void }) {
  return (
    <button onClick={onClick} className="font-semibold underline" style={{ color: p2p.danger }}>
      Raise a dispute
    </button>
  );
}

function TxHash({ label, hash, walletMode }: { label: string; hash: string; walletMode: string }) {
  return (
    <p className="text-sm">
      <span className="text-gray-500">{label}: </span>
      <a
        href={`${getExplorerBaseUrl(walletMode)}${hash}`}
        target="_blank"
        rel="noopener noreferrer"
        className="font-mono underline"
        style={{ color: p2p.brandDark }}
      >
        {hash.substring(0, 10)}...{hash.substring(hash.length - 6)}
      </a>
    </p>
  );
}
