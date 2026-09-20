import { ReactNode, useEffect, useMemo, useRef, useState } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import Button from '../../../components/button';
import ButtonSecondary from '../../../components/buttonSecondary';
import Modal from '../../../components/modal';
import Tabs from '../../../components/tabs';
import { CustomTable } from '../../../components';
import Dropdown from '../../../components/dropdown';
import { RootState } from '../../../store/reduxStore';
import { setUser } from '../../../store/authSlice';
import { useLazyGetUserQuery } from '../../../store/api/authApi';
import {
  useLazyGetApprovalsQuery,
  useDisableSharedAccessMutation,
  useAddSharedAccessMutation,
  useUpdateSharedAccessMutation,
  useLazyCheckUsernameQuery,
  useApproveSharedAccessMutation,
  useRejectSharedAccessMutation,
} from '../../../store/api/sharedAccessApis';
import { Permission, PermissionState } from '../../../types/permission';
import { Wallet } from '../../../types/wallet';
import { deserializeUserData } from '../../../utils/deserializeAndStoreUserData';
import { Encryptor } from '../../../utils/encryptor';
import { signBase64Txn } from '../../../utils/trovoSDK';
import { showNotification, toggleLoader } from '../../../utils/showToaster';
import styles from './landing.module.css';
// import { getExplorerBaseUrl } from '../../../utils/utilities';

type SharedAccessRow = {
  wallet: string;
  owner: string;
  permissions: string[];
  description: string;
  address: string;
};

type SharedAccessActivityRow = {
  wallet: string;
  transactionType: string;
  initiatedBy: string;
  date: string;
  transactionStatus: 'Pending' | 'Completed' | 'Rejected';
  approvalStatus: string;
};

type SharedAccessApprovalRecord = {
  alias?: string;
  walletAlias?: string;
  transactionType?: string;
  initiator?: string;
  createdAt?: string;
  transactionStatus?: string;
  approvalsGotten?: number;
  approvalsNeeded?: number;
  approvalStatus?: string;
};

const APPROVALS_PER_PAGE = 20;

const accessMode = [
  { text: 'Access granted to me', value: 'Access granted to me' },
  { text: 'Access granted by me', value: 'Access granted by me' },
];

const filter = [
  { text: 'All', value: 'All' },
  { text: 'Viewer', value: 'Viewer' },
  { text: 'Initiator', value: 'Initiator' },
  { text: 'Approver', value: 'Approver' },
];

const permissionToLabelMap: Record<string, string> = {
  'VIEW-ONLY': 'Viewer',
  VIEWER: 'Viewer',
  INITIATOR: 'Initiator',
  APPROVER: 'Approver',
};

const filterToPermissionCodeMap: Record<string, string | null> = {
  All: null,
  Viewer: 'VIEW-ONLY',
  Initiator: 'INITIATOR',
  Approver: 'APPROVER',
};

const formatPermissionLabel = (permission: string) => {
  const normalizedPermission = permission.toUpperCase();
  return permissionToLabelMap[normalizedPermission] ?? permission;
};

const toTitleCase = (value: string) =>
  value
    .toLowerCase()
    .split(' ')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ');

const normalizeTransactionType = (value?: string) => {
  if (!value) return '-';
  return toTitleCase(value.replaceAll('_', ' '));
};

const normalizeTransactionStatus = (
  value?: string,
): SharedAccessActivityRow['transactionStatus'] => {
  const normalizedValue = (value ?? '').toUpperCase();
  if (normalizedValue === 'COMPLETED') return 'Completed';
  if (normalizedValue === 'REJECTED') return 'Rejected';
  return 'Pending';
};

const formatApprovalDate = (value?: string) => {
  if (!value) return '-';
  const parsedDate = new Date(value);
  if (Number.isNaN(parsedDate.getTime())) return '-';

  return new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: true,
  }).format(parsedDate);
};

const formatDateRangeLabel = (startDate: string, endDate: string) => {
  const fmt = (value: string) => {
    const parsed = new Date(value);
    if (Number.isNaN(parsed.getTime())) return value;
    return new Intl.DateTimeFormat('en-GB', {
      day: '2-digit',
      month: '2-digit',
      year: '2-digit',
    }).format(parsed);
  };
  return `${fmt(startDate)} – ${fmt(endDate)}`;
};

const transactionStatusOptions = [
  { text: 'All', value: 'All' },
  { text: 'Pending', value: 'Pending' },
  { text: 'Completed', value: 'Completed' },
  { text: 'Rejected', value: 'Rejected' },
];

const transactionTypeBaseOptions = [
  { text: 'All', value: 'All' },
  { text: 'Swap', value: 'SWAP' },
  { text: 'Payment', value: 'PAYMENT' },
  { text: 'Modify Shared Access', value: 'MODIFY SHARED ACCESS' },
  { text: 'Disable Shared Access', value: 'DISABLE SHARED ACCESS' },
  { text: 'Opt In Asset', value: 'OPT IN ASSET' },
  { text: 'Opt Out Asset', value: 'OPT OUT ASSET' },
  { text: 'Accept Pending Asset', value: 'ACCEPT PENDING ASSET' },
  { text: 'Reject Pending Asset', value: 'REJECT PENDING ASSET' },
  { text: 'Make Market Offer', value: 'MAKE MARKET OFFER' },
  { text: 'Delete Market Offer', value: 'DELETE MARKET OFFER' },
  { text: 'Modify Market Offer', value: 'MODIFY MARKET OFFER' },
  { text: 'Mint Token', value: 'MINT TOKEN' },
  { text: 'Burn Token', value: 'BURN TOKEN' },
  { text: 'Bulk Payment', value: 'BULK PAYMENT' },
];

