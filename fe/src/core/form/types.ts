import type { LocalizedText } from "@root/core/i18n/localized-text";

export type FieldKind =
  | "text"
  | "password"
  | "new-password"
  | "change-password"
  | "email"
  | "textarea"
  | "date"
  | "datetime"
  | "color"
  | "currency"
  | "currency-equation"
  | "select"
  | "checkbox"
  | "switch"
  | "number"
  | "multiselect"
  | "autocomplete"
  | "fileupload"
  | "imageupload"
  | "qr"
  | "custom"
  | "searchlist"
  | "searchsingle"
  | "metadata"
  | "relation";

export type DeriveMode = "always" | "whenEmpty" | "untilManual";

export type FormContext = {
  formSessionId: string | null;
  metadataBlocks: { meta: FieldDef; fields: FieldDef[]; deps: string[]; collections: string[] }[];
  values: Record<string, any>;
  setValue: (name: string, v: any) => void;
  setAllValues: (obj: Record<string, any>) => void;
  setFieldError: (name: string, msg: string | null) => void;
  reset: () => void;
  setInitial: (obj: Record<string, any>) => void;
  clear: () => void;
  emit: (event: string, payload?: any) => void;
  on: (event: string, handler: (payload: any) => void) => void;
  off: (event: string, handler: (payload: any) => void) => void;
};

export type PasswordRules = {
  minLength?: number;
  maxLength?: number;
  requireUpper?: boolean;
  requireLower?: boolean;
  requireDigit?: boolean;
  requireSymbol?: boolean;
  disallowSpaces?: boolean;
  disallowReuseCurrent?: boolean;
  custom?: (pw: string, allValues: Record<string, any>) => string | null | undefined;
};

export type FieldRules = {
  required?: boolean | string;
  minLength?: number;
  maxLength?: number;
  min?: number;
  max?: number;
  pattern?: RegExp | { regex: RegExp; message?: string };
  minDateTime?: string;
  maxDateTime?: string;
  custom?: (value: any) => string | null | undefined;
  async?: (value: any, allValues: Record<string, any>) => Promise<string | null | undefined>;
};

export type Option = {
  label: string;
  value: string | number | boolean;
};

export type SearchSingleValidateTrigger = "blur" | "select" | "input" | "clear";

export type QROptions = {
  size?: number;
  tooltipSize?: number;
  level?: "L" | "M" | "Q" | "H";
  fgColor?: string;
  bgColor?: string;
};

export type CustomRenderCtx = {
  value: any;
  setValue: (v: any) => void;
  error?: string | null;
  field: FieldDef;
  values: Record<string, any>;
  ctx?: FormContext | null;
};

export type SearchListSearchFn = (keyword: string, ctx?: FormContext) => Promise<any[]>;
export type SearchListSearchPageFn = (
  keyword: string,
  page: number,
  limit: number,
  ctx?: FormContext
) => Promise<any[]>;
export type SearchListFetchListFn = (
  values: Record<string, any>,
  ctx?: FormContext
) => Promise<any[]>;
export type SearchListHydrateFn = (
  ids: Array<string | number>,
  values: Record<string, any>,
  ctx?: FormContext
) => Promise<any[]>;

export type MiniFieldOverride = {
  name: string;
  col?: number;
  label?: LocalizedText;
  placeholder?: LocalizedText;
  helperText?: LocalizedText;
  defaultNowOnInteract?: boolean;
  rules?: FieldRules;
  showIf?: (values: Record<string, any>, ctx?: FormContext) => boolean;
  disableIf?: (values: Record<string, any>, ctx?: FormContext) => boolean;
  asText?: boolean;

  onBlur?: (text: string, matched: any, ctx?: FormContext | null) => void;
  onSelect?: (item: any) => void;
  onChange?: (value: any, ctx?: FormContext) => void;
  onInputChange?: (text: string) => void;
  validate?: (
    input: string,
    matched: any | null,
    ctx?: FormContext | null
  ) => string | null | undefined;
  validateAsync?: (
    input: string,
    matched: any | null,
    ctx?: FormContext | null
  ) => Promise<string | null | undefined>;
  validateOn?: SearchSingleValidateTrigger | SearchSingleValidateTrigger[];
  onValidate?: (
    message: string | null,
    input: string,
    matched: any | null,
    ctx?: FormContext | null
  ) => void;
  onDragEnd?: (items: any[]) => void;

  where?: (values: Record<string, any>, ctx?: FormContext) => string[];
  searchPage?: SearchListSearchPageFn;
  fetchOne?: (values: Record<string, any>) => Promise<any | null>;
  hydrateById?: (id: string | number, values: Record<string, any>) => Promise<any | null>;
  hydrateOrderField?: string;
  getOptionLabel?: (item: any) => string;
  getInputLabel?: (item: any) => string;
  renderItem?: (item: any, index?: number) => React.ReactNode;
};

