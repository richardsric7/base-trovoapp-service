import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useP2PIdentity } from '../../hooks/useP2PIdentity';
import { useGetP2PMerchantStatusQuery, useRequestP2PMerchantStatusMutation } from '../../store/api/p2pApis';
import { p2p } from './P2PTheme';

// P2PMerchantGate wraps a merchant-only page (create/edit offer, my
// offers): a caller who isn't a merchant yet sees a notice explaining
// that only merchants can do this, with a button to request merchant
// status right there, instead of the page's own content. The backend
// enforces the same rule itself (CreateOffer rejects a non-merchant) -
// this is the UX half of that rule, not a substitute for it.
export default function P2PMerchantGate({ children }: { children: React.ReactNode }) {
  const navigate = useNavigate();
  const { creds, ready } = useP2PIdentity();
  const { data: status, isLoading, isFetching } = useGetP2PMerchantStatusQuery({ creds }, { skip: !ready });
  const [requestMerchantStatus, { isLoading: requesting }] = useRequestP2PMerchantStatusMutation();
  const [error, setError] = useState('');

  if (!ready || isLoading) {
    return <div className="text-center py-16 text-gray-400">Loading...</div>;
  }

  if (status?.isMerchant) {
    return <>{children}</>;
  }

  const requestStatus = async () => {
    setError('');
    try {
      await requestMerchantStatus({ creds }).unwrap();
    } catch (e: any) {
      setError(e?.data?.message || e?.data?.error || 'Could not request merchant status. Please try again.');
    }
  };

  return (
    <div className="flex flex-col items-center justify-center py-16 px-6 text-center max-w-md mx-auto">
      <div className="text-4xl mb-4" style={{ color: p2p.warning }}>
        🛈
      </div>
      <p className="text-lg font-bold text-primary-800 mb-2">Only merchants can do this</p>
      <p className="text-gray-500 mb-6">
        You need to become a P2P merchant before you can create and manage offers. It only takes a moment - if
        you've completed KYC level 2, you'll be approved right away.
      </p>
      {error && <p className="text-sm mb-4" style={{ color: p2p.danger }}>{error}</p>}
      <button
        disabled={requesting || isFetching}
        onClick={requestStatus}
        className="w-full py-4 rounded-2xl font-semibold text-white disabled:opacity-50"
        style={{ backgroundColor: p2p.brandDark }}
      >
        {requesting ? 'Requesting...' : 'Become a merchant'}
      </button>
      <button onClick={() => navigate('/dashboard/p2p/marketplace')} className="mt-4 text-sm font-semibold underline text-gray-500">
        Back to marketplace
      </button>
    </div>
  );
}
