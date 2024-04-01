import React, { useRef, useState } from 'react';

type Props = {
  numberOfDigits: number;
};

function OtpInputWithValidation({ numberOfDigits }: Props) {
  const [otp, setOtp] = useState(new Array(numberOfDigits).fill(''));
  const otpBoxReference = useRef<HTMLInputElement[]>([]);

  function handleChange(value: string, index: number) {
    const newArr = [...otp];
    newArr[index] = value;
    setOtp(newArr);

    console.log('fsds', value, index);
    if (value && index < numberOfDigits - 1) {
      otpBoxReference.current[index + 1].focus();
    }
  }

  function handleBackspaceAndEnter(
    e: React.KeyboardEvent<HTMLInputElement>,
    index: number,
  ) {
    const target = e.target as HTMLInputElement;
    if (e.key === 'Backspace' && !target.value && index > 0) {
      otpBoxReference.current[index - 1].focus();
    }
    if (e.key === 'Enter' && target.value && index < numberOfDigits - 1) {
      otpBoxReference.current[index + 1].focus();
    }
  }

  const inputElements = [];

  for (let i = 0; i < numberOfDigits; i += 1) {
    inputElements.push(
      <input
        key={i}
        type="number"
        value={otp[i]}
        maxLength={1}
        onChange={(e) => handleChange(e.target.value, i)}
        onKeyUp={(e) => handleBackspaceAndEnter(e, i)}
        ref={(reference) => {
          otpBoxReference.current[i] = reference!;
        }}
        className="border remove-arrow w-10 h-10 md:w-14 md:h-auto border-2 p-3 text-center rounded-md block focus:border-2 focus:outline-none appearance-none"
      />,
    );
  }

  return <div className="flex items-center gap-4">{inputElements}</div>;
}

export default OtpInputWithValidation;