export type MetadataFieldPlacement = {
  group: string;
  section?: string;
  fields?: string[];
};

export type FieldDef = {
  name: string;
  altName?: string;
  label: LocalizedText;
  kind: FieldKind;
  group?: string;
  section?: string;
  col?: number;
  placeholder?: LocalizedText;
  rows?: number;
  defaultValue?: any;
  defaultNowOnInteract?: boolean;
  helperText?: LocalizedText;
  fullWidth?: boolean;
  size?: "small" | "medium";
  rules?: FieldRules;
  qr?: QROptions;
  step?: number;
  showIf?: (values: Record<string, any>, ctx?: FormContext) => boolean;
  disableIf?: (values: Record<string, any>, ctx?: FormContext) => boolean;
  asTextFn?: (values: Record<string, any>, ctx?: FormContext) => boolean;
  asText?: boolean;

  options?: Option[];
  loadOptions?: (keyword: string) => Promise<Option[]>;
  freeSolo?: boolean;
  multiple?: boolean;
  debounceMs?: number;

  accept?: string;
  uploader?: (files: File[]) => Promise<string[]>;
  maxFiles?: number;
  multipleFiles?: boolean;
  imagePreviewAspectRatio?: string;
  imagePreviewHeight?: number;

  passwordRules?: PasswordRules;
  currentLabel?: LocalizedText;
  newLabel?: LocalizedText;
  confirmLabel?: LocalizedText;

  render?: (ctx: CustomRenderCtx) => React.ReactNode;
  normalizeInitial?: (value: any, allValues?: Record<string, any>) => any;

  derive?: {
    field: string;
    map: (sourceValue: any, values: Record<string, any>) => any;
    mode?: DeriveMode;
  };

  where?: (values: Record<string, any>, ctx?: FormContext) => string[];
  search?: SearchListSearchFn;
  searchPage?: SearchListSearchPageFn;
  fetchList?: SearchListFetchListFn;
  hydrateByIds?: SearchListHydrateFn;
  resolveDefaultInput?: (
    values: Record<string, any>,
    ctx?: FormContext
  ) => Promise<{ inputValue?: string; value?: any | null } | null>;
  onSelect?: (item: any) => void;
  onBlur?: (text: string, matched: any, ctx?: FormContext | null) => void;
  onAdd?: (item: any) => Promise<void> | void;
  onDelete?: (item: any) => Promise<void> | void;
  onDragEnd?: (items: any[]) => void;
  getOptionLabel?: (item: any, items?: any[]) => string;
  getOptionValue?: (item: any) => string | number;
  getInputLabel?: (item: any) => string;
  prop?: string;

  allowUnmatched?: boolean;
  fetchOne?: (values: Record<string, any>) => Promise<any | null>;
  hydrateById?: (id: string | number, values: Record<string, any>) => Promise<any | null>;
  validate?: (
    input: string,
    matched: any | null,
    ctx?: FormContext | null
  ) => string | null | undefined;
  validateAsync?: (
    input: string,
    matched: any | null,
    ctx?: FormContext | null
  ) => Promise<string | null | undefined>;
  validateOn?: SearchSingleValidateTrigger | SearchSingleValidateTrigger[];
  onValidate?: (
    message: string | null,
    input: string,
    matched: any | null,
    ctx?: FormContext | null
  ) => void;

  metadata?: {
    collection?: string;
    collectionFn?: (ctx: FormContext) => string;
    group?: string;
    tag?: string | null;
    mode?: "whole" | "partial";
    fields?: string[];
    ignoreFields?: string[];
    showIfFields?: string[];
    groups?: MetadataFieldPlacement[];
    def?: MiniFieldOverride[];
  };

  renderItem?: (item: any, index?: number) => React.ReactNode;

  allowDuplicate?: boolean;
  dedupeFn?: (a: any, b: any) => boolean;
  maxItems?: number;
  disableDelete?: (item: any) => boolean;

  onOpenCreate?: () => void;
  refreshKey?: any;
  autoLoadAllOnMount?: boolean;
  pageLimit?: number;
  hydrateOrderField?: string;

  currencyEquation?: string;
};

export type AutoFormOptions = {
  asyncValidate?: (
    values: Record<string, any>
  ) => Promise<Partial<Record<string, string | null>>>;
  asyncDebounceMs?: number;
};
