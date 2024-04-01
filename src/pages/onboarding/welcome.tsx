import React from 'react';
import { useNavigate } from 'react-router-dom';
import Button from '../../components/button';
import ButtonSecondary from '../../components/buttonSecondary';
import TrovoBrand from '../../components/trovoBrand';

function Welcome() {
  const navigate = useNavigate();
  return (
    <div className="flex h-full items-center justify-center ">
      <div className="hidden md:block w-3/5 h-full p-3">
        <div className="flex h-full space-y-3 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
          <TrovoBrand />
          <p className="text-primary-800 font-matahariExtended text-center text-3xl xl:text-4xl font-bold">
            What would you
          </p>
          <p className="text-primary-700 font-matahariExtended text-center text-3xl xl:text-4xl font-bold">
            like to do?
          </p>
          <img
            className="w-100 h-100"
            src="/images/welcome.png"
            alt="Welcome 1"
          />
        </div>
      </div>
      <div className="w-full md:w-2/5 h-full">
        <div className="flex flex-col space-y-5 h-full items-center md:justify-center">
          <div className="md:hidden h-2/4 w-full p-3">
            <div className="flex h-full w-full space-y-2 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
              <TrovoBrand />
              <p className="text-primary-800 font-matahariExtended text-center text-2xl xl:text-4xl font-bold">
                What would you
              </p>
              <p className="text-primary-700 font-matahariExtended text-center text-2xl xl:text-4xl font-bold">
                like to do?
              </p>
              <img
                className="w-2/4"
                src="/images/welcome.png"
                alt="Welcome 1"
              />
            </div>
          </div>
          <div className="w-3/4">
            <Button
              label="Get Started"
              onclick={() => {
                navigate('/register');
              }}
            />
          </div>
          <div className="w-3/4">
            <ButtonSecondary
              label="Import Wallet"
              onclick={() => {
                navigate('/import');
              }}
            />
          </div>
          <div className="w-3/4">
            <ButtonSecondary
              label="Recover Account"
              additionalClasses="bg-primary-600 text-white ring-primary-600"
              onclick={() => {
                navigate('/recovery');
              }}
            />
          </div>
          <a className="underline text-primary-800" href="/dashboard/home">
            Go to Dashboard
          </a>
        </div>
      </div>
    </div>
  );
}

export default Welcome;
