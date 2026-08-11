import { TokenizedAsset } from '../types/tokenizedAsset';
import Button from './button';

type Props = {
  asset: TokenizedAsset;
  onclick: () => void;
};

function TokenizationAssetListItem({ asset, onclick }: Props) {
  const getTokenizationStatus = (status: number) => {
    switch (status) {
      case 0:
        return 'Continue';
      case 1:
        return 'Awaiting Fee';
      case 4:
        return 'Approved';
      case 5:
        return 'Primary Sales';
      case 6:
        return 'Secondary Market';
      case 7:
        return 'Liquidated';
      case 8:
        return 'Refunded';
      default:
        return 'Processing';
    }
  };
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
          <p className="font-semibold md:text-xl">
            {asset.assetName.trim().length <= 0
              ? 'No name yet'
              : `${asset.assetName} (${asset.assetCode.toUpperCase()})`}
          </p>
          <p className="text-xs">{asset.assetSector}</p>
        </div>
      </button>
      <div className="w-1/6">
        <Button
          label={getTokenizationStatus(asset.assetTokenizationStatus)}
          onclick={onclick}
        />
      </div>
    </div>
  );
}

export default TokenizationAssetListItem;
