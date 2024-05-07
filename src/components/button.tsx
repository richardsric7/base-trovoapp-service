type Props = {
  additionalClasses?: string;
  label: string;
  disabled?: boolean;
  onclick: () => void;
  type?: 'button' | 'submit' | 'reset' | undefined;
};

function Button({
  onclick,
  label,
  disabled = false,
  type = 'button',
  additionalClasses = '',
}: Props) {
  const classes = `bg-primary-800 rounded-lg w-full text-white h-12 disabled:bg-primary-400 ${additionalClasses}`;
  return (
    <button
      disabled={disabled}
      className={classes}
      type={type}
      onClick={onclick}
    >
      {label}
    </button>
  );
}

export default Button;
