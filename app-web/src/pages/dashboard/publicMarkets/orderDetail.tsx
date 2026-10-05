import { useNavigate, useParams } from 'react-router-dom';
import Header from '../../../components/header';
import { P2PEmptyState } from '../../../components/p2p/P2PTheme';
import { useP2PIdentity } from '../../../hooks/useP2PIdentity';
import { publicMarketsApi, useGetPMOrderQuery } from '../../../store/api/publicMarketsApis';
import { PM_FINAL_STATES } from '../../../types/publicMarkets';
import { Card, Loading, money, Pill, qty, Row, stateColor, stateLabel, when } from './ui';

const PATHS: Record<string, string> = { FAST: 'Filled from inventory', SLOW: 'Executed in the next session', NETTED: 'Netted with other orders' };

// PublicMarketsOrder follows one order until it is final, refreshing every
// five seconds while it is in progress.
export default function PublicMarketsOrder() {
  const { orderId = '' } = useParams();
  const navigate = useNavigate();
  const { creds, ready } = useP2PIdentity();
  const args = { creds, id: orderId };
  // the cached order, read without subscribing, stops the polling once final
  const { data: cached } = publicMarketsApi.endpoints.getPMOrder.useQueryState(args, { skip: !ready });
  const final = !!cached && PM_FINAL_STATES.includes(cached.state);
  const { data: o, isLoading, isError, refetch } = useGetPMOrderQuery(args, { skip: !ready, pollingInterval: final ? 0 : 5000 });

  let body: React.ReactNode;
  if (!ready || isLoading) body = <Loading />;
  else if (isError || !o) body = <P2PEmptyState message="Could not load this order." ctaLabel="Retry" onCta={refetch} />;
  else {
    const buy = o.type === 'CREATION';
    body = (
      <>
        <Card>
          <div className="flex items-center justify-between">
            <p className="text-lg font-bold text-primary-800">
              {buy ? 'Buy' : 'Sell'} {o.assetCode}
            </p>
            <Pill label={stateLabel(o.state)} color={stateColor(o.state)} />
          </div>
          {o.note && <p className="text-sm text-gray-500 mt-2">{o.note}</p>}
        </Card>
        <Card title="Details">
          <Row label={buy ? 'Amount paid' : 'Proceeds'} value={money(buy ? o.amount : o.netAmount, o.fundingAssetCode)} />
          <Row label="Trovo fee" value={money(o.fee, o.fundingAssetCode)} />
          <Row label="Quantity" value={`${qty(o.quantity)} ${o.assetCode}`} />
          <Row label="Quoted price" value={money(o.referencePrice)} />
          {o.executedPrice && <Row label="Executed price" value={money(o.executedPrice)} />}
          <Row label="How it fills" value={PATHS[o.path] ?? (o.path || '—')} />
          <Row label="Wallet" value={o.walletAlias || o.walletAddress} />
          <Row label="Placed" value={when(o.createdAt)} />
          <Row label="Order" value={<span className="font-mono text-xs">{o.id}</span>} />
          {o.paymentTxHash && <Row label="Payment" value={<span className="font-mono text-xs">{o.paymentTxHash}</span>} />}
          {o.tokenTxHash && <Row label="Tokens" value={<span className="font-mono text-xs">{o.tokenTxHash}</span>} />}
          {o.payoutTxHash && <Row label="Payout" value={<span className="font-mono text-xs">{o.payoutTxHash}</span>} />}
        </Card>
        {(o.events ?? []).length > 0 && (
          <Card title="Progress">
            <ol className="space-y-3">
              {o.events!.map((e, i) => (
                <li key={i} className="flex gap-3">
                  <span className="mt-1 h-3 w-3 rounded-full shrink-0" style={{ backgroundColor: stateColor(e.state) }} />
                  <div>
                    <p className="font-semibold">{stateLabel(e.state)}</p>
                    {e.note && <p className="text-sm text-gray-500">{e.note}</p>}
                    <p className="text-xs text-gray-400">{when(e.at)}</p>
                  </div>
                </li>
              ))}
            </ol>
          </Card>
        )}
      </>
    );
  }

  return (
    <div className="flex flex-col space-y-5 p-3 max-w-2xl">
      <Header />
      <div className="flex gap-4 px-2">
        <button className="text-primary-800 font-semibold" onClick={() => navigate('/dashboard/public-markets/portfolio')}>
          ‹ My Stocks
        </button>
      </div>
      {body}
    </div>
  );
}
