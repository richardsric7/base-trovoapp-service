import React from 'react';
import { FaCircleCheck, FaCircleExclamation, FaStore } from 'react-icons/fa6';

// A small, real design-token set scoped to the new P2P pages only (Plan
// Section 94). Neither app-web nor app-mobile has a consistently-enforced
// design system today, so rather than perpetuate the existing ad hoc
// spacing/colors or invent unrelated new ones, this formalizes values
// already in use somewhere in the product: p2pBrandDark matches app-web's
// own history page's undocumented second brand blue, which is also
// app-mobile's `trovoblue` - the two platforms read as one product.
export const p2p = {
  brandDark: '#00225A',
  primary: '#004988', // existing Tailwind primary-800
  success: '#00A859',
  danger: '#BE3800', // existing Tailwind trovored
  warning: '#B7791F', // new - neither app has an existing amber token
  neutralBg: '#F2F6F9', // existing Tailwind primary-100
  // Matches the history page's own shadow value exactly, so P2P cards read
  // as the same product on both web and mobile.
  cardShadow: '0 18px 60px rgba(0,34,90,0.06)',
};

export function statusColor(status: string): string {
  switch (status) {
    case 'AWAITING_ESCROW_DEPOSIT':
    case 'AWAITING_PAYMENT':
      return p2p.warning;
    case 'AWAITING_PAYMENT_CONFIRMATION':
      return p2p.primary;
    case 'COMPLETED':
      return p2p.success;
    case 'REJECTED':
    case 'CANCELLED':
    case 'EXPIRED':
      return '#9EA3AE';
    default:
      return '#9EA3AE';
  }
}

export function statusLabel(status: string): string {
  switch (status) {
    case 'AWAITING_APPROVAL':
      return 'Awaiting approval';
    case 'AWAITING_ESCROW_DEPOSIT':
      return 'Awaiting escrow deposit';
    case 'AWAITING_PAYMENT':
      return 'Awaiting payment';
    case 'AWAITING_PAYMENT_CONFIRMATION':
      return 'Awaiting confirmation';
    case 'COMPLETED':
      return 'Completed';
    case 'REJECTED':
      return 'Rejected';
    case 'CANCELLED':
      return 'Cancelled';
    case 'EXPIRED':
      return 'Expired';
    default:
      return status;
  }
}

export function P2PStatusPill({ status, disputed }: { status: string; disputed?: boolean }) {
  const color = disputed ? p2p.danger : statusColor(status);
  const label = disputed ? 'Disputed' : statusLabel(status);
  return (
    <span
      className="px-3 py-1 rounded-full text-xs font-semibold"
      style={{ color, backgroundColor: `${color}1F` }}
    >
      {label}
    </span>
  );
}

export function P2PListCard({
  children,
  onClick,
}: {
  children: React.ReactNode;
  onClick?: () => void;
}) {
  return (
    <div
      onClick={onClick}
      className={`bg-white rounded-2xl p-4 mb-3 ${onClick ? 'cursor-pointer hover:bg-primary-100' : ''}`}
      style={{ boxShadow: p2p.cardShadow }}
    >
      {children}
    </div>
  );
}

export function P2PEmptyState({
  icon,
  message,
  ctaLabel,
  onCta,
}: {
  icon?: React.ReactNode;
  message: string;
  ctaLabel?: string;
  onCta?: () => void;
}) {
  return (
    <div className="flex flex-col items-center justify-center py-16 px-6 text-center">
      <div className="text-4xl mb-4 opacity-40">{icon ?? <FaStore />}</div>
      <p className="text-gray-500 mb-4">{message}</p>
      {ctaLabel && (
        <button onClick={onCta} className="text-primary-800 font-semibold underline">
          {ctaLabel}
        </button>
      )}
    </div>
  );
}

export function P2PBanner({ children, tone = 'danger' }: { children: React.ReactNode; tone?: 'danger' | 'success' }) {
  const color = tone === 'danger' ? p2p.danger : p2p.success;
  const Icon = tone === 'danger' ? FaCircleExclamation : FaCircleCheck;
  return (
    <div
      className="flex items-center gap-2 rounded-2xl p-3 mb-4"
      style={{ backgroundColor: `${color}14`, border: `1px solid ${color}4D`, color }}
    >
      <Icon />
      <span className="font-semibold">{children}</span>
    </div>
  );
}
