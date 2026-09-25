"use client";

import Link from "next/link";
import React, { useState } from "react";
import styled from "styled-components";
import { FaAngleUp, FaArrowLeft } from "react-icons/fa6";
import { DatePicker } from "antd";
import type { Dayjs } from "dayjs";
import { DropdownSelect } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import AddFieldModal from "../components/AddFieldModal";
import AddSummary from "../components/AddSummary";

const { RangePicker } = DatePicker;

type DayRange = [Dayjs, Dayjs] | null;

type CostItem = { label: string; amount: string };

type FormShape = {
  asset_name: string;
  reporting_period: DayRange;
  report_title: string;

  revenue_gen: string;
  currency_name: string;
  revenue_notes: string;

  costs: CostItem[];
  cost_notes: string;
};

const formatMoneyInput = (raw: string) => {
  // allow digits + single dot
  const cleaned = raw.replace(/[^\d.]/g, "");
  const parts = cleaned.split(".");
  const intPart = parts[0] ?? "";
  const decPart = parts[1];

  const withCommas = intPart.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
  return decPart !== undefined ? `${withCommas}.${decPart}` : withCommas;
};

const parseMoney = (s: string): number | null => {
  const trimmed = s.trim();
  if (!trimmed) return null;
  const n = Number(trimmed.replace(/,/g, ""));
  return Number.isNaN(n) ? null : n;
};

const defaultCosts: CostItem[] = [
  { label: "Platform Fee", amount: "" },
  { label: "Transfer Fee", amount: "" },
  { label: "Asset Manager Fee", amount: "" },
  { label: "Custodian Fee", amount: "" },
  { label: "Regulatory Fee", amount: "" },
  { label: "Maintenance and Operations", amount: "" },
  { label: "Insurance Premium", amount: "" },
  { label: "Other Cost", amount: "" },
];

// FIX: currency symbol map
const CURRENCY_SYMBOLS = {
  CNGN: "₦",
  USD: "$",
  EUR: "€",
} as const;

