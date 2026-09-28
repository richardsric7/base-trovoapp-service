import { getAssetCode } from "./utilities";

export const normalizeTransactionType = (
    value: string | undefined,
    item?: Record<string, any>,
    activeAddress?: string,
  ) => {
    const type = `${value ?? ''}`.toLowerCase();

    if (type.includes('swap')) return 'Swap';
    if (type.includes('sent') || type.includes('send')) return 'Sent';
    if (type.includes('received') || type.includes('receive'))
      return 'Received';
    if (type.includes('payment')) {
      const fromAddress = `${item?.fromAddress ?? item?.senderAddress ?? ''}`;
      const toAddress = `${item?.toAddress ?? item?.receiverAddress ?? ''}`;

      if (activeAddress && fromAddress === activeAddress) return 'Sent';
      if (activeAddress && toAddress === activeAddress) return 'Received';
    }

    return 'Received';
  };

export const parseAmountValue = (amount: unknown) => {
  if (typeof amount === 'number') return amount;

  if (typeof amount === 'string') {
    const cleaned = amount.replace(/[^0-9.-]/g, '');
    const parsed = Number(cleaned);
    return Number.isNaN(parsed) ? 0 : Math.abs(parsed);
  }

  return 0;
};

export const formatAmount = (amount: number, assetCode: string) => {
  return `${amount.toLocaleString('en-NG', {
    minimumFractionDigits: 4,
    maximumFractionDigits: 4,
  })} ${assetCode}`;
};

export const parseDate = (value: unknown) => {
  if (!value) return null;
  const parsed = new Date(String(value));
  return Number.isNaN(parsed.getTime()) ? null : parsed;
};

export const formatRelativeTime = (date: Date | null) => {
  if (!date) return '-';

  const now = new Date();
  const diffInMs = now.getTime() - date.getTime();
  const diffInDays = Math.floor(diffInMs / (1000 * 60 * 60 * 24));

  if (diffInDays <= 0) return 'Today';
  if (diffInDays === 1) return '1 day ago';

  return `${diffInDays} days ago`;
};

// getTransactionAssetCode returns the asset that arrived at the recipient
// (destinationAssetCode). The payment-history API now splits a payment into
// a source and destination side (see getTransactionSourceAssetCode) - for a
// plain payment they're identical, so most display code that only shows one
// asset code should keep using this one.
export const getTransactionAssetCode = (item: Record<string, any>) => {
  if (item.destinationAssetCode !== undefined && item.destinationAssetCode !== '') {
    return getAssetCode(item.destinationAssetCode);
  }

  if (item.assetCode !== undefined) {
    return getAssetCode(item.assetCode);
  }

  if (item.asset?.assetCode !== undefined) {
    return getAssetCode(item.asset.assetCode);
  }

  if (typeof item.amount === 'string') {
    const segments = item.amount.trim().split(/\s+/);
    const lastSegment = segments[segments.length - 1];

    if (/[A-Za-z]/.test(lastSegment)) {
      return lastSegment.toUpperCase();
    }
  }

  return 'ETH';
};

// getTransactionSourceAssetCode returns the asset that left the sender. For
// a plain payment this is the same as getTransactionAssetCode; when it
// differs, the row is a swap.
export const getTransactionSourceAssetCode = (item: Record<string, any>) => {
  if (item.sourceAssetCode !== undefined && item.sourceAssetCode !== '') {
    return getAssetCode(item.sourceAssetCode);
  }
  return getTransactionAssetCode(item);
};

// getTransactionAmount/getTransactionSourceAmount read the split
// destination/source amount fields, falling back to the legacy flat
// `amount` field for any caller that hasn't migrated.
export const getTransactionAmount = (item: Record<string, any>) =>
  item.destinationAmount ?? item.amount;

export const getTransactionSourceAmount = (item: Record<string, any>) =>
  item.sourceAmount ?? item.amount;

export const getTransactionSourceNetwork = (item: Record<string, any>) =>
  `${item.sourceNetwork ?? 'base'}`;

export const getTransactionDestinationNetwork = (item: Record<string, any>) =>
  `${item.destinationNetwork ?? 'base'}`;

// isTransactionSwap mirrors the backend's definition: a row where the
// source and destination side differ (asset or network) is a swap.
export const isTransactionSwap = (item: Record<string, any>) => {
  const sourceAsset = getTransactionSourceAssetCode(item);
  const destinationAsset = getTransactionAssetCode(item);
  const sourceNetwork = getTransactionSourceNetwork(item);
  const destinationNetwork = getTransactionDestinationNetwork(item);
  return sourceAsset !== destinationAsset || sourceNetwork !== destinationNetwork;
};