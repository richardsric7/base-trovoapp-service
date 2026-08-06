import { useEffect, useState } from 'react';

type CountdownProps = {
  startDate: Date;
  isColumn?: boolean;
  isDark?: boolean;
};

export default function Countdown({
  startDate,
  isColumn = false,
  isDark = false,
}: CountdownProps) {
  const [now, setNow] = useState(new Date());

  // Timer (equivalent to Timer.periodic)
  useEffect(() => {
    const diffDays =
      (startDate.getTime() - new Date().getTime()) / (1000 * 60 * 60 * 24);

    // Only run interval when past date (like your Flutter logic)
    if (diffDays <= 0) {
      const interval = setInterval(() => {
        setNow(new Date());
      }, 1000);

      return () => clearInterval(interval); // cleanup like dispose()
    }
  }, [startDate]);

  const diff = Math.floor(
    (startDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24),
  );

  const formatTime = (durationMs: number) => {
    const totalSeconds = Math.floor(durationMs / 1000);
    const hours = Math.floor(totalSeconds / 3600);
    const minutes = Math.floor((totalSeconds % 3600) / 60);
    const seconds = totalSeconds % 60;

    return [
      hours.toString(),
      minutes.toString().padStart(2, '0'),
      seconds.toString().padStart(2, '0'),
    ];
  };

  const cardStyle = 'px-[5px] py-[5px] text-[12px] rounded-md border shadow-sm';

  const cardColor = isDark
    ? 'bg-blue-900 border-blue-200 text-blue-200'
    : 'bg-gray-200 border-blue-300 text-blue-600';

  // === CASE 1: DAYS LEFT ===
  if (diff > 0) {
    const digits = diff.toString().split('');

    return isColumn ? (
      <div className="flex flex-col items-start">
        <div className="flex items-center">
          {digits.map((d, i) => (
            <div key={i} className="flex items-center">
              <div className={`${cardStyle} ${cardColor}`}>{d}</div>
              <div className="w-[2px]" />
            </div>
          ))}
        </div>

        <div className="h-[5px]" />

        <span className="text-[15px] text-blue-500">days left</span>
      </div>
    ) : (
      <div className="flex items-center">
        {digits.map((d, i) => (
          <div key={i} className="flex items-center">
            <div className={`${cardStyle} ${cardColor}`}>{d}</div>
            <div className="w-[2px]" />
          </div>
        ))}

        <div className="w-[4px]" />

        <span className="text-[12px] text-black">days left</span>
      </div>
    );
  }

  // === CASE 2: HOURS/MINUTES/SECONDS ===
  const endTime = new Date(
    now.getFullYear(),
    now.getMonth(),
    now.getDate() + 1,
    0,
    0,
    0,
  );

  let remaining = endTime.getTime() - now.getTime();
  if (remaining < 0) remaining = 0;

  const timeParts = formatTime(remaining);

  const renderTime = () =>
    timeParts.map((t, i) => (
      <div key={i} className="flex items-center">
        <div className={`${cardStyle} ${cardColor}`}>{t}</div>

        {i < timeParts.length - 1 && (
          <>
            <div className="w-[2px]" />
            <span className="text-[12px] text-black">:</span>
            <div className="w-[2px]" />
          </>
        )}
      </div>
    ));

  return isColumn ? (
    <div className="flex flex-col items-center">
      <div className="flex items-center">{renderTime()}</div>

      <div className="h-[5px]" />

      <span className="text-[15px] text-blue-500">hours left</span>
    </div>
  ) : (
    <div className="flex items-center">
      {renderTime()}

      <div className="w-[4px]" />

      <span className="text-[12px] text-black">hours left</span>
    </div>
  );
}