const EditReportPage = () => {
  const [openAccordian, setOpenAccordian] = useState<string[]>([
    "Period & Asset Details",
    "Income Summary",
    "Cost & Breakdown Expense",
  ]);

  const [form, setForm] = useState<FormShape>({
    asset_name: "",
    reporting_period: null,
    report_title: "",

    revenue_gen: "",
    currency_name: "",
    revenue_notes: "",

    costs: defaultCosts,
    cost_notes: "",
  });
  const [isOpenAdd, setIsOpenAdd] = useState(false);
  const handleOpen = (section: string) => {
    setOpenAccordian((prev) =>
      prev.includes(section)
        ? prev.filter((s) => s !== section)
        : [...prev, section]
    );
  };

  const setField = <K extends keyof FormShape>(key: K, value: FormShape[K]) =>
    setForm((prev) => ({ ...prev, [key]: value }));

  const onSubmit = () => {
    const payload = {
      ...form,
      // FIX: safe ISO conversion without depending on dayjs plugins
      reporting_period: form.reporting_period
        ? [
            form.reporting_period[0].toDate().toISOString(),
            form.reporting_period[1].toDate().toISOString(),
          ]
        : null,
      revenue_gen: parseMoney(form.revenue_gen),
      costs: form.costs.map((c) => ({
        label: c.label.trim(),
        amount: parseMoney(c.amount),
      })),
    };
    console.log("Submit payload:", payload);
    // TODO: send to API
  };

  const updateCost = (idx: number, patch: Partial<CostItem>) =>
    setForm((p) => {
      const next = [...p.costs];
      next[idx] = { ...next[idx], ...patch };
      return { ...p, costs: next };
    });

  const currencySymbol =
    CURRENCY_SYMBOLS[form.currency_name as keyof typeof CURRENCY_SYMBOLS] ??
    "₦";
  return (
    <EditIncomeContainer>
      <BackButtonLink href="/assetincomereport">
        <FaArrowLeft />
      </BackButtonLink>

      <ReportFormCard>
        <Title>Asset Income Report</Title>

        {/* Period & Asset Details */}
        <Section>
          <AccordianHeader onClick={() => handleOpen("Period & Asset Details")}>
            <SectionTitle>Period & Asset Details</SectionTitle>
            <Angle $open={openAccordian.includes("Period & Asset Details")}>
              <FaAngleUp />
            </Angle>
          </AccordianHeader>

          {openAccordian.includes("Period & Asset Details") && (
            <>
              <FormGroup>
                <Label>Asset Name</Label>
                <TextInput
                  placeholder="Atlantis"
                  value={form.asset_name}
                  onChange={(e) => setField("asset_name", e.target.value)}
                />
              </FormGroup>

              <FormGroup>
                <Label>Reporting Period</Label>
                <RangePicker
                  style={{ width: "100%" }}
                  value={form.reporting_period} // FIX: now [Dayjs, Dayjs] | null
                  onChange={(range) =>
                    setField("reporting_period", range as DayRange)
                  }
                />
              </FormGroup>

              <FormGroup>
                <Label>Report Title</Label>
                <TextInput
                  placeholder="Enter"
                  value={form.report_title}
                  onChange={(e) => setField("report_title", e.target.value)}
                />
              </FormGroup>
            </>
          )}
        </Section>

        {/* Income Summary */}
        <Section>
          <AccordianHeader onClick={() => handleOpen("Income Summary")}>
            <SectionTitle>Income Summary</SectionTitle>
            <Angle $open={openAccordian.includes("Income Summary")}>
              <FaAngleUp />
            </Angle>
          </AccordianHeader>

          {openAccordian.includes("Income Summary") && (
            <>
              <FormGroup>
                <Label>Total Revenue Generated</Label>
                <FeeInputWrapper>
                  <DividerGroup>
                    <FeeLabel>{currencySymbol}</FeeLabel>
                    <VerticalDivider />
                  </DividerGroup>
                  <FeeInputField
                    placeholder="100,000,000.00"
                    value={form.revenue_gen}
                    onChange={(e) =>
                      setField("revenue_gen", formatMoneyInput(e.target.value))
                    }
                  />
                </FeeInputWrapper>
              </FormGroup>

              <FormGroup>
                <Label>Currency</Label>
                <DropdownSelect
                  options={["CNGN", "USD", "EUR"]}
                  placeholder="Select"
                  labelText=""
                  labelColor="#828282"
                  placeholderColor="#00225A"
                  value={form.currency_name}
                  onSelect={(item: string) => setField("currency_name", item)}
                />
              </FormGroup>

              <FormGroup>
                <Label>Revenue Notes</Label>
                <TextAreaInput
                  rows={3}
                  value={form.revenue_notes}
                  onChange={(e) => setField("revenue_notes", e.target.value)}
                />
              </FormGroup>
            </>
          )}
        </Section>

        {/* Cost & Breakdown Expense */}
        <Section>
          <AccordianHeader
            onClick={() => handleOpen("Cost & Breakdown Expense")}
          >
            <SectionTitle>Cost & Breakdown Expense</SectionTitle>
            <Angle $open={openAccordian.includes("Cost & Breakdown Expense")}>
              <FaAngleUp />
            </Angle>
          </AccordianHeader>

          {openAccordian.includes("Cost & Breakdown Expense") && (
            <>
              <AddFieldBar>
                <AddButton type="button" onClick={() => setIsOpenAdd(true)}>
                  + Add Field
                </AddButton>
              </AddFieldBar>

              {form.costs.map((c, idx) => (
                <CostRow key={`${c.label}-${idx}`}>
                  <CostGrid>
                    <Label> {c.label}</Label>

                    <FeeInputWrapper>
                      <DividerGroup>
                        <FeeLabel>{currencySymbol}</FeeLabel>
                        <VerticalDivider />
                      </DividerGroup>
                      <FeeInputField
                        placeholder="0"
                        value={c.amount}
                        onChange={(e) =>
                          updateCost(idx, {
                            amount: formatMoneyInput(e.target.value),
                          })
                        }
                      />
                    </FeeInputWrapper>
                  </CostGrid>
                </CostRow>
              ))}

              <FormGroup>
                <Label>Cost Notes</Label>
                <TextAreaInput
                  rows={3}
                  value={form.cost_notes}
                  onChange={(e) => setField("cost_notes", e.target.value)}
                />
              </FormGroup>
            </>
          )}
        </Section>

        <AddSummary />

        <PrimaryButton onClick={onSubmit} buttonStyle={{ width: "100%" }}>
          Submmit
        </PrimaryButton>
      </ReportFormCard>

      {isOpenAdd && (
        <AddFieldModal isOpen={isOpenAdd} setIsOpenAdd={setIsOpenAdd} />
      )}
    </EditIncomeContainer>
  );
};

