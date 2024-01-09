import React from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import Welcome from '../pages/onboarding/welcome';
import Welcome1 from '../pages/onboarding/welcome1';
import Welcome2 from '../pages/onboarding/welcome2';
import Welcome3 from '../pages/onboarding/welcome3';
import ScrollToTop from '../components/scrollToTop';
import CreateAccount from '../pages/createAccount/main';
import RegistrationForm from '../pages/createAccount/registrationForm';
import CreatePassword from '../pages/createAccount/createPassword';
import AccountVerification from '../pages/createAccount/accountVerification';
import Backup from '../pages/createAccount/backup';
import ImportWallet from '../pages/importWallet/importWallet';
import RecoveryMain from '../pages/accountRecovery/main';
import Login from '../pages/login';

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
        <Route path="/register" element={<CreateAccount />}>
          <Route index element={<CreatePassword />} />
          <Route path="create-password" element={<CreatePassword />} />
          <Route path="form" element={<RegistrationForm />} />
        </Route>
        <Route
          path="/register/verification"
          element={<AccountVerification />}
        />
        <Route path="/register/backup" element={<Backup />} />
        <Route path="/import" element={<ImportWallet />} />
        <Route path="/recovery" element={<RecoveryMain />} />
        <Route path="/login" element={<Login />} />
        {/* <Route path="*" element={<Page404 />} /> */}
      </Routes>
    </BrowserRouter>
  );
}
