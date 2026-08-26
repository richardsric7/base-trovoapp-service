import { useEffect, useState } from 'react';
import { TokenizedAsset } from '../types/tokenizedAsset';
import Button from './button';
import LinkButton from './linkButton';
import Modal from './modal';
import TextInput from './textInput';
import { showNotification, toggleLoader } from '../utils/showToaster';
import { ErrorResponse } from '../store/api/baseapi/axiosBaseQuery';
import { useSelector } from 'react-redux';
import { RootState } from '../store/reduxStore';
import { Encryptor } from '../utils/encryptor';
import {
  useBuyTokenizedAssetsMutation,
  useSubscribeTokenizedAssetsMutation,
} from '../store/api/walletApis';
import Dropdown from './dropdown';
import { formatToDecimal } from '../utils/utilities';
import { signBase64Txn } from '../utils/trovoSDK';

type Props = {
  show: boolean;
  onClose: () => void;
  asset?: TokenizedAsset;
  onSuccess?: () => void;
};

export function ExpressInterestModal({
  show,
  onClose,
  onSuccess,
  asset,
}: Props) {
  const [amount, setAmount] = useState(0);
  const [error, setError] = useState('');
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const primaryWallet = appUser.userWallets.find((w) => w.primaryWallet)!;
  const [activeWallet] = useState(primaryWallet);
  const [secretKey, setSecretKey] = useState('');
  const [subscribeToken] = useSubscribeTokenizedAssetsMutation();

  useEffect(() => {
    (async () => {
      const encryptor = new Encryptor();
      setSecretKey(await encryptor.getSecretKey(appUser));
    })();
  }, []);
  return (
    <Modal showModal={show} onClose={onClose}>
      <div className="flex flex-col items-center w-full px-5 md:px-20 space-y-5 py-5 justify-center">
        <div className="w-full flex justify-center">
          <img src="/images/thinkingMan.png" alt="thinking man" />
        </div>
        {(asset?.expressedInterestAmount ?? 0 > 0) ? (
          <p>
            You have already indicated to invest{' '}
            {asset?.expressedInterestAmount ?? 0} CNGN on {asset!.assetCode}{' '}
            tokens when primary sale starts. Do you want to update it?
          </p>
        ) : (
          <p>
            Please enter the amount you want to invest in {asset!.assetCode}{' '}
            token when primary sale starts
          </p>
        )}
        <div className="w-full space-y-2">
          <TextInput
            inputType="text"
            label=""
            placeholder="100"
            trailingText="CNGN"
            onInputChange={(value) => {
              try {
                setAmount(Number(value));
              } catch (error) {
                setError(error?.toString() ?? '');
              }
            }}
          />
          <p className="text-red-500">{error}</p>
        </div>
        <div className="w-full md:w-2/4 space-y-3">
          <Button
            label="Express Interest"
            additionalClasses="font-montserratSemiBold"
            onclick={async () => {
              console.log('amount ===> ', amount);
              if (amount <= 0 || isNaN(amount)) {
                setError('Please enter a valid amount.');
                return;
              }
              const payload = {
                signer: activeWallet.signer,
                publicKey: activeWallet.publicKey,
                secretKey: secretKey,
                body: {
                  isSharedWallet: activeWallet.sharedAccessEnabled,
                  amount: amount,
                  assetId: asset?.id,
                },
              };

              toggleLoader();
              const res = await subscribeToken(payload);
              console.log('response', res);

              if ('data' in res) {
                console.log('response', res);
                onSuccess!();
                setError('');
              } else if ('error' in res) {
                const errorResponse = res.error as ErrorResponse;
                showNotification(
                  'error',
                  errorResponse.data.message ??
                    'Something went wrong. Please try again.',
                );
              }
              toggleLoader();
            }}
          />
        </div>
        <div />
      </div>
    </Modal>
  );
}

export function P2pComingSoonModal({ show, onClose }: Props) {
  return (
    <Modal showModal={show} onClose={onClose}>
      <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
        <img src="/images/buy_with_fiat.png" alt="success" />
        <div className="flex flex-col text-center space-y-5 items-center w-2/3 md:px-10 justify-center">
          <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
            TrovoP2P will be launching soon.
          </p>
        </div>
      </div>
    </Modal>
  );
}

export function KycUnverifiedModal({ show, onClose }: Props) {
  return (
    <Modal showModal={show} onClose={onClose}>
      <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
        <img src="/images/buy_with_fiat.png" alt="success" />
        <div className="flex flex-col text-center space-y-5 items-center w-2/3 md:px-10 justify-center">
          <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
            You must have your KYC verification done before you can proceed with
            this action.
          </p>
        </div>
        <div className="w-3/6">
          <LinkButton
            label="Verify Account"
            link="#"
            additionalClasses="bg-primary-600 text-white font-montserratSemiBold"
          />
        </div>
      </div>
    </Modal>
  );
}
export function SoldOutModal({ show, onClose }: Props) {
  return (
    <Modal showModal={show} onClose={onClose}>
      <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
        <img src="/images/buy_with_fiat.png" alt="success" />
        <div className="flex flex-col text-center space-y-5 items-center w-2/3 md:px-10 justify-center">
          <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
            This token has already sold out.
          </p>
        </div>
      </div>
    </Modal>
  );
}

