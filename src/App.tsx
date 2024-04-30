import './App.css';
import AppRouter from './routingSetup/appRouter';
import Toaster from './components/toaster';
import Loader from './components/loader';
import { useNavigate } from 'react-router-dom';
import { RootState } from './store/reduxStore';
import { useDispatch, useSelector } from 'react-redux';
import { hideToaster } from './store/sidebarSlice';
import { useEffect, useState } from 'react';
import { Encryptor } from './utils/encryptor';
import { getStorage } from './utils/storage';
import { USER_DETAILS } from './store/constants';
import { User } from './types/user';
import { setTempUser } from './store/authSlice';

function App() {
  const dispatch = useDispatch();
  const [isLoading, setIsLoading] = useState(true);

  const toasterInfo = useSelector((state: RootState) => {
    return state.sidebarSlice.showToaster;
  });

  const loaderState = useSelector((state: RootState) => {
    return state.sidebarSlice.showLoader;
  });

  useEffect(() => {
    const encryptor = new Encryptor();
    const storedInfo = getStorage(USER_DETAILS);
    if (!storedInfo) {
      setIsLoading(false);
      return;
    }
    encryptor
      .decryptData(
        storedInfo.__slw31H408,
        storedInfo.__39deR7sx4,
        storedInfo.__i34dcY9Mn,
      )
      .then((result) => {
        const user = JSON.parse(result) as User;
        dispatch(
          setTempUser({
            ...user,
          }),
        );
        setIsLoading(false);
      });
  });

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
        onClose={() => {
          dispatch(hideToaster());
        }}
      />
      <Loader showLoader={loaderState} />
    </div>
  );
}

export default App;
