type Props = {
  localCurrencyBalance: string;
  usdBalance: string;
  label: string;
};

function WidgetCard({ localCurrencyBalance, usdBalance, label }: Props) {
  return (
    <div
      className="max-w-xl w-full flex text-white font-matahariRegular
     rounded-2xl bg-primary-800 items-center justify-between space-x-3 md:mt-5 md:mb-10 py-5 md:py-10 px-5"
    >
      <div className="flex flex-col space-y-2">
        <p className="md:text-lg">{label}</p>
        <p className="text-xl md:text-2xl font-semibold">
          {localCurrencyBalance}
        </p>
        <p className="text-sm">{usdBalance}</p>
      </div>
      <img className="w-20" src="/images/trovoWhite.png" alt="trovo logo" />
    </div>
  );
}

export default WidgetCard;
