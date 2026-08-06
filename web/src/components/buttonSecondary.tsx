type Props = {
  onclick: () => void;
  label: string;
  additionalClasses?: string;
};

function ButtonSecondary({
  onclick,
  label,
  additionalClasses = 'bg-white',
}: Props) {
  const classes = `ring-1 ring-primary-800 rounded-md w-full text-primary-800 h-12 ${additionalClasses}`;
  return (
    <button className={classes} type="button" onClick={onclick}>
      {label}
    </button>
  );
}

export default ButtonSecondary;
