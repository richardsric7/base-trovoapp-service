import { BrowserRouter, Route, Routes } from "react-router-dom";
import { ScrollToTop } from "../components";
import LandingPage from "../pages/landingPage/index";

export const AppRouter = () => {
  return (
    <BrowserRouter>
      <ScrollToTop />
      <Routes>
        <Route
          path="/"
          element={
              <LandingPage />
          }
        />
        {/* <Route path="*" element={<Page404 />} /> */}
      </Routes>
    </BrowserRouter>
  );
};