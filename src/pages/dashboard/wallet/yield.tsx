import { useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import Modal from '../../../components/modal';
import Header from '../../../components/header';
import { Asset } from '../../../types/asset';
import { useSelector } from 'react-redux';
import { RootState } from '../../../store/reduxStore';
import { Wallet } from '../../../types/wallet';
import Tabs from '../../../components/tabs';
import { CustomTable } from '../../../components';
import ButtonSecondary from '../../../components/buttonSecondary';
import { showNotification } from '../../../utils/showToaster';

type Row = {
  amount: string;
  wallet: string;
  publicKey: string;
  asset?: Asset;
  transactionId: string;
  blockchainProof: string;
  date: string;
};

export default function Yield() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const appUser = useSelector((state: RootState) => state.auth.user!);

  const [wallet] = useState<Wallet>(
    appUser.userWallets.find((w) => w.publicKey == searchParams.get('wallet'))!,
  );
  const [asset] = useState<Asset>(
    wallet?.claimedAssets.find(
      (a) =>
        a.assetCode == searchParams.get('assetCode') &&
        a.assetIssuer == searchParams?.get('assetIssuer'),
    )!,
  );
  const [activeRow, setActiveRow] = useState<Row>();
  const [showSuccessModal, setShowSuccessModal] = useState(false);

  const column = [
    {
      header: 'Amount',
      key: 'amount',
    },
    {
      header: 'Wallet',
      key: 'wallet',
    },
    {
      header: 'Transaction ID',
      key: 'transactionId',
    },
    {
      header: 'Date',
      key: 'date',
    },
  ];
  const dataSource = [
    {
      amount: '$100',
      wallet: 'Received',
      publicKey: 'GDZT5...LY320',
      asset: undefined,
      transactionId: 'Deposit to wallet',
      blockchainProof: 'string',
      date: '2023-10-01',
    },
    {
      amount: '$50',
      wallet: 'Sent',
      publicKey: 'GDZT5...LY320',
      asset: undefined,
      transactionId: 'Withdrawal from wallet',
      blockchainProof: 'string',
      date: '2023-10-02',
    },
    {
      amount: '$50',
      wallet: 'Sent',
      publicKey: 'GDZT5...LY320',
      asset: undefined,
      transactionId: 'Withdrawal from wallet',
      blockchainProof: 'string',
      date: '2023-10-02',
    },
    {
      amount: '$50',
      wallet: 'Swap',
      publicKey: 'GDZT5...LY320',
      asset: undefined,
      transactionId: 'Withdrawal from wallet',
      blockchainProof: 'string',
      date: '2023-10-02',
    },
  ];

  return (
    <div className="flex text-primary-800 text-sm md:text-md flex-col space-y-5 p-3">
      <Header />
      <div className="w-full items-center">
        <div className="flex px-5 mb-10 space-x-5 w-full items-center">
          <button type="button" onClick={() => navigate(-1)}>
            <img src="/images/arrowBack.png" alt="arrow back" />
          </button>
          <p className="font-montserratSemiBold text-lg xl:text-xl">
            Dividend and Interest
          </p>
        </div>
        <div className="flex md:h-screen w-full items-start">
          <div className="w-2/5 md:p-3">
            <div className="flex md:pt-5 p-5 rounded-xl flex-col items-center bg-primary-100 max-w-5xl min-h-full">
              <div className="w-full space-y-6 px-3 md:px-5 md:py-5 flex flex-col justify-center space-y-5 md:space-y-3">
                <div className="flex flex justify-center items-center space-x-1 mb-10">
                  <img
                    src={asset?.imageUrl}
                    alt="asset logo"
                    className="h-10 w-10 rounded-full"
                  />
                  <p className="font-montserratSemiBold text-lg xl:text-xl">
                    {asset?.assetCode}
                  </p>
                </div>
                <div className="flex flex-col justify-between w-full space-y-12">
                  <div id="dividends" className="w-full space-y-5">
                    <p className="font-montserratSemiBold text-lg xl:text-xl">
                      Dividend and Interest
                    </p>
                    <p className="flex w-full justify-between">
                      <span>Dividend frequency</span>
                      <span className="font-montserratSemiBold text-md">
                        Quarterly
                      </span>
                    </p>
                    <p className="flex justify-between">
                      <span>Total dividends received</span>
                      <span className="font-montserratSemiBold text-md">
                        N45,000
                      </span>
                    </p>
                    <p className="flex justify-between">
                      <span>Last dividend amount</span>
                      <span className="font-montserratSemiBold text-md">
                        N7,000
                      </span>
                    </p>
                    <p className="flex justify-between">
                      <span>Last dividend date</span>
                      <span className="font-montserratSemiBold text-md">
                        25 June, 2025
                      </span>
                    </p>
                  </div>
                  <div id="yields" className="w-full space-y-5">
                    <p className="font-montserratSemiBold text-lg xl:text-xl">
                      Yield Details
                    </p>
                    <p className="flex w-full justify-between">
                      <span>Interest frequency</span>
                      <span className="font-montserratSemiBold text-md">
                        Quarterly
                      </span>
                    </p>
                    <p className="flex justify-between">
                      <span>Total interest earned</span>
                      <span className="font-montserratSemiBold text-md">
                        N45,000
                      </span>
                    </p>
                    <p className="flex justify-between">
                      <span>Last payment date</span>
                      <span className="font-montserratSemiBold text-md">
                        25 June, 2025
                      </span>
                    </p>
                    <p className="flex justify-between">
                      <span>Accrued interest (unpaid)</span>
                      <span className="font-montserratSemiBold text-md">
                        N1500
                      </span>
                    </p>
                    <p className="flex justify-between">
                      <span>Next payment date</span>
                      <span className="font-montserratSemiBold text-md">
                        25 June, 2025
                      </span>
                    </p>
                    <p className="flex justify-between">
                      <span>YTD Effective Yield</span>
                      <span className="font-montserratSemiBold text-md">
                        10.20%
                      </span>
                    </p>
                    <p className="flex justify-between">
                      <span>Expected annual yield</span>
                      <span className="font-montserratSemiBold text-md">
                        10.50%
                      </span>
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div className="h-full w-full md:p-3">
            <div className="flex md:pt-5 p-5 rounded-lg flex-col items-start bg-primary-100 max-w-5xl min-h-full">
              <div className="flex justify-between items-start space-x-5">
                <Tabs
                  tabList={['Dividend History', 'Yield History']}
                  onTabChanged={(index) => {}}
                />
              </div>
              <div className="w-full">
                <CustomTable
                  columns={column}
                  data={dataSource}
                  onRowClick={(row) => {
                    setShowSuccessModal(true);
                    setActiveRow(row);
                  }}
                />
              </div>
              <Modal
                showModal={showSuccessModal}
                onClose={() => {
                  setShowSuccessModal(false);
                }}
              >
                <div className="flex flex-col items-center w-full px-5 md:px-20 space-y-5 py-5 justify-center">
                  <div className="flex flex-col text-center space-y-5 items-center w-2/3 md:px-10 justify-center">
                    <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
                      Dividend Payment Details
                    </p>
                    <p className="text-green-500 text-md xl:text-lg font-montserratSemiBold">
                      + {activeRow?.amount ?? 0.0}
                    </p>
                  </div>
                  <>
                    <div className="px-5 w-full bg-primary-100 rounded-xl py-5 space-y-2">
                      <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
                        Receiving Wallet
                      </p>
                      <div className="flex space-x-3 w-full justify-between items-between">
                        <p>{activeRow?.wallet}</p>{' '}
                        <button
                          type="button"
                          onClick={() =>
                            navigator.clipboard
                              .writeText(activeRow?.wallet ?? 'Wallet')
                              .then(() => {
                                showNotification('info', 'Username copied!');
                              })
                          }
                        >
                          <img src="/images/copy.png" alt="copy" />
                        </button>
                      </div>
                      <div className="flex space-x-3 w-full justify-between items-between">
                        <p>{activeRow?.publicKey}</p>{' '}
                        <button
                          type="button"
                          onClick={() =>
                            navigator.clipboard
                              .writeText(activeRow?.wallet ?? 'Wallet')
                              .then(() => {
                                showNotification('info', 'Username copied!');
                              })
                          }
                        >
                          <img src="/images/copy.png" alt="copy" />
                        </button>
                      </div>
                      <hr className="border-1" />
                      <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
                        Asset Token
                      </p>
                      <div className="flex space-x-3 w-full justify-between items-between">
                        <p>{activeRow?.publicKey}</p>{' '}
                      </div>
                      <hr className="border-1" />
                      <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
                        Blockchain Proof (Transaction ID)
                      </p>
                      <div className="flex space-x-3 w-full justify-between items-between">
                        <a
                          className="underline"
                          // href={`${getExplorerBaseUrl(appState.walletMode)}${
                          //   formData.transactionData?.transactionId
                          // }`}
                          href="#"
                          target="_blank"
                        >
                          {activeRow?.blockchainProof}
                        </a>
                        <button
                          type="button"
                          onClick={() =>
                            navigator.clipboard
                              .writeText(activeRow?.wallet ?? 'Wallet')
                              .then(() => {
                                showNotification('info', 'Username copied!');
                              })
                          }
                        >
                          <img src="/images/copy.png" alt="copy" />
                        </button>
                      </div>
                      <hr className="border-1" />
                      <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
                        Date
                      </p>
                      <div className="flex space-x-3 w-full justify-between items-between">
                        <a
                          className="underline"
                          // href={`${getExplorerBaseUrl(appState.walletMode)}${
                          //   formData.transactionData?.transactionId
                          // }`}
                          href="#"
                          target="_blank"
                        >
                          {activeRow?.date}
                        </a>
                      </div>
                    </div>
                    <div className="w-full space-y-3">
                      <a
                        href={`#`} // since we use hash router
                        // href={`/#/send-asset-receipt?q=${receiptQuery}`} // since we use hash router
                        target="_blank"
                        className="bg-primary-800 block rounded-lg w-full text-center text-white h-12 py-4 px-5"
                      >
                        Generate receipt
                      </a>
                      <ButtonSecondary
                        label="Close"
                        additionalClasses="font-montserratSemiBold"
                        onclick={async () => {
                          setShowSuccessModal(false);
                        }}
                      />
                    </div>
                  </>
                  <div />
                </div>
              </Modal>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
