import { React } from 'react';
import TrovoBrand from '../../components/trovoBrand';

export default function Wallet() {
  return (
    <div className="flex md:h-screen items-center justify-center">
      <div className="hidden md:block w-3/5 h-full p-3">
        <div className="flex h-full space-y-3 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
          <TrovoBrand />
          <p className="text-primary-800 font-matahariExtended text-center text-3xl xl:text-4xl font-bold">
            Wallet
          </p>
          <img
            className="w-100 h-100"
            src="/images/OTPsecurity.png"
            alt="verification"
          />
        </div>
      </div>
      <div className="w-full md:w-2/5 pt-20 lg:pt-0">
        <div className="flex flex-col px-5 md:px-20 space-y-6 md:h-full items-center justify-center">
          <p className="text-left w-full text-primary-800 text-2xl font-bold">
            Wallet Wallet
          </p>
        </div>
      </div>
    </div>
  );
}
