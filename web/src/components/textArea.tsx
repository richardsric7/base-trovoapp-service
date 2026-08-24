import { useState } from 'react';

type Props = {
  defaultValue?: string;
  placeholder?: string;
  label: string;
  leadingIcon?: string;
  inputType: string;
  error?: string;
  readonly?: boolean;
  onInputChange: (newValue: string) => void;
  trailingIcon?: string;
  trailingText?: string;
  rows?: number;
};

export default function TextArea({
  defaultValue = '',
  placeholder = 'Enter value',
  label,
  leadingIcon,
  inputType = 'text',
  error = '',
  readonly,
  onInputChange,
  trailingIcon,
  trailingText,
  rows = 2,
}: Props) {
  const [showPlainText, setShowPlainText] = useState(false);

  const [value, setValue] = useState(defaultValue);
  const handleInputChange = (event: { target: { value: any } }) => {
    const newValue = event.target.value;
    setValue(newValue);
    onInputChange(newValue);
  };

  return (
    <div>
      <label className="text-primary-700" htmlFor={label}>
        {label}
      </label>
      <div
        className="mt-2 ring-1 md:ring-2 ring-gray-200 focus-within:ring-primary-600 rounded-md
          w-full py-1 px-2 focus-within:ring-2 flex flex-col"
      >
        <div className="flex items-start">
          {leadingIcon && <img src={leadingIcon} alt="" className="mt-2" />}
          <textarea
            name={`${label}-input`}
            className="px-2 focus:outline-none w-full bg-inherit text-gray-700 resize-none"
            placeholder={placeholder}
            rows={rows}
            value={value}
            onChange={handleInputChange}
            readOnly={readonly ?? false}
            autoComplete="off"
          />
          <div className="flex flex-col items-end">
            {trailingIcon && (
              <div className="flex space-x-2 items-center font-montserratSemiBold pr-6 mt-2">
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
                className="pr-2 mt-2"
              >
                <img src="/images/eyeShow.png" alt="" />
              </button>
            )}
          </div>
        </div>
      </div>
      {error && <p className="text-red-500 text-sm mt-1">{error}</p>}
    </div>
  );
}
