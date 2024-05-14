import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import TextInput from '../../components/textInput';
import WidgetCard from '../../components/widgetCard';
import Tabs from '../../components/tabs';
import AssetListItem from '../../components/assetListItem';
import TransactionItem from '../../components/transactionItem';
import Button from '../../components/button';
import ButtonSecondary from '../../components/buttonSecondary';
import Header from '../../components/header';
import Modal from '../../components/modal';
import { useSelector } from 'react-redux';
import { RootState } from '../../store/reduxStore';

function ActivateWalletView() {
  const [showCopiedModal, setShowCopiedModal] = useState(false);
  const [showRequestXBNModal, setShowRequestXBNModal] = useState(false);
  const [confirmRequestXBNModal, setConfirmRequestXBNModal] = useState(false);

  return (
    <div className="w-full">
      <p className="text-primary-800 text-center text-lg md:text-xl font-montserratSemiBold">
        Wallet Activation
      </p>
      <div className="md:px-5 w-full">
        <div className="w-full flex flex-col items-center justify-center space-y-5 md:mb-10">
          <img
            className="self-center"
            src="/images/activateWallet.png"
            alt=""
          />
          <p className="font-montserratSemiBold text-md text-center">
            Your Wallet is ready!
          </p>
          <div className="space-y-2">
            <p className="text-md text-center">
              But you cannot use it for any transaction just yet until it is
              activated with atleast 10 Bantu tokens (XBN)
            </p>
            <p className="text-md text-center">
              You can get Bantu tokens (XBN) for your wallet using either of the
              3 easy ways displayed below
            </p>
          </div>
          <div className="md:w-2/4 px-3 pb-5 space-y-3">
            <Button
              label="Request XBN from Trovo User"
              additionalClasses="font-montserratSemiBold"
              onclick={() => {
                setShowRequestXBNModal(true);
              }}
            />
            <ButtonSecondary
              label="Send XBN to your wallet"
              additionalClasses="bg-primary-600 text-white font-montserratSemiBold ring-primary-600"
              onclick={() => {
                setShowCopiedModal(true);
              }}
            />
            <ButtonSecondary
              label="Buy XBN on TrovoP2P"
              additionalClasses="font-montserratSemiBold"
              onclick={() => {
                // navigate('/import');
              }}
            />
          </div>
        </div>
      </div>
      {/* Public key copied modal */}
      <Modal
        showModal={showCopiedModal}
        onClose={() => {
          setShowCopiedModal(false);
        }}
      >
        <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
          <img src="/images/success.png" alt="success" />
          <div className="flex flex-col text-center space-y-5 items-center w-2/3 md:px-10 justify-center">
            <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
              Public Key Copied Successfully!
            </p>
          </div>
        </div>
      </Modal>
      {/* Request XBN modal */}
      <Modal
        showModal={showRequestXBNModal}
        onClose={() => {
          setShowRequestXBNModal(false);
        }}
      >
        <div className="flex flex-col items-center w-full px-5 md:px-20 space-y-5 py-5 justify-center">
          <p className="text-primary-800 w-full mb-5 text-lg md:text-xl font-montserratSemiBold">
            Request XBN
          </p>
          <div className="w-full">
            <TextInput
              inputType="text"
              label="Receiving Wallet"
              placeholder="Efizee"
              onInputChange={() => {
                // console.log('input has changed', newValue);
              }}
            />
          </div>
          <div className="w-full">
            <TextInput
              inputType="text"
              label="Amount"
              placeholder="100"
              onInputChange={() => {
                // console.log('input has changed', newValue);
              }}
            />
            <div className="flex justify-between text-primary-700 mt-2">
              <span>0.0000 XBN</span>
              <span>2,300,320.3214 XBN</span>
            </div>
          </div>
          <div className="w-full">
            <TextInput
              inputType="text"
              label="Add memo (optional)"
              placeholder=""
              onInputChange={() => {
                // console.log('input has changed', newValue);
              }}
            />
            <div className="flex justify-end text-primary-700 mt-2">
              <span>0/28</span>
            </div>
          </div>
          <div className="w-full md:w-2/4 space-y-3">
            <Button
              label="Proceed"
              additionalClasses="font-montserratSemiBold"
              onclick={() => {
                setConfirmRequestXBNModal(true);
                setShowRequestXBNModal(false);
              }}
            />
          </div>
          <div />
        </div>
      </Modal>
      <Modal
        showModal={confirmRequestXBNModal}
        onClose={() => {
          setConfirmRequestXBNModal(false);
        }}
      >
        <div className="flex flex-col items-center w-full px-5 md:px-20 space-y-5 py-5 justify-center">
          <p className="text-primary-800 w-full mb-5 text-lg text-center">
            You are requesting for
          </p>
          <div className="flex w-full py-3 px-2 xl:px-5 xl:space-y-5 rounded-xl flex-col items-center bg-primary-100">
            <div className="flex flex-col space-y-2 items-center w-full justify-between">
              <p className="font-montserratSemiBold">45.0000 xbn</p>
              <p className="text-xs text-primary-400">~3400 NGN</p>
            </div>
          </div>
          <p className="text-primary-800 w-full mb-5 text-center">
            Receiving Wallet
          </p>
          <div className="flex w-full py-3 px-2 xl:px-10 xl:space-y-5 rounded-xl justify-center items-center bg-primary-100">
            <div className="flex flex-col space-y-2 items-center w-full justify-between">
              <p className="font-montserratSemiBold">Obi</p>
              <p className="text-xs text-primary-400">
                0x71C7656EC7ab88b098defB751B7401B5f6d8976F
              </p>
            </div>
            <button type="button">
              <img src="/images/copy.png" alt="copy" />
            </button>
          </div>
          <p className="text-primary-800 w-full mb-5 text-center">
            Description/Memo
          </p>
          <div className="flex w-full py-5 px-2 xl:px-10 xl:space-y-5 rounded-xl mb-10 justify-center items-center bg-primary-100">
            <div className="flex flex-col space-y-2 items-center w-full justify-between">
              <p className="text-xs text-primary-400">test</p>
            </div>
          </div>
          <div className="flex w-full py-5 px-2 xl:px-10 xl:space-y-5 rounded-xl mb-10 justify-center items-center bg-primary-100">
            <div className="flex flex-col space-y-2 items-center w-full justify-between">
              <img className="w-2/5" src="/images/qrcode.png" alt="QR Code" />
            </div>
          </div>
          <div className="w-full flex space-x-3 pt-5">
            <Button
              label="Share"
              additionalClasses="font-montserratSemiBold"
              onclick={() => {
                setConfirmRequestXBNModal(false);
              }}
            />
            <ButtonSecondary
              label="Back"
              additionalClasses="font-montserratSemiBold"
              onclick={() => {
                setShowRequestXBNModal(true);
                setConfirmRequestXBNModal(false);
              }}
            />
          </div>
          <div />
        </div>
      </Modal>
    </div>
  );
}

