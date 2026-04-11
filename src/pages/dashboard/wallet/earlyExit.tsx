import { useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import Header from '../../../components/header';
import { Asset } from '../../../types/asset';
import { useSelector } from 'react-redux';
import { RootState } from '../../../store/reduxStore';
import { Wallet } from '../../../types/wallet';
import TextInput from '../../../components/textInput';
import Dropdown from '../../../components/dropdown';
import Button from '../../../components/button';
import Modal from '../../../components/modal';
import ButtonSecondary from '../../../components/buttonSecondary';

export default function EarlyExitView() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const appUser = useSelector((state: RootState) => state.auth.user!);

  const [password, setPassword] = useState('');
  const [passwordErr, setPasswordErr] = useState('');
  const [showSummary, setShowSummary] = useState(false);
  const [showSuccessModal, setShowSuccessModal] = useState(false);
  const [agreesToTerms, setAgreesToTerms] = useState(false);
  const [wallet] = useState<Wallet>(
    appUser.userWallets.find((w) => w.publicKey == searchParams.get('wallet'))!,
  );
  const [asset] = useState<Asset>(
    wallet?.claimedAssets.find(
      (a) =>
        a.assetCode == searchParams.get('assetCode') &&
        a.assetIssuer == searchParams?.get('assetIssuer'),
    )!,
  );
  const [errorObj, setErrorObj] = useState({
    terms: '',
  });

  return (
    <div className="flex text-primary-800 text-sm md:text-md flex-col space-y-5 p-3">
      <Header />
      <div className="flex flex-col items-center justify-center w-full">
        <div className="flex px-5 mb-10 space-x-5 w-full items-center">
          <button type="button" onClick={() => navigate(-1)}>
            <img src="/images/arrowBack.png" alt="arrow back" />
          </button>
          <p className="font-montserratSemiBold text-lg xl:text-xl">
            Request Early Exit
          </p>
        </div>
        <div className="flex md:h-screen w-4/6 items-start self-center">
          <div className="w-3/5 md:p-3">
            <div className="flex md:pt-5 p-5 rounded-xl flex-col items-center bg-primary-100 max-w-5xl min-h-full">
              <div className="w-full space-y-6 px-3 md:px-5 md:py-5 flex flex-col justify-center space-y-5 md:space-y-3">
                <div className="flex flex justify-center items-center space-x-1 mb-10">
                  <img
                    src={asset?.imageUrl}
                    alt="asset logo"
                    className="h-10 w-10 rounded-full"
                  />
                  <p className="font-montserratSemiBold text-lg xl:text-xl">
                    {asset?.assetCode}
                  </p>
                </div>
                <div className="flex flex-col justify-between w-full space-y-12">
                  <div id="dividends" className="w-full space-y-5">
                    <p className="font-montserratSemiBold text-lg xl:text-xl">
                      Early Exit
                    </p>
                    <p className="flex w-full justify-between">
                      <span>Eligible for early exit</span>
                      <span className="font-montserratSemiBold text-md">
                        Yes
                      </span>
                    </p>
                    <p className="flex justify-between">
                      <span>Currenty buyback NAV</span>
                      <span className="font-montserratSemiBold text-md">
                        N980/token
                      </span>
                    </p>
                    <p className="flex justify-between">
                      <span>Exit charges or penalty</span>
                      <span className="font-montserratSemiBold text-md">
                        2% early exit fee
                      </span>
                    </p>
                    <p className="flex justify-between">
                      <span>Estimated payout (NET)</span>
                      <span className="font-montserratSemiBold text-md">
                        N980/token
                      </span>
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div className="h-full w-full md:p-3">
            {showSummary ? (
              <div className="flex md:pt-5 p-5 rounded-xl flex-col items-center bg-primary-100 max-w-5xl">
                <div className="w-full space-y-6 px-3 md:px-5 md:py-5 flex flex-col justify-center space-y-5 md:space-y-3">
                  <div className="flex flex-col justify-between w-full space-y-12">
                    <div id="summary" className="w-full space-y-5">
                      <p className="flex w-full justify-between">
                        <span>Token Quantity to Exit</span>
                        <span className="font-montserratSemiBold text-md">
                          1,000 TROV
                        </span>
                      </p>
                      <p className="flex justify-between">
                        <span>Currenty NAV/Token</span>
                        <span className="font-montserratSemiBold text-md">
                          N980
                        </span>
                      </p>
                      <p className="flex justify-between">
                        <span>Early Exit Penalty + Fees</span>
                        <span className="font-montserratSemiBold text-md">
                          2%
                        </span>
                      </p>
                      <p className="flex justify-between">
                        <span>Estimated Payout Amount</span>
                        <span className="font-montserratSemiBold text-md">
                          N980,000
                        </span>
                      </p>
                      <p className="flex justify-between">
                        <span>Payout Currency</span>
                        <span className="font-montserratSemiBold text-md">
                          NGN
                        </span>
                      </p>
                      <p className="flex justify-between">
                        <span>Wallet ID</span>
                        <span className="font-montserratSemiBold text-md">
                          Johnnydoe
                        </span>
                      </p>
                      <p className="flex justify-between">
                        <span>Settlement Time Estimate</span>
                        <span className="font-montserratSemiBold text-md">
                          48 hours
                        </span>
                      </p>
                      <div className="space-y-1">
                        <div className="flex space-x-3">
                          <input
                            type="checkbox"
                            name="import"
                            defaultChecked={agreesToTerms}
                            onChange={() => {
                              setAgreesToTerms(!agreesToTerms);
                            }}
                          />
                          <p className="text-gray-500">
                            I understand that early exit may involve penalties
                            and affects my future returns
                          </p>
                        </div>
                        {errorObj.terms && (
                          <p className="text-red-500 text-sm">
                            {errorObj.terms}
                          </p>
                        )}
                      </div>
                      <div className="w-full space-y-1">
                        <TextInput
                          label=""
                          leadingIcon="/images/lock.png"
                          inputType="password"
                          placeholder="Enter answer"
                          onInputChange={(newValue: string) => {
                            setPassword(newValue);
                          }}
                        />
                        {passwordErr && (
                          <p className="text-red-500 text-sm">{passwordErr}</p>
                        )}
                      </div>
                      <div className="w-full space-y-3">
                        <Button
                          label="Authorize with password"
                          additionalClasses="font-montserratSemiBold"
                          onclick={async () => {
                            setShowSuccessModal(true);
                          }}
                        />
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            ) : (
              <div className="flex md:pt-5 p-5 rounded-lg flex-col items-start max-w-5xl min-h-full">
                <form
                  id="early-exit"
                  onSubmit={() => {
                    setShowSummary(true);
                  }}
                  className="w-full space-y-5"
                >
                  <TextInput
                    inputType="text"
                    label="No of tokens to exit"
                    onInputChange={(value) => {}}
                    error={''}
                  />
                  <Dropdown
                    label="Payment method"
                    options={[
                      { text: 'Wallet', value: 'wallet' },
                      { text: 'Bank account', value: 'bank' },
                    ]}
                    onSelect={(asset) => {
                      console.log(asset);
                    }}
                  />
                  <TextInput
                    inputType="text"
                    label="Wallet username/address"
                    onInputChange={(value) => {}}
                    error={''}
                  />
                  <Button type="submit" label="Proceed" onclick={() => {}} />
                </form>
              </div>
            )}
          </div>
        </div>
      </div>
      <Modal
        showModal={showSuccessModal}
        onClose={() => {
          setShowSuccessModal(false);
        }}
      >
        <div className="flex flex-col items-center w-full px-5 md:px-20 space-y-5 py-5 justify-center">
          <div className="flex flex-col text-center space-y-5 items-center w-2/3 md:px-10 justify-center">
            <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
              Your transaction was successful!
            </p>
          </div>
          <img src="/images/success.png" alt="success" />
          <div className="w-full space-y-3">
            <a
              // href={`/#/send-asset-receipt?q=${receiptQuery}`} // since we use hash router
              target="_blank"
              className="bg-primary-800 block rounded-lg w-full text-center text-white h-12 py-4 px-5"
            >
              Generate receipt
            </a>
            <ButtonSecondary
              label="Close"
              additionalClasses="font-montserratSemiBold"
              onclick={async () => {
                setShowSuccessModal(false);
              }}
            />
          </div>
          <div />
        </div>
      </Modal>
    </div>
  );
}