export function ShowExpressInterestSuccessModal({
  show,
  onClose,
  asset,
}: Props) {
  return (
    <Modal showModal={show} onClose={onClose}>
      <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
        <img src="/images/success.png" alt="success" />
        <div className="flex flex-col text-center space-y-5 items-center w-2/3 md:px-10 justify-center">
          <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
            You have successfully expressed interest on {asset!.assetCode}
          </p>
        </div>
      </div>
    </Modal>
  );
}

export function BuyTokenModal({ show, onClose, asset }: Props) {
  return (
    <Modal showModal={show} onClose={onClose}>
      {renderBuyTokenSteps(asset!, onClose)}
    </Modal>
  );
}

export function ShowVerificationDocumentsModal({
  show,
  onClose,
  asset,
}: Props) {
  return (
    <Modal showModal={show} onClose={onClose}>
      <div className="flex flex-col items-center w-full px-5 space-y-5 py-5 justify-center">
        <p className="text-primary-800 w-full text-lg font-montserratSemiBold">
          Asset Documents
        </p>
        <div className="flex w-full py-3 px-2 xl:px-5 xl:space-y-5 rounded-xl flex-col items-center bg-primary-100">
          {asset!.AssetTokenizationDocuments?.length > 0 ? (
            asset?.AssetTokenizationDocuments.map((item, index) => (
              <div
                key={index}
                className="flex flex-col space-y-2 w-full justify-between"
              >
                <a href={item.documentUrl} target="_blank">
                  <p className="font-montserratSemiBold underline">
                    {item.documentTitle}
                  </p>
                </a>
              </div>
            ))
          ) : (
            <></>
          )}
        </div>
        <div className="flex w-3/4 space-x-5">
          <Button
            label="Close"
            additionalClasses="bg-primary-600 text-white font-montserratSemiBold"
            onclick={onClose}
          />
        </div>
      </div>
    </Modal>
  );
}

