import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Button from '../../components/button';
import TextInput from '../../components/textInput';
import { User } from '../../types/user';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../../store/reduxStore';
import { setFormState, setTempUser } from '../../store/authSlice';
import 'react-phone-number-input/style.css';
import PhoneInput, {
  getCountries,
  isValidPhoneNumber,
} from 'react-phone-number-input';
import { useRegisterMutation } from '../../store/api/authApi';
import { createAccount, signHTTP } from '../../utils/trovoSDK';
import { showNotification, toggleLoader } from '../../utils/showToaster';
import {
  ErrorResponse,
  SuccessResponse,
} from '../../store/api/baseapi/axiosBaseQuery';

export default function RegistrationForm() {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const [userRegister] = useRegisterMutation();
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const formInfo = useSelector((state: RootState) => state.auth.regFormInfo);
  const [user, setUser] = useState<User>(appUser);
  const [isCorporate, setIsCorporate] = useState(user.isCorporateUser);
  // const [importExistingWallet, setImportExistingWallet] = useState(
  //   formInfo.importExistingWallet,
  // );
  // const [usePassphrase, setUsePassphrase] = useState(formInfo.usePassphrase);
  const [secretKey, setSecretKey] = useState(formInfo.secretKey);
  const [agreesToTerms, setAgreesToTerms] = useState(formInfo.agreesToTerms);
  const [errorObj, setErrorObj] = useState({
    firstName: '',
    lastName: '',
    username: '',
    email: '',
    phoneNumber: '',
    referrer: '',
    secretKey: '',
    passphrase: '',
    termsOfUse: '',
    general: '',
  });

  useEffect(() => {
    if (!formInfo.password) {
      navigate('/register');
    }
  }, [user, secretKey, errorObj]);

  const validateForm = (): boolean => {
    const usernamePattern = /^(?!.*\.\.)(?!.*\.$)[^\W][\w]{2,16}$/;
    const emailPattern =
      /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    let isValid = true;
    let newObj = errorObj;

    if (!user.firstName) {
      // setErrorObj();
      newObj = { ...newObj, firstName: 'Please enter your first name' };
      isValid = false;
    } else if (user.firstName.trim().replaceAll(' ', '').length < 3) {
      newObj = { ...newObj, firstName: 'Name must be at least 3 characters' };
      isValid = false;
    } else {
      newObj = { ...newObj, firstName: '' };
    }

    if (!user.lastName) {
      newObj = { ...newObj, lastName: 'Please enter your last name' };
      isValid = false;
    } else if (user.lastName.trim().replaceAll(' ', '').length < 3) {
      newObj = {
        ...newObj,
        lastName: 'Name must be at least 3 characters',
      };
      isValid = false;
    } else {
      newObj = {
        ...newObj,
        lastName: '',
      };
    }

    if (!user.email) {
      newObj = { ...newObj, email: 'Please enter an email address' };
      isValid = false;
    } else if (!emailPattern.test(user.email)) {
      newObj = {
        ...newObj,
        email: 'Please enter a valid email address',
      };
      isValid = false;
    } else {
      newObj = {
        ...newObj,
        email: '',
      };
    }

    if (!user.username) {
      newObj = { ...newObj, username: 'Choose a username' };
      isValid = false;
    } else if (
      user.username.trim().replaceAll(' ', '').length < 3 ||
      user.username.trim().replaceAll(' ', '').length > 16
    ) {
      newObj = {
        ...newObj,
        username: 'Username must be between 3 and 16 characters long',
      };
      isValid = false;
    } else if (!usernamePattern.test(user.username)) {
      newObj = {
        ...newObj,
        username: 'Username cannot contain any special character',
      };
      isValid = false;
    } else {
      newObj = {
        ...newObj,
        username: '',
      };
    }

    if (!user.mobile) {
      newObj = { ...newObj, phoneNumber: 'Please enter phone number' };
      isValid = false;
    } else if (!isValidPhoneNumber(user.mobile)) {
      newObj = {
        ...newObj,
        phoneNumber: 'This phone number is invalid.',
      };
      isValid = false;
    } else {
      newObj = {
        ...newObj,
        phoneNumber: '',
      };
    }

    if (
      user.referrer &&
      (user.referrer.trim().replaceAll(' ', '').length < 3 ||
        user.referrer.trim().replaceAll(' ', '').length > 16)
    ) {
      newObj = {
        ...newObj,
        referrer: 'This is not a valid Trovo username',
      };
      isValid = false;
    } else if (user.referrer && !usernamePattern.test(user.referrer)) {
      newObj = {
        ...newObj,
        referrer: 'This is not a valid Trovo username',
      };
      isValid = false;
    } else {
      newObj = {
        ...newObj,
        referrer: '',
      };
    }

    // if (importExistingWallet && !secretKey) {
    //   newObj = {
    //     ...newObj,
    //     secretKey: usePassphrase
    //       ? 'Please enter passphrase'
    //       : 'Please enter secret key',
    //   };
    //   isValid = false;
    // } else if (
    //   importExistingWallet &&
    //   !usePassphrase &&
    //   secretKey.trim().replaceAll(' ', '').length < 56
    // ) {
    //   newObj = {
    //     ...newObj,
    //     secretKey: 'Secret key must be 56 characters long',
    //   };
    //   isValid = false;
    // } else {
    //   newObj = {
    //     ...newObj,
    //     secretKey: '',
    //   };
    // }

    if (!agreesToTerms) {
      newObj = {
        ...newObj,
        termsOfUse: 'You need to accept terms',
      };
      isValid = false;
    } else {
      newObj = {
        ...newObj,
        termsOfUse: '',
      };
    }

    setErrorObj({ ...newObj });
    return isValid;
  };

  const handleSubmit = async (e: React.ChangeEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (validateForm()) {
      try {
        const account = createAccount();
        setSecretKey(account.secretKey);
        dispatch(setTempUser({ ...user, publicKey: account.publicKey }));
        dispatch(
          setFormState({
            ...formInfo,
            // importExistingWallet,
            agreesToTerms,
            // usePassphrase,
            secretKey: account.secretKey,
          }),
        );

        toggleLoader();

        const res = await userRegister({
          signer: account.publicKey,
          publicKey: account.publicKey,
          secretKey: account.secretKey,
          body: user,
        });

        toggleLoader();
        console.log('res', res);

        if ('data' in res) {
          const successResponse = res as SuccessResponse;
          showNotification('success', successResponse.data.message);
          navigate('/register/verification');
        } else if ('error' in res) {
          const errorResponse = res as ErrorResponse;
          showNotification(
            'error',
            errorResponse.error.data.message ??
              'Sorry we could not complete the request. Please try again.',
          );
        }
      } catch (error: any) {
        showNotification(
          'error',
          'Sorry something went wrong. Please try again.',
        );
      }
    }
  };

  return (
    <form
      id="reg-form"
      className="overflow-y-scroll flex flex-col pt-10 md:pt-16 md:pb-28 items-center w-full space-y-8 md:space-y-10 "
      onSubmit={handleSubmit}
    >
      <p className="text-primary-800 text-lg xl:text-xl">
        Fill in the Information to setup your account
      </p>
      <div className="flex w-full px-10 md:px-24 justify-center">
        <button
          type="button"
          onClick={() => {
            setIsCorporate(false);
            const newUser = { ...user, isCorporateUser: false };
            setUser(newUser);
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
            const newUser = { ...user, isCorporateUser: true };
            setUser(newUser);
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
      <div className="w-3/4 space-y-1">
        {isCorporate ? (
          <TextInput
            label="Entity Name"
            leadingIcon="/images/iconUser.png"
            inputType="text"
            defaultValue={user.firstName}
            onInputChange={(newValue) => {
              const newUser = { ...user, firstName: newValue };
              setUser(newUser);
            }}
          />
        ) : (
          <TextInput
            label="First Name"
            leadingIcon="/images/iconUser.png"
            inputType="text"
            defaultValue={user.firstName}
            onInputChange={(newValue) => {
              const newUser = { ...user, firstName: newValue };
              setUser(newUser);
            }}
          />
        )}
        {errorObj.firstName && (
          <p className="text-red-500 text-sm">{errorObj.firstName}</p>
        )}
      </div>
      {!isCorporate && (
        <div className="w-3/4 space-y-1">
          <TextInput
            label="Last Name"
            leadingIcon="/images/iconUser.png"
            inputType="text"
            defaultValue={user.lastName}
            onInputChange={(newValue) => {
              const newUser = { ...user, lastName: newValue };
              setUser(newUser);
            }}
          />
          {errorObj.lastName && (
            <p className="text-red-500 text-sm">{errorObj.lastName}</p>
          )}
        </div>
      )}
      <div className="w-3/4 space-y-1">
        <TextInput
          label="Email Address"
          leadingIcon="/images/email.png"
          inputType="text"
          defaultValue={user.email}
          onInputChange={(newValue) => {
            const newUser = { ...user, email: newValue };
            setUser(newUser);
          }}
        />
        {errorObj.email && (
          <p className="text-red-500 text-sm">{errorObj.email}</p>
        )}
      </div>
      <div className="w-3/4 space-y-1">
        <TextInput
          label="Username"
          leadingIcon="/images/iconUser.png"
          inputType="text"
          defaultValue={user.username}
          onInputChange={(newValue) => {
            const newUser = { ...user, username: newValue };
            setUser(newUser);
          }}
        />
        {errorObj.username && (
          <p className="text-red-500 text-sm">{errorObj.username}</p>
        )}
      </div>
      <div className="w-3/4 space-y-1">
        <label className="text-primary-700" htmlFor="Phone Input">
          Phone Number
        </label>
        <PhoneInput
          placeholder="Enter phone number"
          value={user.mobile}
          international
          onCountryChange={(newValue) => {
            const newUser = {
              ...user,
              mobileCountryCode: newValue?.toString() ?? '',
            };
            setUser(newUser);
          }}
          onChange={(newValue) => {
            const newUser = {
              ...user,
              mobile: newValue?.toString() ?? '',
            };
            setUser(newUser);
          }}
          defaultCountry="NG"
          containerComponentProps={{
            className:
              'mt-2 ring-1 md:ring-2 ring-gray-200 focus-within:ring-primary-600 rounded-md w-full h-12 py-1 px-2 focus-within:ring-2 flex items-center',
          }}
          numberInputProps={{
            required: true,
            className: 'h-full px-2 focus:outline-none w-full',
          }}
        />
        {errorObj.phoneNumber && (
          <p className="text-red-500 text-sm">{errorObj.phoneNumber}</p>
        )}
      </div>
      <div className="w-3/4 space-y-1">
        <TextInput
          label="Referrer (optional)"
          leadingIcon="/images/link.png"
          inputType="text"
          defaultValue={user.referrer}
          onInputChange={(newValue) => {
            const newUser = { ...user, referrer: newValue };
            setUser(newUser);
          }}
        />
        {errorObj.referrer && (
          <p className="text-red-500 text-sm">{errorObj.referrer}</p>
        )}
      </div>
      {/* <div className="w-3/4 flex space-x-3">
        <input
          type="checkbox"
          name="import"
          defaultChecked={importExistingWallet}
          onChange={() => {
            setImportExistingWallet(!importExistingWallet);
          }}
        />
        <p className="text-gray-500">Import existing wallet</p>
      </div>
      {importExistingWallet && (
        <div className="w-3/4 space-y-5">
          {usePassphrase ? (
            <div className="space-y-1">
              <textarea
                className="mt-2 ring-2 ring-gray-200 focus-within:ring-primary-600 rounded-md
              w-full h-12 py-1 px-2 focus-within:ring-2 flex items-center focus:outline-none"
                placeholder="Enter Passphrase"
                rows={8}
                defaultValue={secretKey}
                onChange={(evt) => {
                  setSecretKey(evt.target.value);
                }}
              />
              {errorObj.secretKey && (
                <p className="text-red-500 text-sm">{errorObj.secretKey}</p>
              )}
            </div>
          ) : (
            <div className="space-y-1">
              <TextInput
                label="Secret Key"
                leadingIcon="/images/iconUser.png"
                inputType="text"
                defaultValue={secretKey}
                onInputChange={(newValue) => {
                  setSecretKey(newValue);
                }}
              />
              {errorObj.secretKey && (
                <p className="text-red-500 text-sm">{errorObj.secretKey}</p>
              )}
            </div>
          )}
          <div className="w-3/4 flex space-x-3">
            <input
              type="checkbox"
              name="import"
              defaultChecked={usePassphrase}
              onChange={() => {
                setUsePassphrase(!usePassphrase);
              }}
            />
            <p className="text-gray-500">Use passphrase instead </p>
          </div>
        </div>
      )} */}
      <div className="w-3/4 space-y-1">
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
        {errorObj.termsOfUse && (
          <p className="text-red-500 text-sm">{errorObj.termsOfUse}</p>
        )}
      </div>
      <div className="w-3/4">
        <Button type="submit" label="Create Account" onclick={() => {}} />
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
    </form>
  );
}
