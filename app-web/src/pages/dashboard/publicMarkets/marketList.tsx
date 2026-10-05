import { useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Header from '../../../components/header';
import { P2PEmptyState, P2PListCard } from '../../../components/p2p/P2PTheme';
import { useListPMAssetsQuery } from '../../../store/api/publicMarketsApis';
import { PMAsset } from '../../../types/publicMarkets';
import { AssetLogo, changeColor, Chips, Loading, money, num, pct, StatusBadge } from './ui';

type Filter = 'all' | 'EQUITY' | 'BOND' | 'gainers';

const FILTERS: { key: Filter; label: string }[] = [
  { key: 'all', label: 'All' },
  { key: 'EQUITY', label: 'Equities' },
  { key: 'BOND', label: 'Bonds' },
  { key: 'gainers', label: 'Top gainers' },
];

// PublicMarketsList is the catalogue of tokenized NGX equities and FMDQ
// bonds: search, filter and open an asset to see its detail and trade.
export default function PublicMarketsList() {
  const navigate = useNavigate();
  const [search, setSearch] = useState('');
  const [filter, setFilter] = useState<Filter>('all');
  const { data, isLoading, isError, refetch } = useListPMAssetsQuery({ search: search.trim() || undefined }, { pollingInterval: 30000 });

  const assets = useMemo(() => {
    let list = data?.assets ?? [];
    if (filter === 'EQUITY' || filter === 'BOND') list = list.filter((a) => a.assetType === filter);
    if (filter === 'gainers') list = [...list].filter((a) => num(a.dayChangePercent) > 0).sort((a, b) => num(b.dayChangePercent) - num(a.dayChangePercent));
    return list;
  }, [data, filter]);
  const anyOpen = (data?.assets ?? []).some((a) => a.session?.open);

  return (
    <div className="flex flex-col space-y-5 p-3 max-w-3xl">
      <Header />
      <div className="flex flex-wrap items-center justify-between gap-3 px-2">
        <div>
          <h1 className="text-xl font-bold text-primary-800">Public Markets</h1>
          <p className="text-gray-500 text-sm">
            Own Nigerian stocks and bonds as tokens in your wallet, backed one-to-one by shares held for you at the CSCS.
          </p>
        </div>
        <div className="flex gap-2">
          <button className="px-4 py-2 rounded-xl font-semibold bg-white text-primary-800" onClick={() => navigate('/dashboard/public-markets/portfolio')}>
            My Stocks
          </button>
          <button className="px-4 py-2 rounded-xl font-semibold bg-white text-primary-800" onClick={() => navigate('/dashboard/public-markets/dividends')}>
            Dividends
          </button>
        </div>
      </div>
      <div className="px-2 space-y-3">
        <input
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="Search by name or ticker"
          className="w-full rounded-xl border border-gray-200 px-4 py-3 bg-white"
        />
        <Chips items={FILTERS} value={filter} onChange={setFilter} />
        {data && !anyOpen && (
          <p className="text-sm text-gray-500">
            The market is closed. Orders placed now are queued and filled when the next session opens.
          </p>
        )}
      </div>
      <div className="px-2">
        {isLoading ? (
          <Loading />
        ) : isError ? (
          <P2PEmptyState message="Could not load the market." ctaLabel="Retry" onCta={refetch} />
        ) : assets.length === 0 ? (
          <P2PEmptyState message={search ? 'No asset matches your search.' : 'No assets listed yet.'} />
        ) : (
          assets.map((a) => <AssetRow key={a.assetCode} asset={a} onClick={() => navigate(`/dashboard/public-markets/asset/${encodeURIComponent(a.assetCode)}`)} />)
        )}
      </div>
    </div>
  );
}

function AssetRow({ asset, onClick }: { asset: PMAsset; onClick: () => void }) {
  return (
    <P2PListCard onClick={onClick}>
      <div className="flex items-center gap-3">
        <AssetLogo asset={asset} />
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <p className="font-bold text-primary-800">{asset.ticker}</p>
            <span className="text-xs text-gray-400">{asset.market}</span>
          </div>
          <p className="text-sm text-gray-500 truncate">{asset.shortName || asset.name}</p>
        </div>
        <div className="text-right">
          {asset.status === 'open' ? (
            <>
              <p className="font-bold">{money(asset.price, asset.currency)}</p>
              <p className="text-sm font-semibold" style={{ color: changeColor(asset.dayChangePercent) }}>
                {asset.assetType === 'BOND' && asset.coupon ? `${asset.coupon}% coupon` : pct(asset.dayChangePercent)}
              </p>
            </>
          ) : (
            <StatusBadge asset={asset} />
          )}
        </div>
      </div>
    </P2PListCard>
  );
}
