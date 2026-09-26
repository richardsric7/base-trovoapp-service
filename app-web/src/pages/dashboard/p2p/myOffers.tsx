import { useNavigate } from 'react-router-dom';
import Header from '../../../components/header';
import {
  useListMyP2POffersQuery,
  useActivateP2POfferMutation,
  usePauseP2POfferMutation,
  useCloseP2POfferMutation,
} from '../../../store/api/p2pApis';
import { useP2PIdentity } from '../../../hooks/useP2PIdentity';
import { P2PListCard, P2PEmptyState, p2p } from '../../../components/p2p/P2PTheme';

// P2PMyOffers (Plan Section 96.11): a merchant's own offers, with
// online/offline toggle.
export default function P2PMyOffers() {
  const navigate = useNavigate();
  const { creds, ready } = useP2PIdentity();
  const { data, isLoading, refetch } = useListMyP2POffersQuery({ creds }, { skip: !ready });
  const [activate] = useActivateP2POfferMutation();
  const [pause] = usePauseP2POfferMutation();
  const [closeOffer] = useCloseP2POfferMutation();

  const toggle = async (offerId: string, isOnline: boolean) => {
    if (isOnline) await pause({ creds, offerId });
    else await activate({ creds, offerId });
    refetch();
  };

  const close = async (offerId: string) => {
    if (!window.confirm('Close this offer? This cannot be undone.')) return;
    await closeOffer({ creds, offerId });
    refetch();
  };

  return (
    <div className="flex flex-col space-y-5 p-3 max-w-2xl">
      <Header />
      <div className="flex items-center justify-between px-2">
        <h1 className="text-xl font-bold text-primary-800">My offers</h1>
        <button
          onClick={() => navigate('/dashboard/p2p/create-offer')}
          className="text-sm font-semibold text-white px-4 py-2 rounded-xl"
          style={{ backgroundColor: p2p.brandDark }}
        >
          + New offer
        </button>
      </div>

      {isLoading ? (
        <div className="text-center py-16 text-gray-400">Loading...</div>
      ) : !data?.data?.length ? (
        <P2PEmptyState
          message="You have no offers yet."
          ctaLabel="Create your first offer"
          onCta={() => navigate('/dashboard/p2p/create-offer')}
        />
      ) : (
        data.data.map((offer) => (
          <P2PListCard key={offer.id}>
            <div className="flex items-center justify-between">
              <div>
                <p className="font-bold">
                  {offer.offerType} {offer.asset}
                </p>
                <p className="text-gray-500">
                  {offer.price} {offer.currency}
                </p>
                <p className="text-xs text-gray-400">{offer.status}</p>
              </div>
              {offer.status !== 'CLOSED' && (
                <label className="inline-flex items-center cursor-pointer">
                  <input
                    type="checkbox"
                    className="sr-only"
                    checked={offer.availabilityStatus === 'ONLINE'}
                    onChange={() => toggle(offer.id, offer.availabilityStatus === 'ONLINE')}
                  />
                  <span
                    className="w-11 h-6 rounded-full relative transition-colors"
                    style={{ backgroundColor: offer.availabilityStatus === 'ONLINE' ? p2p.success : '#D1D5DB' }}
                  >
                    <span
                      className="absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full transition-transform"
                      style={{ transform: offer.availabilityStatus === 'ONLINE' ? 'translateX(20px)' : 'none' }}
                    />
                  </span>
                </label>
              )}
            </div>
            {offer.status !== 'CLOSED' && (
              <div className="flex gap-3 mt-3">
                <button
                  onClick={() => navigate(`/dashboard/p2p/offer/${offer.id}/edit`)}
                  className="text-sm font-semibold underline"
                  style={{ color: p2p.brandDark }}
                >
                  Edit
                </button>
                <button onClick={() => close(offer.id)} className="text-sm font-semibold underline" style={{ color: p2p.danger }}>
                  Close
                </button>
              </div>
            )}
          </P2PListCard>
        ))
      )}
    </div>
  );
}
