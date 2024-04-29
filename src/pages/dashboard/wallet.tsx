import Button from '../../components/button';
import TransactionItem from '../../components/transactionItem';
import WidgetCard from '../../components/widgetCard';
import Carousel from '../../components/carousel';
import Header from '../../components/header';
import Tabs from '../../components/tabs';

export default function Wallet() {
  // const style = {
  //   backfaceVisibility: 'hidden',
  // };

  return (
    <div className="flex text-primary-800 text-sm overflow-x-hidden md:text-md flex-col space-y-5 p-3">
      <Header
        fullName="Obi Enechi"
        avatar="/images/avatar.png"
        email="obienechi@gmail.com"
        isHomeView
      />
      <div className="flex md:h-screen w-full items-center justify-center">
        <div className="h-full w-full md:p-3">
          <div className="md:hidden items-center px-5 mb-5 justify-center flex">
            <WidgetCard />
          </div>
          <div className="flex space-y-3 mb-20 md:pt-5 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
            <div className="w-full flex justify-between px-5">
              <Tabs tabList={['My Wallets', 'Shared Wallets', 'All Wallets']} />
              <div className="w-1/4">
                <Button
                  label="Add Subwallet"
                  onclick={() => {
                    /* */
                  }}
                />
              </div>
            </div>
            {/* <div className="hidden md:flex w-[900px] md:w-[1000px] items-center
            space-x-5 overflow-x-auto px-5 justify-center">
              <WalletCard />
              <WalletCard />
              <WalletCard />
            </div> */}
            <Carousel />
          </div>
        </div>
        <div className="hidden md:block w-3/6 flex flex-col space-y-5 h-full py-3 lg:px-3">
          <div className="flex space-y-3 rounded-lg py-5 flex-col items-center bg-primary-100">
            <p>Current Plan</p>
            <div className="flex space-x-1">
              <img src="/images/gold.png" alt="gold" />
              <p className="font-bold">Gold Patron (Monthly)</p>
            </div>
            <div className="w-full px-3 lg:px-10">
              <div
                className="max-w-xl w-full flex flex-col xl:flex-row text-white font-matahariRegular text-xl xl:text-2xl
              rounded-2xl bg-primary-800 items-center justify-center space-y-3 xl:space-y-0 xl:items-start xl:justify-between xl:space-x-3 mt-5 xl:mb-5 py-5 xl:py-10 xl:px-5"
              >
                <img
                  className="w-10"
                  src="/images/platinum.png"
                  alt="platinum"
                />
                <div className="flex flex-col space-y-2 text-center xl:text-start">
                  <p className="text-lg">Try the Diamond Plan</p>
                  <p className="text-sm">Upgrade to unlock more features</p>
                </div>
                <img className="w-10" src="/images/go.png" alt="platinum" />
              </div>
            </div>
          </div>
          <div className="flex w-full space-y-3 py-7 px-2 xl:px-5 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
            <div className="flex items-center w-full justify-between">
              <p className="font-semibold">Recent Transactions</p>
              <p className="test-xs text-primary-400">view all</p>
            </div>
            <TransactionItem
              addressOrUsername="Tom"
              transactionType={1}
              amount="3,000.000"
              assetCode="XBN"
              date={new Date().toLocaleString()}
              onclick={() => {
                // setHasAssets(!hasAssets);
              }}
            />
            <TransactionItem
              addressOrUsername="GH3Y8...ZX89O"
              transactionType={0}
              amount="0.5000"
              assetCode="TROV"
              date={new Date().toLocaleString()}
              onclick={() => {
                /* */
              }}
            />
            <TransactionItem
              addressOrUsername="Tom"
              transactionType={1}
              amount="3,000.000"
              assetCode="XBN"
              date={new Date().toLocaleString()}
              onclick={() => {
                /* */
              }}
            />
            <TransactionItem
              addressOrUsername="GH3Y8...ZX89O"
              transactionType={0}
              amount="0.5000"
              assetCode="TROV"
              date={new Date().toLocaleString()}
              onclick={() => {
                /* */
              }}
            />
          </div>
          <div className="flex w-full space-y-3 py-7 px-5 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
            <img src="/images/bankNotes.png" alt="bank notes" />
            <p className="text-center font-semibold">
              Make Deposits and Withdrawals on your assets with ease on Trovo
              Wallet
            </p>
            <div className="w-3/4">
              <Button
                label="Deposit/Withdraw"
                onclick={() => {
                  /* */
                }}
              />
            </div>
          </div>
          <div />
          <div />
          <div />
          <div />
        </div>
      </div>
    </div>
  );
}
