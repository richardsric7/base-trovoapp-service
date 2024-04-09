import React from 'react';

type Props = {
  additionalClasses?: string;
  label: string;
  onclick: () => void;
  type?: 'button' | 'submit' | 'reset' | undefined;
};

function Button({
  onclick,
  label,
  type = 'button',
  additionalClasses = '',
}: Props) {
  const classes = `bg-primary-800 rounded-lg w-full text-white h-12 ${additionalClasses}`;
  return (
    <button className={classes} type={type} onClick={onclick}>
      {label}
    </button>
  );
}

export default Button;
