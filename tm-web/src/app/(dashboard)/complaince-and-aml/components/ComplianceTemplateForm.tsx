"use client";

import styled from "styled-components";
import { FieldArray, FormikProvider, useFormik } from "formik";
import * as yup from "yup";

import { DropdownSelect, showErrorToast, showSuccessToast } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import { ORG_TYPES } from "@/constants/organizationTypes";
import {
  ComplianceInputType,
  ComplianceTemplate,
} from "@/redux/api/compliance/interface";
import {
  useCreateComplianceTemplateMutation,
  useUpdateComplianceTemplateMutation,
} from "@/redux/api/admin";

const INPUT_TYPE_OPTIONS: { value: string; label: string }[] = [
  { value: ComplianceInputType.DocumentUpload, label: "Document Upload" },
  { value: ComplianceInputType.Text, label: "Text" },
  { value: ComplianceInputType.StructuredForm, label: "Structured Form" },
];

const LEVEL_OPTIONS = ["1", "2", "3"];

interface TemplateItemFormValue {
  name: string;
  category: string;
  description: string;
  input_type: string;
  required: boolean;
}

interface TemplateFormValues {
  org_type: string[];
  level: string;
  items: TemplateItemFormValue[];
}

const emptyItem: TemplateItemFormValue = {
  name: "",
  category: "",
  description: "",
  input_type: "",
  required: true,
};

const validationSchema = yup.object({
  org_type: yup.array().of(yup.string()).min(1, "Select at least one org type"),
  level: yup.string().required("Level is required"),
  items: yup
    .array()
    .of(
      yup.object({
        name: yup.string().trim().required("Name is required"),
        category: yup.string().trim().required("Category is required"),
        input_type: yup.string().required("Input type is required"),
      }),
    )
    .min(1, "Add at least one requirement item"),
});

export interface ComplianceTemplateFormProps {
  editingTemplate?: ComplianceTemplate | null;
  onSaved?: () => void;
}

