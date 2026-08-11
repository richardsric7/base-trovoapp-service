import { LandingPage } from "./trovotech-main";
import {
  // BrowserRouter,
  HashRouter,
  Navigate,
  Route,
  Routes,
} from "react-router-dom";
import { PrivacyPolicy } from "./trovotech-main/privacy-policy";
import { TermsOfUse } from "./trovotech-main/TermsOfUse";

import AccountDeletion from "./trovotech-main/help-center/account-deletion";
import ScrollToTop from "./trovotech-main/SccrollToTop";

function App() {
  return (
    <HashRouter>
      <ScrollToTop />
      <Routes>
        <Route path={"/"} element={<LandingPage />} />
        <Route path="/features" element={<LandingPage scrollTo="features" />} />
        <Route path="/download" element={<LandingPage scrollTo="download" />} />
        <Route path="/contact" element={<LandingPage scrollTo="contact" />} />

        <Route path={"/terms-of-use"} element={<TermsOfUse />} />
        <Route path={"/privacy-policy"} element={<PrivacyPolicy />} />
        <Route
          path={"/help-center/account-deletion"}
          element={<AccountDeletion />}
        />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </HashRouter>
  );
}

export default App;
