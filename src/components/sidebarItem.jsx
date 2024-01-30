import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';
// eslint-disable-next-line
function SideBarItem({ label, icon, url, isActive, onSidebarClicked }) {
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

SideBarItem.propTypes = {
  label: PropTypes.string.isRequired,
  icon: PropTypes.string.isRequired,
  isActive: PropTypes.bool.isRequired,
  url: PropTypes.string.isRequired,
  onSidebarClicked: PropTypes.func.isRequired,
};

export default SideBarItem;
