import React from 'react';

type Props = {
  children: React.ReactNode;
  showModal: boolean;
  onClose: () => void;
};

export default function Modal({ children, showModal, onClose }: Props) {
  return (
    <div>
      {showModal ? (
        <>
          <div className="justify-center items-center flex overflow-x-hidden overflow-y-auto fixed inset-0 z-50 outline-none focus:outline-none">
            <div className="relative w-full my-6 mx-5 md:mx-auto max-w-3xl">
              {/* content */}
              <div className="border-0 rounded-lg shadow-lg relative flex flex-col w-full bg-white outline-none focus:outline-none">
                <div className="w-full py-2 px-5">
                  <button
                    className="font-semibold text-lg"
                    type="button"
                    onClick={() => {
                      onClose();
                    }}
                  >
                    X
                  </button>
                </div>
                {/* body */}
                {children}
              </div>
            </div>
          </div>
          <div className="opacity-25 fixed inset-0 z-40 bg-black" />
        </>
      ) : null}
    </div>
  );
}
