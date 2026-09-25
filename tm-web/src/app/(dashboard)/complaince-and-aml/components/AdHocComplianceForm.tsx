"use client";

import { useMemo } from "react";
import styled from "styled-components";
import { useFormik } from "formik";
import * as yup from "yup";

import { DropdownSelect, showErrorToast, showSuccessToast } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import { useGetOrganizationsListQuery } from "@/redux/api/organizations";
import { useCreateComplianceRequirementMutation } from "@/redux/api/admin";
import { ComplianceInputType } from "@/redux/api/compliance/interface";

const INPUT_TYPE_OPTIONS: { value: string; label: string }[] = [
  { value: ComplianceInputType.DocumentUpload, label: "Document Upload" },
  { value: ComplianceInputType.Text, label: "Text" },
  { value: ComplianceInputType.StructuredForm, label: "Structured Form" },
];

interface ComplianceFormValues {
  orgId: string;
  category: string;
  requirement: string;
  description: string;
  inputType: string;
  required: boolean;
  dueDate: string;
}

const initialValues: ComplianceFormValues = {
  orgId: "",
  category: "",
  requirement: "",
  description: "",
  inputType: "",
  required: true,
  dueDate: "",
};

const complianceValidationSchema = yup.object({
  orgId: yup.string().required("Organisation is required"),
  category: yup.string().trim().required("Category is required"),
  requirement: yup.string().trim().required("Requirement is required"),
  inputType: yup.string().required("Input type is required"),
  dueDate: yup.string(),
});

// Ad-hoc assignment: an admin assigns a one-off requirement to any
// organisation outside the standard template, but it flows into the same
// review queue as templated items (PRD §6.2) via the shared instance table.
const AdHocComplianceForm = ({ onCreated }: { onCreated?: () => void }) => {
  const { data, isLoading: isLoadingOrgs } = useGetOrganizationsListQuery({
    page: 1,
    pageSize: 100,
  });
  const [createRequirement, { isLoading: isSubmitting }] =
    useCreateComplianceRequirementMutation();

  const organisations = useMemo(() => data?.organizations ?? [], [data]);

  const { values, errors, touched, handleChange, setFieldValue, handleSubmit } =
    useFormik<ComplianceFormValues>({
      initialValues,
      validationSchema: complianceValidationSchema,
      onSubmit: async (formValues, { resetForm: reset }) => {
        try {
          await createRequirement({
            org_id: formValues.orgId,
            category: formValues.category.trim(),
            requirement: formValues.requirement.trim(),
            description: formValues.description.trim() || undefined,
            input_type: formValues.inputType,
            required: formValues.required,
            ...(formValues.dueDate ? { due_date: formValues.dueDate } : {}),
          }).unwrap();

          showSuccessToast("Compliance requirement created successfully.");
          reset();
          onCreated?.();
        } catch (err: any) {
          showErrorToast(
            err?.data?.message || err?.error || "Failed to create requirement.",
          );
        }
      },
    });

  const selectedOrg = useMemo(
    () => organisations.find((org) => org.id === values.orgId),
    [organisations, values.orgId],
  );

  return (
    <Form onSubmit={handleSubmit}>
      <FormGrid>
        <FormGroup>
          <DropdownSelect
            labelText="Organisation *"
            placeholder={
              isLoadingOrgs ? "Loading organisations..." : "Select organisation"
            }
            value={selectedOrg?.name || ""}
            options={organisations.map((org) => org.name)}
            onSelect={(name) => {
              const found = organisations.find((org) => org.name === name);
              setFieldValue("orgId", found?.id || "");
            }}
            errorMessage={touched.orgId ? errors.orgId : undefined}
          />
        </FormGroup>

        <FormGroup>
          <Label>Category *</Label>
          <Input
            name="category"
            placeholder="e.g. KYC, AML Screening, Custody Audit"
            value={values.category}
            onChange={handleChange}
          />
          {touched.category && errors.category && (
            <ErrorMessage>{errors.category}</ErrorMessage>
          )}
        </FormGroup>

        <FormGroup>
          <DropdownSelect
            labelText="Input Type *"
            placeholder="Select input type"
            value={
              INPUT_TYPE_OPTIONS.find((opt) => opt.value === values.inputType)
                ?.label || ""
            }
            options={INPUT_TYPE_OPTIONS.map((opt) => opt.label)}
            onSelect={(label) => {
              const found = INPUT_TYPE_OPTIONS.find(
                (opt) => opt.label === label,
              );
              setFieldValue("inputType", found?.value || "");
            }}
            errorMessage={touched.inputType ? errors.inputType : undefined}
          />
        </FormGroup>

        <FormGroup>
          <Label>Due Date</Label>
          <Input
            type="date"
            name="dueDate"
            value={values.dueDate}
            onChange={handleChange}
          />
        </FormGroup>

        <FormGroup>
          <CheckboxLabel>
            <input
              type="checkbox"
              checked={values.required}
              onChange={(e) => setFieldValue("required", e.target.checked)}
            />
            Required
          </CheckboxLabel>
        </FormGroup>

        <FormGroup style={{ gridColumn: "1 / -1" }}>
          <Label>Requirement *</Label>
          <TextArea
            name="requirement"
            placeholder="Describe what the organisation needs to do or provide"
            value={values.requirement}
            onChange={handleChange}
            rows={4}
          />
          {touched.requirement && errors.requirement && (
            <ErrorMessage>{errors.requirement}</ErrorMessage>
          )}
        </FormGroup>

        <FormGroup style={{ gridColumn: "1 / -1" }}>
          <Label>Description / Instructions</Label>
          <TextArea
            name="description"
            placeholder="Shown to the organisation at submission time"
            value={values.description}
            onChange={handleChange}
            rows={2}
          />
        </FormGroup>
      </FormGrid>

      <PrimaryButton
        disabled={isSubmitting}
        buttonStyle={{ margin: "24px 0 0 0", width: "260px" }}
      >
        {isSubmitting ? "Creating..." : "Create Requirement"}
      </PrimaryButton>
    </Form>
  );
};

export default AdHocComplianceForm;

const Form = styled.form``;

const ErrorMessage = styled.p`
  font-size: 12px;
  color: #be3800;
  margin: 6px 0 0 0;
`;

const FormGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
`;

const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
`;

const Label = styled.label`
  font-weight: 500;
  font-size: 14px;
  line-height: 24px;
  margin-bottom: 8px;
  color: #00225a;
`;

const CheckboxLabel = styled.label`
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #00225a;
  margin-top: auto;
  height: 44px;
`;

const Input = styled.input`
  height: 44px;
  border-radius: 10px;
  border: 1px solid #e0e0e0;
  padding: 0 12px;
  font-size: 14px;
  color: #00225a;
  font-family: inherit;
  outline: none;

  &:focus {
    border-color: #007cdf;
  }
`;

const TextArea = styled.textarea`
  border-radius: 10px;
  border: 1px solid #e0e0e0;
  padding: 12px;
  font-size: 14px;
  color: #00225a;
  font-family: inherit;
  resize: vertical;
  outline: none;

  &:focus {
    border-color: #007cdf;
  }
`;
