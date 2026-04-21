import { TokenizedAsset } from '../types/tokenizedAsset';
import Countdown from './countdown';

type Props = {
  asset: TokenizedAsset;
  isSubscribed?: boolean;
  isSecondary?: boolean;
  onclick: () => void;
};

function AssetListItem({
  asset,
  isSecondary = false,
  isSubscribed = false,
  onclick,
}: Props) {
  var tokensRemaining =
    asset.numberOfTokenToBeSold! - asset.quantityOfTokensSold!;
  return (
    <button
      className="rounded-2xl w-full bg-white flex justify-between items-center px-2 md:px-7 py-3 md:py-5"
      onClick={onclick}
      type="button"
    >
      <div className="flex space-x-5">
        <div className="relative inline-block">
          <div className="w-[38px] h-[38px] rounded-full overflow-hidden">
            <img
              src={asset.assetLogo}
              alt="logo"
              className="w-full h-full object-fill"
            />
          </div>

          <div
            className={`absolute top-0 right-0 w-[10px] h-[10px] rounded-full border border-white 
              ${asset.assetAlreadyExists ? 'bg-green-500' : 'bg-yellow-400'}`}
          ></div>
        </div>
        <div className="flex items-start flex-col">
          <p className="font-semibold md:text-xl">{asset.assetName}</p>
          <p className="text-xs">{asset.assetSector}</p>
          {!isSecondary && asset.assetTokenizationStatus == 4 && (
            <div className="mt-3">
              <Countdown startDate={new Date(asset.salesStart)} />
            </div>
          )}
        </div>
      </div>
      {isSecondary ? (
        <></>
      ) : (
        <>
          {asset.assetTokenizationStatus == 5 ? (
            <>
              {tokensRemaining > 0 && (
                <div className="flex items-center space-x-3">
                  <p>Buy</p>
                  <img
                    src="/images/buy_with_fiat.png"
                    height={40}
                    width={40}
                    alt="buy"
                  />
                </div>
              )}
            </>
          ) : (
            <div className="row justify-between">
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
            </div>
          )}
        </>
      )}
    </button>
  );
}

export default AssetListItem;
