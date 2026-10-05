import React from 'react';
import { p2p } from '../../../components/p2p/P2PTheme';
import { PMAsset } from '../../../types/publicMarkets';

// Shared pieces of the Public Markets pages.

export const errorMessage = (e: any, fallback: string) => e?.data?.message || e?.data?.error || e?.message || fallback;

export const num = (v?: string | number) => {
  const n = Number(v ?? '');
  return Number.isFinite(n) ? n : 0;
};

export const money = (v?: string | number, code = 'NGN') => {
  if (v === undefined || v === null || v === '') return '—';
  const n = Number(v);
  if (!Number.isFinite(n)) return String(v);
  const s = n.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 });
  return code === 'NGN' || code === 'CNGN' ? `₦${s}` : `${s} ${code}`;
};

export const qty = (v?: string | number) => {
  if (v === undefined || v === null || v === '') return '—';
  const n = Number(v);
  return Number.isFinite(n) ? n.toLocaleString(undefined, { maximumFractionDigits: 6 }) : String(v);
};

export const pct = (v?: string | number) => {
  if (v === undefined || v === null || v === '') return '—';
  const n = Number(v);
  return Number.isFinite(n) ? `${n >= 0 ? '+' : ''}${n.toFixed(2)}%` : String(v);
};

export const when = (v?: string) => (v ? new Date(v).toLocaleString() : '—');

export const changeColor = (v?: string | number) => (num(v) > 0 ? p2p.success : num(v) < 0 ? p2p.danger : '#667085');

export const ORDER_STATE_LABELS: Record<string, string> = {
  'awaiting-payment': 'Awaiting payment',
  queued: 'Queued for the next session',
  'filled-from-inventory': 'Filled from inventory',
  'pending-execution': 'Awaiting execution',
  executed: 'Executed on the exchange',
  submitted: 'Submitted to the chain',
  chain_final: 'Delivered to your wallet',
  settlement_final: 'Settled',
  complete: 'Complete',
  rejected: 'Rejected',
  failed: 'Failed',
  cancelled: 'Cancelled',
};

export const stateLabel = (s: string) => ORDER_STATE_LABELS[s] ?? s.replace(/[_-]/g, ' ');

export const stateColor = (s: string) =>
  s === 'complete' || s === 'PAID'
    ? p2p.success
    : ['rejected', 'failed', 'cancelled', 'FAILED'].includes(s)
      ? p2p.danger
      : ['awaiting-payment', 'queued', 'pending-execution', 'PENDING'].includes(s)
        ? p2p.warning
        : p2p.primary;

export function Pill({ label, color }: { label: string; color: string }) {
  return (
    <span className="px-3 py-1 rounded-full text-xs font-semibold whitespace-nowrap" style={{ color, backgroundColor: `${color}1F` }}>
      {label}
    </span>
  );
}

export function AssetLogo({ asset, size = 44 }: { asset: Pick<PMAsset, 'logo' | 'ticker'>; size?: number }) {
  const logo = asset.logo ?? {};
  if (logo.url) return <img src={logo.url} alt={asset.ticker} className="rounded-xl object-cover" style={{ width: size, height: size }} />;
  return (
    <div
      className="rounded-xl grid place-items-center font-bold shrink-0"
      style={{ width: size, height: size, background: logo.background || p2p.brandDark, color: logo.foreground || '#fff', fontSize: size / 3 }}
    >
      {logo.initials || asset.ticker.slice(0, 2)}
    </div>
  );
}

export function StatusBadge({ asset }: { asset: PMAsset }) {
  if (asset.status === 'halted') return <Pill label="Halted" color={p2p.danger} />;
  if (asset.status === 'coming-soon') return <Pill label="Coming soon" color="#667085" />;
  return asset.session?.open ? <Pill label="Market open" color={p2p.success} /> : <Pill label="Market closed" color={p2p.warning} />;
}

export function Row({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div className="flex justify-between gap-4 py-1">
      <span className="text-gray-500">{label}</span>
      <span className="font-semibold text-right break-all">{value}</span>
    </div>
  );
}

export function Card({ title, children }: { title?: string; children: React.ReactNode }) {
  return (
    <div className="bg-white rounded-2xl p-4" style={{ boxShadow: p2p.cardShadow }}>
      {title && <p className="font-bold text-primary-800 mb-2">{title}</p>}
      {children}
    </div>
  );
}

export function Chips<T extends string>({ items, value, onChange }: { items: { key: T; label: string }[]; value: T; onChange: (k: T) => void }) {
  return (
    <div className="flex flex-wrap gap-2">
      {items.map((i) => (
        <button
          key={i.key}
          onClick={() => onChange(i.key)}
          className={`px-4 py-2 rounded-xl font-semibold text-sm ${value === i.key ? 'text-white' : 'text-primary-800 bg-white'}`}
          style={value === i.key ? { backgroundColor: p2p.brandDark } : undefined}
        >
          {i.label}
        </button>
      ))}
    </div>
  );
}

// PriceChart draws the price history as an SVG line (app-web has no chart
// library), green when the range ends up and red when it ends down.
export function PriceChart({ points, height = 180 }: { points: { at: string; price: string }[]; height?: number }) {
  const values = points.map((p) => num(p.price));
  if (values.length < 2) return <div className="text-center text-gray-400 py-12">No price history yet</div>;
  const min = Math.min(...values);
  const max = Math.max(...values);
  const span = max - min || 1;
  const w = 600;
  const path = values.map((v, i) => `${i === 0 ? 'M' : 'L'}${((i / (values.length - 1)) * w).toFixed(1)},${(height - 8 - ((v - min) / span) * (height - 16)).toFixed(1)}`).join(' ');
  const color = values[values.length - 1] >= values[0] ? p2p.success : p2p.danger;
  return (
    <svg viewBox={`0 0 ${w} ${height}`} preserveAspectRatio="none" className="w-full" style={{ height }}>
      <path d={`${path} L${w},${height} L0,${height} Z`} fill={`${color}14`} stroke="none" />
      <path d={path} fill="none" stroke={color} strokeWidth={2} vectorEffect="non-scaling-stroke" />
    </svg>
  );
}

export function Loading() {
  return <div className="text-center py-16 text-gray-400">Loading...</div>;
}
