import { useEffect, useMemo, useState } from 'react';
import { useSelector } from 'react-redux';
import { useLocation, useNavigate, useParams } from 'react-router-dom';
import Header from '../../../components/header';
import { P2PEmptyState, P2PListCard, p2p } from '../../../components/p2p/P2PTheme';
import { useP2PIdentity } from '../../../hooks/useP2PIdentity';
import { useGetPMAssetQuery, useGetPMPortfolioQuery, useGetPMQuoteQuery, useTradePMMutation } from '../../../store/api/publicMarketsApis';
import { RootState } from '../../../store/reduxStore';
import { signBase64Txn } from '../../../utils/trovoSDK';
import { showNotification } from '../../../utils/showToaster';
import { AssetLogo, Chips, errorMessage, Loading, money, num, qty, Row } from './ui';

type Side = 'buy' | 'sell';

const PATHS: Record<string, string> = {
  FAST: 'Filled now from inventory',
  SLOW: 'Executed in the next session',
  NETTED: 'Netted with other orders',
};

// PublicMarketsTrade buys (an amount of the funding asset, cNGN) or sells (a
// quantity of tokens) from one of the user's own wallets. The quote is
// refreshed as the amount changes; Review shows the backend's messages,
// then the operation is signed with the wallet key and the order placed.
export default function PublicMarketsTrade() {
  const { code = '' } = useParams();
  const navigate = useNavigate();
  const location = useLocation();
  const [side, setSide] = useState<Side>((location.state as { side?: Side } | null)?.side ?? 'buy');
  const { creds, ready } = useP2PIdentity();
  const appUser = useSelector((state: RootState) => state.auth.user);
  // the user's own wallets (shared wallets with approvers cannot trade here)
  const wallets = useMemo(() => (appUser?.userWallets ?? []).filter((w) => !w.isSharedWallet && !w.sharedAccessEnabled), [appUser]);
  // the chosen wallet, the main one until the user picks another
  const [picked, setWalletAddress] = useState('');
  const walletAddress = picked || (wallets.find((w) => w.primaryWallet) ?? wallets[0])?.address || '';
  const wallet = wallets.find((w) => w.address === walletAddress);

  const { data: asset, isLoading } = useGetPMAssetQuery(code);
  const { data: portfolio } = useGetPMPortfolioQuery({ creds }, { skip: !ready });
  const holding = portfolio?.holdings.find((h) => h.asset.assetCode === code);
  const [value, setValue] = useState('');
  const debounced = useDebounced(value, 400);
  const buy = side === 'buy';

  const fundingBalance = num(
    wallet?.claimedAssets?.find((a) => a.assetCode?.toUpperCase() === (asset?.fundingAsset || 'CNGN').toUpperCase())?.amount,
  );
  const heldInWallet = num(holding?.wallets?.[walletAddress] ?? holding?.wallets?.[walletAddress.toLowerCase()]);

  const problem = (() => {
    if (!asset || !value) return '';
    const v = num(value);
    if (!(v > 0)) return 'Enter an amount above zero.';
    if (buy && num(asset.minimumBuy) > 0 && v < num(asset.minimumBuy)) return `The smallest buy is ${money(asset.minimumBuy, asset.fundingAsset)}.`;
    if (buy && v > fundingBalance) return `Not enough ${asset.fundingAsset} in this wallet.`;
    if (!buy && v > heldInWallet) return `This wallet holds ${qty(heldInWallet)} ${asset.ticker}.`;
    return '';
  })();

  const { data: quote, isFetching: quoting, error: quoteError } = useGetPMQuoteQuery(
    { creds, code, side, value: debounced },
    { skip: !ready || !debounced || num(debounced) <= 0 || !!problem },
  );
  const [trade, { isLoading: placing }] = useTradePMMutation();

  if (isLoading) return <Page code={code}><Loading /></Page>;
  if (!asset) return <Page code={code}><P2PEmptyState message="Could not load this asset." /></Page>;
  if (asset.status !== 'open') return <Page code={code}><P2PEmptyState message="This asset cannot be traded right now." /></Page>;
  if (!wallets.length) return <Page code={code}><P2PEmptyState message="You need a wallet of your own to trade." /></Page>;

  const presets = buy
    ? [1000, 5000, 10000, 50000].map((n) => ({ label: money(n, asset.fundingAsset), value: String(n) }))
    : [0.25, 0.5, 1].map((f) => ({ label: f === 1 ? 'All' : `${f * 100}%`, value: String(+(heldInWallet * f).toFixed(asset.tokenDecimals || 6)) }));

  const submit = async () => {
    if (problem || !num(value)) {
      showNotification('error', problem || 'Enter an amount.');
      return;
    }
    const body = { walletAddress, ...(buy ? { amount: value } : { quantity: value }) };
    try {
      const built = await trade({ creds, code, side, body }).unwrap();
      if (!built.transaction) throw new Error('Could not prepare the order');
      const q = built.quote;
      const summary = buy
        ? [`Buy about ${qty(q.quantity)} ${asset.ticker} for ${money(q.amount, q.fundingAsset)}`, `Trovo fee: ${money(q.fee, q.fundingAsset)} (${q.feePercent}%)`]
        : [`Sell ${qty(q.quantity)} ${asset.ticker}`, `Trovo fee: ${money(q.fee, q.fundingAsset)} (${q.feePercent}%)`, `You receive about ${money(q.netAmount, q.fundingAsset)}`];
      const lines = [...summary, q.note || PATHS[q.path] || '', '', ...(built.messages ?? []), '', buy ? 'Place this buy order?' : 'Place this sell order?'];
      if (!window.confirm(lines.filter((l, i, a) => l !== '' || a[i - 1] !== '').join('\n'))) return;
      const transactionSignature = signBase64Txn(creds.secretKey, built.transaction, '');
      const placed = await trade({ creds, code, side, body: { ...body, transaction: built.transaction, transactionSignature } }).unwrap();
      showNotification('success', buy ? 'Buy order placed.' : 'Sell order placed.');
      if (placed.order) navigate(`/dashboard/public-markets/order/${encodeURIComponent(placed.order.id)}`, { replace: true });
    } catch (e) {
      showNotification('error', errorMessage(e, 'The order could not be placed'));
    }
  };

  return (
    <Page code={code}>
      <div className="flex items-center gap-3 px-2">
        <AssetLogo asset={asset} />
        <div className="flex-1">
          <p className="font-bold text-primary-800">{asset.ticker}</p>
          <p className="text-sm text-gray-500">{money(asset.price, asset.currency)} per unit</p>
        </div>
      </div>
      <div className="px-2">
        <Chips
          items={[
            { key: 'buy' as Side, label: 'Buy' },
            { key: 'sell' as Side, label: 'Sell' },
          ]}
          value={side}
          onChange={(s) => {
            setSide(s);
            setValue('');
          }}
        />
      </div>
      <P2PListCard>
        <div className="flex flex-col space-y-4">
          <label className="flex flex-col space-y-1">
            <span className="text-sm font-semibold text-gray-600">Wallet</span>
            <select value={walletAddress} onChange={(e) => setWalletAddress(e.target.value)} className="rounded-xl border border-gray-200 px-3 py-3 bg-white">
              {wallets.map((w) => (
                <option key={w.address} value={w.address}>
                  {w.alias || (w.primaryWallet ? 'Main wallet' : `${w.address.slice(0, 6)}…${w.address.slice(-4)}`)}
                </option>
              ))}
            </select>
          </label>
          <label className="flex flex-col space-y-1">
            <span className="text-sm font-semibold text-gray-600">{buy ? `Amount (${asset.fundingAsset})` : `Quantity of ${asset.ticker}`}</span>
            <input
              value={value}
              onChange={(e) => setValue(e.target.value.replace(/[^0-9.]/g, ''))}
              inputMode="decimal"
              placeholder="0"
              className="rounded-xl border border-gray-200 px-3 py-3 text-lg"
            />
            <span className="text-xs text-gray-500">
              Available: {buy ? money(fundingBalance, asset.fundingAsset) : `${qty(heldInWallet)} ${asset.ticker}`}
            </span>
          </label>
          <div className="flex flex-wrap gap-2">
            {presets.map((p) => (
              <button key={p.label} onClick={() => setValue(p.value)} className="px-3 py-1 rounded-lg bg-gray-100 text-sm font-semibold">
                {p.label}
              </button>
            ))}
          </div>
          {problem && <p className="text-sm" style={{ color: p2p.danger }}>{problem}</p>}
          {!problem && quoteError && <p className="text-sm" style={{ color: p2p.danger }}>{errorMessage(quoteError, 'Could not get a quote')}</p>}
          {!problem && quote && (
            <div className="rounded-xl bg-gray-50 p-4">
              {buy ? (
                <>
                  <Row label="You pay" value={money(quote.amount, quote.fundingAsset)} />
                  <Row label="Trovo fee" value={`${money(quote.fee, quote.fundingAsset)} (${quote.feePercent}%)`} />
                  <Row label="You get about" value={`${qty(quote.quantity)} ${asset.ticker}`} />
                </>
              ) : (
                <>
                  <Row label="You sell" value={`${qty(quote.quantity)} ${asset.ticker}`} />
                  <Row label="Trovo fee" value={`${money(quote.fee, quote.fundingAsset)} (${quote.feePercent}%)`} />
                  <Row label="You receive about" value={money(quote.netAmount, quote.fundingAsset)} />
                </>
              )}
              <Row label="Price" value={`${money(quote.price, quote.currency)} (${quote.priceSource})`} />
              <Row label="How it fills" value={PATHS[quote.path] ?? quote.path} />
              {quote.note && <p className="text-xs text-gray-500 pt-2">{quote.note}</p>}
              {quote.settlementNote && <p className="text-xs text-gray-500 pt-1">{quote.settlementNote}</p>}
              {quote.custodianName && <p className="text-xs text-gray-400 pt-1">Units held for you by {quote.custodianName}.</p>}
            </div>
          )}
          {quoting && <p className="text-xs text-gray-400">Updating the quote…</p>}
          <button
            onClick={submit}
            disabled={placing || !!problem || !num(value)}
            className="py-3 rounded-xl font-bold text-white disabled:opacity-40"
            style={{ backgroundColor: buy ? p2p.success : p2p.danger }}
          >
            {placing ? 'Placing…' : buy ? `Review buy` : `Review sell`}
          </button>
          <p className="text-xs text-gray-400">
            Prices move: the final quantity or amount is set when the order fills. Tokens are backed one-to-one by units held at the CSCS.
          </p>
        </div>
      </P2PListCard>
    </Page>
  );
}

function useDebounced<T>(value: T, ms: number) {
  const [v, setV] = useState(value);
  useEffect(() => {
    const t = setTimeout(() => setV(value), ms);
    return () => clearTimeout(t);
  }, [value, ms]);
  return v;
}

function Page({ code, children }: { code: string; children: React.ReactNode }) {
  const navigate = useNavigate();
  return (
    <div className="flex flex-col space-y-5 p-3 max-w-2xl">
      <Header />
      <button className="self-start text-primary-800 font-semibold px-2" onClick={() => navigate(`/dashboard/public-markets/asset/${encodeURIComponent(code)}`)}>
        ‹ Back
      </button>
      {children}
    </div>
  );
}
