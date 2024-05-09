import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Button from '../../components/button';
import TextInput from '../../components/textInput';
import OtpInput from '../../components/otpInput';
import Modal from '../../components/modal';
import TrovoBrand from '../../components/trovoBrand';
import { showNotification, toggleLoader } from '../../utils/showToaster';
import ButtonSecondary from '../../components/buttonSecondary';
import {
  useRecoveryOtpMutation,
  useVerifyEmailOtpMutation,
} from '../../store/api/authApi';
import { createAccount } from '../../utils/trovoSDK';
import { ErrorResponse } from '../../store/api/baseapi/axiosBaseQuery';
import { setFormState } from '../../store/authSlice';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../../store/reduxStore';

function RecoveryMain() {
  const [currentView, setCurrentView] = useState(0);
  const [username, setUsername] = useState('');
  const [tempAccount] = useState<Account>(createAccount());
  const [otp, setOtp] = useState('');
  const [showAlreadySentOtp, setShowAlreadySentOtp] = useState(false);
  const [
    showHaveYouSetupSecurityQuestion,
    setShowHaveYouSetupSecurityQuestion,
  ] = useState(false);
  const navigate = useNavigate();
  const [verifyEmailOtp] = useVerifyEmailOtpMutation();
  const [recoveryOtp] = useRecoveryOtpMutation();
  const tempData = useSelector((state: RootState) => state.auth.tempData);
  const dispatch = useDispatch();

  const verifyOtp = async (e: React.ChangeEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!otp) {
      showNotification('error', 'Please enter a valid 6 digits otp');
      return;
    }

    try {
      toggleLoader();

      const res = await verifyEmailOtp({
        signer: tempAccount.publicKey,
        publicKey: tempAccount.publicKey,
        secretKey: tempAccount.secretKey,
        body: { username, otp },
      });

      toggleLoader();
      console.log('res', res);
      setShowAlreadySentOtp(true);
      if ('data' in res) {
        setShowHaveYouSetupSecurityQuestion(true);
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

  const requestOtp = async (e: React.ChangeEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!username) {
      showNotification('error', 'Please enter your username to proceed');
      return;
    }

    try {
      toggleLoader();
      const res = await recoveryOtp({
        signer: tempAccount.publicKey,
        publicKey: tempAccount.publicKey,
        secretKey: tempAccount.secretKey,
        body: { username },
      });

      toggleLoader();
      console.log('res', res);
      setShowAlreadySentOtp(true);
      if ('data' in res) {
        setCurrentView(1);
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

  function renderSwitch(view: number) {
    switch (view) {
      case 1:
        return (
          <div className="flex flex-col space-y-6 md:h-full items-center justify-center">
            <p className="text-center w-full text-primary-800 text-2xl font-bold">
              Account Recovery
            </p>
            <div className="w-3/4 text-center px-5 md:px-10 py-5 bg-primary-100 rounded-xl">
              <p className="text-primary-800 text-md">
                An OTP has successfully been sent to your email.
                <span>
                  &nbsp;Please check your inbox (possibly the spam folder) to
                  get the OTP and then enter it into the text fields below.
                </span>
              </p>
            </div>
            <p className="text-primary-800 text-md xl:text-lg">
              Enter OTP for user:
              <span className="font-semibold">&nbsp;{username}</span>
            </p>
            <form
              id="security-answers"
              onSubmit={verifyOtp}
              className="space-y-6"
            >
              <div className="w-auto flex justify-center">
                <OtpInput
                  onInputChange={(newValue) => {
                    setOtp(newValue);
                  }}
                  numberOfDigits={6}
                />
              </div>
              <Button type="submit" label="Verify" onclick={() => {}} />
            </form>
            <div className="w-3/4"></div>
            <Modal
              showModal={showHaveYouSetupSecurityQuestion}
              onClose={() => {}}
            >
              <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
                <div className="flex flex-col text-center space-y-5 items-center w-2/3 mb-5 md:px-10 justify-center">
                  <p className="text-primary-800 text-md xl:text-lg font-bold">
                    Important!
                  </p>
                  <p className="text-primary-800 text-md xl:text-lg">
                    Have you setup security questions for this account?
                  </p>
                </div>
                <div className="w-3/4">
                  <Button
                    label="Yes, I have setup security questions"
                    onclick={() => {
                      dispatch(
                        setFormState({
                          ...tempData,
                          secretKey: tempAccount.secretKey,
                          publicKey: tempAccount.publicKey,
                          emailOtp: otp,
                          username: username,
                        }),
                      );
                      navigate('/answer-security-questions');
                    }}
                  />
                </div>
                <div className="w-full md:w-3/4">
                  <ButtonSecondary
                    label="No, I have not setup security questions"
                    onclick={async () => {
                      dispatch(
                        setFormState({
                          ...tempData,
                          secretKey: tempAccount.secretKey,
                          publicKey: tempAccount.publicKey,
                          username: username,
                        }),
                      );
                      navigate('/recovery');
                    }}
                  />
                </div>
              </div>
            </Modal>
          </div>
        );
      default:
        return (
          <>
            <p className="text-center w-full text-primary-800 text-2xl font-bold">
              Account Recovery
            </p>
            <div className="w-3/4 text-center px-5 md:px-10 py-5 bg-primary-100 rounded-xl">
              <p className="text-primary-800 text-md">
                Please provide accurate information in the fields below to
                ensure a successful account recovery process
              </p>
            </div>
            <form
              id="security-answers"
              onSubmit={requestOtp}
              className="w-3/4 space-y-6"
            >
              <TextInput
                label="Username"
                leadingIcon="/images/iconUser.png"
                inputType="text"
                onInputChange={(newValue) => {
                  setUsername(newValue);
                }}
              />
              <Button type="submit" label="Continue" onclick={() => {}} />
            </form>
            {showAlreadySentOtp && (
              <button
                onClick={() => {
                  setCurrentView(1);
                }}
              >
                Already recieved OTP?
              </button>
            )}
          </>
        );
    }
  }
  return (
    <div className="flex h-screen items-center justify-center ">
      <div className="hidden md:block w-3/5 h-full p-3">
        <div className="flex h-full space-y-3 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
          <TrovoBrand />
          <p className="text-primary-800 font-matahariExtended text-center text-3xl xl:text-4xl font-bold">
            Account
          </p>
          <p className="text-primary-700 font-matahariExtended text-center text-3xl xl:text-4xl font-bold">
            Recovery
          </p>
          <img
            className="w-100 h-100"
            src="/images/OTPsecurity.png"
            alt="Welcome 1"
          />
        </div>
      </div>
      <div className="w-full md:w-2/5 h-full overflow-y-scroll">
        <div className="flex flex-col md:px-20 py-20 md:py-0 space-y-6 h-full items-center md:justify-center">
          {renderSwitch(currentView)}
        </div>
      </div>
    </div>
  );
}

export default RecoveryMain;