export default EditReportPage;

const EditIncomeContainer = styled.section`
  background: #ffffff;
  padding: 32px;
  gap: 20px;
  border-radius: 24px;
`;

const BackButtonLink = styled(Link)`
  text-decoration: none;
  display: block;
  color: #000000;
  width: 20px;
`;

const ReportFormCard = styled.div`
  width: 100%;
  max-width: 560px;
  margin: auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
`;

const Title = styled.h2`
  color: #00225a;
  text-align: center;
  font-weight: 700;
  font-size: 24px;
  line-height: 28px;
  margin-bottom: 8px;
`;

const Section = styled.div`
  background: #fff;
`;

const AccordianHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
`;

const SectionTitle = styled.p`
  color: #00225a;
  font-weight: 600;
  font-size: 16px;
  line-height: 24px;
  padding: 10px 0;
  margin: 0;
`;

const Angle = styled.span<{ $open: boolean }>`
  display: inline-flex;
  transition: transform 0.2s ease;
  transform: rotate(${(p) => (p.$open ? 0 : 180)}deg);
`;

const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Label = styled.label`
  color: #191919;
  font-weight: 500;
  font-size: 16px;
  margin-top: 14px;
`;

const VerticalDivider = styled.div`
  width: 1px;
  height: 24px;
  background-color: #e0e0e0;
`;

const FeeLabel = styled.p`
  font-weight: 500;
  font-size: 14px;
  color: #191919;
`;

const DividerGroup = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;

const FeeInputField = styled.input`
  font-size: 14px;
  width: 100%;
  outline: none;
  font-family: inherit;
  border: none;

  &::placeholder {
    color: #bdbdbd;
    font-family: inherit;
    font-size: 16px;
    padding-left: 8px;
  }
`;

const FeeInputWrapper = styled.div`
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  display: flex;
  align-items: center;
  width: 100%;
  margin-top: 4px;
`;

const TextInput = styled.input`
  padding: 10px 12px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  outline: none;
  font-size: 14px;
  font-family: inherit;

  &::placeholder {
    font-family: inherit;
    color: #bdbdbd;
    font-size: 16px;
  }
`;

const TextAreaInput = styled.textarea<{ rows?: number }>`
  padding: 10px 12px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  outline: none;
  font-size: 14px;
  resize: vertical;
`;

const AddFieldBar = styled.div`
  margin-bottom: 8px;
`;

const CostRow = styled.div`
  margin-bottom: 8px;
`;

const CostGrid = styled.div`
  display: flex;
  flex-direction: column;
  gap: 10px;
`;

const AddButton = styled.button`
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #007cdf;
  background: transparent;
  border-radius: 8px;
  cursor: pointer;
  color: #007cdf;
  font-family: inherit;
  font-weight: 600;
  font-size: 14px;
`;
