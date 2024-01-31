import PropTypes from 'prop-types';
import React from 'react';

function TransactionItem({
  addressOrUsername,
  transactionType,
  amount,
  assetCode,
  date,
  onclick,
}) {
  const amountClasses = `text-sm ${
    transactionType === 0 ? 'text-red-500' : 'text-green-500'
  }`;
  return (
    <button
      className="rounded-2xl w-full bg-white px-2 py-3"
      onClick={onclick}
      type="button"
    >
      <div className="flex text-sm justify-between items-center">
        <div className="flex space-x-2 w-full text-start items-center">
          <img
            className="w-7"
            src={
              transactionType === 0
                ? '/images/sent.png'
                : '/images/received.png'
            }
            alt="direction"
          />
          <div className="flex items-start flex-col">
            <p className="font-semibold text-sm">
              <span className="hidden xl:inline">
                {transactionType === 0 ? 'To ' : 'From '}
              </span>
              {addressOrUsername}
            </p>
            <p className="text-xs">
              {transactionType === 0 ? 'Sent' : 'Received'}
            </p>
          </div>
        </div>
        <div className="flex w-full text-end flex-col items-end space-x-3">
          <p className={amountClasses}>
            {`${transactionType === 0 ? '-' : '+'} ${amount} ${assetCode}`}
          </p>
          <p className="hidden xl:block text-xs">{date}</p>
        </div>
      </div>
      <p className="w-full xl:hidden text-end text-xs">{date}</p>
    </button>
  );
}

TransactionItem.propTypes = {
  onclick: PropTypes.func.isRequired,
  addressOrUsername: PropTypes.string.isRequired,
  transactionType: PropTypes.number.isRequired,
  amount: PropTypes.string.isRequired,
  date: PropTypes.bool.isRequired,
  assetCode: PropTypes.bool.isRequired,
};

export default TransactionItem;
