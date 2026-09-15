type Props = {
  label: string;
  tone?: 'primary' | 'success' | 'warning' | 'danger';
};

export default function PermissionChip({ label, tone = 'primary' }: Props) {
  const toneClasses = {
    primary: 'bg-primary-100 text-primary-700',
    success: 'bg-emerald-100 text-emerald-700',
    warning: 'bg-amber-100 text-amber-700',
    danger: 'bg-rose-100 text-rose-700',
  };

  return (
    <span
      className={`rounded-full px-3 py-1 text-sm font-medium ${toneClasses[tone]}`}
    >
      {label}
    </span>
  );
}
