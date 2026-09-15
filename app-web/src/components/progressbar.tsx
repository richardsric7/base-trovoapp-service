import { formatHistoryNumber } from '../utils/fomatNumber';

type Props = {
  normalizedProgress: number;
  normalizedDaysProgress: number;
  tokensRemaining: number;
  tokenizedAsset: any;
  totalDays: number;
  daysProgress: number;
};

const ProgressBar = ({ value, color }: { value: number; color: string }) => (
  <div className="w-full h-2.5 bg-gray-300 rounded-full overflow-hidden">
    <div
      className="h-full rounded-full transition-all duration-300"
      style={{
        width: `${value * 100}%`,
        backgroundColor: color,
      }}
    />
  </div>
);

export default function AssetProgressCard({
  normalizedProgress,
  normalizedDaysProgress,
  tokensRemaining,
  tokenizedAsset,
  totalDays,
  daysProgress,
}: Props) {
  return (
    <div className="w-full px-3 py-2.5">
      <div className=" px-12 py-4">
        <div className="flex flex-col gap-3">
          {/* Tokens Remaining */}
          <div>
            <p className="text-sm font-semibold">Tokens Remaining</p>

            <ProgressBar value={normalizedProgress} color={'#108606'} />

            <p className="text-xs mt-1">
              <span>
                {formatHistoryNumber(tokensRemaining, 6)}{' '}
                {tokenizedAsset.assetCode}
              </span>{' '}
              remaining out of{' '}
              <span>
                {formatHistoryNumber(tokenizedAsset.numberOfTokenToBeSold, 6)}{' '}
                {tokenizedAsset.assetCode}
              </span>
            </p>
          </div>

          {/* Days Remaining */}
          <div>
            <p className="text-sm font-semibold">Days Remaining</p>

            {daysProgress >= 0 && (
              <>
                <ProgressBar
                  value={normalizedDaysProgress}
                  color="#60a5fa" // Tailwind blue-400
                />

                <p className="text-xs mt-1">
                  <span>{totalDays - daysProgress}</span> days remaining out of{' '}
                  <span>{totalDays} days</span>
                </p>
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
