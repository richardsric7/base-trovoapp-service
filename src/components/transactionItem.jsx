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
  const amountClasses = `text-xl ${
    transactionType === 0 ? 'text-red-500' : 'text-green-500'
  }`;
  return (
    <button
      className="rounded-2xl w-full bg-white flex justify-between items-center px-7 py-5"
      onClick={onclick}
      type="button"
    >
      <div className="flex space-x-5 items-center">
        <img
          src={
            transactionType === 0 ? '/images/sent.png' : '/images/received.png'
          }
          alt="direction"
        />
        <div className="flex items-start flex-col">
          <p className="font-semibold text-xl">
            {transactionType === 0 ? 'To ' : 'From '}
            {addressOrUsername}
          </p>
          <p>{transactionType === 0 ? 'Sent' : 'Received'}</p>
        </div>
      </div>
      <div className="flex flex-col items-end space-x-3">
        <p className={amountClasses}>
          {`${transactionType === 0 ? '-' : '+'} ${amount} ${assetCode}`}
        </p>
        <p className="text-sm">{date}</p>
      </div>
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
