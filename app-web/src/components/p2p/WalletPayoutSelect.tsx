import { useSelector } from 'react-redux';
import { RootState } from '../../store/reduxStore';

// WalletPayoutSelect is the BUY-offer wallet dropdown (Plan: "an ordered
// list of their wallets showing the wallet alias only"). Sourced from the
// same userWallets data already loaded into Redux for the wallet
// switcher, rather than a second P2P-scoped endpoint.
export default function WalletPayoutSelect({
  value,
  onChange,
}: {
  value: string;
  onChange: (address: string) => void;
}) {
  const wallets = useSelector((state: RootState) => state.auth.user?.userWallets ?? []);

  return (
    <div>
      <label className="text-sm text-gray-500 block mb-1">Receive asset into</label>
      <select
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="w-full p-3 rounded-xl bg-white outline-none"
      >
        <option value="" disabled>
          Select a wallet
        </option>
        {wallets.map((w) => (
          <option key={w.address} value={w.address}>
            {w.alias}
          </option>
        ))}
      </select>
    </div>
  );
}