const getDateRangeQueryValue = (dateRange: string) => {
  if (dateRange === 'All') return null;

  if (dateRange === 'Custom date') return null;

  const now = new Date();
  const endDate = new Date(now.getFullYear(), now.getMonth(), now.getDate());
  const startDate = new Date(endDate);

  if (dateRange === 'Today') {
    // startDate already points to today.
  } else if (dateRange === 'Last 7 days') {
    startDate.setDate(endDate.getDate() - 6);
  } else if (dateRange === 'Last 30 days') {
    startDate.setDate(endDate.getDate() - 29);
  } else {
    return null;
  }

  const formatAsYmd = (value: Date) => {
    const year = value.getFullYear();
    const month = String(value.getMonth() + 1).padStart(2, '0');
    const day = String(value.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  };

  return `${formatAsYmd(startDate)}|${formatAsYmd(endDate)}`;
};

export default function SharedAccessLanding() {
  const appUser = useSelector((state: RootState) => state.auth.user);
  // const appState = useSelector((state: RootState) => state.appState!);
  const [fetchApprovals] = useLazyGetApprovalsQuery();
  const [selectedFilter, setSelectedFilter] = useState<string>(filter[0].value);
  const [selectedAccessMode, setSelectedAccessMode] = useState<string>(
    accessMode[0].value,
  );
  const [selectedTransactionType, setSelectedTransactionType] =
    useState<string>('All');
  const [selectedDateRange, setSelectedDateRange] = useState<string>('All');
  const [selectedTransactionStatus, setSelectedTransactionStatus] =
    useState<string>('Pending');
  const [selectedInitiator, setSelectedInitiator] = useState<string>('');
  const [selectedWalletAlias, setSelectedWalletAlias] = useState<string>('');
  const [selectedWalletAddress, setSelectedWalletAddress] =
    useState<string>('');
  const [selectedDescription, setSelectedDescription] = useState<string>('');
  const [draftInitiator, setDraftInitiator] = useState<string>('');
  const [draftWalletAlias, setDraftWalletAlias] = useState<string>('');
  const [draftWalletAddress, setDraftWalletAddress] = useState<string>('');
  const [draftDescription, setDraftDescription] = useState<string>('');
  const [customStartDate, setCustomStartDate] = useState<string>('');
  const [customEndDate, setCustomEndDate] = useState<string>('');
  const [draftCustomStartDate, setDraftCustomStartDate] = useState<string>('');
  const [draftCustomEndDate, setDraftCustomEndDate] = useState<string>('');
  const [isApprovalsLoading, setIsApprovalsLoading] = useState(false);
  const [approvalRecords, setApprovalRecords] = useState<
    SharedAccessApprovalRecord[]
  >([]);
  const [approvalPage, setApprovalPage] = useState(1);
  const [, setApprovalTotalRecords] = useState(0);
  const pendingFallbackAppliedRef = useRef(false);
  const navigate = useNavigate();
  const [showSharedAccessView, setShowSharedAccessView] = useState(true);
  const [activeModal, setActiveModal] = useState<
    | 'grant'
    | 'confirm'
    | 'success'
    | 'details'
    | 'modify'
    | 'update'
    | 'initiators'
    | 'confirmModify'
    | 'modifySuccess'
    | 'dateRange'
    | 'moreFilters'
    | 'disableWarn'
    | 'disableConfirm'
    | 'disableSuccess'
    | 'approvalDetails'
    | null
  >(null);
  const [selectedRow, setSelectedRow] = useState<SharedAccessRow | null>(null);
  const [selectedActivityRecord, setSelectedActivityRecord] =
    useState<SharedAccessApprovalRecord | null>(null);
  const [approvalPassword, setApprovalPassword] = useState('');
  const [approvalPasswordErr, setApprovalPasswordErr] = useState('');
  const [approvalReason, setApprovalReason] = useState('');
  const [isApproving, setIsApproving] = useState(false);
  const [isRejecting, setIsRejecting] = useState(false);
  const [showApprovalPassword, setShowApprovalPassword] = useState(false);
  const [showRejectReason, setShowRejectReason] = useState(false);
  const [includeSignedTransactions, setIncludeSignedTransactions] =
    useState(false);

  const [disableSharedAccess] = useDisableSharedAccessMutation();
  const [approveSharedAccess] = useApproveSharedAccessMutation();
  const [rejectSharedAccess] = useRejectSharedAccessMutation();
  const [getUser] = useLazyGetUserQuery();
  const dispatch = useDispatch();

  const [disableTxData, setDisableTxData] = useState<any>(null);
  const [disablePassword, setDisablePassword] = useState('');
  const [disablePasswordErr, setDisablePasswordErr] = useState('');
  const [isDisabling, setIsDisabling] = useState(false);
  const [disableSuccessMessage, setDisableSuccessMessage] = useState('');

  const [addSharedAccess, { isLoading: isAdding }] =
    useAddSharedAccessMutation();
  const [updateSharedAccess, { isLoading: isUpdating }] =
    useUpdateSharedAccessMutation();
  const [checkUsernameQuery] = useLazyCheckUsernameQuery();

  // Selector for shareable wallets
  const shareableWallets = useMemo(() => {
    if (!appUser?.userWallets) return [];
    return appUser.userWallets.filter(
      (wallet) =>
        (wallet.walletType === 0 ||
          wallet.walletType === 1 ||
          wallet.walletType === 2) &&
        !wallet.isSharedWallet,
    );
  }, [appUser]);

  // Grant Access States
  const [grantWallet, setGrantWallet] = useState<Wallet | null>(null);
  const [grantViewers, setGrantViewers] = useState<string[]>([]);
  const [grantApprovers, setGrantApprovers] = useState<string[]>([]);
  const [grantInitiators, setGrantInitiators] = useState<string[]>([]);
  const [grantUserFullnames, setGrantUserFullnames] = useState<
    Record<string, string>
  >({});
  const [grantAddApprovers, setGrantAddApprovers] = useState(false);
  const [grantNoOfApprovers, setGrantNoOfApprovers] = useState(3);
  const [grantNoOfApprovalsNeeded, setGrantNoOfApprovalsNeeded] = useState(2);
  const [grantStep, setGrantStep] = useState(0); // 0 = Viewers, 1 = Approvers, 2 = Initiators
  const [grantViewerInput, setGrantViewerInput] = useState('');
  const [grantApproverInput, setGrantApproverInput] = useState('');
  const [grantInitiatorInput, setGrantInitiatorInput] = useState('');
  const [grantPassword, setGrantPassword] = useState('');
  const [grantPasswordErr, setGrantPasswordErr] = useState('');
  const [showGrantPassword, setShowGrantPassword] = useState(false);
  const [grantViewerErr, setGrantViewerErr] = useState('');
  const [grantApproverErr, setGrantApproverErr] = useState('');
  const [grantInitiatorErr, setGrantInitiatorErr] = useState('');

  // Modify/Update Shared Access States
  const [modifyViewers, setModifyViewers] = useState<Permission[]>([]);
  const [modifyApprovers, setModifyApprovers] = useState<Permission[]>([]);
  const [modifyInitiators, setModifyInitiators] = useState<Permission[]>([]);
  const [modifyNoOfApprovalsNeeded, setModifyNoOfApprovalsNeeded] = useState(2);
  const [modifyNoOfApprovers, setModifyNoOfApprovers] = useState(3);
  const [modifyAddApprovers, setModifyAddApprovers] = useState(false);
  const [modifyViewerInput, setModifyViewerInput] = useState('');
  const [modifyApproverInput, setModifyApproverInput] = useState('');
  const [modifyInitiatorInput, setModifyInitiatorInput] = useState('');
  const [modifyPassword, setModifyPassword] = useState('');
  const [modifyPasswordErr, setModifyPasswordErr] = useState('');
  const [showModifyPassword, setShowModifyPassword] = useState(false);
  const [modifyViewerErr, setModifyViewerErr] = useState('');
  const [modifyApproverErr, setModifyApproverErr] = useState('');
  const [modifyInitiatorErr, setModifyInitiatorErr] = useState('');

  // Modal State Initialization Effects
  useEffect(() => {
    if (activeModal === 'grant') {
      if (shareableWallets.length > 0) {
        setGrantWallet(shareableWallets[0]);
      } else {
        setGrantWallet(null);
      }
      setGrantViewers([]);
      if (appUser) {
        setGrantApprovers([appUser.username]);
        setGrantInitiators([appUser.username]);
        setGrantUserFullnames({
          [appUser.username]: `${appUser.firstName} ${appUser.lastName}`,
        });
      } else {
        setGrantApprovers([]);
        setGrantInitiators([]);
        setGrantUserFullnames({});
      }
      setGrantAddApprovers(false);
      setGrantNoOfApprovers(3);
      setGrantNoOfApprovalsNeeded(2);
      setGrantStep(0);
      setGrantViewerInput('');
      setGrantApproverInput('');
      setGrantInitiatorInput('');
      setGrantPassword('');
      setGrantPasswordErr('');
      setShowGrantPassword(false);
      setGrantViewerErr('');
      setGrantApproverErr('');
      setGrantInitiatorErr('');
    }
  }, [activeModal, shareableWallets, appUser]);

  useEffect(() => {
    if (activeModal === 'modify' && selectedRow && appUser) {
      const wallet = appUser.userWallets.find(
        (w) => w.address === selectedRow.address,
      );
      if (wallet) {
        const rawPermissions = wallet.permissions ?? [];
        const viewersList: Permission[] = rawPermissions
          .filter(
            (p) => p.permission === 'VIEW-ONLY' || p.permission === 'VIEWER',
          )
          .map((p) => ({ ...p, permissionState: null }));
        const approversList: Permission[] = rawPermissions
          .filter((p) => p.permission === 'APPROVER')
          .map((p) => ({ ...p, permissionState: null }));
        const initiatorsList: Permission[] = rawPermissions
          .filter((p) => p.permission === 'INITIATOR')
          .map((p) => ({ ...p, permissionState: null }));

        setModifyViewers(viewersList);
        setModifyApprovers(approversList);
        setModifyInitiators(initiatorsList);

        const needed = wallet.numberOfApprovalsNeeded ?? 2;
        setModifyNoOfApprovalsNeeded(needed);
        setModifyNoOfApprovers(approversList.length || 3);
        setModifyAddApprovers(approversList.length > 0 && needed > 0);
      }
      setModifyViewerInput('');
      setModifyApproverInput('');
      setModifyInitiatorInput('');
      setModifyPassword('');
      setModifyPasswordErr('');
      setShowModifyPassword(false);
      setModifyViewerErr('');
      setModifyApproverErr('');
      setModifyInitiatorErr('');
    }
  }, [activeModal, selectedRow, appUser]);

  useEffect(() => {
    toggleLoader(isApprovalsLoading || isDisabling || isAdding || isUpdating);

    return () => {
      toggleLoader(false);
    };
  }, [isApprovalsLoading, isDisabling, isAdding, isUpdating]);

  const checkUsername = async (username: string, activeWallet: Wallet) => {
    try {
      const encryptor = new Encryptor();
      const decryptedSecretKey = await encryptor.getSecretKey(appUser!);
      const payload = {
        signer: activeWallet.signer || appUser!.primarySigner,
        address: activeWallet.address,
        secretKey: decryptedSecretKey,
        body: { username },
      };

      setIsApprovalsLoading(true);
      const response = await checkUsernameQuery(payload);
      setIsApprovalsLoading(false);

      if ('data' in response) {
        return (response.data as any).userData;
      }
      return null;
    } catch (e: any) {
      setIsApprovalsLoading(false);
      showNotification('error', e.message || 'Error checking username');
      return null;
    }
  };

  // Grant Access Helpers
  const handleAddGrantViewer = async () => {
    setGrantViewerErr('');
    const username = grantViewerInput.trim();
    if (!username) {
      setGrantViewerErr('Enter username');
      return;
    }

    if (appUser?.username === username) {
      showNotification(
        'error',
        'You cannot add yourself as a viewer on this wallet because as the owner of this wallet you already have view access',
      );
      return;
    }

    if (grantViewers.includes(username)) {
      setGrantViewerErr('Username already added');
      return;
    }

    if (grantApprovers.includes(username)) {
      showNotification(
        'error',
        'This user is already an Approver on this wallet. An Approver cannot be added as a Viewer.',
      );
      return;
    }

    if (grantInitiators.includes(username)) {
      showNotification(
        'error',
        'This user is already an Initiator on this wallet. An Initiator cannot be added as a Viewer.',
      );
      return;
    }

    const userInfo = await checkUsername(username, grantWallet!);
    if (!userInfo) {
      setGrantViewerErr(
        'This username does not exist. Please enter a valid username.',
      );
      return;
    }

    setGrantUserFullnames((prev) => ({
      ...prev,
      [username]: `${userInfo.firstName} ${userInfo.lastName}`,
    }));
    setGrantViewers((prev) => [...prev, username]);
    setGrantViewerInput('');
  };

  const handleAddGrantApprover = async () => {
    setGrantApproverErr('');
    const username = grantApproverInput.trim();
    if (!username) {
      setGrantApproverErr('Enter username');
      return;
    }

    if (grantApprovers.includes(username)) {
      setGrantApproverErr('Username already added');
      return;
    }

    if (grantApprovers.length === grantNoOfApprovers) {
      showNotification(
        'error',
        'The number of approvers cannot exceed the specified amount.',
      );
      return;
    }

    if (grantViewers.includes(username)) {
      const confirmed = window.confirm(
        'This user is currently a Viewer. Adding them as an Approver will remove them from the Viewer list. Proceed?',
      );
      if (confirmed) {
        setGrantViewers((prev) => prev.filter((u) => u !== username));
        setGrantApprovers((prev) => [...prev, username]);
        setGrantInitiators((prev) => [...prev, username]);
        setGrantApproverInput('');
      }
      return;
    }

    if (grantInitiators.includes(username)) {
      setGrantApprovers((prev) => [...prev, username]);
      setGrantApproverInput('');
      return;
    }

    const userInfo = await checkUsername(username, grantWallet!);
    if (!userInfo) {
      setGrantApproverErr(
        'This username does not exist. Please enter a valid username.',
      );
      return;
    }

    setGrantUserFullnames((prev) => ({
      ...prev,
      [username]: `${userInfo.firstName} ${userInfo.lastName}`,
    }));
    setGrantApprovers((prev) => [...prev, username]);
    setGrantInitiators((prev) => [...prev, username]);
    setGrantApproverInput('');
  };

  const handleAddGrantInitiator = async () => {
    setGrantInitiatorErr('');
    const username = grantInitiatorInput.trim();
    if (!username) {
      setGrantInitiatorErr('Enter username');
      return;
    }

    if (grantInitiators.includes(username)) {
      setGrantInitiatorErr('Username already added');
      return;
    }

    if (grantViewers.includes(username)) {
      const confirmed = window.confirm(
        'This user is currently a Viewer. Adding them as an Initiator will remove them from the Viewer list. Proceed?',
      );
      if (confirmed) {
        setGrantViewers((prev) => prev.filter((u) => u !== username));
        setGrantInitiators((prev) => [...prev, username]);
        setGrantInitiatorInput('');
      }
      return;
    }

    if (grantApprovers.includes(username)) {
      setGrantInitiators((prev) => [...prev, username]);
      setGrantInitiatorInput('');
      return;
    }

    const userInfo = await checkUsername(username, grantWallet!);
    if (!userInfo) {
      setGrantInitiatorErr(
        'This username does not exist. Please enter a valid username.',
      );
      return;
    }

    setGrantUserFullnames((prev) => ({
      ...prev,
      [username]: `${userInfo.firstName} ${userInfo.lastName}`,
    }));
    setGrantInitiators((prev) => [...prev, username]);
    setGrantInitiatorInput('');
  };

  // Modify Shared Access Helpers
  const handleAddModifyViewer = async () => {
    console.log('got here 1');
    setModifyViewerErr('');
    const username = modifyViewerInput.trim();
    if (!username) {
      setModifyViewerErr('Enter username');
      return;
    }

    if (appUser?.username === username) {
      showNotification(
        'error',
        'You cannot add yourself as a viewer on this wallet because as the owner of this wallet you already have view access',
      );
      return;
    }

    if (modifyViewers.some((v) => v.targetUsername === username)) {
      setModifyViewerErr('Username already added');
      return;
    }

    if (
      modifyApprovers.some(
        (a) =>
          a.targetUsername === username &&
          a.permissionState !== PermissionState.Revoked,
      )
    ) {
      showNotification(
        'error',
        'This user is already an Approver on this wallet. An Approver cannot be added as a Viewer.',
      );
      return;
    }

    if (
      modifyInitiators.some(
        (i) =>
          i.targetUsername === username &&
          i.permissionState !== PermissionState.Revoked,
      )
    ) {
      showNotification(
        'error',
        'This user is already an Initiator on this wallet. An Initiator cannot be added as a Viewer.',
      );
      return;
    }

    const wallet = appUser?.userWallets.find(
      (w) => w.address === selectedRow?.address,
    );
    if (!wallet) return;

    console.log('got here 2');
    const userInfo = await checkUsername(username, wallet);
    console.log('got here 3', userInfo);
    if (!userInfo) {
      setModifyViewerErr(
        'This username does not exist. Please enter a valid username.',
      );
      return;
    }

    const newPerm: Permission = {
      walletAddress: wallet.address,
      targetUsername: username,
      fullName: `${userInfo.firstName} ${userInfo.lastName}`,
      permission: 'VIEW-ONLY',
      permissionState: PermissionState.Added,
    };
    setModifyViewers((prev) => [...prev, newPerm]);
    setModifyViewerInput('');
  };

  const handleAddModifyApprover = async () => {
    setModifyApproverErr('');
    const username = modifyApproverInput.trim();
    if (!username) {
      setModifyApproverErr('Enter username');
      return;
    }

    if (modifyApprovers.some((a) => a.targetUsername === username)) {
      setModifyApproverErr('Username already added');
      return;
    }

    const activeApproversCount = modifyApprovers.filter(
      (a) => a.permissionState !== PermissionState.Revoked,
    ).length;
    if (activeApproversCount === modifyNoOfApprovers) {
      showNotification(
        'error',
        'The number of approvers cannot exceed the specified amount.',
      );
      return;
    }

    const wallet = appUser?.userWallets.find(
      (w) => w.address === selectedRow?.address,
    );
    if (!wallet) return;

    const existingViewerIndex = modifyViewers.findIndex(
      (v) =>
        v.targetUsername === username &&
        v.permissionState !== PermissionState.Revoked,
    );
    if (existingViewerIndex !== -1) {
      const viewer = modifyViewers[existingViewerIndex];
      const confirmed = window.confirm(
        'This user is currently a Viewer. Adding them as an Approver will revoke their Viewer access. Proceed?',
      );
      if (confirmed) {
        setModifyApprovers((prev) => [
          ...prev,
          {
            walletAddress: wallet.address,
            targetUsername: username,
            fullName: viewer.fullName,
            permission: 'APPROVER',
            permissionState: PermissionState.Added,
          },
        ]);
        setModifyInitiators((prev) => [
          ...prev,
          {
            walletAddress: wallet.address,
            targetUsername: username,
            fullName: viewer.fullName,
            permission: 'INITIATOR',
            permissionState: PermissionState.Added,
          },
        ]);
        setModifyApproverInput('');

        if (viewer.permissionState === PermissionState.Added) {
          setModifyViewers((prev) =>
            prev.filter((v) => v.targetUsername !== username),
          );
        } else {
          setModifyViewers((prev) =>
            prev.map((v) =>
              v.targetUsername === username
                ? { ...v, permissionState: PermissionState.Revoked }
                : v,
            ),
          );
        }
      }
      return;
    }

    const existingInitiator = modifyInitiators.find(
      (i) => i.targetUsername === username,
    );
    if (existingInitiator) {
      setModifyApprovers((prev) => [
        ...prev,
        {
          walletAddress: wallet.address,
          targetUsername: username,
          fullName: existingInitiator.fullName,
          permission: 'APPROVER',
          permissionState: PermissionState.Added,
        },
      ]);
      setModifyApproverInput('');
      return;
    }

    const userInfo = await checkUsername(username, wallet);
    if (!userInfo) {
      setModifyApproverErr(
        'This username does not exist. Please enter a valid username.',
      );
      return;
    }

    setModifyApprovers((prev) => [
      ...prev,
      {
        walletAddress: wallet.address,
        targetUsername: username,
        fullName: `${userInfo.firstName} ${userInfo.lastName}`,
        permission: 'APPROVER',
        permissionState: PermissionState.Added,
      },
    ]);
    setModifyInitiators((prev) => [
      ...prev,
      {
        walletAddress: wallet.address,
        targetUsername: username,
        fullName: `${userInfo.firstName} ${userInfo.lastName}`,
        permission: 'INITIATOR',
        permissionState: PermissionState.Added,
      },
    ]);
    setModifyApproverInput('');
  };

  const handleAddModifyInitiator = async () => {
    setModifyInitiatorErr('');
    const username = modifyInitiatorInput.trim();
    if (!username) {
      setModifyInitiatorErr('Enter username');
      return;
    }

    if (modifyInitiators.some((i) => i.targetUsername === username)) {
      setModifyInitiatorErr('Username already added');
      return;
    }

    const wallet = appUser?.userWallets.find(
      (w) => w.address === selectedRow?.address,
    );
    if (!wallet) return;

    const existingViewerIndex = modifyViewers.findIndex(
      (v) =>
        v.targetUsername === username &&
        v.permissionState !== PermissionState.Revoked,
    );
    if (existingViewerIndex !== -1) {
      const viewer = modifyViewers[existingViewerIndex];
      const confirmed = window.confirm(
        'This user is currently a Viewer. Adding them as an Initiator will revoke their Viewer access. Proceed?',
      );
      if (confirmed) {
        setModifyInitiators((prev) => [
          ...prev,
          {
            walletAddress: wallet.address,
            targetUsername: username,
            fullName: viewer.fullName,
            permission: 'INITIATOR',
            permissionState: PermissionState.Added,
          },
        ]);
        setModifyInitiatorInput('');

        if (viewer.permissionState === PermissionState.Added) {
          setModifyViewers((prev) =>
            prev.filter((v) => v.targetUsername !== username),
          );
        } else {
          setModifyViewers((prev) =>
            prev.map((v) =>
              v.targetUsername === username
                ? { ...v, permissionState: PermissionState.Revoked }
                : v,
            ),
          );
        }
      }
      return;
    }

    const existingApprover = modifyApprovers.find(
      (a) => a.targetUsername === username,
    );
    if (existingApprover) {
      setModifyInitiators((prev) => [
        ...prev,
        {
          walletAddress: wallet.address,
          targetUsername: username,
          fullName: existingApprover.fullName,
          permission: 'INITIATOR',
          permissionState: PermissionState.Added,
        },
      ]);
      setModifyInitiatorInput('');
      return;
    }

    const userInfo = await checkUsername(username, wallet);
    if (!userInfo) {
      setModifyInitiatorErr(
        'This username does not exist. Please enter a valid username.',
      );
      return;
    }

    setModifyInitiators((prev) => [
      ...prev,
      {
        walletAddress: wallet.address,
        targetUsername: username,
        fullName: `${userInfo.firstName} ${userInfo.lastName}`,
        permission: 'INITIATOR',
        permissionState: PermissionState.Added,
      },
    ]);
    setModifyInitiatorInput('');
  };

  const handlePermissionItemClick = (
    permList: Permission[],
    setPermList: React.Dispatch<React.SetStateAction<Permission[]>>,
    index: number,
    relation: string,
  ) => {
    const perm = permList[index];
    if (perm.permissionState === PermissionState.Revoked) {
      setPermList((prev) =>
        prev.map((p, i) => (i === index ? { ...p, permissionState: null } : p)),
      );
    } else if (perm.permissionState === PermissionState.Added) {
      setPermList((prev) => prev.filter((_, i) => i !== index));
    } else {
      const confirmed = window.confirm(
        `You are about to revoke ${relation.toLowerCase()} access for ${perm.targetUsername}. Proceed?`,
      );
      if (confirmed) {
        setPermList((prev) =>
          prev.map((p, i) =>
            i === index
              ? { ...p, permissionState: PermissionState.Revoked }
              : p,
          ),
        );
      }
    }
  };

  const handleModifyAddApproversChange = (checked: boolean) => {
    if (!checked && modifyApprovers.length > 0) {
      const confirmed = window.confirm(
        'All existing approvers and initiators will be revoked/removed. Proceed?',
      );
      if (confirmed) {
        setModifyApprovers((prev) =>
          prev.map((p) => ({ ...p, permissionState: PermissionState.Revoked })),
        );
        setModifyInitiators((prev) =>
          prev.map((p) => ({ ...p, permissionState: PermissionState.Revoked })),
        );
        setModifyAddApprovers(false);
      }
    } else {
      setModifyAddApprovers(checked);
      if (checked) {
        setModifyApprovers((prev) =>
          prev.map((p) =>
            p.permissionState === PermissionState.Revoked
              ? { ...p, permissionState: null }
              : p,
          ),
        );
        setModifyInitiators((prev) =>
          prev.map((p) =>
            p.permissionState === PermissionState.Revoked
              ? { ...p, permissionState: null }
              : p,
          ),
        );
        setActiveModal('update');
      }
    }
  };

  const isModifyValid = () => {
    const activeApprovers = modifyApprovers.filter(
      (p) => p.permissionState !== PermissionState.Revoked,
    );
    const activeInitiators = modifyInitiators.filter(
      (p) => p.permissionState !== PermissionState.Revoked,
    );

    if (modifyAddApprovers) {
      if (activeApprovers.length === 0) {
        showNotification('error', 'Please add approvers');
        setActiveModal('update');
        return false;
      }

      if (activeApprovers.length !== modifyNoOfApprovers) {
        showNotification(
          'error',
          'Please add the required number of approvers',
        );
        setActiveModal('update');
        return false;
      }

      if (activeInitiators.length === 0) {
        showNotification('error', 'Please add at least one initiator');
        setActiveModal('initiators');
        return false;
      }
    }
    return true;
  };

  const refreshUserInfo = async (
    decryptedSecretKey: string,
    passwordUsed: string,
  ) => {
    try {
      const encryptor = new Encryptor();
      const importedPayload = {
        signer: appUser!.address,
        address: appUser!.address,
        secretKey: appUser!.secretKeys[0],
        body: { userId: appUser!.username, import: 1 },
      };
      const { data: importData } = await getUser(importedPayload);
      if (importData) {
        const userData = deserializeUserData(importData);
        const base64EncryptedSecretKey = await encryptor.encryptData(
          decryptedSecretKey,
          passwordUsed,
          appUser!.address,
        );
        const passwordHash = await encryptor.createHash(userData.username);
        const encryptedPassword = await encryptor.encryptData(
          passwordUsed,
          passwordHash,
          appUser!.address,
        );
        const updatedUser = {
          ...userData,
          password: encryptedPassword,
          isLoggedIn: true,
          currency: 'USD',
          secretKeys: [
            ...appUser!.secretKeys.filter((k) => k !== appUser!.secretKeys[0]),
            base64EncryptedSecretKey,
          ],
        };
        const userHash = await encryptor.createHash(updatedUser.username);
        const base64EncryptedUserData = await encryptor.encryptData(
          JSON.stringify(updatedUser),
          userHash,
          updatedUser.address,
        );
        dispatch(
          setUser({
            key: userHash,
            user: updatedUser,
            encryptedUser: base64EncryptedUserData,
          }),
        );
      }
    } catch (err) {
      console.error('Failed to refresh user info', err);
    }
  };

  const handleGrantSharedAccessSubmit = async () => {
    if (!grantWallet || !appUser) return;
    if (!grantPassword) {
      setGrantPasswordErr('Please enter your password!');
      return;
    }

    setIsApprovalsLoading(true);
    setGrantPasswordErr('');

    try {
      const encryptor = new Encryptor();
      let decryptedSecretKey: string;
      try {
        decryptedSecretKey = await encryptor.decryptData(
          appUser.secretKeys[0],
          grantPassword,
          appUser.primarySigner,
        );
      } catch (err) {
        setGrantPasswordErr('Invalid password!');
        setIsApprovalsLoading(false);
        return;
      }

      if (!decryptedSecretKey) {
        setGrantPasswordErr('Invalid password!');
        setIsApprovalsLoading(false);
        return;
      }

      const permissions = [
        ...grantViewers.map((username) => ({
          targetUsername: username,
          permission: 'VIEW-ONLY',
        })),
        ...(grantAddApprovers
          ? [
              ...grantApprovers.map((username) => ({
                targetUsername: username,
                permission: 'APPROVER',
              })),
              ...grantInitiators.map((username) => ({
                targetUsername: username,
                permission: 'INITIATOR',
              })),
            ]
          : []),
      ];

      const body = grantAddApprovers
        ? { numberOfApprovalsNeeded: grantNoOfApprovalsNeeded, permissions }
        : { permissions };

      const payload = {
        signer: grantWallet.signer || appUser.primarySigner,
        address: grantWallet.address,
        secretKey: decryptedSecretKey,
        body,
      };

      const response = await addSharedAccess(payload);
      if ('data' in response && response.data) {
        const txData = response.data as any;
        const signature = signBase64Txn(
          decryptedSecretKey,
          txData.transaction,
          txData.networkPassPhrase,
        );

        const secondBody = {
          ...txData,
          transactionId: '',
          transactionSignature: signature,
        };

        const secondPayload = {
          signer: grantWallet.signer || appUser.primarySigner,
          address: grantWallet.address,
          secretKey: decryptedSecretKey,
          body: secondBody,
        };

        const secondResponse = await addSharedAccess(secondPayload);
        if ('data' in secondResponse) {
          showNotification('success', 'Shared access enabled successfully.');
          await refreshUserInfo(decryptedSecretKey, grantPassword);
          setGrantPassword('');
          setActiveModal('success');
        } else {
          const err = (secondResponse as any).error;
          showNotification(
            'error',
            err?.data?.message || 'Unable to submit signed transaction',
          );
        }
      } else {
        const err = (response as any).error;
        showNotification(
          'error',
          err?.data?.message || 'Unable to enable shared access',
        );
      }
    } catch (error: any) {
      showNotification(
        'error',
        error?.message || 'Unable to grant shared access',
      );
    } finally {
      setIsApprovalsLoading(false);
    }
  };

  const handleModifySharedAccessSubmit = async () => {
    if (!selectedRow || !appUser) return;
    if (!modifyPassword) {
      setModifyPasswordErr('Please enter your password!');
      return;
    }

    setIsApprovalsLoading(true);
    setModifyPasswordErr('');

    try {
      const encryptor = new Encryptor();
      let decryptedSecretKey: string;
      try {
        decryptedSecretKey = await encryptor.decryptData(
          appUser.secretKeys[0],
          modifyPassword,
          appUser.primarySigner,
        );
      } catch (err) {
        setModifyPasswordErr('Invalid password!');
        setIsApprovalsLoading(false);
        return;
      }

      if (!decryptedSecretKey) {
        setModifyPasswordErr('Invalid password!');
        setIsApprovalsLoading(false);
        return;
      }

      const wallet = appUser.userWallets.find(
        (w) => w.address === selectedRow.address,
      );
      if (!wallet) {
        setIsApprovalsLoading(false);
        return;
      }

      const addedPermissions: any[] = [];
      const revokedPermissions: any[] = [];

      modifyViewers.forEach((v) => {
        if (v.permissionState === PermissionState.Added) {
          addedPermissions.push({
            targetUsername: v.targetUsername,
            permission: 'VIEW-ONLY',
          });
        } else if (v.permissionState === PermissionState.Revoked) {
          revokedPermissions.push({
            targetUsername: v.targetUsername,
            permission: 'VIEW-ONLY',
          });
        }
      });

      modifyApprovers.forEach((a) => {
        if (a.permissionState === PermissionState.Added) {
          addedPermissions.push({
            targetUsername: a.targetUsername,
            permission: 'APPROVER',
          });
        } else if (a.permissionState === PermissionState.Revoked) {
          revokedPermissions.push({
            targetUsername: a.targetUsername,
            permission: 'APPROVER',
          });
        }
      });

      modifyInitiators.forEach((i) => {
        if (i.permissionState === PermissionState.Added) {
          addedPermissions.push({
            targetUsername: i.targetUsername,
            permission: 'INITIATOR',
          });
        } else if (i.permissionState === PermissionState.Revoked) {
          revokedPermissions.push({
            targetUsername: i.targetUsername,
            permission: 'INITIATOR',
          });
        }
      });

      const body = {
        numberOfApprovalsNeeded: modifyAddApprovers
          ? modifyNoOfApprovalsNeeded
          : 0,
        revokedPermissions,
        addedPermissions,
      };

      const payload = {
        signer: appUser.primarySigner,
        address: wallet.address,
        secretKey: decryptedSecretKey,
        body,
      };

      const response = await updateSharedAccess(payload);
      if ('data' in response && response.data) {
        const txData = { ...(response.data as any) };

        if (txData.signatureRequired === 1) {
          txData.commit = 0;
          const signature = signBase64Txn(
            decryptedSecretKey,
            txData.transaction,
            txData.networkPassPhrase,
          );
          txData.transactionId = '';
          txData.transactionSignature = signature;
        } else {
          txData.commit = 1;
        }

        const secondPayload = {
          signer: appUser.primarySigner,
          address: wallet.address,
          secretKey: decryptedSecretKey,
          body: txData,
        };

        const secondResponse = await updateSharedAccess(secondPayload);
        if ('data' in secondResponse) {
          showNotification('success', 'Shared access updated successfully.');
          await refreshUserInfo(decryptedSecretKey, modifyPassword);
          setModifyPassword('');
          setActiveModal('modifySuccess');
        } else {
          const err = (secondResponse as any).error;
          showNotification(
            'error',
            err?.data?.message || 'Unable to submit update transaction',
          );
        }
      } else {
        const err = (response as any).error;
        showNotification(
          'error',
          err?.data?.message || 'Unable to update shared access',
        );
      }
    } catch (error: any) {
      showNotification(
        'error',
        error?.message || 'Unable to modify shared access',
      );
    } finally {
      setIsApprovalsLoading(false);
    }
  };

  const handleDisableSharedAccessInit = async () => {
    console.log('handleDisableSharedAccessInit');
    if (!selectedRow || !appUser) return;
    if (!disablePassword) {
      setDisablePasswordErr('Please enter your password!');
      return;
    }

    setIsDisabling(true);
    setDisablePasswordErr('');

    try {
      const encryptor = new Encryptor();
      let decryptedSecretKey: string;
      try {
        decryptedSecretKey = await encryptor.decryptData(
          appUser.secretKeys[0],
          disablePassword,
          appUser.primarySigner,
        );
      } catch (err) {
        setDisablePasswordErr('Invalid password!');
        setIsDisabling(false);
        return;
      }

      if (!decryptedSecretKey) {
        setDisablePasswordErr('Invalid password!');
        setIsDisabling(false);
        return;
      }

      const payload = {
        signer: appUser.primarySigner,
        address: selectedRow.address,
        secretKey: decryptedSecretKey,
        body: {},
      };

      const response = await disableSharedAccess(payload);
      if ('data' in response) {
        console.log(response.data, 'response.data');
        setDisableTxData(response.data);
        setDisablePasswordErr('');
        setActiveModal('disableConfirm');
      } else if ('error' in response) {
        const err = response.error as any;
        showNotification(
          'error',
          err?.data?.message || 'Unable to disable shared access',
        );
      }
    } catch (error: any) {
      showNotification(
        'error',
        error?.message || 'Unable to disable shared access',
      );
      setActiveModal('details');
    } finally {
      setIsDisabling(false);
    }
  };

  const handleDisableSharedAccessSubmit = async () => {
    if (!selectedRow || !disableTxData || !appUser) return;
    if (!disablePassword) {
      setDisablePasswordErr('Please enter your password!');
      return;
    }

    setIsDisabling(true);
    setDisablePasswordErr('');

    try {
      const encryptor = new Encryptor();
      let decryptedSecretKey: string;
      try {
        decryptedSecretKey = await encryptor.decryptData(
          appUser.secretKeys[0],
          disablePassword,
          appUser.primarySigner,
        );
      } catch (err) {
        setDisablePasswordErr('Invalid password!');
        setIsDisabling(false);
        return;
      }

      if (!decryptedSecretKey) {
        setDisablePasswordErr('Invalid password!');
        setIsDisabling(false);
        return;
      }

      const signature = signBase64Txn(
        decryptedSecretKey,
        disableTxData.transaction,
        disableTxData.networkPassPhrase,
      );

      const secondBody = {
        ...disableTxData,
        transactionId: '',
        transactionSignature: signature,
        commit: 1,
      };

      const payload = {
        signer: appUser.primarySigner,
        address: selectedRow.address,
        secretKey: decryptedSecretKey,
        body: secondBody,
      };

      const response = await disableSharedAccess(payload);
      if ('data' in response) {
        const viewOnly = selectedAccessMode === 'Access granted to me';
        const walletLabel = selectedRow.wallet || selectedRow.address;
        const msg = viewOnly
          ? `Shared access has successfully been disabled on this wallet [${walletLabel}]`
          : `Your request to disable shared access on wallet [${walletLabel}] has been submitted. This transaction will be completed when it gets the required number of approvals.`;

        setDisableSuccessMessage(msg);
        setDisablePassword('');

        const importedPayload = {
          signer: appUser.address,
          address: appUser.address,
          secretKey: appUser.secretKeys[0],
          body: { userId: appUser.username, import: 1 },
        };
        const { data: importData } = await getUser(importedPayload);
        if (importData) {
          const userData = deserializeUserData(importData);
          const base64EncryptedSecretKey = await encryptor.encryptData(
            decryptedSecretKey,
            disablePassword,
            appUser.address,
          );
          const passwordHash = await encryptor.createHash(userData.username);
          const encryptedPassword = await encryptor.encryptData(
            disablePassword,
            passwordHash,
            appUser.address,
          );
          const updatedUser = {
            ...userData,
            password: encryptedPassword,
            isLoggedIn: true,
            currency: 'USD',
            secretKeys: [
              ...appUser.secretKeys.filter((k) => k !== appUser.secretKeys[0]),
              base64EncryptedSecretKey,
            ],
          };
          const userHash = await encryptor.createHash(updatedUser.username);
          const base64EncryptedUserData = await encryptor.encryptData(
            JSON.stringify(updatedUser),
            userHash,
            updatedUser.address,
          );
          dispatch(
            setUser({
              key: userHash,
              user: updatedUser,
              encryptedUser: base64EncryptedUserData,
            }),
          );
        }

        setActiveModal('disableSuccess');
      } else if ('error' in response) {
        const err = response.error as any;
        showNotification(
          'error',
          err?.data?.message || 'Unable to submit disable transaction',
        );
      }
    } catch (error: any) {
      showNotification(
        'error',
        error?.message || 'Unable to submit disable transaction',
      );
    } finally {
      setIsDisabling(false);
    }
  };

  const handleApproveSharedAccess = async () => {
    if (!selectedActivityRecord || !appUser) return;
    if (!approvalPassword) {
      setApprovalPasswordErr('Please enter your password');
      return;
    }

    setIsApproving(true);
    setApprovalPasswordErr('');

    try {
      const encryptor = new Encryptor();
      let decryptedSecretKey: string;
      try {
        decryptedSecretKey = await encryptor.decryptData(
          appUser.secretKeys[0],
          approvalPassword,
          appUser.primarySigner,
        );
      } catch (err) {
        setApprovalPasswordErr('Invalid password!');
        setIsApproving(false);
        return;
      }

      if (!decryptedSecretKey) {
        setApprovalPasswordErr('Invalid password!');
        setIsApproving(false);
        return;
      }

      const payload = {
        signer: appUser.primarySigner,
        address: appUser.primarySigner, // Uses primary signer as per dart
        secretKey: decryptedSecretKey,
        body: { id: (selectedActivityRecord as any).id }, // Ensure ID is passed
      };

      const response = await approveSharedAccess(payload);

      if ('data' in response && response.data) {
        const txData = response.data as any;
        const signature = signBase64Txn(
          decryptedSecretKey,
          txData.transaction,
          txData.networkPassPhrase,
        );

        const secondPayload = {
          signer: appUser.primarySigner,
          address: appUser.primarySigner,
          secretKey: decryptedSecretKey,
          body: {
            id: (selectedActivityRecord as any).id,
            ...txData,
            transactionSignature: signature,
          },
        };

        const secondResponse = await approveSharedAccess(secondPayload);
        if ('data' in secondResponse) {
          showNotification(
            'success',
            'Transaction approval submitted successfully',
          );
          setActiveModal(null);
          setApprovalPassword('');
          setApprovalPasswordErr('');
          setShowRejectReason(false);
          // Refresh list locally
          setApprovalRecords((prev) =>
            prev.filter(
              (record) =>
                (record as any).id !== (selectedActivityRecord as any).id,
            ),
          );
          // Refresh list
          fetchApprovals({
            signer: appUser.primarySigner,
            address: appUser.address,
            secretKey: appUser.secretKeys[0],
            body: {
              limit: APPROVALS_PER_PAGE,
              excludeUserApproved: includeSignedTransactions ? 0 : 1,
              page: approvalPage,
            },
          });
        } else {
          const err = (secondResponse as any).error;
          showNotification(
            'error',
            err?.data?.message || 'Unable to submit approval transaction',
          );
        }
      } else {
        const err = (response as any).error;
        showNotification(
          'error',
          err?.data?.message || 'Unable to get transaction for approval',
        );
      }
    } catch (error: any) {
      showNotification('error', error?.message || 'Unable to submit approval');
    } finally {
      setIsApproving(false);
    }
  };

  const handleRejectSharedAccess = async () => {
    if (!selectedActivityRecord || !appUser) return;
    if (!showRejectReason) {
      setShowRejectReason(true);
      return;
    }
    if (!approvalPassword) {
      setApprovalPasswordErr('Please enter your password');
      return;
    }
    if (!approvalReason.trim()) {
      showNotification('error', 'Please enter a reason for rejection');
      return;
    }

    setIsRejecting(true);
    setApprovalPasswordErr('');

    try {
      const encryptor = new Encryptor();
      let decryptedSecretKey: string;
      try {
        decryptedSecretKey = await encryptor.decryptData(
          appUser.secretKeys[0],
          approvalPassword,
          appUser.primarySigner,
        );
      } catch (err) {
        setApprovalPasswordErr('Invalid password!');
        setIsRejecting(false);
        return;
      }

      if (!decryptedSecretKey) {
        setApprovalPasswordErr('Invalid password!');
        setIsRejecting(false);
        return;
      }

      const payload = {
        signer: appUser.primarySigner,
        address: appUser.primarySigner,
        secretKey: decryptedSecretKey,
        body: {
          id: (selectedActivityRecord as any).id,
          rejectionReason: approvalReason,
        },
      };

      const response = await rejectSharedAccess(payload);

      if ('data' in response) {
        showNotification('success', 'Rejection submitted successfully');
        setActiveModal(null);
        setApprovalPassword('');
        setApprovalPasswordErr('');
        setApprovalReason('');
        setShowRejectReason(false);
        // Refresh list locally
        setApprovalRecords((prev) =>
          prev.filter(
            (record) =>
              (record as any).id !== (selectedActivityRecord as any).id,
          ),
        );
        // Refresh list
        fetchApprovals({
          signer: appUser.primarySigner,
          address: appUser.address,
          secretKey: appUser.secretKeys[0],
          body: {
            limit: APPROVALS_PER_PAGE,
            excludeUserApproved: includeSignedTransactions ? 0 : 1,
            page: approvalPage,
          },
        });
      } else {
        const err = (response as any).error;
        showNotification(
          'error',
          err?.data?.message || 'Unable to reject transaction',
        );
      }
    } catch (error: any) {
      showNotification('error', error?.message || 'Unable to submit rejection');
    } finally {
      setIsRejecting(false);
    }
  };

  const debouncedInitiator = selectedInitiator;
  const debouncedWalletAlias = selectedWalletAlias;
  const debouncedWalletAddress = selectedWalletAddress;
  const debouncedDescription = selectedDescription;

  const toggleView = () => {
    setShowSharedAccessView((prev) => !prev);
  };

  const transactionTypeOptions = useMemo(() => {
    const baseValues = new Set(
      transactionTypeBaseOptions.map((option) => option.value),
    );
    const dynamicValues = new Set(
      approvalRecords
        .map((record) => record.transactionType?.toUpperCase())
        .filter((value): value is string => Boolean(value))
        .filter((value) => !baseValues.has(value)),
    );
    const dynamicOptions = Array.from(dynamicValues, (value) => ({
      text: normalizeTransactionType(value),
      value,
    }));

    return [...transactionTypeBaseOptions, ...dynamicOptions];
  }, [approvalRecords]);

  const dateRangeLabel = useMemo(() => {
    if (
      selectedDateRange === 'Custom date' &&
      customStartDate &&
      customEndDate
    ) {
      return formatDateRangeLabel(customStartDate, customEndDate);
    }
    return selectedDateRange === 'All' ? 'Date range' : selectedDateRange;
  }, [customEndDate, customStartDate, selectedDateRange]);

  const resetApprovalFilters = () => {
    pendingFallbackAppliedRef.current = false;
    setApprovalPage(1);
    setSelectedTransactionType('All');
    setSelectedDateRange('All');
    setSelectedTransactionStatus('Pending');
    setSelectedInitiator('');
    setSelectedWalletAlias('');
    setSelectedWalletAddress('');
    setSelectedDescription('');
    setCustomStartDate('');
    setCustomEndDate('');
    setDraftCustomStartDate('');
    setDraftCustomEndDate('');
    setDraftInitiator('');
    setDraftWalletAlias('');
    setDraftWalletAddress('');
    setDraftDescription('');
  };

  const activeFilters = useMemo(() => {
    const chips: Array<{ label: string; onRemove: () => void }> = [];
    if (selectedTransactionType !== 'All') {
      chips.push({
        label: `Type: ${normalizeTransactionType(selectedTransactionType)}`,
        onRemove: () => setSelectedTransactionType('All'),
      });
    }
    if (selectedDateRange !== 'All') {
      chips.push({
        label: `Date: ${dateRangeLabel}`,
        onRemove: () => {
          setSelectedDateRange('All');
          setCustomStartDate('');
          setCustomEndDate('');
          setDraftCustomStartDate('');
          setDraftCustomEndDate('');
        },
      });
    }
    if (selectedTransactionStatus !== 'All') {
      chips.push({
        label: `Status: ${selectedTransactionStatus}`,
        onRemove: () => setSelectedTransactionStatus('All'),
      });
    }
    if (selectedInitiator) {
      chips.push({
        label: `Initiator: ${selectedInitiator}`,
        onRemove: () => {
          setSelectedInitiator('');
          setDraftInitiator('');
        },
      });
    }
    if (selectedWalletAlias) {
      chips.push({
        label: `Wallet: ${selectedWalletAlias}`,
        onRemove: () => {
          setSelectedWalletAlias('');
          setDraftWalletAlias('');
        },
      });
    }
    if (selectedWalletAddress) {
      chips.push({
        label: `Key: ${selectedWalletAddress}`,
        onRemove: () => {
          setSelectedWalletAddress('');
          setDraftWalletAddress('');
        },
      });
    }
    if (selectedDescription) {
      chips.push({
        label: `Desc: ${selectedDescription}`,
        onRemove: () => {
          setSelectedDescription('');
          setDraftDescription('');
        },
      });
    }
    return chips;
  }, [
    dateRangeLabel,
    selectedTransactionType,
    selectedDateRange,
    selectedTransactionStatus,
    selectedInitiator,
    selectedWalletAlias,
    selectedWalletAddress,
    selectedDescription,
  ]);

  useEffect(() => {
    setApprovalPage(1);
  }, [
    customEndDate,
    customStartDate,
    debouncedDescription,
    debouncedInitiator,
    debouncedWalletAlias,
    debouncedWalletAddress,
    selectedDateRange,
    selectedTransactionStatus,
    selectedTransactionType,
  ]);

  useEffect(() => {
    if (!appUser) {
      setApprovalRecords([]);
      setApprovalTotalRecords(0);
      return;
    }

    let isCurrentRequest = true;

    const loadApprovals = async () => {
      const queryParts: string[] = [];

      setIsApprovalsLoading(true);

      if (selectedTransactionStatus !== 'All') {
        queryParts.push(
          `transactionStatus=${encodeURIComponent(selectedTransactionStatus.toUpperCase())}`,
        );
      }

      if (selectedTransactionType !== 'All') {
        queryParts.push(
          `transactionType=${encodeURIComponent(selectedTransactionType)}`,
        );
      }

      if (debouncedInitiator.trim()) {
        queryParts.push(
          `initiator=${encodeURIComponent(debouncedInitiator.trim())}`,
        );
      }

      if (debouncedWalletAlias.trim()) {
        queryParts.push(
          `walletAlias=${encodeURIComponent(debouncedWalletAlias.trim())}`,
        );
      }

      if (debouncedWalletAddress.trim()) {
        queryParts.push(
          `walletAddress=${encodeURIComponent(debouncedWalletAddress.trim())}`,
        );
      }

      if (debouncedDescription.trim()) {
        queryParts.push(
          `description=${encodeURIComponent(debouncedDescription.trim())}`,
        );
      }

      const dateRangeQueryValue = getDateRangeQueryValue(selectedDateRange);
      if (dateRangeQueryValue) {
        queryParts.push(
          `dateBetween=${encodeURIComponent(dateRangeQueryValue)}`,
        );
      }

      if (
        selectedDateRange === 'Custom date' &&
        customStartDate &&
        customEndDate
      ) {
        queryParts.push(
          `dateBetween=${encodeURIComponent(`${customStartDate}|${customEndDate}`)}`,
        );
      }

      const query = queryParts.length > 0 ? `&${queryParts.join('&')}` : '';

      try {
        const result = await fetchApprovals({
          signer: appUser.primarySigner,
          address: appUser.address,
          secretKey: appUser.secretKeys[0],
          body: {
            limit: APPROVALS_PER_PAGE,
            excludeUserApproved: includeSignedTransactions ? 0 : 1,
            page: approvalPage,
            query,
          },
        });

        if (!isCurrentRequest) return;

        if ('error' in result) {
          showNotification('error', 'Unable to fetch pending approvals');
          return;
        }

        const responseData = result.data as any;
        const records = Array.isArray(responseData)
          ? responseData
          : Array.isArray(responseData?.records)
            ? responseData.records
            : Array.isArray(responseData?.approvals)
              ? responseData.approvals
              : [];

        const totalRecords =
          typeof responseData?.totalRecords === 'number'
            ? responseData.totalRecords
            : typeof responseData?.pagination?.totalRecords === 'number'
              ? responseData.pagination.totalRecords
              : records.length;

        if (
          selectedTransactionStatus === 'Pending' &&
          totalRecords === 0 &&
          !pendingFallbackAppliedRef.current
        ) {
          pendingFallbackAppliedRef.current = true;
          setApprovalPage(1);
          setSelectedTransactionStatus('All');
          return;
        }

        pendingFallbackAppliedRef.current = false;
        setApprovalRecords(records);
        setApprovalTotalRecords(totalRecords);
      } catch (error) {
        if (isCurrentRequest) {
          showNotification('error', 'Unable to fetch pending approvals');
        }
      } finally {
        if (isCurrentRequest) {
          setIsApprovalsLoading(false);
        }
      }
    };

    loadApprovals();

    return () => {
      isCurrentRequest = false;
    };
  }, [
    approvalPage,
    appUser,
    customEndDate,
    customStartDate,
    debouncedDescription,
    debouncedInitiator,
    debouncedWalletAlias,
    debouncedWalletAddress,
    fetchApprovals,
    selectedTransactionStatus,
    selectedTransactionType,
    selectedDateRange,
    includeSignedTransactions,
  ]);

  const activityRows = useMemo<
    (SharedAccessActivityRow & { _record: SharedAccessApprovalRecord })[]
  >(() => {
    return approvalRecords.map((record) => {
      const approvalsGotten = Number(record.approvalsGotten ?? 0);
      const approvalsNeeded = Number(record.approvalsNeeded ?? 0);
      const fallbackApprovalStatus =
        approvalsNeeded > 0
          ? `${approvalsGotten}/${approvalsNeeded} approved`
          : '-';

      return {
        wallet: record.alias || record.walletAlias || '-',
        transactionType: normalizeTransactionType(record.transactionType),
        initiatedBy: record.initiator || '-',
        date: formatApprovalDate(record.createdAt),
        transactionStatus: normalizeTransactionStatus(record.transactionStatus),
        approvalStatus: record.approvalStatus || fallbackApprovalStatus,
        _record: record,
      };
    });
  }, [approvalRecords]);

  const sharedWalletRows = useMemo<SharedAccessRow[]>(() => {
    if (!appUser) {
      return [];
    }

    const selectedPermissionCode = filterToPermissionCodeMap[selectedFilter];

    if (selectedAccessMode === 'Access granted to me') {
      const walletsMap = new Map<string, SharedAccessRow>();

      for (const wallet of appUser.userWallets) {
        if (!wallet.isSharedWallet) {
          continue;
        }

        const permissionCodes = new Set<string>();

        if (wallet.permission) {
          permissionCodes.add(wallet.permission.toUpperCase());
        }

        for (const permission of wallet.permissions ?? []) {
          if (permission.targetUsername === appUser.username) {
            permissionCodes.add(permission.permission.toUpperCase());
          }
        }

        if (
          selectedPermissionCode &&
          !permissionCodes.has(selectedPermissionCode)
        ) {
          continue;
        }

        const walletKey = wallet.alias || wallet.address;
        const permissionLabels = Array.from(permissionCodes).map(
          formatPermissionLabel,
        );

        if (!walletsMap.has(walletKey)) {
          walletsMap.set(walletKey, {
            wallet: wallet.alias || wallet.address,
            owner: wallet.owner,
            permissions: permissionLabels,
            description: wallet.description || '-',
            address: wallet.address,
          });
          continue;
        }

        const existing = walletsMap.get(walletKey)!;
        existing.permissions = Array.from(
          new Set([...existing.permissions, ...permissionLabels]),
        );
      }

      return Array.from(walletsMap.values()).sort((rowA, rowB) =>
        rowA.wallet.localeCompare(rowB.wallet),
      );
    }

    return appUser.userWallets
      .filter((wallet) => {
        const isOwnedByUser =
          !wallet.owner || wallet.owner === appUser.username;
        return (
          wallet.isSharedWallet &&
          wallet.walletThreshold === 1 &&
          (wallet.permissions?.length ?? 0) > 0 &&
          isOwnedByUser
        );
      })
      .map((wallet) => {
        const permissionCodes = Array.from(
          new Set(
            (wallet.permissions ?? []).map((permission) =>
              permission.permission.toUpperCase(),
            ),
          ),
        );

        const filteredPermissionCodes = selectedPermissionCode
          ? permissionCodes.filter(
              (permissionCode) => permissionCode === selectedPermissionCode,
            )
          : permissionCodes;

        return {
          wallet: wallet.alias || wallet.address,
          owner: wallet.owner || appUser.username,
          permissions: filteredPermissionCodes.map(formatPermissionLabel),
          description: wallet.description || '-',
          address: wallet.address,
        };
      })
      .filter((row) => row.permissions.length > 0)
      .sort((rowA, rowB) => rowA.wallet.localeCompare(rowB.wallet));
  }, [appUser, selectedAccessMode, selectedFilter]);

  const renderPermissionChip = (
    perm: Permission,
    index: number,
    permList: Permission[],
    setPermList: React.Dispatch<React.SetStateAction<Permission[]>>,
    relation: string,
  ) => {
    let chipClass =
      'inline-flex items-center gap-1 rounded px-2 py-1 text-sm font-medium transition cursor-pointer ';
    if (perm.permissionState === PermissionState.Revoked) {
      chipClass += 'bg-red-500 text-white hover:bg-red-600 line-through';
    } else if (perm.permissionState === PermissionState.Added) {
      chipClass += 'bg-green-500 text-white hover:bg-green-600';
    } else {
      chipClass += 'bg-primary-800 text-white hover:bg-white-700';
    }

    const isOwner = perm.targetUsername === selectedRow?.owner;

    return (
      <span
        key={`${perm.targetUsername}-${index}`}
        className={chipClass}
        onClick={() => {
          if (
            isOwner &&
            (relation === 'Approver' || relation === 'Initiator')
          ) {
            showNotification(
              'error',
              'The owner cannot be removed/revoked from Approver or Initiator lists.',
            );
            return;
          }
          handlePermissionItemClick(permList, setPermList, index, relation);
        }}
        title={
          isOwner && (relation === 'Approver' || relation === 'Initiator')
            ? 'Owner cannot be revoked'
            : perm.permissionState === PermissionState.Revoked
              ? 'Click to restore access'
              : 'Click to revoke access'
        }
      >
        {perm.targetUsername} [{perm.fullName}]
        {!(isOwner && (relation === 'Approver' || relation === 'Initiator')) &&
          (perm.permissionState === PermissionState.Revoked ? ' ⟲' : '×')}
      </span>
    );
  };

  const renderModalContent = () => {
    if (activeModal === 'approvalDetails' && selectedActivityRecord) {
      const record = selectedActivityRecord;
      const isPending =
        normalizeTransactionStatus(record.transactionStatus) === 'Pending';
      console.log('record', record);

      const approvalsGotten = Number(record.approvalsGotten ?? 0);
      const approvalsNeeded = Number(record.approvalsNeeded ?? 0);
      const approvalStatusText =
        approvalsNeeded > 0
          ? `${approvalsGotten} out of ${approvalsNeeded} approvals received`
          : '-';

      const hasApproved = (record as any).approvedBy?.includes(
        appUser?.username || '',
      );
      const hasRejected = (record as any).rejectedBy?.includes(
        appUser?.username || '',
      );
      const hasSigned = hasApproved || hasRejected;

      const relatedWallet = appUser?.userWallets.find(
        (w) => w.address === record.walletAlias || w.alias === record.alias,
      );
      const isApprover = relatedWallet?.isApprover ?? false;

      let headline = '';
      if (isPending) {
        headline = hasSigned
          ? 'You have already signed'
          : 'Your approval is requested';
      } else if (
        normalizeTransactionStatus(record.transactionStatus) === 'Rejected'
      ) {
        headline = 'This transaction is rejected';
      } else if (
        normalizeTransactionStatus(record.transactionStatus) === 'Completed'
      ) {
        headline = 'Transaction completed';
      }

      return (
        <div className="p-6 md:p-8 max-h-[85vh] overflow-y-auto">
          <div className="space-y-5">
            <div className="flex items-center justify-center">
              <h3 className="text-2xl font-montserratSemiBold text-primary-800 text-center">
                Approve Request
              </h3>
            </div>

            {headline && (
              <p className="text-sm font-medium text-primary-700 mb-2 text-center">
                {headline}
              </p>
            )}

            <div className="rounded-xl border border-primary-200 bg-[#f5f9fc] p-5 space-y-4">
              <div className="flex flex-col text-center">
                <span className="text-sm font-bold text-primary-800">
                  Wallet:
                </span>
                <span className="text-sm text-primary-700 break-words">
                  {record.alias || record.walletAlias || '-'}
                </span>
              </div>
              <div className="flex flex-col text-center">
                <span className="text-sm font-bold text-primary-800">
                  Transaction Type:
                </span>
                <span className="text-sm text-primary-700">
                  {normalizeTransactionType(record.transactionType)}
                </span>
              </div>
              <div className="flex flex-col text-center">
                <span className="text-sm font-bold text-primary-800">
                  Initiator:
                </span>
                <span className="text-sm text-primary-700 break-words">
                  {record.initiator || '-'}
                </span>
              </div>
              <div className="flex flex-col text-center">
                <span className="text-sm font-bold text-primary-800">
                  Initiated:
                </span>
                <span className="text-sm text-primary-700">
                  {formatApprovalDate(record.createdAt)}
                </span>
              </div>
              <div className="flex flex-col text-center">
                <span className="text-sm font-bold text-primary-800">
                  Description:
                </span>
                <span className="text-sm text-primary-700 whitespace-pre-wrap break-words">
                  {(record as any).description || '-'}
                </span>
              </div>
              <div className="flex flex-col text-center">
                <span className="text-sm font-bold text-primary-800">
                  Approval Status:
                </span>
                <span className="text-sm text-primary-700">
                  {approvalStatusText}
                </span>
              </div>

              {(record as any).approvedBy && (
                <div className="flex flex-col text-center">
                  <span className="text-sm font-bold text-primary-800">
                    Approved By:
                  </span>
                  <span className="text-sm text-primary-700 break-words">
                    {(record as any).approvedBy}
                  </span>
                </div>
              )}

              {(record as any).rejectedBy && (
                <>
                  <div className="flex flex-col text-center">
                    <span className="text-sm font-bold text-primary-800">
                      Rejected By:
                    </span>
                    <span className="text-sm text-primary-700 break-words">
                      {(record as any).rejectedBy}
                    </span>
                  </div>
                  <div className="flex flex-col text-center">
                    <span className="text-sm font-bold text-primary-800">
                      Reason for Rejection:
                    </span>
                    <span className="text-sm text-primary-700 break-words">
                      {(record as any).reasonForRejection || '-'}
                    </span>
                  </div>
                </>
              )}

              <div className="flex flex-col text-center">
                <span className="text-sm font-bold text-primary-800">
                  Transaction Status:
                </span>
                <span className="text-sm text-primary-700">
                  {normalizeTransactionStatus(record.transactionStatus)}
                </span>
              </div>

              {/* <div className="flex flex-col text-center">
                <span className="text-sm font-bold text-primary-800">Blockchain Proof:</span>
                <div className="flex items-center justify-center gap-2">
                  <a
                    href={`${getExplorerBaseUrl(appState.walletMode)}${(record as any).id}`}
                    target="_blank"
                    rel="noreferrer"
                    className="text-sm text-primary-600 underline break-all hover:text-primary-800"
                  >
                    {(record as any).id}
                  </a>
                  <button
                    onClick={() => {
                      navigator.clipboard.writeText((record as any).id);
                      showNotification('success', 'Transaction ID copied!');
                    }}
                    className="text-primary-600 hover:text-primary-800 shrink-0"
                    title="Copy Transaction ID"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                      <path d="M8 3a1 1 0 011-1h2a1 1 0 110 2H9a1 1 0 01-1-1z" />
                      <path d="M6 3a2 2 0 00-2 2v11a2 2 0 002 2h8a2 2 0 002-2V5a2 2 0 00-2-2 3 3 0 01-3 3H9a3 3 0 01-3-3z" />
                    </svg>
                  </button>
                </div>
              </div> */}
            </div>

            {isPending && !hasSigned && isApprover ? (
              <div className="space-y-4 pt-4 border-t border-primary-100">
                <div className="space-y-1">
                  <label className="block text-sm font-medium text-primary-700 text-center">
                    Password
                  </label>
                  <div className="relative">
                    <input
                      type={showApprovalPassword ? 'text' : 'password'}
                      className="h-12 w-full rounded-md border border-[#cfe0ee] bg-white px-10 text-primary-800 outline-none"
                      placeholder="Enter your password"
                      value={approvalPassword}
                      onChange={(e) => setApprovalPassword(e.target.value)}
                    />
                    <svg
                      className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-[#88a5bc]"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                      />
                    </svg>
                    <button
                      type="button"
                      onClick={() =>
                        setShowApprovalPassword(!showApprovalPassword)
                      }
                      className="absolute right-3 top-1/2 -translate-y-1/2 text-sm font-medium text-primary-600 hover:text-primary-800"
                    >
                      {showApprovalPassword ? 'Hide' : 'Show'}
                    </button>
                  </div>
                  {approvalPasswordErr && (
                    <p className="text-xs text-red-500">
                      {approvalPasswordErr}
                    </p>
                  )}
                </div>

                {showRejectReason && (
                  <div className="space-y-1">
                    <label className="block text-sm font-medium text-primary-700 text-center">
                      Rejection Reason (Required)
                    </label>
                    <textarea
                      className="w-full rounded-md border border-[#cfe0ee] bg-white p-3 text-primary-800 outline-none resize-none"
                      placeholder="Enter reason if rejecting"
                      rows={2}
                      value={approvalReason}
                      onChange={(e) => setApprovalReason(e.target.value)}
                    />
                  </div>
                )}

                <div className="grid grid-cols-2 gap-4 pt-2">
                  <ButtonSecondary
                    label={isRejecting ? 'Rejecting...' : 'Reject'}
                    onclick={handleRejectSharedAccess}
                    disabled={isApproving || isRejecting}
                    additionalClasses="w-full border-red-500 text-red-500 hover:bg-red-50"
                  />
                  <Button
                    label={isApproving ? 'Approving...' : 'Approve'}
                    onclick={handleApproveSharedAccess}
                    disabled={isApproving || isRejecting}
                  />
                </div>
              </div>
            ) : (
              <div className="pt-4 flex justify-end">
                <Button label="Done" onclick={() => setActiveModal(null)} />
              </div>
            )}
          </div>
        </div>
      );
    }

    if (activeModal === 'grant') {
      return (
        <div className="p-6 md:p-8">
          {grantStep === 0 && (
            <div className="space-y-5">
              <div className="flex items-center justify-between">
                <h3 className="text-2xl font-montserratSemiBold text-primary-800">
                  Add viewer access
                </h3>
              </div>

              <div className="rounded-xl border border-primary-200 bg-[#f5f9fc] p-4">
                <label className="mb-2 block text-sm font-medium text-primary-700">
                  Choose Wallet
                </label>
                <select
                  className="h-12 w-full rounded-md border border-[#b5cfe4] bg-white px-3 text-base text-primary-800 outline-none"
                  value={grantWallet?.address || ''}
                  onChange={(e) => {
                    const val = e.target.value;
                    const found =
                      shareableWallets.find((w) => w.address === val) || null;
                    setGrantWallet(found);
                    setGrantAddApprovers(false);
                    setGrantViewers([]);
                    if (appUser) {
                      setGrantApprovers([appUser.username]);
                      setGrantInitiators([appUser.username]);
                      setGrantUserFullnames({
                        [appUser.username]: `${appUser.firstName} ${appUser.lastName}`,
                      });
                    }
                  }}
                >
                  {shareableWallets.map((wallet) => (
                    <option key={wallet.address} value={wallet.address}>
                      {wallet.alias || wallet.address}
                    </option>
                  ))}
                </select>
              </div>

              <div className="space-y-2">
                <label className="block text-sm font-medium text-primary-700">
                  Enter the usernames of all accounts that require viewer access
                  to this wallet
                </label>
                <div className="flex flex-wrap gap-2 rounded-lg border border-dashed border-[#b7d0df] bg-[#f7fbff] p-3 min-h-[3rem] items-center">
                  {grantViewers.length > 0 ? (
                    grantViewers.map((username, idx) => (
                      <span
                        key={username}
                        className="inline-flex items-center gap-1 rounded bg-[#dfeaf4] px-2 py-1 text-sm text-primary-800 font-medium"
                      >
                        {username} [{grantUserFullnames[username] || ''}]
                        <button
                          type="button"
                          onClick={() =>
                            setGrantViewers((prev) =>
                              prev.filter((_, i) => i !== idx),
                            )
                          }
                          className="text-red-500 hover:text-red-700 font-bold ml-1 leading-none"
                        >
                          ×
                        </button>
                      </span>
                    ))
                  ) : (
                    <span className="text-sm text-primary-400">
                      Name of viewers appear here
                    </span>
                  )}
                </div>
              </div>

              <div className="space-y-3">
                <div className="flex gap-2">
                  <input
                    type="text"
                    value={grantViewerInput}
                    onChange={(e) => setGrantViewerInput(e.target.value)}
                    placeholder="Viewer username"
                    className="h-12 w-full rounded-md border border-[#cfe0ee] bg-white px-3 text-primary-800 outline-none"
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') {
                        e.preventDefault();
                        handleAddGrantViewer();
                      }
                    }}
                  />
                  <button
                    type="button"
                    onClick={handleAddGrantViewer}
                    className="h-12 rounded-md bg-primary-800 px-5 text-sm font-medium text-white hover:bg-primary-700 transition"
                  >
                    Add
                  </button>
                </div>
                {grantViewerErr && (
                  <p className="text-xs text-red-500">{grantViewerErr}</p>
                )}
                {grantWallet &&
                  !grantWallet.isPrimaryWallet &&
                  !grantWallet.primaryWallet &&
                  grantWallet.walletType !== 2 && (
                    <label className="flex items-center gap-2 text-sm text-primary-700 cursor-pointer">
                      <input
                        type="checkbox"
                        checked={grantAddApprovers}
                        onChange={(e) => {
                          setGrantAddApprovers(e.target.checked);
                          if (e.target.checked) {
                            setGrantStep(1);
                          }
                        }}
                        className="h-4 w-4 rounded border-[#b5cfe4] text-primary-800 focus:ring-primary-800"
                      />
                      Grant approver access
                    </label>
                  )}
              </div>

              <Button
                label="Proceed"
                onclick={() => {
                  if (!grantAddApprovers && grantViewers.length === 0) {
                    showNotification(
                      'error',
                      'Please add at least one user to grant access',
                    );
                    return;
                  }
                  if (grantAddApprovers) {
                    setGrantStep(1);
                  } else {
                    setActiveModal('confirm');
                  }
                }}
                additionalClasses="mt-2"
              />
            </div>
          )}

          {grantStep === 1 && (
            <div className="space-y-5">
              <div className="flex items-center justify-between">
                <h3 className="text-2xl font-montserratSemiBold text-primary-800">
                  Add approver access
                </h3>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="rounded-xl border border-primary-200 bg-[#f5f9fc] p-3">
                  <label className="mb-1 block text-xs font-medium text-primary-700">
                    Approvals Needed
                  </label>
                  <select
                    className="h-10 w-full rounded-md border border-[#b5cfe4] bg-white px-2 text-sm text-primary-800 outline-none"
                    value={grantNoOfApprovalsNeeded}
                    onChange={(e) =>
                      setGrantNoOfApprovalsNeeded(Number(e.target.value))
                    }
                  >
                    {Array.from(
                      { length: grantNoOfApprovers - 1 },
                      (_, i) => i + 2,
                    ).map((val) => (
                      <option key={val} value={val}>
                        {val}
                      </option>
                    ))}
                  </select>
                </div>

                <div className="rounded-xl border border-primary-200 bg-[#f5f9fc] p-3">
                  <label className="mb-1 block text-xs font-medium text-primary-700">
                    Total Approvers
                  </label>
                  <select
                    className="h-10 w-full rounded-md border border-[#b5cfe4] bg-white px-2 text-sm text-primary-800 outline-none"
                    value={grantNoOfApprovers}
                    onChange={(e) => {
                      const val = Number(e.target.value);
                      setGrantNoOfApprovers(val);
                      if (grantNoOfApprovalsNeeded >= val) {
                        setGrantNoOfApprovalsNeeded(val - 1);
                      }
                    }}
                  >
                    {[3, 4, 5, 6, 7, 8, 9, 10].map((val) => (
                      <option key={val} value={val}>
                        {val}
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              <div className="rounded-xl bg-[#f5f9fc] border border-primary-200 p-4 text-center text-xs text-primary-700 font-medium">
                {grantNoOfApprovalsNeeded} approvals required out of{' '}
                {grantNoOfApprovers} approvers
              </div>

              <div className="space-y-2">
                <label className="block text-sm font-medium text-primary-700">
                  Enter the usernames of all accounts that require approver
                  access to this wallet
                </label>
                <div className="flex flex-wrap gap-2 rounded-lg border border-dashed border-[#b7d0df] bg-[#f7fbff] p-3 min-h-[3rem] items-center">
                  {grantApprovers.length > 0 ? (
                    grantApprovers.map((username, idx) => (
                      <span
                        key={username}
                        className="inline-flex items-center gap-1 rounded bg-[#dfeaf4] px-2 py-1 text-sm text-primary-800 font-medium"
                      >
                        {username} [{grantUserFullnames[username] || ''}]
                        <button
                          type="button"
                          onClick={() =>
                            setGrantApprovers((prev) =>
                              prev.filter((_, i) => i !== idx),
                            )
                          }
                          className="text-red-500 hover:text-red-700 font-bold ml-1 leading-none"
                        >
                          ×
                        </button>
                      </span>
                    ))
                  ) : (
                    <span className="text-sm text-primary-400">
                      Name of approvers appear here
                    </span>
                  )}
                </div>
              </div>

              <div className="space-y-3">
                <div className="flex gap-2">
                  <input
                    type="text"
                    value={grantApproverInput}
                    onChange={(e) => setGrantApproverInput(e.target.value)}
                    placeholder="Approver username"
                    className="h-12 w-full rounded-md border border-[#cfe0ee] bg-white px-3 text-primary-800 outline-none"
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') {
                        e.preventDefault();
                        handleAddGrantApprover();
                      }
                    }}
                  />
                  <button
                    type="button"
                    onClick={handleAddGrantApprover}
                    className="h-12 rounded-md bg-primary-800 px-5 text-sm font-medium text-white hover:bg-primary-700 transition"
                  >
                    Add
                  </button>
                </div>
                {grantApproverErr && (
                  <p className="text-xs text-red-500">{grantApproverErr}</p>
                )}
              </div>

              <div className="flex gap-3 mt-4">
                <ButtonSecondary
                  label="Back"
                  onclick={() => setGrantStep(0)}
                  additionalClasses="flex-1"
                />
                <Button
                  label="Proceed"
                  onclick={() => {
                    if (grantApprovers.length === 0) {
                      showNotification(
                        'error',
                        'Please add at least one user to grant access',
                      );
                      return;
                    }
                    if (grantApprovers.length < grantNoOfApprovers) {
                      showNotification(
                        'error',
                        'Please add the required number of approvers',
                      );
                      return;
                    }
                    setGrantStep(2);
                  }}
                  additionalClasses="flex-1"
                />
              </div>
            </div>
          )}

          {grantStep === 2 && (
            <div className="space-y-5">
              <div className="flex items-center justify-between">
                <h3 className="text-2xl font-montserratSemiBold text-primary-800">
                  Add initiator access
                </h3>
              </div>

              <div className="space-y-2">
                <label className="block text-sm font-medium text-primary-700">
                  Enter the usernames of all accounts that require initiator
                  access to this wallet
                </label>
                <div className="flex flex-wrap gap-2 rounded-lg border border-dashed border-[#b7d0df] bg-[#f7fbff] p-3 min-h-[3rem] items-center">
                  {grantInitiators.length > 0 ? (
                    grantInitiators.map((username, idx) => (
                      <span
                        key={username}
                        className="inline-flex items-center gap-1 rounded bg-[#dfeaf4] px-2 py-1 text-sm text-primary-800 font-medium"
                      >
                        {username} [{grantUserFullnames[username] || ''}]
                        {username !== appUser?.username && (
                          <button
                            type="button"
                            onClick={() =>
                              setGrantInitiators((prev) =>
                                prev.filter((_, i) => i !== idx),
                              )
                            }
                            className="text-red-500 hover:text-red-700 font-bold ml-1 leading-none"
                          >
                            ×
                          </button>
                        )}
                      </span>
                    ))
                  ) : (
                    <span className="text-sm text-primary-400">
                      Name of initiators appear here
                    </span>
                  )}
                </div>
              </div>

              <div className="space-y-3">
                <div className="flex gap-2">
                  <input
                    type="text"
                    value={grantInitiatorInput}
                    onChange={(e) => setGrantInitiatorInput(e.target.value)}
                    placeholder="Initiator username"
                    className="h-12 w-full rounded-md border border-[#cfe0ee] bg-white px-3 text-primary-800 outline-none"
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') {
                        e.preventDefault();
                        handleAddGrantInitiator();
                      }
                    }}
                  />
                  <button
                    type="button"
                    onClick={handleAddGrantInitiator}
                    className="h-12 rounded-md bg-primary-800 px-5 text-sm font-medium text-white hover:bg-primary-700 transition"
                  >
                    Add
                  </button>
                </div>
                {grantInitiatorErr && (
                  <p className="text-xs text-red-500">{grantInitiatorErr}</p>
                )}
              </div>

              <div className="flex gap-3 mt-4">
                <ButtonSecondary
                  label="Back"
                  onclick={() => setGrantStep(1)}
                  additionalClasses="flex-1"
                />
                <Button
                  label="Proceed"
                  onclick={() => {
                    if (grantInitiators.length === 0) {
                      showNotification(
                        'error',
                        'Please add at least one user to grant access',
                      );
                      return;
                    }
                    if (grantApprovers.length < grantNoOfApprovers) {
                      showNotification(
                        'error',
                        'Please add the required number of approvers',
                      );
                      return;
                    }
                    if (grantApprovers.length === 0) {
                      showNotification(
                        'error',
                        'You cannot have initiators without approvers',
                      );
                      return;
                    }
                    setActiveModal('confirm');
                  }}
                  additionalClasses="flex-1"
                />
              </div>
            </div>
          )}
        </div>
      );
    }

    if (activeModal === 'confirm') {
      return (
        <div className="p-6 md:p-8">
          <div className="space-y-5">
            <div className="flex items-center justify-between">
              <h3 className="text-[2rem] font-montserratSemiBold text-primary-800">
                Confirm Request
              </h3>
            </div>

            <p className="text-center text-base text-primary-700 font-medium">
              You&rsquo;re about to grant access to the following users
            </p>

            <div className="rounded-xl border border-primary-100 bg-primary-100 p-4 text-center">
              <label className="mb-2 block text-sm font-semibold text-primary-700">
                Wallet
              </label>
              <div className="rounded-lg border border-[#bed1e3] bg-[#e9eff7] px-3 py-3 text-center text-base font-semibold text-primary-800">
                {grantWallet?.alias || grantWallet?.address}
              </div>
            </div>

            <div className="space-y-3">
              {grantViewers.length > 0 && (
                <div className="rounded-xl border border-primary-100 bg-primary-100 p-4">
                  <label className="mb-2 block text-sm font-semibold text-primary-700">
                    Viewer Access
                  </label>
                  <div className="flex flex-wrap gap-2 text-sm">
                    {grantViewers.map((username) => (
                      <span
                        key={username}
                        className="rounded bg-[#dfeaf4] px-2 py-1 text-primary-800 font-medium"
                      >
                        {username} [{grantUserFullnames[username] || ''}]
                      </span>
                    ))}
                  </div>
                </div>
              )}

              {grantAddApprovers && grantApprovers.length > 0 && (
                <div className="rounded-xl border border-primary-100 bg-primary-100 p-4">
                  <label className="mb-2 block text-sm font-semibold text-primary-700">
                    Approver Access
                  </label>
                  <div className="flex flex-wrap gap-2 text-sm">
                    {grantApprovers.map((username) => (
                      <span
                        key={username}
                        className="rounded bg-[#dfeaf4] px-2 py-1 text-primary-800 font-medium"
                      >
                        {username} [{grantUserFullnames[username] || ''}]
                      </span>
                    ))}
                  </div>
                </div>
              )}

              {grantAddApprovers && grantInitiators.length > 0 && (
                <div className="rounded-xl border border-primary-100 bg-primary-100 p-4">
                  <label className="mb-2 block text-sm font-semibold text-primary-700">
                    Initiator Access
                  </label>
                  <div className="flex flex-wrap gap-2 text-sm">
                    {grantInitiators.map((username) => (
                      <span
                        key={username}
                        className="rounded bg-[#dfeaf4] px-2 py-1 text-primary-800 font-medium"
                      >
                        {username} [{grantUserFullnames[username] || ''}]
                      </span>
                    ))}
                  </div>
                </div>
              )}

              {grantAddApprovers && (
                <div className="rounded-xl border border-primary-100 bg-primary-100 p-4 text-center">
                  <label className="mb-2 block text-sm font-semibold text-primary-700">
                    No. of Approvals Required
                  </label>
                  <div className="text-base text-primary-800 font-semibold">
                    {grantNoOfApprovalsNeeded} out of {grantNoOfApprovers}
                  </div>
                </div>
              )}

              <div className="bg-white p-2">
                <label className="mb-1 block text-sm text-primary-600 font-medium px-1">
                  Enter Password to Authorize
                </label>
                <div className="flex items-center gap-3">
                  <input
                    type={showGrantPassword ? 'text' : 'password'}
                    value={grantPassword}
                    onChange={(e) => setGrantPassword(e.target.value)}
                    placeholder="Enter your wallet password"
                    className="h-12 flex-1 rounded-md border border-[#cfe0ee] bg-white px-3 text-primary-800 outline-none"
                  />
                </div>
                {grantPasswordErr && (
                  <p className="text-xs text-red-500 mt-1 px-1 font-semibold">
                    {grantPasswordErr}
                  </p>
                )}
              </div>
            </div>

            <Button
              label="Authorize"
              onclick={handleGrantSharedAccessSubmit}
              additionalClasses="mt-2"
            />
          </div>
        </div>
      );
    }

    if (activeModal === 'success') {
      return (
        <div className="p-6 md:p-8">
          <div className="space-y-5 text-center">
            <h3 className="text-[2rem] font-montserratSemiBold text-primary-800">
              Success!
            </h3>

            <div className="flex justify-center">
              <div className="flex h-28 w-28 items-center justify-center rounded-full bg-[#dfeaf4] text-5xl text-primary-800 font-bold">
                ✓
              </div>
            </div>

            <p className="text-xl font-medium text-primary-800">
              Shared access enabled successfully.
            </p>

            <p className="text-base text-primary-700">
              You have successfully enabled shared access on your wallet{' '}
              <span className="font-semibold text-primary-800">
                {grantWallet?.alias || grantWallet?.address}
              </span>
            </p>

            <Button
              label="Done"
              onclick={() => setActiveModal(null)}
              additionalClasses="mt-2 max-w-xs mx-auto"
            />
          </div>
        </div>
      );
    }

    if (activeModal === 'details') {
      return (
        <div className="p-6 md:p-8">
          <div className="space-y-5">
            <div className="flex items-center justify-between">
              <h3 className="text-center w-full text-xl font-montserratSemiBold text-primary-800">
                Shared Access
              </h3>
            </div>

            <div className="space-y-3 rounded-xl border border-[#d2e2f1] bg-[#f1f6fb] space-y-6 p-4 text-primary-800">
              <div className="flex flex-col items-center justify-between space-y-4">
                <span className="font-medium font-montserratSemiBold">
                  Wallet
                </span>
                <span className="text-primary-700">
                  {selectedRow?.wallet ?? '-'}
                </span>
              </div>
              <div className="flex flex-col items-center justify-between space-y-4">
                <span className="font-medium font-montserratSemiBold">
                  Description
                </span>
                <span className="text-primary-700">
                  {selectedRow?.description ?? '-'}
                </span>
              </div>
              <div className="flex flex-col items-center justify-between space-y-4">
                <span className="font-medium font-montserratSemiBold">
                  Owner
                </span>
                <span className="text-primary-700">
                  {selectedRow?.owner ?? '-'}
                </span>
              </div>
              <div className="flex flex-col items-center justify-between space-y-4">
                <span className="font-medium font-montserratSemiBold">
                  Permissions
                </span>
                <span className="text-primary-700">
                  {selectedRow && selectedRow.permissions.length > 0
                    ? `You have ${selectedRow.permissions.map((p) => p.toLowerCase()).join(' and ')} access on this wallet`
                    : 'No access permissions configured'}
                </span>
              </div>
            </div>

            <div className="space-y-2">
              <Button
                label="View Wallet"
                onclick={() => {
                  if (selectedRow) {
                    setActiveModal(null);
                    navigate(
                      `/dashboard/wallet?wallet=${selectedRow.address}&rel=shared`,
                    );
                  }
                }}
                additionalClasses="w-full hover:bg-primary-900"
              />
              <ButtonSecondary
                label="View Transaction History"
                onclick={() => {
                  setActiveModal(null);
                  if (selectedRow) {
                    navigate(
                      `/dashboard/history?wallet=${selectedRow.address}&rel=shared`,
                    );
                  } else {
                    navigate('/dashboard/history');
                  }
                }}
                additionalClasses="w-full bg-primary-600 text-white hover:bg-primary-700"
              />
              <button
                type="button"
                onClick={() => setActiveModal('modify')}
                className="w-full rounded-md border border-[#c7d8e7] bg-primary-200 px-4 py-3 text-primary-800 hover:bg-primary-300"
              >
                Modify Shared Access
              </button>
              <button
                type="button"
                onClick={() => {
                  if (selectedRow) {
                    setActiveModal('disableWarn');
                  }
                }}
                className="w-full rounded-md bg-red-500 px-4 py-3 text-white hover:bg-red-600"
              >
                Disable Shared Access
              </button>
            </div>
          </div>
        </div>
      );
    }

    if (activeModal === 'modify') {
      return (
        <div className="p-6 md:p-8">
          <div className="space-y-5">
            <div className="flex items-center justify-between">
              <h3 className="text-center w-full text-xl font-montserratSemiBold text-primary-800">
                Update Shared Access
              </h3>
            </div>

            <div className="rounded-xl  p-4">
              {modifyAddApprovers && (
                <div className="flex items-center justify-center w-full gap-4 border-b border-[#d7e5f3] pb-3 mb-4">
                  <button
                    type="button"
                    className="border-b-2 border-primary-700 px-2 py-1 text-primary-800 font-montserratSemiBold"
                  >
                    Viewers
                  </button>
                  <button
                    type="button"
                    onClick={() => setActiveModal('update')}
                    className="px-2 py-1 text-primary-500 hover:text-primary-800 font-montserratSemiBold transition"
                  >
                    Approvers
                  </button>
                  <button
                    type="button"
                    onClick={() => setActiveModal('initiators')}
                    className="px-2 py-1 text-primary-500 hover:text-primary-800 font-montserratSemiBold transition"
                  >
                    Initiators
                  </button>
                </div>
              )}

              <div className="space-y-3">
                <p className="text-sm text-primary-700">
                  Enter the usernames of all accounts that require viewer access
                  to this wallet
                </p>

                <div className="flex gap-2">
                  <input
                    type="text"
                    value={modifyViewerInput}
                    onChange={(e) => setModifyViewerInput(e.target.value)}
                    placeholder="Viewer username"
                    className="h-12 flex-1 rounded-md border border-[#cfe0ee] bg-white px-3 text-primary-800 outline-none"
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') {
                        e.preventDefault();
                        handleAddModifyViewer();
                      }
                    }}
                  />
                  <button
                    type="button"
                    onClick={handleAddModifyViewer}
                    className="h-12 rounded-md bg-primary-800 px-5 text-sm font-medium text-white hover:bg-primary-700 transition"
                  >
                    Add
                  </button>
                </div>
                {modifyViewerErr && (
                  <p className="text-xs text-red-500">{modifyViewerErr}</p>
                )}

                <div className="flex flex-wrap gap-2 rounded-lg border border-dashed border-[#b3d0e7] bg-primary-100 p-3 min-h-[3rem] items-center">
                  {modifyViewers.length > 0 ? (
                    modifyViewers.map((perm, idx) =>
                      renderPermissionChip(
                        perm,
                        idx,
                        modifyViewers,
                        setModifyViewers,
                        'Viewer',
                      ),
                    )
                  ) : (
                    <span className="text-sm text-primary-800">
                      Name of viewers appear here
                    </span>
                  )}
                </div>
              </div>

              {(() => {
                const walletObj = appUser?.userWallets.find(
                  (w) => w.address === selectedRow?.address,
                );
                if (
                  walletObj &&
                  !walletObj.isPrimaryWallet &&
                  !walletObj.primaryWallet &&
                  walletObj.walletType !== 2
                ) {
                  return (
                    <div className="mt-4">
                      <label className="flex items-center gap-2 text-sm text-primary-700 cursor-pointer">
                        <input
                          type="checkbox"
                          checked={modifyAddApprovers}
                          onChange={(e) =>
                            handleModifyAddApproversChange(e.target.checked)
                          }
                          className="h-4 w-4 rounded border-[#b5cfe4] text-primary-800 focus:ring-primary-800"
                        />
                        Grant approver access
                      </label>
                    </div>
                  );
                }
                return null;
              })()}
            </div>

            <Button
              label="Proceed"
              onclick={() => {
                if (modifyAddApprovers) {
                  setActiveModal('update');
                } else {
                  if (isModifyValid()) {
                    setActiveModal('confirmModify');
                  }
                }
              }}
              additionalClasses="mt-2"
            />
          </div>
        </div>
      );
    }

    if (activeModal === 'update') {
      return (
        <div className="p-6 md:p-8">
          <div className="space-y-5">
            <div className="flex items-center justify-between">
              <h3 className="text-center w-full text-xl font-montserratSemiBold text-primary-800">
                Update Shared Access
              </h3>
            </div>

            <div className="rounded-xl p-4">
              <div className="flex items-center justify-center w-full gap-4 border-b border-[#d7e5f3] pb-3 mb-4">
                <button
                  type="button"
                  onClick={() => setActiveModal('modify')}
                  className="px-2 py-1 text-primary-500 hover:text-primary-800 font-montserratSemiBold transition"
                >
                  Viewers
                </button>
                <button
                  type="button"
                  className="border-b-2 border-primary-700 px-2 py-1 text-primary-800 font-montserratSemiBold"
                >
                  Approvers
                </button>
                <button
                  type="button"
                  onClick={() => setActiveModal('initiators')}
                  className="px-2 py-1 text-primary-500 hover:text-primary-800 font-montserratSemiBold transition"
                >
                  Initiators
                </button>
              </div>

              <div className="space-y-4">
                <div>
                  <div className="flex items-center justify-center gap-4">
                    <div className="rounded-xl p-3">
                      <label className="mb-1 block font-medium text-primary-700">
                        Approvals
                      </label>
                      <select
                        className="h-10 w-full rounded-md border border-[#b5cfe4] bg-white px-2 text-sm text-primary-800 outline-none"
                        value={modifyNoOfApprovalsNeeded}
                        onChange={(e) =>
                          setModifyNoOfApprovalsNeeded(Number(e.target.value))
                        }
                      >
                        {Array.from(
                          { length: modifyNoOfApprovers - 1 },
                          (_, i) => i + 2,
                        ).map((val) => (
                          <option key={val} value={val}>
                            {val}
                          </option>
                        ))}
                      </select>
                    </div>
                    <label className="mb-1 block font-medium text-primary-700">
                      out of
                    </label>
                    <div className="rounded-xl p-3">
                      <label className="mb-1 block font-medium text-primary-700">
                        Approvers
                      </label>
                      <select
                        className="h-10 w-full rounded-md border border-[#b5cfe4] bg-white px-2 text-sm text-primary-800 outline-none"
                        value={modifyNoOfApprovers}
                        onChange={(e) => {
                          const val = Number(e.target.value);
                          setModifyNoOfApprovers(val);
                          if (modifyNoOfApprovalsNeeded >= val) {
                            setModifyNoOfApprovalsNeeded(val - 1);
                          }
                        }}
                      >
                        {[3, 4, 5, 6, 7, 8, 9, 10].map((val) => (
                          <option key={val} value={val}>
                            {val}
                          </option>
                        ))}
                      </select>
                    </div>
                  </div>
                  <div className="text-center text-xs text-primary-700 font-medium">
                    {modifyNoOfApprovalsNeeded} approvals required out of{' '}
                    {modifyNoOfApprovers} approvers
                  </div>
                </div>
                <div></div>

                <div className="space-y-3">
                  <label className="block text-sm font-medium text-primary-700">
                    Enter the usernames of all accounts that require approver
                    access to this wallet
                  </label>
                  <div className="flex gap-2">
                    <input
                      type="text"
                      value={modifyApproverInput}
                      onChange={(e) => setModifyApproverInput(e.target.value)}
                      placeholder="Approver"
                      className="h-12 flex-1 rounded-lg border border-[#cfe0ee] bg-white px-3 text-primary-800 outline-none"
                      onKeyDown={(e) => {
                        if (e.key === 'Enter') {
                          e.preventDefault();
                          handleAddModifyApprover();
                        }
                      }}
                    />
                    <button
                      type="button"
                      onClick={handleAddModifyApprover}
                      className="h-12 rounded-md bg-primary-800 px-5 text-sm font-medium text-white hover:bg-primary-700 transition"
                    >
                      Add
                    </button>
                  </div>
                  {modifyApproverErr && (
                    <p className="text-xs text-red-500">{modifyApproverErr}</p>
                  )}

                  <div className="flex flex-wrap gap-2 rounded-lg border border-dashed border-[#b3d0e7] bg-primary-100 p-3 min-h-[3rem] items-center">
                    {modifyApprovers.length > 0 ? (
                      modifyApprovers.map((perm, idx) =>
                        renderPermissionChip(
                          perm,
                          idx,
                          modifyApprovers,
                          setModifyApprovers,
                          'Approver',
                        ),
                      )
                    ) : (
                      <span className="text-sm text-primary-400">
                        Name of approvers appear here
                      </span>
                    )}
                  </div>
                </div>
              </div>
            </div>

            <Button
              label="Proceed"
              onclick={() => {
                const activeApprovers = modifyApprovers.filter(
                  (p) => p.permissionState !== PermissionState.Revoked,
                );
                if (activeApprovers.length === 0) {
                  showNotification('error', 'Please add approvers');
                  return;
                }
                if (activeApprovers.length < modifyNoOfApprovers) {
                  showNotification(
                    'error',
                    'Please add the required number of approvers',
                  );
                  return;
                }
                setActiveModal('initiators');
              }}
              additionalClasses="mt-2"
            />
          </div>
        </div>
      );
    }

    if (activeModal === 'initiators') {
      return (
        <div className="p-6 md:p-8">
          <div className="space-y-5">
            <div className="flex items-center justify-between">
              <h3 className="text-center w-full text-xl font-montserratSemiBold text-primary-800">
                Update Shared Access
              </h3>
            </div>

            <div className="rounded-xl p-4">
              <div className="flex items-center justify-center w-full gap-4 border-b border-[#d7e5f3] pb-3 mb-4">
                <button
                  type="button"
                  onClick={() => setActiveModal('modify')}
                  className="px-2 py-1 text-primary-500 hover:text-primary-800 font-montserratSemiBold transition"
                >
                  Viewers
                </button>
                <button
                  type="button"
                  onClick={() => setActiveModal('update')}
                  className="px-2 py-1 text-primary-500 hover:text-primary-800 font-montserratSemiBold transition"
                >
                  Approvers
                </button>
                <button
                  type="button"
                  className="border-b-2 border-primary-700 px-2 py-1 text-primary-800 font-montserratSemiBold"
                >
                  Initiators
                </button>
              </div>

              <div className="mt-4 space-y-3">
                <p className="text-sm text-primary-700">
                  Enter the usernames of all accounts that require initiator
                  access to this wallet
                </p>

                <div className="flex gap-2">
                  <input
                    type="text"
                    value={modifyInitiatorInput}
                    onChange={(e) => setModifyInitiatorInput(e.target.value)}
                    placeholder="Initiator username"
                    className="h-12 flex-1 rounded-md border border-[#cfe0ee] bg-white px-3 text-primary-800 outline-none"
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') {
                        e.preventDefault();
                        handleAddModifyInitiator();
                      }
                    }}
                  />
                  <button
                    type="button"
                    onClick={handleAddModifyInitiator}
                    className="h-12 rounded-md bg-primary-800 px-5 text-sm font-medium text-white hover:bg-primary-700 transition"
                  >
                    Add
                  </button>
                </div>
                {modifyInitiatorErr && (
                  <p className="text-xs text-red-500">{modifyInitiatorErr}</p>
                )}

                <div className="flex flex-wrap gap-2 rounded-lg border border-dashed border-[#b3d0e7] bg-primary-100 p-3 min-h-[3rem] items-center">
                  {modifyInitiators.length > 0 ? (
                    modifyInitiators.map((perm, idx) =>
                      renderPermissionChip(
                        perm,
                        idx,
                        modifyInitiators,
                        setModifyInitiators,
                        'Initiator',
                      ),
                    )
                  ) : (
                    <span className="text-sm text-primary-400">
                      Name of initiators appear here
                    </span>
                  )}
                </div>
              </div>
            </div>

            <Button
              label="Proceed"
              onclick={() => {
                if (isModifyValid()) {
                  setActiveModal('confirmModify');
                }
              }}
              additionalClasses="mt-2"
            />
          </div>
        </div>
      );
    }

    if (activeModal === 'confirmModify') {
      return (
        <div className="p-6 md:p-8">
          <div className="space-y-5">
            <div className="flex items-center justify-between">
              <h3 className="w-full text-center text-xl font-montserratSemiBold text-primary-800">
                Confirm Request
              </h3>
            </div>

            <p className="text-center text-sm text-primary-700 font-medium">
              You&rsquo;re about to make the following modifications to your
              shared access
            </p>

            <div className="space-y-4">
              <div className="rounded-xl border border-primary-100 bg-primary-100 p-4 text-center">
                <div className="mb-2 font-montserratSemiBold text-primary-700">
                  Wallet
                </div>
                <div className="text-base text-primary-800">
                  {selectedRow?.wallet || selectedRow?.address}
                </div>
              </div>

              {modifyViewers.length > 0 && (
                <div className="flex flex-col items-center space-y-6 rounded-xl border border-primary-100 bg-primary-100 p-4">
                  <div className="mb-3 text-center font-montserratSemiBold text-primary-700">
                    Viewer Access
                  </div>
                  <div className="space-y-2 text-sm w-2/3 text-primary-800">
                    {modifyViewers.map((v) => (
                      <div
                        key={v.targetUsername}
                        className="flex items-center justify-between font-medium"
                      >
                        <span>
                          {v.targetUsername} [{v.fullName}]
                        </span>
                        <span
                          className={
                            v.permissionState === PermissionState.Added
                              ? 'text-green-500 font-semibold'
                              : v.permissionState === PermissionState.Revoked
                                ? 'text-red-500 font-semibold'
                                : 'text-primary-500'
                          }
                        >
                          {v.permissionState || 'Active'}
                        </span>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {modifyAddApprovers && modifyApprovers.length > 0 && (
                <div className="flex flex-col items-center space-y-6 rounded-xl border border-primary-100 bg-primary-100 p-4">
                  <div className="mb-3 text-center font-montserratSemiBold text-primary-700">
                    Approver Access
                  </div>
                  <div className="space-y-2 text-sm w-2/3 text-primary-800">
                    {modifyApprovers.map((a) => (
                      <div
                        key={a.targetUsername}
                        className="flex items-center justify-between font-medium"
                      >
                        <span>
                          {a.targetUsername} [{a.fullName}]
                        </span>
                        <span
                          className={
                            a.permissionState === PermissionState.Added
                              ? 'text-green-500 font-semibold'
                              : a.permissionState === PermissionState.Revoked
                                ? 'text-red-500 font-semibold'
                                : 'text-primary-500'
                          }
                        >
                          {a.permissionState || 'Active'}
                        </span>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {modifyAddApprovers && modifyInitiators.length > 0 && (
                <div className="flex flex-col items-center space-y-6 rounded-xl border border-primary-100 bg-primary-100 p-4">
                  <div className="mb-3 text-center font-montserratSemiBold text-primary-700">
                    Initiator Access
                  </div>
                  <div className="space-y-2 text-sm w-2/3 text-primary-800">
                    {modifyInitiators.map((i) => (
                      <div
                        key={i.targetUsername}
                        className="flex items-center justify-between font-medium"
                      >
                        <span>
                          {i.targetUsername} [{i.fullName}]
                        </span>
                        <span
                          className={
                            i.permissionState === PermissionState.Added
                              ? 'text-green-500 font-semibold'
                              : i.permissionState === PermissionState.Revoked
                                ? 'text-red-500 font-semibold'
                                : 'text-primary-500'
                          }
                        >
                          {i.permissionState || 'Active'}
                        </span>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              <div className="rounded-xl border border-primary-100 bg-primary-100 p-4 text-center">
                <div className="mb-2 font-montserratSemiBold text-primary-700">
                  No. of Approvals Required
                </div>
                <div className="text-base text-primary-800 font-semibold">
                  {modifyAddApprovers
                    ? `${modifyNoOfApprovalsNeeded}/${modifyNoOfApprovers}`
                    : '0'}
                </div>
              </div>

              <div className="bg-white p-2">
                <label className="mb-1 block text-sm text-primary-600 font-medium px-1">
                  Enter Password to Authorize
                </label>
                <div className="flex items-center gap-3">
                  <input
                    type={showModifyPassword ? 'text' : 'password'}
                    value={modifyPassword}
                    onChange={(e) => setModifyPassword(e.target.value)}
                    placeholder="Enter your wallet password"
                    className="h-12 flex-1 rounded-md border border-[#cfe0ee] bg-white px-3 text-primary-800 outline-none"
                  />
                </div>
                {modifyPasswordErr && (
                  <p className="text-xs text-red-500 mt-1 px-1 font-semibold">
                    {modifyPasswordErr}
                  </p>
                )}
              </div>
            </div>

            <Button
              label="Authorize"
              onclick={handleModifySharedAccessSubmit}
              additionalClasses="mt-2"
            />
          </div>
        </div>
      );
    }

    if (activeModal === 'modifySuccess') {
      return (
        <div className="p-6 md:p-8">
          <div className="space-y-5 text-center">
            <h3 className="text-[2rem] font-montserratSemiBold text-primary-800">
              Success!
            </h3>

            <div className="flex justify-center">
              <div className="flex h-28 w-28 items-center justify-center rounded-full bg-[#dfeaf4] text-5xl text-primary-800 font-bold">
                ✓
              </div>
            </div>

            <div className="rounded-xl border border-[#dbe6f0] bg-[#edf4f9] px-6 py-5 text-center">
              <p className="text-xl font-medium text-primary-800">
                Request Successfully Submitted
              </p>
              <p className="mt-3 text-base text-primary-700">
                Your request to modify shared access on wallet{' '}
                <span className="font-semibold text-primary-800">
                  {selectedRow?.wallet || selectedRow?.address}
                </span>{' '}
                has been successfully submitted.
              </p>
              <p className="mt-2 text-base text-primary-700">
                This transaction will be completed when it gets the required
                number of approvals
              </p>
            </div>

            <Button
              label="Done"
              onclick={() => setActiveModal(null)}
              additionalClasses="mt-2 max-w-xs mx-auto"
            />
          </div>
        </div>
      );
    }

    if (activeModal === 'dateRange') {
      return (
        <div className={styles.modalBody}>
          <div className={styles.modalHeaderBlock}>
            <h2 className={styles.largeTitle}>Select Date Range</h2>
          </div>

          <div className={styles.dateSection}>
            <p className={styles.sectionTitle}>Quick Select</p>
            <div className={styles.presetRow}>
              <PresetButton
                label="Today"
                active={selectedDateRange === 'Today'}
                onClick={() => {
                  setSelectedDateRange('Today');
                  setCustomStartDate('');
                  setCustomEndDate('');
                  setDraftCustomStartDate('');
                  setDraftCustomEndDate('');
                  setActiveModal(null);
                }}
              />
              <PresetButton
                label="Last 7 days"
                active={selectedDateRange === 'Last 7 days'}
                onClick={() => {
                  setSelectedDateRange('Last 7 days');
                  setCustomStartDate('');
                  setCustomEndDate('');
                  setDraftCustomStartDate('');
                  setDraftCustomEndDate('');
                  setActiveModal(null);
                }}
              />
              <PresetButton
                label="Last 30 days"
                active={selectedDateRange === 'Last 30 days'}
                onClick={() => {
                  setSelectedDateRange('Last 30 days');
                  setCustomStartDate('');
                  setCustomEndDate('');
                  setDraftCustomStartDate('');
                  setDraftCustomEndDate('');
                  setActiveModal(null);
                }}
              />
            </div>

            <p className={styles.sectionTitle}>Custom Range</p>
            <div className={styles.dateInputs}>
              <DateInput
                label="Start Date"
                placeholder="From"
                value={draftCustomStartDate}
                onChange={setDraftCustomStartDate}
              />
              <DateInput
                label="End Date"
                placeholder="To"
                value={draftCustomEndDate}
                onChange={setDraftCustomEndDate}
              />
            </div>
          </div>

          <div className={styles.singleAction}>
            <button
              type="button"
              className={styles.primaryActionWide}
              onClick={() => {
                if (!draftCustomStartDate || !draftCustomEndDate) {
                  showNotification(
                    'error',
                    'Please select both start and end dates',
                  );
                  return;
                }
                setCustomStartDate(draftCustomStartDate);
                setCustomEndDate(draftCustomEndDate);
                setSelectedDateRange('Custom date');
                setActiveModal(null);
              }}
            >
              Apply
            </button>
          </div>
        </div>
      );
    }

    if (activeModal === 'moreFilters') {
      return (
        <div className={styles.modalBody}>
          <div className={styles.modalHeaderBlock}>
            <h2 className={styles.modalTitle}>Filter By:</h2>
            <p className={styles.modalSubtitle}>
              Select the options under the filters you wish to apply
            </p>
          </div>

          <div className={styles.formGrid}>
            <LabeledInput
              label="Initiator username"
              placeholder="Enter initiator username"
              value={draftInitiator}
              onChange={(value) => setDraftInitiator(value)}
            />
            <LabeledInput
              label="Wallet alias"
              placeholder="Select or enter wallet alias"
              value={draftWalletAlias}
              onChange={(value) => setDraftWalletAlias(value)}
            />
            <LabeledInput
              label="Wallet public key"
              placeholder="Select or enter wallet public key"
              value={draftWalletAddress}
              onChange={(value) => setDraftWalletAddress(value)}
            />
            <LabeledInput
              label="Description"
              placeholder="Select or enter description"
              value={draftDescription}
              onChange={(value) => setDraftDescription(value)}
            />
          </div>

          <div className={styles.modalActions}>
            <button
              type="button"
              className={styles.secondaryAction}
              onClick={() => {
                setDraftInitiator('');
                setDraftWalletAlias('');
                setDraftWalletAddress('');
                setDraftDescription('');
                setSelectedInitiator('');
                setSelectedWalletAlias('');
                setSelectedWalletAddress('');
                setSelectedDescription('');
                setActiveModal(null);
              }}
            >
              Clear All Filters
            </button>
            <button
              type="button"
              className={styles.primaryAction}
              onClick={() => {
                setSelectedInitiator(draftInitiator);
                setSelectedWalletAlias(draftWalletAlias);
                setSelectedWalletAddress(draftWalletAddress);
                setSelectedDescription(draftDescription);
                setActiveModal(null);
              }}
            >
              Apply
            </button>
          </div>
        </div>
      );
    }

    if (activeModal === 'disableWarn') {
      return (
        <div className="p-6 md:p-8">
          <div className="space-y-5 text-center">
            <h3 className="text-xl font-semibold text-primary-800">
              Important
            </h3>
            <p className="text-base text-red-600">
              Are you sure you want to disable shared access on this wallet?
            </p>
            <div className="space-y-1">
              <label className="block text-left text-sm font-medium text-primary-700">
                Password
              </label>
              <input
                type="password"
                placeholder="Enter password"
                value={disablePassword}
                onChange={(e) => setDisablePassword(e.target.value)}
                className="h-12 w-full rounded-md border border-[#cfe0ee] bg-white px-3 text-primary-800 outline-none"
              />
              {disablePasswordErr && (
                <p className="text-left text-sm text-red-500">
                  {disablePasswordErr}
                </p>
              )}
            </div>
            <div className="flex flex-col gap-2 pt-2">
              <Button
                label={isDisabling ? 'Processing...' : 'Disable'}
                onclick={handleDisableSharedAccessInit}
                additionalClasses="w-full bg-red-600 hover:bg-red-700"
              />
              <ButtonSecondary
                label="Cancel"
                onclick={() => {
                  setDisablePassword('');
                  setDisablePasswordErr('');
                  setActiveModal('details');
                }}
                additionalClasses="w-full"
              />
            </div>
          </div>
        </div>
      );
    }

    if (activeModal === 'disableConfirm') {
      return (
        <div className="p-6 md:p-8">
          <div className="space-y-5 text-center">
            <h3 className="text-xl font-semibold text-primary-800">
              Confirm Disable
            </h3>
            {disableTxData?.messages && disableTxData.messages.length > 0 && (
              <div className="rounded-xl border border-primary-100 bg-primary-100 p-4 text-left">
                {disableTxData.messages.map((msg: string, idx: number) => (
                  <p key={idx} className="text-sm text-red-800 mb-1">
                    {msg}
                  </p>
                ))}
              </div>
            )}
            <div className="flex flex-col gap-2 pt-4">
              <Button
                label={isDisabling ? 'Authorizing...' : 'Confirm'}
                onclick={handleDisableSharedAccessSubmit}
                additionalClasses="w-full"
              />
              <ButtonSecondary
                label="Cancel"
                onclick={() => {
                  setDisablePassword('');
                  setDisablePasswordErr('');
                  setActiveModal('details');
                }}
                additionalClasses="w-full"
              />
            </div>
          </div>
        </div>
      );
    }

    if (activeModal === 'disableSuccess') {
      return (
        <div className="p-6 md:p-8">
          <div className="space-y-5 text-center">
            <h3 className="text-[2rem] font-montserratSemiBold text-primary-800">
              Success!
            </h3>

            <div className="flex justify-center">
              <div className="flex h-28 w-28 items-center justify-center rounded-full bg-[#dfeaf4] text-5xl">
                ✓
              </div>
            </div>

            <div className="rounded-xl border border-[#dbe6f0] bg-[#edf4f9] px-6 py-5 text-center">
              <p className="text-xl font-medium text-primary-800">
                Request Completed
              </p>
              <p className="mt-3 text-base text-primary-700">
                {disableSuccessMessage}
              </p>
            </div>

            <Button
              label="Done"
              onclick={() => setActiveModal(null)}
              additionalClasses="mt-2 max-w-xs mx-auto"
            />
          </div>
        </div>
      );
    }

    return null;
  };

  const renderWelcomeView = () => (
    <div className="min-h-screen bg-white px-4 py-5 md:px-8">
      <div className="w-full rounded-lg bg-white p-10">
        <div className="mb-8 flex items-center justify-between gap-4">
          <h1 className="text-xl font-montserratSemiBold tracking-tight text-primary-800">
            Shared Access
          </h1>

          <button
            type="button"
            onClick={() => navigate('/dashboard/shared-access/add')}
            className="inline-flex items-center gap-2 rounded-xl bg-primary-800 px-5 py-3 text-sm font-medium text-white shadow-sm transition hover:bg-primary-700"
          >
            <span className="text-lg leading-none">+</span>
            <span>Grant Access</span>
          </button>
        </div>

        <div className="rounded-[2rem] bg-[#edf3f8] px-4 pb-6 pt-2">
          <div className="flex justify-center">
            <img src="/images/shared-access-rafiki.png" alt="illustration" />
          </div>
        </div>

        <div className="mt-6 rounded-lg bg-white px-10 py-6 text-primary-800 shadow-sm">
          <p className="mb-4 text-base">
            Here you can give others various access rights to your wallet.
          </p>

          <div className="space-y-4 leading-7">
            <p>
              <span className="font-montserratSemiBold">Viewer Access:</span>{' '}
              Enables other users to view your wallet balance, receive payment
              into your wallet and view your wallet history.
            </p>

            <p>
              <span className="font-montserratSemiBold">Initiator Access:</span>{' '}
              Enables other users in addition to viewer access, to initiate a
              transaction from your wallet and pass it to the appropriate
              approvers to approve.
            </p>

            <p>
              <span className="font-montserratSemiBold">Approver Access:</span>{' '}
              Enables other users in addition to viewer access, to become
              approvers on your wallet. This means that whenever a transaction
              is initiated from your wallet by the initiators, it must be
              approved by the required number of approvers out of the added
              approvers for the transaction to successfully go through.
            </p>

            <p>
              <span className="font-montserratSemiBold">Note:</span> You must be
              careful when giving approver access because once you give approver
              access on any of your wallets, the wallet seizes to be your sole
              wallet, it now becomes a jointly owned wallet that must get the
              approval of all the required approvers for any transaction to take
              place on it successfully. If you happen to be an initiator on the
              wallet with approver access, even though the wallet is originally
              your wallet, you will no longer be able to initiate transactions
              from the wallet.
            </p>

            <p>
              Also if you happen to be an Approver on the wallet with approver
              access enabled, you cannot approve any transaction on the wallet
              as well, you can only view the wallet going forward.
            </p>

            <p>
              Individuals can use Approver Access, however, it is best suited
              for organizations, businesses, associations, and any other use
              case where more than 1 person is required to operate a wallet.
            </p>
          </div>
        </div>
      </div>
    </div>
  );

  const renderSharedAccessView = () => (
    <div className="min-h-screen bg-white px-4 py-5 md:px-8">
      <div className="w-full rounded-lg bg-white-100 p-10">
        <div className="mb-8 flex items-center justify-between gap-4">
          <h1 className="text-xl font-montserratSemiBold tracking-tight text-primary-800">
            Shared Access
          </h1>

          <button
            type="button"
            onClick={() => setActiveModal('grant')}
            className="inline-flex items-center gap-2 rounded-xl bg-primary-800 px-5 py-3 text-sm font-medium text-white shadow-sm transition hover:bg-primary-700"
          >
            <span className="text-lg leading-none">+</span>
            <span>Grant Access</span>
          </button>
        </div>

        {/* <div className="mb-5 flex flex-wrap gap-3">
          <button
            type="button"
            onClick={() => setActiveModal('grant')}
            className="rounded-lg border border-primary-800 bg-white px-4 py-2 text-sm font-medium text-primary-800"
          >
            Grant Access Modal
          </button>
          <button
            type="button"
            onClick={() => setActiveModal('confirm')}
            className="rounded-lg border border-primary-800 bg-white px-4 py-2 text-sm font-medium text-primary-800"
          >
            Confirm Request
          </button>
          <button
            type="button"
            onClick={() => setActiveModal('success')}
            className="rounded-lg border border-primary-800 bg-white px-4 py-2 text-sm font-medium text-primary-800"
          >
            Success
          </button>
          <button
            type="button"
            onClick={() => setActiveModal('details')}
            className="rounded-lg border border-primary-800 bg-white px-4 py-2 text-sm font-medium text-primary-800"
          >
            Shared Access Details
          </button>
          <button
            type="button"
            onClick={() => setActiveModal('modify')}
            className="rounded-lg border border-primary-800 bg-white px-4 py-2 text-sm font-medium text-primary-800"
          >
            Modify Access
          </button>
          <button
            type="button"
            onClick={() => setActiveModal('update')}
            className="rounded-lg border border-primary-800 bg-white px-4 py-2 text-sm font-medium text-primary-800"
          >
            Update Access
          </button>
          <button
            type="button"
            onClick={() => setActiveModal('initiators')}
            className="rounded-lg border border-primary-800 bg-white px-4 py-2 text-sm font-medium text-primary-800"
          >
            Initiator Access
          </button>
          <button
            type="button"
            onClick={() => setActiveModal('confirmModify')}
            className="rounded-lg border border-primary-800 bg-white px-4 py-2 text-sm font-medium text-primary-800"
          >
            Confirm Modify
          </button>
          <button
            type="button"
            onClick={() => setActiveModal('modifySuccess')}
            className="rounded-lg border border-primary-800 bg-white px-4 py-2 text-sm font-medium text-primary-800"
          >
            Modify Success
          </button>
        </div> */}

        <div className="mb-2 flex items-center justify-between gap-3 border-b border-[#bdd0dd] px-2 pb-2 pt-2">
          <div className="flex-1">
            <Tabs tabList={['Shared Wallets', 'Access History']}>
              <div className="w-full">
                <div className="flex space-x-6">
                  <Dropdown
                    label={selectedAccessMode}
                    options={accessMode}
                    onSelect={(selectedItem) => {
                      setSelectedAccessMode(
                        selectedItem?.value ?? accessMode[0].value,
                      );
                    }}
                  />
                  <Dropdown
                    label={selectedFilter}
                    options={filter}
                    onSelect={(selectedItem) => {
                      setSelectedFilter(selectedItem?.value ?? filter[0].value);
                    }}
                  />
                </div>

                <div className="mt-4 overflow-hidden rounded-xl border border-[#d3dfe9] bg-white">
                  <CustomTable<SharedAccessRow>
                    className="p-0 bg-white spacing"
                    data={sharedWalletRows}
                    paginate
                    defaultItemsPerPage={20}
                    onRowClick={(row) => {
                      setSelectedRow(row);
                      setActiveModal('details');
                    }}
                    columns={[
                      {
                        key: 'wallet',
                        header: 'Wallet',
                        width: '30%',
                        className:
                          'text-[#295778] font-medium text-[1.02rem] bg-white',
                        render: (row) => <span>{row.wallet}</span>,
                      },
                      {
                        key: 'owner',
                        header: 'Owner',
                        width: '20%',
                        className: 'text-[#3a6787] text-[1.02rem] bg-white',
                      },
                      {
                        key: 'permissions',
                        header: 'Permissions',
                        width: '25%',
                        className:
                          'text-[#295778] font-medium text-[1.02rem] bg-white',
                        render: (row) => (
                          <div className="flex flex-wrap gap-2">
                            {row.permissions.map((permission) => (
                              <span
                                key={`${row.wallet}-${permission}`}
                                className="inline-flex rounded-full bg-[#d5e9f5] px-3 py-1 text-[0.9rem] font-medium text-[#2f5e85]"
                              >
                                {permission}
                              </span>
                            ))}
                          </div>
                        ),
                      },
                      {
                        key: 'description',
                        header: 'Description',
                        width: '25%',
                        className: 'text-[#3c6788] text-[1.02rem] bg-white',
                      },
                    ]}
                    emptyMessage="No shared wallets found"
                  />
                </div>
              </div>

              <div className="w-full">
                <div className={styles.filtersBar}>
                  <div className={styles.filterGrid}>
                    <HistorySelect
                      value={selectedTransactionType}
                      onChange={setSelectedTransactionType}
                      options={transactionTypeOptions.map((opt) => ({
                        label: opt.text,
                        value: opt.value,
                      }))}
                    />
                    <HistorySelect
                      value={selectedTransactionStatus}
                      onChange={setSelectedTransactionStatus}
                      options={transactionStatusOptions.map((opt) => ({
                        label: opt.text,
                        value: opt.value,
                      }))}
                    />
                    <DateFilterTrigger
                      label={dateRangeLabel}
                      icon={<CalendarIcon />}
                      customStyle={styles.DateFilterTrigger}
                      onClick={() => {
                        setDraftCustomStartDate(customStartDate);
                        setDraftCustomEndDate(customEndDate);
                        setActiveModal('dateRange');
                      }}
                    />
                    <FilterTrigger
                      label={`More Filters${[selectedInitiator, selectedWalletAlias, selectedWalletAddress, selectedDescription].filter(Boolean).length > 0 ? ` (${[selectedInitiator, selectedWalletAlias, selectedWalletAddress, selectedDescription].filter(Boolean).length})` : ''}`}
                      icon={<FilterIcon />}
                      onClick={() => {
                        setDraftInitiator(selectedInitiator);
                        setDraftWalletAlias(selectedWalletAlias);
                        setDraftWalletAddress(selectedWalletAddress);
                        setDraftDescription(selectedDescription);
                        setActiveModal('moreFilters');
                      }}
                    />
                    <ButtonSecondary
                      label="Reset Filters"
                      onclick={resetApprovalFilters}
                      additionalClasses={styles.resetFiltersButton}
                    />
                  </div>
                  <div className="flex space-x-2 justify-between w-full">
                    {activeFilters.length > 0 && (
                      <div className={styles.activeFiltersRow}>
                        <span className={styles.activeFiltersLabel}>
                          Active:
                        </span>
                        {activeFilters.map((chip) => (
                          <span key={chip.label} className={styles.filterChip}>
                            {chip.label}
                            <button
                              type="button"
                              className={styles.chipRemove}
                              onClick={chip.onRemove}
                              aria-label={`Remove ${chip.label} filter`}
                            >
                              ×
                            </button>
                          </span>
                        ))}
                      </div>
                    )}
                    <label className="flex items-center space-x-2 text-sm text-primary-800">
                      <input
                        type="checkbox"
                        checked={includeSignedTransactions}
                        onChange={(e) =>
                          setIncludeSignedTransactions(e.target.checked)
                        }
                        className="form-checkbox h-4 w-4 text-primary-600 rounded border-[#bdd0dd]"
                      />
                      <span>Include transaction I have already signed</span>
                    </label>
                  </div>
                </div>

                <div className="overflow-hidden rounded-xl border border-[#d3dfe9] bg-white">
                  <CustomTable<
                    SharedAccessActivityRow & {
                      _record: SharedAccessApprovalRecord;
                    }
                  >
                    className="p-0 bg-white spacing"
                    data={activityRows}
                    onRowClick={(row) => {
                      setSelectedActivityRecord(row._record);
                      setActiveModal('approvalDetails');
                      setApprovalPassword('');
                      setApprovalPasswordErr('');
                      setApprovalReason('');
                      setShowRejectReason(false);
                    }}
                    columns={[
                      {
                        key: 'wallet',
                        header: 'Wallet',
                        width: '18%',
                        className:
                          'text-[#295778] font-medium text-[1rem] bg-white',
                      },
                      {
                        key: 'transactionType',
                        header: 'Transaction Type',
                        width: '20%',
                        className: 'text-[#295778] text-[1rem] bg-white',
                      },
                      {
                        key: 'initiatedBy',
                        header: 'Initiated By',
                        width: '15%',
                        className: 'text-[#3a6787] text-[1rem] bg-white',
                      },
                      {
                        key: 'date',
                        header: 'Date',
                        width: '18%',
                        className: 'text-[#3a6787] text-[1rem] bg-white',
                      },
                      {
                        key: 'transactionStatus',
                        header: 'Transaction Status',
                        width: '15%',
                        className: 'text-[#3a6787] text-[1rem] bg-white',
                        render: (row) => {
                          const statusClasses = {
                            Pending: 'bg-[#dfeaf4] text-[#2f5e85]',
                            Completed: 'bg-[#dff5e8] text-[#2d7a58]',
                            Rejected: 'bg-[#f9e0d9] text-[#b35b42]',
                          };

                          return (
                            <span
                              className={`inline-flex rounded-full px-3 py-1 text-[0.9rem] font-medium ${statusClasses[row.transactionStatus]}`}
                            >
                              {row.transactionStatus}
                            </span>
                          );
                        },
                      },
                      {
                        key: 'approvalStatus',
                        header: 'Approval Status',
                        width: '14%',
                        className: 'text-[#3a6787] text-[1rem] bg-white',
                      },
                    ]}
                    emptyMessage="No activity history found"
                  />
                </div>
              </div>
            </Tabs>
          </div>
        </div>
      </div>
    </div>
  );

  return (
    <div className="min-h-screen bg-[#edf3f8]">
      <div className="flex justify-end px-4 pt-4 md:px-8">
        <button
          type="button"
          onClick={toggleView}
          className="inline-flex items-center rounded-xl border border-primary-800 bg-white px-4 py-2 text-sm font-medium text-primary-800 shadow-sm transition hover:bg-primary-50"
        >
          {showSharedAccessView ? 'Show existing view' : 'Show screenshot view'}
        </button>
      </div>
      {showSharedAccessView ? renderSharedAccessView() : renderWelcomeView()}
      {activeModal ? (
        <Modal
          showModal={Boolean(activeModal)}
          onClose={() => setActiveModal(null)}
        >
          {renderModalContent()}
        </Modal>
      ) : null}
    </div>
  );
}

const HistorySelect = ({
  value,
  onChange,
  options,
}: {
  value: string;
  onChange: (value: string) => void;
  options: Array<{ label: string; value: string }>;
}) => {
  return (
    <div className={styles.selectWrap}>
      <select
        className={styles.select}
        value={value}
        onChange={(event) => onChange(event.target.value)}
      >
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
      <span className={styles.selectIcon}>
        <ChevronDownIcon />
      </span>
    </div>
  );
};

const FilterTrigger = ({
  label,
  icon,
  onClick,
}: {
  label: string;
  icon: ReactNode;
  onClick: () => void;
}) => {
  return (
    <button type="button" className={styles.trigger} onClick={onClick}>
      <span className={styles.triggerIcon}>{icon}</span>
      <span className={styles.triggerLabel}>{label}</span>
    </button>
  );
};

const DateFilterTrigger = ({
  label,
  icon,
  onClick,
  customStyle,
}: {
  label: string;
  icon: ReactNode;
  onClick: () => void;
  customStyle: string;
}) => {
  return (
    <button
      type="button"
      className={`${styles.trigger} ${customStyle}`}
      onClick={onClick}
    >
      <span className={styles.triggerIcon}>{icon}</span>
      <span className={styles.triggerLabel}>{label}</span>
    </button>
  );
};

const LabeledInput = ({
  label,
  placeholder,
  value,
  onChange,
  type = 'text',
}: {
  label: string;
  placeholder: string;
  value: string;
  onChange: (value: string) => void;
  type?: string;
}) => {
  return (
    <label className={styles.inputGroup}>
      <span className={styles.inputLabel}>{label}</span>
      <input
        className={styles.input}
        type={type}
        placeholder={placeholder}
        value={value}
        onChange={(event) => onChange(event.target.value)}
      />
    </label>
  );
};

const DateInput = ({
  label,
  placeholder,
  value,
  onChange,
}: {
  label: string;
  placeholder: string;
  value: string;
  onChange: (value: string) => void;
}) => {
  return (
    <label className={styles.inputGroup}>
      <span className={styles.inputLabel}>{label}</span>
      <div className={styles.dateInputWrap}>
        <input
          className={`${styles.input}`}
          type="date"
          aria-label={label}
          placeholder={placeholder}
          value={value}
          onChange={(event) => onChange(event.target.value)}
        />
      </div>
    </label>
  );
};

const PresetButton = ({
  label,
  active,
  onClick,
}: {
  label: string;
  active: boolean;
  onClick: () => void;
}) => {
  return (
    <button
      type="button"
      className={`${styles.presetButton} ${active ? styles.presetButtonActive : ''}`}
      onClick={onClick}
    >
      {label}
    </button>
  );
};

const CalendarIcon = () => (
  <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
    <path
      d="M7 2V5M17 2V5M3 9H21M5 4H19C20.1046 4 21 4.89543 21 6V19C21 20.1046 20.1046 21 19 21H5C3.89543 21 3 20.1046 3 19V6C3 4.89543 3.89543 4 5 4Z"
      stroke="#336DA0"
      strokeWidth="2.4"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

const ChevronDownIcon = () => (
  <svg width="14" height="8" viewBox="0 0 14 8" fill="none">
    <path
      d="M1 1L7 7L13 1"
      stroke="#163C7A"
      strokeWidth="1.8"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

const FilterIcon = () => (
  <img
    src="/images/filterWhite.png"
    alt="filter"
    style={{ width: '16px', height: '16px' }}
  />
);
