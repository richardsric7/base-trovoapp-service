import { useEffect, useState } from 'react';
import { useSelector } from 'react-redux';
import { useNavigate, useParams } from 'react-router-dom';
import Header from '../../../components/header';
import Button from '../../../components/button';
import ButtonSecondary from '../../../components/buttonSecondary';
import { RootState } from '../../../store/reduxStore';
import {
  useApproveSharedAccessMutation,
  useRejectSharedAccessMutation,
} from '../../../store/api/sharedAccessApis';
import { showNotification } from '../../../utils/showToaster';

type ApprovalDetail = {
  id: string;
  alias?: string;
  description?: string;
  status?: string;
};

export default function SharedAccessApprovalDetails() {
  const navigate = useNavigate();
  const params = useParams();
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const [approve] = useApproveSharedAccessMutation();
  const [reject] = useRejectSharedAccessMutation();
  const [approval, setApproval] = useState<ApprovalDetail | null>(null);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (!params.id) return;
    setApproval({
      id: params.id,
      alias: 'Pending shared access change',
      description: 'This request is awaiting your decision.',
      status: 'PENDING',
    });
  }, [params.id]);

  const handleDecision = async (action: 'approve' | 'reject') => {
    if (!params.id) return;
    setSubmitting(true);
    try {
      const payload = {
        signer: appUser.primarySigner,
        address: appUser.address,
        secretKey: appUser.secretKeys[0],
        body: {
          id: params.id,
          rejectionReason: action === 'reject' ? 'Rejected from UI' : undefined,
        },
      };
      const response =
        action === 'approve' ? await approve(payload) : await reject(payload);
      if ('data' in response && response.data) {
        showNotification(
          'success',
          action === 'approve' ? 'Approval submitted' : 'Rejection submitted',
        );
        navigate('/dashboard/shared-access/approvals');
      } else {
        showNotification('error', 'Unable to complete action');
      }
    } catch (error) {
      showNotification('error', 'Unable to complete action');
    } finally {
      setSubmitting(false);
    }
  };

  if (!approval) {
    return <div className="p-6">Loading approval...</div>;
  }

  return (
    <div className="flex flex-col gap-5 p-4 md:p-6">
      <Header isHomeView={false} />
      <div className="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm">
        <p className="text-sm font-semibold uppercase tracking-wide text-primary-600">
          Approval request
        </p>
        <h2 className="mt-2 text-2xl font-semibold text-primary-800">
          {approval.alias}
        </h2>
        <p className="mt-2 text-sm text-gray-500">{approval.description}</p>
        <div className="mt-6 flex flex-wrap gap-3">
          <Button
            label={submitting ? 'Submitting...' : 'Approve'}
            onclick={() => handleDecision('approve')}
            additionalClasses="h-11 md:w-auto px-8"
          />
          <ButtonSecondary
            label={submitting ? 'Submitting...' : 'Reject'}
            onclick={() => handleDecision('reject')}
            additionalClasses="h-11 md:w-auto px-8"
          />
        </div>
      </div>
    </div>
  );
}
