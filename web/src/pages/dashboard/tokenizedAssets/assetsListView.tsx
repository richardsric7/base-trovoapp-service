import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import Header from '../../../components/header';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../../../store/reduxStore';
import { Wallet } from '../../../types/wallet';
import {
  useFetchExpressedInterestsQuery,
  useFetchTokenizedAssetsQuery,
} from '../../../store/api/walletApis';
import { Encryptor } from '../../../utils/encryptor';
import { TokenizedAsset } from '../../../types/tokenizedAsset';
import AssetListItem from '../../../components/assetListItem';
import {
  BuyTokenModal,
  ExpressInterestModal,
  KycUnverifiedModal,
  P2pComingSoonModal,
  SoldOutModal,
} from '../../../components/tokenizedAssetActionModals';
import { setActiveTokenizedAsset } from '../../../store/appStateSlice';
import { ExpressedInterest } from '../../../types/expressedInterest';

export default function TokenizedAssetsListView() {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const [showSubscribeModal, setShowSubscribeModal] = useState(false);
  const [showSoldOutModal, setShowSoldOutModal] = useState(false);
  const [showKycUnverifiedModal, setShowKycUnverifiedModal] = useState(false);
  const [showBuyTokenModal, setShowBuyTokenModal] = useState(false);
  const [showP2pComingSoonModal, setShowP2pComingSoonModal] = useState(false);
  const [asset, setAsset] = useState(
    useSelector((state: RootState) => state.appState.activeTokenizedAsset),
  );

  useEffect(() => {
    (async () => {
      const encryptor = new Encryptor();
      setSecretKey(await encryptor.getSecretKey(appUser));
    })();
  }, []);

  const appUser = useSelector((state: RootState) => state.auth.user!);
  const [secretKey, setSecretKey] = useState('');

  const [primaryWallet] = useState<Wallet>(
    appUser.userWallets.find((w) => w.primaryWallet)!,
  );

  const listType =
    searchParams.get('rel') == 'primary'
      ? { label: 'Primary Offers', code: 0 }
      : { label: 'Secondary Listing', code: 1 };

  const { data, isLoading } = useFetchTokenizedAssetsQuery(
    {
      signer: primaryWallet?.signer ?? '',
      publicKey: primaryWallet.publicKey,
      secretKey,
      body: { status: listType.code },
    },
    {
      skip: !secretKey || !primaryWallet,
    },
  );

  var typedAssets: TokenizedAsset[] = [];

  if (listType.code == 0) {
    const { data: expressedInterest } = useFetchExpressedInterestsQuery(
      {
        signer: primaryWallet?.signer ?? '',
        publicKey: primaryWallet.publicKey,
        secretKey,
        body: { status: 1 },
      },
      {
        skip: !secretKey || !primaryWallet,
      },
    );

    console.log('expressed interests', expressedInterest);

    data?.records.map((x: TokenizedAsset) => {
      var record = expressedInterest?.records?.find(
        (a: ExpressedInterest) => a.assetCode == x.assetCode,
      );
      var y = {
        ...x,
        expressedInterestAmount: record != null ? record.amount : 0,
      };
      typedAssets = [...typedAssets, y];
    });
  }

  return (
    <div className="flex text-primary-800 text-sm md:text-md flex-col space-y-5 p-3">
      <Header />
      <div className="flex flex-col items-center justify-center w-full">
        <div className="flex px-5 mb-10 space-x-5 w-full items-center">
          <button type="button" onClick={() => navigate(-1)}>
            <img src="/images/arrowBack.png" alt="arrow back" />
          </button>
          <p className="font-montserratSemiBold text-lg xl:text-xl">
            {listType.label}
          </p>
        </div>
        <div className="flex min-h-screen w-full justify-center">
          <div className="flex md:pt-5 p-5 w-full rounded-xl flex-col items-center bg-primary-100 max-w-5xl min-h-full">
            <div className="w-full space-y-6 px-3 md:px-5 md:py-5 flex flex-col justify-center space-y-5 md:space-y-3">
              {isLoading ? (
                <p className="text-center">Loading primary offers...</p>
              ) : typedAssets.length == 0 ? (
                <div className="w-full">
                  <div className="md:px-5 w-full">
                    <div className="w-full flex flex-col justify-center space-y-5 md:mb-10">
                      <img
                        className="self-center"
                        src="/images/empty.png"
                        alt=""
                      />
                      <p className="font-montserratSemiBold text-md text-center">
                        No Assets Yet
                      </p>
                    </div>
                  </div>
                </div>
              ) : (
                <div className="space-y-5">
                  {typedAssets.map((asset, index) => (
                    <AssetListItem
                      key={index}
                      asset={asset}
                      isSubscribed={asset.expressedInterestAmount > 0}
                      onclick={() => {
                        dispatch(setActiveTokenizedAsset(asset));
                        navigate(
                          `/dashboard/tokenized-asset?salesList=0&assetCode=${asset.assetCode}`,
                        );
                      }}
                      onBuy={() => {
                        setAsset(asset);
                        setShowBuyTokenModal(true);
                      }}
                      onExpressInterest={() => {
                        setAsset(asset);
                        setShowSubscribeModal(true);
                      }}
                    />
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
      {showSubscribeModal && (
        <ExpressInterestModal
          show={showSubscribeModal}
          onClose={() => {
            setShowSubscribeModal(false);
          }}
          asset={asset}
          onSuccess={() => {
            setShowSubscribeModal(false);
          }}
        />
      )}
      {/* Public key copied modal */}
      <P2pComingSoonModal
        show={showP2pComingSoonModal}
        onClose={() => {
          setShowP2pComingSoonModal(false);
        }}
      />
      <KycUnverifiedModal
        show={showKycUnverifiedModal}
        onClose={() => {
          setShowKycUnverifiedModal(false);
        }}
      />
      <SoldOutModal
        show={showSoldOutModal}
        onClose={() => {
          setShowSoldOutModal(false);
        }}
      />
      {showBuyTokenModal && (
        <BuyTokenModal
          show={showBuyTokenModal}
          onClose={() => {
            setShowBuyTokenModal(false);
          }}
          asset={asset}
        />
      )}
    </div>
  );
}
