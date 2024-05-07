import { useState } from 'react';

type Props = {
  defaultValue?: string;
  placeholder?: string;
  label: string;
  leadingIcon?: string;
  inputType: string;
  error?: string;
  onInputChange: (newValue: string) => void;
  trailingIcon?: string;
  trailingText?: string;
};

export default function TextInput({
  defaultValue = '',
  placeholder = 'Enter value',
  label,
  leadingIcon,
  inputType = 'text',
  error = '',
  onInputChange,
  trailingIcon,
  trailingText,
}: Props) {
  const [showPlainText, setShowPlainText] = useState(false);

  const [value, setValue] = useState(defaultValue);
  const handleInputChange = (event: { target: { value: any } }) => {
    const newValue = event.target.value;
    setValue(newValue);
    onInputChange(newValue);
  };

  return (
    <>
      <label className="text-primary-700" htmlFor={label}>
        {label}
      </label>
      <div
        className="mt-2 ring-1 md:ring-2 ring-gray-200 focus-within:ring-primary-600 rounded-md
          w-full h-12 py-1 px-2 focus-within:ring-2 flex items-center"
      >
        {leadingIcon && <img src={leadingIcon} alt="" />}
        <input
          name={`${label}-input`}
          className="h-full px-2 focus:outline-none w-full"
          placeholder={placeholder}
          type={inputType !== 'password' || showPlainText ? 'text' : 'password'}
          value={value}
          onChange={handleInputChange}
          autoComplete="off"
        />
        {trailingIcon && (
          <div className="flex space-x-2 items-center font-montserratSemiBold pr-6">
            <img src={trailingIcon} alt="" />
            <span>{trailingText}</span>
          </div>
        )}
        {inputType === 'password' && (
          <button
            type="button"
            aria-label="show password"
            onClick={() => {
              setShowPlainText(!showPlainText);
            }}
          >
            <img src="/images/eyeShow.png" alt="" />
          </button>
        )}
      </div>
      {error && <p className="text-red-500 text-sm mt-1">{error}</p>}
    </>
  );
}
