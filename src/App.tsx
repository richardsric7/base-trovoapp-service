import './App.css';
import AppRouter from './routingSetup/appRouter';
import Toaster from './components/toaster';
import Loader from './components/loader';
import { RootState } from './store/reduxStore';
import { useDispatch, useSelector } from 'react-redux';
import { hideToaster } from './store/sidebarSlice';

function App() {
  const dispatch = useDispatch();

  const toasterInfo = useSelector((state: RootState) => {
    return state.sidebarSlice.showToaster;
  });

  const loaderState = useSelector((state: RootState) => {
    return state.sidebarSlice.showLoader;
  });

  return (
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
