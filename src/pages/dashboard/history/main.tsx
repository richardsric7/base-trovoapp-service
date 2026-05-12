import { ReactNode, useEffect, useState } from 'react';
import { useSelector } from 'react-redux';
import {
  Document,
  Page,
  Text,
  View,
  StyleSheet,
  pdf,
} from '@react-pdf/renderer';
import Header from '../../../components/header';
import Button from '../../../components/button';
import ButtonSecondary from '../../../components/buttonSecondary';
import { RootState } from '../../../store/reduxStore';
import { useFetchFiatPaymentsQuery } from '../../../store/api/walletApis';
import { Encryptor } from '../../../utils/encryptor';
import { getAssetCode } from '../../../utils/utilities';
import styles from './main.module.css';
import { transactionTypeConfig } from './data';
import {
  formatAmount,
  formatRelativeTime,
  getTransactionAssetCode,
  normalizeTransactionType,
  parseAmountValue,
  parseDate,
} from '../../../utils/transactionUtils';

type FilterState = {
  username: string;
  fromPublicKey: string;
  toPublicKey: string;
  memo: string;
};

type DateRangeKey = 'none' | 'week' | 'month' | 'quarter' | 'custom';

type HistoryRow = {
  id: string;
  type: string;
  assetCode: string;
  amountValue: number;
  price: string;
  description: string;
  dateLabel: string;
  dateValue: Date | null;
  dateRaw: string;
  walletKey: string;
  walletLabel: string;
  username: string;
  fromName: string;
  toName: string;
  fromPublicKey: string;
  toPublicKey: string;
  transactionId: string;
  memo: string;
};

const defaultFilters: FilterState = {
  username: '',
  fromPublicKey: '',
  toPublicKey: '',
  memo: '',
};

