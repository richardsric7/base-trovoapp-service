type Props = {
  additionalClasses?: string;
  label: string;
  link: string;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
};

function LinkButton({
  label,
  link,
  additionalClasses = '',
  leftIcon,
  rightIcon,
}: Props) {
  const classes = `bg-primary-800 rounded-lg flex items-center justify-center gap-3 w-full text-white h-12 disabled:bg-primary-400 ${additionalClasses}`;
  return (
    <a href={link} className={classes}>
      {leftIcon}
      {label}
      {rightIcon}
    </a>
  );
}

export default LinkButton;
