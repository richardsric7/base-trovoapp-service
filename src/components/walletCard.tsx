import { forwardRef, useState } from 'react';

type Props = {
  localCurrencyBalance: string;
  usdBalance: string;
  currency: string;
  alias: string;
};

const WalletCard = forwardRef<HTMLDivElement, Props>(
  ({ localCurrencyBalance, usdBalance, alias, currency }: Props, ref) => {
    const [hideBalance, setHideBalance] = useState(false);
    return (
      <div
        className="min-w-[900px] flex text-white font-matahariRegular
     rounded-2xl bg-primary-800 items-center justify-between space-x-3 md:mt-5 md:mb-10 py-5 px-5"
        ref={ref}
      >
        <div className="flex flex-col space-y-2">
          <div className="flex space-x-10 items-center">
            <p className="md:text-lg">{alias}</p>
            <button onClick={() => setHideBalance(!hideBalance)}>
              <img
                width="25"
                height="3"
                src={hideBalance ? '/images/show.png' : '/images/hide.png'}
                alt="trovo logo"
              />
            </button>
          </div>
          <div>Total Balance</div>
          <p className="text-xl md:text-2xl font-semibold">
            {hideBalance ? '**********' : `${localCurrencyBalance} ${currency}`}
          </p>
          {currency.toLowerCase().includes('usd') ? (
            ''
          ) : (
            <p className="text-sm">
              {hideBalance ? '**********' : `${usdBalance} USD`}
            </p>
          )}
        </div>
        <img className="w-20" src="/images/trovoWhite.png" alt="trovo logo" />
      </div>
    );
  },
);

export default WalletCard;
