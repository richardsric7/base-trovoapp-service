import { useEffect, useState } from 'react';
import { useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import Header from '../../../components/header';
import Button from '../../../components/button';
import { RootState } from '../../../store/reduxStore';
import { useLazyGetApprovalsQuery } from '../../../store/api/sharedAccessApis';
import { showNotification } from '../../../utils/showToaster';

type ApprovalItem = {
  id: string;
  alias?: string;
  status?: string;
  description?: string;
};

export default function SharedAccessApprovals() {
  const navigate = useNavigate();
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const [fetchApprovals] = useLazyGetApprovalsQuery();
  const [approvals, setApprovals] = useState<ApprovalItem[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const load = async () => {
      setLoading(true);
      try {
        const result = await fetchApprovals({
          signer: appUser.primarySigner,
          address: appUser.address,
          secretKey: appUser.secretKeys[0],
          body: {},
        });
        if ('data' in result && result.data) {
          const data = result.data as any;
          setApprovals(Array.isArray(data) ? data : (data.approvals ?? []));
        }
      } catch (error) {
        showNotification('error', 'Unable to fetch approvals');
      } finally {
        setLoading(false);
      }
    };

    load();
  }, [appUser, fetchApprovals]);

  return (
    <div className="flex flex-col gap-5 p-4 md:p-6">
      <Header isHomeView={false} />
      <div className="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm">
        <h2 className="text-2xl font-semibold text-primary-800">
          Pending approvals
        </h2>
        <p className="mt-2 text-sm text-gray-500">
          Review requests that require your approval before shared access
          changes are committed.
        </p>
        {loading ? (
          <div className="mt-6 text-sm text-gray-500">Loading approvals...</div>
        ) : approvals.length === 0 ? (
          <div className="mt-6 text-sm text-gray-500">
            There are no pending approvals right now.
          </div>
        ) : (
          <div className="mt-6 flex flex-col gap-3">
            {approvals.map((item) => (
              <div
                key={item.id}
                className="rounded-2xl border border-gray-200 p-4"
              >
                <div className="flex flex-wrap items-center justify-between gap-3">
                  <div>
                    <p className="text-base font-semibold text-primary-800">
                      {item.alias || 'Shared access request'}
                    </p>
                    <p className="mt-1 text-sm text-gray-500">
                      {item.description || 'Approval pending'}
                    </p>
                  </div>
                  <div className="rounded-full bg-amber-100 px-3 py-1 text-sm font-semibold text-amber-700">
                    {item.status || 'PENDING'}
                  </div>
                </div>
                <div className="mt-4">
                  <Button
                    label="Review"
                    onclick={() =>
                      navigate(`/dashboard/shared-access/approvals/${item.id}`)
                    }
                    additionalClasses="h-11 md:w-auto px-8"
                  />
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
