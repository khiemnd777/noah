import React from "react";
import toast from "react-hot-toast";
import CloudUploadOutlinedIcon from "@mui/icons-material/CloudUploadOutlined";
import DownloadOutlinedIcon from "@mui/icons-material/DownloadOutlined";
import AddIcon from "@mui/icons-material/Add";
import EditOutlinedIcon from "@mui/icons-material/EditOutlined";
import DeleteOutlineIcon from "@mui/icons-material/DeleteOutline";
import SaveOutlinedIcon from "@mui/icons-material/SaveOutlined";
import {
  Alert,
  Autocomplete,
  Box,
  Button,
  Chip,
  Divider,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  Typography,
} from "@mui/material";
import { useParams } from "react-router-dom";
import { IfPermission } from "@root/core/auth/if-permission";
import { usePermissionChecks } from "@root/core/auth/rbac-utils";
import { AutoForm } from "@root/core/form/auto-form";
import { FormDialog } from "@root/core/form/form-dialog";
import type { AutoFormRef } from "@root/core/form/form.types";
import { openFormDialog } from "@root/core/form/form-dialog.service";
import { useI18n } from "@root/core/i18n/use-i18n";
import { registerSlot } from "@root/core/module/registry";
import { reloadTable } from "@root/core/table/table-reload";
import { exportLanguageXml, getLanguageById, listLanguages } from "@features/languages/api/language.api";
import type { LanguageModel, LanguageResourceModel } from "@features/languages/model/language.model";
import { ConfirmDialog } from "@shared/components/dialog/confirm-dialog";
import { SectionCard } from "@shared/components/ui/section-card";
import { SafeButton } from "@shared/components/button/safe-button";
import { TabContainer } from "@shared/components/ui/tab-container";
import { bytesToBlob } from "@shared/utils/file.utils";

const ADMIN_RESOURCE_KEY_REGEX = /^admin\.[^.]+\..+$/;
const DEFAULT_LANGUAGE_FETCH_LIMIT = 200;

function resolveModules(resources: LanguageResourceModel[]): string[] {
  const modules = new Set<string>();
  for (const resource of resources) {
    const matched = resource.key.match(/^admin\.([^.]+)\./);
    if (matched?.[1]) modules.add(matched[1]);
  }

  if (modules.size === 0) modules.add("general");
  return Array.from(modules).sort((a, b) => a.localeCompare(b));
}

function nextResourceId(resources: LanguageResourceModel[]): number {
  const currentIds = resources.map((item) => Number(item.id ?? 0)).filter((item) => Number.isFinite(item));
  const min = currentIds.length > 0 ? Math.min(...currentIds, 0) : 0;
  return min <= 0 ? min - 1 : -1;
}

function moduleFromKey(key: string): string | null {
  const matched = key.trim().match(/^admin\.([^.]+)\./);
  return matched?.[1] ?? null;
}

function findValidationMessage(
  draft: LanguageResourceModel,
  resources: LanguageResourceModel[],
  t: (key: string, fallback?: string) => string,
  ignoreIndex?: number,
): string | null {
  const key = String(draft.key ?? "").trim();
  if (!key) {
    return t("admin.languages.validation.resource_key_required");
  }
  if (!ADMIN_RESOURCE_KEY_REGEX.test(key)) {
    return t("admin.languages.validation.resource_key_invalid").replace("{key}", key);
  }

  const duplicated = resources.some((resource, index) => {
    if (ignoreIndex !== undefined && index === ignoreIndex) return false;
    return String(resource.key ?? "").trim() === key;
  });

  if (duplicated) {
    return t("admin.languages.validation.resource_key_duplicate").replace("{key}", key);
  }

  return null;
}

function readResourceValue(resource?: LanguageResourceModel | null): string {
  return resource?.value ?? "";
}

