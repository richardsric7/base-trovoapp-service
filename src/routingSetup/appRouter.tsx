import { HashRouter, Route, Routes } from "react-router-dom";
import Welcome from "../pages/onboarding/welcome";
import Welcome1 from "../pages/onboarding/welcome1";
import Welcome2 from "../pages/onboarding/welcome2";
import Welcome3 from "../pages/onboarding/welcome3";
import ScrollToTop from "../components/scrollToTop";
import CreateAccount from "../pages/createAccount/main";
import RegistrationForm from "../pages/createAccount/registrationForm";
import CreatePassword from "../pages/createAccount/createPassword";
import AccountVerification from "../pages/createAccount/accountVerification";
import Backup from "../pages/createAccount/backup";
import ImportWallet from "../pages/importWallet/importWallet";
import RecoveryMain from "../pages/accountRecovery/main";
import Login from "../pages/login";
import Dashboard from "../pages/dashboard/main";
import Home from "../pages/dashboard/home";
import TokenizedAsset from "../pages/dashboard/tokenizedAsset";
import { ProtectedRoutes } from "./routeGuard";
import AnswerSecurityQuestions from "../pages/accountRecovery/answerSecurityQuestions";
import SetupSecurityQuestions from "../pages/accountRecovery/setupSecurityQuestions";
import RestoreInactiveAccount from "../pages/accountRecovery/restoreInactiveAccount";
import WalletView from "../pages/dashboard/walletView";
import SendAssetReceipt from "../pages/pdfPages/sendAssetReceipt";
import { History } from "../pages/dashboard/history";
import { Tokenize, TokenizeAssetForm } from "../pages/dashboard";

export default function AppRouter() {
  return (
    <HashRouter>
      <ScrollToTop />
      <Routes>
        <Route path="/" element={<Welcome1 />} />
        <Route path="/welcome" element={<Welcome />} />
        <Route path="/welcome1" element={<Welcome1 />} />
        <Route path="/welcome2" element={<Welcome2 />} />
        <Route path="/welcome3" element={<Welcome3 />} />
        <Route path="/login" element={<Login />} />
        <Route path="/import" element={<ImportWallet />} />
        <Route path="/recovery" element={<RecoveryMain />} />
        <Route path="/send-asset-receipt" element={<SendAssetReceipt />} />
        <Route
          path="/answer-security-questions"
          element={<AnswerSecurityQuestions />}
        />
        <Route
          path="/restore-inactive-account"
          element={<RestoreInactiveAccount />}
        />
        <Route path="/register" element={<CreateAccount />}>
          <Route index element={<CreatePassword />} />
          <Route path="create-password" element={<CreatePassword />} />
          <Route path="form" element={<RegistrationForm />} />
        </Route>
        <Route
          path="/register/verification"
          element={<AccountVerification />}
        />
        <Route element={<ProtectedRoutes />}>
          <Route path="/backup" element={<Backup />} />
          <Route path="/dashboard" element={<Dashboard />}>
            <Route index element={<Home />} />
            <Route path="tokenized-asset" element={<TokenizedAsset />} />
            <Route path="wallet" element={<WalletView />} />
            <Route path="history" element={<History />} />
            <Route path="tokenize" element={<Tokenize />}>
            </Route>
              <Route path="tokenize/apply" element={<TokenizeAssetForm />} />
            <Route
              path="setup-security-questions"
              element={<SetupSecurityQuestions />}
            />
          </Route>
        </Route>
        {/* <Route path="*" element={<Page404 />} /> */}
      </Routes>
    </HashRouter>
  );
}
