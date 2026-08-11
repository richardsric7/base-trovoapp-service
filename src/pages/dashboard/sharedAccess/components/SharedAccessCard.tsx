import { Wallet } from '../../../../types/wallet';

type Props = {
  wallet: Wallet;
  onClick: () => void;
};

export default function SharedAccessCard({ wallet, onClick }: Props) {
  return (
    <button
      type="button"
      className="w-full rounded-2xl border border-gray-200 bg-white p-4 text-left shadow-sm transition hover:-translate-y-0.5 hover:shadow-md"
      onClick={onClick}
    >
      <div className="flex items-center justify-between gap-3">
        <div>
          <p className="text-lg font-semibold text-primary-800">
            {wallet.alias}
          </p>
          <p className="mt-1 text-sm text-gray-500">
            {wallet.description || 'Shared wallet'}
          </p>
        </div>
        <span className="rounded-full bg-primary-100 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-primary-700">
          {wallet.permission || 'Shared'}
        </span>
      </div>
      <div className="mt-4 flex flex-wrap gap-2 text-sm text-gray-600">
        <span className="rounded-full bg-gray-100 px-2 py-1">
          Owner: {wallet.owner || 'You'}
        </span>
        {wallet.numberOfApprovalsNeeded ? (
          <span className="rounded-full bg-gray-100 px-2 py-1">
            Approvals: {wallet.numberOfApprovalsNeeded}
          </span>
        ) : null}
      </div>
    </button>
  );
}
