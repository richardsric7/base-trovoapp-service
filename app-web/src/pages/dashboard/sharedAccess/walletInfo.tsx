import { useEffect, useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useSelector } from 'react-redux';
import Header from '../../../components/header';
import Button from '../../../components/button';
import ButtonSecondary from '../../../components/buttonSecondary';
import PermissionChip from './components/PermissionChip';
import { RootState } from '../../../store/reduxStore';
import { Wallet } from '../../../types/wallet';
import { useLazyGetWalletBalanceQuery } from '../../../store/api/sharedAccessApis';
import { showNotification } from '../../../utils/showToaster';

export default function SharedAccessWalletInfo() {
  const navigate = useNavigate();
  const params = useParams();
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const [fetchBalance] = useLazyGetWalletBalanceQuery();
  const [wallet, setWallet] = useState<Wallet | null>(null);
  const [balance, setBalance] = useState<any>(null);
  const [loading, setLoading] = useState(false);

  const currentWallet = useMemo(() => {
    const walletFromState = appUser.userWallets.find(
      (entry) => entry.address === params.address,
    );
    return walletFromState ?? wallet;
  }, [appUser.userWallets, params.address, wallet]);

  useEffect(() => {
    if (!params.address) return;
    const foundWallet = appUser.userWallets.find(
      (entry) => entry.address === params.address,
    );
    if (foundWallet) {
      setWallet(foundWallet);
    }

    const load = async () => {
      setLoading(true);
      try {
        const result = await fetchBalance({
          signer: appUser.primarySigner,
          address: appUser.address,
          secretKey: appUser.secretKeys[0],
          body: {},
        });
        if ('data' in result && result.data) {
          setBalance(result.data);
        }
      } catch (error) {
        showNotification('error', 'Unable to load wallet balance');
      } finally {
        setLoading(false);
      }
    };

    load();
  }, [appUser, fetchBalance, params.address]);

  if (!currentWallet) {
    return <div className="p-6">Loading shared wallet details...</div>;
  }

  return (
    <div className="flex flex-col gap-5 p-4 md:p-6">
      <Header isHomeView={false} />
      <div className="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm">
        <div className="flex flex-wrap items-start justify-between gap-3">
          <div>
            <p className="text-sm font-semibold uppercase tracking-wide text-primary-600">
              Shared wallet
            </p>
            <h2 className="text-2xl font-semibold text-primary-800">
              {currentWallet.alias}
            </h2>
            <p className="mt-2 text-sm text-gray-500">
              {currentWallet.description || 'No description set.'}
            </p>
          </div>
          <div className="rounded-full bg-primary-100 px-3 py-1 text-sm font-semibold text-primary-700">
            {currentWallet.permission || 'Shared access'}
          </div>
        </div>
        <div className="mt-6 grid gap-4 md:grid-cols-2">
          <div className="rounded-2xl bg-gray-50 p-4">
            <p className="text-sm font-semibold text-gray-600">Owner</p>
            <p className="mt-1 text-base text-gray-800">
              {currentWallet.owner || 'You'}
            </p>
          </div>
          <div className="rounded-2xl bg-gray-50 p-4">
            <p className="text-sm font-semibold text-gray-600">
              Approvals needed
            </p>
            <p className="mt-1 text-base text-gray-800">
              {currentWallet.numberOfApprovalsNeeded ?? 'Not configured'}
            </p>
          </div>
        </div>
        <div className="mt-6">
          <p className="text-sm font-semibold text-gray-600">Permissions</p>
          <div className="mt-3 flex flex-wrap gap-2">
            {(currentWallet.permissions ?? []).length === 0 ? (
              <span className="text-sm text-gray-500">
                No permission metadata is available yet.
              </span>
            ) : (
              currentWallet.permissions.map((permission) => (
                <PermissionChip
                  key={`${permission.targetUsername}-${permission.permission}`}
                  label={`${permission.targetUsername}: ${permission.permission}`}
                />
              ))
            )}
          </div>
        </div>
        <div className="mt-6">
          <p className="text-sm font-semibold text-gray-600">
            Balance overview
          </p>
          {loading ? (
            <div className="mt-2 text-sm text-gray-500">Loading balance...</div>
          ) : (
            <pre className="mt-2 whitespace-pre-wrap break-all rounded-2xl bg-primary-50 p-4 text-sm text-primary-800">
              {JSON.stringify(balance, null, 2)}
            </pre>
          )}
        </div>
        <div className="mt-6 flex flex-wrap gap-3">
          <Button
            label="Update access"
            onclick={() => navigate('/dashboard/shared-access/update')}
            additionalClasses="h-11 md:w-auto px-8"
          />
          <ButtonSecondary
            label="Back"
            onclick={() => navigate('/dashboard/shared-access')}
            additionalClasses="h-11 md:w-auto px-8"
          />
        </div>
      </div>
    </div>
  );
}
