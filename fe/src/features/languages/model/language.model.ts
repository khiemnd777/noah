export interface LanguageResourceModel {
  id?: number;
  key: string;
  value: string;
  createdAt?: string | null;
  updatedAt?: string | null;
}

export interface LanguageModel {
  id?: number;
  code: string;
  name: string;
  nativeName?: string | null;
  isDefault?: boolean;
  active?: boolean;
  createdAt?: string | null;
  updatedAt?: string | null;
  resources?: LanguageResourceModel[];
}
