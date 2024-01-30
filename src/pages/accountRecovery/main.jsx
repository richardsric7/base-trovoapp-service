import { React, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import PropTypes from 'prop-types';
import Button from '../../components/button';
import TextInput from '../../components/textInput';
import OtpInput from '../../components/otpInput';
import Modal from '../../components/modal';
import TrovoBrand from '../../components/trovoBrand';

function RecoveryMain() {
  const [currentView, setCurrentView] = useState(0);
  function renderSwitch(view) {
    switch (view) {
      case 1:
        return (
          <ShowEnterUsername
            onDone={() => {
              setCurrentView(2);
            }}
          />
        );
      case 2:
        return (
          <ShowEnterOTP
            onDone={() => {
              setCurrentView(3);
            }}
          />
        );
      case 3:
        return <ShowAnswerSecurityQuestions />;
      default:
        return (
          <ShowPreliminary
            onDone={() => {
              setCurrentView(1);
            }}
          />
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

function ShowPreliminary({ onDone }) {
  return (
    <>
      <div className="px-10 text-primary-800 space-y-5 text-md text-justify">
        <p className="text-center w-full text-primary-800 text-2xl font-bold">
          Account Recovery
        </p>
        <p>
          You are about to opt in for and enable account recovery service, once
          account recovery service is enabled on your account, Trovotech will be
          able to recover your account should you loose your secret key within
          the time the recovery service is active on your account.
        </p>
        <p>
          The account recovery is a paid service and is only limited to recovery
          of your account and the wallet(s) on your account, the payment for
          this service is renewable on an annual basis.
        </p>
        <p>
          Trovotech des not have access to your secret key, therefore will not
          be liable for any missing asset in your wallet(s).
        </p>
        <p>
          As long as your secret key remains safe on your side, your assets are
          safe. Therefore, ensure that you keep your secret key safe always.
        </p>
      </div>
      <div className="w-full px-10 pt-10">
        <Button
          label="Continue"
          onclick={() => {
            onDone();
          }}
        />
      </div>
    </>
  );
}

function ShowEnterUsername({ onDone }) {
  return (
    <>
      <p className="text-center w-full text-primary-800 text-2xl font-bold">
        Account Recovery
      </p>
      <div className="w-3/4 text-center px-5 md:px-10 py-5 bg-primary-100 rounded-xl">
        <p className="text-primary-800 text-md">
          Please provide accurate information in the fields below to ensure a
          successful account recovery process
        </p>
      </div>
      <div className="w-3/4">
        <TextInput
          label="Username"
          leadingIcon="/images/iconUser.png"
          inputType="text"
          onInputChange={(newValue) => {
            console.log('input has changed', newValue);
          }}
        />
      </div>
      <div className="w-3/4">
        <Button
          label="Continue"
          onclick={() => {
            onDone();
          }}
        />
      </div>
    </>
  );
}

function ShowEnterOTP({ onDone }) {
  return (
    <div className="flex flex-col space-y-6 md:h-full items-center justify-center">
      <p className="text-center w-full text-primary-800 text-2xl font-bold">
        Account Recovery
      </p>
      <div className="w-3/4 text-center px-5 md:px-10 py-5 bg-primary-100 rounded-xl">
        <p className="text-primary-800 text-md">
          An OTP has successfully been sent to your email&nbsp;
          <span className="font-semibold"> obixxx@gmail.com.</span>
          <span>
            &nbsp;Please check your inbox (possibly the spam folder) to get the
            OTP and then enter it into the text fields below.
          </span>
        </p>
      </div>
      <p className="text-primary-800 text-md xl:text-lg">
        Enter OTP for user:
        <span className="font-semibold">&nbsp;nancy</span>
      </p>
      <div className="w-auto flex justify-center">
        <OtpInput numberOfDigits={6} />
      </div>
      <div className="flex justify-center text-primary-800 w-full">
        <button type="button">Resend OTP?</button>
      </div>
      <div className="w-3/4">
        <Button
          label="Verify"
          onclick={() => {
            onDone();
          }}
        />
      </div>
    </div>
  );
}

function ShowAnswerSecurityQuestions() {
  const [showModal, setShowModal] = useState(false);
  const navigate = useNavigate();

  return (
    <div className="flex flex-col space-y-6 md:h-full items-center justify-center">
      <p className="text-center w-full text-primary-800 text-2xl font-bold">
        Answer Security Questions
      </p>
      <div className="w-3/4 text-center px-5 md:px-10 py-5 bg-primary-100 rounded-xl">
        <p className="text-primary-800 text-md">
          Please answer the following security questions to proceed to the next
          step.
        </p>
      </div>
      <div className="w-3/4">
        <TextInput
          label="What is your mother’s maiden name?"
          leadingIcon="/images/question.png"
          inputType="text"
          placeholder="Enter answer"
          onInputChange={(newValue) => {
            console.log('input has changed', newValue);
          }}
        />
      </div>
      <div className="w-3/4">
        <TextInput
          label="What is the name of the first street you lived at?"
          leadingIcon="/images/question.png"
          inputType="text"
          placeholder="Enter answer"
          onInputChange={(newValue) => {
            console.log('input has changed', newValue);
          }}
        />
      </div>
      <div className="w-3/4">
        <TextInput
          label="What is your mother’s maiden name?"
          leadingIcon="/images/question.png"
          inputType="text"
          placeholder="Enter answer"
          onInputChange={(newValue) => {
            console.log('input has changed', newValue);
          }}
        />
      </div>
      <Modal showModal={showModal}>
        <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
          <img src="/images/launch.png" alt="success" />
          <div className="flex flex-col text-center space-y-5 items-center w-2/3 mb-5 md:px-10 justify-center">
            <p className="text-primary-800 text-md xl:text-lg font-semibold">
              Congratulations!
            </p>
            <p className="text-primary-800 text-md xl:text-lg">
              You have successfully recovered your account. Please copy your
              secret key below to import your wallet afresh from your device.
            </p>
          </div>
          <div className="w-3/4 text-justify space-y-5 px-5 md:px-10 mb-5 py-5 bg-primary-100 rounded-xl">
            <p className="text-primary-800 text-md xl:text-lg font-semibold">
              Secret key for Ogbonge Wallet
            </p>
            <p className="text-primary-800 text-md break-all">
              {/* eslint-disable-next-line */}
              ASKONRWINOT-949I0IWRINEKLNFSKNFSDPG0J3-9JGRNKNGKN0J34INORGSDFW4W4WWEKLNDKLSFWKLNREKONELN4T4U48T53UONGKNGKLNK34T34WRKLGNKLGNREGKLNRKENGKLRNGEL
            </p>
            <div className="flex justify-between w-full">
              <button type="button">
                <img src="/images/copy.png" alt="copy" />
              </button>
              <button type="button">
                <img src="/images/eyeShow.png" alt="show/hide" />
              </button>
            </div>
          </div>
          <div className="w-3/4">
            <Button
              label="Go to Dashboard"
              onclick={() => {
                navigate('/home');
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
  );
}

ShowPreliminary.propTypes = {
  onDone: PropTypes.func.isRequired,
};

ShowEnterUsername.propTypes = {
  onDone: PropTypes.func.isRequired,
};

ShowEnterOTP.propTypes = {
  onDone: PropTypes.func.isRequired,
};
export default RecoveryMain;
