import type { StaticImageData } from "next/image";
import Dashboard from "@/assets/images/dashboardIcon.svg";
import assetIcon from "@/assets/images/chainlink-(link).svg";
import teamIcon from "@/assets/images/buildings-2.svg";
import complianceIcon from "@/assets/images/complianceIcon.svg";
import reportIcon from "@/assets/images/reportsIcon.svg";
import revenueIcon from "@/assets/images/revenueIcon.svg";
import auditTrailIcon from "@/assets/images/AuditTrailIcon.svg";
import notificationIcon from "@/assets/images/notificationBellIcon.svg";
import vaultSignerIcon from "@/assets/images/key-square.svg";

export interface OrganizationNavigationItem {
  title: string;
  icon: StaticImageData;
  link: string;
  exact?: boolean;
}

const sharedNavigation: OrganizationNavigationItem[] = [
  {
    title: "Team Members",
    icon: teamIcon,
    link: "/organisation/teammembers",
  },
  {
    title: "Vault Signer",
    icon: vaultSignerIcon,
    link: "/organisation/vault-signer",
  },
  {
    title: "Compliance",
    icon: complianceIcon,
    link: "/organisation/compliance",
  },
  {
    title: "Audit Trail",
    icon: auditTrailIcon,
    link: "/organisation/audit-trail",
  },
  {
    title: "Notifications",
    icon: notificationIcon,
    link: "/organisation/notification",
  },
];

const defaultRoleNavigation: OrganizationNavigationItem[] = [
  {
    title: "Dashboard",
    icon: Dashboard,
    link: "/organisation",
    exact: true,
  },
  {
    title: "Tokenized Asset",
    icon: assetIcon,
    link: "/organisation/tokenizedasset",
  },
];

const roleNavigation: Record<string, OrganizationNavigationItem[]> = {
  rating_agency: [
    {
      title: "Dashboard",
      icon: Dashboard,
      link: "/organisation/ratingagency",
      exact: true,
    },
    {
      title: "Tokenized Asset",
      icon: assetIcon,
      link: "/organisation/ratingagency/tokenizedasset",
    },
  ],
  asset_issuing_house: [
    {
      title: "Dashboard",
      icon: Dashboard,
      link: "/organisation/issuinghouse",
      exact: true,
    },
    {
      title: "Tokenized Asset",
      icon: assetIcon,
      link: "/organisation/issuinghouse/tokenizedasset",
    },
  ],
  asset_manager: [
    {
      title: "Dashboard",
      icon: Dashboard,
      link: "/organisation/assetmanager",
      exact: true,
    },
    {
      title: "Tokenized Asset",
      icon: assetIcon,
      link: "/organisation/assetmanager/tokenizedasset",
    },
    {
      title: "Reports",
      icon: reportIcon,
      link: "/organisation/assetmanager/reports",
    },
    {
      title: "Revenue",
      icon: revenueIcon,
      link: "/organisation/assetmanager/revenue",
    },
    {
      title: "Fund Management",
      icon: assetIcon,
      link: "/organisation/assetmanager/fundmanagement",
    },
  ],
  approved_asset_custodian: [
    {
      title: "Dashboard",
      icon: Dashboard,
      link: "/organisation/assetcustodian",
      exact: true,
    },
    {
      title: "Tokenized Asset",
      icon: assetIcon,
      link: "/organisation/assetcustodian/tokenizedasset",
    },
    {
      title: "Fund Management",
      icon: assetIcon,
      link: "/organisation/assetcustodian/fundmanagement",
    },
  ],
  asset_custodian: [
    {
      title: "Dashboard",
      icon: Dashboard,
      link: "/organisation/assetcustodian",
      exact: true,
    },
    {
      title: "Tokenized Asset",
      icon: assetIcon,
      link: "/organisation/assetcustodian/tokenizedasset",
    },
    {
      title: "Fund Management",
      icon: assetIcon,
      link: "/organisation/assetcustodian/fundmanagement",
    },
  ],
  trustees: [
    {
      title: "Dashboard",
      icon: Dashboard,
      link: "/organisation/trustee",
      exact: true,
    },
    {
      title: "Tokenized Asset",
      icon: assetIcon,
      link: "/organisation/trustee/tokenizedasset",
    },
    {
      title: "Fund Management",
      icon: assetIcon,
      link: "/organisation/trustee/fundmanagement",
    },
    {
      title: "Distributions",
      icon: revenueIcon,
      link: "/organisation/trustee/distributions",
    },
  ],
  trustee: [
    {
      title: "Dashboard",
      icon: Dashboard,
      link: "/organisation/trustee",
      exact: true,
    },
    {
      title: "Tokenized Asset",
      icon: assetIcon,
      link: "/organisation/trustee/tokenizedasset",
    },
    {
      title: "Fund Management",
      icon: assetIcon,
      link: "/organisation/trustee/fundmanagement",
    },
    {
      title: "Distributions",
      icon: revenueIcon,
      link: "/organisation/trustee/distributions",
    },
  ],
  legal_adviser: [
    {
      title: "Dashboard",
      icon: Dashboard,
      link: "/organisation/legaladviser",
      exact: true,
    },
  ],
  financial_adviser: [
    {
      title: "Dashboard",
      icon: Dashboard,
      link: "/organisation/financialadviser",
      exact: true,
    },
  ],
};

export const getOrganizationNavigation = (
  stakeholderType?: string | null,
): OrganizationNavigationItem[] => {
  const normalizedType = stakeholderType?.trim().toLowerCase() ?? "";
  const roleItems = roleNavigation[normalizedType] ?? defaultRoleNavigation;

  return [...roleItems, ...sharedNavigation];
};
