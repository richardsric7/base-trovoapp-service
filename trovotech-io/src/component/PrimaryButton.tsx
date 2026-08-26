import Link from "next/link";
import React from "react";

interface PrimaryButtonProps {
  label: string;
  href?: string;
  className?: string;
  onClick?: (e: React.MouseEvent<HTMLAnchorElement>) => void;
}

const PrimaryButton = ({
  label,
  href = "#",
  className = "",
  onClick,
}: PrimaryButtonProps) => {
  return (
    <Link
      href={href}
      onClick={onClick}
      className={`px-6 py-3 rounded-lg font-medium text-center transition-all bg-primary hover:bg-primary-dark hover:-translate-y-0.5 text-white inline-block ${className}`.trim()}
    >
      {label}
    </Link>
  );
};

export default PrimaryButton;