export const History = () => {
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const [secretKey, setSecretKey] = useState('');
  const [selectedWallet, setSelectedWallet] = useState('');
  const [selectedAsset, setSelectedAsset] = useState('');
  const [selectedType, setSelectedType] = useState('');
  const [showAmountModal, setShowAmountModal] = useState(false);
  const [showDateModal, setShowDateModal] = useState(false);
  const [showMoreFiltersModal, setShowMoreFiltersModal] = useState(false);
  const [minAmountInput, setMinAmountInput] = useState('');
  const [maxAmountInput, setMaxAmountInput] = useState('');
  const [appliedMinAmount, setAppliedMinAmount] = useState<number | null>(null);
  const [appliedMaxAmount, setAppliedMaxAmount] = useState<number | null>(null);
  const [dateRangeKey, setDateRangeKey] = useState<DateRangeKey>('none');
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [draftStartDate, setDraftStartDate] = useState('');
  const [draftEndDate, setDraftEndDate] = useState('');
  const [draftDateRangeKey, setDraftDateRangeKey] =
    useState<DateRangeKey>('none');
  const [filters, setFilters] = useState<FilterState>(defaultFilters);
  const [draftFilters, setDraftFilters] = useState<FilterState>(defaultFilters);
  const [selectedHistoryRow, setSelectedHistoryRow] =
    useState<HistoryRow | null>(null);
  const [shareLoading, setShareLoading] = useState<
    'image' | 'pdf' | 'text' | null
  >(null);
  const [shareFeedback, setShareFeedback] = useState('');

  useEffect(() => {
    const encryptor = new Encryptor();
    encryptor.getSecretKey(appUser).then(setSecretKey);
  }, [appUser]);

  const wallets = appUser?.userWallets ?? [];
  const primaryWallet =
    wallets.find((wallet) => wallet.primaryWallet) ?? wallets?.[0];

  const forPublicKey =
    selectedWallet !== '' ? selectedWallet : (primaryWallet?.publicKey ?? '');

  const apiDateRange = buildDateRange(dateRangeKey, startDate, endDate);

  const transactionTypeParam =
    selectedType === '' ? '' : selectedType === 'Swap' ? 'swap' : 'payment';

  const queryParams: string[] = [];
  if (filters.username)
    queryParams.push(`&name=${encodeURIComponent(filters.username)}`);
  if (filters.memo)
    queryParams.push(`&memo=${encodeURIComponent(filters.memo)}`);
  if (filters.fromPublicKey)
    queryParams.push(
      `&fromPublicKey=${encodeURIComponent(filters.fromPublicKey)}`,
    );
  if (filters.toPublicKey)
    queryParams.push(`&toPublicKey=${encodeURIComponent(filters.toPublicKey)}`);
  if (appliedMinAmount !== null || appliedMaxAmount !== null) {
    const min = appliedMinAmount !== null ? appliedMinAmount : '';
    const max = appliedMaxAmount !== null ? appliedMaxAmount : '';
    queryParams.push(`&amount=${encodeURIComponent(`${min}%${max}`)}`);
  }
  if (apiDateRange.start && apiDateRange.end) {
    queryParams.push(
      `&dateBetween=${encodeURIComponent(
        `${apiDateRange.start}%${apiDateRange.end}`,
      )}`,
    );
  }
  queryParams.push(`&transactionType=${transactionTypeParam}`);

  const { data, isLoading } = useFetchFiatPaymentsQuery(
    {
      signer: primaryWallet?.signer ?? '',
      publicKey: forPublicKey,
      secretKey,
      body: { limit: 50, query: queryParams.join('') },
    },
    { skip: !secretKey || !primaryWallet },
  );

  const walletOptions = [
    { label: 'All wallets', value: '' },
    ...wallets.map((wallet) => ({
      label: wallet.alias,
      value: wallet.publicKey,
    })),
  ];

  const assetSet = new Set<string>();
  wallets.forEach((wallet) => {
    wallet.claimedAssets.forEach((asset) => {
      assetSet.add(getAssetCode(asset.assetCode));
    });
  });

  const assetOptions = [
    { label: 'All assets', value: '' },
    ...Array.from(assetSet).map((assetCode) => ({
      label: assetCode,
      value: assetCode,
    })),
  ];

  const transactionOptions = [
    { label: 'All transactions', value: '' },
    { label: 'Received', value: 'Received' },
    { label: 'Sent', value: 'Sent' },
    { label: 'Swap', value: 'Swap' },
  ];

  const rawRecords = data?.data?.records ?? data?.records ?? [];

  const rows: HistoryRow[] = rawRecords.map(
    (item: Record<string, any>, index: number) => {
      const type = normalizeTransactionType(
        item.transactionType ?? item.type,
        item,
        forPublicKey,
      );
      const amountValue = parseAmountValue(item.amount);
      const assetCode = getTransactionAssetCode(item);
      const amountPrefix =
        type === 'Sent' ? '-' : type === 'Received' ? '+' : '+';
      const amountText =
        typeof item.amount === 'string' && item.amount.trim().length > 0
          ? item.amount
          : formatAmount(amountValue, assetCode);
      const createdAt = parseDate(
        item.transactionDate ?? item.createdAt ?? item.date,
      );
      const description =
        item.description ??
        item.narration ??
        item.memo ??
        (type === 'Received'
          ? `Received ${assetCode}`
          : type === 'Sent'
            ? `Sent ${assetCode}`
            : `Swapped asset to ${assetCode}`);

      const walletMatch =
        wallets.find(
          (wallet) =>
            wallet.publicKey === item.publicKey ||
            wallet.publicKey === item.walletPublicKey ||
            wallet.alias === item.walletAlias,
        ) ?? primaryWallet;

      return {
        id: `${item.transactionId ?? item.id ?? item._id ?? item.reference ?? index}`,
        type,
        assetCode,
        amountValue,
        price: `${amountPrefix}${amountText}`,
        description,
        dateLabel: formatRelativeTime(createdAt),
        dateValue: createdAt,
        dateRaw: `${item.transactionDate ?? item.createdAt ?? item.date ?? ''}`,
        walletKey: walletMatch?.publicKey ?? 'unknown',
        walletLabel: walletMatch?.alias ?? 'Primary wallet',
        username: `${item.username ?? item.fullName ?? item.name ?? ''}`.trim(),
        fromName: `${item.from ?? item.senderName ?? ''}`.trim(),
        toName: `${item.to ?? item.receiverName ?? ''}`.trim(),
        fromPublicKey: `${item.fromPublicKey ?? item.senderPublicKey ?? ''}`,
        toPublicKey: `${item.toPublicKey ?? item.receiverPublicKey ?? ''}`,
        transactionId: `${item.transactionId ?? item.id ?? item._id ?? item.reference ?? ''}`,
        memo: `${item.memo ?? item.narration ?? ''}`.trim(),
      };
    },
  );

  const filteredRows = rows.filter((row) => {
    if (selectedAsset !== '' && row.assetCode !== selectedAsset) {
      return false;
    }
    if (selectedType === 'Sent' && row.type !== 'Sent') return false;
    if (selectedType === 'Received' && row.type !== 'Received') return false;
    if (
      !matchesDateRange(row.dateValue, apiDateRange.start, apiDateRange.end)
    ) {
      return false;
    }
    return true;
  });

  const amountLabel = getAmountLabel(appliedMinAmount, appliedMaxAmount);
  const dateLabel = getDateLabel(dateRangeKey, startDate, endDate);

  const openAmountModal = () => {
    setMinAmountInput(appliedMinAmount?.toString() ?? '');
    setMaxAmountInput(appliedMaxAmount?.toString() ?? '');
    setShowAmountModal(true);
  };

  const openDateModal = () => {
    setDraftDateRangeKey(dateRangeKey);
    setDraftStartDate(startDate);
    setDraftEndDate(endDate);
    setShowDateModal(true);
  };

  const openMoreFiltersModal = () => {
    setDraftFilters(filters);
    setShowMoreFiltersModal(true);
  };

  const applyAmountFilter = () => {
    setAppliedMinAmount(
      minAmountInput.trim().length > 0 ? Number(minAmountInput) : null,
    );
    setAppliedMaxAmount(
      maxAmountInput.trim().length > 0 ? Number(maxAmountInput) : null,
    );
    setShowAmountModal(false);
  };

  const applyDateFilter = () => {
    setDateRangeKey(draftDateRangeKey);
    setStartDate(draftStartDate);
    setEndDate(draftEndDate);
    setShowDateModal(false);
  };

  const applyMoreFilters = () => {
    setFilters(draftFilters);
    setShowMoreFiltersModal(false);
  };

  const resetAllFilters = () => {
    setSelectedWallet('');
    setSelectedAsset('');
    setSelectedType('');
    setMinAmountInput('');
    setMaxAmountInput('');
    setAppliedMinAmount(null);
    setAppliedMaxAmount(null);
    setDateRangeKey('none');
    setStartDate('');
    setEndDate('');
    setDraftDateRangeKey('none');
    setDraftStartDate('');
    setDraftEndDate('');
    setFilters(defaultFilters);
    setDraftFilters(defaultFilters);
  };

  const handleShareText = async () => {
    if (!selectedHistoryRow) return;
    setShareLoading('text');
    setShareFeedback('');

    const message = buildHistoryShareText(selectedHistoryRow);

    try {
      if (navigator.share) {
        await navigator.share({
          title: 'Trovo Payment Details',
          text: message,
        });
        setShareFeedback('Details shared successfully.');
      } else {
        await navigator.clipboard.writeText(message);
        setShareFeedback('Details copied to clipboard.');
      }
    } catch (error: any) {
      if (error?.name !== 'AbortError') {
        setShareFeedback('Unable to share text right now.');
      }
    } finally {
      setShareLoading(null);
    }
  };

  const handleShareImage = async () => {
    if (!selectedHistoryRow) return;
    setShareLoading('image');
    setShareFeedback('');

    try {
      const imageBlob = await generateHistoryImageBlob(selectedHistoryRow);
      const fileName = `trovo-payment-${safeFileDate(selectedHistoryRow.dateValue)}.png`;
      const file = new File([imageBlob], fileName, { type: 'image/png' });
      const shared = await shareFileWithFallback(
        file,
        'Trovo Payment Details',
        'Payment details image',
      );
      setShareFeedback(
        shared
          ? 'Image shared successfully.'
          : 'Image downloaded successfully.',
      );
    } catch (error) {
      setShareFeedback('Unable to generate payment image.');
    } finally {
      setShareLoading(null);
    }
  };

  const handleSharePdf = async () => {
    if (!selectedHistoryRow) return;
    setShareLoading('pdf');
    setShareFeedback('');

    try {
      const doc = <HistoryReceiptPdf row={selectedHistoryRow} />;
      const pdfBlob = await pdf(doc).toBlob();
      const fileName = `trovo-payment-${safeFileDate(selectedHistoryRow.dateValue)}.pdf`;
      const file = new File([pdfBlob], fileName, { type: 'application/pdf' });
      const shared = await shareFileWithFallback(
        file,
        'Trovo Payment Details',
        'Payment details PDF',
      );
      setShareFeedback(
        shared ? 'PDF shared successfully.' : 'PDF downloaded successfully.',
      );
    } catch (error) {
      setShareFeedback('Unable to generate payment PDF.');
    } finally {
      setShareLoading(null);
    }
  };

  return (
    <div className={styles.pageShell}>
      <Header />
      <div className={styles.pageBody}>
        <section className={styles.panel}>
          <div className={styles.topRow}>
            <h3 className={styles.headerText}>Transaction History</h3>
          </div>

          <div className={styles.filtersBar}>
            <div className={styles.filterGrid}>
              <HistorySelect
                value={selectedWallet}
                onChange={setSelectedWallet}
                options={walletOptions}
              />
              <HistorySelect
                value={selectedAsset}
                onChange={setSelectedAsset}
                options={assetOptions}
              />
              <HistorySelect
                value={selectedType}
                onChange={setSelectedType}
                options={transactionOptions}
              />
              <FilterTrigger
                label={amountLabel}
                icon={<MoneyIcon />}
                onClick={openAmountModal}
              />
              <FilterTrigger
                label={dateLabel}
                icon={<CalendarIcon />}
                onClick={openDateModal}
              />
              <div className={styles.filterActionCell}>
                <Button
                  label="More Filters"
                  onclick={openMoreFiltersModal}
                  leftIcon={<img src="/images/filterWhite.png" alt="filter" />}
                  additionalClasses={styles.moreFiltersCompact}
                />
              </div>
              <div className={styles.filterActionCell}>
                <ButtonSecondary
                  label="Reset Filters"
                  onclick={resetAllFilters}
                  additionalClasses={styles.resetFiltersButton}
                />
              </div>
            </div>
          </div>

          <div className={styles.tableWrap}>
            <div className={styles.tableHeader}>
              <span>Price</span>
              <span>Transaction type</span>
              <span>Description</span>
              <span>Date</span>
            </div>

            {isLoading ? (
              <div className={styles.emptyState}>
                Loading transaction history...
              </div>
            ) : filteredRows.length === 0 ? (
              <div className={styles.emptyState}>
                No transactions match the selected filters.
              </div>
            ) : (
              <div className={styles.rows}>
                {filteredRows.map((row) => (
                  <button
                    type="button"
                    className={`${styles.row} ${styles.rowButton}`}
                    key={row.id}
                    onClick={() => setSelectedHistoryRow(row)}
                  >
                    <div
                      className={`${styles.price} ${
                        row.type === 'Sent'
                          ? styles.priceNegative
                          : styles.pricePositive
                      }`}
                    >
                      {row.price}
                    </div>
                    <TransactionType tx={row} />
                    <div className={styles.description}>{row.description}</div>
                    <div className={styles.date}>{row.dateLabel}</div>
                  </button>
                ))}
              </div>
            )}
          </div>
        </section>
      </div>

      <OverlayCard
        show={Boolean(selectedHistoryRow)}
        onClose={() => setSelectedHistoryRow(null)}
        cardClassName={styles.detailsOverlayCard}
        hideCloseButton
      >
        {selectedHistoryRow ? (
          <div className={styles.detailsModalBody}>
            <button
              type="button"
              className={styles.detailsCloseButton}
              onClick={() => {
                setSelectedHistoryRow(null);
                setShareFeedback('');
              }}
              aria-label="Close payment details"
            >
              <CloseIcon />
            </button>

            <div className={styles.detailsHeader}>
              <img
                src="/images/trovoHorizontalLogo.png"
                alt="Trovotech"
                className={styles.detailsLogo}
              />
              <h2 className={styles.detailsTitle}>Payment Details</h2>
              <p className={styles.detailsGeneratedText}>
                Generated from Trovo-App on {formatModalDateTime(new Date())}
              </p>
            </div>

            <div className={styles.detailsCard}>
              <DetailRow
                title="Received On"
                name={
                  selectedHistoryRow.toName ||
                  selectedHistoryRow.username ||
                  selectedHistoryRow.walletLabel
                }
                value={shortenKey(selectedHistoryRow.toPublicKey)}
              />
              <DetailRow
                title="From"
                name={selectedHistoryRow.fromName || '-'}
                value={shortenKey(selectedHistoryRow.fromPublicKey)}
              />
              <DetailRow
                title="Blockchain Proof (Transaction ID)"
                name={selectedHistoryRow.transactionId || '-'}
              />
              <DetailRow
                title="Date"
                name={formatModalDateTime(
                  selectedHistoryRow.dateValue,
                  selectedHistoryRow.dateRaw,
                )}
              />
            </div>

            <div className={styles.detailsActions}>
              <button
                type="button"
                className={styles.shareImageButton}
                onClick={handleShareImage}
                disabled={shareLoading !== null}
              >
                {shareLoading === 'image'
                  ? 'Preparing image...'
                  : 'Share Image'}
              </button>
              <button
                type="button"
                className={styles.sharePdfButton}
                onClick={handleSharePdf}
                disabled={shareLoading !== null}
              >
                {shareLoading === 'pdf' ? 'Preparing PDF...' : 'Share PDF'}
              </button>
              <button
                type="button"
                className={styles.shareTextButton}
                onClick={handleShareText}
                disabled={shareLoading !== null}
              >
                {shareLoading === 'text' ? 'Preparing text...' : 'Share Text'}
              </button>
              {shareFeedback ? (
                <p className={styles.detailsShareFeedback}>{shareFeedback}</p>
              ) : null}
            </div>
          </div>
        ) : null}
      </OverlayCard>

      <OverlayCard
        show={showMoreFiltersModal}
        onClose={() => setShowMoreFiltersModal(false)}
      >
        <div className={styles.modalBody}>
          <div className={styles.modalHeaderBlock}>
            <h2 className={styles.modalTitle}>Filter By:</h2>
            <p className={styles.modalSubtitle}>
              Select the options under the filters you wish to apply
            </p>
          </div>

          <div className={styles.formGrid}>
            <LabeledInput
              label="Username or Full Name"
              placeholder="Enter username"
              value={draftFilters.username}
              onChange={(value) =>
                setDraftFilters((current) => ({ ...current, username: value }))
              }
            />
            <LabeledInput
              label="“From” Public Key"
              placeholder="“From” Public Key"
              value={draftFilters.fromPublicKey}
              onChange={(value) =>
                setDraftFilters((current) => ({
                  ...current,
                  fromPublicKey: value,
                }))
              }
            />
            <LabeledInput
              label="“To” Public Key"
              placeholder="“To” Public Key"
              value={draftFilters.toPublicKey}
              onChange={(value) =>
                setDraftFilters((current) => ({
                  ...current,
                  toPublicKey: value,
                }))
              }
            />
            <LabeledInput
              label="Memo Text"
              placeholder="Enter Memo"
              value={draftFilters.memo}
              onChange={(value) =>
                setDraftFilters((current) => ({ ...current, memo: value }))
              }
            />
          </div>

          <div className={styles.modalActions}>
            <button
              type="button"
              className={styles.secondaryAction}
              onClick={() => {
                setDraftFilters(defaultFilters);
                setFilters(defaultFilters);
                setShowMoreFiltersModal(false);
              }}
            >
              Clear All Filters
            </button>
            <button
              type="button"
              className={styles.primaryAction}
              onClick={applyMoreFilters}
            >
              Apply
            </button>
          </div>
        </div>
      </OverlayCard>

      <OverlayCard show={showDateModal} onClose={() => setShowDateModal(false)}>
        <div className={styles.modalBody}>
          <div className={styles.modalHeaderBlock}>
            <h2 className={styles.largeTitle}>Select Date Range</h2>
          </div>

          <div className={styles.dateSection}>
            <p className={styles.sectionTitle}>Enter the Date Range Below</p>
            <div className={styles.presetRow}>
              <PresetButton
                label="Past Week"
                active={draftDateRangeKey === 'week'}
                onClick={() => setDraftDateRangeKey('week')}
              />
              <PresetButton
                label="Past Month"
                active={draftDateRangeKey === 'month'}
                onClick={() => setDraftDateRangeKey('month')}
              />
              <PresetButton
                label="Past 3 Months"
                active={draftDateRangeKey === 'quarter'}
                onClick={() => setDraftDateRangeKey('quarter')}
              />
            </div>

            <div className={styles.dateInputs}>
              <DateInput
                label="Start Date"
                placeholder="From"
                value={draftStartDate}
                onChange={(value) => {
                  setDraftDateRangeKey('custom');
                  setDraftStartDate(value);
                }}
              />
              <DateInput
                label="End Date"
                placeholder="To"
                value={draftEndDate}
                onChange={(value) => {
                  setDraftDateRangeKey('custom');
                  setDraftEndDate(value);
                }}
              />
            </div>
          </div>

          <div className={styles.singleAction}>
            <button
              type="button"
              className={styles.primaryActionWide}
              onClick={applyDateFilter}
            >
              Apply
            </button>
          </div>
        </div>
      </OverlayCard>

      <OverlayCard
        show={showAmountModal}
        onClose={() => setShowAmountModal(false)}
      >
        <div className={styles.modalBody}>
          <div className={styles.modalHeaderBlock}>
            <h2 className={styles.largeTitle}>Enter the Amount Range Below</h2>
          </div>

          <div className={styles.amountInputs}>
            <LabeledInput
              label="Minimum Amount"
              placeholder="Minimum Amount"
              type="number"
              value={minAmountInput}
              onChange={setMinAmountInput}
            />
            <LabeledInput
              label="Maximum Amount"
              placeholder="Maximum Amount"
              type="number"
              value={maxAmountInput}
              onChange={setMaxAmountInput}
            />
          </div>

          <div className={styles.singleAction}>
            <button
              type="button"
              className={styles.primaryActionWide}
              onClick={applyAmountFilter}
            >
              Apply
            </button>
          </div>
        </div>
      </OverlayCard>
    </div>
  );
};

