import { useNavigate, useParams } from 'react-router-dom';
import Header from '../../../components/header';
import { useGetP2POfferQuery, useGetMerchantP2PPerformanceQuery } from '../../../store/api/p2pApis';
import { useP2PIdentity } from '../../../hooks/useP2PIdentity';
import { P2PListCard, P2PEmptyState, p2p } from '../../../components/p2p/P2PTheme';

// P2POfferDetail (Plan Section 96.2): full offer terms + a Buy/Sell CTA
// into order creation.
export default function P2POfferDetail() {
  const { offerId } = useParams();
  const navigate = useNavigate();
  const { creds, ready } = useP2PIdentity();

  const { data: offer, isLoading } = useGetP2POfferQuery({ creds, offerId: offerId! }, { skip: !ready || !offerId });
  const { data: perf } = useGetMerchantP2PPerformanceQuery(
    { creds, merchantId: offer?.merchantUserId ?? '' },
    { skip: !ready || !offer?.merchantUserId },
  );

  if (isLoading) return <div className="p-8 text-center text-gray-400">Loading...</div>;
  if (!offer) return <P2PEmptyState message="This offer could not be found." />;

  const isBuyingFromOffer = offer.offerType === 'BUY';

  return (
    <div className="flex flex-col space-y-5 p-3 max-w-2xl">
      <Header />
      <h1 className="text-xl font-bold text-primary-800 px-2">Offer details</h1>
      <P2PListCard>
        <div className="flex items-center gap-3 mb-4">
          <div
            className="w-10 h-10 rounded-full flex items-center justify-center text-white font-bold"
            style={{ backgroundColor: p2p.brandDark }}
          >
            {offer.merchantUsername?.substring(0, 1).toUpperCase()}
          </div>
          <span className="font-bold text-lg">{offer.merchantUsername}</span>
        </div>
        <Row label="Price" value={`${offer.price} ${offer.currency}`} />
        <Row label="Available" value={`${offer.availableLiquidity} ${offer.asset}`} />
        <Row label="Limits" value={`${offer.minOrderAmount} - ${offer.maxOrderAmount} ${offer.asset}`} />
        <Row label="Payment method" value={offer.paymentMethod?.paymentChannel || '-'} />
        <Row label="Provider" value={offer.paymentMethod?.provider || '-'} />
        {offer.remark && <p className="text-gray-500 mt-2">{offer.remark}</p>}
      </P2PListCard>

      {perf && perf.completedTrades > 0 && (
        <P2PListCard>
          <p className="font-bold mb-2">Merchant performance</p>
          <Row label="Completed trades" value={String(perf.completedTrades)} />
          <Row label="Completion rate" value={`${perf.completionRate}%`} />
          {perf.disputesResolvedAgainstMerchant > 0 && (
            <Row label="Disputes resolved against merchant" value={String(perf.disputesResolvedAgainstMerchant)} />
          )}
        </P2PListCard>
      )}

      <button
        className="w-full py-4 rounded-2xl font-semibold text-white"
        style={{ backgroundColor: p2p.brandDark }}
        onClick={() => navigate(`/dashboard/p2p/create-order/${offer.id}`)}
      >
        {isBuyingFromOffer ? 'Sell to this offer' : 'Buy from this offer'}
      </button>
    </div>
  );
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex justify-between py-1">
      <span className="text-gray-500">{label}</span>
      <span className="font-semibold">{value}</span>
    </div>
  );
}
