import { useNavigate, useParams } from 'react-router-dom';
import { useGetP2POrderQuery } from '../../../store/api/p2pApis';
import { useP2PIdentity } from '../../../hooks/useP2PIdentity';
import { p2p } from '../../../components/p2p/P2PTheme';

// P2PEscrowShare (Plan Section 96.6/97.7): the feature's "money moment" -
// the page most likely seen by a third party who isn't even a Trovo user
// yet (Plan Section 33), so it gets the most visual care in the whole
// feature.
export default function P2PEscrowShare() {
  const { orderId } = useParams();
  const navigate = useNavigate();
  const { creds, ready } = useP2PIdentity();
  const { data: order } = useGetP2POrderQuery({ creds, orderId: orderId! }, { skip: !ready || !orderId });

  if (!order) return <div className="p-8 text-center text-gray-400">Loading...</div>;

  return (
    <div className="min-h-screen flex flex-col items-center justify-center p-8" style={{ backgroundColor: p2p.brandDark }}>
      <p className="text-white text-3xl font-bold">
        {order.sellerEscrowAssetAmount} {order.asset}
      </p>
      <p className="text-white/70 mt-1 mb-8 text-center">
        Anyone with this link or QR can fund this deposit.
      </p>

      {order.escrowDepositQrCode && (
        <div className="bg-white rounded-3xl p-4 mb-6">
          <img src={order.escrowDepositQrCode} width={220} height={220} alt="Escrow deposit QR" />
        </div>
      )}

      <div className="bg-white rounded-2xl p-4 flex items-center gap-3 w-full max-w-md mb-4">
        <span className="flex-1 text-sm truncate">{order.escrowDepositShortlink}</span>
        <button
          onClick={() => {
            navigator.clipboard.writeText(order.escrowDepositShortlink || '');
          }}
          className="text-primary-800 font-semibold"
        >
          Copy
        </button>
      </div>

      <button
        onClick={() => {
          if (navigator.share) {
            navigator.share({ url: order.escrowDepositShortlink || '' });
          } else {
            navigator.clipboard.writeText(order.escrowDepositShortlink || '');
          }
        }}
        className="w-full max-w-md py-4 rounded-2xl font-semibold mb-3"
        style={{ backgroundColor: 'white', color: p2p.brandDark }}
      >
        Share
      </button>
      <button onClick={() => navigate(`/dashboard/p2p/order/${order.id}`)} className="text-white/70">
        Back to order
      </button>
    </div>
  );
}
