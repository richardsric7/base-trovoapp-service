import React from 'react';
import { useSelector } from 'react-redux';
import { Outlet } from 'react-router-dom';
import SideBar from '../../components/sideBar';

function Dashboard() {
  const sidebarState = useSelector((state) => {
    console.log('state', state.sidebar.value);
    return state.sidebar.value;
  });

  return (
    <div className="flex bg-gray-100 h-screen items-start xl:space-x-5 xl:p-6 w-full">
      <div className="xl:block hidden h-full w-1/6">
        <SideBar />
      </div>
      <div className="flex w-4/6 w-full xl:rounded-3xl overflow-x-hidden h-full bg-white">
        {sidebarState && (
          <div className="xl:hidden fixed 2md:relative h-full w-3/4 sm:w-1/4">
            <SideBar mobileMode />
          </div>
        )}
        <div className="flex flex-col overflow-y-scroll w-full">
          <Outlet />
        </div>
      </div>
    </div>
  );
}

export default Dashboard;
