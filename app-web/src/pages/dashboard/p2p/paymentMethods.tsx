import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Header from '../../../components/header';
import {
  useListMyP2PPaymentMethodsQuery,
  useCreateP2PPaymentMethodMutation,
  useUpdateP2PPaymentMethodMutation,
  useSetP2PPaymentMethodActiveMutation,
} from '../../../store/api/p2pApis';
import { useP2PIdentity } from '../../../hooks/useP2PIdentity';
import { P2PListCard, P2PEmptyState, p2p } from '../../../components/p2p/P2PTheme';
import P2PMerchantGate from '../../../components/p2p/P2PMerchantGate';
import { P2PMerchantPaymentMethod } from '../../../types/p2p';

// P2PPaymentMethods manages a merchant's own saved fiat settlement
// channels: create, edit details (propagates to any live offer using it),
// and activate/deactivate. Never a delete - a payment method still
// referenced by a live offer cannot be deactivated either, the backend
// enforces that and this page just surfaces the resulting error.
export default function P2PPaymentMethods() {
  const navigate = useNavigate();
  const { creds, ready } = useP2PIdentity();
  const { data, isLoading, refetch } = useListMyP2PPaymentMethodsQuery({ creds }, { skip: !ready });
  const [createPaymentMethod, { isLoading: creating }] = useCreateP2PPaymentMethodMutation();
  const [updatePaymentMethod] = useUpdateP2PPaymentMethodMutation();
  const [setActive] = useSetP2PPaymentMethodActiveMutation();

  const [showNew, setShowNew] = useState(false);
  const [newChannel, setNewChannel] = useState('');
  const [newProvider, setNewProvider] = useState('');
  const [newAccount, setNewAccount] = useState('');
  const [error, setError] = useState('');

  const [editingId, setEditingId] = useState('');
  const [editChannel, setEditChannel] = useState('');
  const [editProvider, setEditProvider] = useState('');
  const [editAccount, setEditAccount] = useState('');

  const startEdit = (pm: P2PMerchantPaymentMethod) => {
    setEditingId(pm.id);
    setEditChannel(pm.paymentChannel);
    setEditProvider(pm.provider);
    setEditAccount(pm.account);
  };

  const saveEdit = async () => {
    setError('');
    try {
      await updatePaymentMethod({
        creds,
        paymentMethodId: editingId,
        body: { paymentChannel: editChannel, provider: editProvider, account: editAccount },
      }).unwrap();
      setEditingId('');
    } catch (e: any) {
      setError(e?.data?.message || 'Could not save changes.');
    }
  };

  const createNew = async () => {
    setError('');
    try {
      await createPaymentMethod({
        creds,
        body: { paymentChannel: newChannel, provider: newProvider, account: newAccount },
      }).unwrap();
      setShowNew(false);
      setNewChannel('');
      setNewProvider('');
      setNewAccount('');
    } catch (e: any) {
      setError(e?.data?.message || 'Could not save this payment method.');
    }
  };

  const toggleActive = async (pm: P2PMerchantPaymentMethod) => {
    setError('');
    try {
      await setActive({ creds, paymentMethodId: pm.id, active: !pm.isActive }).unwrap();
      refetch();
    } catch (e: any) {
      setError(e?.data?.message || 'This payment method is in use by a live offer.');
    }
  };

  return (
    <div className="flex flex-col space-y-5 p-3 max-w-2xl">
      <Header />
      <div className="flex items-center justify-between px-2">
        <h1 className="text-xl font-bold text-primary-800">Payment methods</h1>
        <button
          onClick={() => navigate('/dashboard/p2p/my-offers')}
          className="text-sm font-semibold text-primary-800 underline"
        >
          Back to my offers
        </button>
      </div>

      <P2PMerchantGate>
        {error && <div className="text-red-600 px-2">{error}</div>}

        {!showNew ? (
          <button
            onClick={() => setShowNew(true)}
            className="text-sm font-semibold text-white px-4 py-3 rounded-xl w-fit"
            style={{ backgroundColor: p2p.brandDark }}
          >
            + Add new payment method
          </button>
        ) : (
          <div className="p-3 rounded-xl bg-white space-y-2">
            <input
              className="w-full p-3 rounded-xl bg-gray-50 outline-none"
              placeholder="Payment channel (e.g. BANK_TRANSFER)"
              value={newChannel}
              onChange={(e) => setNewChannel(e.target.value)}
            />
            <input
              className="w-full p-3 rounded-xl bg-gray-50 outline-none"
              placeholder="Provider (e.g. GTBANK)"
              value={newProvider}
              onChange={(e) => setNewProvider(e.target.value)}
            />
            <input
              className="w-full p-3 rounded-xl bg-gray-50 outline-none"
              placeholder="Account details"
              value={newAccount}
              onChange={(e) => setNewAccount(e.target.value)}
            />
            <div className="flex gap-2">
              <button
                disabled={creating || !newChannel || !newAccount}
                onClick={createNew}
                className="flex-1 py-2 rounded-xl font-semibold text-white disabled:opacity-50"
                style={{ backgroundColor: p2p.brandDark }}
              >
                Save
              </button>
              <button onClick={() => setShowNew(false)} className="px-4 py-2 rounded-xl font-semibold text-gray-600">
                Cancel
              </button>
            </div>
          </div>
        )}

        {isLoading ? (
          <div className="text-center py-16 text-gray-400">Loading...</div>
        ) : !data?.data?.length ? (
          <P2PEmptyState message="You have no saved payment methods yet." />
        ) : (
          data.data.map((pm) => (
            <P2PListCard key={pm.id}>
              {editingId === pm.id ? (
                <div className="space-y-2">
                  <input
                    className="w-full p-3 rounded-xl bg-gray-50 outline-none"
                    value={editChannel}
                    onChange={(e) => setEditChannel(e.target.value)}
                  />
                  <input
                    className="w-full p-3 rounded-xl bg-gray-50 outline-none"
                    value={editProvider}
                    onChange={(e) => setEditProvider(e.target.value)}
                  />
                  <input
                    className="w-full p-3 rounded-xl bg-gray-50 outline-none"
                    value={editAccount}
                    onChange={(e) => setEditAccount(e.target.value)}
                  />
                  <div className="flex gap-2">
                    <button
                      onClick={saveEdit}
                      className="flex-1 py-2 rounded-xl font-semibold text-white"
                      style={{ backgroundColor: p2p.brandDark }}
                    >
                      Save
                    </button>
                    <button onClick={() => setEditingId('')} className="px-4 py-2 rounded-xl font-semibold text-gray-600">
                      Cancel
                    </button>
                  </div>
                </div>
              ) : (
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-bold">
                      {pm.paymentChannel} - {pm.provider}
                    </p>
                    <p className="text-gray-500">{pm.account}</p>
                    <p className="text-xs" style={{ color: pm.isActive ? p2p.success : p2p.danger }}>
                      {pm.isActive ? 'Active' : 'Disabled'}
                    </p>
                    <div className="flex gap-3 mt-2">
                      <button onClick={() => startEdit(pm)} className="text-sm font-semibold underline" style={{ color: p2p.brandDark }}>
                        Edit
                      </button>
                      <button onClick={() => toggleActive(pm)} className="text-sm font-semibold underline text-gray-600">
                        {pm.isActive ? 'Disable' : 'Enable'}
                      </button>
                    </div>
                  </div>
                </div>
              )}
            </P2PListCard>
          ))
        )}
      </P2PMerchantGate>
    </div>
  );
}
