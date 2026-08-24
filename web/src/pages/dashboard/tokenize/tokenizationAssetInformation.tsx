import { useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Button from '../../../components/button';
import ButtonSecondary from '../../../components/buttonSecondary';
import TextArea from '../../../components/textArea';
import TextInput from '../../../components/textInput';
import {
  hideLoader,
  showLoader,
  showNotification,
} from '../../../utils/showToaster';
import { useSelector } from 'react-redux';
import { RootState } from '../../../store/reduxStore';
import { Encryptor } from '../../../utils/encryptor';
import {
  TokenizedAsset,
  excludedSubmissionFields,
} from '../../../types/tokenizedAsset';
import {
  useLazyFetchFormFieldsQuery,
  useSubmitTokenizationMutation,
} from '../../../store/api/tokenizationApis';
import { useActiveTokenizedAsset } from './useActiveTokenizedAsset';
import { useTokenizationData } from './useTokenizationData';

type FormCondition = {
  type?: string;
  section?: string | number;
  targetName?: string;
  value?: string | number | boolean | null;
};

type FormField = {
  widgetType?: string;
  type?: string;
  label?: string;
  placeholderText?: string;
  required?: boolean;
  value?: unknown;
  defaultValue?: unknown;
  options?:
    | Array<{ text: string; value: string | number | boolean }>
    | string[];
  when?: FormCondition;
  maxWords?: number;
};

type FormSection = {
  title?: string;
  description?: string;
  body: Record<string, FormField>;
};

type FormDefinition = Record<string, FormSection>;

type FormPayload = Record<string, unknown>;
type TokenizedAssetFormPayload = Partial<TokenizedAsset>;

type BehaviorEvent = 'onInit' | 'onChange';

type BehaviorCondition = {
  op?: string;
  args?: BehaviorCondition[];
  field?: string;
  value?: unknown;
};

type BehaviorValueExpression = {
  op?: string;
  args?: Array<BehaviorValueExpression | number | string | boolean | null>;
  value?: BehaviorValueExpression | number | string | boolean | null;
  field?: string;
  default?: unknown;
  digits?: number;
};

type BehaviorAction = {
  type?: string;
  field?: string;
  value?: unknown;
  valueExpr?: BehaviorValueExpression;
};

type FormBehavior = {
  id?: string;
  event?: BehaviorEvent;
  watch?: string[];
  when?: BehaviorCondition;
  actions?: BehaviorAction[];
};

const fallbackFormDefinition: FormDefinition = {
  'Asset Information': {
    title: 'Asset Information',
    description:
      'Provide the basic information about the asset you want to tokenize.',
    body: {
      assetName: {
        widgetType: 'text',
        label: 'Asset Name',
        placeholderText: 'Enter the asset name',
        required: true,
      },
      assetDescription: {
        widgetType: 'longtext',
        label: 'Asset Description',
        placeholderText: 'Describe the asset',
        required: true,
      },
      assetLocation: {
        widgetType: 'dropdown',
        label: 'Asset Location',
        placeholderText: 'Select the country',
        options: [
          { text: 'United States', value: 'US' },
          { text: 'United Kingdom', value: 'UK' },
          { text: 'Nigeria', value: 'NG' },
        ],
        required: true,
      },
    },
  },
  'Asset Status': {
    title: 'Asset Status',
    description: 'Select the current status that applies to the asset.',
    body: {
      assetStatus: {
        widgetType: 'radio',
        label: 'Asset Status',
        options: [
          { text: 'Asset is already existing', value: 'existing' },
          { text: 'Asset is yet to be built/acquired', value: 'toBuild' },
        ],
        required: true,
      },
    },
  },
  'Funding Structure': {
    title: 'Funding Structure',
    description: 'Choose the funding structure that applies to the asset.',
    body: {
      fundingStructure: {
        widgetType: 'radio',
        label: 'Funding Structure',
        options: [
          { text: 'Equity', value: 'equity' },
          { text: 'Debt', value: 'debt' },
          { text: 'Hybrid', value: 'hybrid' },
        ],
        required: true,
      },
    },
  },
  'Offering Details': {
    title: 'Offering Details',
    description:
      'Choose the offering type and confirm the document requirements.',
    body: {
      offeringType: {
        widgetType: 'radio',
        label: 'Offering Type',
        options: [
          { text: 'Public', value: 'public' },
          { text: 'Private', value: 'private' },
        ],
        required: true,
      },
      hasDocuments: {
        widgetType: 'radio',
        label: 'Do you have the required documents?',
        options: [
          { text: 'Yes', value: 'yes' },
          { text: 'No', value: 'no' },
        ],
        required: true,
      },
      confirmCustodian: {
        widgetType: 'checkbox',
        label:
          'I confirm that this asset will be transferred to a licensed Asset Custodian.',
        required: true,
      },
    },
  },
};

function parseApiFormPayload(payload: unknown): unknown {
  if (typeof payload === 'string') {
    const trimmedPayload = payload.trim();

    if (!trimmedPayload) {
      return null;
    }

    try {
      return parseApiFormPayload(JSON.parse(trimmedPayload));
    } catch {
      return payload;
    }
  }

  if (payload && typeof payload === 'object' && !Array.isArray(payload)) {
    const source = payload as Record<string, unknown>;

    if (source.formString) {
      return parseApiFormPayload(source.formString);
    }

    if (
      source.data &&
      typeof source.data === 'object' &&
      !Array.isArray(source.data)
    ) {
      return parseApiFormPayload(
        (source.data as Record<string, unknown>).formString,
      );
    }
  }

  return payload;
}

function normalizeFormDefinition(payload: unknown): FormDefinition {
  const parsedPayload = parseApiFormPayload(payload);

  if (
    !parsedPayload ||
    typeof parsedPayload !== 'object' ||
    Array.isArray(parsedPayload)
  ) {
    return fallbackFormDefinition;
  }

  const source = parsedPayload as Record<string, unknown>;

  if (
    source.fields &&
    typeof source.fields === 'object' &&
    !Array.isArray(source.fields)
  ) {
    return Object.entries(
      source.fields as Record<string, unknown>,
    ).reduce<FormDefinition>((accumulator, [sectionKey, sectionValue]) => {
      if (
        !sectionValue ||
        typeof sectionValue !== 'object' ||
        Array.isArray(sectionValue)
      ) {
        return accumulator;
      }

      const section = sectionValue as Record<string, unknown>;
      const sectionName = String(
        section.sectionName ?? section.title ?? sectionKey,
      );
      const bodySource =
        section.body &&
        typeof section.body === 'object' &&
        !Array.isArray(section.body)
          ? (section.body as Record<string, unknown>)
          : {};

      const body = Object.entries(bodySource).reduce<Record<string, FormField>>(
        (bodyAccumulator, [fieldName, fieldValue]) => {
          if (
            !fieldValue ||
            typeof fieldValue !== 'object' ||
            Array.isArray(fieldValue)
          ) {
            bodyAccumulator[fieldName] = { value: '' };
            return bodyAccumulator;
          }

          const fieldEntry = fieldValue as Record<string, unknown>;
          const optionsValue = fieldEntry.options;
          const normalizedOptions = Array.isArray(optionsValue)
            ? optionsValue.map((option) => {
                if (typeof option === 'string') {
                  return option;
                }

                if (option && typeof option === 'object') {
                  const optionEntry = option as Record<string, unknown>;
                  return {
                    text: String(optionEntry.text ?? ''),
                    value: optionEntry.value ?? '',
                  };
                }

                return { text: String(option), value: String(option) };
              })
            : [];

          bodyAccumulator[fieldName] = {
            widgetType: String(fieldEntry.widgetType ?? ''),
            type: String(fieldEntry.type ?? ''),
            label: fieldEntry.label ? String(fieldEntry.label) : fieldName,
            placeholderText: fieldEntry.placeholderText
              ? String(fieldEntry.placeholderText)
              : undefined,
            required: Boolean(fieldEntry.required),
            value: fieldEntry.value ?? undefined,
            defaultValue:
              fieldEntry.defaultValue ?? fieldEntry.value ?? undefined,
            options: normalizedOptions as FormField['options'],
            when:
              fieldEntry.when && typeof fieldEntry.when === 'object'
                ? {
                    type: String(
                      (fieldEntry.when as Record<string, unknown>).type ?? '',
                    ),
                    section: (fieldEntry.when as Record<string, unknown>)
                      .section as string | number | undefined,
                    targetName: (fieldEntry.when as Record<string, unknown>)
                      .targetName
                      ? String(
                          (fieldEntry.when as Record<string, unknown>)
                            .targetName,
                        )
                      : undefined,
                    value: (fieldEntry.when as Record<string, unknown>)
                      .value as string | number | boolean | null | undefined,
                  }
                : undefined,
            maxWords: fieldEntry.maxWords
              ? Number(fieldEntry.maxWords)
              : undefined,
          };

          return bodyAccumulator;
        },
        {},
      );

      accumulator[sectionKey] = {
        title: String(sectionName),
        description: String(section.description ?? ''),
        body,
      };

      return accumulator;
    }, {});
  }

  if (source.sections && Array.isArray(source.sections)) {
    return source.sections.reduce<FormDefinition>(
      (accumulator, sectionEntry) => {
        if (!sectionEntry || typeof sectionEntry !== 'object') {
          return accumulator;
        }

        const section = sectionEntry as Record<string, unknown>;
        const sectionName = String(section.name ?? section.title ?? 'Section');
        const body =
          section.body &&
          typeof section.body === 'object' &&
          !Array.isArray(section.body)
            ? (section.body as Record<string, FormField>)
            : {};

        accumulator[sectionName] = {
          title: String(section.title ?? sectionName),
          description: String(section.description ?? ''),
          body,
        };

        return accumulator;
      },
      {},
    );
  }

  return Object.entries(source).reduce<FormDefinition>(
    (accumulator, [sectionName, sectionValue]) => {
      if (
        !sectionValue ||
        typeof sectionValue !== 'object' ||
        Array.isArray(sectionValue)
      ) {
        return accumulator;
      }

      const section = sectionValue as Record<string, unknown>;
      const body =
        section.body &&
        typeof section.body === 'object' &&
        !Array.isArray(section.body)
          ? (section.body as Record<string, FormField>)
          : {};

      accumulator[sectionName] = {
        title: String(section.title ?? sectionName),
        description: String(section.description ?? ''),
        body,
      };

      return accumulator;
    },
    {},
  );
}

function normalizeFormBehaviors(payload: unknown): FormBehavior[] {
  const parsedPayload = parseApiFormPayload(payload);

  if (
    !parsedPayload ||
    typeof parsedPayload !== 'object' ||
    Array.isArray(parsedPayload)
  ) {
    return [];
  }

  const source = parsedPayload as Record<string, unknown>;
  const behaviors = source.behaviors;

  if (!Array.isArray(behaviors)) {
    return [];
  }

  return behaviors
    .filter((entry) => entry && typeof entry === 'object')
    .map((entry) => {
      const behavior = entry as Record<string, unknown>;
      const watch = Array.isArray(behavior.watch)
        ? behavior.watch
            .filter((item) => typeof item === 'string')
            .map((item) => String(item))
        : [];

      return {
        id: behavior.id ? String(behavior.id) : undefined,
        event:
          behavior.event === 'onInit' || behavior.event === 'onChange'
            ? behavior.event
            : undefined,
        watch,
        when:
          behavior.when && typeof behavior.when === 'object'
            ? (behavior.when as BehaviorCondition)
            : undefined,
        actions: Array.isArray(behavior.actions)
          ? (behavior.actions as BehaviorAction[])
          : [],
      };
    });
}

function getFieldValueByName(
  definition: FormDefinition,
  fieldName: string,
): unknown {
  for (const section of Object.values(definition)) {
    if (fieldName in section.body) {
      return section.body[fieldName]?.value;
    }
  }

  return undefined;
}

function isValuePresent(value: unknown): boolean {
  if (value === null || value === undefined) {
    return false;
  }

  if (typeof value === 'string') {
    return value.trim().length > 0;
  }

  if (Array.isArray(value)) {
    return value.length > 0;
  }

  return true;
}

function parseBehaviorNumericValue(value: unknown): number {
  if (typeof value === 'number') {
    return value;
  }

  if (typeof value === 'string') {
    const normalized = value.replace(/,/g, '').trim();
    if (!normalized) {
      return Number.NaN;
    }

    return Number(normalized);
  }

  return Number(value);
}

function evaluateBehaviorCondition(
  condition: BehaviorCondition | undefined,
  definition: FormDefinition,
): boolean {
  if (!condition || !condition.op) {
    return true;
  }

  switch (condition.op) {
    case 'or':
      return (condition.args ?? []).some((item) =>
        evaluateBehaviorCondition(item, definition),
      );
    case 'and':
      return (condition.args ?? []).every((item) =>
        evaluateBehaviorCondition(item, definition),
      );
    case 'notEmpty':
      if (!condition.field) {
        return false;
      }
      return isValuePresent(getFieldValueByName(definition, condition.field));
    case 'isEmpty':
      if (!condition.field) {
        return false;
      }
      return !isValuePresent(getFieldValueByName(definition, condition.field));
    case 'eq':
      if (!condition.field) {
        return false;
      }
      return (
        String(getFieldValueByName(definition, condition.field) ?? '') ===
        String(condition.value ?? '')
      );
    case 'gt':
      if (!condition.field) {
        return false;
      }
      return (
        parseBehaviorNumericValue(
          getFieldValueByName(definition, condition.field),
        ) > parseBehaviorNumericValue(condition.value)
      );
    case 'gte':
      if (!condition.field) {
        return false;
      }
      return (
        parseBehaviorNumericValue(
          getFieldValueByName(definition, condition.field),
        ) >= parseBehaviorNumericValue(condition.value)
      );
    case 'lt':
      if (!condition.field) {
        return false;
      }
      return (
        parseBehaviorNumericValue(
          getFieldValueByName(definition, condition.field),
        ) < parseBehaviorNumericValue(condition.value)
      );
    case 'lte':
      if (!condition.field) {
        return false;
      }
      return (
        parseBehaviorNumericValue(
          getFieldValueByName(definition, condition.field),
        ) <= parseBehaviorNumericValue(condition.value)
      );
    default:
      return false;
  }
}

function evaluateBehaviorValueExpression(
  expression: BehaviorValueExpression | number | string | boolean | null,
  definition: FormDefinition,
): unknown {
  if (
    expression === null ||
    typeof expression === 'number' ||
    typeof expression === 'string' ||
    typeof expression === 'boolean'
  ) {
    return expression;
  }

  if (!expression || typeof expression !== 'object' || !expression.op) {
    return null;
  }

  switch (expression.op) {
    case 'toNumber': {
      const raw = expression.field
        ? getFieldValueByName(definition, expression.field)
        : expression.default;
      const parsed = parseBehaviorNumericValue(raw ?? expression.default ?? 0);
      return Number.isNaN(parsed) ? Number(expression.default ?? 0) : parsed;
    }
    case 'mul': {
      const values = (expression.args ?? []).map((item) =>
        Number(evaluateBehaviorValueExpression(item, definition) ?? 0),
      );
      return values.reduce((accumulator, current) => accumulator * current, 1);
    }
    case 'div': {
      const left = Number(
        evaluateBehaviorValueExpression(
          expression.args?.[0] ?? 0,
          definition,
        ) ?? 0,
      );
      const right = Number(
        evaluateBehaviorValueExpression(
          expression.args?.[1] ?? 1,
          definition,
        ) ?? 1,
      );

      if (right === 0) {
        return 0;
      }

      return left / right;
    }
    case 'round': {
      const digits = Number(expression.digits ?? 0);
      const base = Number(
        evaluateBehaviorValueExpression(expression.value ?? 0, definition) ?? 0,
      );

      if (Number.isNaN(base)) {
        return 0;
      }

      const factor = 10 ** (Number.isNaN(digits) ? 0 : digits);
      return Math.round(base * factor) / factor;
    }
    default:
      return null;
  }
}

function setFieldValueByName(
  definition: FormDefinition,
  fieldName: string,
  value: unknown,
): { definition: FormDefinition; changed: boolean } {
  for (const [sectionName, section] of Object.entries(definition)) {
    if (!(fieldName in section.body)) {
      continue;
    }

    const targetField = section.body[fieldName];
    const currentValue = targetField?.value;
    const nextValue = normalizeExistingFieldValue(targetField, value);

    if (currentValue === nextValue) {
      return { definition, changed: false };
    }

    return {
      definition: {
        ...definition,
        [sectionName]: {
          ...section,
          body: {
            ...section.body,
            [fieldName]: {
              ...targetField,
              value: nextValue,
            },
          },
        },
      },
      changed: true,
    };
  }

  return { definition, changed: false };
}

function applyBehaviorEngine(
  definition: FormDefinition,
  behaviors: FormBehavior[],
  event: BehaviorEvent,
  changedFieldNames: string[] = [],
): FormDefinition {
  if (!behaviors.length) {
    return definition;
  }

  let nextDefinition = definition;
  const triggerFields = new Set(changedFieldNames);

  for (const behavior of behaviors) {
    if (behavior.event !== event) {
      continue;
    }

    if (
      event === 'onChange' &&
      behavior.watch &&
      behavior.watch.length > 0 &&
      !behavior.watch.some((watchedField) => triggerFields.has(watchedField))
    ) {
      continue;
    }

    if (!evaluateBehaviorCondition(behavior.when, nextDefinition)) {
      continue;
    }

    for (const action of behavior.actions ?? []) {
      if (action.type !== 'setValue' || !action.field) {
        continue;
      }

      const resolvedValue =
        action.valueExpr !== undefined
          ? evaluateBehaviorValueExpression(action.valueExpr, nextDefinition)
          : action.value;

      const result = setFieldValueByName(
        nextDefinition,
        action.field,
        resolvedValue,
      );

      if (result.changed) {
        nextDefinition = result.definition;
      }
    }
  }

  return nextDefinition;
}

function seedFormValues(definition: FormDefinition): FormDefinition {
  return Object.entries(definition).reduce<FormDefinition>(
    (accumulator, [sectionName, section]) => {
      accumulator[sectionName] = {
        ...section,
        body: Object.entries(section.body).reduce<Record<string, FormField>>(
          (bodyAccumulator, [fieldName, field]) => {
            let initialValue: unknown = field.value ?? field.defaultValue;

            if (initialValue === undefined || initialValue === null) {
              if (field.widgetType === 'checkbox') {
                initialValue = false;
              } else if (field.widgetType === 'list') {
                initialValue = '';
              } else if (
                field.widgetType === 'dropdown' ||
                field.widgetType === 'radio'
              ) {
                initialValue = '';
              } else {
                initialValue = '';
              }
            }

            bodyAccumulator[fieldName] = {
              ...field,
              value: initialValue,
            };
            return bodyAccumulator;
          },
          {},
        ),
      };

      return accumulator;
    },
    {},
  );
}

function normalizeExistingFieldValue(
  field: FormField,
  value: unknown,
): unknown {
  if (value === null || value === undefined) {
    return value;
  }

  if (field.widgetType === 'checkbox') {
    if (typeof value === 'boolean') {
      return value;
    }

    if (typeof value === 'number') {
      return value === 1;
    }

    if (typeof value === 'string') {
      const normalizedValue = value.trim().toLowerCase();
      return (
        normalizedValue === '1' ||
        normalizedValue === 'true' ||
        normalizedValue === 'yes'
      );
    }

    return Boolean(value);
  }

  if (field.widgetType === 'datetime') {
    const parsedDate = new Date(String(value));

    if (!Number.isNaN(parsedDate.getTime())) {
      return parsedDate.toISOString().slice(0, 10);
    }
  }

  if (field.widgetType === 'list') {
    if (Array.isArray(value)) {
      return value.join(', ');
    }

    return String(value);
  }

  return value;
}

function hydrateFormWithExistingData(
  definition: FormDefinition,
  existingData?: Record<string, unknown>,
): FormDefinition {
  if (!existingData) {
    return definition;
  }

  return Object.entries(definition).reduce<FormDefinition>(
    (accumulator, [sectionName, section]) => {
      accumulator[sectionName] = {
        ...section,
        body: Object.entries(section.body).reduce<Record<string, FormField>>(
          (bodyAccumulator, [fieldName, field]) => {
            const matchingValue = existingData[fieldName];
            bodyAccumulator[fieldName] = {
              ...field,
              value:
                matchingValue === undefined
                  ? field.value
                  : normalizeExistingFieldValue(field, matchingValue),
            };

            return bodyAccumulator;
          },
          {},
        ),
      };

      return accumulator;
    },
    {},
  );
}

function resolveFieldValue(field: FormField): string {
  if (field.value === null || field.value === undefined) {
    return '';
  }

  if (typeof field.value === 'boolean') {
    return field.value ? 'true' : 'false';
  }

  if (Array.isArray(field.value)) {
    return field.value.join(', ');
  }

  return String(field.value);
}

function evaluateCondition(
  condition: FormCondition | undefined,
  definition: FormDefinition | null,
) {
  if (
    !condition ||
    !definition ||
    !condition.section ||
    !condition.targetName
  ) {
    return true;
  }

  const sectionKey = String(condition.section);
  const targetSection = definition[sectionKey];
  const targetField = targetSection?.body?.[condition.targetName];
  const targetValue = targetField?.value ?? targetField?.defaultValue;
  const normalizedTargetValue = String(targetValue ?? '').toLowerCase();
  const normalizedExpectedValue = String(condition.value ?? '').toLowerCase();

  if (condition.type === 'hasvalue') {
    return normalizedTargetValue.length > 0;
  }

  if (condition.type === 'isEmpty') {
    return normalizedTargetValue.length === 0;
  }

  return normalizedTargetValue === normalizedExpectedValue;
}

function buildPayload(definition: FormDefinition | null): FormPayload {
  if (!definition) {
    return {};
  }

  return Object.entries(definition).reduce<FormPayload>(
    (accumulator, [, section]) => {
      Object.entries(section.body).forEach(([fieldName, field]) => {
        accumulator[fieldName] = serializeSubmissionValue(field, field.value);
      });
      return accumulator;
    },
    {},
  );
}

function serializeSubmissionValue(field: FormField, value: unknown): unknown {
  if (field.type === 'double') {
    if (typeof value === 'number') {
      return value;
    }

    const parsedValue = Number(String(value).replace(/,/g, '').trim());
    return Number.isNaN(parsedValue) ? value : parsedValue;
  }

  if (field.type === 'int') {
    if (typeof value === 'number') {
      return Math.trunc(value);
    }

    const parsedValue = Number.parseInt(String(value).trim(), 10);
    return Number.isNaN(parsedValue) ? value : parsedValue;
  }

  if (
    field.widgetType === 'datetime' &&
    typeof value === 'string' &&
    value.trim()
  ) {
    const parsedDate = new Date(value);

    if (!Number.isNaN(parsedDate.getTime())) {
      const year = parsedDate.getUTCFullYear();
      const month = String(parsedDate.getUTCMonth() + 1).padStart(2, '0');
      const day = String(parsedDate.getUTCDate()).padStart(2, '0');
      const hours = String(parsedDate.getUTCHours()).padStart(2, '0');
      const minutes = String(parsedDate.getUTCMinutes()).padStart(2, '0');
      const seconds = String(parsedDate.getUTCSeconds()).padStart(2, '0');
      const milliseconds = String(parsedDate.getUTCMilliseconds()).padStart(
        3,
        '0',
      );
      const fraction = `${milliseconds}000`.slice(0, 6);

      return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}.${fraction}Z`;
    }
  }

  if (typeof value === 'boolean') return Number(value);

  return value;
}

function validateRequiredFields(
  definition: FormDefinition | null,
): Record<string, string> {
  if (!definition) {
    return {};
  }

  const errors: Record<string, string> = {};

  Object.entries(definition).forEach(([sectionName, section]) => {
    Object.entries(section.body).forEach(([fieldName, field]) => {
      const isVisible = evaluateCondition(field.when, definition);

      if (!field.required || !isVisible) {
        return;
      }

      const value = field.value;
      const isEmpty =
        value === null ||
        value === undefined ||
        value === '' ||
        (Array.isArray(value) && value.length === 0) ||
        (typeof value === 'boolean' && !value);

      if (isEmpty) {
        errors[`${sectionName}.${fieldName}`] = `${
          field.label ?? fieldName
        } is required`;
      }
    });
  });

  return errors;
}

export function TokenizationAssetInformation() {
  const [getForm, {}] = useLazyFetchFormFieldsQuery();
  const [submitTokenization] = useSubmitTokenizationMutation();
  const navigate = useNavigate();
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const { activeTokenizedAsset: activeTokenizationRecord, assetId } =
    useActiveTokenizedAsset();
  useTokenizationData();
  const availableFinancialAssetTypes = useSelector(
    (state: RootState) => state.appState.availableFinancialAssetTypes,
  );
  const [formDefinition, setFormDefinition] = useState<FormDefinition | null>(
    null,
  );
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [formBehaviors, setFormBehaviors] = useState<FormBehavior[]>([]);
  const onChangeBehaviorTimerRef = useRef<ReturnType<typeof setTimeout> | null>(
    null,
  );
  const pendingBehaviorFieldsRef = useRef<Set<string>>(new Set());
  const formBehaviorsRef = useRef<FormBehavior[]>([]);
  const fieldRefs = useRef<Map<string, HTMLElement>>(new Map());

  useEffect(() => {
    formBehaviorsRef.current = formBehaviors;
  }, [formBehaviors]);

  useEffect(() => {
    return () => {
      if (onChangeBehaviorTimerRef.current) {
        clearTimeout(onChangeBehaviorTimerRef.current);
      }
    };
  }, []);

  useEffect(() => {
    const loadFormDefinition = async () => {
      setIsLoading(true);
      setErrorMessage(null);

      try {
        const encryptor = new Encryptor();
        const currentSecretKey = await encryptor.getSecretKey(appUser);

        let resolvedPayload: unknown = null;
        showLoader();
        try {
          const assetType = activeTokenizationRecord?.assetType;
          const assetTypeIsAvailable = availableFinancialAssetTypes.includes(
            String(assetType ?? ''),
          );
          const formId = assetTypeIsAvailable
            ? String(assetType)
            : (activeTokenizationRecord?.assetAlreadyExists ?? 0) === 1
              ? 1
              : 2;

          const payload = {
            signer: appUser.primarySigner,
            publicKey: appUser.publicKey,
            secretKey: currentSecretKey,
            body: { formId },
          };

          const res = await getForm(payload);
          const responsePayload =
            res && typeof res === 'object' && 'data' in res && res.data
              ? (res.data as Record<string, unknown>).formString
              : null;

          resolvedPayload = parseApiFormPayload(responsePayload);
        } catch {
          // continue to the next endpoint if the request fails
        } finally {
          hideLoader();
        }

        const normalizedBehaviors = normalizeFormBehaviors(resolvedPayload);
        const normalizedDefinition = hydrateFormWithExistingData(
          seedFormValues(normalizeFormDefinition(resolvedPayload)),
          activeTokenizationRecord as unknown as Record<string, unknown>,
        );
        const withInitBehaviors = applyBehaviorEngine(
          normalizedDefinition,
          normalizedBehaviors,
          'onInit',
        );
        setFormBehaviors(normalizedBehaviors);
        setFormDefinition(withInitBehaviors);
      } catch {
        setErrorMessage(
          'Unable to load the form definition. Showing a fallback template instead.',
        );
        setFormBehaviors([]);
        setFormDefinition(
          hydrateFormWithExistingData(
            seedFormValues(fallbackFormDefinition),
            activeTokenizationRecord as unknown as Record<string, unknown>,
          ),
        );
      } finally {
        setIsLoading(false);
      }
    };

    void loadFormDefinition();
  }, [
    activeTokenizationRecord,
    appUser,
    getForm,
    availableFinancialAssetTypes,
  ]);

  const updateField = (
    sectionName: string,
    fieldName: string,
    value: unknown,
  ) => {
    setFormDefinition((currentDefinition) => {
      if (!currentDefinition) {
        return currentDefinition;
      }

      const nextDefinition = { ...currentDefinition };
      const section = nextDefinition[sectionName];
      if (!section) {
        return currentDefinition;
      }

      nextDefinition[sectionName] = {
        ...section,
        body: {
          ...section.body,
          [fieldName]: {
            ...section.body[fieldName],
            value,
          },
        },
      };

      setErrors((currentErrors) => {
        const nextErrors = { ...currentErrors };
        delete nextErrors[`${sectionName}.${fieldName}`];
        return nextErrors;
      });

      return nextDefinition;
    });

    pendingBehaviorFieldsRef.current.add(fieldName);

    if (onChangeBehaviorTimerRef.current) {
      clearTimeout(onChangeBehaviorTimerRef.current);
    }

    onChangeBehaviorTimerRef.current = setTimeout(() => {
      const changedFields = Array.from(pendingBehaviorFieldsRef.current);
      pendingBehaviorFieldsRef.current.clear();

      if (!changedFields.length) {
        return;
      }

      setFormDefinition((currentDefinition) => {
        if (!currentDefinition) {
          return currentDefinition;
        }

        return applyBehaviorEngine(
          currentDefinition,
          formBehaviorsRef.current,
          'onChange',
          changedFields,
        );
      });
    }, 1000);
  };

  const renderInputField = (
    sectionName: string,
    fieldName: string,
    field: FormField,
  ) => {
    const value = resolveFieldValue(field);

    switch (field.widgetType) {
      case 'datetime':
        return (
          <input
            type="date"
            className="mt-2 h-12 w-full rounded-md border border-gray-200 bg-primary-100 px-3 text-gray-700 outline-none focus:border-primary-600"
            value={value}
            onChange={(event) =>
              updateField(sectionName, fieldName, event.target.value)
            }
          />
        );
      case 'dropdown':
        return (
          <select
            className="mt-2 h-12 w-full rounded-md border border-gray-200 bg-primary-100 px-3 text-gray-700 outline-none focus:border-primary-600"
            value={value}
            onChange={(event) =>
              updateField(sectionName, fieldName, event.target.value)
            }
          >
            <option value="">
              {field.placeholderText ?? 'Select an option'}
            </option>
            {(field.options ?? []).map((option, index) => {
              if (typeof option === 'string') {
                return (
                  <option key={`${fieldName}-${index}`} value={option}>
                    {option}
                  </option>
                );
              }

              return (
                <option
                  key={`${fieldName}-${index}`}
                  value={String(option.value)}
                >
                  {option.text}
                </option>
              );
            })}
          </select>
        );
      case 'checkbox':
        return (
          <label className="mt-2 flex items-center gap-2 text-sm text-gray-600">
            <input
              type="checkbox"
              className="accent-primary-700"
              checked={Boolean(field.value)}
              onChange={(event) =>
                updateField(sectionName, fieldName, event.target.checked)
              }
            />
            <span>{field.label}</span>
          </label>
        );
      case 'radio':
        return (
          <div className="mt-3 space-y-2">
            {(field.options ?? []).map((option, index) => {
              if (typeof option === 'string') {
                return (
                  <label
                    key={`${fieldName}-${index}`}
                    className="flex items-center gap-2 cursor-pointer text-sm text-gray-700"
                  >
                    <input
                      type="radio"
                      className="accent-primary-700"
                      checked={String(field.value) === option}
                      onChange={() =>
                        updateField(sectionName, fieldName, option)
                      }
                    />
                    {option}
                  </label>
                );
              }

              return (
                <label
                  key={`${fieldName}-${index}`}
                  className="flex items-center gap-2 cursor-pointer text-sm text-gray-700"
                >
                  <input
                    type="radio"
                    className="accent-primary-700"
                    checked={String(field.value) === String(option.value)}
                    onChange={() =>
                      updateField(sectionName, fieldName, option.value)
                    }
                  />
                  {option.text}
                </label>
              );
            })}
          </div>
        );
      case 'list':
      case 'longtext':
        return (
          <TextArea
            inputType="textarea"
            label=""
            rows={4}
            defaultValue={value}
            placeholder={
              field.placeholderText ?? 'Enter comma-separated values'
            }
            onInputChange={(value) =>
              updateField(sectionName, fieldName, value)
            }
          />
        );
      default:
        if (field.type === 'double' || field.type === 'int') {
          return (
            <input
              type="number"
              step={field.type === 'double' ? '0.01' : '1'}
              className="mt-2 h-12 w-full rounded-md border border-gray-200 bg-primary-100 px-3 text-gray-700 outline-none focus:border-primary-600"
              value={value}
              placeholder={field.placeholderText ?? 'Enter a number'}
              onChange={(event) =>
                updateField(sectionName, fieldName, event.target.value)
              }
            />
          );
        }

        return (
          <TextInput
            defaultValue={value}
            placeholder={field.placeholderText ?? 'Enter a value'}
            label=""
            inputType="text"
            onInputChange={(nextValue) =>
              updateField(sectionName, fieldName, nextValue)
            }
          />
        );
    }
  };

  function excludeSubmissionFields(payload: any) {
    return Object.entries(payload).reduce((accumulator: any, [key, value]) => {
      if (!excludedSubmissionFields.has(key)) {
        accumulator[key] = value as never;
      }

      return accumulator;
    }, {});
  }

  const scrollToFirstError = (validationErrors: Record<string, string>) => {
    const firstErrorKey = Object.keys(validationErrors)[0];
    if (!firstErrorKey) {
      return;
    }

    const targetElement = fieldRefs.current.get(firstErrorKey);
    if (targetElement) {
      targetElement.scrollIntoView({
        behavior: 'smooth',
        block: 'center',
      });

      const focusable = targetElement.querySelector<HTMLElement>(
        'input, select, textarea, button',
      );
      focusable?.focus();
    }
  };

  const handleSubmit = async () => {
    const validationErrors = validateRequiredFields(formDefinition);
    setErrors(validationErrors);

    if (Object.keys(validationErrors).length > 0) {
      setErrorMessage('Please complete the required fields before continuing.');
      scrollToFirstError(validationErrors);
      return;
    }

    setIsSubmitting(true);
    setErrorMessage(null);
    showLoader();

    try {
      const encryptor = new Encryptor();
      const currentSecretKey = await encryptor.getSecretKey(appUser);
      const formPayload = buildPayload(formDefinition);
      const existingPayload: TokenizedAssetFormPayload = {
        ...(activeTokenizationRecord ?? {}),
      };

      const mergedPayload = {
        ...existingPayload,
        ...formPayload,
        id: existingPayload.id,
      };

      const payload = excludeSubmissionFields(mergedPayload);

      const res = await submitTokenization({
        signer: appUser.primarySigner,
        publicKey: appUser.publicKey,
        secretKey: currentSecretKey,
        body: payload,
      });

      hideLoader();

      if ('data' in res) {
        showNotification('success', 'Asset information saved successfully.');
        navigate(`/dashboard/tokenize/apply/${assetId}/asset-documents`);
      } else if ('error' in res) {
        const errorData = (res as any).error;
        showNotification(
          'error',
          errorData?.data?.message ?? 'Unable to save the form right now.',
        );
      } else {
        showNotification('error', 'Unable to save the form right now.');
      }
    } catch (submitError: any) {
      hideLoader();
      showNotification(
        'error',
        submitError?.message ?? 'Unable to save the form right now.',
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="bg-white p-6 w-full mx-auto">
      <div className="space-y-6">
        {isLoading ? (
          <div className="rounded-lg border border-gray-200 bg-primary-100 p-6 text-sm text-gray-600">
            Loading form fields...
          </div>
        ) : null}

        {errorMessage ? (
          <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-sm text-red-600">
            {errorMessage}
          </div>
        ) : null}

        {formDefinition
          ? Object.entries(formDefinition).map(
              ([sectionName, section], index, sections) => (
                <div key={sectionName}>
                  <div className="flex flex-col gap-6 md:flex-row md:gap-x-6 pb-8">
                    <div className="w-full md:w-2/6">
                      <label className="block font-medium text-[#0E1F51] mb-1">
                        {section.title ?? sectionName}
                      </label>
                      {section.description ? (
                        <p className="text-sm text-gray-500 mb-2 max-w-sm">
                          {section.description}
                        </p>
                      ) : null}
                    </div>

                    <div className="flex-1 w-full md:w-4/6 space-y-4">
                      {Object.entries(section.body).map(
                        ([fieldName, field]) => {
                          const isVisible = evaluateCondition(
                            field.when,
                            formDefinition,
                          );
                          if (!isVisible) {
                            return null;
                          }

                          return (
                            <div
                              key={`${sectionName}-${fieldName}`}
                              className="w-full"
                              ref={(node) => {
                                if (node) {
                                  fieldRefs.current.set(
                                    `${sectionName}.${fieldName}`,
                                    node,
                                  );
                                } else {
                                  fieldRefs.current.delete(
                                    `${sectionName}.${fieldName}`,
                                  );
                                }
                              }}
                            >
                              {field.widgetType !== 'checkbox' && (
                                <label className="block font-medium text-[#0E1F51] mb-1">
                                  {field.label ?? fieldName}
                                  {field.required ? (
                                    <span className="ml-1 text-red-500">*</span>
                                  ) : null}
                                </label>
                              )}
                              {renderInputField(sectionName, fieldName, field)}
                              {errors[`${sectionName}.${fieldName}`] ? (
                                <p className="mt-1 text-sm text-red-500">
                                  {errors[`${sectionName}.${fieldName}`]}
                                </p>
                              ) : null}
                            </div>
                          );
                        },
                      )}
                    </div>
                  </div>

                  {index < sections.length - 1 ? (
                    <hr className="border-gray-200" />
                  ) : null}
                </div>
              ),
            )
          : null}
      </div>

      <div className="w-full flex justify-end mt-10">
        <div className="flex self-end w-full max-w-2xl space-x-6">
          <div className="w-2/4">
            <ButtonSecondary
              label="Back"
              onclick={() => {
                navigate(-1);
              }}
            />
          </div>
          <div className="w-3/4">
            <Button
              label={isSubmitting ? 'Saving...' : 'Save and Continue'}
              onclick={() => {
                void handleSubmit();
              }}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
