import { ReactNode, useEffect, useState } from "react";
import { useSelector } from "react-redux";
import Header from "../../../components/header";
import Button from "../../../components/button";
import ButtonSecondary from "../../../components/buttonSecondary";
import { RootState } from "../../../store/reduxStore";
import { useFetchFiatPaymentsQuery } from "../../../store/api/walletApis";
import { Encryptor } from "../../../utils/encryptor";
import { getAssetCode } from "../../../utils/utilities";
import styles from "./main.module.css";
import { transactionTypeConfig } from "./data";

type FilterState = {
  username: string;
  fromPublicKey: string;
  toPublicKey: string;
  memo: string;
};

type DateRangeKey = "none" | "week" | "month" | "quarter" | "custom";

type HistoryRow = {
  id: string;
  type: string;
  assetCode: string;
  amountValue: number;
  price: string;
  description: string;
  dateLabel: string;
  dateValue: Date | null;
  walletKey: string;
  walletLabel: string;
  username: string;
  fromPublicKey: string;
  toPublicKey: string;
  memo: string;
};

const defaultFilters: FilterState = {
  username: "",
  fromPublicKey: "",
  toPublicKey: "",
  memo: "",
};

export const History = () => {
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const [secretKey, setSecretKey] = useState("");
  const [selectedWallet, setSelectedWallet] = useState("");
  const [selectedAsset, setSelectedAsset] = useState("");
  const [selectedType, setSelectedType] = useState("");
  const [showAmountModal, setShowAmountModal] = useState(false);
  const [showDateModal, setShowDateModal] = useState(false);
  const [showMoreFiltersModal, setShowMoreFiltersModal] = useState(false);
  const [minAmountInput, setMinAmountInput] = useState("");
  const [maxAmountInput, setMaxAmountInput] = useState("");
  const [appliedMinAmount, setAppliedMinAmount] = useState<number | null>(null);
  const [appliedMaxAmount, setAppliedMaxAmount] = useState<number | null>(null);
  const [dateRangeKey, setDateRangeKey] = useState<DateRangeKey>("none");
  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState("");
  const [draftStartDate, setDraftStartDate] = useState("");
  const [draftEndDate, setDraftEndDate] = useState("");
  const [draftDateRangeKey, setDraftDateRangeKey] =
    useState<DateRangeKey>("none");
  const [filters, setFilters] = useState<FilterState>(defaultFilters);
  const [draftFilters, setDraftFilters] = useState<FilterState>(defaultFilters);

  useEffect(() => {
    const encryptor = new Encryptor();
    encryptor.getSecretKey(appUser).then(setSecretKey);
  }, [appUser]);

  const wallets = appUser?.userWallets ?? [];
  const primaryWallet =
    wallets.find((wallet) => wallet.primaryWallet) ?? wallets?.[0];

  const forPublicKey =
    selectedWallet !== "" ? selectedWallet : (primaryWallet?.publicKey ?? "");

  const apiDateRange = buildDateRange(dateRangeKey, startDate, endDate);

  const transactionTypeParam =
    selectedType === ""
      ? ""
      : selectedType === "Swap"
      ? "swap"
      : "payment";

  const queryParams: string[] = [];
  if (filters.username) queryParams.push(`&name=${encodeURIComponent(filters.username)}`);
  if (filters.memo) queryParams.push(`&memo=${encodeURIComponent(filters.memo)}`);
  if (filters.fromPublicKey) queryParams.push(`&fromPublicKey=${encodeURIComponent(filters.fromPublicKey)}`);
  if (filters.toPublicKey) queryParams.push(`&toPublicKey=${encodeURIComponent(filters.toPublicKey)}`);
  if (appliedMinAmount !== null || appliedMaxAmount !== null) {
    const min = appliedMinAmount !== null ? appliedMinAmount : "";
    const max = appliedMaxAmount !== null ? appliedMaxAmount : "";
    queryParams.push(`&amount=${encodeURIComponent(`${min}%${max}`)}`);
  }
  if (apiDateRange.start && apiDateRange.end) {
    queryParams.push(
      `&dateBetween=${encodeURIComponent(
        `${apiDateRange.start}%${apiDateRange.end}`
      )}`
    );
  }
  queryParams.push(`&transactionType=${transactionTypeParam}`);

  const { data, isLoading } = useFetchFiatPaymentsQuery(
    {
      signer: primaryWallet?.signer ?? "",
      publicKey: forPublicKey,
      secretKey,
      body: { limit: 50, query: queryParams.join("") },
    },
    { skip: !secretKey || !primaryWallet }
  );
console.log({data});

  const walletOptions = [
    { label: "All wallets", value: "" },
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
    { label: "All assets", value: "" },
    ...Array.from(assetSet).map((assetCode) => ({
      label: assetCode,
      value: assetCode,
    })),
  ];

  const transactionOptions = [
    { label: "All transactions", value: "" },
    { label: "Received", value: "Received" },
    { label: "Sent", value: "Sent" },
    { label: "Swap", value: "Swap" },
  ];

  const rawRecords = data?.data?.records ?? data?.records ?? [];

  const rows: HistoryRow[] = rawRecords.map(
    (item: Record<string, any>, index: number) => {
      const type = normalizeTransactionType(
        item.transactionType ?? item.type,
        item,
        forPublicKey
      );
      const amountValue = parseAmountValue(item.amount);
      const assetCode = getTransactionAssetCode(item);
      const amountPrefix =
        type === "Sent" ? "-" : type === "Received" ? "+" : "+";
      const amountText =
        typeof item.amount === "string" && item.amount.trim().length > 0
          ? item.amount
          : formatAmount(amountValue, assetCode);
      const createdAt = parseDate(
        item.transactionDate ?? item.createdAt ?? item.date
      );
      const description =
        item.description ??
        item.narration ??
        item.memo ??
        (type === "Received"
          ? `Received ${assetCode}`
          : type === "Sent"
          ? `Sent ${assetCode}`
          : `Swapped asset to ${assetCode}`);

      const walletMatch =
        wallets.find(
          (wallet) =>
            wallet.publicKey === item.publicKey ||
            wallet.publicKey === item.walletPublicKey ||
            wallet.alias === item.walletAlias
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
        walletKey: walletMatch?.publicKey ?? "unknown",
        walletLabel: walletMatch?.alias ?? "Primary wallet",
        username: `${item.username ?? item.fullName ?? item.name ?? ""}`.trim(),
        fromPublicKey: `${item.fromPublicKey ?? item.senderPublicKey ?? ""}`,
        toPublicKey: `${item.toPublicKey ?? item.receiverPublicKey ?? ""}`,
        memo: `${item.memo ?? item.narration ?? ""}`.trim(),
      };
    }
  );

  const filteredRows = rows.filter((row) => {
    if (selectedAsset !== "" && row.assetCode !== selectedAsset) {
      return false;
    }
    if (selectedType === "Sent" && row.type !== "Sent") return false;
    if (selectedType === "Received" && row.type !== "Received") return false;
    if (
      !matchesDateRange(
        row.dateValue,
        apiDateRange.start,
        apiDateRange.end
      )
    ) {
      return false;
    }
    return true;
  });

  const amountLabel = getAmountLabel(appliedMinAmount, appliedMaxAmount);
  const dateLabel = getDateLabel(dateRangeKey, startDate, endDate);

  const openAmountModal = () => {
    setMinAmountInput(appliedMinAmount?.toString() ?? "");
    setMaxAmountInput(appliedMaxAmount?.toString() ?? "");
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
      minAmountInput.trim().length > 0 ? Number(minAmountInput) : null
    );
    setAppliedMaxAmount(
      maxAmountInput.trim().length > 0 ? Number(maxAmountInput) : null
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
    setSelectedWallet("");
    setSelectedAsset("");
    setSelectedType("");
    setMinAmountInput("");
    setMaxAmountInput("");
    setAppliedMinAmount(null);
    setAppliedMaxAmount(null);
    setDateRangeKey("none");
    setStartDate("");
    setEndDate("");
    setDraftDateRangeKey("none");
    setDraftStartDate("");
    setDraftEndDate("");
    setFilters(defaultFilters);
    setDraftFilters(defaultFilters);
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
              <div className={styles.emptyState}>Loading transaction history...</div>
            ) : filteredRows.length === 0 ? (
              <div className={styles.emptyState}>
                No transactions match the selected filters.
              </div>
            ) : (
              <div className={styles.rows}>
                {filteredRows.map((row) => (
                  <article className={styles.row} key={row.id}>
                    <div
                      className={`${styles.price} ${
                        row.type === "Sent"
                          ? styles.priceNegative
                          : styles.pricePositive
                      }`}
                    >
                      {row.price}
                    </div>
                    <TransactionType tx={row} />
                    <div className={styles.description}>{row.description}</div>
                    <div className={styles.date}>{row.dateLabel}</div>
                  </article>
                ))}
              </div>
            )}
          </div>
        </section>
      </div>

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
                setDraftFilters((current) => ({ ...current, toPublicKey: value }))
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
                active={draftDateRangeKey === "week"}
                onClick={() => setDraftDateRangeKey("week")}
              />
              <PresetButton
                label="Past Month"
                active={draftDateRangeKey === "month"}
                onClick={() => setDraftDateRangeKey("month")}
              />
              <PresetButton
                label="Past 3 Months"
                active={draftDateRangeKey === "quarter"}
                onClick={() => setDraftDateRangeKey("quarter")}
              />
            </div>

            <div className={styles.dateInputs}>
              <DateInput
                label="Start Date"
                placeholder="From"
                value={draftStartDate}
                onChange={(value) => {
                  setDraftDateRangeKey("custom");
                  setDraftStartDate(value);
                }}
              />
              <DateInput
                label="End Date"
                placeholder="To"
                value={draftEndDate}
                onChange={(value) => {
                  setDraftDateRangeKey("custom");
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

const OverlayCard = ({
  children,
  show,
  onClose,
}: {
  children: ReactNode;
  show: boolean;
  onClose: () => void;
}) => {
  if (!show) return null;

  return (
    <div className={styles.overlay}>
      <div className={styles.overlayBackdrop} onClick={onClose} />
      <div className={styles.overlayCard}>
        <button
          type="button"
          className={styles.closeButton}
          onClick={onClose}
          aria-label="Close modal"
        >
          <CloseIcon />
        </button>
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
  type = "text",
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
      className={`${styles.presetButton} ${active ? styles.presetButtonActive : ""}`}
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
    <rect x="3" y="6" width="18" height="12" rx="2.5" stroke="#1B4E91" strokeWidth="2.2" />
    <circle cx="12" cy="12" r="2.6" stroke="#1B4E91" strokeWidth="2.2" />
    <path d="M6.5 10.5H6.51M17.5 13.5H17.51" stroke="#1B4E91" strokeWidth="2.6" strokeLinecap="round" />
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
  end: string
): { start: string; end: string } => {
  if (rangeKey === "none") {
    return { start: "", end: "" };
  }

  const now = new Date();
  const fmt = (d: Date) => d.toISOString().slice(0, 10);

  if (rangeKey === "week") {
    const s = new Date(now);
    s.setDate(now.getDate() - 7);
    return { start: fmt(s), end: fmt(now) };
  }
  if (rangeKey === "month") {
    const s = new Date(now);
    s.setDate(now.getDate() - 30);
    return { start: fmt(s), end: fmt(now) };
  }
  if (rangeKey === "quarter") {
    const s = new Date(now);
    s.setDate(now.getDate() - 90);
    return { start: fmt(s), end: fmt(now) };
  }
  return { start, end };
};

const normalizeTransactionType = (
  value: string | undefined,
  item?: Record<string, any>,
  activePublicKey?: string
) => {
  const type = `${value ?? ""}`.toLowerCase();

  if (type.includes("swap")) return "Swap";
  if (type.includes("sent") || type.includes("send")) return "Sent";
  if (type.includes("received") || type.includes("receive")) return "Received";
  if (type.includes("payment")) {
    const fromPublicKey = `${item?.fromPublicKey ?? item?.senderPublicKey ?? ""}`;
    const toPublicKey = `${item?.toPublicKey ?? item?.receiverPublicKey ?? ""}`;

    if (activePublicKey && fromPublicKey === activePublicKey) return "Sent";
    if (activePublicKey && toPublicKey === activePublicKey) return "Received";
  }

  return "Received";
};

const parseAmountValue = (amount: unknown) => {
  if (typeof amount === "number") return amount;

  if (typeof amount === "string") {
    const cleaned = amount.replace(/[^0-9.-]/g, "");
    const parsed = Number(cleaned);
    return Number.isNaN(parsed) ? 0 : Math.abs(parsed);
  }

  return 0;
};

const formatAmount = (amount: number, assetCode: string) => {
  return `${amount.toLocaleString("en-NG", {
    minimumFractionDigits: 4,
    maximumFractionDigits: 4,
  })} ${assetCode}`;
};

const parseDate = (value: unknown) => {
  if (!value) return null;
  const parsed = new Date(String(value));
  return Number.isNaN(parsed.getTime()) ? null : parsed;
};

const formatRelativeTime = (date: Date | null) => {
  if (!date) return "-";

  const now = new Date();
  const diffInMs = now.getTime() - date.getTime();
  const diffInDays = Math.floor(diffInMs / (1000 * 60 * 60 * 24));

  if (diffInDays <= 0) return "Today";
  if (diffInDays === 1) return "1 day ago";

  return `${diffInDays} days ago`;
};

const getTransactionAssetCode = (item: Record<string, any>) => {
  if (item.assetCode !== undefined) {
    return getAssetCode(item.assetCode);
  }

  if (item.asset?.assetCode !== undefined) {
    return getAssetCode(item.asset.assetCode);
  }

  if (typeof item.amount === "string") {
    const segments = item.amount.trim().split(/\s+/);
    const lastSegment = segments[segments.length - 1];

    if (/[A-Za-z]/.test(lastSegment)) {
      return lastSegment.toUpperCase();
    }
  }

  return "XBN";
};

const getAmountLabel = (minAmount: number | null, maxAmount: number | null) => {
  if (minAmount === null && maxAmount === null) return "Amount";
  if (minAmount !== null && maxAmount !== null) {
    return `${minAmount} - ${maxAmount}`;
  }
  if (minAmount !== null) return `From ${minAmount}`;
  return `Up to ${maxAmount}`;
};

const getDateLabel = (
  dateRangeKey: DateRangeKey,
  startDate: string,
  endDate: string
) => {
  if (dateRangeKey === "week") return "Past Week";
  if (dateRangeKey === "month") return "Past Month";
  if (dateRangeKey === "quarter") return "Past 3 Months";
  if (startDate && endDate) return `${startDate} - ${endDate}`;
  if (startDate) return `From ${startDate}`;
  if (endDate) return `To ${endDate}`;
  return "Select Date";
};

const matchesDateRange = (
  date: Date | null,
  startDate: string,
  endDate: string
) => {
  if (!startDate || !endDate) return true;
  if (!date) return false;

  const dateOnly = date.toISOString().slice(0, 10);
  return dateOnly >= startDate && dateOnly <= endDate;
};
