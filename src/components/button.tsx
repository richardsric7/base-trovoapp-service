import React from 'react';

type Props = {
  additionalClasses?: string;
  label: string;
  onclick: () => void;
};

function Button({ onclick, label, additionalClasses = '' }: Props) {
  const classes = `bg-primary-800 rounded-lg w-full text-white h-12 ${additionalClasses}`;
  return (
    <button className={classes} type="button" onClick={onclick}>
      {label}
    </button>
  );
}

export default Button;
