import React from 'react';
import { TEToast } from 'tw-elements-react';

type Props = {
  type: 'success' | 'info' | 'error' | undefined;
  message: string;
  onClose: () => void;
};

export default function Toaster({
  type,
  onClose,
  message = '',
}: Props): JSX.Element {
  const getTypeHeader = (): string => {
    switch (type) {
      case 'success':
        return 'Success';
      case 'info':
        return 'Info';
      case 'error':
        return 'Error';
      default:
        return '';
    }
  };

  return (
    <div>
      <TEToast
        open={type !== undefined}
        delay={3000}
        autohide={true}
        onClose={onClose}
        color={`bg-${type === 'info' ? 'primary' : type}-100 text-${
          type === 'info' ? 'primary' : type
        }-700`}
      >
        <div
          className={`flex items-center justify-between rounded-t-lg border-b-2 border-${
            type === 'info' ? 'primary' : type
          }/20 border-opacity-100 bg-clip-padding px-4 pb-2 pt-2.5`}
        >
          <p className="font-bold">{getTypeHeader()}</p>
          <div className="flex items-center">
            <button
              type="button"
              className="ml-2 box-content rounded-none border-none opacity-80 hover:no-underline hover:opacity-75 focus:opacity-100 focus:shadow-none focus:outline-none"
              onClick={onClose}
              aria-label="Close"
            >
              <span className="w-[1em] focus:opacity-100 disabled:pointer-events-none disabled:select-none disabled:opacity-25 [&.disabled]:pointer-events-none [&.disabled]:select-none [&.disabled]:opacity-25">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth="1.5"
                  stroke="currentColor"
                  className="h-6 w-6"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </span>
            </button>
          </div>
        </div>
        <div className="break-words rounded-b-lg px-4 py-4">{message}</div>
      </TEToast>
    </div>
  );
}
