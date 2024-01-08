import React from 'react';
import { Outlet, useLocation } from 'react-router-dom';
// import CreatePassword from './createPassword';
// import RegistrationForm from './registrationForm';

function AuthMain() {
  return (
    <div className="flex h-screen items-center justify-center ">
      <div className="hidden md:block w-3/5 h-full p-3">
        <div className="flex h-full space-y-3 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
          <div className="flex items-center space-x-3 mt-5 mb-10 px-3 w-full">
            <img
              className="h-10 w-10"
              src="/images/trovoLogo.png"
              alt="trovo logo"
            />
            <span className="text-primary-800 font-montserratMedium text-xl xl:text-2xl">
              Trovo App
            </span>
          </div>
          <p className="text-primary-800 font-matahariExtended text-center text-3xl xl:text-4xl font-bold">
            Let&apos;s get you
          </p>
          <p className="text-primary-700 font-matahariExtended text-center text-3xl xl:text-4xl font-bold">
            started!
          </p>
          <img
            className="w-100 h-100"
            src="/images/welcome.png"
            alt="Welcome 1"
          />
        </div>
      </div>
      <div className="w-full md:w-2/5 h-full overflow-y-scroll">
        <div className="flex flex-col space-y-6 h-full overflow-y-scroll items-center md:justify-center">
          <ShowImageOnMobile />
          <Outlet />
        </div>
      </div>
    </div>
  );
}

function ShowImageOnMobile() {
  const location = useLocation();
  const isFormView = location.pathname.includes('/register/form');
  const className = 'h-2/4 w-full p-3';

  return (
    <div className={`${className} ${isFormView ? 'hidden' : 'md:hidden'}`}>
      <div className="flex h-full w-full space-y-2 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
        <div className="flex items-center self-start space-x-3 mt-5 mb-10 px-3 w-full">
          <img
            className="h-10 w-10"
            src="/images/trovoLogo.png"
            alt="trovo logo"
          />
          <span className="text-primary-800 font-montserratMedium text-xl xl:text-2xl">
            Trovo App
          </span>
        </div>
        <p className="text-primary-800 font-matahariExtended text-center text-2xl xl:text-4xl font-bold">
          Let&apos;s get you
        </p>
        <p className="text-primary-700 font-matahariExtended text-center text-2xl xl:text-4xl font-bold">
          started!
        </p>
        <img className="w-2/4" src="/images/welcome.png" alt="Welcome 1" />
      </div>
    </div>
  );
}

export default AuthMain;