const HistorySelect = ({
  value,
  onChange,
  options,
}: {
  value: string;
  onChange: (value: string) => void;
  options: Array<{ label: string; value: string }>;
}) => {
  return (
    <div className={styles.selectWrap}>
      <select
        className={styles.select}
        value={value}
        onChange={(event) => onChange(event.target.value)}
      >
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
      <span className={styles.selectIcon}>
        <ChevronDownIcon />
      </span>
    </div>
  );
};

const FilterTrigger = ({
  label,
  icon,
  onClick,
}: {
  label: string;
  icon: ReactNode;
  onClick: () => void;
}) => {
  return (
    <button type="button" className={styles.trigger} onClick={onClick}>
      <span className={styles.triggerIcon}>{icon}</span>
      <span className={styles.triggerLabel}>{label}</span>
    </button>
  );
};

const TransactionType = ({ tx }: { tx: { type: string } }) => {
  const config =
    transactionTypeConfig[tx.type as keyof typeof transactionTypeConfig];

  if (!config) {
    return <span className={styles.txTypeFallback}>{tx.type}</span>;
  }

  return (
    <div className={styles.txType}>
      <img className={styles.txIcon} src={config.icon} alt={tx.type} />
      <span className={config.color}>{tx.type}</span>
    </div>
  );
};

