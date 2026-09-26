import { useState } from 'react';
import {
  useListMyP2PPaymentMethodsQuery,
  useCreateP2PPaymentMethodMutation,
  P2PCreds,
} from '../../store/api/p2pApis';
import { p2p } from './P2PTheme';

const ADD_NEW = '__add_new__';

// PaymentMethodSelect is the SELL-offer payment-method dropdown: "Add new
// payment method" is always the first item, and picking it opens an
// inline quick-create form - on save, the new payment method is selected
// automatically and the merchant continues filling out the rest of the
// offer without leaving the page.
export default function PaymentMethodSelect({
  creds,
  ready,
  value,
  onChange,
}: {
  creds: P2PCreds;
  ready: boolean;
  value: string;
  onChange: (paymentMethodId: string) => void;
}) {
  const { data } = useListMyP2PPaymentMethodsQuery({ creds }, { skip: !ready });
  const [createPaymentMethod, { isLoading: creating }] = useCreateP2PPaymentMethodMutation();
  const [addingNew, setAddingNew] = useState(false);
  const [paymentChannel, setPaymentChannel] = useState('');
  const [provider, setProvider] = useState('');
  const [account, setAccount] = useState('');
  const [error, setError] = useState('');

  const activeMethods = (data?.data ?? []).filter((m) => m.isActive);

  const handleSelect = (v: string) => {
    if (v === ADD_NEW) {
      setAddingNew(true);
      return;
    }
    onChange(v);
  };

  const saveNew = async () => {
    setError('');
    try {
      const pm = await createPaymentMethod({ creds, body: { paymentChannel, provider, account } }).unwrap();
      setAddingNew(false);
      setPaymentChannel('');
      setProvider('');
      setAccount('');
      onChange(pm.id);
    } catch (e: any) {
      setError(e?.data?.message || 'Could not save this payment method.');
    }
  };

  return (
    <div>
      <label className="text-sm text-gray-500 block mb-1">Payment method</label>
      <select
        value={addingNew ? ADD_NEW : value}
        onChange={(e) => handleSelect(e.target.value)}
        className="w-full p-3 rounded-xl bg-white outline-none"
      >
        <option value="" disabled>
          Select a payment method
        </option>
        <option value={ADD_NEW}>+ Add new payment method</option>
        {activeMethods.map((m) => (
          <option key={m.id} value={m.id}>
            {m.paymentChannel} - {m.provider} ({m.account})
          </option>
        ))}
      </select>

      {addingNew && (
        <div className="mt-3 p-3 rounded-xl bg-white space-y-2">
          <input
            className="w-full p-3 rounded-xl bg-gray-50 outline-none"
            placeholder="Payment channel (e.g. BANK_TRANSFER)"
            value={paymentChannel}
            onChange={(e) => setPaymentChannel(e.target.value)}
          />
          <input
            className="w-full p-3 rounded-xl bg-gray-50 outline-none"
            placeholder="Provider (e.g. GTBANK)"
            value={provider}
            onChange={(e) => setProvider(e.target.value)}
          />
          <input
            className="w-full p-3 rounded-xl bg-gray-50 outline-none"
            placeholder="Account details"
            value={account}
            onChange={(e) => setAccount(e.target.value)}
          />
          {error && <div className="text-red-600 text-sm">{error}</div>}
          <div className="flex gap-2">
            <button
              type="button"
              disabled={creating || !paymentChannel || !account}
              onClick={saveNew}
              className="flex-1 py-2 rounded-xl font-semibold text-white disabled:opacity-50"
              style={{ backgroundColor: p2p.brandDark }}
            >
              Save payment method
            </button>
            <button
              type="button"
              onClick={() => setAddingNew(false)}
              className="px-4 py-2 rounded-xl font-semibold text-gray-600"
            >
              Cancel
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
