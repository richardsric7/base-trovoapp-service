import { useState } from 'react';
import { useDispatch } from 'react-redux';
import TrovoBrand from './trovoBrand';
import SideBarItem from './sidebarItem';
import { toggleSidebar } from '../store/sidebarSlice';
import { Encryptor } from '../utils/encryptor';

type Props = {
  mobileMode?: boolean;
};

function SideBar({ mobileMode = false }: Props) {
  const [activeItem, setActiveItem] = useState(1);
  const dispatch = useDispatch();
  const classes = `flex flex-col pb-10 overflow-y-scroll h-full bg-primary-800 ${
    mobileMode ? 'rounded-tr-3xl rounded-br-3xl' : 'rounded-3xl'
  }`;
  return (
    <div className={classes}>
      <div className="m-3 flex justify-between items-center">
        <TrovoBrand textColor="xl:text-sm 2xl:text-lg text-white" />
        <button
          className="mt-5 mb-10 xl:hidden w-1/4"
          type="button"
          onClick={() => dispatch(toggleSidebar())}
        >
          <span className="text-white text-xl">X</span>
        </button>
      </div>
      <div className="flex flex-col justify-between h-screen">
        <div className="mx-3 2xl:mx-7 space-y-4 flex flex-col">
          <SideBarItem
            label="Home"
            icon={
              activeItem === 1 ? '/images/home.svg' : '/images/homeGray.svg'
            }
            url="/dashboard"
            isActive={activeItem === 1}
            onSidebarClicked={() => {
              setActiveItem(1);
              console.log('activeItem set to', activeItem);
            }}
          />
          <SideBarItem
            label="Wallets"
            icon={
              activeItem === 2
                ? '/images/walletWhite.svg'
                : '/images/wallet.svg'
            }
            url="/dashboard/wallet"
            isActive={activeItem === 2}
            onSidebarClicked={() => {
              setActiveItem(2);
              console.log('activeItem set to', activeItem);
            }}
          />
          <SideBarItem
            label="History"
            icon={
              activeItem === 3
                ? '/images/historyWhite.svg'
                : '/images/history.svg'
            }
            url="/dashboard/history"
            isActive={activeItem === 3}
            onSidebarClicked={() => {
              setActiveItem(3);
              console.log('activeItem set to', activeItem);
            }}
          />
          <SideBarItem
            label="Tokenize"
            icon={
              activeItem === 4
                ? '/images/tokenizeWhite.svg'
                : '/images/tokenize.svg'
            }
            url="/dashboard/tokenize"
            isActive={activeItem === 4}
            onSidebarClicked={() => {
              setActiveItem(4);
            }}
          />
          <SideBarItem
            label="Shared Access"
            icon={
              activeItem === 5
                ? '/images/accessWhite.svg'
                : '/images/access.svg'
            }
            url="/dashboard/shared-access"
            isActive={activeItem === 5}
            onSidebarClicked={() => {
              setActiveItem(5);
            }}
          />
          {/* <SideBarItem
          label="Market Trade"
          icon={
            activeItem === 6
              ? '/images/marketTradeWhite.svg'
              : '/images/marketTrade.svg'
          }
          url="/dashboard/market-trade"
          isActive={activeItem === 6}
          onSidebarClicked={() => {
            setActiveItem(6);
          }}
        /> */}
          <SideBarItem
            label="Trovo Patron"
            icon={
              activeItem === 7
                ? '/images/trovoPatronWhite.svg'
                : '/images/trovoPatron.svg'
            }
            url="/dashboard/trovo-patron"
            isActive={activeItem === 7}
            onSidebarClicked={() => {
              setActiveItem(7);
            }}
          />
          <SideBarItem
            label="Add/Remove Asset"
            icon={
              activeItem === 8
                ? '/images/addAssetWhite.svg'
                : '/images/addAsset.svg'
            }
            url="/dashboard/add-remove-assets"
            isActive={activeItem === 8}
            onSidebarClicked={() => {
              setActiveItem(8);
            }}
          />
          <SideBarItem
            label="Import Wallet"
            icon={
              activeItem === 9
                ? '/images/importWhite.svg'
                : '/images/import.svg'
            }
            url="/import"
            isActive={activeItem === 9}
            onSidebarClicked={() => {
              setActiveItem(9);
            }}
          />
          <SideBarItem
            label="Backup Wallet"
            icon={
              activeItem === 10
                ? '/images/backupWhite.svg'
                : '/images/backup.svg'
            }
            url="/backup"
            isActive={activeItem === 10}
            onSidebarClicked={() => {
              setActiveItem(10);
            }}
          />
          <SideBarItem
            label="Account Recovery"
            icon={
              activeItem === 11
                ? '/images/historyWhite.svg'
                : '/images/history.svg'
            }
            url="/recovery"
            isActive={activeItem === 11}
            onSidebarClicked={() => {
              setActiveItem(11);
            }}
          />
          {/* <SideBarItem
          label="Closed Groups"
          icon={
            activeItem === 12
              ? '/images/closedGroupsWhite.svg'
              : '/images/closedGroups.svg'
          }
          url="/dashboard/closed-groups"
          isActive={activeItem === 12}
          onSidebarClicked={() => {
            setActiveItem(12);
          }}
        /> */}
          {/* <SideBarItem
          label="Settings"
          icon={
            activeItem === 13
              ? '/images/settingsWhite.svg'
              : '/images/settings.svg'
          }
          url="/dashboard/settings"
          isActive={activeItem === 13}
          onSidebarClicked={() => {
            setActiveItem(13);
          }}
        /> */}
        </div>
        <div className="mx-3 2xl:mx-7 space-y-4">
          <SideBarItem
            label="Logout"
            icon={
              activeItem === 1 ? '/images/logout.png' : '/images/logout.png'
            }
            url="/login"
            isActive={false}
            onSidebarClicked={async () => {
              const encryptor = new Encryptor();
              let user = await encryptor.decryptUserData();
              if (user?.isLoggedIn) {
                user = { ...user, isLoggedIn: false };
                const encryptor = new Encryptor();
                await encryptor.encryptUserData(user);
                console.log('logged out');
              } else {
                console.log('user is not logged in.');
              }
            }}
          />
        </div>
      </div>
    </div>
  );
}

export default SideBar;
