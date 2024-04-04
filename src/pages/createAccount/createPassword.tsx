import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Button from '../../components/button';
import TextInput from '../../components/textInput';

export default function CreatePassword() {
  const navigate = useNavigate();
  const [userName, setUsername] = useState(false);

  return (
    <>
      <p className="text-primary-800 text-lg xl:text-xl">
        Create password to secure your account
      </p>
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
      <div className="w-3/4">
        <Button
          label="Continue"
          onclick={() => {
            navigate('/register/form');
          }}
        />
      </div>
    </>
  );
}
