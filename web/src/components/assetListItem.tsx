import { TokenizedAsset } from '../types/tokenizedAsset';
import Button from './button';
import Countdown from './countdown';

type Props = {
  asset: TokenizedAsset;
  isSubscribed?: boolean;
  isSecondary?: boolean;
  onclick: () => void;
  onBuy?: () => void;
  onExpressInterest?: () => void;
};

function AssetListItem({
  asset,
  isSecondary = false,
  isSubscribed = false,
  onclick,
  onBuy,
  onExpressInterest,
}: Props) {
  var tokensRemaining =
    asset.numberOfTokenToBeSold! - asset.quantityOfTokensSold!;
  return (
    <div className="rounded-2xl w-full bg-white flex justify-between items-center px-2 md:px-7 py-3 md:py-5">
      <button onClick={onclick} type="button" className="flex w-full space-x-5">
        <div className="relative inline-block">
          <div className="w-[38px] h-[38px] rounded-full overflow-hidden">
            <img
              src={asset.assetLogo ?? '/images/trovoLogo.png'}
              alt="logo"
              className="w-full h-full object-fill"
              onError={(e) => {
                (e.target as HTMLImageElement).src = '/images/trovoLogo.png';
              }}
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
      </button>
      {isSecondary ? (
        <></>
      ) : (
        <>
          {asset.assetTokenizationStatus == 5 ? (
            <>
              {tokensRemaining > 0 && (
                <div className="w-1/6">
                  <Button
                    label="Buy"
                    onclick={onBuy!}
                    rightIcon={
                      <img
                        src="/images/buy_with_fiat.png"
                        height={25}
                        width={25}
                        alt="buy"
                      />
                    }
                  />
                </div>
              )}
            </>
          ) : (
            <div className="w-1/4 flex justify-end">
              {isSubscribed ? (
                <button
                  onClick={onExpressInterest}
                  className="w-1/4 py-3 px-1 bg-primary-100 rounded-md hover:bg-primary-200"
                >
                  <div className="flex justify-center space-x-3">
                    <img src="/images/tick.png" alt="add" />
                  </div>
                </button>
              ) : (
                <button
                  onClick={onExpressInterest}
                  className="w-full py-3 px-1 bg-primary-100 rounded-md hover:bg-primary-200"
                >
                  <div className="flex justify-center items-center space-x-3">
                    <p>Express Interest</p>
                    <img src="/images/plus.png" alt="add" />
                  </div>
                </button>
              )}
            </div>
          )}
        </>
      )}
    </div>
  );
}

export default AssetListItem;
