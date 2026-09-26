import { useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Header from '../../../components/header';
import {
  useListP2PMarketplaceOffersQuery,
  useGetMerchantP2PPerformanceQuery,
  useGetP2PAssetClassesQuery,
  useGetP2PMarketplaceFacetsQuery,
  P2PCreds,
} from '../../../store/api/p2pApis';
import { useP2PIdentity } from '../../../hooks/useP2PIdentity';
import { P2PListCard, P2PEmptyState, p2p } from '../../../components/p2p/P2PTheme';

// P2PMarketplace (Plan Section 96.1): a Buy/Sell toggle + filterable offer
// list, mobile-first regardless of the rest of app-web's desktop-primary
// posture (Plan Section 92's deliberate improvement), built at the history
// page's level of care, not the wallet page's.
export default function P2PMarketplace() {
  const navigate = useNavigate();
  const { creds, ready } = useP2PIdentity();
  const [offerType, setOfferType] = useState<'BUY' | 'SELL'>('BUY');
  const [asset, setAsset] = useState('');
  const [currency, setCurrency] = useState('');
  const [assetClassId, setAssetClassId] = useState('');

  const { data: assetClasses } = useGetP2PAssetClassesQuery({ creds }, { skip: !ready });
  const { data: facets } = useGetP2PMarketplaceFacetsQuery({ creds }, { skip: !ready });

  // Asset is a dependent list box for category: picking a category
  // narrows the asset options down to just that category's assets.
  const assetOptions = useMemo(() => {
    if (!assetClassId) return facets?.assets ?? [];
    return (facets?.assets ?? []).filter((a) => String(a.assetClassId) === assetClassId);
  }, [facets, assetClassId]);

  const selectCategory = (value: string) => {
    setAssetClassId(value);
    // Currently selected asset may no longer belong to the new category -
    // reset it back to "All assets" rather than leaving a stale, now-
    // invalid selection in place.
    const stillValid = !value || (facets?.assets ?? []).some((a) => a.asset === asset && String(a.assetClassId) === value);
    if (!stillValid) setAsset('');
  };

  const { data, isLoading, isError, refetch } = useListP2PMarketplaceOffersQuery(
    {
      creds,
      offerType,
      asset: asset || undefined,
      currency: currency || undefined,
      assetClassId: assetClassId ? Number(assetClassId) : undefined,
      page: 1,
      pageSize: 20,
    },
    { skip: !ready },
  );

  return (
    <div className="flex flex-col space-y-5 p-3">
      <Header />
      <div className="flex items-center justify-between px-2">
        <h1 className="text-xl font-bold text-primary-800">P2P Marketplace</h1>
        <div className="flex gap-3">
          <button
            className="text-sm font-semibold text-primary-800 underline"
            onClick={() => navigate('/dashboard/p2p/my-orders')}
          >
            My orders
          </button>
          <button
            className="text-sm font-semibold text-primary-800 underline"
            onClick={() => navigate('/dashboard/p2p/my-offers')}
          >
            My offers
          </button>
          <button
            className="text-sm font-semibold text-primary-800 underline"
            onClick={() => navigate('/dashboard/p2p/my-refunds')}
          >
            My refunds
          </button>
        </div>
      </div>

      <div className="flex gap-2 max-w-md">
        {(['BUY', 'SELL'] as const).map((t) => (
          <button
            key={t}
            onClick={() => setOfferType(t)}
            className="flex-1 py-3 rounded-2xl font-semibold"
            style={{
              backgroundColor: offerType === t ? p2p.brandDark : 'white',
              color: offerType === t ? 'white' : '#374151',
            }}
          >
            {t === 'BUY' ? 'Buy' : 'Sell'}
          </button>
        ))}
      </div>

      <div className="flex flex-wrap gap-2 max-w-2xl">
        <select
          value={assetClassId}
          onChange={(e) => selectCategory(e.target.value)}
          className="px-3 py-2 rounded-xl bg-white text-sm border border-gray-200"
        >
          <option value="">All categories</option>
          {assetClasses?.data?.map((c) => (
            <option key={c.id} value={c.id}>
              {c.assetClass}
            </option>
          ))}
        </select>
        <select
          value={asset}
          onChange={(e) => setAsset(e.target.value)}
          className="px-3 py-2 rounded-xl bg-white text-sm border border-gray-200"
        >
          <option value="">All assets</option>
          {assetOptions.map((a) => (
            <option key={a.asset} value={a.asset}>
              {a.asset}
            </option>
          ))}
        </select>
        <select
          value={currency}
          onChange={(e) => setCurrency(e.target.value)}
          className="px-3 py-2 rounded-xl bg-white text-sm border border-gray-200"
        >
          <option value="">All currencies</option>
          {facets?.currencies?.map((c) => (
            <option key={c} value={c}>
              {c}
            </option>
          ))}
        </select>
      </div>

      <div className="max-w-2xl w-full">
        {isLoading ? (
          <div className="text-center py-16 text-gray-400">Loading offers...</div>
        ) : isError ? (
          <P2PEmptyState message="Could not load offers. Please try again." ctaLabel="Retry" onCta={refetch} />
        ) : !data?.data?.length ? (
          <P2PEmptyState message="No offers match your filters right now." ctaLabel="Refresh" onCta={refetch} />
        ) : (
          data.data.map((offer) => (
            <P2PListCard key={offer.id} onClick={() => navigate(`/dashboard/p2p/offer/${offer.id}`)}>
              <div className="flex items-center gap-3">
                <div
                  className="w-10 h-10 rounded-full flex items-center justify-center text-white font-bold"
                  style={{ backgroundColor: p2p.brandDark }}
                >
                  {offer.merchantUsername?.substring(0, 1).toUpperCase()}
                </div>
                <div className="flex-1">
                  <div className="flex items-center justify-between">
                    <span className="font-semibold">{offer.merchantUsername}</span>
                    <span
                      className="w-2 h-2 rounded-full"
                      style={{ backgroundColor: offer.availabilityStatus === 'ONLINE' ? p2p.success : '#9CA3AF' }}
                    />
                  </div>
                  <div className="text-sm text-gray-500">
                    {offer.price} {offer.currency} / {offer.asset}
                  </div>
                  <div className="text-xs text-gray-400">
                    Limit: {offer.minOrderAmount} - {offer.maxOrderAmount} {offer.asset}
                  </div>
                </div>
                {offer.paymentMethod?.paymentChannel && (
                  <span className="text-xs px-2 py-1 rounded-lg bg-primary-100">
                    {offer.paymentMethod.paymentChannel}
                  </span>
                )}
              </div>
              <MerchantPerformanceBadge creds={creds} merchantId={offer.merchantUserId} ready={ready} />
            </P2PListCard>
          ))
        )}
      </div>
    </div>
  );
}

// MerchantPerformanceBadge shows the marketplace trust signal (Plan
// Section 8/26) on each offer card - RTK Query dedupes this per
// merchantId, so browsing a page of offers from the same handful of
// merchants only fetches each merchant's performance once.
function MerchantPerformanceBadge({ creds, merchantId, ready }: { creds: P2PCreds; merchantId: string; ready: boolean }) {
  const { data: perf } = useGetMerchantP2PPerformanceQuery({ creds, merchantId }, { skip: !ready || !merchantId });
  if (!perf || perf.completedTrades === 0) return null;
  return (
    <p className="text-xs text-gray-400 mt-1">
      {perf.completedTrades} trade{perf.completedTrades === 1 ? '' : 's'} · {perf.completionRate}% completion
    </p>
  );
}