const ComplianceTemplateForm = ({
  editingTemplate,
  onSaved,
}: ComplianceTemplateFormProps) => {
  const [createTemplate, { isLoading: isCreating }] =
    useCreateComplianceTemplateMutation();
  const [updateTemplate, { isLoading: isUpdating }] =
    useUpdateComplianceTemplateMutation();

  const isEditing = !!editingTemplate;
  const isSubmitting = isCreating || isUpdating;

  const initialValues: TemplateFormValues = editingTemplate
    ? {
        org_type: [editingTemplate.org_type],
        level: String(editingTemplate.level),
        items: editingTemplate.items.map((item) => ({
          name: item.name,
          category: item.category,
          description: item.description || "",
          input_type: item.input_type,
          required: item.required,
        })),
      }
    : {
        org_type: [],
        level: "",
        items: [{ ...emptyItem }],
      };

  const formik = useFormik<TemplateFormValues>({
    initialValues,
    enableReinitialize: true,
    validationSchema,
    onSubmit: async (values, { resetForm: reset }) => {
      const items = values.items.map((item) => ({
        name: item.name.trim(),
        category: item.category.trim(),
        description: item.description.trim() || undefined,
        input_type: item.input_type,
        required: item.required,
      }));

      try {
        if (isEditing && editingTemplate) {
          await updateTemplate({
            id: editingTemplate.id,
            payload: {
              org_type: values.org_type[0],
              level: Number(values.level),
              items,
            },
          }).unwrap();
          showSuccessToast(
            "Template updated. Existing organisation requirements already assigned are unaffected.",
          );
        } else {
          // The backend stores one org_type per template row, so a
          // multi-select creates one template per selected org type.
          const results = await Promise.allSettled(
            values.org_type.map((orgType) =>
              createTemplate({
                org_type: orgType,
                level: Number(values.level),
                items,
              }).unwrap(),
            ),
          );
          const failures = results.filter((r) => r.status === "rejected");
          if (failures.length === 0) {
            showSuccessToast("Compliance template created successfully.");
            reset();
          } else if (failures.length < results.length) {
            showErrorToast(
              `Created for ${results.length - failures.length} of ${results.length} org types — some failed (they may already have a template at this level).`,
            );
          } else {
            showErrorToast(
              "Failed to create template — an org type may already have a template at this level.",
            );
          }
        }
        onSaved?.();
      } catch (err: any) {
        showErrorToast(
          err?.data?.message || err?.error || "Failed to save template.",
        );
      }
    },
  });

  const { values, errors, touched, handleChange, setFieldValue, handleSubmit } =
    formik;

  const toggleOrgType = (value: string) => {
    const next = values.org_type.includes(value)
      ? values.org_type.filter((v) => v !== value)
      : [...values.org_type, value];
    setFieldValue("org_type", next);
  };

  return (
    <FormikProvider value={formik}>
      <Form onSubmit={handleSubmit}>
        {isEditing && (
          <Notice>
            Saving increments this template&apos;s version. Organisation
            requirements already assigned keep the version they were created
            under and are not changed retroactively.
          </Notice>
        )}

        <FormGroup>
          <Label>Applicable Org Types *</Label>
          <CheckboxGrid>
            {ORG_TYPES.map((orgType) => (
              <CheckboxLabel key={orgType.value}>
                <input
                  type="checkbox"
                  checked={values.org_type.includes(orgType.value)}
                  onChange={() => toggleOrgType(orgType.value)}
                />
                {orgType.label}
              </CheckboxLabel>
            ))}
          </CheckboxGrid>
          {touched.org_type && errors.org_type && (
            <ErrorMessage>{errors.org_type as string}</ErrorMessage>
          )}
        </FormGroup>

        <FormGroup>
          <DropdownSelect
            labelText="Level *"
            placeholder="Select level"
            value={values.level}
            options={LEVEL_OPTIONS}
            onSelect={(item) => setFieldValue("level", item)}
            errorMessage={touched.level ? (errors.level as string) : undefined}
          />
        </FormGroup>

        <FieldArray name="items">
          {({ push, remove }) => (
            <ItemsSection>
              <ItemsHeader>
                <Label style={{ margin: 0 }}>Requirement Items *</Label>
                <AddItemButton
                  type="button"
                  onClick={() => push({ ...emptyItem })}
                >
                  + Add item
                </AddItemButton>
              </ItemsHeader>

              {typeof errors.items === "string" && (
                <ErrorMessage>{errors.items}</ErrorMessage>
              )}

              {values.items.map((item, index) => {
                const itemErrors = (errors.items?.[index] || {}) as Partial<
                  Record<keyof TemplateItemFormValue, string>
                >;
                const itemTouched = touched.items?.[index] as
                  | Partial<Record<keyof TemplateItemFormValue, boolean>>
                  | undefined;

                return (
                  <ItemRow key={index}>
                    <ItemGrid>
                      <FormGroup>
                        <Label>Name *</Label>
                        <Input
                          name={`items.${index}.name`}
                          placeholder="e.g. Certificate of Incorporation"
                          value={item.name}
                          onChange={handleChange}
                        />
                        {itemTouched?.name && itemErrors.name && (
                          <ErrorMessage>{itemErrors.name}</ErrorMessage>
                        )}
                      </FormGroup>

                      <FormGroup>
                        <Label>Category *</Label>
                        <Input
                          name={`items.${index}.category`}
                          placeholder="e.g. KYB"
                          value={item.category}
                          onChange={handleChange}
                        />
                        {itemTouched?.category && itemErrors.category && (
                          <ErrorMessage>{itemErrors.category}</ErrorMessage>
                        )}
                      </FormGroup>

                      <FormGroup>
                        <DropdownSelect
                          labelText="Input Type *"
                          placeholder="Select input type"
                          value={
                            INPUT_TYPE_OPTIONS.find(
                              (opt) => opt.value === item.input_type,
                            )?.label || ""
                          }
                          options={INPUT_TYPE_OPTIONS.map((opt) => opt.label)}
                          onSelect={(label) => {
                            const found = INPUT_TYPE_OPTIONS.find(
                              (opt) => opt.label === label,
                            );
                            setFieldValue(
                              `items.${index}.input_type`,
                              found?.value || "",
                            );
                          }}
                          errorMessage={
                            itemTouched?.input_type
                              ? itemErrors.input_type
                              : undefined
                          }
                        />
                      </FormGroup>

                      <FormGroup>
                        <CheckboxLabel>
                          <input
                            type="checkbox"
                            checked={item.required}
                            onChange={(e) =>
                              setFieldValue(
                                `items.${index}.required`,
                                e.target.checked,
                              )
                            }
                          />
                          Required
                        </CheckboxLabel>
                      </FormGroup>

                      <FormGroup style={{ gridColumn: "1 / -1" }}>
                        <Label>Description / Instructions</Label>
                        <TextArea
                          name={`items.${index}.description`}
                          placeholder="Shown to the organisation at submission time"
                          value={item.description}
                          onChange={handleChange}
                          rows={2}
                        />
                      </FormGroup>
                    </ItemGrid>

                    {values.items.length > 1 && (
                      <RemoveItemButton
                        type="button"
                        onClick={() => remove(index)}
                      >
                        Remove item
                      </RemoveItemButton>
                    )}
                  </ItemRow>
                );
              })}
            </ItemsSection>
          )}
        </FieldArray>

        <PrimaryButton
          disabled={isSubmitting}
          buttonStyle={{ margin: "24px 0 0 0", width: "260px" }}
        >
          {isSubmitting
            ? "Saving..."
            : isEditing
              ? "Save Template"
              : "Create Template"}
        </PrimaryButton>
      </Form>
    </FormikProvider>
  );
};

export default ComplianceTemplateForm;

const Form = styled.form``;

const Notice = styled.p`
  background: #fff4e0;
  color: #b25e00;
  font-size: 13px;
  padding: 12px 16px;
  border-radius: 10px;
  margin: 0 0 20px 0;
`;

const ErrorMessage = styled.p`
  font-size: 12px;
  color: #be3800;
  margin: 6px 0 0 0;
`;

const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  margin-bottom: 20px;
`;

const Label = styled.label`
  font-weight: 500;
  font-size: 14px;
  line-height: 24px;
  margin-bottom: 8px;
  color: #00225a;
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

const CheckboxGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
`;

const CheckboxLabel = styled.label`
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #00225a;
`;

const ItemsSection = styled.div`
  margin-bottom: 20px;
`;

const ItemsHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
`;

const AddItemButton = styled.button`
  border: 1px dashed #007cdf;
  background: transparent;
  color: #007cdf;
  border-radius: 8px;
  padding: 6px 12px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
`;

const ItemRow = styled.div`
  border: 1px solid #f2f2f2;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 12px;
`;

const ItemGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
`;

const RemoveItemButton = styled.button`
  border: none;
  background: transparent;
  color: #eb5757;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  margin-top: 8px;
  padding: 0;
`;
