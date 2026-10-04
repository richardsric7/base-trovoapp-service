import { useEffect, useState } from 'react';
import Header from '../../../components/header';
import Button from '../../../components/button';
import { useP2PIdentity } from '../../../hooks/useP2PIdentity';
import {
  useGetGasFeeAssetsQuery,
  useSetGasFeeAssetMutation,
} from '../../../store/api/settingsApis';
import { showNotification } from '../../../utils/showToaster';

// Settings: the asset the user's wallets pay network fees in. Every send
// is paid for by the wallet itself - in the chosen stablecoin when the
// wallet holds enough of it, otherwise in ETH.
export default function Settings() {
  const { creds, ready } = useP2PIdentity();
  const { data, isLoading, refetch } = useGetGasFeeAssetsQuery(
    { creds },
    { skip: !ready },
  );
  const [setGasFeeAsset, { isLoading: saving }] = useSetGasFeeAssetMutation();
  const [selected, setSelected] = useState('');

  useEffect(() => {
    setSelected(data?.gasFeeAsset ?? '');
  }, [data?.gasFeeAsset]);

  const save = async () => {
    try {
      await setGasFeeAsset({ creds, assetCode: selected }).unwrap();
      showNotification('success', 'Network fee preference saved.');
      refetch();
    } catch (e: any) {
      showNotification(
        'error',
        e?.data?.message ?? 'Could not save your preference. Please try again.',
      );
    }
  };

  const options = [
    { value: '', label: 'ETH (default)' },
    ...(data?.assets ?? []).map((a) => ({
      value: a.assetCode,
      label: a.assetName ? `${a.assetCode} - ${a.assetName}` : a.assetCode,
    })),
  ];

  return (
    <div className="flex text-primary-800 text-sm md:text-md flex-col space-y-5 p-5">
      <Header isHomeView />
      <div className="flex flex-col p-5 rounded-lg bg-primary-100 max-w-3xl space-y-4">
        <p className="text-lg font-montserratSemiBold">Network fees</p>
        <p>
          Your wallets pay their own network fees. Choose a stablecoin to pay
          them in; a wallet that does not hold enough of it pays in ETH
          instead. The fee is shown before you confirm each transaction.
        </p>
        {isLoading || !ready ? (
          <p>Loading...</p>
        ) : (
          <div className="flex flex-col space-y-2">
            {options.map((o) => (
              <label
                key={o.value || 'eth'}
                className="flex items-center space-x-2 cursor-pointer"
              >
                <input
                  type="radio"
                  name="gasFeeAsset"
                  className="w-4 h-4 accent-primary-800 cursor-pointer"
                  checked={selected === o.value}
                  onChange={() => setSelected(o.value)}
                />
                <span>{o.label}</span>
              </label>
            ))}
          </div>
        )}
        <div className="w-full md:w-1/3">
          <Button
            label={saving ? 'Saving...' : 'Save'}
            onclick={save}
            disabled={saving || !ready || selected === (data?.gasFeeAsset ?? '')}
          />
        </div>
      </div>
    </div>
  );
}
