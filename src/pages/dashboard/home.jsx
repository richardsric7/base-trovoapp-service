import { React } from 'react';
import TextInput from '../../components/textInput';
import WidgetCard from '../../components/widgetCard';
import Tabs from '../../components/tabs';
import AssetListItem from '../../components/assetListItem';
import TransactionItem from '../../components/transactionItem';
import Button from '../../components/button';

export default function Home() {
  return (
    <div className="flex text-primary-800 font-matahariRegular flex-col space-y-5 p-3">
      <div className="flex w-full p-3 justify-between items-center">
        <p>
          Good evening
          <span className="font-semibold text-lg"> Osondu</span>
        </p>
        <div className="w-1/4">
          <TextInput
            leadingIcon="/images/search.png"
            inputType="text"
            placeholder="Search"
            onInputChange={(newValue) => {
              console.log('input has changed', newValue);
            }}
          />
        </div>
        <div className="flex space-x-3 items-center">
          <img src="/images/notification.png" alt="notification bell" />
          <img src="/images/avatar.png" alt="avatar" />
          <div className="flex flex-col space-y-2">
            <p className="font-semibold">Obi Enechi</p>
            <p>obienechi@gmail.com</p>
          </div>
        </div>
      </div>
      <div className="flex md:h-screen w-full items-center justify-center">
        <div className="h-full w-full p-3">
          <div className="md:hidden items-center px-5 justify-center flex">
            <WidgetCard />
          </div>
          <div className="flex space-y-3 mb-20 pt-5 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
            <div className="hidden w-full items-center px-5 justify-center md:flex">
              <WidgetCard />
            </div>
            <p className="text-primary-800 text-center text-xl font-semibold">
              Asset Offering
            </p>
            <div className="px-5  w-full">
              <Tabs tabList={['Primary Listing', 'Secondary Listing']}>
                <div className="flex flex-col space-y-3">
                  <AssetListItem
                    image="/images/avatar.png"
                    assetName="Atlantis 1"
                    assetClass="Real Estate"
                    isSubscribed
                    onclick={() => console.log('clicked!')}
                  />
                  <AssetListItem
                    image="/images/avatar.png"
                    assetName="Orchard Estate"
                    assetClass="Real Estate"
                    onclick={() => console.log('clicked!')}
                  />
                  <AssetListItem
                    image="/images/avatar.png"
                    assetName="Atlantis 1"
                    assetClass="Real Estate"
                    isSubscribed
                    onclick={() => console.log('clicked!')}
                  />
                  <AssetListItem
                    image="/images/avatar.png"
                    assetName="Orchard Estate"
                    assetClass="Real Estate"
                    onclick={() => console.log('clicked!')}
                  />
                  <AssetListItem
                    image="/images/avatar.png"
                    assetName="Atlantis 1"
                    assetClass="Real Estate"
                    isSubscribed
                    onclick={() => console.log('clicked!')}
                  />
                  <AssetListItem
                    image="/images/avatar.png"
                    assetName="Orchard Estate"
                    assetClass="Real Estate"
                    onclick={() => console.log('clicked!')}
                  />
                  <AssetListItem
                    image="/images/avatar.png"
                    assetName="Atlantis 1"
                    assetClass="Real Estate"
                    isSubscribed
                    onclick={() => console.log('clicked!')}
                  />
                  <AssetListItem
                    image="/images/avatar.png"
                    assetName="Orchard Estate"
                    assetClass="Real Estate"
                    onclick={() => console.log('clicked!')}
                  />
                  <AssetListItem
                    image="/images/avatar.png"
                    assetName="Atlantis 1"
                    assetClass="Real Estate"
                    isSubscribed
                    onclick={() => console.log('clicked!')}
                  />
                  <AssetListItem
                    image="/images/avatar.png"
                    assetName="Orchard Estate"
                    assetClass="Real Estate"
                    onclick={() => console.log('clicked!')}
                  />
                </div>
                <div className="tab-2">
                  <div className="flex flex-col space-y-3">
                    <AssetListItem
                      image="/images/avatar.png"
                      assetName="Metro Railings"
                      assetClass="Infrastructure"
                      isSubscribed
                      onclick={() => console.log('clicked!')}
                    />
                    <AssetListItem
                      image="/images/avatar.png"
                      assetName="Mount"
                      assetClass="Infrastructure"
                      onclick={() => console.log('clicked!')}
                    />
                  </div>
                </div>
              </Tabs>
            </div>
          </div>
        </div>
        <div className="hidden md:block w-2/6 flex flex-col space-y-5 h-full p-3">
          <div className="flex space-y-3 xl:space-y-5 rounded-lg py-5 flex-col items-center bg-primary-100">
            <p>Current Plan</p>
            <div className="flex space-x-1">
              <img src="/images/gold.png" alt="gold" />
              <p className="font-bold">Gold Patron (Monthly)</p>
            </div>
            <div className="w-full px-10">
              <div
                className="max-w-xl w-full flex text-white font-matahariRegular text-xl xl:text-2xl
              rounded-2xl bg-primary-800 items-start justify-between space-x-3 mt-5 mb-5 py-10 px-5"
              >
                <img
                  className="w-10"
                  src="/images/platinum.png"
                  alt="platinum"
                />
                <div className="flex flex-col space-y-2">
                  <p className="text-lg">Try the Diamond Plan</p>
                  <p className="text-sm">Upgrade to unlock more features</p>
                </div>
                <img className="w-10" src="/images/go.png" alt="platinum" />
              </div>
            </div>
          </div>
          <div className="flex w-full space-y-3 py-7 px-5 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
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
              onclick={() => console.log('clicked!')}
            />
            <TransactionItem
              addressOrUsername="GH3Y8...ZX89O"
              transactionType={0}
              amount="0.5000"
              assetCode="TROV"
              date={new Date().toLocaleString()}
              onclick={() => console.log('clicked!')}
            />
            <TransactionItem
              addressOrUsername="Tom"
              transactionType={1}
              amount="3,000.000"
              assetCode="XBN"
              date={new Date().toLocaleString()}
              onclick={() => console.log('clicked!')}
            />
            <TransactionItem
              addressOrUsername="GH3Y8...ZX89O"
              transactionType={0}
              amount="0.5000"
              assetCode="TROV"
              date={new Date().toLocaleString()}
              onclick={() => console.log('clicked!')}
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
                  console.log('clicked!');
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
