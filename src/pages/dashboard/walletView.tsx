import Button from '../../components/button';
import TransactionItem from '../../components/transactionItem';
import Header from '../../components/header';
import Tabs from '../../components/tabs';
import { useEffect, useRef, useState } from 'react';
import WalletCard from '../../components/walletCard';
import { useSelector } from 'react-redux';
import { RootState } from '../../store/reduxStore';
import React from 'react';
import {
  calculateFiatValue,
  formatToDecimal,
  totalWalletBalanceInCurrency,
} from '../../utils/utilities';
import { Wallet } from '../../types/wallet';
import AssetTokenItem from '../../components/assetTokenListItem';
import {
  TEDropdown,
  TEDropdownItem,
  TEDropdownMenu,
  TEDropdownToggle,
  TERipple,
} from 'tw-elements-react';
import AssetListItem from '../../components/assetListItem';
import Dropdown from '../../components/dropdown';
import TextInput from '../../components/textInput';

export default function WalletView() {
  // const style = {
  //   backfaceVisibility: 'hidden',
  // };

  const ref = useRef<HTMLDivElement>(null);
  const appUser = useSelector((state: RootState) => state.auth.user!);
  let mutableWalletArrary = [...appUser.userWallets];
  const [wallets] = useState(
    mutableWalletArrary.sort((w) => (w.primaryWallet ? 0 : 1)),
  );
  const [activeWalletIndex, setActiveWalletIndex] = useState(0);
  const [activeWallet, setActiveWallet] = useState<Wallet>(wallets[0]);
  const [assetFilterMode, setAssetFilterMode] = useState('Asset Tokens');
  const [walletActionMode, setWalletActionMode] = useState(0);
  const itemRefs = useRef<HTMLDivElement[]>([]);
  const fiatRates = useSelector((state: RootState) => state.cache.fiatRates);
  const [gas, setGas] = useState(
    activeWallet?.claimedAssets.find(
      (a) => a.assetCode === '' && a.assetIssuer === '',
    )?.amount!,
  );

  useEffect(() => {
    setGas(
      activeWallet?.claimedAssets.find(
        (a) => a.assetCode === '' && a.assetIssuer === '',
      )?.amount!,
    );
  }, [activeWallet]);

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
      key={index}
      ref={(el) => (itemRefs.current[index] = el!)}
    />
  ));

  const dropdownItems = [
    {
      text: 'Asset Tokens',
      value: 'Asset Tokens',
    },
    {
      text: 'Other Tokens',
      value: 'Other Tokens',
    },
  ].map((item) => (
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
    <div className="flex text-primary-800 text-sm overflow-x-hidden md:text-md flex-col space-y-5 p-3">
      <Header isHomeView />
      <div className="flex md:h-screen w-full items-center justify-center">
        <div className="h-full w-full md:p-3">
          <div className="flex md:pt-5 p-5 rounded-lg flex-col items-center bg-primary-100 max-w-5xl min-h-full">
            <div className="w-full flex justify-between">
              <Tabs tabList={['My Wallets', 'Shared Wallets', 'All Wallets']} />
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
                {userWallets}
              </div>
              {activeWalletIndex > 0 && (
                <button
                  onClick={() => {
                    activeWalletIndex > 0 &&
                      setActiveWalletIndex(activeWalletIndex - 1);
                    setActiveWallet(wallets[activeWalletIndex - 1]);

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
              <p>{gas} XBN</p>
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
                  {assetFilterMode}
                </span>
              </p>
            </div>
            <div className="w-full mt-4 space-y-4">
              {activeWallet.claimedAssets.length > 1 ? ( // since every wallet must have an XBN asset check length greater than 1
                // if filter is set on other tokens
                assetFilterMode === 'Other Tokens' ? (
                  activeWallet.claimedAssets.map(
                    (asset) =>
                      asset.assetCode !== '' && (
                        <AssetTokenItem
                          image={asset.imageUrl}
                          assetCode={asset.assetCode}
                          usdPrice={asset.usdPrice.toString()}
                          amount={asset.amount.toString()}
                          valueInFiat={`${calculateFiatValue(
                            asset.amount,
                            (fiatRates as Record<string, number>)[
                              appUser.currency.toUpperCase()
                            ] || 0,
                            asset.usdPrice,
                          )} ${appUser.currency.toUpperCase()}`}
                          onclick={() => {
                            // navigate('/dashboard/tokenized-asset');
                          }}
                        />
                      ),
                  )
                ) : (
                  activeWallet.claimedAssets.map(
                    (asset) =>
                      asset.assetCode !== '' && (
                        <AssetListItem
                          image="/images/avatar.png"
                          assetName="Atlantis 1"
                          assetClass="Real Estate"
                          isSubscribed
                          onclick={() => {
                            // navigate('/dashboard/tokenized-asset');
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
        <div className="hidden md:block w-3/6 flex flex-col space-y-5 h-full py-3 lg:px-3">
          <div className="flex space-y-6 rounded-lg py-5 px-6 flex-col items-center bg-primary-100">
            <div className="flex justify-between items-center space-x-5">
              <button
                className={`flex space-y-3 ring-1 rounded-full px-6 py-2 flex-col items-center ${
                  walletActionMode === 0
                    ? 'ring-primary bg-primary-200 font-bold font-montserratSemiBold'
                    : 'ring-gray-400 text-gray-400'
                }`}
                onClick={() => setWalletActionMode(0)}
              >
                Send
              </button>
              <button
                className={`flex space-y-3 ring-1 rounded-full px-6 py-2 flex-col items-center ${
                  walletActionMode === 1
                    ? 'ring-primary bg-primary-200 font-bold font-montserratSemiBold'
                    : 'ring-gray-400 text-gray-400'
                }`}
                onClick={() => setWalletActionMode(1)}
              >
                Receive
              </button>
              <button
                className={`flex space-y-3 ring-1 rounded-full px-6 py-2 flex-col items-center ${
                  walletActionMode === 2
                    ? 'ring-primary bg-primary-200 font-bold font-montserratSemiBold'
                    : 'ring-gray-400 text-gray-400'
                }`}
                onClick={() => setWalletActionMode(2)}
              >
                Swap
              </button>
            </div>
            {(() => {
              switch (walletActionMode) {
                case 0:
                  return <SendView />;
                case 1:
                  return <ReceiveView />;
                case 2:
                  return <SwapView />;
              }
            })()}
          </div>
          <div className="flex w-full space-y-3 py-7 px-5 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
            <img src="/images/bankNotes.png" alt="bank notes" />
            <p className="text-center font-semibold">
              Make Deposits and Withdrawals on your assets with ease on Trovo
              Wallet
            </p>
            <div className="w-3/4">
              <Button
                label="Deposit/Withdraw"
                onclick={() => {
                  /* */
                }}
              />
            </div>
          </div>
          <div />
          <div />
          <div />
          <div />
        </div>
      </div>
    </div>
  );

  function SendView() {
    return (
      <div className="w-full space-y-5">
        <div className="space-y-3 w-full">
          <p>Select wallet</p>
          <Dropdown
            label={activeWallet.alias}
            options={[
              ...wallets.map((w) => ({
                text: w.alias,
                value: w.publicKey,
              })),
            ]}
            onSelect={() => {
              //
            }}
          />
        </div>
        <div className="space-y-3 w-full">
          <TextInput
            inputType="text"
            label="Send to"
            onInputChange={() => {}}
          />
        </div>
        <div className="w-full space-y-2">
          <div className="flex items-center space-x-3 space-y-3 w-full">
            <div className="w-4/6">
              <TextInput
                inputType="text"
                label="Token quantity/Amount"
                onInputChange={() => {}}
              />
            </div>
            <div className="w-2/6 self-end">
              <Dropdown
                label={
                  activeWallet.claimedAssets[0].assetCode === ''
                    ? 'XBN'
                    : activeWallet.claimedAssets[0].assetCode
                }
                options={[
                  ...activeWallet.claimedAssets.map((w) => ({
                    text: w.assetCode === '' ? 'XBN' : w.assetCode,
                    value: w.assetCode,
                  })),
                ]}
                onSelect={() => {
                  //
                }}
              />
            </div>
          </div>
          <div className="w-full flex justify-between items-center">
            <span>0.00000 XBN</span>
            <span>0.2343567 XBN</span>
          </div>
        </div>
        <div className="space-y-3 w-full">
          <TextInput
            inputType="text"
            label="Add memo (optional)"
            onInputChange={() => {}}
          />
          <div className="w-full flex justify-end items-center">
            <span>0/28</span>
          </div>
        </div>
        <Button label="Proceed" onclick={() => {}} />
      </div>
    );
  }

  function ReceiveView() {
    return (
      <div className="w-full space-y-5">
        <div className="space-y-3 w-full">
          <p>Receiving wallet</p>
          <Dropdown
            label={activeWallet.alias}
            options={[
              ...wallets.map((w) => ({
                text: w.alias,
                value: w.publicKey,
              })),
            ]}
            onSelect={() => {
              //
            }}
          />
        </div>
        <div className="w-full space-y-2">
          <div className="flex items-center space-x-3 space-y-3 w-full">
            <div className="w-4/6">
              <TextInput
                inputType="text"
                label="Token quantity/Amount"
                onInputChange={() => {}}
              />
            </div>
            <div className="w-2/6 self-end">
              <Dropdown
                label={
                  activeWallet.claimedAssets[0].assetCode === ''
                    ? 'XBN'
                    : activeWallet.claimedAssets[0].assetCode
                }
                options={[
                  ...activeWallet.claimedAssets.map((w) => ({
                    text: w.assetCode === '' ? 'XBN' : w.assetCode,
                    value: w.assetCode,
                  })),
                ]}
                onSelect={() => {
                  //
                }}
              />
            </div>
          </div>
          <div className="w-full flex justify-between items-center">
            <span>0.00000 XBN</span>
            <span>0.2343567 XBN</span>
          </div>
        </div>
        <div className="space-y-3 w-full">
          <TextInput
            inputType="text"
            label="Add memo (optional)"
            onInputChange={() => {}}
          />
          <div className="w-full flex justify-end items-center">
            <span>0/28</span>
          </div>
        </div>
        <Button label="Proceed" onclick={() => {}} />
      </div>
    );
  }

  function SwapView() {
    return (
      <div className="w-full space-y-5">
        <div className="space-y-3 w-full">
          <p>Select wallet</p>
          <Dropdown
            label={activeWallet.alias}
            options={[
              ...wallets.map((w) => ({
                text: w.alias,
                value: w.publicKey,
              })),
            ]}
            onSelect={() => {
              //
            }}
          />
        </div>
        <div className="w-full space-y-2">
          <p>Swap from</p>
          <Dropdown
            label={
              activeWallet.claimedAssets[0].assetCode === ''
                ? 'XBN'
                : activeWallet.claimedAssets[0].assetCode
            }
            options={[
              ...activeWallet.claimedAssets.map((w) => ({
                text: w.assetCode === '' ? 'XBN' : w.assetCode,
                value: w.assetCode,
              })),
            ]}
            onSelect={() => {
              //
            }}
          />
        </div>
        <div className="w-full space-y-2">
          <p>Swap to</p>
          <Dropdown
            label={
              activeWallet.claimedAssets[0].assetCode === ''
                ? 'XBN'
                : activeWallet.claimedAssets[0].assetCode
            }
            options={[
              ...activeWallet.claimedAssets.map((w) => ({
                text: w.assetCode === '' ? 'XBN' : w.assetCode,
                value: w.assetCode,
              })),
            ]}
            onSelect={() => {
              //
            }}
          />
        </div>
        <div className="w-full space-y-2">
          <div className="flex items-center space-x-3 space-y-3 w-full">
            <div className="w-4/6">
              <TextInput
                inputType="text"
                label="Token quantity/Amount"
                onInputChange={() => {}}
              />
            </div>
            <div className="w-2/6 self-end">
              <Dropdown
                label={
                  activeWallet.claimedAssets[0].assetCode === ''
                    ? 'XBN'
                    : activeWallet.claimedAssets[0].assetCode
                }
                options={[
                  ...activeWallet.claimedAssets.map((w) => ({
                    text: w.assetCode === '' ? 'XBN' : w.assetCode,
                    value: w.assetCode,
                  })),
                ]}
                onSelect={() => {
                  //
                }}
              />
            </div>
          </div>
          <div className="w-full flex justify-between items-center">
            <span>0.00000 XBN</span>
            <span>0.2343567 XBN</span>
          </div>
        </div>
        <Button label="Proceed" onclick={() => {}} />
      </div>
    );
  }
}
