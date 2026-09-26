import Header from '../../../components/header';
import { useListMyP2PRefundsQuery, useClaimP2PRefundMutation } from '../../../store/api/p2pApis';
import { useP2PIdentity } from '../../../hooks/useP2PIdentity';
import { P2PListCard, P2PEmptyState, P2PStatusPill, p2p } from '../../../components/p2p/P2PTheme';

// P2PMyRefunds: a refund is recorded (overpayment, an order that expired
// before payment, or a dispute resolved in the buyer's favor) but funds
// only actually move back to the depositor when they claim it here - this
// page is that missing claim step (Plan Section 45).
export default function P2PMyRefunds() {
  const { creds, ready } = useP2PIdentity();
  const { data, isLoading, refetch } = useListMyP2PRefundsQuery({ creds }, { skip: !ready });
  const [claim, { isLoading: claiming }] = useClaimP2PRefundMutation();

  const doClaim = async (refundId: string) => {
    try {
      await claim({ creds, refundId }).unwrap();
      refetch();
    } catch (e) {
      // surfaced via the refund row staying unclaimed; kept minimal here
      // since there is no existing app-web toast wired into this page yet
    }
  };

  return (
    <div className="flex flex-col space-y-5 p-3 max-w-2xl">
      <Header />
      <h1 className="text-xl font-bold text-primary-800 px-2">My refunds</h1>

      {isLoading ? (
        <div className="text-center py-16 text-gray-400">Loading...</div>
      ) : !data?.data?.length ? (
        <P2PEmptyState message="You have no refunds." />
      ) : (
        data.data.map((refund) => (
          <P2PListCard key={refund.id}>
            <div className="flex items-center justify-between">
              <div>
                <p className="font-bold">
                  {refund.amount} {refund.token}
                </p>
                <p className="text-sm text-gray-500">{refund.reason}</p>
              </div>
              {refund.claimed ? (
                <P2PStatusPill status="COMPLETED" />
              ) : (
                <button
                  disabled={claiming}
                  onClick={() => doClaim(refund.id)}
                  className="px-4 py-2 rounded-xl font-semibold text-white disabled:opacity-50"
                  style={{ backgroundColor: p2p.brandDark }}
                >
                  Claim
                </button>
              )}
            </div>
          </P2PListCard>
        ))
      )}
    </div>
  );
}
