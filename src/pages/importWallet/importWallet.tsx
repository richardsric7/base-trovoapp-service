import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Button from '../../components/button';
import TextInput from '../../components/textInput';
import TrovoBrand from '../../components/trovoBrand';
import ButtonSecondary from '../../components/buttonSecondary';
import { FieldState, FormFieldGuide } from '../../types/user';
import { getCredsFromPassPhrase, parseSecretKey } from '../../utils/trovoSDK';
import { showNotification, toggleLoader } from '../../utils/showToaster';
import { useLazyGetUserQuery } from '../../store/api/authApi';
import { useDispatch } from 'react-redux';
import { Encryptor } from '../../utils/encryptor';
import { setUser } from '../../store/authSlice';
import { ErrorResponse } from '../../store/api/baseapi/axiosBaseQuery';
import { deserializeUserData } from '../../utils/deserializeAndStoreUserData';

export default function ImportWallet() {
  const [usePassphrase, setUsePassphrase] = useState(false);
  const [agreesToTerms, setAgreesToTerms] = useState(true);

  const [tempData, setTempData] = useState({
    userId: '',
    secret: '',
    password: '',
    confirmPassword: '',
  });

  const [getUser, {}] = useLazyGetUserQuery();

  const dispatch = useDispatch();
  const [errorObj, setErrorObj] = useState({
    userId: '',
    secret: '',
    password: '',
    termsOfUse: '',
    confirmPassword: '',
  });
  const navigate = useNavigate();

  const initialGuidesState = [
    {
      info: 'Password must not contain whitespaces',
      fieldState: FieldState.pristine,
    },
    {
      info: 'Password must have at least one uppercase character',
      fieldState: FieldState.pristine,
    },
    {
      info: 'Password must have at least one lowercase character',
      fieldState: FieldState.pristine,
    },
    {
      info: 'Password must contain at least one digit',
      fieldState: FieldState.pristine,
    },
    {
      info: 'Password must contain at least one special symbol',
      fieldState: FieldState.pristine,
    },
    {
      info: 'Password length be minimum of 6 and maximimum of 16',
      fieldState: FieldState.pristine,
    },
  ];
  const [passwordGuides, setPasswordGuide] =
    useState<FormFieldGuide[]>(initialGuidesState);

  const validatePassword = (): boolean => {
    let isValid = true;
    let newGuides = [...passwordGuides];

    const isWhitespace = /^(?=.*\s)/;
    if (isWhitespace.test(tempData.password)) {
      newGuides[0].fieldState = FieldState.error;
      isValid = false;
    } else {
      newGuides[0].fieldState = FieldState.ok;
    }

    const isContainsUppercase = /^(?=.*[A-Z])/;
    if (!isContainsUppercase.test(tempData.password)) {
      newGuides[1].fieldState = FieldState.error;
      isValid = false;
    } else {
      newGuides[1].fieldState = FieldState.ok;
    }

    const isContainsLowercase = /^(?=.*[a-z])/;
    if (!isContainsLowercase.test(tempData.password)) {
      newGuides[2].fieldState = FieldState.error;
      isValid = false;
    } else {
      newGuides[2].fieldState = FieldState.ok;
    }

    const isContainsNumber = /^(?=.*[0-9])/;
    if (!isContainsNumber.test(tempData.password)) {
      newGuides[3].fieldState = FieldState.error;
      isValid = false;
    } else {
      newGuides[3].fieldState = FieldState.ok;
    }

    // eslint-disable-next-line
    const isContainsSymbol = /^(?=.*[~`!@#$%^&*()--+={}\[\]|\\:;"'<>,.?/_₹])/;
    if (!isContainsSymbol.test(tempData.password)) {
      newGuides[4].fieldState = FieldState.error;
      isValid = false;
    } else {
      newGuides[4].fieldState = FieldState.ok;
    }

    // const isValidLength = /^.{6,16}$/;
    if (tempData.password.length < 6 || tempData.password.length > 16) {
      newGuides[5].fieldState = FieldState.error;
      isValid = false;
    } else {
      newGuides[5].fieldState = FieldState.ok;
    }

    setPasswordGuide(newGuides);
    return isValid;
  };

  const guides = passwordGuides.map((guide) => {
    const additionalClasses =
      guide.fieldState == FieldState.error
        ? 'text-red-500'
        : guide.fieldState == FieldState.ok
        ? 'text-green-500'
        : '';
    return (
      <p
        key={`${new Date().getTime()}${guide.info.replace(' ', '')}`}
        className={`text-gray-400 text-left mt-1 ${additionalClasses}`}
      >
        - {guide.info}
      </p>
    );
  });

  useEffect(() => {
    if (tempData.password) {
      validatePassword();
    } else {
      setPasswordGuide([...initialGuidesState]);
    }
  }, [tempData, errorObj]);

  const validateForm = (): boolean => {
    let isValid = true;
    let newObj = errorObj;

    if (!tempData.userId) {
      newObj = {
        ...newObj,
        userId: 'Please enter your username or email address.',
      };
      isValid = false;
    } else if (
      tempData.userId.trim().replaceAll(' ', '').length < 3 ||
      tempData.userId.trim().replaceAll(' ', '').length > 16
    ) {
      newObj = {
        ...newObj,
        userId: 'Username/Email must be between 3 and 16 characters long',
      };
      isValid = false;
    } else {
      newObj = {
        ...newObj,
        userId: '',
      };
    }
    if (!tempData.secret) {
      newObj = {
        ...newObj,
        secret: usePassphrase
          ? 'Please enter passphrase'
          : 'Please enter secret key',
      };
      isValid = false;
    } else {
      newObj = {
        ...newObj,
        secret: '',
      };
    }

    if (!validatePassword()) {
      isValid = false;
    }

    if (!tempData.password) {
      newObj = {
        ...newObj,
        password: 'Please enter your password!',
      };
      isValid = false;
    }

    if (!tempData.confirmPassword) {
      newObj = {
        ...newObj,
        confirmPassword: 'Please confirm your password!',
      };
      isValid = false;
    }

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

  const getAccountFromExistingInfo = () =>
    usePassphrase
      ? getCredsFromPassPhrase(tempData.secret)!
      : parseSecretKey(tempData.secret);

  const handleSubmit = async (e: React.ChangeEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (validateForm()) {
      try {
        const account = getAccountFromExistingInfo();

        const payload = {
          signer: account.publicKey,
          publicKey: account.publicKey,
          secretKey: account.secretKey,
          body: { userId: tempData.userId, import: 1 },
        };

        toggleLoader();
        const { data, error } = await getUser(payload);
        toggleLoader();

        if (data) {
          try {
            showNotification('success', 'Wallet successfully imported!');
            const userData = deserializeUserData(data);

            const encryptor = new Encryptor();
            const base64EncryptedSecretKey = await encryptor.encryptData(
              account.secretKey,
              tempData.password,
              account.publicKey,
            );

            const user = {
              ...userData,
              isLoggedIn: true,
              secretKeys: [base64EncryptedSecretKey],
            };

            const hash = await encryptor.createHash(user.username);
            const base64EncryptedUserData = await encryptor.encryptData(
              JSON.stringify(user),
              hash,
              user.publicKey,
            );
            dispatch(
              setUser({
                key: hash,
                user,
                encryptedUser: base64EncryptedUserData,
              }),
            );

            navigate('/dashboard');
          } catch (error: any) {
            console.log('err', error);
          }
        } else if (error) {
          const err = error as ErrorResponse;
          showNotification(
            'error',
            err.data.message ??
              'Sorry we could not complete the request. Please try again.',
          );
        }
      } catch (error: any) {
        showNotification(
          'error',
          'Sorry something went wrong. Please check your inputs try again.',
        );
      }
    }
  };

  return (
    <div className="flex h-full items-center justify-center ">
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
      <div className="w-full md:w-2/5 h-full overflow-y-scroll">
        <form
          id="import-form"
          onSubmit={handleSubmit}
          className="flex flex-col space-y-6 h-full overflow-y-scroll items-center md:justify-center"
        >
          <div className="w-3/4 space-y-1">
            <TextInput
              label="Username or Email Address"
              leadingIcon="/images/email.png"
              inputType="text"
              onInputChange={(newValue) => {
                setTempData({ ...tempData, userId: newValue });
              }}
            />
            {errorObj.userId && (
              <p className="text-red-500 text-sm">{errorObj.userId}</p>
            )}
          </div>
          <div className="w-3/4 space-y-3">
            {usePassphrase ? (
              <div className="space-y-2">
                <label className="text-primary-700" htmlFor="Phone Input">
                  Enter Passphrase
                </label>
                <textarea
                  className="mt-2 ring-2 ring-gray-200 focus-within:ring-primary-600 rounded-md
              w-full h-12 py-1 px-2 focus-within:ring-2 flex items-center focus:outline-none"
                  placeholder="Enter Pass phrase"
                  rows={6}
                  onChange={(evt) => {
                    setTempData({ ...tempData, secret: evt.target.value });
                  }}
                />
              </div>
            ) : (
              <TextInput
                label="Secret Key"
                leadingIcon="/images/lock.png"
                inputType="password"
                onInputChange={(newValue) => {
                  setTempData({ ...tempData, secret: newValue });
                }}
              />
            )}
            {errorObj.secret && (
              <p className="text-red-500 text-sm">{errorObj.secret}</p>
            )}
            <div className="w-3/4 flex space-x-3">
              <input
                type="checkbox"
                onChange={() => {
                  setUsePassphrase(!usePassphrase);
                  setTempData({ ...tempData, secret: '' });
                }}
                name="import"
              />
              <p className="text-gray-500">Enter pass phrase instead</p>
            </div>
          </div>
          <div className="w-3/4 space-y-1">
            <TextInput
              label="Password"
              leadingIcon="/images/lock.png"
              inputType="password"
              onInputChange={(newValue) => {
                setTempData({ ...tempData, password: newValue });
              }}
            />
            {errorObj.password && (
              <p className="text-red-500 text-sm">{errorObj.password}</p>
            )}
          </div>
          <div className="w-3/4">{guides}</div>
          <div className="w-3/4 space-y-1">
            <TextInput
              label="Confirm Password"
              leadingIcon="/images/lock.png"
              inputType="password"
              onInputChange={(newValue) => {
                setTempData({ ...tempData, confirmPassword: newValue });
              }}
            />
            {errorObj.confirmPassword && (
              <p className="text-red-500 text-sm">{errorObj.confirmPassword}</p>
            )}
          </div>
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
              <p className="text-red-500 w-3/4 text-sm">
                {errorObj.termsOfUse}
              </p>
            )}
          </div>
          <div className="w-3/4 space-y-3">
            <Button
              type="submit"
              label="Continue"
              onclick={() => {
                // validateForm();
              }}
            />
            <ButtonSecondary
              label="Back"
              onclick={() => {
                navigate(-1);
              }}
            />
          </div>
        </form>
      </div>
    </div>
  );
}
