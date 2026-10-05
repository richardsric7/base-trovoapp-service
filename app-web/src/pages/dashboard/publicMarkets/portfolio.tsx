import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Header from '../../../components/header';
import { P2PEmptyState, P2PListCard } from '../../../components/p2p/P2PTheme';
import { useP2PIdentity } from '../../../hooks/useP2PIdentity';
import { useGetPMOrdersQuery, useGetPMPortfolioQuery } from '../../../store/api/publicMarketsApis';
import { PMOrder } from '../../../types/publicMarkets';
import { AssetLogo, Card, changeColor, Chips, Loading, money, num, pct, Pill, qty, stateColor, stateLabel, when } from './ui';

type Tab = 'holdings' | 'orders' | 'activity';

// PublicMarketsPortfolio ("My Stocks") is the user's Public Markets
// holdings across all their wallets: value, returns and income, each
// holding, orders and activity.
export default function PublicMarketsPortfolio() {
  const navigate = useNavigate();
  const { creds, ready } = useP2PIdentity();
  const [tab, setTab] = useState<Tab>('holdings');
  const { data: p, isLoading, isError, refetch } = useGetPMPortfolioQuery({ creds }, { skip: !ready, pollingInterval: 30000 });
  const { data: orders } = useGetPMOrdersQuery({ creds }, { skip: !ready || tab !== 'orders' });

  let body: React.ReactNode;
  if (!ready || isLoading) body = <Loading />;
  else if (isError || !p) body = <P2PEmptyState message="Could not load your portfolio." ctaLabel="Retry" onCta={refetch} />;
  else
    body = (
      <>
        <Card>
          <p className="text-sm text-gray-500">Portfolio value</p>
          <p className="text-3xl font-bold">{money(p.value)}</p>
          <p className="font-semibold" style={{ color: changeColor(p.todayChange) }}>
            {money(p.todayChange)} today
          </p>
          <div className="grid grid-cols-3 gap-3 mt-4 text-sm">
            <div>
              <p className="text-gray-500">Invested</p>
              <p className="font-semibold">{money(p.costBasis)}</p>
            </div>
            <div>
              <p className="text-gray-500">Total return</p>
              <p className="font-semibold" style={{ color: changeColor(p.totalReturn) }}>
                {money(p.totalReturn)} ({pct(p.returnPercent)})
              </p>
            </div>
            <div>
              <p className="text-gray-500">Income</p>
              <p className="font-semibold">{money(p.income)}</p>
            </div>
          </div>
        </Card>
        {p.openOrders.length > 0 && (
          <Card title="In progress">
            {p.openOrders.map((o) => (
              <OrderRow key={o.id} order={o} onClick={() => navigate(`/dashboard/public-markets/order/${encodeURIComponent(o.id)}`)} />
            ))}
          </Card>
        )}
        <div className="px-2">
          <Chips
            items={[
              { key: 'holdings' as Tab, label: 'Holdings' },
              { key: 'orders' as Tab, label: 'Orders' },
              { key: 'activity' as Tab, label: 'Activity' },
            ]}
            value={tab}
            onChange={setTab}
          />
        </div>
        {tab === 'holdings' &&
          (p.holdings.length === 0 ? (
            <P2PEmptyState message="You don't own any stocks or bonds yet." ctaLabel="Browse the market" onCta={() => navigate('/dashboard/public-markets')} />
          ) : (
            <div>
              {p.holdings.map((h) => (
                <P2PListCard key={h.asset.assetCode} onClick={() => navigate(`/dashboard/public-markets/asset/${encodeURIComponent(h.asset.assetCode)}`)}>
                  <div className="flex items-center gap-3">
                    <AssetLogo asset={h.asset} />
                    <div className="flex-1 min-w-0">
                      <p className="font-bold text-primary-800">{h.asset.ticker}</p>
                      <p className="text-sm text-gray-500">
                        {qty(h.quantity)} units · avg {money(h.averageCost)}
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="font-bold">{money(h.marketValue)}</p>
                      <p className="text-sm font-semibold" style={{ color: changeColor(h.totalReturn) }}>
                        {pct(h.returnPercent)}
                      </p>
                    </div>
                  </div>
                  {num(h.incomeReceived) > 0 && <p className="text-xs text-gray-500 mt-2">Income received: {money(h.incomeReceived)}</p>}
                </P2PListCard>
              ))}
            </div>
          ))}
        {tab === 'orders' &&
          ((orders?.orders ?? []).length === 0 ? (
            <P2PEmptyState message="No orders yet." />
          ) : (
            <Card>
              {orders!.orders.map((o) => (
                <OrderRow key={o.id} order={o} onClick={() => navigate(`/dashboard/public-markets/order/${encodeURIComponent(o.id)}`)} />
              ))}
            </Card>
          ))}
        {tab === 'activity' &&
          (p.activity.length === 0 ? (
            <P2PEmptyState message="No activity yet." />
          ) : (
            <Card>
              {p.activity.map((a, i) => (
                <div key={`${a.reference}-${i}`} className="flex justify-between gap-3 py-2 border-b last:border-0 border-gray-100">
                  <div className="min-w-0">
                    <p className="font-semibold">{a.title}</p>
                    <p className="text-sm text-gray-500">{a.detail}</p>
                    <p className="text-xs text-gray-400">{when(a.at)}</p>
                  </div>
                  <div className="text-right text-sm">
                    {a.received && <p className="font-semibold text-green-700">+{a.received}</p>}
                    {a.spent && <p className="text-gray-600">−{a.spent}</p>}
                  </div>
                </div>
              ))}
            </Card>
          ))}
      </>
    );

  return (
    <div className="flex flex-col space-y-5 p-3 max-w-3xl">
      <Header />
      <div className="flex items-center justify-between px-2">
        <h1 className="text-xl font-bold text-primary-800">My Stocks</h1>
        <div className="flex gap-3">
          <button className="text-primary-800 font-semibold underline" onClick={() => navigate('/dashboard/public-markets/dividends')}>
            Dividends
          </button>
          <button className="text-primary-800 font-semibold underline" onClick={() => navigate('/dashboard/public-markets')}>
            Market
          </button>
        </div>
      </div>
      {body}
    </div>
  );
}

function OrderRow({ order: o, onClick }: { order: PMOrder; onClick: () => void }) {
  const buy = o.type === 'CREATION';
  return (
    <div onClick={onClick} className="flex items-center justify-between gap-3 py-2 border-b last:border-0 border-gray-100 cursor-pointer">
      <div>
        <p className="font-semibold">
          {buy ? 'Buy' : 'Sell'} {o.assetCode}
        </p>
        <p className="text-xs text-gray-400">
          {buy ? money(o.amount, o.fundingAssetCode) : `${qty(o.quantity)} units`} · {when(o.createdAt)}
        </p>
      </div>
      <Pill label={stateLabel(o.state)} color={stateColor(o.state)} />
    </div>
  );
}
