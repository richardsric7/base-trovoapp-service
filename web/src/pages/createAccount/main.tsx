import { Outlet, useLocation } from 'react-router-dom';
import TrovoBrand from '../../components/trovoBrand';

function AuthMain() {
  return (
    <div className="flex h-full items-center justify-center ">
      <div className="hidden md:block w-3/5 h-full p-3">
        <div className="flex h-full space-y-3 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
          <TrovoBrand />
          <p className="text-primary-800 font-matahariExtended text-center text-3xl xl:text-4xl font-bold">
            Let&apos;s get you
          </p>
          <p className="text-primary-700 font-matahariExtended text-center text-3xl xl:text-4xl font-bold">
            started!
          </p>
          <img
            className="w-100 h-100"
            src="/images/welcome.png"
            alt="Welcome 1"
          />
        </div>
      </div>
      <div className="w-full md:w-2/5 h-full overflow-y-scroll">
        <div className="flex flex-col space-y-6 h-full overflow-y-scroll items-center md:justify-center">
          <ShowImageOnMobile />
          <Outlet />
        </div>
      </div>
    </div>
  );
}

function ShowImageOnMobile() {
  const location = useLocation();
  const isFormView = location.pathname.includes('/register/form');
  const className = 'h-2/4 w-full p-3';

  return (
    <div className={`${className} ${isFormView ? 'hidden' : 'md:hidden'}`}>
      <div className="flex h-full w-full space-y-2 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
        <TrovoBrand />
        <p className="text-primary-800 font-matahariExtended text-center text-2xl xl:text-4xl font-bold">
          Let&apos;s get you
        </p>
        <p className="text-primary-700 font-matahariExtended text-center text-2xl xl:text-4xl font-bold">
          started!
        </p>
        <img className="w-2/4" src="/images/welcome.png" alt="Welcome 1" />
      </div>
    </div>
  );
}

export default AuthMain;
