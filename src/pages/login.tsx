import React from 'react';
import Button from '../components/button';
import TextInput from '../components/textInput';
import ButtonSecondary from '../components/buttonSecondary';
import TrovoBrand from '../components/trovoBrand';

export default function Login() {
  return (
    <div className="flex h-full items-center justify-center ">
      <div className="hidden md:block w-3/5 h-full p-3">
        <div className="flex h-full space-y-3 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
          <TrovoBrand />
          <p className="text-primary-800 font-matahariExtended text-center text-3xl xl:text-4xl font-bold">
            Welcome back
          </p>
          <img
            className="w-100 h-100"
            src="/images/OTPsecurity.png"
            alt="Welcome 1"
          />
        </div>
      </div>
      <div className="w-full md:w-2/5 h-full overflow-y-scroll">
        <div className="flex flex-col  py-20 md:py-0 space-y-6 h-full items-center md:justify-center">
          <div className="flex flex-col px-10 w-full space-y-6 md:h-full items-center md:items-start justify-center">
            <p className="text-left w-full text-primary-800 text-3xl font-bold">
              Welcome back
              <span className="text-primary-700">&nbsp;Obi!</span>
            </p>
            <p className="text-left w-full text-primary-800 text-md">
              You have been missed
            </p>
            <div className="md:w-3/4">
              <TextInput
                label="Password"
                leadingIcon="/images/lock.png"
                inputType="password"
                placeholder="Enter answer"
                onInputChange={(newValue: string) => {
                  console.log('input has changed', newValue);
                }}
              />
            </div>
            <div className="flex justify-end w-full md:w-3/4 text-primary-800">
              <button type="button">Forgot Password?</button>
            </div>
            <div className="w-full md:w-3/4">
              <Button
                label="Sign In"
                onclick={() => {
                  // navigate('/register/backup');
                }}
              />
            </div>
            <div className="w-full md:w-3/4">
              <ButtonSecondary
                label="Create Account"
                onclick={() => {
                  // navigate('/register/backup');
                }}
              />
            </div>
            <div className="w-full md:w-3/4">
              <ButtonSecondary
                label="Recover Account"
                additionalClasses="bg-primary-600 text-white ring-primary-600"
                onclick={() => {
                  // navigate('/register/backup');
                }}
              />
            </div>
            <div className="flex justify-center w-full md:w-3/4 text-gray-500">
              <button type="button">Version 0.0.51</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