const DetailRow = ({
  title,
  name,
  value,
}: {
  title: string;
  name: string;
  value?: string;
}) => {
  return (
    <div className={styles.detailRow}>
      <p className={styles.detailRowTitle}>{title}</p>
      <p className={styles.detailRowName}>{name}</p>
      {value ? <p className={styles.detailRowValue}>{value}</p> : null}
    </div>
  );
};

const HistoryReceiptPdf = ({ row }: { row: HistoryRow }) => {
  return (
    <Document>
      <Page size="A4" style={receiptPdfStyles.page}>
        <View style={receiptPdfStyles.header}>
          <Text style={receiptPdfStyles.brand}>trovotech</Text>
          <Text style={receiptPdfStyles.title}>Payment Details</Text>
          <Text style={receiptPdfStyles.generated}>
            Generated from Trovo-App on {formatModalDateTime(new Date())}
          </Text>
        </View>

        <View style={receiptPdfStyles.card}>
          <PdfDetailRow
            title="Received On"
            lineOne={row.toName || row.username || row.walletLabel}
            lineTwo={shortenKey(row.toPublicKey)}
          />
          <PdfDetailRow
            title="From"
            lineOne={row.fromName || '-'}
            lineTwo={shortenKey(row.fromPublicKey)}
          />
          <PdfDetailRow
            title="Blockchain Proof (Transaction ID)"
            lineOne={row.transactionId || '-'}
          />
          <PdfDetailRow
            title="Date"
            lineOne={formatModalDateTime(row.dateValue, row.dateRaw)}
            isLast
          />
        </View>
      </Page>
    </Document>
  );
};

