import { useMemo, useState } from 'react';
import { useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import Header from '../../../components/header';
import Button from '../../../components/button';
import ButtonSecondary from '../../../components/buttonSecondary';
import TextInput from '../../../components/textInput';
import { RootState } from '../../../store/reduxStore';
import { useAddSharedAccessMutation } from '../../../store/api/sharedAccessApis';
import { showNotification } from '../../../utils/showToaster';
import { parsePermissionInput } from './utils';

export default function SharedAccessAdd() {
  const navigate = useNavigate();
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const [addSharedAccess] = useAddSharedAccessMutation();
  const [walletAlias, setWalletAlias] = useState('');
  const [walletDescription, setWalletDescription] = useState('');
  const [viewers, setViewers] = useState('');
  const [approvers, setApprovers] = useState('');
  const [initiators, setInitiators] = useState('');
  const [approvalsNeeded, setApprovalsNeeded] = useState('2');
  const [saving, setSaving] = useState(false);

  const permissions = useMemo(() => {
    return [
      ...parsePermissionInput(viewers).map((username) => ({
        targetUsername: username,
        permission: 'VIEW-ONLY',
      })),
      ...parsePermissionInput(approvers).map((username) => ({
        targetUsername: username,
        permission: 'APPROVER',
      })),
      ...parsePermissionInput(initiators).map((username) => ({
        targetUsername: username,
        permission: 'INITIATOR',
      })),
    ];
  }, [approvers, initiators, viewers]);

  const handleSubmit = async () => {
    if (!walletAlias || !walletDescription) {
      showNotification(
        'error',
        'Please provide a wallet alias and description',
      );
      return;
    }

    setSaving(true);
    try {
      const payload = {
        signer: appUser.primarySigner,
        publicKey: appUser.publicKey,
        secretKey: appUser.secretKeys[0],
        body: {
          walletAlias,
          walletDescription,
          numberOfApprovalsNeeded: Number(approvalsNeeded) || 2,
          permissions,
        },
      };
      const response = await addSharedAccess(payload);
      if ('data' in response && response.data) {
        showNotification(
          'success',
          'Shared access request submitted successfully',
        );
        navigate('/dashboard/shared-access');
      } else {
        showNotification('error', 'Unable to submit shared access request');
      }
    } catch (error) {
      showNotification('error', 'Unable to submit shared access request');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="flex flex-col gap-5 p-4 md:p-6">
      <Header isHomeView={false} />
      <div className="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm">
        <h2 className="text-2xl font-semibold text-primary-800">
          Grant shared access
        </h2>
        <p className="mt-2 text-sm text-gray-500">
          Create a shared wallet and assign viewers, approvers, and initiators.
        </p>
        <div className="mt-6 grid gap-4 md:grid-cols-2">
          <TextInput
            inputType="text"
            label="Wallet alias"
            placeholder="e.g. Treasury"
            onInputChange={(value) => setWalletAlias(value)}
          />
          <TextInput
            inputType="text"
            label="Wallet description"
            placeholder="Describe the wallet"
            onInputChange={(value) => setWalletDescription(value)}
          />
          <TextInput
            inputType="text"
            label="Viewers"
            placeholder="user1, user2"
            onInputChange={(value) => setViewers(value)}
          />
          <TextInput
            inputType="text"
            label="Approvers"
            placeholder="user1, user2"
            onInputChange={(value) => setApprovers(value)}
          />
          <TextInput
            inputType="text"
            label="Initiators"
            placeholder="user1, user2"
            onInputChange={(value) => setInitiators(value)}
          />
          <TextInput
            inputType="text"
            label="Approvals needed"
            placeholder="2"
            defaultValue="2"
            onInputChange={(value) => setApprovalsNeeded(value)}
          />
        </div>
        <div className="mt-6 rounded-2xl bg-primary-50 p-4 text-sm text-primary-700">
          Permissions to be submitted:{' '}
          {permissions.length > 0
            ? permissions
                .map((entry) => `${entry.targetUsername} (${entry.permission})`)
                .join(', ')
            : 'None yet'}
        </div>
        <div className="mt-6 flex flex-wrap gap-3">
          <Button
            label={saving ? 'Submitting...' : 'Submit'}
            onclick={handleSubmit}
            additionalClasses="h-11 md:w-auto px-8"
          />
          <ButtonSecondary
            label="Cancel"
            onclick={() => navigate('/dashboard/shared-access')}
            additionalClasses="h-11 md:w-auto px-8"
          />
        </div>
      </div>
    </div>
  );
}
