import PropTypes from 'prop-types';
import { React, useState } from 'react';

export default function PasswordInput({
  defaultValue,
  placeholder,
  label,
  leadingIcon,
  inputType,
  onInputChange,
}) {
  const [showPlainText, setShowPlainText] = useState(false);

  const trailingIcon = '/images/eyeShow.png';

  const [value, setValue] = useState(defaultValue);
  const handleInputChange = (event) => {
    const newValue = event.target.value;
    setValue(newValue);
    console.log('input changed...', newValue);
    onInputChange(newValue);
  };

  return (
    <>
      <label className="text-primary-800" htmlFor={label}>
        {label}
      </label>
      <div
        className="mt-2 ring-2 ring-gray-200 focus-within:ring-primary-600 rounded-md
          w-full h-12 py-1 px-2 focus-within:ring-2 flex items-center"
      >
        {leadingIcon && <img src={leadingIcon} alt="" />}
        <input
          name={label}
          className="h-full px-2 focus:outline-none w-full"
          placeholder={placeholder}
          type={inputType !== 'password' || showPlainText ? 'text' : 'password'}
          value={value}
          onChange={handleInputChange}
        />
        {inputType === 'password' && (
          <button
            type="button"
            aria-label="show password"
            onClick={() => {
              setShowPlainText(!showPlainText);
            }}
          >
            <img src={trailingIcon} alt="" />
          </button>
        )}
      </div>
    </>
  );
}

PasswordInput.propTypes = {
  placeholder: PropTypes.string,
  inputType: PropTypes.string,
  leadingIcon: PropTypes.string,
  label: PropTypes.string.isRequired,
  defaultValue: PropTypes.string,
  onInputChange: PropTypes.func.isRequired,
};

PasswordInput.defaultProps = {
  placeholder: 'Enter value',
  leadingIcon: null,
  inputType: 'text',
  defaultValue: '',
};
