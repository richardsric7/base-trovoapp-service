"use client";
import React, { useEffect, useState } from "react";
import { showErrorToast, showSuccessToast } from "@/components";
import {
  useGetPayoutEngineQuery,
  useGetPayoutFeeConfigQuery,
  usePayoutEngineActionMutation,
  useSetPayoutFeeConfigMutation,
} from "@/redux/api/proceedPayouts";
import { Button, Card, date, errorMessage, Grid, Input, Mono, Muted, Note, Row, Select, Stat, SubTitle } from "./ui";

// The payout fee wallet and default fee, and payout-engine's state with its
// kill switch and sweep.
const FeeAndEngine = () => {
  const { data: cfgData } = useGetPayoutFeeConfigQuery();
  const { data: engineData } = useGetPayoutEngineQuery(undefined, { pollingInterval: 15000 });
  const [saveConfig, { isLoading: saving }] = useSetPayoutFeeConfigMutation();
  const [engineAction, { isLoading: acting }] = usePayoutEngineActionMutation();
  const cfg = cfgData?.data;
  const engine = engineData?.data;

  const [feeWallet, setFeeWallet] = useState("");
  const [feeType, setFeeType] = useState<"FIXED" | "PERCENT">("FIXED");
  const [feeValue, setFeeValue] = useState("0");
  const [feeCap, setFeeCap] = useState("0");
  const [sweepToken, setSweepToken] = useState("");

  useEffect(() => {
    if (!cfg) return;
    setFeeWallet(cfg.feeWallet || "");
    setFeeType(cfg.feeType || "FIXED");
    setFeeValue(cfg.feeValue || "0");
    setFeeCap(cfg.feeCap || "0");
  }, [cfg]);

  const save = async () => {
    try {
      await saveConfig({ feeWallet: feeWallet.trim(), feeType, feeValue: feeValue.trim() || "0", feeCap: feeType === "PERCENT" ? feeCap.trim() || "0" : "0" }).unwrap();
      showSuccessToast("Payout fee configuration saved");
    } catch (err) {
      showErrorToast(errorMessage(err, "Could not save the configuration"));
    }
  };

  const run = async (action: "halt" | "unhalt" | "sweep") => {
    let reason: string | undefined;
    if (action === "halt") {
      reason = window.prompt("Why stop all payouts? (recorded in the audit trail)") ?? undefined;
      if (!reason?.trim()) return;
    }
    if (action === "unhalt" && !window.confirm("Let payout-engine work again?")) return;
    if (action === "sweep" && !window.confirm(`Move the payout Safe's whole balance of ${sweepToken} to the sweep address?`)) return;
    try {
      await engineAction({ action, reason, token: sweepToken.trim() }).unwrap();
      showSuccessToast(action === "halt" ? "Payout engine stopped" : action === "unhalt" ? "Payout engine resumed" : "Sweep requested");
    } catch (err) {
      showErrorToast(errorMessage(err, "The request failed"));
    }
  };

  return (
    <>
      <Card>
        <SubTitle>Payout processing fee</SubTitle>
        <Muted>
          Each payout carries its own fee, set on the payout and approved with its schedule; new payouts start from
          the defaults below. FIXED is an amount of the payout token; PERCENT is a share of the payout, capped by the
          cap when it is above 0. VAT at the asset country&apos;s rate is charged on the fee. The fee and VAT come out
          of the payout and are paid to the fee wallet and the VAT wallet with the holders.
        </Muted>
        <Row>
          <Input $wide placeholder="Fee wallet address (0x…)" value={feeWallet} onChange={(e) => setFeeWallet(e.target.value)} />
          <Select value={feeType} onChange={(e) => setFeeType(e.target.value as "FIXED" | "PERCENT")}>
            <option value="FIXED">Fixed fee</option>
            <option value="PERCENT">Percent fee</option>
          </Select>
          <Input placeholder={feeType === "PERCENT" ? "Percent" : "Amount"} value={feeValue} onChange={(e) => setFeeValue(e.target.value)} />
          {feeType === "PERCENT" && <Input placeholder="Cap (0 = none)" value={feeCap} onChange={(e) => setFeeCap(e.target.value)} />}
          <Button $variant="primary" disabled={saving} onClick={save}>
            {saving ? "Saving…" : "Save"}
          </Button>
        </Row>
        <Grid>
          <Stat label="VAT wallet" value={cfg?.vatWallet ? <Mono>{cfg.vatWallet}</Mono> : "VAT_WALLET (engine environment)"} />
          <Stat label="Last changed by" value={cfg?.lastUpdatedBy || "—"} />
        </Grid>
      </Card>

      <Card>
        <SubTitle>Payout engine</SubTitle>
        {engine ? (
          <>
            <Grid>
              <Stat label="Status" value={engine.halted ? "Stopped (kill switch)" : engine.online ? "Running" : "Not responding"} />
              <Stat label="Last heartbeat" value={date(engine.heartbeatAt)} />
              <Stat label="Instance" value={engine.instance || "—"} />
              <Stat label="Version" value={engine.version || "—"} />
              <Stat label="Doing" value={engine.activity || "—"} />
            </Grid>
            {engine.halted && (
              <Note $tone="warn">
                Stopped by {engine.haltedBy} on {date(engine.haltedAt)}: {engine.haltReason}
              </Note>
            )}
            {engine.lastError && <Note $tone="warn">Last error: {engine.lastError}</Note>}
            {engine.sweepToken && <Note>Sweep of {engine.sweepToken} requested by {engine.sweepRequestedBy}; waiting for the engine.</Note>}
            {!engine.sweepToken && engine.sweepResult && <Note>Last sweep: {engine.sweepResult}</Note>}
          </>
        ) : (
          <Muted>The engine has not reported yet.</Muted>
        )}
        <Row>
          {engine?.halted ? (
            <Button disabled={acting} onClick={() => run("unhalt")}>
              Resume engine
            </Button>
          ) : (
            <Button $variant="danger" disabled={acting} onClick={() => run("halt")}>
              Stop all payouts
            </Button>
          )}
        </Row>
        <Row>
          <Input $wide placeholder="Token contract to sweep from the payout Safe (0x…)" value={sweepToken} onChange={(e) => setSweepToken(e.target.value)} />
          <Button disabled={acting || !sweepToken.trim()} onClick={() => run("sweep")}>
            Sweep
          </Button>
        </Row>
      </Card>
    </>
  );
};

export default FeeAndEngine;