function WalletView({ hasAssets }: { hasAssets: boolean }) {
  const navigate = useNavigate();

  if (!hasAssets) {
    return (
      <div className="w-full">
        <p className="text-primary-800 text-center text-lg md:text-xl font-montserratSemiBold">
          Asset Offering
        </p>
        <div className="md:px-5 w-full">
          <div className="w-full flex flex-col justify-center space-y-5 md:mb-10">
            <img className="self-center" src="/images/empty.png" alt="" />
            <p className="font-montserratSemiBold text-md text-center">
              No Assets Yet
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="w-full">
      <p className="text-primary-800 text-center text-lg md:text-xl font-montserratSemiBold">
        Asset Offering
      </p>
      <div className="md:px-5 w-full">
        <Tabs tabList={['Primary Listing', 'Secondary Listing']}>
          <div className="tab-1 flex flex-col space-y-3">
            <AssetListItem
              image="/images/avatar.png"
              assetName="Atlantis 1"
              assetClass="Real Estate"
              isSubscribed
              onclick={() => {
                navigate('/dashboard/tokenized-asset');
              }}
            />
            <AssetListItem
              image="/images/avatar.png"
              assetName="Orchard Estate"
              assetClass="Real Estate"
              onclick={() => {
                /* */
              }}
            />
            <AssetListItem
              image="/images/avatar.png"
              assetName="Atlantis 1"
              assetClass="Real Estate"
              isSubscribed
              onclick={() => {
                /* */
              }}
            />
            <AssetListItem
              image="/images/avatar.png"
              assetName="Orchard Estate"
              assetClass="Real Estate"
              onclick={() => {
                /* */
              }}
            />
            <AssetListItem
              image="/images/avatar.png"
              assetName="Atlantis 1"
              assetClass="Real Estate"
              isSubscribed
              onclick={() => {
                /* */
              }}
            />
            <AssetListItem
              image="/images/avatar.png"
              assetName="Orchard Estate"
              assetClass="Real Estate"
              onclick={() => {
                /* */
              }}
            />
            <AssetListItem
              image="/images/avatar.png"
              assetName="Atlantis 1"
              assetClass="Real Estate"
              isSubscribed
              onclick={() => {
                /* */
              }}
            />
            <AssetListItem
              image="/images/avatar.png"
              assetName="Orchard Estate"
              assetClass="Real Estate"
              onclick={() => {
                /* */
              }}
            />
            <AssetListItem
              image="/images/avatar.png"
              assetName="Atlantis 1"
              assetClass="Real Estate"
              isSubscribed
              onclick={() => {
                /* */
              }}
            />
            <AssetListItem
              image="/images/avatar.png"
              assetName="Orchard Estate"
              assetClass="Real Estate"
              onclick={() => {
                /* */
              }}
            />
          </div>
          <div className="tab-2">
            <div className="flex flex-col space-y-3">
              <AssetListItem
                image="/images/avatar.png"
                assetName="Metro Railings"
                assetClass="Infrastructure"
                isSubscribed
                onclick={() => {
                  /* */
                }}
              />
              <AssetListItem
                image="/images/avatar.png"
                assetName="Mount"
                assetClass="Infrastructure"
                onclick={() => {
                  /* */
                }}
              />
            </div>
          </div>
        </Tabs>
      </div>
    </div>
  );
}

export default function Home() {
  const [hasAssets, setHasAssets] = useState(false);
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const primaryWallet = appUser.userWallets.find((w) => w.primaryWallet);

  // const [isActivated] = useState(
  //   Number(
  //     primaryWallet?.claimedAssets.find((a) => !a.assetCode && !a.assetIssuer)
  //       ?.amount,
  //   ) !== 0,
  // );

  // console.log('primary wallet ', appUser);

  return (
    <div className="flex text-primary-800 text-sm md:text-md flex-col space-y-5 p-3">
      <Header isHomeView />
      <div className="flex md:h-full w-full items-center justify-center">
        <div className="h-full w-full md:p-3">
          {!appUser.hasSecurityQuestions && (
            <div
              className="bg-trovored-light flex justify-between mb-2 rounded-md ring-1 ring-trovored-primary
            text-trovored-primary px-4 py-2 font-semibold"
            >
              <div className="flex space-x-5 items-center">
                <Link className="flex space-x-5 items-center" to={'/welcome'}>
                  <img src="/images/alert.png" alt="" />
                  <p className="text-xs md:text-md">
                    You have not setup security questions yet. Tap to setup
                    security questions
                  </p>
                </Link>
              </div>
            </div>
          )}
          <div className="md:hidden items-center px-5 mb-5 justify-center flex">
            <WidgetCard />
          </div>
          <div className="flex space-y-3 mb-20 md:pt-5 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
            <div className="hidden w-full items-center px-5 justify-center md:flex">
              <WidgetCard />
            </div>
            {/* {isActivated ? (
              <WalletView hasAssets={hasAssets} />
            ) : (
              <ActivateWalletView />
            )} */}
          </div>
        </div>
        <div className="hidden md:block w-3/6 flex flex-col space-y-5 h-full py-3 lg:px-3">
          <div className="flex space-y-3 rounded-lg py-5 flex-col items-center bg-primary-100">
            <p>Current Plan</p>
            <div className="flex space-x-1">
              <img src="/images/gold.png" alt="gold" />
              <p className="font-bold">Gold Patron (Monthly)</p>
            </div>
            <div className="w-full px-3 lg:px-10">
              <div
                className="max-w-xl w-full flex flex-col xl:flex-row text-white font-matahariRegular text-xl xl:text-2xl
              rounded-2xl bg-primary-800 items-center justify-center space-y-3 xl:space-y-0 xl:items-start xl:justify-between xl:space-x-3 mt-5 xl:mb-5 py-5 xl:py-10 xl:px-5"
              >
                <img
                  className="w-10"
                  src="/images/platinum.png"
                  alt="platinum"
                />
                <div className="flex flex-col space-y-2 text-center xl:text-start">
                  <p className="text-lg">Try the Diamond Plan</p>
                  <p className="text-sm">Upgrade to unlock more features</p>
                </div>
                <img className="w-10" src="/images/go.png" alt="platinum" />
              </div>
            </div>
          </div>
          <div className="flex w-full space-y-3 py-7 px-2 xl:px-5 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
            <div className="flex items-center w-full justify-between">
              <p className="font-semibold">Recent Transactions</p>
              <p className="test-xs text-primary-400">view all</p>
            </div>
            <TransactionItem
              addressOrUsername="Tom"
              transactionType={1}
              amount="3,000.000"
              assetCode="XBN"
              date={new Date().toLocaleString()}
              onclick={() => {
                setHasAssets(!hasAssets);
              }}
            />
            <TransactionItem
              addressOrUsername="GH3Y8...ZX89O"
              transactionType={0}
              amount="0.5000"
              assetCode="TROV"
              date={new Date().toLocaleString()}
              onclick={() => {
                /* */
              }}
            />
            <TransactionItem
              addressOrUsername="Tom"
              transactionType={1}
              amount="3,000.000"
              assetCode="XBN"
              date={new Date().toLocaleString()}
              onclick={() => {
                /* */
              }}
            />
            <TransactionItem
              addressOrUsername="GH3Y8...ZX89O"
              transactionType={0}
              amount="0.5000"
              assetCode="TROV"
              date={new Date().toLocaleString()}
              onclick={() => {
                /* */
              }}
            />
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
}
