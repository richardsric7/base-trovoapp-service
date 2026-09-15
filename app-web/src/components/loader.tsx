import { Oval } from 'react-loader-spinner';

type Props = {
  showLoader: boolean;
};

export default function Loader({ showLoader }: Props) {
  return (
    <div>
      {showLoader ? (
        <>
          <div className="justify-center items-center flex overflow-x-hidden overflow-y-auto fixed inset-0 z-50 outline-none focus:outline-none">
            <Oval
              visible={true}
              height="80"
              width="80"
              color="#F2F6F9"
              secondaryColor="#336DA0"
              ariaLabel="oval-loading"
            />
          </div>
          <div className="opacity-25 fixed inset-0 z-40 bg-black" />
        </>
      ) : null}
    </div>
  );
}
