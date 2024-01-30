import React from 'react';
// import PropTypes from 'prop-types';

function WidgetCard() {
  return (
    <div
      className="max-w-xl w-full flex text-white font-matahariRegular text-xl xl:text-2xl
     rounded-2xl bg-primary-800 items-center justify-between space-x-3 mt-5 mb-10 py-10 px-5"
    >
      <div className="flex flex-col space-y-2">
        <p className="text-lg">Total Account Balance</p>
        <p className="text-2xl font-semibold">2,082,898 NGN</p>
        <p className="text-sm">4,014 USD</p>
      </div>
      <img className="w-20" src="/images/trovoWhite.png" alt="trovo logo" />
    </div>
  );
}

// WidgetCard.propTypes = {
//   textColor: PropTypes.string,
// };

// WidgetCard.defaultProps = {
//   textColor: 'text-primary-800',
// };

export default WidgetCard;
