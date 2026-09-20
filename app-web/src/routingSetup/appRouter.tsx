import { HashRouter, Route, Routes } from 'react-router-dom';
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
import Dashboard from '../pages/dashboard/main';
import Home from '../pages/dashboard/home';
import { ProtectedRoutes } from './routeGuard';
import AnswerSecurityQuestions from '../pages/accountRecovery/answerSecurityQuestions';
import SetupSecurityQuestions from '../pages/accountRecovery/setupSecurityQuestions';
import RestoreInactiveAccount from '../pages/accountRecovery/restoreInactiveAccount';
import WalletView from '../pages/dashboard/wallet/walletView';
import SendAssetReceipt from '../pages/pdfPages/sendAssetReceipt';
import { History } from '../pages/dashboard/history';
import {
  Tokenize,
  TokenizationSetupAndCompliance,
  ConfirmTokenizationDetails,
} from '../pages/dashboard';
import AssetDetail from '../pages/dashboard/wallet/assetDetails';
import YieldView from '../pages/dashboard/wallet/yield';
import EarlyExitView from '../pages/dashboard/wallet/earlyExit';
import TokenizedAssetsListView from '../pages/dashboard/tokenizedAssets/assetsListView';
import TokenizedAssetDetailsView from '../pages/dashboard/tokenizedAssets/tokenizedAsset';
import AssetsList from '../pages/dashboard/addOrRemoveAssets/assetsList';
import { TokenizationApplication } from '../pages/dashboard/tokenize/tokenizationApplication';
import { TokenizationAssetInformation } from '../pages/dashboard/tokenize/tokenizationAssetInformation';
import { TokenizationAssetDocuments } from '../pages/dashboard/tokenize/tokenizationAssetDocuments';
import { TokenizationAssetTokenInformation } from '../pages/dashboard/tokenize/tokenizationAssetTokenInformation';
import TokenizationAssetDashboard from '../pages/dashboard/tokenize/tokenizationAssetDashboard';
import { TokenizationFeePayment } from '../pages/dashboard/tokenize/tokenizationFeePayment';
import {
  SharedAccessLanding,
  SharedAccessWalletInfo,
  SharedAccessUpdate,
  SharedAccessApprovals,
  SharedAccessApprovalDetails,
  SharedAccessAdd,
} from '../pages/dashboard/sharedAccess';

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
            <Route path="yield" element={<YieldView />} />
            <Route path="early-exit" element={<EarlyExitView />} />
            <Route path="asset-details" element={<AssetDetail />} />
            <Route
              path="tokenized-asset"
              element={<TokenizedAssetDetailsView />}
            />
            <Route
              path="tokenized-asset-list"
              element={<TokenizedAssetsListView />}
            />
            <Route path="add-remove-assets" element={<AssetsList />} />
            <Route path="wallet" element={<WalletView />} />
            <Route path="history" element={<History />} />
            <Route path="shared-access" element={<SharedAccessLanding />} />
            <Route
              path="shared-access/wallets"
              element={<SharedAccessLanding />}
            />
            <Route
              path="shared-access/wallets/:address"
              element={<SharedAccessWalletInfo />}
            />
            <Route
              path="shared-access/update"
              element={<SharedAccessUpdate />}
            />
            <Route
              path="shared-access/approvals"
              element={<SharedAccessApprovals />}
            />
            <Route
              path="shared-access/approvals/:id"
              element={<SharedAccessApprovalDetails />}
            />
            <Route path="shared-access/add" element={<SharedAccessAdd />} />
            <Route path="tokenize" element={<Tokenize />}></Route>
            <Route
              path="tokenize/confirm-details"
              element={<ConfirmTokenizationDetails />}
            />
            <Route
              path="tokenize/confirm-details/:id"
              element={<ConfirmTokenizationDetails />}
            />
            <Route
              path="tokenize/asset-dashboard/:id"
              element={<TokenizationAssetDashboard />}
            />
            <Route
              path="tokenize/fee-payment/:id"
              element={<TokenizationFeePayment />}
            />
            <Route path="tokenize/apply" element={<TokenizationApplication />}>
              <Route index element={<TokenizationSetupAndCompliance />} />
              <Route path=":id" element={<TokenizationSetupAndCompliance />} />
              <Route
                path=":id/asset-information"
                element={<TokenizationAssetInformation />}
              />
              <Route
                path=":id/asset-documents"
                element={<TokenizationAssetDocuments />}
              />
              <Route
                path=":id/asset-token-information"
                element={<TokenizationAssetTokenInformation />}
              />
            </Route>
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
