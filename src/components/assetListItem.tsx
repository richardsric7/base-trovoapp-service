type Props = {
  image: string;
  assetName: string;
  assetClass: string;
  isSubscribed?: boolean;
  onclick: () => void;
};

function AssetListItem({
  image,
  assetName,
  assetClass,
  isSubscribed = false,
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
          <p className="font-semibold md:text-xl">{assetName}</p>
          <p className="text-xs">{assetClass}</p>
        </div>
      </div>
      {isSubscribed ? (
        <div className="flex items-center space-x-3">
          <p>Subscribed</p>
          <img src="/images/tick.png" alt="add" />
        </div>
      ) : (
        <div className="flex items-center space-x-3">
          <p>Subscribe</p>
          <img src="/images/plus.png" alt="add" />
        </div>
      )}
    </button>
  );
}

export default AssetListItem;
