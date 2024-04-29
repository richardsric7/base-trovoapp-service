import React, { useEffect } from 'react';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import OtpInput from '../../components/otpInput';
import Button from '../../components/button';
import Modal from '../../components/modal';
import ButtonSecondary from '../../components/buttonSecondary';
import TrovoBrand from '../../components/trovoBrand';
import { showNotification, toggleLoader } from '../../utils/showToaster';
import { useDispatch, useSelector } from 'react-redux';
import {
  useLazyGetUserQuery,
  useRegisterMutation,
} from '../../store/api/authApi';
import { RootState } from '../../store/reduxStore';
import { ErrorResponse } from '../../store/api/baseapi/axiosBaseQuery';
import { setFormState, setUser } from '../../store/authSlice';
import { Encryptor } from '../../types/encryptor';
import { User } from '../../types/user';

function AccountVerification() {
  const navigate = useNavigate();
  const [getUser] = useLazyGetUserQuery();
  const [verificationCode, setVerificationCode] = useState('');
  const [showModal, setShowModal] = useState(false);
  const dispatch = useDispatch();
  const [userRegister] = useRegisterMutation();
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const formInfo = useSelector((state: RootState) => state.auth.regFormInfo);

  useEffect(() => {
    if (!formInfo.password) {
      navigate('/register');
    }
  }, []);

  const handleSubmit = async () => {
    if (verificationCode.length !== 6 || isNaN(Number(verificationCode))) {
      showNotification('error', 'Please enter a valid 6 digits otp');
      return;
    }

    try {
      toggleLoader();

      const res = await userRegister({
        signer: appUser.publicKey,
        publicKey: appUser.publicKey,
        secretKey: formInfo.secretKey,
        body: { ...appUser, verificationCode },
      });

      toggleLoader();
      console.log('res', res);

      if ('data' in res) {
        const payload = {
          signer: appUser.publicKey,
          publicKey: appUser.publicKey,
          secretKey: formInfo.secretKey,
          body: { userId: appUser.username, import: 1 },
        };

        const { data, error } = await getUser(payload);

        if (data) {
          const response = data.userData as unknown as User;
          const encryptor = new Encryptor();
          const base64EncryptedData = await encryptor.encryptData(
            formInfo.secretKey,
            formInfo.password,
            appUser.publicKey,
          );
          dispatch(
            setUser({
              ...response,
              isLoggedIn: true,
              secretKeys: [base64EncryptedData],
            }),
          );

          dispatch(
            setFormState({
              ...formInfo,
              importExistingWallet: false,
              agreesToTerms: false,
              usePassphrase: false,
              secretKey: '',
              password: '',
              passphrase: '',
            }),
          );
          setShowModal(true);
        } else if (error) {
          const err = error as ErrorResponse;
          showNotification(
            'error',
            err.data.message ??
              'Sorry we could not complete the request. Please try again.',
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
    } catch (error: any) {
      console.log(error);
    }
  };

  return (
    <div className="flex h-full items-center justify-center ">
      <div className="hidden md:block w-3/5 h-full p-3">
        <div className="flex h-full space-y-3 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
          <TrovoBrand />
          <p className="text-primary-800 font-matahariExtended text-center text-3xl xl:text-4xl font-bold">
            Account
          </p>
          <p className="text-primary-700 font-matahariExtended text-center text-3xl xl:text-4xl font-bold">
            Verification
          </p>
          <img
            className="w-100 h-100"
            src="/images/OTPsecurity.png"
            alt="verification"
          />
        </div>
      </div>
      <div className="w-full md:w-2/5 h-3/5">
        <div className="flex flex-col space-y-6 md:h-full items-center justify-center">
          <div className="w-3/4 text-center px-5 md:px-10 py-5 bg-primary-100 rounded-xl">
            <p className="text-primary-800 text-md">
              Enter 6-digit code we just sent to your email address:
              <span className="font-semibold"> {appUser.email}</span>
            </p>
          </div>
          <p className="text-primary-800 text-md xl:text-lg font-semibold">
            Enter OTP
          </p>
          <div className="w-auto flex justify-center">
            <OtpInput
              numberOfDigits={6}
              onInputChange={(newValue) => {
                console.log('input has changed', newValue);
                setVerificationCode(newValue);
              }}
            />
          </div>
          <Modal
            showModal={showModal}
            onClose={() => {
              setShowModal(false);
            }}
          >
            <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
              <img src="/images/launch.png" alt="success" />
              <div className="flex flex-col text-center space-y-5 items-center w-2/3 md:px-10 justify-center">
                <p className="text-primary-800 text-md xl:text-lg font-semibold">
                  Congratulations!
                </p>
                <p className="text-primary-800 text-md xl:text-lg">
                  Your wallet has been successfully created
                </p>
                <p className="text-primary-800 text-md xl:text-lg">
                  We strongly recommend that you backup your wallet before
                  proceeding
                </p>
                <p className="text-primary-800 text-md xl:text-lg">
                  Backing up your wallet is a way to restore your wallet if you
                  loose your device
                </p>
              </div>
              <div className="w-3/4">
                <Button
                  label="Backup"
                  onclick={() => {
                    navigate('/backup');
                  }}
                />
              </div>
              <div className="w-3/4">
                <ButtonSecondary
                  label="Skip"
                  onclick={() => {
                    setShowModal(false);
                    navigate('/dashboard');
                  }}
                />
              </div>
            </div>
          </Modal>
          <div className="w-3/4">
            <Button
              label="Verify"
              onclick={() => {
                handleSubmit();
              }}
            />
          </div>
        </div>
      </div>
    </div>
  );
}

export default AccountVerification;
