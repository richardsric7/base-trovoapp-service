import { useLocation, useNavigate } from 'react-router-dom';
import Header from '../../../components/header';
import { P2PEmptyState, P2PListCard } from '../../../components/p2p/P2PTheme';
import { useP2PIdentity } from '../../../hooks/useP2PIdentity';
import { useGetPMDividendsQuery } from '../../../store/api/publicMarketsApis';
import { Loading, money, Pill, qty, Row, stateColor, when } from './ui';

const STATUS: Record<string, string> = { PENDING: 'Upcoming', QUEUED: 'Being paid', PAID: 'Paid', FAILED: 'Being retried' };

// PublicMarketsDividends lists the dividends and coupons the user is
// entitled to as a token holder on the record date, paid in cNGN to the
// wallet that held the tokens (net of withholding tax).
export default function PublicMarketsDividends() {
  const navigate = useNavigate();
  const assetCode = (useLocation().state as { assetCode?: string } | null)?.assetCode;
  const { creds, ready } = useP2PIdentity();
  const { data, isLoading, isError, refetch } = useGetPMDividendsQuery({ creds, assetCode }, { skip: !ready });
  const list = data?.dividends ?? [];

  let body: React.ReactNode;
  if (!ready || isLoading) body = <Loading />;
  else if (isError) body = <P2PEmptyState message="Could not load your payments." ctaLabel="Retry" onCta={refetch} />;
  else if (list.length === 0)
    body = <P2PEmptyState message="No dividends or coupons yet. Hold a stock or bond on its record date to receive its payments." />;
  else
    body = (
      <div>
        {list.map((d) => (
          <P2PListCard key={d.id}>
            <div className="flex items-center justify-between mb-2">
              <p className="font-bold text-primary-800">
                {d.assetCode} · {d.eventType === 'COUPON' ? 'Coupon' : 'Dividend'}
              </p>
              <Pill label={STATUS[d.status] ?? d.status} color={stateColor(d.status)} />
            </div>
            {d.description && <p className="text-sm text-gray-500 mb-2">{d.description}</p>}
            <Row label="Units on the record date" value={qty(d.units)} />
            <Row label="Per unit" value={money(d.amountPerUnit)} />
            <Row label="Gross" value={money(d.grossAmount)} />
            <Row label={`Withholding tax (${d.whtPercent}%)`} value={money(d.whtAmount)} />
            <Row label="You receive" value={money(d.netAmount)} />
            <Row label="Record date" value={d.recordDate} />
            <Row label="Pay date" value={d.payDate || '—'} />
            {d.paidAt && <Row label="Paid" value={when(d.paidAt)} />}
            {d.txHash && <Row label="Transaction" value={<span className="font-mono text-xs">{d.txHash}</span>} />}
            {d.custodianName && <p className="text-xs text-gray-400 pt-2">Received from the registrar by {d.custodianName}.</p>}
          </P2PListCard>
        ))}
      </div>
    );

  return (
    <div className="flex flex-col space-y-5 p-3 max-w-2xl">
      <Header />
      <div className="flex items-center justify-between px-2">
        <h1 className="text-xl font-bold text-primary-800">Dividends &amp; coupons{assetCode ? ` · ${assetCode}` : ''}</h1>
        <button className="text-primary-800 font-semibold underline" onClick={() => navigate('/dashboard/public-markets/portfolio')}>
          My Stocks
        </button>
      </div>
      {body}
    </div>
  );
}
