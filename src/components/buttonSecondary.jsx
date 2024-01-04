import PropTypes from 'prop-types';
import React from 'react';

function ButtonSecondary({ onclick, label, additionalClasses }) {
  const classes = `ring-1 ring-primary-800 rounded-md w-full text-primary-800 h-12 ${additionalClasses}`;
  return (
    <button className={classes} type="button" onClick={onclick}>
      {label}
    </button>
  );
}

ButtonSecondary.propTypes = {
  onclick: PropTypes.func.isRequired,
  label: PropTypes.string.isRequired,
  additionalClasses: PropTypes.string,
};

ButtonSecondary.defaultProps = {
  additionalClasses: 'bg-white',
};

export default ButtonSecondary;
