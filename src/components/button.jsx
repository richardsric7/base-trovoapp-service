import PropTypes from 'prop-types';
import React from 'react';

function Button({ onclick, label }) {
  return (
    <button
      className="bg-primary-800 rounded-md w-full text-white h-12"
      type="button"
      onClick={onclick}
    >
      {label}
    </button>
  );
}

Button.propTypes = {
  onclick: PropTypes.func.isRequired,
  label: PropTypes.string.isRequired,
};

export default Button;
