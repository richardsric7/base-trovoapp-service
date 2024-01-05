import React from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import Welcome from '../pages/onboarding/welcome';
import Welcome1 from '../pages/onboarding/welcome1';
import Welcome2 from '../pages/onboarding/welcome2';
import Welcome3 from '../pages/onboarding/welcome3';
import ScrollToTop from '../components/scrollToTop';

export default function AppRouter() {
  return (
    <BrowserRouter>
      <ScrollToTop />
      <Routes>
        <Route path="/" element={<Welcome1 />} />
        <Route path="/welcome" element={<Welcome />} />
        <Route path="/welcome1" element={<Welcome1 />} />
        <Route path="/welcome2" element={<Welcome2 />} />
        <Route path="/welcome3" element={<Welcome3 />} />
        {/* <Route path="*" element={<Page404 />} /> */}
      </Routes>
    </BrowserRouter>
  );
}
