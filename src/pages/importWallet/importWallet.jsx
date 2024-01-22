import { React, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Button from '../../components/button';
import TextInput from '../../components/textInput';
import TrovoBrand from '../../components/trovoBrand';

export default function ImportWallet() {
  const [usePassphrase, setUsePassphrase] = useState(false);
  const navigate = useNavigate();
  return (
    <div className="flex md:h-screen items-center justify-center">
      <div className="hidden md:block w-3/5 h-full p-3">
        <div className="flex h-full space-y-3 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
          <TrovoBrand />
          <p className="text-primary-800 font-matahariExtended text-center text-3xl xl:text-4xl font-bold">
            Import
          </p>
          <p className="text-primary-700 font-matahariExtended text-center text-3xl xl:text-4xl font-bold">
            Wallet
          </p>
          <div />
          <div />
          <img
            className="w-100 h-100"
            src="/images/rafiki.png"
            alt="verification"
          />
        </div>
      </div>
      <div className="w-full md:w-2/5 pt-20">
        <div className="flex flex-col space-y-6 h-full overflow-y-scroll items-center md:justify-center">
          <div className="w-3/4">
            <TextInput
              label="Username or Email Address"
              leadingIcon="/images/email.png"
              inputType="text"
              onInputChange={(newValue) => {
                console.log('input has changed', newValue);
              }}
            />
          </div>
          <div className="w-3/4 flex space-x-3">
            <input
              type="checkbox"
              onChange={() => {
                setUsePassphrase(!usePassphrase);
              }}
              name="import"
            />
            <p className="text-gray-500">Enter pass phrase instead</p>
          </div>
          <div className="w-3/4">
            {usePassphrase ? (
              <textarea
                className="mt-2 ring-2 ring-gray-200 focus-within:ring-primary-600 rounded-md
              w-full h-12 py-1 px-2 focus-within:ring-2 flex items-center focus:outline-none"
                placeholder="Enter Pass phrase"
                rows="6"
                onChange={(newValue) => {
                  console.log('input has changed', newValue);
                }}
              />
            ) : (
              <TextInput
                label="Secret Key"
                leadingIcon="/images/lock.png"
                inputType="password"
                onInputChange={(newValue) => {
                  console.log('input has changed', newValue);
                }}
              />
            )}
          </div>
          <div className="w-3/4">
            <TextInput
              label="Password"
              leadingIcon="/images/lock.png"
              inputType="password"
              onInputChange={(newValue) => {
                console.log('input has changed', newValue);
              }}
            />
          </div>
          <div className="w-3/4">
            <TextInput
              label="Confirm Password"
              leadingIcon="/images/lock.png"
              inputType="password"
              onInputChange={(newValue) => {
                console.log('input has changed', newValue);
              }}
            />
          </div>
          <div className="w-3/4 flex space-x-3">
            <input type="checkbox" name="import" />
            <p className="text-gray-500">
              I agree to the Trovotech
              <a
                className="text-primary-800"
                href="https://trovotech.io/terms.html"
              >
                &nbsp;Terms of Service
              </a>
              &nbsp;and&nbsp;
              <a
                className="text-primary-800"
                href="https://trovotech.io/privacy-policy.html"
              >
                Privacy Policy
              </a>
            </p>
          </div>
          <div className="w-3/4">
            <Button
              label="Continue"
              onclick={() => {
                navigate('/register/verification');
              }}
            />
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
