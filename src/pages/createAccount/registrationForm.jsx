import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Button from '../../components/button';
import TextInput from '../../components/textInput';

export default function RegistrationForm() {
  const navigate = useNavigate();
  const [isCorporate, setIsCorporate] = useState(false);

  return (
    <div className="overflow-y-scroll flex flex-col pt-10 md:pt-16 md:pb-28 items-center w-full space-y-8 md:space-y-10 ">
      <p className="text-primary-800 text-lg xl:text-xl">
        Fill in the Information to setup your account
      </p>
      <div className="flex w-full px-10 md:px-24 justify-center">
        <button
          type="button"
          onClick={() => {
            setIsCorporate(false);
          }}
          className={`px-2 py-3 w-full rounded-l-md ${
            isCorporate
              ? 'bg-primary-100 text-primary-800'
              : 'bg-primary-800 text-white'
          } `}
        >
          Individual
        </button>
        <button
          type="button"
          onClick={() => {
            setIsCorporate(true);
          }}
          className={`px-2 py-3 w-full rounded-r-md ${
            isCorporate
              ? 'bg-primary-800 text-white'
              : 'bg-primary-100 text-primary-800'
          }`}
        >
          Corporate
        </button>
      </div>
      {isCorporate ? (
        <div className="w-3/4">
          <TextInput
            label="Entity Name"
            leadingIcon="/images/iconUser.png"
            inputType="text"
            onInputChange={(newValue) => {
              console.log('input has changed', newValue);
            }}
          />
        </div>
      ) : (
        <div className="w-3/4">
          <TextInput
            label="First Name"
            leadingIcon="/images/iconUser.png"
            inputType="text"
            onInputChange={(newValue) => {
              console.log('input has changed', newValue);
            }}
          />
        </div>
      )}
      {!isCorporate && (
        <div className="w-3/4">
          <TextInput
            label="Last Name"
            leadingIcon="/images/iconUser.png"
            inputType="text"
            onInputChange={(newValue) => {
              console.log('input has changed', newValue);
            }}
          />
        </div>
      )}
      <div className="w-3/4">
        <TextInput
          label="Email Address"
          leadingIcon="/images/email.png"
          inputType="text"
          onInputChange={(newValue) => {
            console.log('input has changed', newValue);
          }}
        />
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
        <TextInput
          label="Phone Number"
          leadingIcon="/images/lock.png"
          inputType="text"
          onInputChange={(newValue) => {
            console.log('input has changed', newValue);
          }}
        />
      </div>
      <div className="w-3/4">
        <TextInput
          label="Referrer (optional)"
          leadingIcon="/images/link.png"
          inputType="text"
          onInputChange={(newValue) => {
            console.log('input has changed', newValue);
          }}
        />
      </div>
      <div className="w-3/4 flex space-x-3">
        <input type="checkbox" name="import" />
        <p className="text-gray-500">Import existing wallet</p>
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
          label="Create Account"
          onclick={() => {
            navigate('/register/verification');
          }}
        />
      </div>
      <p className="text-gray-500">
        Already have an account?
        <a className="text-primary-800" href="https://trovotech.io/terms.html">
          &nbsp;Sign In
        </a>
      </p>
      <div />
      <div />
      <div />
      <div />
    </div>
  );
}