const PdfDetailRow = ({
  title,
  lineOne,
  lineTwo,
  isLast = false,
}: {
  title: string;
  lineOne: string;
  lineTwo?: string;
  isLast?: boolean;
}) => {
  return (
    <View
      style={
        isLast
          ? [receiptPdfStyles.detailRow, receiptPdfStyles.noBorder]
          : receiptPdfStyles.detailRow
      }
    >
      <Text style={receiptPdfStyles.detailTitle}>{title}</Text>
      <Text style={receiptPdfStyles.detailValue}>{lineOne}</Text>
      {lineTwo ? (
        <Text style={receiptPdfStyles.detailValue}>{lineTwo}</Text>
      ) : null}
    </View>
  );
};

const OverlayCard = ({
  children,
  show,
  onClose,
  cardClassName = '',
  hideCloseButton = false,
}: {
  children: ReactNode;
  show: boolean;
  onClose: () => void;
  cardClassName?: string;
  hideCloseButton?: boolean;
}) => {
  if (!show) return null;

  return (
    <div className={styles.overlay}>
      <div className={styles.overlayBackdrop} onClick={onClose} />
      <div className={`${styles.overlayCard} ${cardClassName}`.trim()}>
        {!hideCloseButton ? (
          <button
            type="button"
            className={styles.closeButton}
            onClick={onClose}
            aria-label="Close modal"
          >
            <CloseIcon />
          </button>
        ) : null}
        {children}
      </div>
    </div>
  );
};

