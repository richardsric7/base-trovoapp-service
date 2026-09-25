"use client";

import React from "react";
import EmptyState from "@/components/EmptyState";
import { FaLock } from "react-icons/fa6";

// The backend gates every admin/vault-signer/* route behind AuthenticateSuperAdmin
// specifically (not "any authenticated admin") — this codebase doesn't currently persist a
// reliable "is this admin a SuperAdmin" flag on the client (the wallet-connect login response
// only carries adminLevel: 0 | 1, not the Role string other admins are listed with), so rather
// than fabricate a client-side role check against a field that may not correspond to reality,
// this surfaces the backend's own 403 plainly instead of a blank/broken table.
const AccessDenied = () => (
  <EmptyState
    icon={<FaLock color="#007CDF" size={28} />}
    title="Super Admin access required"
    message="Vault signer management is restricted to Trovo Super Admins. Contact a Super Admin if you need access."
  />
);

export default AccessDenied;
