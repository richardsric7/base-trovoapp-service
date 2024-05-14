import './App.css';
import AppRouter from './routingSetup/appRouter';
import Toaster from './components/toaster';
import Loader from './components/loader';
import { RootState } from './store/reduxStore';
import { useDispatch, useSelector } from 'react-redux';
import { hideToaster } from './store/sidebarSlice';
import { useEffect, useState } from 'react';
import { Encryptor } from './utils/encryptor';
import Modal from './components/modal';
import { User } from './types/user';
import useIdle from './utils/useIdleTimeout';

function App() {
  const dispatch = useDispatch();
  const [showPromptModal, setShowPromptModal] = useState(false);
  const [count, setCount] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const { idleTimer } = useIdle({
    onIdle: async () => {
      let user = await decryptData();
      if (user?.isLoggedIn) {
        user = { ...user, isLoggedIn: false };
        const encryptor = new Encryptor();
        await encryptor.encryptUserData(user);

        console.log('logged out');
        setShowPromptModal(false);
      } else {
        console.log('user is not logged in.');
      }
    },
    onPrompt: async () => {
      let user = await decryptData();
      console.log('appUser', user?.isLoggedIn);
      if (user?.isLoggedIn) {
        setCount(Math.floor(idleTimer.getRemainingTime() / 1000));
        setShowPromptModal(true);
      } else {
        console.log('user is not logged in.');
      }
    },
    idleTime: 260,
  });

  const toasterInfo = useSelector((state: RootState) => {
    return state.sidebarSlice.showToaster;
  });

  const loaderState = useSelector((state: RootState) => {
    return state.sidebarSlice.showLoader;
  });

  const decryptData = async (): Promise<User | undefined> => {
    const encryptor = new Encryptor();
    return encryptor.decryptUserData();
  };
  useEffect(() => {
    const encryptor = new Encryptor();
    encryptor.decryptUserData().then(() => {
      setIsLoading(false);
    });
  }, [isLoading]);

  useEffect(() => {
    if (count > 0) {
      //Implementing the setInterval method
      const interval = setInterval(() => {
        setCount(Math.floor(idleTimer.getRemainingTime() / 1000));
      }, 1000);

      //Clearing the interval
      return () => clearInterval(interval);
    }
  }, [count]);

  return isLoading ? (
    <div className="App min-h-[900px] h-screen">
      <Loader showLoader={isLoading} />
    </div>
  ) : (
    <div className="App min-h-[900px] h-screen">
      <AppRouter />
      <Toaster
        type={toasterInfo.type}
        message={toasterInfo.message}
        delay={toasterInfo.delay}
        onClose={() => {
          dispatch(hideToaster());
        }}
      />
      <Loader showLoader={loaderState} />
      <Modal
        showModal={showPromptModal}
        onClose={() => {
          setShowPromptModal(false);
          idleTimer.reset();
        }}
      >
        <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
          <img src="/images/success.png" alt="success" />
          <div className="flex flex-col text-center space-y-5 items-center w-2/3 md:px-10 justify-center">
            <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
              You will be logged out in... {count} seconds!
            </p>
          </div>
        </div>
      </Modal>
    </div>
  );
}

export default App;