const LabeledInput = ({
  label,
  placeholder,
  value,
  onChange,
  type = 'text',
}: {
  label: string;
  placeholder: string;
  value: string;
  onChange: (value: string) => void;
  type?: string;
}) => {
  return (
    <label className={styles.inputGroup}>
      <span className={styles.inputLabel}>{label}</span>
      <input
        className={styles.input}
        type={type}
        placeholder={placeholder}
        value={value}
        onChange={(event) => onChange(event.target.value)}
      />
    </label>
  );
};

const DateInput = ({
  label,
  placeholder,
  value,
  onChange,
}: {
  label: string;
  placeholder: string;
  value: string;
  onChange: (value: string) => void;
}) => {
  return (
    <label className={styles.inputGroup}>
      <span className={styles.inputLabel}>{label}</span>
      <div className={styles.dateInputWrap}>
        <input
          className={`${styles.input} ${styles.dateInput}`}
          type="date"
          aria-label={label}
          placeholder={placeholder}
          value={value}
          onChange={(event) => onChange(event.target.value)}
        />
        <span className={styles.dateIcon}>
          <CalendarIcon />
        </span>
      </div>
    </label>
  );
};

const PresetButton = ({
  label,
  active,
  onClick,
}: {
  label: string;
  active: boolean;
  onClick: () => void;
}) => {
  return (
    <button
      type="button"
      className={`${styles.presetButton} ${active ? styles.presetButtonActive : ''}`}
      onClick={onClick}
    >
      {label}
    </button>
  );
};

