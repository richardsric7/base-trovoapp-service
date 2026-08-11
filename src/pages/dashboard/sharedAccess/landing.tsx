import { useEffect, useMemo, useState } from 'react';
import { useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import Header from '../../../components/header';
import Button from '../../../components/button';
import { RootState } from '../../../store/reduxStore';
import { Wallet } from '../../../types/wallet';
import SharedAccessCard from './components/SharedAccessCard';
import { useLazyGetSharedAccessWalletsQuery } from '../../../store/api/sharedAccessApis';
import { showNotification } from '../../../utils/showToaster';
import { useDispatch } from 'react-redux';
import { setActiveWallet } from '../../../store/appStateSlice';

export default function SharedAccessLanding() {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const [fetchWallets] = useLazyGetSharedAccessWalletsQuery();
  const [wallets, setWallets] = useState<Wallet[]>([]);
  const [loading, setLoading] = useState(false);

  const sharedWallets = useMemo(() => {
    return wallets.filter(
      (wallet) => wallet.sharedAccessEnabled || wallet.permission,
    );
  }, [wallets]);

  useEffect(() => {
    const load = async () => {
      setLoading(true);
      const payload = {
        signer: appUser.primarySigner,
        publicKey: appUser.publicKey,
        secretKey: appUser.secretKeys[0],
        body: { userId: appUser.username },
      };
      try {
        const result = await fetchWallets(payload);
        if ('data' in result && result.data) {
          const data = result.data as any;
          const nextWallets = [
            ...(data.wallets ?? []),
            ...(appUser.userWallets ?? []),
          ]
            .filter(
              (wallet: Wallet) =>
                wallet.sharedAccessEnabled || wallet.permission,
            )
            .map((wallet: Wallet) => ({ ...wallet }));
          setWallets(nextWallets);
        }
      } catch (error) {
        showNotification('error', 'Unable to load shared access wallets');
      } finally {
        setLoading(false);
      }
    };

    load();
  }, [appUser, fetchWallets]);

  return (
    <div className="flex flex-col gap-5 p-4 md:p-6">
      <Header isHomeView={false} />
      <div className="rounded-3xl bg-primary-50 p-6 shadow-sm">
        <h2 className="text-2xl font-semibold text-primary-800">
          Shared Access
        </h2>
        <p className="mt-2 text-sm text-primary-700">
          Manage shared wallets, review pending approvals, and update
          permissions for your wallet network.
        </p>
      </div>
      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
        <div className="rounded-3xl border border-gray-200 bg-white p-5">
          <h3 className="text-lg font-semibold text-gray-700">
            Manage wallets
          </h3>
          <p className="mt-2 text-sm text-gray-500">
            Review the shared wallets available to you and open their details.
          </p>
          <div className="mt-4">
            <Button
              label="View shared wallets"
              onclick={() => navigate('/dashboard/shared-access/wallets')}
              additionalClasses="h-11"
            />
          </div>
        </div>
        <div className="rounded-3xl border border-gray-200 bg-white p-5">
          <h3 className="text-lg font-semibold text-gray-700">
            Pending approvals
          </h3>
          <p className="mt-2 text-sm text-gray-500">
            Review requests that need your approval before shared-access changes
            are committed.
          </p>
          <div className="mt-4">
            <Button
              label="Open approvals"
              onclick={() => navigate('/dashboard/shared-access/approvals')}
              additionalClasses="h-11"
            />
          </div>
        </div>
        <div className="rounded-3xl border border-gray-200 bg-white p-5">
          <h3 className="text-lg font-semibold text-gray-700">Grant access</h3>
          <p className="mt-2 text-sm text-gray-500">
            Create or update permission groups for viewers, approvers, and
            initiators.
          </p>
          <div className="mt-4">
            <Button
              label="Add shared access"
              onclick={() => navigate('/dashboard/shared-access/add')}
              additionalClasses="h-11"
            />
          </div>
        </div>
      </div>
      <div className="rounded-3xl border border-gray-200 bg-white p-5">
        <div className="flex items-center justify-between gap-3">
          <h3 className="text-lg font-semibold text-gray-700">
            Available shared wallets
          </h3>
          <span className="text-sm text-gray-500">
            {sharedWallets.length} found
          </span>
        </div>
        {loading ? (
          <div className="mt-4 text-sm text-gray-500">
            Loading shared wallets...
          </div>
        ) : sharedWallets.length === 0 ? (
          <div className="mt-4 text-sm text-gray-500">
            No shared wallets are available yet.
          </div>
        ) : (
          <div className="mt-4 flex flex-col gap-3">
            {sharedWallets.map((wallet, index) => (
              <SharedAccessCard
                key={`${wallet.publicKey}-${index}`}
                wallet={wallet}
                onClick={() => {
                  dispatch(setActiveWallet(wallet));
                  navigate(
                    `/dashboard/shared-access/wallets/${wallet.publicKey}`,
                  );
                }}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
