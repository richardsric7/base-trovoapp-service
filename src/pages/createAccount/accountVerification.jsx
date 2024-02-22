import { React, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import OtpInput from '../../components/otpInput';
import Button from '../../components/button';
import Modal from '../../components/modal';
import ButtonSecondary from '../../components/buttonSecondary';
import TrovoBrand from '../../components/trovoBrand';

function AccountVerification() {
  const navigate = useNavigate();
  const [showModal, setShowModal] = useState(false);
  return (
    <div className="flex h-screen items-center justify-center ">
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
              <span className="font-semibold"> nancy@gmail.com</span>
            </p>
          </div>
          <p className="text-primary-800 text-md xl:text-lg font-semibold">
            Enter OTP
          </p>
          <div className="w-auto flex justify-center">
            <OtpInput numberOfDigits={6} />
          </div>
          <Modal showModal={showModal}>
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
                    navigate('/register/backup');
                    // setShowModal(true);
                  }}
                />
              </div>
              <div className="w-3/4">
                <ButtonSecondary
                  label="Skip"
                  onclick={() => {
                    // navigate('/register/backup');
                    setShowModal(false);
                  }}
                />
              </div>
            </div>
          </Modal>
          <div className="w-3/4">
            <Button
              label="Verify"
              onclick={() => {
                // navigate('/register/backup');
                setShowModal(true);
              }}
            />
          </div>
        </div>
      </div>
    </div>
  );
}

export default AccountVerification;
