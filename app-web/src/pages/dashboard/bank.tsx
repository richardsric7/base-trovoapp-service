import { useMemo, useState } from 'react';
import { useLocation } from 'react-router-dom';
import Header from '../../components/header';
import { P2PEmptyState, P2PListCard, p2p } from '../../components/p2p/P2PTheme';
import { useP2PIdentity } from '../../hooks/useP2PIdentity';
import {
  useDepositNairaMutation,
  useGetBankWithdrawalsQuery,
  useGetStablerailBanksQuery,
  useGetStablerailProfileQuery,
  useWithdrawToBankMutation,
} from '../../store/api/bankApis';
import { BankWithdrawal, VirtualAccount } from '../../types/bank';
import { signBase64Txn } from '../../utils/trovoSDK';
import { showNotification } from '../../utils/showToaster';

type Tab = 'deposit' | 'withdraw' | 'history';

const errorMessage = (e: any, fallback: string) => e?.data?.message || e?.data?.error || e?.message || fallback;

// BankView moves Naira between a Nigerian bank account and the wallet's cNGN
// through Stablerail: Deposit (pay into a virtual account, cNGN arrives in
// the wallet), Withdraw (cNGN leaves the wallet, Naira is paid to the bank
// account) and the withdrawals' history. Opened from the wallet's
// Deposit/Withdraw button; location.state.walletAddress picks the wallet.
// Users are onboarded by the backend when KYC level 1 (BVN) completes in the
// mobile app.
export default function BankView() {
  const location = useLocation();
  const { creds: baseCreds, ready } = useP2PIdentity();
  const walletAddress = (location.state as { walletAddress?: string } | null)?.walletAddress;
  const creds = useMemo(() => ({ ...baseCreds, address: walletAddress || baseCreds.address }), [baseCreds, walletAddress]);
  const [tab, setTab] = useState<Tab>((location.state as { tab?: Tab } | null)?.tab ?? 'deposit');

  const { data: profile, isLoading } = useGetStablerailProfileQuery({ creds }, { skip: !ready });
  const onboarded = !!profile?.enabled && !!profile?.onboarded;
  const { data: banks } = useGetStablerailBanksQuery({ creds }, { skip: !ready || !onboarded });
  const { data: history, refetch } = useGetBankWithdrawalsQuery({ creds }, { skip: !ready || !onboarded });

  let body: React.ReactNode;
  if (!ready || isLoading) {
    body = <div className="text-center py-16 text-gray-400">Loading...</div>;
  } else if (!profile?.enabled) {
    body = <P2PEmptyState message="Bank deposits and withdrawals are not available right now." />;
  } else if (!profile.onboarded) {
    body = (
      <P2PEmptyState
        message={
          profile.onboardingStatus === 'failed'
            ? 'Your BVN could not be verified for bank transactions. Please contact support.'
            : profile.onboardingStatus
              ? 'Your BVN is being verified for bank transactions. We will notify you when it is done.'
              : 'Complete KYC level 1 (BVN verification) in the Trovo mobile app to deposit from and withdraw to your bank account.'
        }
      />
    );
  } else {
    body = (
      <>
        <div className="flex space-x-2 px-2">
          {(['deposit', 'withdraw', 'history'] as Tab[]).map((t) => (
            <button
              key={t}
              onClick={() => setTab(t)}
              className={`px-4 py-2 rounded-xl font-semibold capitalize ${tab === t ? 'text-white' : 'text-primary-800 bg-white'}`}
              style={tab === t ? { backgroundColor: p2p.brandDark } : undefined}
            >
              {t}
            </button>
          ))}
        </div>
        {tab === 'deposit' && <DepositTab creds={creds} />}
        {tab === 'withdraw' && (
          <WithdrawTab
            creds={creds}
            banks={banks ?? []}
            history={history?.withdrawals ?? []}
            minimum={profile.minimumWithdrawal}
            onSent={() => {
              refetch();
              setTab('history');
            }}
          />
        )}
        {tab === 'history' && <HistoryTab withdrawals={history?.withdrawals ?? []} />}
      </>
    );
  }

  return (
    <div className="flex flex-col space-y-5 p-3 max-w-2xl">
      <Header />
      <h1 className="text-xl font-bold text-primary-800 px-2">Bank (NGN)</h1>
      {body}
    </div>
  );
}

