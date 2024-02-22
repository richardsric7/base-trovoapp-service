import PropTypes from 'prop-types';
import React from 'react';

function Button({ onclick, label, additionalClasses }) {
  const classes = `bg-primary-800 rounded-lg w-full text-white h-12 ${additionalClasses}`;
  return (
    <button className={classes} type="button" onClick={onclick}>
      {label}
    </button>
  );
}

Button.propTypes = {
  onclick: PropTypes.func.isRequired,
  label: PropTypes.string.isRequired,
  additionalClasses: PropTypes.string,
};

Button.defaultProps = {
  additionalClasses: '',
};

export default Button;
