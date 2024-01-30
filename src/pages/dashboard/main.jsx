import React from 'react';
import { Outlet } from 'react-router-dom';
import SideBar from '../../components/sideBar';

function Dashboard() {
  return (
    <div className="flex bg-gray-100 h-screen items-start xl:space-x-5 xl:p-6 w-full">
      <div className="xl:block hidden h-full w-1/6">
        <SideBar />
      </div>
      <div className="flex flex-col overflow-y-scroll w-4/6 w-full xl:rounded-3xl h-full bg-white">
        <Outlet />
      </div>
    </div>
  );
}

export default Dashboard;
