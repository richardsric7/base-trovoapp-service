type Props = {
  onclick: () => void;
  label: string;
  additionalClasses?: string;
  disabled?: boolean;
};

function ButtonSecondary({
  onclick,
  label,
  additionalClasses = 'bg-white',
  disabled = false,
}: Props) {
  const classes = `ring-1 ring-primary-800 rounded-md w-full text-primary-800 h-12 disabled:opacity-50 disabled:cursor-not-allowed ${additionalClasses}`;
  return (
    <button disabled={disabled} className={classes} type="button" onClick={onclick}>
      {label}
    </button>
  );
}

export default ButtonSecondary;
