import { env } from "@core/config/env";
import { mapper } from "@core/mapper/auto-mapper";
import { apiClient } from "@core/network/api-client";
import type { FetchTableOpts } from "@core/table/table.types";
import type { ListResult } from "@core/types/list-result";
import type { LanguageModel } from "@features/languages/model/language.model";

const languageBasePath = `${env.apiBasePath}/i18n/languages`;

type LanguageResourceDto = {
  id?: number;
  key: string;
  value: string;
  created_at?: string | null;
  updated_at?: string | null;
};

type LanguageDto = {
  id?: number;
  code: string;
  name: string;
  native_name?: string | null;
  is_default?: boolean;
  active?: boolean;
  created_at?: string | null;
  updated_at?: string | null;
  resources?: LanguageResourceDto[];
};

type LanguageListDto = ListResult<LanguageDto>;

function languageDetailTag(id?: number | string) {
  return `languages:detail:${id ?? "unknown"}`;
}

export async function listLanguages(tableOpts: FetchTableOpts): Promise<ListResult<LanguageModel>> {
  const { data } = await apiClient.getTable<LanguageListDto>(languageBasePath, tableOpts, {
    cacheMode: "stale-while-revalidate",
    cacheTags: ["languages:list"],
  });
  return mapper.map<LanguageListDto, ListResult<LanguageModel>>("Language", data, "dto_to_model");
}

export async function getLanguageById(id: number | string): Promise<LanguageModel> {
  const { data } = await apiClient.get<LanguageDto>(`${languageBasePath}/${id}`, {
    cacheMode: "stale-while-revalidate",
    cacheTags: [languageDetailTag(id)],
  });
  return mapper.map<LanguageDto, LanguageModel>("Language", data, "dto_to_model");
}

export async function createLanguage(model: LanguageModel): Promise<LanguageModel> {
  const dto = mapper.map<LanguageModel, LanguageDto>("Language", model, "model_to_dto");
  const { data } = await apiClient.post<LanguageDto>(languageBasePath, dto, {
    invalidateTagPrefixes: ["languages:"],
  });
  return mapper.map<LanguageDto, LanguageModel>("Language", data, "dto_to_model");
}

export async function updateLanguage(model: LanguageModel): Promise<LanguageModel> {
  const dto = mapper.map<LanguageModel, LanguageDto>("Language", model, "model_to_dto");
  const { data } = await apiClient.put<LanguageDto>(`${languageBasePath}/${model.id}`, dto, {
    invalidateTagPrefixes: ["languages:"],
  });
  return mapper.map<LanguageDto, LanguageModel>("Language", data, "dto_to_model");
}

export async function deleteLanguage(id: number): Promise<{ success: boolean }> {
  const { data } = await apiClient.delete<{ success: boolean }>(`${languageBasePath}/${id}`, {
    invalidateTagPrefixes: ["languages:"],
  });
  return data;
}

export async function importLanguageXml(id: number, file: File): Promise<void> {
  const formData = new FormData();
  formData.append("file", file);

  await apiClient.post(`${languageBasePath}/${id}/import-xml`, formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
    invalidateTagPrefixes: ["languages:"],
  });
}

export async function exportLanguageXml(id: number): Promise<{ bytes: Uint8Array; filename: string; contentType: string }> {
  const res = await apiClient.get<ArrayBuffer>(`${languageBasePath}/${id}/export-xml`, {
    responseType: "arraybuffer",
  });

  const contentType =
    (res.headers["content-type"] as string | undefined) ??
    (res.headers["Content-Type"] as string | undefined) ??
    "application/xml";

  const disposition =
    (res.headers["content-disposition"] as string | undefined) ??
    (res.headers["Content-Disposition"] as string | undefined) ??
    "";

  const matched = disposition.match(/filename\*?=(?:UTF-8''|")?([^";]+)/i);
  const filename = matched?.[1] ? decodeURIComponent(matched[1].replace(/"/g, "")) : `language-${id}.xml`;

  return {
    bytes: new Uint8Array(res.data),
    filename,
    contentType,
  };
}
