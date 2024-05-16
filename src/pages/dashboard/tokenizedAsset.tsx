import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import ButtonSecondary from '../../components/buttonSecondary';
import Button from '../../components/button';
import Dropdown from '../../components/dropdown';
import TextInput from '../../components/textInput';
import Modal from '../../components/modal';
import Header from '../../components/header';

type GridItemProps = {
  title: string;
  value: string;
  value2: string;
};

function GridItem({ title, value, value2 }: GridItemProps) {
  return (
    <div className="flex h-full w-full py-5 px-2 xl:px-5 xl:space-y-5 rounded-xl mb-10 bg-primary-100">
      <div className="flex  flex-col space-y-2 w-full">
        <p className="text-primary-400">{title}</p>
        <p className="font-montserratSemiBold text-lg md:text-xl">{value}</p>
        <p className="text-sm">{value2}</p>
      </div>
    </div>
  );
}

type AssetDetailItemProps = {
  title: string;
  value: string;
  value2?: string;
  value3?: string;
};

// eslint-disable-next-line
function AssetDetailItem({
  title,
  value,
  value2,
  value3,
}: AssetDetailItemProps) {
  return (
    <div className="rounded-2xl w-full bg-white px-4 py-3">
      <div className="flex  flex-col space-y-2 w-full">
        <p className="font-montserratSemiBold text-primary-400 md:text-lg">
          {title}
        </p>
        <p className="text-sm">{value}</p>
        <p className="text-sm">{value2}</p>
        <p className="text-sm">{value3}</p>
      </div>
    </div>
  );
}

AssetDetailItem.defaultProps = {
  value2: '',
  value3: '',
};

