import { Link, useNavigate } from 'react-router-dom';
import { useListPMAssetsQuery } from '../../../store/api/publicMarketsApis';
import { AssetLogo, changeColor, money, pct, StatusBadge } from './ui';

// PMHomeSection is the home page's Public Markets strip: the first few
// listed stocks and bonds, each opening its page, and View all.
export default function PMHomeSection() {
  const navigate = useNavigate();
  const { data } = useListPMAssetsQuery({});
  const assets = (data?.assets ?? []).slice(0, 4);
  if (!assets.length) return null;
  return (
    <div className="md:px-5 pb-10 w-full space-y-5">
      <div className="flex items-center w-full justify-between">
        <p className="font-montserratSemiBold text-md">PUBLIC MARKETS</p>
        <Link className="flex space-x-5 items-center" to={'public-markets'}>
          <p className="test-xs text-primary-400 underline">View all</p>
        </Link>
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
        {assets.map((a) => (
          <button
            key={a.assetCode}
            onClick={() => navigate(`/dashboard/public-markets/asset/${encodeURIComponent(a.assetCode)}`)}
            className="flex items-center gap-3 bg-white rounded-2xl p-3 text-left hover:bg-primary-100"
          >
            <AssetLogo asset={a} size={40} />
            <div className="flex-1 min-w-0">
              <p className="font-bold text-primary-800">{a.ticker}</p>
              <p className="text-xs text-gray-500 truncate">{a.shortName || a.name}</p>
            </div>
            {a.status === 'open' ? (
              <div className="text-right">
                <p className="font-semibold text-sm">{money(a.price, a.currency)}</p>
                <p className="text-xs font-semibold" style={{ color: changeColor(a.dayChangePercent) }}>
                  {pct(a.dayChangePercent)}
                </p>
              </div>
            ) : (
              <StatusBadge asset={a} />
            )}
          </button>
        ))}
      </div>
    </div>
  );
}
