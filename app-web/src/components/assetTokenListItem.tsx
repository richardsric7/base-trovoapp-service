type Props = {
  image: string;
  assetCode: string;
  usdPrice: string;
  amount?: string;
  valueInFiat?: string;
  onclick: () => void;
};

function AssetTokenItem({
  image,
  assetCode,
  usdPrice,
  amount,
  valueInFiat,
  onclick,
}: Props) {
  return (
    <button
      className="rounded-2xl w-full bg-white flex justify-between items-center px-2 md:px-7 py-3 md:py-5"
      onClick={onclick}
      type="button"
    >
      <div className="flex space-x-5 items-center">
        <img className="h-10" src={image} alt="atlantis 1" />
        <div className="flex items-start flex-col">
          <p className="font-bold font-montserratSemiBold">{assetCode}</p>
          <p className="text-xs">{usdPrice}</p>
        </div>
      </div>
      <div className="flex items-start flex-col">
        <p className="font-bold font-montserratSemiBold">{amount}</p>
        <p className="text-xs">{valueInFiat}</p>
      </div>
    </button>
  );
}

export default AssetTokenItem;
