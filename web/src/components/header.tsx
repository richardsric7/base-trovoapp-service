import { useDispatch, useSelector } from 'react-redux';
import { toggleSidebar } from '../store/sidebarSlice';
import TextInput from './textInput';
import { RootState } from '../store/reduxStore';
import capitalizeFirstLetter from '../utils/capitalizeFirst';

type Props = {
  isHomeView?: boolean;
};

// eslint-disable-next-line
function Header({ isHomeView = false }: Props) {
  const dispatch = useDispatch();
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const classes = `flex w-full p-3 items-center ${
    isHomeView ? 'justify-between' : 'justify-between xl:justify-end'
  }`;
  return (
    <div className={classes}>
      {isHomeView ? (
        <>
          <div className="flex space-x-4">
            <button
              className="xl:hidden"
              type="button"
              onClick={() => dispatch(toggleSidebar())}
            >
              <img src="/images/hamburgerMenu.png" alt="copy" className="w-8" />
            </button>
            <p>
              Good day,
              <span className="font-semibold text-lg">
                {' '}
                {capitalizeFirstLetter(appUser.firstName)}
              </span>
            </p>
          </div>
          <div className="hidden md:block w-1/4">
            <TextInput
              leadingIcon="/images/search.png"
              inputType="text"
              label=""
              placeholder="Search"
              onInputChange={(newValue: string) => {
                console.log('input has changed', newValue);
              }}
            />
          </div>
        </>
      ) : (
        <button
          className="xl:hidden"
          type="button"
          onClick={() => dispatch(toggleSidebar())}
        >
          <img src="/images/hamburgerMenu.png" alt="copy" className="w-8" />
        </button>
      )}
      <div className="flex space-x-3 w-1/5 items-center">
        <img
          className="h-6"
          src="/images/notification.png"
          alt="notification bell"
        />
        <img
          className="h-10 rounded-full"
          src={
            appUser.imageThumbnailURL.length > 0
              ? appUser.imageThumbnailURL
              : '/images/avatar.png'
          }
          alt="avatar"
        />
        <div className="flex w-64 flex-col hidden w-56 md:block space-y-2">
          <p className="font-semibold truncate">
            {capitalizeFirstLetter(appUser.firstName)}{' '}
            {capitalizeFirstLetter(appUser.lastName)}
          </p>
          <p className="truncate">{appUser.email}</p>
        </div>
      </div>
    </div>
  );
}

export default Header;
