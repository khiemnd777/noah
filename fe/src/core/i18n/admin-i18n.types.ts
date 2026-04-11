export interface AdminLanguageOption {
  code: string;
  name: string;
  nativeName?: string | null;
  isDefault?: boolean;
  active?: boolean;
}

export interface AdminLanguagePreference {
  code: string | null;
}

export type AdminResourceDictionary = Record<string, string>;
