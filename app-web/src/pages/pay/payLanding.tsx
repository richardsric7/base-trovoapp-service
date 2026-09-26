import { useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { RootState } from '../../store/reduxStore';
import { useSendAssetMutation } from '../../store/api/walletApis';
import { useP2PIdentity } from '../../hooks/useP2PIdentity';
import { signBase64Txn } from '../../utils/trovoSDK';
import { p2p } from '../../components/p2p/P2PTheme';

// PayLanding is a new page with no existing app-web precedent (Plan
// Section 92/96.6's flagged gap): app-web today only ever *produces*
// action=payment shortlinks (via the receive-QR flow), it never consumes
// one. When a resolved escrow-deposit shortlink is opened on a desktop
// browser or by a non-mobile client (Plan Section 33's third-party payer,
// who may not even be a Trovo user yet), it lands here instead of in a
// native app's deep-link handler. Deliberately outside ProtectedRoutes,
// like /login - a visitor with no Trovo account yet still needs to see
// what they're being asked to pay before being asked to sign up.
export default function PayLanding() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const appUser = useSelector((state: RootState) => state.auth.user);
  const { creds, ready } = useP2PIdentity();
  const [sendAsset] = useSendAssetMutation();
  const [step, setStep] = useState<'review' | 'signing' | 'done' | 'error'>('review');
  const [error, setError] = useState('');

  const destination = searchParams.get('paymentDestination') || '';
  const assetCode = searchParams.get('assetCode') || '';
  const contractAddress = searchParams.get('contractAddress') || '';
  const amount = searchParams.get('amount') || '';
  const memo = searchParams.get('memo') || '';
  const action = searchParams.get('action');

  const pay = async () => {
    if (!ready) return;
    setStep('signing');
    setError('');
    try {
      const built = await sendAsset({
        signer: creds.signer,
        address: creds.address,
        secretKey: creds.secretKey,
        body: { destination, memo, amount, assetCode, contractAddress },
      }).unwrap();

      if (!built.transaction) {
        setError('Could not build this transaction.');
        setStep('error');
        return;
      }
      const transactionSignature = signBase64Txn(creds.secretKey, built.transaction, built.networkPassPhrase);
      await sendAsset({
        signer: creds.signer,
        address: creds.address,
        secretKey: creds.secretKey,
        body: { ...built, commit: 1, transactionSignature },
      }).unwrap();
      setStep('done');
    } catch (e: any) {
      setError(e?.data?.message || 'Payment failed. Please try again.');
      setStep('error');
    }
  };

  if (action !== 'payment' || !destination) {
    return (
      <div className="min-h-screen flex items-center justify-center p-8">
        <p className="text-gray-500">This link is not a valid payment link.</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex flex-col items-center justify-center p-8" style={{ backgroundColor: p2p.brandDark }}>
      <div className="bg-white rounded-3xl p-8 max-w-md w-full">
        <p className="text-gray-500 mb-1">You're being asked to pay</p>
        <p className="text-3xl font-bold mb-6" style={{ color: p2p.brandDark }}>
          {amount} {assetCode}
        </p>
        <Row label="To" value={destination} mono />
        {memo && <Row label="Reference" value={memo} mono />}

        {!appUser ? (
          <>
            <p className="text-gray-500 mt-6 mb-3">
              Log in to your Trovo account to complete this payment, then come back to this link.
            </p>
            <button
              onClick={() => navigate('/login')}
              className="w-full py-4 rounded-2xl font-semibold text-white mt-2"
              style={{ backgroundColor: p2p.brandDark }}
            >
              Log in to pay
            </button>
          </>
        ) : step === 'done' ? (
          <p className="font-semibold mt-6" style={{ color: p2p.success }}>
            Payment submitted successfully.
          </p>
        ) : (
          <>
            {error && <p className="text-red-600 mt-4">{error}</p>}
            <button
              disabled={!ready || step === 'signing'}
              onClick={pay}
              className="w-full py-4 rounded-2xl font-semibold text-white mt-6 disabled:opacity-50"
              style={{ backgroundColor: p2p.brandDark }}
            >
              {step === 'signing' ? 'Submitting...' : 'Approve and pay'}
            </button>
          </>
        )}
      </div>
    </div>
  );
}

function Row({ label, value, mono }: { label: string; value: string; mono?: boolean }) {
  return (
    <div className="flex justify-between py-2 border-t border-gray-100">
      <span className="text-gray-500">{label}</span>
      <span className={`font-semibold text-right ${mono ? 'font-mono text-sm' : ''}`}>{value}</span>
    </div>
  );
}