function buildResourceMap(resources: LanguageResourceModel[]): Map<string, LanguageResourceModel> {
  const nextMap = new Map<string, LanguageResourceModel>();
  for (const resource of resources) {
    nextMap.set(resource.key, resource);
  }
  return nextMap;
}

const compareFieldSx = {
  "& .MuiFormHelperText-root": {
    minHeight: 22,
  },
} as const;

type ResourceDialogProps = {
  open: boolean;
  mode: "add" | "edit";
  currentLanguageCode: string;
  defaultResources: Map<string, LanguageResourceModel>;
  defaultResourceOptions: string[];
  initialValue: LanguageResourceModel;
  resources: LanguageResourceModel[];
  editIndex?: number;
  onClose: () => void;
  onSubmit: (resource: LanguageResourceModel) => void;
};

function ResourceEditorDialog({
  open,
  mode,
  currentLanguageCode,
  defaultResources,
  defaultResourceOptions,
  initialValue,
  resources,
  editIndex,
  onClose,
  onSubmit,
}: ResourceDialogProps) {
  const { t } = useI18n();
  const [draft, setDraft] = React.useState<LanguageResourceModel>(initialValue);
  const [error, setError] = React.useState<string | null>(null);

  React.useEffect(() => {
    setDraft(initialValue);
    setError(null);
  }, [initialValue, open]);

  const handleChange = React.useCallback((patch: Partial<LanguageResourceModel>) => {
    setDraft((current) => ({ ...current, ...patch }));
    setError(null);
  }, []);

  const handleSubmit = React.useCallback(() => {
    const validationMessage = findValidationMessage(draft, resources, t, mode === "edit" ? editIndex : undefined);
    if (validationMessage) {
      setError(validationMessage);
      return;
    }

    onSubmit({
      ...draft,
      key: String(draft.key ?? "").trim(),
    });
  }, [draft, editIndex, mode, onSubmit, resources, t]);

  const defaultResource = React.useMemo(() => defaultResources.get(String(draft.key ?? "").trim()) ?? null, [defaultResources, draft.key]);
  const defaultKey = defaultResource?.key ?? draft.key ?? "";
  const defaultValue = readResourceValue(defaultResource);

  return (
    <FormDialog
      open={open}
      title={
        mode === "edit"
          ? t("admin.languages.resources.edit_title")
          : t("admin.languages.resources.add_title")
      }
      confirmText={mode === "edit" ? t("admin.general.update_button") : t("admin.general.create_button")}
      cancelText={t("admin.general.cancel_button")}
      onClose={onClose}
      onSubmit={handleSubmit}
      maxWidth="lg"
    >
      <Stack spacing={2}>
        {error ? <Alert severity="error">{error}</Alert> : null}
        <Box
          sx={{
            display: "grid",
            gap: 3,
            gridTemplateColumns: {
              xs: "1fr",
              md: "repeat(2, minmax(0, 1fr))",
            },
            alignItems: "start",
          }}
        >
          <Stack spacing={2.5} sx={{ minWidth: 0 }}>
            <Typography variant="subtitle1" fontWeight={600}>
              {t("admin.languages.resources.default_column")}
            </Typography>
            {mode === "add" ? (
              <Autocomplete
                freeSolo
                options={defaultResourceOptions}
                value={null}
                inputValue={draft.key}
                onInputChange={(_, value, reason) => {
                  if (reason === "reset") return;
                  handleChange({ key: value });
                }}
                onChange={(_, value) => {
                  handleChange({ key: typeof value === "string" ? value : value ?? "" });
                }}
                renderInput={(params) => (
                  <TextField
                    {...params}
                    label={t("admin.languages.resources.default_search_key")}
                    helperText={t("admin.languages.resources.default_search_helper")}
                    sx={compareFieldSx}
                  />
                )}
              />
            ) : null}
            <TextField
              fullWidth
              label={t("admin.languages.resources.key_label")}
              value={defaultKey}
              InputProps={{ readOnly: true }}
              helperText=" "
              sx={compareFieldSx}
            />
            <TextField
              fullWidth
              multiline
              minRows={6}
              label={t("admin.languages.resources.value_label")}
              value={defaultValue}
              InputProps={{ readOnly: true }}
              helperText={
                defaultResource
                  ? undefined
                  : t("admin.languages.resources.default_missing")
              }
              sx={compareFieldSx}
            />
          </Stack>

          <Stack spacing={2.5} sx={{ minWidth: 0 }}>
            <Typography variant="subtitle1" fontWeight={600}>
              {t("admin.languages.resources.current_column").replace("{code}", currentLanguageCode)}
            </Typography>
            <TextField
              fullWidth
              label={t("admin.languages.resources.key_label")}
              value={draft.key}
              onChange={(event) => handleChange({ key: event.target.value })}
              helperText={t("admin.languages.resources.key_helper")}
              sx={compareFieldSx}
            />
            <TextField
              fullWidth
              multiline
              minRows={6}
              label={t("admin.languages.resources.value_label")}
              value={draft.value}
              onChange={(event) => handleChange({ value: event.target.value })}
              helperText=" "
              sx={compareFieldSx}
            />
          </Stack>
        </Box>
      </Stack>
    </FormDialog>
  );
}

