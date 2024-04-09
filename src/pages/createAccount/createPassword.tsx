import React, { useEffect, useState, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import Button from '../../components/button';
import TextInput from '../../components/textInput';
import { useDispatch, useSelector } from 'react-redux';
import { setTempUser } from '../../store/authSlice';
import { RootState } from '../../store/reduxStore';

enum FieldState {
  error,
  pristine,
  ok,
}

interface FormFieldGuide {
  info: string;
  fieldState: FieldState;
}

export default function CreatePassword() {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const isFirst = useRef(true);
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
  const user = useSelector((state: RootState) => state.auth.user);
  const [password, setPassword] = useState(user.password);
  const [confirmPassword, setConfirmPassword] = useState(user.password);
  const [confirmPasswordErr, setConfirmPasswordErr] = useState('');
  const [passwordGuides, setPasswordGuide] =
    useState<FormFieldGuide[]>(initialGuidesState);

  const validatePassword = (): boolean => {
    let isValid = true;
    const newGuides = [...passwordGuides];

    const isWhitespace = /^(?=.*\s)/;
    if (isWhitespace.test(password)) {
      newGuides[0].fieldState = FieldState.error;
      isValid = false;
    } else {
      newGuides[0].fieldState = FieldState.ok;
    }

    const isContainsUppercase = /^(?=.*[A-Z])/;
    if (!isContainsUppercase.test(password)) {
      newGuides[1].fieldState = FieldState.error;
      isValid = false;
    } else {
      newGuides[1].fieldState = FieldState.ok;
    }

    const isContainsLowercase = /^(?=.*[a-z])/;
    if (!isContainsLowercase.test(password)) {
      newGuides[2].fieldState = FieldState.error;
      isValid = false;
    } else {
      newGuides[2].fieldState = FieldState.ok;
    }

    const isContainsNumber = /^(?=.*[0-9])/;
    if (!isContainsNumber.test(password)) {
      newGuides[3].fieldState = FieldState.error;
      isValid = false;
    } else {
      newGuides[3].fieldState = FieldState.ok;
    }

    // eslint-disable-next-line
    const isContainsSymbol = /^(?=.*[~`!@#$%^&*()--+={}\[\]|\\:;"'<>,.?/_₹])/;
    if (!isContainsSymbol.test(password)) {
      newGuides[4].fieldState = FieldState.error;
      isValid = false;
    } else {
      newGuides[4].fieldState = FieldState.ok;
    }

    // const isValidLength = /^.{6,16}$/;
    if (password.length < 6 || password.length > 16) {
      newGuides[5].fieldState = FieldState.error;
      isValid = false;
    } else {
      newGuides[5].fieldState = FieldState.ok;
    }

    setPasswordGuide(newGuides);
    return isValid;
  };

  const guides = passwordGuides.map((guide, index) => {
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
    console.log('lsdksds', isFirst.current);
    if (password) {
      validatePassword();
    } else {
      setPasswordGuide([...initialGuidesState]);
    }
  }, [password]);

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
          defaultValue={password}
          onInputChange={(newValue) => {
            setPassword(newValue);
          }}
        />
      </div>
      <div className="w-3/4">{guides}</div>
      <div className="w-3/4">
        <TextInput
          label="Confirm Password"
          leadingIcon="/images/lock.png"
          inputType="password"
          defaultValue={confirmPassword}
          onInputChange={(newValue) => {
            setConfirmPasswordErr('');
            setConfirmPassword(newValue);
          }}
        />
        <p className="text-red-500 text-sm">{confirmPasswordErr}</p>
      </div>
      <div className="w-3/4">
        <Button
          label="Continue"
          onclick={() => {
            let isValid = true;
            setConfirmPasswordErr('');

            if (!confirmPassword) {
              setConfirmPasswordErr('Re-enter your password');
              isValid = false;
            } else if (password != confirmPassword) {
              setConfirmPasswordErr('Both passwords didn’t match. Try again.');
              isValid = false;
            }

            if (validatePassword() && isValid) {
              const newUser = { ...user, password: password };
              dispatch(setTempUser(newUser));
              navigate('/register/form');
            }
          }}
        />
      </div>
    </>
  );
}
