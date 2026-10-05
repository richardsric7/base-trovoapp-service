import { useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import Header from '../../../components/header';
import { P2PEmptyState, p2p } from '../../../components/p2p/P2PTheme';
import { useP2PIdentity } from '../../../hooks/useP2PIdentity';
import { useGetPMAssetQuery, useGetPMPortfolioQuery, useGetPMPricesQuery } from '../../../store/api/publicMarketsApis';
import { PMAsset, PMPriceRange } from '../../../types/publicMarkets';
import { AssetLogo, Card, changeColor, Chips, Loading, money, num, pct, PriceChart, qty, Row, StatusBadge, when } from './ui';

const RANGES: { key: PMPriceRange; label: string }[] = ['1D', '1W', '1M', '3M', '1Y', 'All'].map((r) => ({ key: r as PMPriceRange, label: r }));

// PublicMarketsAsset shows one asset: price and chart, the user's position,
// key figures, how the tokens are backed (custody chain), the trading
// session and corporate actions, with Buy and Sell.
export default function PublicMarketsAsset() {
  const { code = '' } = useParams();
  const navigate = useNavigate();
  const [range, setRange] = useState<PMPriceRange>('1D');
  const { creds, ready } = useP2PIdentity();
  const { data: asset, isLoading, isError, refetch } = useGetPMAssetQuery(code, { pollingInterval: 30000 });
  const { data: prices } = useGetPMPricesQuery({ code, range });
  const { data: portfolio } = useGetPMPortfolioQuery({ creds }, { skip: !ready });
  const holding = portfolio?.holdings.find((h) => h.asset.assetCode === code);

  if (isLoading) return <Page><Loading /></Page>;
  if (isError || !asset) return <Page><P2PEmptyState message="Could not load this asset." ctaLabel="Retry" onCta={refetch} /></Page>;

  const trade = (side: 'buy' | 'sell') => navigate(`/dashboard/public-markets/asset/${encodeURIComponent(code)}/trade`, { state: { side } });
  const bond = asset.assetType === 'BOND';

  return (
    <Page>
      <div className="flex items-center gap-3 px-2">
        <AssetLogo asset={asset} size={52} />
        <div className="flex-1">
          <p className="text-xl font-bold text-primary-800">{asset.ticker}</p>
          <p className="text-gray-500 text-sm">{asset.name}</p>
        </div>
        <StatusBadge asset={asset} />
      </div>

      <Card>
        <div className="flex items-end justify-between flex-wrap gap-2">
          <div>
            <p className="text-3xl font-bold">{money(asset.price, asset.currency)}</p>
            <p className="font-semibold" style={{ color: changeColor(asset.dayChangePercent) }}>
              {pct(asset.dayChangePercent)} today
            </p>
          </div>
          <p className="text-xs text-gray-400">
            {asset.priceLive ? 'Live' : 'Last'} price · {asset.priceSource || '—'} · {when(asset.priceAt)}
          </p>
        </div>
        <div className="my-4">
          <PriceChart points={prices?.points ?? []} />
        </div>
        <Chips items={RANGES} value={range} onChange={setRange} />
      </Card>

      {holding && num(holding.quantity) > 0 && (
        <Card title="Your position">
          <Row label="Units" value={`${qty(holding.quantity)} ${asset.ticker}`} />
          <Row label="Market value" value={money(holding.marketValue)} />
          <Row label="Average cost" value={money(holding.averageCost)} />
          <Row
            label="Total return"
            value={<span style={{ color: changeColor(holding.totalReturn) }}>{money(holding.totalReturn)} ({pct(holding.returnPercent)})</span>}
          />
          {num(holding.incomeReceived) > 0 && <Row label={bond ? 'Coupons received' : 'Dividends received'} value={money(holding.incomeReceived)} />}
        </Card>
      )}

      <div className="flex gap-3 px-2">
        <button
          disabled={asset.status !== 'open'}
          onClick={() => trade('buy')}
          className="flex-1 py-3 rounded-xl font-bold text-white disabled:opacity-40"
          style={{ backgroundColor: p2p.success }}
        >
          Buy
        </button>
        <button
          disabled={asset.status !== 'open' || !holding || num(holding.quantity) <= 0}
          onClick={() => trade('sell')}
          className="flex-1 py-3 rounded-xl font-bold text-white disabled:opacity-40"
          style={{ backgroundColor: p2p.danger }}
        >
          Sell
        </button>
      </div>
      {asset.status === 'halted' && <p className="px-2 text-sm text-gray-500">Trading in this asset is paused. Your tokens stay in your wallet.</p>}

      <Session asset={asset} />

      <Card title="Key figures">
        <Row label="Previous close" value={money(asset.previousClose, asset.currency)} />
        <Row label="Day range" value={`${money(asset.dayLow, asset.currency)} – ${money(asset.dayHigh, asset.currency)}`} />
        <Row label="Volume" value={qty(asset.dayVolume)} />
        {bond ? (
          <>
            <Row label="Coupon" value={asset.coupon ? `${asset.coupon}%` : '—'} />
            <Row label="Maturity" value={asset.maturityDate ? new Date(asset.maturityDate).toLocaleDateString() : '—'} />
          </>
        ) : (
          <>
            <Row label="Market cap" value={money(asset.marketCap, asset.currency)} />
            <Row label="P/E ratio" value={asset.peRatio || '—'} />
            <Row label="Dividend yield" value={asset.dividendYield ? `${asset.dividendYield}%` : '—'} />
          </>
        )}
        <Row label="Sector" value={asset.sector || '—'} />
        <Row label="ISIN" value={asset.isin || '—'} />
        <Row label="Trovo fee" value={`${asset.feePercent}%`} />
        <Row label="Minimum buy" value={money(asset.minimumBuy, asset.fundingAsset)} />
      </Card>

      {asset.description && (
        <Card title="About">
          <p className="text-gray-600 whitespace-pre-line">{asset.description}</p>
          {asset.unitDescription && <p className="text-sm text-gray-500 mt-2">{asset.unitDescription}</p>}
        </Card>
      )}

      {asset.custody && (
        <Card title="How your tokens are backed">
          <Row label="Custodian" value={asset.custody.custodian || '—'} />
          <Row label="Nominee" value={asset.custody.nominee || '—'} />
          <Row label="Dealing member" value={asset.custody.dealingMember || '—'} />
          <Row label="Depository" value={asset.custody.depository || '—'} />
          <Row label="Settlement" value={asset.custody.settlement || '—'} />
          <Row label="Units held at the CSCS" value={qty(asset.custody.unitsHeld)} />
          <Row label="Tokens in circulation" value={qty(asset.tokensInCirculation)} />
          <Row label="Position as of" value={when(asset.custody.positionAsOf)} />
          <Row
            label="Last reconciliation"
            value={asset.custody.lastReconciliation?.at ? `${asset.custody.lastReconciliation.result} · ${when(asset.custody.lastReconciliation.at)}` : '—'}
          />
          <Row label="Token contract" value={<span className="font-mono text-xs">{asset.contractAddress || '—'}</span>} />
        </Card>
      )}

      {(asset.corporateActions ?? []).length > 0 && (
        <Card title="Corporate actions">
          {asset.corporateActions!.map((ca) => (
            <div key={ca.id} className="py-2 border-b last:border-0 border-gray-100">
              <div className="flex justify-between">
                <span className="font-semibold">{ca.eventType}</span>
                <span className="text-sm text-gray-500">{ca.status}</span>
              </div>
              <p className="text-sm text-gray-600">{ca.description}</p>
              <p className="text-xs text-gray-400">
                Record date {ca.recordDate}
                {ca.payDate && ` · Pay date ${ca.payDate}`}
                {ca.amountPerUnit && ` · ${money(ca.amountPerUnit, ca.currency || 'NGN')} per unit`}
              </p>
            </div>
          ))}
          <button className="text-primary-800 font-semibold underline mt-2" onClick={() => navigate('/dashboard/public-markets/dividends', { state: { assetCode: asset.assetCode } })}>
            Your payments
          </button>
        </Card>
      )}
    </Page>
  );
}

function Session({ asset }: { asset: PMAsset }) {
  const s = asset.session;
  if (!s) return null;
  const elapsed = s.open && s.minutesTotal ? Math.max(0, Math.min(1, 1 - (s.minutesRemaining ?? 0) / s.minutesTotal)) : 0;
  return (
    <Card title={`${s.market || asset.market} trading session`}>
      {s.open ? (
        <>
          <p className="text-sm text-gray-600 mb-2">
            Open · closes {when(s.closesAt)} ({s.minutesRemaining} min left). Buys are usually filled straight away from inventory.
          </p>
          <div className="h-2 rounded-full bg-gray-100 overflow-hidden">
            <div className="h-full" style={{ width: `${elapsed * 100}%`, backgroundColor: p2p.success }} />
          </div>
        </>
      ) : (
        <p className="text-sm text-gray-600">
          Closed{s.opensAt ? ` · opens ${when(s.opensAt)}` : ''}. Orders placed now are queued and executed in the next session.
        </p>
      )}
    </Card>
  );
}

function Page({ children }: { children: React.ReactNode }) {
  const navigate = useNavigate();
  return (
    <div className="flex flex-col space-y-5 p-3 max-w-3xl">
      <Header />
      <button className="self-start text-primary-800 font-semibold px-2" onClick={() => navigate('/dashboard/public-markets')}>
        ‹ Public Markets
      </button>
      {children}
    </div>
  );
}
