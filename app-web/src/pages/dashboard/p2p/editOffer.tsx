import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import Header from '../../../components/header';
import { useGetP2POfferQuery, useUpdateP2POfferMutation } from '../../../store/api/p2pApis';
import { useP2PIdentity } from '../../../hooks/useP2PIdentity';
import { p2p } from '../../../components/p2p/P2PTheme';
import P2PMerchantGate from '../../../components/p2p/P2PMerchantGate';
import PaymentMethodSelect from '../../../components/p2p/PaymentMethodSelect';
import WalletPayoutSelect from '../../../components/p2p/WalletPayoutSelect';

// P2PEditOffer lets a merchant edit an existing offer's terms (Plan Section
// 13's implied listing-management control). offerType/asset are fixed at
// creation (they drive curated-asset validation and the role mapping in
// Plan Section 11) so they're shown read-only here, not editable fields.
export default function P2PEditOffer() {
  const { offerId } = useParams();
  const navigate = useNavigate();
  const { creds, ready } = useP2PIdentity();
  const { data: offer, isLoading } = useGetP2POfferQuery({ creds, offerId: offerId! }, { skip: !ready || !offerId });
  const [updateOffer] = useUpdateP2POfferMutation();
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  const [price, setPrice] = useState('');
  const [minOrderAmount, setMinOrderAmount] = useState('');
  const [maxOrderAmount, setMaxOrderAmount] = useState('');
  const [addLiquidity, setAddLiquidity] = useState('');
  const [paymentMethodId, setPaymentMethodId] = useState('');
  const [merchantPayoutAddress, setMerchantPayoutAddress] = useState('');
  const [remark, setRemark] = useState('');

  useEffect(() => {
    if (!offer) return;
    setPrice(offer.price);
    setMinOrderAmount(offer.minOrderAmount);
    setMaxOrderAmount(offer.maxOrderAmount);
    setPaymentMethodId(offer.paymentMethodId ?? '');
    setMerchantPayoutAddress(offer.merchantPayoutAddress ?? '');
    setRemark(offer.remark ?? '');
  }, [offer]);

  if (isLoading || !offer) return <div className="p-8 text-center text-gray-400">Loading...</div>;

  const submit = async () => {
    setSubmitting(true);
    setError('');
    try {
      await updateOffer({
        creds,
        offerId: offer.id,
        body: {
          offerType: offer.offerType,
          asset: offer.asset,
          paymentMethodId: offer.offerType === 'SELL' ? paymentMethodId : undefined,
          merchantPayoutAddress: offer.offerType === 'BUY' ? merchantPayoutAddress : undefined,
          country: offer.country,
          countryCode: offer.countryCode,
          currency: offer.currency,
          priceType: offer.priceType,
          price,
          minOrderAmount,
          maxOrderAmount,
          availableLiquidity: addLiquidity,
          remark,
        },
      }).unwrap();
      navigate('/dashboard/p2p/my-offers');
    } catch (e: any) {
      setError(e?.data?.message || 'Please check your inputs and try again.');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="flex flex-col space-y-4 p-3 max-w-2xl">
      <Header />
      <h1 className="text-xl font-bold text-primary-800 px-2">Edit offer</h1>
      <p className="text-sm text-gray-400 px-2">
        {offer.offerType} {offer.asset} ({offer.currency})
      </p>

      <P2PMerchantGate>
      <Field label="Price per unit" value={price} onChange={setPrice} numeric />
      <Field label="Min order amount" value={minOrderAmount} onChange={setMinOrderAmount} numeric />
      <Field label="Max order amount" value={maxOrderAmount} onChange={setMaxOrderAmount} numeric />
      <Field
        label={`Add liquidity (current: ${offer.availableLiquidity} ${offer.asset})`}
        value={addLiquidity}
        onChange={setAddLiquidity}
        numeric
      />
      {offer.offerType === 'SELL' ? (
        <PaymentMethodSelect creds={creds} ready={ready} value={paymentMethodId} onChange={setPaymentMethodId} />
      ) : (
        <WalletPayoutSelect value={merchantPayoutAddress} onChange={setMerchantPayoutAddress} />
      )}
      <Field label="Remark (optional)" value={remark} onChange={setRemark} />

      {error && <div className="text-red-600">{error}</div>}
      <button
        disabled={submitting}
        onClick={submit}
        className="w-full py-4 rounded-2xl font-semibold text-white disabled:opacity-50"
        style={{ backgroundColor: p2p.brandDark }}
      >
        Save changes
      </button>
      </P2PMerchantGate>
    </div>
  );
}

function Field({
  label,
  value,
  onChange,
  numeric,
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
  numeric?: boolean;
}) {
  return (
    <div>
      <label className="text-sm text-gray-500 block mb-1">{label}</label>
      <input
        type="text"
        inputMode={numeric ? 'decimal' : 'text'}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="w-full p-3 rounded-xl bg-white outline-none"
      />
    </div>
  );
}