function renderBuyTokenSteps(asset: TokenizedAsset, onClose: () => void) {
  const [buyTokenStep, setBuyTokenStep] = useState(1);
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const primaryWallet = appUser.userWallets.find((w) => w.primaryWallet)!;
  const [activeWallet, setActiveWallet] = useState(primaryWallet);
  const [amount, setAmount] = useState(0);
  const [quantity, setQuantity] = useState(0);
  const [error, setError] = useState('');
  const [secretKey, setSecretKey] = useState('');
  const [password, setPassword] = useState('');
  const [passwordErr, setPasswordErr] = useState('');
  const [buyToken] = useBuyTokenizedAssetsMutation();

  useEffect(() => {
    (async () => {
      const encryptor = new Encryptor();
      setSecretKey(await encryptor.getSecretKey(appUser));
    })();
  }, []);

  const isValidPassword = async () => {
    try {
      const encryptor = new Encryptor();
      const key = await encryptor.decryptData(
        appUser?.secretKeys[0],
        password,
        appUser?.primarySigner,
      );
      setSecretKey(key);
      return true;
    } catch (error: any) {
      return false;
    }
  };

  switch (buyTokenStep) {
    case 1:
      return (
        <div className="flex flex-col items-center w-full px-5 md:px-20 space-y-5 py-5 justify-center">
          <div className="w-full flex justify-center">
            <img src="/images/thinkingMan.png" alt="thinking man" />
          </div>
          <p>
            Do you want to buy {asset?.assetCode} Tokens with the following
            wallet?
          </p>
          <Dropdown
            label={activeWallet.alias ?? 'Choose wallet'}
            options={[
              ...appUser.userWallets.map((w, index) => ({
                text: w.alias,
                value: w,
                index,
              })),
            ]}
            onSelect={(selectedItem) => {
              console.log(selectedItem);
              setActiveWallet(selectedItem.value);
            }}
            key={`${activeWallet.publicKey}`}
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
            label={'CNGN'}
            defaultValue={{
              text: 'CNGN',
              value: 'CNGN',
            }}
            options={[
              {
                text: 'CNGN',
                value: 'CNGN',
              },
            ]}
            onSelect={(selectedItem) => {
              console.log(selectedItem);
            }}
            key={`${activeWallet.publicKey}`}
          />
          <div className="w-full space-y-2">
            <TextInput
              inputType="text"
              label="Amount"
              defaultValue={amount > 0 ? amount.toString() : ''}
              placeholder="5,000.0000"
              trailingIcon="/images/trovIcon.png"
              trailingText="CNGN"
              onInputChange={(value) => {
                var a = Number(value);
                if (a != null) {
                  setAmount(a);
                  var pricePerToken = asset?.pricePerToken!;
                  setQuantity(a / pricePerToken);
                } else {
                  setAmount(0);
                }
              }}
            />
            <p className="text-red-500">{error}</p>
          </div>
          <div className="space-y-3">
            <p>Quantity of {asset?.assetCode}</p>
            <div className="flex w-full py-5 px-4 xl:space-y-5 rounded-xl mb-10 items-center bg-primary-100">
              <p>{formatToDecimal(quantity)}</p>
            </div>
          </div>
          <div className="w-full self-center md:w-2/4 space-y-3">
            <Button
              label="Buy"
              additionalClasses="font-montserratSemiBold"
              onclick={() => {
                if (amount <= 0 || isNaN(Number(amount))) {
                  setError('Please enter a valid amount.');
                  return;
                }
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
              <p className="font-montserratSemiBold">
                {formatToDecimal(quantity)} Tokens
              </p>
              <p className="text-xs text-primary-800">
                of {asset?.assetName} Asset
              </p>
            </div>
          </div>
          <p className="text-primary-800 w-full mb-5 text-center">Amount</p>
          <div className="flex w-full py-3 px-2 xl:px-10 xl:space-y-5 rounded-xl justify-center items-center bg-primary-100">
            <div className="flex flex-col space-y-2 items-center w-full justify-between">
              <p className="font-montserratSemiBold">{amount} CNGN</p>
              {/* <p className="text-xs text-primary-800">$357.14</p> */}
            </div>
          </div>
          <p className="text-primary-800 w-full mb-5 text-center">Pay with</p>
          <div className="flex w-full py-5 px-2 xl:px-10 xl:space-y-5 rounded-xl mb-10 justify-center items-center bg-primary-100">
            <div className="flex space-x-2 items-center">
              <img src="/images/trovIcon.png" alt="" />
              <p className="text-xs text-primary-800">{activeWallet.alias}</p>
            </div>
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
                setPasswordErr('');

                if (!password) {
                  setPasswordErr('Please enter a password!');
                  return;
                }

                if (!(await isValidPassword())) {
                  setPasswordErr('Password is invalid!');
                  return;
                }

                const payload = {
                  signer: activeWallet.signer,
                  publicKey: activeWallet.publicKey,
                  secretKey: secretKey,
                  body: {
                    isSharedWallet: activeWallet.sharedAccessEnabled,
                    data: { amount: amount },
                    assetId: asset?.id,
                  },
                };

                toggleLoader();
                const res = await buyToken(payload);
                console.log('response', res);

                if ('data' in res) {
                  console.log('response', res);
                  const body = {
                    data: {
                      isSharedWallet: activeWallet.sharedAccessEnabled,
                      ...res.data,
                      commit: res.data.signatureRequired == 1 ? 0 : 1,
                      transactionSignature:
                        res.data.signatureRequired == 1
                          ? signBase64Txn(
                              secretKey,
                              res.data.transaction,
                              res.data.networkPassPhrase,
                            )
                          : '',
                    },
                    assetId: asset?.id,
                  };

                  const payload = {
                    signer: activeWallet.signer,
                    publicKey: activeWallet.publicKey,
                    secretKey: secretKey,
                    body,
                  };
                  console.log('body', payload);
                  const res2 = await buyToken(payload);
                  if ('data' in res2) {
                    console.log('response2', res2);
                    setBuyTokenStep(5);
                  } else if ('error' in res2) {
                    const errorResponse = res2.error as ErrorResponse;
                    showNotification(
                      'error',
                      errorResponse.data.message ??
                        'Something went wrong. Please try again.',
                    );
                  }
                } else if ('error' in res) {
                  const errorResponse = res.error as ErrorResponse;
                  showNotification(
                    'error',
                    errorResponse.data.message ??
                      'Something went wrong. Please try again.',
                  );
                }
                toggleLoader();
              }}
            />
          </div>
          {/* <div className="w-full flex space-x-3 pt-5">
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
            </div> */}
          <div />
        </div>
      );
    default:
      return (
        <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
          <img src="/images/success.png" alt="success" />
          <div className="flex flex-col text-center space-y-5 items-center w-2/3 md:px-10 justify-center">
            <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
              You have successfully purchased {amount} {asset.assetCode}!
            </p>
            <Button
              label="Done"
              additionalClasses="font-montserratSemiBold"
              onclick={async () => {
                onClose();
              }}
            />
          </div>
        </div>
      );
  }
}
