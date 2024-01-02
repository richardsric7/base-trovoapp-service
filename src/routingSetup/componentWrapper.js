import React from "react";
// import { useSelector } from "react-redux";
// import { Navigate, useLocation } from "react-router-dom";

export const ComponentWrapper = ({ children }) => {
//   const { user } = useSelector((state) => state.auth);
//   const location = useLocation();
//   if (!user && ["/buy-assets", "/sell-assets", '/'].includes(location.pathname)) {
//     return <Navigate to={"/login"} />;
//   }
  return <>{children}</>;
};