function DepositTab({ creds }: { creds: any }) {
  const [amount, setAmount] = useState('');
  const [account, setAccount] = useState<VirtualAccount | null>(null);
  const [deposit, { isLoading }] = useDepositNairaMutation();

  const submit = async () => {
    if (!(Number(amount) > 0)) {
      showNotification('error', 'Enter the amount in Naira.');
      return;
    }
    try {
      const r = await deposit({ creds, amount }).unwrap();
      setAccount(r.data.virtualAccount);
    } catch (e) {
      showNotification('error', errorMessage(e, 'Could not start the deposit'));
    }
  };

  return (
    <P2PListCard>
      <div className="flex flex-col space-y-4">
        <p className="text-gray-600">
          Pay Naira into a dedicated bank account and receive the same amount of cNGN in this wallet, less the provider&apos;s fee.
        </p>
        <Field label="Amount (NGN)" value={amount} onChange={setAmount} inputMode="decimal" />
        <Submit label="Get account to pay into" onClick={submit} busy={isLoading} />
        {account && (
          <div className="rounded-xl bg-gray-50 p-4 space-y-1">
            <p className="font-bold">Pay exactly this amount into:</p>
            <Row label="Account number" value={account.accountNumber} />
            <Row label="Bank" value={account.bankName} />
            <Row label="Account name" value={account.accountName} />
            <Row label="Amount to pay" value={`NGN ${account.amount}`} />
            <p className="text-xs text-gray-500 pt-2">
              The account is valid for 30 minutes. cNGN arrives in your wallet once the payment is confirmed.
            </p>
          </div>
        )}
      </div>
    </P2PListCard>
  );
}

function WithdrawTab({
  creds,
  banks,
  history,
  minimum,
  onSent,
}: {
  creds: any;
  banks: { bank_code: string; bank_name: string }[];
  history: BankWithdrawal[];
  minimum: string;
  onSent: () => void;
}) {
  const [bankCode, setBankCode] = useState('');
  const [accountNumber, setAccountNumber] = useState('');
  const [amount, setAmount] = useState('');
  const [withdraw, { isLoading }] = useWithdrawToBankMutation();
  const sortedBanks = useMemo(() => [...banks].sort((a, b) => a.bank_name.localeCompare(b.bank_name)), [banks]);
  // the accounts paid before, most recent first, to pick again
  const recent = useMemo(() => {
    const seen = new Map<string, BankWithdrawal>();
    history.forEach((w) => {
      const key = `${w.bankCode}:${w.accountNumber}`;
      if (w.accountNumber && !seen.has(key)) seen.set(key, w);
    });
    return Array.from(seen.values()).slice(0, 5);
  }, [history]);

  const submit = async () => {
    if (!bankCode || accountNumber.length !== 10 || !(Number(amount) > 0)) {
      showNotification('error', 'Choose a bank, enter the 10-digit account number and the amount.');
      return;
    }
    const body = { amount, accountNumber, bankCode };
    try {
      const built = await withdraw({ creds, body }).unwrap();
      if (!built.transaction) throw new Error('Could not prepare the withdrawal');
      if (!window.confirm([...(built.messages ?? []), '', 'Withdraw now?'].join('\n'))) return;
      const transactionSignature = signBase64Txn(creds.secretKey, built.transaction, '');
      await withdraw({ creds, body: { ...body, transaction: built.transaction, transactionSignature } }).unwrap();
      setAmount('');
      showNotification('success', 'Withdrawal on its way. We will notify you when it reaches your bank account.');
      onSent();
    } catch (e) {
      showNotification('error', errorMessage(e, 'Withdrawal failed'));
    }
  };

  return (
    <P2PListCard>
      <div className="flex flex-col space-y-4">
        <p className="text-gray-600">
          Withdraw cNGN from this wallet to a Nigerian bank account. The smallest withdrawal is NGN {minimum}.
        </p>
        {recent.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {recent.map((w) => (
              <button
                key={w.id}
                onClick={() => {
                  setBankCode(banks.some((b) => b.bank_code === w.bankCode) ? w.bankCode : '');
                  setAccountNumber(w.accountNumber);
                }}
                className="px-3 py-1 rounded-full bg-gray-100 text-sm"
              >
                {w.bankName || w.bankCode} {w.accountNumber}
              </button>
            ))}
          </div>
        )}
        <label className="flex flex-col text-sm text-gray-600">
          Bank
          <select value={bankCode} onChange={(e) => setBankCode(e.target.value)} className="mt-1 p-3 rounded-xl border bg-white">
            <option value="">Choose a bank</option>
            {sortedBanks.map((b) => (
              <option key={b.bank_code} value={b.bank_code}>
                {b.bank_name}
              </option>
            ))}
          </select>
        </label>
        <Field label="Account number" value={accountNumber} onChange={(v) => setAccountNumber(v.replace(/\D/g, '').slice(0, 10))} inputMode="numeric" />
        <Field label="Amount (NGN)" value={amount} onChange={setAmount} inputMode="decimal" />
        <Submit label="Withdraw to bank" onClick={submit} busy={isLoading} />
      </div>
    </P2PListCard>
  );
}

