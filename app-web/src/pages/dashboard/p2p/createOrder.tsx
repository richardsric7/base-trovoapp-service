import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import Header from '../../../components/header';
import {
  useGetP2POfferQuery,
  useQuoteP2POrderFeesQuery,
  useCreateP2POrderMutation,
} from '../../../store/api/p2pApis';
import { useP2PIdentity } from '../../../hooks/useP2PIdentity';
import { P2PListCard, p2p } from '../../../components/p2p/P2PTheme';

// P2PCreateOrder (Plan Section 96.3): amount entry with a live, real (not
// estimated) itemized fee breakdown fetched from the backend's quote
// endpoint before the customer commits (Plan Section 97.6).
export default function P2PCreateOrder() {
  const { offerId } = useParams();
  const navigate = useNavigate();
  const { creds, ready } = useP2PIdentity();
  const [amount, setAmount] = useState('');
  const [debouncedAmount, setDebouncedAmount] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  const { data: offer } = useGetP2POfferQuery({ creds, offerId: offerId! }, { skip: !ready || !offerId });

  useEffect(() => {
    const t = setTimeout(() => setDebouncedAmount(amount), 500);
    return () => clearTimeout(t);
  }, [amount]);

  const { data: quote, isFetching: quoting, isError: quoteError } = useQuoteP2POrderFeesQuery(
    { creds, offerId: offerId!, amount: debouncedAmount },
    { skip: !ready || !offerId || !debouncedAmount },
  );

  const [createOrder] = useCreateP2POrderMutation();

  const submit = async () => {
    if (!offerId || !quote) return;
    setSubmitting(true);
    setError('');
    try {
      const order = await createOrder({ creds, offerId, specifiedAssetAmount: amount }).unwrap();
      navigate(`/dashboard/p2p/order/${order.id}`);
    } catch (e: any) {
      setError(e?.data?.message || 'Could not create order. Please try again.');
    } finally {
      setSubmitting(false);
    }
  };

  if (!offer) return <div className="p-8 text-center text-gray-400">Loading...</div>;

  const isBuy = offer.offerType === 'BUY';

  return (
    <div className="flex flex-col space-y-5 p-3 max-w-2xl">
      <Header />
      <h1 className="text-xl font-bold text-primary-800 px-2">Create order</h1>

      <div>
        <label className="font-semibold block mb-2">Amount ({offer.asset})</label>
        <input
          type="text"
          inputMode="decimal"
          value={amount}
          onChange={(e) => setAmount(e.target.value)}
          placeholder={`Between ${offer.minOrderAmount} and ${offer.maxOrderAmount}`}
          className="w-full p-3 rounded-2xl bg-white outline-none"
        />
      </div>

      {quoting && <div className="text-gray-400">Calculating fees...</div>}
      {quoteError && <div className="text-red-600">This amount is outside the offer's allowed range.</div>}
      {quote && (
        <P2PListCard>
          <p className="font-bold mb-3">Order breakdown</p>
          <Row label="You pay (fiat)" value={`${quote.paymentAmount} ${offer.currency}`} />
          <hr className="my-3" />
          {isBuy ? (
            <>
              <Row label="You receive" value={`${quote.buyerNetAssetAmount} ${offer.asset}`} />
              <Row label="Platform + regulatory fee" value={`${quote.buyerTotalFees} ${offer.asset}`} />
              <Row label="VAT on fee" value={`${quote.buyerTotalVat} ${offer.asset}`} />
            </>
          ) : (
            <>
              <Row label="Escrow deposit required" value={`${quote.sellerEscrowAssetAmount} ${offer.asset}`} />
              <Row label="Platform + regulatory fee" value={`${quote.sellerTotalFees} ${offer.asset}`} />
              <Row label="VAT on fee" value={`${quote.sellerTotalVat} ${offer.asset}`} />
            </>
          )}
        </P2PListCard>
      )}
      {error && <div className="text-red-600">{error}</div>}

      <button
        disabled={!quote || submitting}
        onClick={submit}
        className="w-full py-4 rounded-2xl font-semibold text-white disabled:opacity-50"
        style={{ backgroundColor: p2p.brandDark }}
      >
        Review and submit order
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
