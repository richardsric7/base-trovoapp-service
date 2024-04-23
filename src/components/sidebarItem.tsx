import React from 'react';
import { Link } from 'react-router-dom';

type Props = {
  label: string;
  icon: string;
  isActive: boolean;
  url: string;
  onSidebarClicked: () => void;
};

function SideBarItem({ label, icon, url, isActive, onSidebarClicked }: Props) {
  const classes = `flex items-center rounded-md px-3 py-2 space-x-3 text-sm hover:bg-primary-700 ${
    isActive ? 'bg-primary-700 text-white' : 'text-gray-300'
  }`;
  return (
    <Link to={url} className={classes} onClick={() => onSidebarClicked()}>
      <img src={icon} alt="home" className="h-5" />
      <span>{label}</span>
    </Link>
  );
}

export default SideBarItem;