// bank withdrawal statuses: ours while the cNGN moves, then Stablerail's
function statusOf(status: string): { label: string; color: string } {
  const s = status.toLowerCase();
  if (['completed', 'complete', 'success', 'successful'].includes(s)) return { label: 'Paid', color: p2p.success };
  if (s.includes('fail') || ['cancelled', 'canceled', 'rejected', 'reversed'].includes(s)) return { label: 'Failed', color: p2p.danger };
  return { label: 'Processing', color: p2p.warning };
}

function HistoryTab({ withdrawals }: { withdrawals: BankWithdrawal[] }) {
  if (!withdrawals.length) return <P2PEmptyState message="No bank withdrawals yet." />;
  return (
    <>
      {withdrawals.map((w) => {
        const st = statusOf(w.status);
        return (
          <P2PListCard key={w.id}>
            <div className="flex items-center justify-between">
              <div>
                <p className="font-bold">NGN {w.amount}</p>
                <p className="text-sm text-gray-500">
                  {w.bankName || w.bankCode} {w.accountNumber}
                </p>
                <p className="text-xs text-gray-400">{new Date(w.createdAt).toLocaleString()}</p>
              </div>
              <span className="px-3 py-1 rounded-full text-xs font-semibold" style={{ color: st.color, backgroundColor: `${st.color}1f` }}>
                {st.label}
              </span>
            </div>
          </P2PListCard>
        );
      })}
    </>
  );
}

function Field({
  label,
  value,
  onChange,
  inputMode,
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
  inputMode?: 'decimal' | 'numeric';
}) {
  return (
    <label className="flex flex-col text-sm text-gray-600">
      {label}
      <input value={value} inputMode={inputMode} onChange={(e) => onChange(e.target.value)} className="mt-1 p-3 rounded-xl border bg-white" />
    </label>
  );
}

function Submit({ label, onClick, busy }: { label: string; onClick: () => void; busy: boolean }) {
  return (
    <button
      disabled={busy}
      onClick={onClick}
      className="px-4 py-3 rounded-xl font-semibold text-white disabled:opacity-50"
      style={{ backgroundColor: p2p.brandDark }}
    >
      {busy ? 'Please wait...' : label}
    </button>
  );
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex justify-between items-center">
      <span className="text-gray-500">{label}</span>
      <button className="font-semibold" title="Copy" onClick={() => navigator.clipboard?.writeText(value.replace(/^NGN /, ''))}>
        {value}
      </button>
    </div>
  );
}
