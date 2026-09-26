import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Header from '../../../components/header';
import { useCreateP2POfferMutation } from '../../../store/api/p2pApis';
import { useP2PIdentity } from '../../../hooks/useP2PIdentity';
import { p2p } from '../../../components/p2p/P2PTheme';
import P2PMerchantGate from '../../../components/p2p/P2PMerchantGate';

// P2PCreateOffer (Plan Section 96.12): merchant offer creation - per Plan
// Section 92, whether this introduces react-hook-form/zod or follows the
// existing manual-useState pattern is a genuine open implementation
// choice; this follows the manual pattern for consistency with the rest
// of app-web today.
export default function P2PCreateOffer() {
  const navigate = useNavigate();
  const { creds, ready } = useP2PIdentity();
  const [createOffer] = useCreateP2POfferMutation();
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  const [offerType, setOfferType] = useState<'BUY' | 'SELL'>('SELL');
  const [asset, setAsset] = useState('');
  const [price, setPrice] = useState('');
  const [currency, setCurrency] = useState('');
  const [minOrderAmount, setMinOrderAmount] = useState('');
  const [maxOrderAmount, setMaxOrderAmount] = useState('');
  const [availableLiquidity, setAvailableLiquidity] = useState('');
  const [country, setCountry] = useState('');
  const [countryCode, setCountryCode] = useState('');
  const [paymentChannel, setPaymentChannel] = useState('');
  const [provider, setProvider] = useState('');
  const [account, setAccount] = useState('');
  const [remark, setRemark] = useState('');

  const submit = async () => {
    if (!ready) return;
    setSubmitting(true);
    setError('');
    try {
      await createOffer({
        creds,
        body: {
          offerType,
          asset: asset.toUpperCase(),
          paymentMethod: { paymentChannel, provider, account },
          country,
          countryCode: countryCode.toUpperCase(),
          currency: currency.toUpperCase(),
          priceType: 'FIXED',
          price,
          minOrderAmount,
          maxOrderAmount,
          availableLiquidity,
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
      <h1 className="text-xl font-bold text-primary-800 px-2">Create offer</h1>
      <P2PMerchantGate>
      <div className="flex gap-2 max-w-xs">
        {(['SELL', 'BUY'] as const).map((t) => (
          <button
            key={t}
            onClick={() => setOfferType(t)}
            className="flex-1 py-2 rounded-xl font-semibold"
            style={{ backgroundColor: offerType === t ? p2p.brandDark : 'white', color: offerType === t ? 'white' : '#374151' }}
          >
            {t === 'SELL' ? 'Sell' : 'Buy'}
          </button>
        ))}
      </div>

      <Field label="Asset code (e.g. USDC)" value={asset} onChange={setAsset} />
      <Field label="Price per unit" value={price} onChange={setPrice} numeric />
      <Field label="Currency (e.g. NGN)" value={currency} onChange={setCurrency} />
      <Field label="Min order amount" value={minOrderAmount} onChange={setMinOrderAmount} numeric />
      <Field label="Max order amount" value={maxOrderAmount} onChange={setMaxOrderAmount} numeric />
      <Field label="Available liquidity" value={availableLiquidity} onChange={setAvailableLiquidity} numeric />
      <Field label="Country" value={country} onChange={setCountry} />
      <Field label="Country code (ISO-3)" value={countryCode} onChange={setCountryCode} />
      <p className="font-semibold">Payment method</p>
      <Field label="Payment channel (e.g. BANK_TRANSFER)" value={paymentChannel} onChange={setPaymentChannel} />
      <Field label="Provider (e.g. GTBANK)" value={provider} onChange={setProvider} />
      <Field label="Account details" value={account} onChange={setAccount} />
      <Field label="Remark (optional)" value={remark} onChange={setRemark} />

      {error && <div className="text-red-600">{error}</div>}
      <button
        disabled={submitting}
        onClick={submit}
        className="w-full py-4 rounded-2xl font-semibold text-white disabled:opacity-50"
        style={{ backgroundColor: p2p.brandDark }}
      >
        Create offer
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
