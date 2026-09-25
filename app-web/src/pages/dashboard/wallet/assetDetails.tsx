import { useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import ButtonSecondary from '../../../components/buttonSecondary';
import Button from '../../../components/button';
import Dropdown from '../../../components/dropdown';
import TextInput from '../../../components/textInput';
import Modal from '../../../components/modal';
import Header from '../../../components/header';
import { Asset } from '../../../types/asset';
import { useSelector } from 'react-redux';
import { RootState } from '../../../store/reduxStore';
import { formatToDecimal, getAssetCode } from '../../../utils/utilities';
import { truncateAddress } from '../../../utils/truncateValues';
import { showNotification } from '../../../utils/showToaster';
import WalletOperations from '../../../components/walletOperations';
import { Wallet } from '../../../types/wallet';

type AssetDetailItemProps = {
  asset: Asset;
};

// eslint-disable-next-line
function AssetDetailItem({ asset }: AssetDetailItemProps) {
  return (
    <div className="rounded-2xl w-full bg-white px-4 py-3">
      <div className="flex  flex-col space-y-2 w-full">
        <p className="font-montserratSemiBold text-primary-400 md:text-lg">
          {asset.assetCode}
        </p>
      </div>
    </div>
  );
}

AssetDetailItem.defaultProps = { asset: '' };

export default function AssetDetail() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const appUser = useSelector((state: RootState) => state.auth.user!);

  const [wallet, setWallet] = useState<Wallet>(
    appUser.userWallets.find((w) => w.address == searchParams.get('wallet'))!,
  );
  const [asset, setAsset] = useState<Asset>(
    wallet?.claimedAssets.find(
      (a) =>
        a.assetCode == searchParams.get('assetCode') &&
        a.contractAddress == searchParams?.get('contractAddress'),
    )!,
  );
  let mutableWalletArray = [...appUser.userWallets];
  const [wallets] = useState(
    mutableWalletArray.sort((w) => (w.primaryWallet ? 0 : 1)),
  );
  const fiatRates = useSelector((state: RootState) => state.cache.fiatRates);
  const fiatRate =
    (fiatRates as Record<string, number>)[appUser.currency.toUpperCase()] || 0;
  const curatedAsset = appUser.curatedSwapList.find(
    (a) =>
      a.assetCode == searchParams.get('assetCode') &&
      a.contractAddress == searchParams?.get('contractAddress'),
  );
  const [showSubscribeModal, setShowSubscribeModal] = useState(false);
  const [showSuccessModal, setShowSuccessModal] = useState(false);

  return (
    <div className="flex text-primary-800 text-sm md:text-md flex-col space-y-5 p-3">
      <Header />
      <div className="flex md:h-screen w-full items-center">
        <div className="h-full w-full md:p-3">
          <div className="flex md:pt-5 p-5 rounded-lg flex-col items-center bg-primary-100 max-w-5xl min-h-full">
            <div className="w-full px-3 md:px-5 md:pt-5 flex flex-col md:flex-row space-y-5 md:space-y-0 items-center justify-between">
              <div className="flex space-x-5 w-full md:w-auto items-center">
                <button type="button" onClick={() => navigate(-1)}>
                  <img src="/images/arrowBack.png" alt="arrow back" />
                </button>
                <p className="font-montserratSemiBold text-lg xl:text-xl">
                  {getAssetCode(asset?.assetCode!)}
                </p>
              </div>
              <div className="flex w-full md:w-2/6 space-x-5">
                <ButtonSecondary
                  label="Yield"
                  additionalClasses="bg-primary-600 text-white font-montserratSemiBold"
                  onclick={() => {
                    navigate(
                      `/dashboard/yield?wallet=${wallet.address}&assetCode=${asset.assetCode}&contractAddress=${asset.contractAddress}`,
                    );
                  }}
                />
                <ButtonSecondary
                  label="Early Exit"
                  additionalClasses="font-montserratSemiBold"
                  onclick={() => {
                    navigate(
                      `/dashboard/early-exit?wallet=${wallet.address}&assetCode=${asset.assetCode}&contractAddress=${asset.contractAddress}`,
                    );
                  }}
                />
              </div>
            </div>
            <div className="w-full space-y-6 px-3 md:px-5 md:pt-5 flex flex-col justify-center md:space-x-3 space-y-5 md:space-y-3">
              <div className="flex flex justify-center items-center space-x-1">
                <img
                  src={asset?.imageUrl}
                  alt="asset logo"
                  className="h-12 w-12 rounded-full"
                />
                <p className="font-montserratSemiBold text-lg xl:text-xl">
                  {asset?.assetCode}
                </p>
              </div>
              <div className="w-full flex flex-col items-center space-y-4">
                <p className="flex items-center space-x-2">
                  <span>www.atlantis.com</span>
                </p>
                {curatedAsset != null && (
                  <p className="text-center">{curatedAsset.description}</p>
                )}
              </div>
              <div className="w-full flex flex-col items-center space-y-2">
                <p className="font-montserratSemiBold flex items-center">
                  Price
                </p>
                <p className="flex items-center space-x-2">
                  {asset?.tokenizedAsset ? (
                    <span>
                      1 {asset?.assetCode} ={' '}
                      {formatToDecimal(fiatRate * asset?.usdPrice!)} NGN
                    </span>
                  ) : (
                    <span>
                      1 {asset?.assetCode} ={' '}
                      {formatToDecimal(fiatRate * asset?.usdPrice!)}{' '}
                      {appUser.currency}
                    </span>
                  )}
                </p>
              </div>
              <div className="w-full flex flex-col items-center space-y-2">
                <p className="font-montserratSemiBold flex items-center">
                  Token Contract Address
                </p>
                <p className="flex items-center space-x-2">
                  <span>{truncateAddress(asset?.contractAddress ?? '')}</span>
                  <button
                    type="button"
                    onClick={() =>
                      navigator.clipboard
                        .writeText(asset?.contractAddress ?? '')
                        .then(() => {
                          showNotification('info', 'Address copied!');
                        })
                    }
                  >
                    <img
                      src="/images/copy.png"
                      alt="copy"
                      className="text-primary-500"
                    />
                  </button>
                </p>
              </div>
              {curatedAsset != null && curatedAsset.contactEmail && (
                <div className="w-full flex flex-col items-center space-y-2">
                  <p className="font-montserratSemiBold flex items-center">
                    Contact Email
                  </p>
                  <p className="flex items-center space-x-2">
                    {curatedAsset.contactEmail}
                  </p>
                </div>
              )}
            </div>

            <Modal
              showModal={showSubscribeModal}
              onClose={() => {
                setShowSubscribeModal(false);
              }}
            >
              <div className="flex flex-col items-center w-full px-5 md:px-20 space-y-5 py-5 justify-center">
                <div className="w-full flex justify-center">
                  <img src="/images/thinkingMan.png" alt="thinking man" />
                </div>
                <p>
                  Please enter the amount you would like spend on AE1 tokens
                  when primary sale starts
                </p>
                <div className="w-full">
                  <TextInput
                    inputType="text"
                    label=""
                    placeholder="100"
                    trailingIcon="/images/trovIcon.png"
                    trailingText="CNGN"
                    onInputChange={() => {
                      // console.log('input has changed', newValue);
                    }}
                  />
                </div>
                <p>
                  Are you sure you want to subscribe to the Atlantis 1 Asset
                  with the following wallet?
                </p>
                <Dropdown
                  label="Choose wallet"
                  options={
                    [
                      // 'atlantis_issuer',
                      // 'atlantis_distro',
                      // 'Ephizee',
                      // 'Odogwu_assets',
                    ]
                  }
                  onSelect={(selectedItem) => {
                    console.log('selected item here...', selectedItem);
                  }}
                />
                <div className="w-full md:w-2/4 space-y-3">
                  <Button
                    label="Proceed"
                    additionalClasses="font-montserratSemiBold"
                    onclick={() => {
                      setShowSuccessModal(true);
                      setShowSubscribeModal(false);
                    }}
                  />
                </div>
                <div />
              </div>
            </Modal>
            {/* Address copied modal */}
            <Modal
              showModal={showSuccessModal}
              onClose={() => {
                setShowSuccessModal(false);
              }}
            >
              <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
                <img src="/images/success.png" alt="success" />
                <div className="flex flex-col text-center space-y-5 items-center w-2/3 md:px-10 justify-center">
                  <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
                    You have successfully subscribed to Atlantis 1
                  </p>
                </div>
              </div>
            </Modal>
          </div>
        </div>
        <WalletOperations
          wallets={wallets}
          selectedAsset={asset!}
          activeWallet={wallet!}
          onWalletChanged={(wallet, _) => {
            setWallet(wallet);
            setAsset(
              wallet?.claimedAssets.find(
                (a) => a.assetCode === '' && a.contractAddress === '',
              )!,
            );
          }}
          onAssetChanged={(asset) => {
            setAsset(asset);
          }}
          isWalletDetailsPage={true}
        />
      </div>
    </div>
  );
}
