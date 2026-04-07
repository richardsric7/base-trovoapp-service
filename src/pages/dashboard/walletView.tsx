import Button from '../../components/button';
import Header from '../../components/header';
import Tabs from '../../components/tabs';
import { useNavigate } from 'react-router-dom';
import { useEffect, useRef, useState } from 'react';
import WalletCard from '../../components/walletCard';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../../store/reduxStore';
import {
  setActiveWallet as setActiveWalletToStore,
  setActiveAsset as setActiveAssetToStore,
} from '../../store/appStateSlice';
import React from 'react';
import {
  formatToDecimal,
  canInitiate,
  totalWalletBalanceInCurrency,
} from '../../utils/utilities';
import { Wallet } from '../../types/wallet';
import {
  TEDropdown,
  TEDropdownItem,
  TEDropdownMenu,
  TEDropdownToggle,
  TERipple,
} from 'tw-elements-react';
import { Asset } from '../../types/asset';
import AssetItem from '../../components/assetItem';
import WalletOperations from '../../components/walletOperations';

export default function WalletView() {
  const ref = useRef<HTMLDivElement>(null);
  const dispatch = useDispatch();
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const appState = useSelector((state: RootState) => state.appState!);
  const navigate = useNavigate();
  let mutableWalletArray = [...appUser.userWallets];
  const [wallets, setWallets] = useState(
    mutableWalletArray.sort((w) => (w.primaryWallet ? 0 : 1)),
  );
  const [activeWalletIndex, setActiveWalletIndex] = useState(0);
  const [activeWallet, setActiveWallet] = useState<Wallet>(wallets[0]);
  const [assetFilterMode, setAssetFilterMode] = useState(0);
  const [walletActionMode, setWalletActionMode] = useState(1);
  const itemRefs = useRef<HTMLDivElement[]>([]);
  const fiatRates = useSelector((state: RootState) => state.cache.fiatRates);
  const [selectedAsset, setSelectedAsset] = useState<Asset>(
    activeWallet?.claimedAssets.find(
      (a) => a.assetCode === '' && a.assetIssuer === '',
    )!,
  );
  const [currentTabIndex, setCurrentTabIndex] = useState(1);
  const [gas, setGas] = useState(
    activeWallet?.claimedAssets.find(
      (a) => a.assetCode === '' && a.assetIssuer === '',
    )?.amount!,
  );

  useEffect(() => {
    console.log('using effect here!');
    setWallets(
      mutableWalletArray
        .filter((w) => {
          if (currentTabIndex === 1) {
            return w;
          } else if (
            currentTabIndex === 2 &&
            (w.owner === appUser.username ||
              !w.sharedAccessEnabled ||
              w.walletThreshold === 2)
          ) {
            return w;
          } else if (currentTabIndex === 3 && w.sharedAccessEnabled) {
            return w;
          }
        })
        .sort((w) => (w.primaryWallet ? 0 : 1)),
    );

    setActiveWalletIndex(0);
    const ref = itemRefs.current[0];
    if (ref) {
      // Change background color
      ref.scrollIntoView({
        behavior: 'smooth',
        block: 'nearest',
        inline: 'center',
      });
    }
  }, [currentTabIndex]);

  useEffect(() => {
    setGas(
      activeWallet?.claimedAssets.find(
        (a) => a.assetCode === '' && a.assetIssuer === '',
      )?.amount!,
    );
  }, [activeWallet, walletActionMode]);

  // Ensure `itemRefs` has refs for each item on each render
  if (itemRefs.current.length !== wallets.length) {
    // Adjust the array size based on items length
    itemRefs.current = Array(wallets.length)
      .fill(null)
      .map(
        (_, i) =>
          itemRefs.current[i] || React.createRef<HTMLDivElement>().current,
      );
  }

  const userWallets = wallets.map((wallet, index) => (
    // const walletRef = useRef<HTMLDivElement>(null);
    <WalletCard
      localCurrencyBalance={formatToDecimal(
        totalWalletBalanceInCurrency(
          wallet,
          (fiatRates as Record<string, number>)[
            appUser.currency.toUpperCase()
          ] || 0,
        ),
      )}
      usdBalance={formatToDecimal(
        totalWalletBalanceInCurrency(
          wallet,
          (fiatRates as { USD: number }).USD,
        ),
      )}
      currency={appUser.currency.toUpperCase()}
      alias={wallet.alias}
      isSharedAccess={wallet.sharedAccessEnabled}
      walletType={wallet.walletType!}
      key={index}
      ref={(el) => (itemRefs.current[index] = el!)}
    />
  ));

  const assetClassOptions = [
    {
      text: 'Asset Tokens',
      value: 0,
    },
    {
      text: 'Other Tokens',
      value: 1,
    },
  ];

  const dropdownItems = assetClassOptions.map((item) => (
    <TEDropdownItem
      key={`${new Date().getTime()}${item.text.replace(' ', '')}`}
      id={`${new Date().getTime()}${item.text.replace(' ', '')}`}
      className="bg-primary-200"
    >
      <button
        type="button"
        className="block w-full cursor-pointer hover:bg-primary-200 bg-primary-100 whitespace-nowrap px-8 py-2 text-sm text-left font-normal pointer-events-auto active:text-primary-800 focus:hover:bg-primary-200 focus:text-primary-800 focus:outline-none active:no-underline"
        onClick={(e: React.MouseEvent<HTMLButtonElement>) => {
          e.preventDefault;
          setAssetFilterMode(item.value);
        }}
      >
        {item.text}
      </button>
    </TEDropdownItem>
  ));

  return (
    <div className="flex text-primary-800 text-sm overflow-x-hidden md:text-md flex-col space-y-5 p-5">
      <Header isHomeView />
      <div className="flex md:h-screen w-full items-center justify-center">
        <div className="h-full w-full md:p-3">
          <div className="flex md:pt-5 p-5 rounded-lg flex-col items-center bg-primary-100 max-w-5xl min-h-full">
            <div className="w-full flex justify-between space-x-5">
              <Tabs
                tabList={['All Wallets', 'My Wallets', 'Shared Wallets']}
                onTabChanged={(index) => {
                  setCurrentTabIndex(index);
                }}
              />
              <div className="w-1/4">
                <Button
                  label="Add Subwallet"
                  onclick={() => {
                    /* */
                  }}
                />
              </div>
            </div>
            <div className="relative max-w-full">
              <div
                className="flex overflow-x-hidden whitespace-nowrap w-full space-x-4 "
                ref={ref}
              >
                {userWallets.length > 0 && userWallets}
              </div>
              {activeWalletIndex > 0 && (
                <button
                  onClick={() => {
                    activeWalletIndex > 0 &&
                      setActiveWalletIndex(activeWalletIndex - 1);
                    setActiveWallet(wallets[activeWalletIndex - 1]);
                    if (!canInitiate(wallets[activeWalletIndex - 1])) {
                      setWalletActionMode(1);
                    }
                    setSelectedAsset(
                      activeWallet?.claimedAssets.find(
                        (a) => a.assetCode === '' && a.assetIssuer === '',
                      )!,
                    );
                    const ref = itemRefs.current[activeWalletIndex - 1];
                    if (ref) {
                      // Change background color
                      ref.scrollIntoView({
                        behavior: 'smooth',
                        block: 'nearest',
                        inline: 'center',
                      });
                    }
                  }}
                  className="absolute ml-2 text-3xl top-1/2 left-0 transform -translate-y-1/2 bg-gray-200 w-10  h-10 rounded-full hover:bg-gray-300"
                >
                  ←
                </button>
              )}

              {activeWalletIndex + 1 <= itemRefs.current.length - 1 && (
                <button
                  onClick={() => {
                    activeWalletIndex <= itemRefs.current.length &&
                      setActiveWalletIndex(activeWalletIndex + 1);
                    setActiveWallet(wallets[activeWalletIndex + 1]);
                    setSelectedAsset(
                      activeWallet?.claimedAssets.find(
                        (a) => a.assetCode === '' && a.assetIssuer === '',
                      )!,
                    );

                    if (!canInitiate(wallets[activeWalletIndex + 1])) {
                      setWalletActionMode(1);
                    }
                    // resetForm();

                    const ref = itemRefs.current[activeWalletIndex + 1];
                    if (ref) {
                      // Change background color
                      ref.scrollIntoView({
                        behavior: 'smooth',
                        block: 'nearest',
                        inline: 'center',
                      });
                    }
                  }}
                  className="absolute mr-2 text-3xl top-1/2 right-0 transform -translate-y-1/2 bg-gray-200 w-10  h-10 rounded-full hover:bg-gray-300"
                >
                  →
                </button>
              )}
            </div>
            <div className="w-full ring-2 ring-primary-600 rounded-lg p-3">
              <div className="flex justify-between w-full">
                <div className="flex space-x-2 items-center">
                  <img src="/images/gas.png" alt="gas" width="30" height="30" />
                  <span>Gas</span>
                </div>
                <img
                  src="/images/export.png"
                  alt="export"
                  width="30"
                  height="30"
                />
              </div>
              <p>{formatToDecimal(gas)} XBN</p>
            </div>
            <div className="w-full p-3 mt-7">
              <div className="flex justify-between w-full">
                <p className="font-bold font-montserratSemiBold text-lg">
                  My Assets
                </p>
                <TEDropdown className="flex justify-center">
                  <TERipple className="w-full" rippleColor="light">
                    <TEDropdownToggle className="flex items-center whitespace-nowrap rounded bg-primary-100 hover:bg-primary-200 text-primary-800 p-2 justify-between rounded-xl w-full">
                      <img src="/images/filter.png" alt="filter" />
                    </TEDropdownToggle>
                  </TERipple>
                  <TEDropdownMenu className="w-full bg-primary-100">
                    {dropdownItems}
                  </TEDropdownMenu>
                </TEDropdown>
              </div>
              <p>
                <span>Sort by:</span>{' '}
                <span className="font-bold font-montserratSemiBold">
                  {assetClassOptions[assetFilterMode].text}
                </span>
              </p>
            </div>
            <div className="w-full mt-4 space-y-4">
              {activeWallet.claimedAssets.filter((a) => {
                if (assetFilterMode === 0) return a.tokenizedAsset;
                else return !a.tokenizedAsset;
              }).length > 1 ? ( // since every wallet must have an XBN asset check length greater than 1
                // if filter is set on other tokens
                assetFilterMode === 0 ? (
                  activeWallet.claimedAssets.map(
                    (asset, index) =>
                      asset.tokenizedAsset == true && (
                        <AssetItem
                          key={index}
                          image={asset.imageUrl}
                          assetCode={asset.assetCode}
                          amount={asset.amount.toString()}
                          usdPrice={asset.usdPrice.toString()}
                          nativePrice={asset.nativePrice.toString()}
                          currency={appUser.currency}
                          onclick={() => {
                            dispatch(setActiveWalletToStore(activeWallet));
                            dispatch(
                              setActiveAssetToStore(
                                `${asset.assetCode}|${asset.assetIssuer}`,
                              ),
                            );
                            navigate('/dashboard/asset-details');
                          }}
                        />
                      ),
                  )
                ) : (
                  activeWallet.claimedAssets.map(
                    (asset, index) =>
                      !asset.tokenizedAsset && (
                        <AssetItem
                          key={index}
                          image={asset.imageUrl}
                          assetCode={asset.assetCode}
                          amount={asset.amount.toString()}
                          usdPrice={asset.usdPrice.toString()}
                          nativePrice={asset.nativePrice.toString()}
                          currency={appUser.currency}
                          onclick={() => {
                            dispatch(setActiveWalletToStore(activeWallet));
                            dispatch(
                              setActiveAssetToStore(
                                `${asset.assetCode}|${asset.assetIssuer}`,
                              ),
                            );
                            console.log(appState);
                            navigate('/dashboard/asset-details');
                          }}
                        />
                      ),
                  )
                )
              ) : (
                <div className="flex items-center justify-center w-full h-60">
                  <p className="text-lg">Nothing to show here</p>
                </div>
              )}
            </div>
          </div>
        </div>
        <WalletOperations
          wallets={wallets}
          selectedAsset={selectedAsset}
          activeWallet={activeWallet}
          onWalletChanged={(wallet, index) => {
            setActiveWallet(wallet);
            setActiveWalletIndex(index);
            setSelectedAsset(
              activeWallet?.claimedAssets.find(
                (a) => a.assetCode === '' && a.assetIssuer === '',
              )!,
            );
            const ref = itemRefs.current[index];
            if (ref) {
              // Change background color
              ref.scrollIntoView({
                behavior: 'smooth',
                block: 'nearest',
                inline: 'center',
              });
            }
          }}
          onAssetChanged={(asset) => {
            setSelectedAsset(asset);
          }}
          isWalletDetailsPage={false}
        />
      </div>
    </div>
  );
}
