import { useSelector } from 'react-redux';
import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { RootState } from '../store/reduxStore';

export const ProtectedRoutes = () => {
  const { user } = useSelector((state: RootState) => state.auth);
  const location = useLocation();

  return user?.isLoggedIn ? (
    <Outlet />
  ) : (
    <Navigate
      to={user ? '/login' : '/welcome'}
      state={{ from: location }}
      replace
    />
  );
};
