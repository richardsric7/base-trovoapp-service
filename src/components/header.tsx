import React from 'react';
import { useDispatch } from 'react-redux';
import { toggle } from '../store/sidebarSlice';
import TextInput from './textInput';

type Props = {
  fullName: string;
  email: string;
  avatar: string;
  isHomeView?: boolean;
};

// eslint-disable-next-line
function Header({ fullName, email, avatar, isHomeView = false }: Props) {
  const dispatch = useDispatch();
  const classes = `flex w-full p-3 items-center ${
    isHomeView ? 'justify-between' : 'justify-between xl:justify-end'
  }`;
  return (
    <div className={classes}>
      {isHomeView ? (
        <>
          <div className="flex space-x-4">
            <button
              className="xl:hidden"
              type="button"
              onClick={() => dispatch(toggle())}
            >
              <img src="/images/hamburgerMenu.png" alt="copy" className="w-8" />
            </button>
            <p>
              Good day,
              <span className="font-semibold text-lg"> Osondu</span>
            </p>
          </div>
          <div className="hidden md:block w-1/4">
            <TextInput
              leadingIcon="/images/search.png"
              inputType="text"
              label=""
              placeholder="Search"
              onInputChange={(newValue: string) => {
                // console.log('input has changed', newValue);
              }}
            />
          </div>
        </>
      ) : (
        <button
          className="xl:hidden"
          type="button"
          onClick={() => dispatch(toggle())}
        >
          <img src="/images/hamburgerMenu.png" alt="copy" className="w-8" />
        </button>
      )}
      <div className="flex space-x-3 items-center">
        <img
          className="h-6"
          src="/images/notification.png"
          alt="notification bell"
        />
        <img className="h-10" src={avatar} alt="avatar" />
        <div className="flex flex-col hidden md:block space-y-2">
          <p className="font-semibold">{fullName}</p>
          <p>{email}</p>
        </div>
      </div>
    </div>
  );
}

export default Header;
