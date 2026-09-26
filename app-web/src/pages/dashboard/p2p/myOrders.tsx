import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Header from '../../../components/header';
import { useListMyP2POrdersQuery } from '../../../store/api/p2pApis';
import { useP2PIdentity } from '../../../hooks/useP2PIdentity';
import { isCustomer } from '../../../types/p2p';
import { P2PListCard, P2PEmptyState, P2PStatusPill, p2p } from '../../../components/p2p/P2PTheme';

// P2PMyOrders (Plan Section 96.10): My Orders/My Trades, built on the
// history page's pattern (filter bar + rounded-panel list) rather than
// CustomTable, since that page is app-web's actual highest-polish
// precedent for exactly this shape (Plan Section 92).
export default function P2PMyOrders() {
  const navigate = useNavigate();
  const { creds, username, ready } = useP2PIdentity();
  const [role, setRole] = useState<'' | 'customer' | 'merchant'>('');

  const { data, isLoading, refetch } = useListMyP2POrdersQuery(
    { creds, role: role || undefined, page: 1, pageSize: 20 },
    { skip: !ready },
  );

  return (
    <div className="flex flex-col space-y-5 p-3 max-w-2xl">
      <Header />
      <h1 className="text-xl font-bold text-primary-800 px-2">My orders</h1>

      <div className="flex gap-2">
        {[
          { value: '' as const, label: 'All' },
          { value: 'customer' as const, label: 'As customer' },
          { value: 'merchant' as const, label: 'As merchant' },
        ].map((tab) => (
          <button
            key={tab.value}
            onClick={() => setRole(tab.value)}
            className="flex-1 py-2 rounded-xl text-sm font-semibold"
            style={{ backgroundColor: role === tab.value ? p2p.brandDark : 'white', color: role === tab.value ? 'white' : '#374151' }}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {isLoading ? (
        <div className="text-center py-16 text-gray-400">Loading...</div>
      ) : !data?.data?.length ? (
        <P2PEmptyState
          message="You have no orders yet."
          ctaLabel="Browse the marketplace"
          onCta={() => navigate('/dashboard/p2p')}
        />
      ) : (
        data.data.map((order) => {
          const counterparty = isCustomer(order, username) ? order.merchantUsername : order.customerUsername;
          return (
            <P2PListCard key={order.id} onClick={() => navigate(`/dashboard/p2p/order/${order.id}`)}>
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-bold">
                    {order.specifiedAssetAmount} {order.asset}
                  </p>
                  <p className="text-sm text-gray-500">with {counterparty}</p>
                </div>
                <P2PStatusPill status={order.orderStatus} disputed={order.isDisputed} />
              </div>
            </P2PListCard>
          );
        })
      )}
      <button onClick={() => refetch()} className="text-sm text-primary-800 underline">
        Refresh
      </button>
    </div>
  );
}