const CalendarIcon = () => (
  <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
    <path
      d="M7 2V5M17 2V5M3 9H21M5 4H19C20.1046 4 21 4.89543 21 6V19C21 20.1046 20.1046 21 19 21H5C3.89543 21 3 20.1046 3 19V6C3 4.89543 3.89543 4 5 4Z"
      stroke="#1B4E91"
      strokeWidth="2.4"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

const MoneyIcon = () => (
  <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
    <rect
      x="3"
      y="6"
      width="18"
      height="12"
      rx="2.5"
      stroke="#1B4E91"
      strokeWidth="2.2"
    />
    <circle cx="12" cy="12" r="2.6" stroke="#1B4E91" strokeWidth="2.2" />
    <path
      d="M6.5 10.5H6.51M17.5 13.5H17.51"
      stroke="#1B4E91"
      strokeWidth="2.6"
      strokeLinecap="round"
    />
  </svg>
);

const ChevronDownIcon = () => (
  <svg width="14" height="8" viewBox="0 0 14 8" fill="none">
    <path
      d="M1 1L7 7L13 1"
      stroke="#163C7A"
      strokeWidth="1.8"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

const CloseIcon = () => (
  <svg width="26" height="26" viewBox="0 0 24 24" fill="none">
    <path
      d="M6 6L18 18M18 6L6 18"
      stroke="#1C3E7A"
      strokeWidth="2.3"
      strokeLinecap="round"
    />
  </svg>
);

const buildDateRange = (
  rangeKey: DateRangeKey,
  start: string,
  end: string,
): { start: string; end: string } => {
  if (rangeKey === 'none') {
    return { start: '', end: '' };
  }

  const now = new Date();
  const fmt = (d: Date) => d.toISOString().slice(0, 10);

  if (rangeKey === 'week') {
    const s = new Date(now);
    s.setDate(now.getDate() - 7);
    return { start: fmt(s), end: fmt(now) };
  }
  if (rangeKey === 'month') {
    const s = new Date(now);
    s.setDate(now.getDate() - 30);
    return { start: fmt(s), end: fmt(now) };
  }
  if (rangeKey === 'quarter') {
    const s = new Date(now);
    s.setDate(now.getDate() - 90);
    return { start: fmt(s), end: fmt(now) };
  }
  return { start, end };
};

const getAmountLabel = (minAmount: number | null, maxAmount: number | null) => {
  if (minAmount === null && maxAmount === null) return 'Amount';
  if (minAmount !== null && maxAmount !== null) {
    return `${minAmount} - ${maxAmount}`;
  }
  if (minAmount !== null) return `From ${minAmount}`;
  return `Up to ${maxAmount}`;
};

const getDateLabel = (
  dateRangeKey: DateRangeKey,
  startDate: string,
  endDate: string,
) => {
  if (dateRangeKey === 'week') return 'Past Week';
  if (dateRangeKey === 'month') return 'Past Month';
  if (dateRangeKey === 'quarter') return 'Past 3 Months';
  if (startDate && endDate) return `${startDate} - ${endDate}`;
  if (startDate) return `From ${startDate}`;
  if (endDate) return `To ${endDate}`;
  return 'Select Date';
};

const matchesDateRange = (
  date: Date | null,
  startDate: string,
  endDate: string,
) => {
  if (!startDate || !endDate) return true;
  if (!date) return false;

  const dateOnly = date.toISOString().slice(0, 10);
  return dateOnly >= startDate && dateOnly <= endDate;
};

const shortenKey = (value: string) => {
  const text = value.trim();
  if (!text) return '-';
  if (text.length <= 12) return text;
  return `${text.slice(0, 5)}...${text.slice(-5)}`;
};

const formatModalDateTime = (date: Date | null, fallbackRaw = ''): string => {
  if (!date) {
    if (!fallbackRaw) return '-';
    const parsed = parseDate(fallbackRaw);
    if (!parsed) return fallbackRaw;
    return formatModalDateTime(parsed);
  }

  return date.toLocaleString('en-US', {
    month: 'long',
    day: 'numeric',
    year: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
    hour12: true,
  });
};

const buildHistoryShareText = (row: HistoryRow) => {
  const receiver = row.toName || row.username || row.walletLabel;
  return [
    'Trovo Payment Details',
    '',
    `Type: ${row.type}`,
    `Amount: ${row.price}`,
    `Asset: ${row.assetCode}`,
    `Received On: ${receiver}`,
    `Receiver Key: ${row.toPublicKey || '-'}`,
    `From: ${row.fromName || '-'}`,
    `Sender Key: ${row.fromPublicKey || '-'}`,
    `Transaction ID: ${row.transactionId || '-'}`,
    `Date: ${formatModalDateTime(row.dateValue, row.dateRaw)}`,
    row.memo ? `Memo: ${row.memo}` : '',
  ]
    .filter(Boolean)
    .join('\n');
};

const safeFileDate = (date: Date | null) => {
  if (!date) return 'transaction';
  return date.toISOString().slice(0, 19).replace(/[:T]/g, '-');
};

const shareFileWithFallback = async (
  file: File,
  title: string,
  text: string,
) => {
  if (
    navigator.share &&
    typeof navigator.canShare === 'function' &&
    navigator.canShare({ files: [file] })
  ) {
    await navigator.share({ title, text, files: [file] });
    return true;
  }

  downloadBlob(file, file.name);
  return false;
};

const downloadBlob = (blob: Blob, fileName: string) => {
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement('a');
  anchor.href = url;
  anchor.download = fileName;
  document.body.appendChild(anchor);
  anchor.click();
  anchor.remove();
  URL.revokeObjectURL(url);
};

const generateHistoryImageBlob = async (row: HistoryRow) => {
  const canvas = document.createElement('canvas');
  const width = 1080;
  const height = 1400;
  canvas.width = width;
  canvas.height = height;
  const ctx = canvas.getContext('2d');
  if (!ctx) throw new Error('Canvas context unavailable');

  ctx.fillStyle = '#FFFFFF';
  ctx.fillRect(0, 0, width, height);

  ctx.fillStyle = '#1B2E5A';
  ctx.font = '700 72px Arial';
  ctx.fillText('trovotech', 110, 140);
  ctx.font = '700 54px Arial';
  ctx.fillText('Payment Details', 110, 230);

  ctx.fillStyle = '#6A7FA8';
  ctx.font = '400 30px Arial';
  ctx.fillText(`Generated on ${formatModalDateTime(new Date())}`, 110, 285);

  const cardX = 90;
  const cardY = 340;
  const cardWidth = 900;
  const cardHeight = 760;
  ctx.fillStyle = '#EDF1F7';
  roundRect(ctx, cardX, cardY, cardWidth, cardHeight, 28);
  ctx.fill();

  const sections = [
    {
      title: 'Received On',
      lineOne: row.toName || row.username || row.walletLabel,
      lineTwo: shortenKey(row.toPublicKey),
    },
    {
      title: 'From',
      lineOne: row.fromName || '-',
      lineTwo: shortenKey(row.fromPublicKey),
    },
    {
      title: 'Blockchain Proof (Transaction ID)',
      lineOne: row.transactionId || '-',
      lineTwo: '',
    },
    {
      title: 'Date',
      lineOne: formatModalDateTime(row.dateValue, row.dateRaw),
      lineTwo: '',
    },
  ];

  const rowHeight = cardHeight / sections.length;
  sections.forEach((section, index) => {
    const baseY = cardY + rowHeight * index;
    if (index > 0) {
      ctx.strokeStyle = '#DDE4EF';
      ctx.lineWidth = 2;
      ctx.beginPath();
      ctx.moveTo(cardX + 28, baseY);
      ctx.lineTo(cardX + cardWidth - 28, baseY);
      ctx.stroke();
    }

    ctx.fillStyle = '#183E7D';
    ctx.font = '700 36px Arial';
    ctx.fillText(section.title, cardX + 48, baseY + 70);
    ctx.fillStyle = '#49618B';
    ctx.font = '400 34px Arial';
    fillWrappedText(
      ctx,
      section.lineOne,
      cardX + 48,
      baseY + 128,
      cardWidth - 96,
      40,
    );
    if (section.lineTwo) {
      fillWrappedText(
        ctx,
        section.lineTwo,
        cardX + 48,
        baseY + 176,
        cardWidth - 96,
        40,
      );
    }
  });

  return await new Promise<Blob>((resolve, reject) => {
    canvas.toBlob((blob) => {
      if (!blob) {
        reject(new Error('Failed to create image blob'));
        return;
      }
      resolve(blob);
    }, 'image/png');
  });
};

const fillWrappedText = (
  ctx: CanvasRenderingContext2D,
  text: string,
  x: number,
  y: number,
  maxWidth: number,
  lineHeight: number,
) => {
  const words = text.split(' ');
  let line = '';
  let lineOffset = 0;

  words.forEach((word) => {
    const testLine = line ? `${line} ${word}` : word;
    const testWidth = ctx.measureText(testLine).width;
    if (testWidth > maxWidth && line) {
      ctx.fillText(line, x, y + lineOffset);
      line = word;
      lineOffset += lineHeight;
      return;
    }
    line = testLine;
  });

  if (line) {
    ctx.fillText(line, x, y + lineOffset);
  }
};

const roundRect = (
  ctx: CanvasRenderingContext2D,
  x: number,
  y: number,
  width: number,
  height: number,
  radius: number,
) => {
  ctx.beginPath();
  ctx.moveTo(x + radius, y);
  ctx.lineTo(x + width - radius, y);
  ctx.quadraticCurveTo(x + width, y, x + width, y + radius);
  ctx.lineTo(x + width, y + height - radius);
  ctx.quadraticCurveTo(x + width, y + height, x + width - radius, y + height);
  ctx.lineTo(x + radius, y + height);
  ctx.quadraticCurveTo(x, y + height, x, y + height - radius);
  ctx.lineTo(x, y + radius);
  ctx.quadraticCurveTo(x, y, x + radius, y);
  ctx.closePath();
};

const receiptPdfStyles = StyleSheet.create({
  page: {
    paddingTop: 40,
    paddingHorizontal: 36,
    backgroundColor: '#FFFFFF',
  },
  header: {
    textAlign: 'center',
    marginBottom: 20,
  },
  brand: {
    fontSize: 34,
    color: '#1B2E5A',
    fontWeight: 700,
    marginBottom: 8,
  },
  title: {
    fontSize: 24,
    color: '#193E7D',
    fontWeight: 700,
    marginBottom: 6,
  },
  generated: {
    fontSize: 11,
    color: '#6A7FA8',
  },
  card: {
    backgroundColor: '#EDF1F7',
    borderRadius: 12,
    overflow: 'hidden',
  },
  detailRow: {
    borderBottomWidth: 1,
    borderBottomColor: '#DDE4EF',
    paddingHorizontal: 14,
    paddingVertical: 12,
  },
  noBorder: {
    borderBottomWidth: 0,
  },
  detailTitle: {
    fontSize: 12,
    fontWeight: 700,
    color: '#183E7D',
    marginBottom: 4,
  },
  detailValue: {
    fontSize: 11,
    color: '#49618B',
    marginBottom: 3,
  },
});