function saveBlob(bytes: Uint8Array, filename: string, mime: string) {
  const blob = bytesToBlob(bytes, mime);
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  a.click();
  URL.revokeObjectURL(url);
}

function getErrorMessage(error: unknown, fallback: string): string {
  if (typeof error === "object" && error !== null) {
    const response = (error as { response?: { data?: { statusMessage?: string } } }).response;
    if (response?.data?.statusMessage) return response.data.statusMessage;
  }
  return fallback;
}

export function LanguageResourcesEditor({
  resources,
  defaultResources,
  currentLanguageCode,
  disabled,
  onChange,
}: {
  resources: LanguageResourceModel[];
  defaultResources: Map<string, LanguageResourceModel>;
  currentLanguageCode: string;
  disabled: boolean;
  onChange: (resources: LanguageResourceModel[]) => void;
}) {
  const { t } = useI18n();
  const modules = React.useMemo(() => resolveModules(resources), [resources]);
  const [activeModule, setActiveModule] = React.useState<string>(modules[0] ?? "general");
  const [editingIndex, setEditingIndex] = React.useState<number | null>(null);
  const [adding, setAdding] = React.useState(false);
  const [deleteIndex, setDeleteIndex] = React.useState<number | null>(null);

  React.useEffect(() => {
    if (!modules.includes(activeModule)) {
      setActiveModule(modules[0] ?? "general");
    }
  }, [activeModule, modules]);

  const handleDelete = React.useCallback((index: number) => {
    const next = resources.filter((_, itemIndex) => itemIndex !== index);
    onChange(next);
  }, [onChange, resources]);

  const deletingResource = deleteIndex !== null ? resources[deleteIndex] ?? null : null;

  const handleConfirmDelete = React.useCallback(() => {
    if (deleteIndex === null) return;
    handleDelete(deleteIndex);
    setDeleteIndex(null);
  }, [deleteIndex, handleDelete]);

  const defaultResourceOptions = React.useMemo(() => Array.from(defaultResources.keys()).sort((a, b) => a.localeCompare(b)), [defaultResources]);

  const editingResource = editingIndex !== null ? resources[editingIndex] : null;
  const handleEditSave = React.useCallback((draft: LanguageResourceModel) => {
    if (editingIndex === null) return;

    const next = resources.map((item, index) => index === editingIndex ? { ...item, ...draft } : item);
    onChange(next);
    setEditingIndex(null);

    const nextModule = moduleFromKey(draft.key);
    if (nextModule) setActiveModule(nextModule);
  }, [editingIndex, onChange, resources]);

  const handleAddSave = React.useCallback((draft: LanguageResourceModel) => {
    const next = [
      ...resources,
      {
        ...draft,
        id: nextResourceId(resources),
      },
    ];
    onChange(next);
    setAdding(false);

    const nextModule = moduleFromKey(draft.key);
    if (nextModule) setActiveModule(nextModule);
  }, [onChange, resources]);

  const tabs = modules.map((moduleName) => {
    const moduleRows = resources
      .map((resource, index) => ({ resource, index }))
      .filter(({ resource }) => resource.key.startsWith(`admin.${moduleName}.`));

    return {
      label: `${moduleName} (${moduleRows.length})`,
      value: moduleName,
      content: (
        <Stack spacing={1.5}>
          {moduleRows.length === 0 ? (
            <Alert severity="info">{t("admin.languages.resources.empty_module")}</Alert>
            
          ) : (
            <TableContainer
              sx={(theme) => ({
                border: `1px solid ${theme.palette.divider}`,
                borderRadius: 2,
                overflowX: "auto",
              })}
            >
              <Table size="small">
                <TableHead>
                  <TableRow>
                    <TableCell width={180}>{t("admin.languages.resources.actions_header")}</TableCell>
                    <TableCell width="35%">{t("admin.languages.resources.key_label")}</TableCell>
                    <TableCell>{t("admin.languages.resources.value_label")}</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {moduleRows.map(({ resource, index }) => (
                    <TableRow key={`${resource.id ?? "new"}-${index}`} hover>
                      <TableCell>
                        <Stack direction="row" spacing={1}>
                          <Button
                            size="small"
                            variant="outlined"
                            startIcon={<EditOutlinedIcon />}
                            disabled={disabled}
                            onClick={() => setEditingIndex(index)}
                          >
                            {t("admin.general.edit_button")}
                          </Button>
                          <Button
                            size="small"
                            color="error"
                            variant="outlined"
                            startIcon={<DeleteOutlineIcon />}
                            disabled={disabled}
                            onClick={() => setDeleteIndex(index)}
                          >
                            {t("admin.general.delete_button")}
                          </Button>
                        </Stack>
                      </TableCell>
                      <TableCell sx={{ verticalAlign: "top", fontFamily: "monospace", whiteSpace: "nowrap" }}>
                        {resource.key}
                      </TableCell>
                      <TableCell sx={{ verticalAlign: "top" }}>
                        <Typography variant="body2" sx={{ whiteSpace: "pre-wrap", wordBreak: "break-word" }}>
                          {resource.value}
                        </Typography>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          )}
        </Stack>
      ),
    };
  });

  return (
    <Stack spacing={2}>
      <Stack direction="row" spacing={1} justifyContent="space-between" alignItems="center">
        <Stack direction="row" spacing={1} alignItems="center">
          <Typography variant="body2" color="text.secondary">
            {t("admin.languages.resources.format_prefix")}
          </Typography>
          <Chip size="small" label={t("admin.languages.resources.format_chip")} />
        </Stack>
        <Button variant="outlined" startIcon={<AddIcon />} disabled={disabled} onClick={() => setAdding(true)}>
          {t("admin.languages.resources.add")}
        </Button>
      </Stack>
      <TabContainer tabs={tabs} defaultValue={activeModule} onChange={setActiveModule} />
      <ResourceEditorDialog
        open={editingIndex !== null && !!editingResource}
        mode="edit"
        currentLanguageCode={currentLanguageCode}
        defaultResources={defaultResources}
        defaultResourceOptions={defaultResourceOptions}
        initialValue={editingResource ?? { key: "", value: "" }}
        resources={resources}
        editIndex={editingIndex ?? undefined}
        onClose={() => setEditingIndex(null)}
        onSubmit={handleEditSave}
      />
      <ResourceEditorDialog
        open={adding}
        mode="add"
        currentLanguageCode={currentLanguageCode}
        defaultResources={defaultResources}
        defaultResourceOptions={defaultResourceOptions}
        initialValue={{
          id: nextResourceId(resources),
          key: `admin.${activeModule || "general"}.new_key`,
          value: "",
        }}
        resources={resources}
        onClose={() => setAdding(false)}
        onSubmit={handleAddSave}
      />
      <ConfirmDialog
        open={deleteIndex !== null && !!deletingResource}
        title={t("admin.languages.resources.delete_confirm_title")}
        content={
          deletingResource
            ? t("admin.languages.resources.delete_confirm_content").replace("{key}", deletingResource.key)
            : t("admin.languages.resources.delete_confirm_fallback")
        }
        confirmText={t("admin.general.delete_button")}
        cancelText={t("admin.general.cancel_button")}
        onClose={() => setDeleteIndex(null)}
        onConfirm={handleConfirmDelete}
      />
    </Stack>
  );
}

export function LanguageDetailWidget() {
  const { t } = useI18n();
  const { languageId } = useParams();
  const { hasPermission } = usePermissionChecks();
  const formRef = React.useRef<AutoFormRef>(null);
  const [loading, setLoading] = React.useState(true);
  const [detail, setDetail] = React.useState<LanguageModel | null>(null);
  const [resources, setResources] = React.useState<LanguageResourceModel[]>([]);
  const [defaultResources, setDefaultResources] = React.useState<Map<string, LanguageResourceModel>>(new Map());

  const numericLanguageId = Number(languageId ?? 0);

  const loadDetail = React.useCallback(async () => {
    if (!numericLanguageId) return;
    setLoading(true);
    try {
      const data = await getLanguageById(numericLanguageId);
      setDetail(data);
      setResources(data.resources ?? []);
    } catch (error: unknown) {
      toast.error(getErrorMessage(error, t("admin.languages.messages.detail_load_failed")));
    } finally {
      setLoading(false);
    }
  }, [numericLanguageId, t]);

  React.useEffect(() => {
    void loadDetail();
  }, [loadDetail]);

  React.useEffect(() => {
    formRef.current?.setValue("resources", resources);
  }, [resources]);

  React.useEffect(() => {
    let active = true;

    async function loadDefaultResources() {
      if (!detail?.id) {
        if (active) setDefaultResources(new Map());
        return;
      }

      try {
        if (detail.isDefault) {
          if (!active) return;
          setDefaultResources(buildResourceMap(resources));
          return;
        }

        let page = 0;
        let defaultLanguageId: number | undefined;
        let total: number | null = null;

        while (defaultLanguageId === undefined && (total === null || page * DEFAULT_LANGUAGE_FETCH_LIMIT < total)) {
          const listing = await listLanguages({
            page,
            limit: DEFAULT_LANGUAGE_FETCH_LIMIT,
          });
          total = listing.total;

          const defaultLanguage = listing.items.find((item) => item.isDefault);
          if (defaultLanguage?.id && defaultLanguage.id !== detail.id) {
            defaultLanguageId = Number(defaultLanguage.id);
            break;
          }

          if (listing.items.length < DEFAULT_LANGUAGE_FETCH_LIMIT) {
            break;
          }

          page += 1;
        }

        if (!defaultLanguageId) {
          if (!active) return;
          setDefaultResources(new Map());
          toast.error(t("admin.languages.messages.default_load_failed"));
          return;
        }

        const source = await getLanguageById(defaultLanguageId);

        if (!active) return;
        setDefaultResources(buildResourceMap(source.resources ?? []));
      } catch (error: unknown) {
        if (!active) return;
        setDefaultResources(new Map());
        toast.error(getErrorMessage(error, t("admin.languages.messages.default_load_failed")));
      }
    }

    void loadDefaultResources();

    return () => {
      active = false;
    };
  }, [detail, resources, t]);

  const validateResources = React.useCallback(() => {
    const seen = new Set<string>();

    for (const resource of resources) {
      const key = String(resource.key ?? "").trim();
      if (!key) {
        toast.error(t("admin.languages.validation.resource_key_required"));
        return false;
      }
      if (!ADMIN_RESOURCE_KEY_REGEX.test(key)) {
        toast.error(
          t("admin.languages.validation.resource_key_invalid").replace("{key}", key)
        );
        return false;
      }
      if (seen.has(key)) {
        toast.error(t("admin.languages.validation.resource_key_duplicate").replace("{key}", key));
        return false;
      }
      seen.add(key);
    }

    return true;
  }, [resources, t]);

  const handleSave = React.useCallback(async () => {
    if (!validateResources()) return;
    formRef.current?.setValue("resources", resources);
    const ok = await formRef.current?.submit();
    if (ok) {
      await loadDetail();
      reloadTable("languages");
    }
  }, [loadDetail, resources, validateResources]);

  const handleExport = React.useCallback(async () => {
    if (!numericLanguageId) return;
    try {
      const file = await exportLanguageXml(numericLanguageId);
      saveBlob(file.bytes, file.filename, file.contentType);
      toast.success(t("admin.languages.messages.export_success"));
    } catch (error: unknown) {
      toast.error(getErrorMessage(error, t("admin.languages.messages.export_failed")));
    }
  }, [numericLanguageId, t]);

  const handleOpenImportDialog = React.useCallback(() => {
    if (!numericLanguageId) return;
    void openFormDialog("language-import", {
      initial: { languageId: numericLanguageId },
      onSaved: async () => {
        await loadDetail();
        reloadTable("languages");
      },
    });
  }, [loadDetail, numericLanguageId]);

  if (!languageId) {
    return <Alert severity="error">{t("admin.languages.messages.missing_id")}</Alert>;
  }

  if (loading) {
    return <Alert severity="info">{t("admin.languages.messages.detail_loading")}</Alert>;
  }

  if (!detail) {
    return <Alert severity="error">{t("admin.languages.messages.detail_not_found")}</Alert>;
  }

  return (
    <Stack spacing={2}>
      <SectionCard
        title={t("admin.languages.detail.metadata_title")}
        extra={
          <IfPermission permissions={["languages.update"]}>
            <SafeButton variant="contained" startIcon={<SaveOutlinedIcon />} onClick={() => void handleSave()}>
              {t("admin.general.save_button")}
            </SafeButton>
          </IfPermission>
        }
      >
        <AutoForm
          key={`language-detail-${detail.id ?? languageId}-${detail.updatedAt ?? "initial"}`}
          name="language-detail"
          ref={formRef}
          initial={detail}
        />
      </SectionCard>

      <SectionCard
        title={t("admin.languages.detail.resources_title")}
        extra={
          <Stack direction={{ xs: "column", sm: "row" }} spacing={1}>
            <IfPermission permissions={["languages.export"]}>
              <Button variant="outlined" startIcon={<DownloadOutlinedIcon />} onClick={() => void handleExport()}>
                {t("admin.general.export_button")}
              </Button>
            </IfPermission>
            <IfPermission permissions={["languages.import"]}>
              <Button variant="outlined" startIcon={<CloudUploadOutlinedIcon />} onClick={handleOpenImportDialog}>
                {t("admin.general.import_button")}
              </Button>
            </IfPermission>
          </Stack>
        }
      >
        <Stack spacing={2}>
          <Alert severity="info">
            {t("admin.languages.resources.info")} <strong>admin.{"{module}"}.*</strong>. {t("admin.languages.resources.example_prefix")} <strong>admin.settings.title</strong>.
          </Alert>
          <Divider />
          <LanguageResourcesEditor
            resources={resources}
            defaultResources={defaultResources}
            currentLanguageCode={detail.code}
            disabled={!hasPermission?.("languages.update")}
            onChange={(next) => {
              setResources(next);
              formRef.current?.setValue("resources", next);
            }}
          />
        </Stack>
      </SectionCard>
    </Stack>
  );
}

registerSlot({
  id: "languages-detail",
  name: "languages-detail:left",
  render: () => <LanguageDetailWidget />,
});
