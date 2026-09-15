type Props = {
  image: string;
  assetCode: string;
  amount: string;
  usdPrice?: string;
  nativePrice?: string;
  currency?: string;
  onclick: () => void;
};

function AssetItem({
  image,
  assetCode,
  amount,
  usdPrice,
  nativePrice,
  currency,
  onclick,
}: Props) {
  return (
    <button
      className="rounded-2xl w-full bg-white flex justify-between items-center px-2 md:px-7 py-3 md:py-5"
      onClick={onclick}
      type="button"
    >
      <div className="flex space-x-5 items-center">
        <img
          className="h-10 w-10 rounded-full"
          src={image}
          alt={assetCode === '' ? 'ETH' : assetCode}
        />
        <div className="flex items-start flex-col">
          <p className="font-semibold md:text-xl">{assetCode}</p>
          <p className="text-xs">
            {nativePrice} {currency}
          </p>
        </div>
      </div>
      <div className="flex items-end flex-col">
        <p className="font-semibold md:text-xl">{amount}</p>
        <p className="text-xs">
          {usdPrice} {currency}
        </p>
      </div>
    </button>
  );
}

export default AssetItem;