export default function TokenizedAsset() {
  const navigate = useNavigate();
  const [showSubscribeModal, setShowSubscribeModal] = useState(false);
  const [showSuccessModal, setShowSuccessModal] = useState(false);
  const [showBuyTokenModal, setShowBuyTokenModal] = useState(false);
  const [buyTokenStep, setBuyTokenStep] = useState(1);

  function renderBuyTokenSteps() {
    switch (buyTokenStep) {
      case 1:
        return (
          <div className="flex flex-col items-center w-full px-5 md:px-20 space-y-5 py-5 justify-center">
            <div className="w-full flex justify-center">
              <img src="/images/thinkingMan.png" alt="thinking man" />
            </div>
            <p>Do you want to buy AE1 Tokens with the following wallet?</p>
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
                label="Proceed to buy tokens"
                additionalClasses="font-montserratSemiBold"
                onclick={() => {
                  setBuyTokenStep(2);
                }}
              />
            </div>
            <div />
          </div>
        );
      case 2:
        return (
          <div className="flex flex-col w-full px-5 md:px-20 space-y-5 py-5 justify-center">
            <p className="font-montserratSemiBold text-lg">Buy Token</p>
            <p className="text-sm">Select Currency</p>
            <Dropdown
              label="CNGN"
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
            <div className="w-full">
              <TextInput
                inputType="text"
                label="Quantity of AE1"
                placeholder="5,000.0000"
                onInputChange={() => {
                  // console.log('input has changed', newValue);
                }}
              />
            </div>
            <div className="w-full">
              <TextInput
                inputType="text"
                label="Amount"
                placeholder="5,000.0000"
                trailingIcon="/images/trovIcon.png"
                trailingText="CNGN"
                onInputChange={() => {
                  // console.log('input has changed', newValue);
                }}
              />
            </div>
            <div className="w-full self-center md:w-2/4 space-y-3">
              <Button
                label="Buy"
                additionalClasses="font-montserratSemiBold"
                onclick={() => {
                  setBuyTokenStep(3);
                }}
              />
            </div>
            <div />
          </div>
        );
      case 3:
        return (
          <div className="flex flex-col items-center w-full px-5 md:px-20 space-y-5 py-5 justify-center">
            <p className="text-primary-800 w-full mb-5 text-lg text-center">
              You are buying
            </p>
            <div className="flex w-full py-3 px-2 xl:px-5 xl:space-y-5 rounded-xl flex-col items-center bg-primary-100">
              <div className="flex flex-col space-y-2 items-center w-full justify-between">
                <p className="font-montserratSemiBold">5,000 AE1</p>
                <p className="text-xs text-primary-800">
                  of Atlantis Estate 1 Asset
                </p>
              </div>
            </div>
            <p className="text-primary-800 w-full mb-5 text-center">Amount</p>
            <div className="flex w-full py-3 px-2 xl:px-10 xl:space-y-5 rounded-xl justify-center items-center bg-primary-100">
              <div className="flex flex-col space-y-2 items-center w-full justify-between">
                <p className="font-montserratSemiBold">500,000 CNGN</p>
                <p className="text-xs text-primary-800">$357.14</p>
              </div>
            </div>
            <p className="text-primary-800 w-full mb-5 text-center">Pay with</p>
            <div className="flex w-full py-5 px-2 xl:px-10 xl:space-y-5 rounded-xl mb-10 justify-center items-center bg-primary-100">
              <div className="flex space-x-2 items-center">
                <img src="/images/trovIcon.png" alt="" />
                <p className="text-xs text-primary-800">Efizee</p>
              </div>
            </div>
            <div className="w-full flex space-x-3 pt-5">
              <Button
                label="Confirm"
                additionalClasses="font-montserratSemiBold"
                onclick={() => {
                  setBuyTokenStep(4);
                }}
              />
              <ButtonSecondary
                label="Back"
                additionalClasses="font-montserratSemiBold"
                onclick={() => {
                  setBuyTokenStep(2);
                }}
              />
            </div>
            <div />
          </div>
        );
      default:
        return (
          <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
            <img src="/images/success.png" alt="success" />
            <div className="flex flex-col text-center space-y-5 items-center w-2/3 md:px-10 justify-center">
              <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
                Purchase Successful!
              </p>
            </div>
          </div>
        );
    }
  }

  return (
    <div className="flex text-primary-800 text-sm md:text-md flex-col space-y-5 p-3">
      <Header isHomeView />
      <div className="w-full px-3 md:px-5 md:pt-5 flex flex-col md:flex-row space-y-5 md:space-y-0 items-center justify-between">
        <div className="flex space-x-5 w-full md:w-auto items-center">
          <button type="button" onClick={() => navigate(-1)}>
            <img src="/images/arrowBack.png" alt="arrow back" />
          </button>
          <p className="font-montserratSemiBold text-lg xl:text-xl">
            Atlantis Estate 1 Asset
          </p>
        </div>
        <div className="flex w-full md:w-2/6 space-x-5">
          <ButtonSecondary
            label="Buy"
            additionalClasses="bg-primary-600 text-white font-montserratSemiBold"
            onclick={() => {
              setShowBuyTokenModal(true);
            }}
          />
          <ButtonSecondary
            label="Subscribe"
            additionalClasses="font-montserratSemiBold"
            onclick={() => {
              setShowSubscribeModal(true);
            }}
          />
        </div>
      </div>
      <div className="w-full px-3 md:px-5 md:pt-5 flex flex-col md:flex-row md:space-x-3 space-y-5 md:space-y-0">
        <div>
          <img src="/images/realEstate.png" alt="asset logo" />
        </div>
        <div className="w-full space-y-4">
          <p className="flex items-center space-x-2">
            <img src="/images/location.png" alt="location" />
            <span>No 10, Maitama, Abuja, Nigeria</span>
          </p>
          <p className="flex items-center space-x-2">
            <img src="/images/tickOutline.png" alt="check" />
            <span>Available in your Country</span>
          </p>
          <p className="flex items-center space-x-2">
            <img src="/images/uiLink.png" alt="website" />
            <span>www.atlantis.com</span>
          </p>
          <p>
            Atlantis Estate 1 tokens - AE1 are fractional tokens that represent
            part ownership (via investment) of our real estate development
            project at Atlantis Estate, Lekki, Lagos, Nigeria. Subscribe to this
            token to earn rental income monthly pushed to your Trovo App. Also
            earn benefit from asset appreciation with the ability to sell or buy
            any portion of your tokens at anytime.
          </p>
        </div>
      </div>
      <div className="px-3 md:px-5 pt-5 w-full grid grid-cols-2 md:grid-cols-3 gap-4">
        <GridItem title="Total Supply" value="10,000,000" value2="" />
        <GridItem title="Total Subscribed" value="800" value2="" />
        <GridItem title="Price per Asset" value="100 CNGN" value2="$2.20" />
        <GridItem title="Funding Currency" value="CNGN" value2="" />
        <GridItem title="Subscription Amount" value="0 CNGN" value2="" />
        <GridItem title="Actual Amount Bought" value="0 CNGN" value2="" />
      </div>
      <div className="px-3 md:px-5">
        <div className="flex w-full space-y-3 py-7 px-2 xl:px-5 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
          <div className="flex items-center w-full justify-between">
            <p className="font-montserratSemiBold text-lg md:text-xl">
              Other Details
            </p>
          </div>
          <AssetDetailItem title="Asset Code" value="ATLANTIS 1" />
          <AssetDetailItem title="Sector" value="Real Estate" />
          <AssetDetailItem title="Sub-Sector" value="Residential Properties" />
          <AssetDetailItem title="Type" value="Single-Family Homes" />
          <AssetDetailItem title="Issuer" value="Atlantis Developers" />
          <AssetDetailItem title="Issuer Website" value="www.atlantis.com" />
          <AssetDetailItem
            title="Sales Window"
            value="12/11/2023 - 15/12/2023"
          />
          <AssetDetailItem title="Cap Quantity" value="500,000 CNGN" />
          <AssetDetailItem
            title="Cap Duration"
            value="12/11/2023 - 15/11/2023"
          />
          <AssetDetailItem title="Proceed Payout Cycle" value="Monthly" />
          <AssetDetailItem title="Countries Exempted" value="See list" />
          <AssetDetailItem
            title="Proof of Existence"
            value="C of O"
            value2="Survey Plan"
            value3="Governor’s Consent"
          />
        </div>
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
            Please enter the amount you would like spend on AE1 tokens when
            primary sale starts
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
            Are you sure you want to subscribe to the Atlantis 1 Asset with the
            following wallet?
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
      {/* Public key copied modal */}
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
      <Modal
        showModal={showBuyTokenModal}
        onClose={() => {
          setShowBuyTokenModal(false);
          setBuyTokenStep(1);
        }}
      >
        {renderBuyTokenSteps()}
      </Modal>
    </div>
  );
}
