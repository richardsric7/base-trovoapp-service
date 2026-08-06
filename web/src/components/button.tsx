type Props = {
  additionalClasses?: string;
  label: string;
  disabled?: boolean;
  onclick: () => void;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
  type?: 'button' | 'submit' | 'reset' | undefined;
};

function Button({
  onclick,
  label,
  disabled = false,
  type = 'button',
  additionalClasses = '',
  leftIcon,
  rightIcon
}: Props) {
  const classes = `bg-primary-800 rounded-lg flex items-center justify-center gap-3 w-full text-white h-12 disabled:bg-primary-400 ${additionalClasses}`;
  return (
    <button
      disabled={disabled}
      className={classes}
      type={type}
      onClick={onclick}
    >
      {leftIcon}
      {label}
      {rightIcon}
    </button>
  );
}

export default Button;
