import React from "react";
import toast from "react-hot-toast";
import CloudUploadOutlinedIcon from "@mui/icons-material/CloudUploadOutlined";
import DownloadOutlinedIcon from "@mui/icons-material/DownloadOutlined";
import AddIcon from "@mui/icons-material/Add";
import DeleteOutlineIcon from "@mui/icons-material/DeleteOutline";
import SaveOutlinedIcon from "@mui/icons-material/SaveOutlined";
import {
  Alert,
  Box,
  Button,
  Chip,
  Divider,
  IconButton,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import { useParams } from "react-router-dom";
import { IfPermission } from "@root/core/auth/if-permission";
import { usePermissionChecks } from "@root/core/auth/rbac-utils";
import { AutoForm } from "@root/core/form/auto-form";
import type { AutoFormRef } from "@root/core/form/form.types";
import { openFormDialog } from "@root/core/form/form-dialog.service";
import { registerSlot } from "@root/core/module/registry";
import { reloadTable } from "@root/core/table/table-reload";
import { exportLanguageXml, getLanguageById } from "@features/languages/api/language.api";
import type { LanguageModel, LanguageResourceModel } from "@features/languages/model/language.model";
import { SectionCard } from "@shared/components/ui/section-card";
import { SafeButton } from "@shared/components/button/safe-button";
import { TabContainer } from "@shared/components/ui/tab-container";
import { bytesToBlob } from "@shared/utils/file.utils";

const ADMIN_RESOURCE_KEY_REGEX = /^admin\.[^.]+\..+$/;

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
  disabled,
  onChange,
}: {
  resources: LanguageResourceModel[];
  disabled: boolean;
  onChange: (resources: LanguageResourceModel[]) => void;
}) {
  const modules = React.useMemo(() => resolveModules(resources), [resources]);
  const [activeModule, setActiveModule] = React.useState<string>(modules[0] ?? "general");

  React.useEffect(() => {
    if (!modules.includes(activeModule)) {
      setActiveModule(modules[0] ?? "general");
    }
  }, [activeModule, modules]);

  const handlePatch = React.useCallback((index: number, patch: Partial<LanguageResourceModel>) => {
    const next = resources.map((item, itemIndex) => itemIndex === index ? { ...item, ...patch } : item);
    onChange(next);
  }, [onChange, resources]);

  const handleDelete = React.useCallback((index: number) => {
    const next = resources.filter((_, itemIndex) => itemIndex !== index);
    onChange(next);
  }, [onChange, resources]);

  const handleAdd = React.useCallback(() => {
    const next = [
      ...resources,
      {
        id: nextResourceId(resources),
        key: `admin.${activeModule || "general"}.new_key`,
        value: "",
      },
    ];
    onChange(next);
  }, [activeModule, onChange, resources]);

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
            <Alert severity="info">Chưa có resource nào cho module này.</Alert>
          ) : (
            moduleRows.map(({ resource, index }) => (
              <Box
                key={`${resource.id ?? "new"}-${index}`}
                sx={(theme) => ({
                  border: `1px solid ${theme.palette.divider}`,
                  borderRadius: 2,
                  p: 2,
                  backgroundColor: theme.palette.background.paper,
                })}
              >
                <Stack direction={{ xs: "column", md: "row" }} spacing={1.5} alignItems={{ md: "flex-start" }}>
                  <TextField
                    fullWidth
                    label="Resource key"
                    value={resource.key}
                    disabled={disabled}
                    onChange={(event) => handlePatch(index, { key: event.target.value })}
                    helperText="Bắt buộc theo format admin.{module}.*"
                  />
                  <IconButton
                    color="error"
                    disabled={disabled}
                    onClick={() => handleDelete(index)}
                    sx={{ alignSelf: { xs: "flex-end", md: "center" } }}
                  >
                    <DeleteOutlineIcon />
                  </IconButton>
                </Stack>
                <TextField
                  fullWidth
                  multiline
                  minRows={3}
                  label="Giá trị"
                  sx={{ mt: 1.5 }}
                  value={resource.value}
                  disabled={disabled}
                  onChange={(event) => handlePatch(index, { value: event.target.value })}
                />
              </Box>
            ))
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
            Resource keys phải theo format
          </Typography>
          <Chip size="small" label="admin.{module}.*" />
        </Stack>
        <Button variant="outlined" startIcon={<AddIcon />} disabled={disabled} onClick={handleAdd}>
          Thêm resource
        </Button>
      </Stack>
      <TabContainer tabs={tabs} defaultValue={activeModule} onChange={setActiveModule} />
    </Stack>
  );
}

export function LanguageDetailWidget() {
  const { languageId } = useParams();
  const { hasPermission } = usePermissionChecks();
  const formRef = React.useRef<AutoFormRef>(null);
  const [loading, setLoading] = React.useState(true);
  const [detail, setDetail] = React.useState<LanguageModel | null>(null);
  const [resources, setResources] = React.useState<LanguageResourceModel[]>([]);

  const numericLanguageId = Number(languageId ?? 0);

  const loadDetail = React.useCallback(async () => {
    if (!numericLanguageId) return;
    setLoading(true);
    try {
      const data = await getLanguageById(numericLanguageId);
      setDetail(data);
      setResources(data.resources ?? []);
    } catch (error: unknown) {
      toast.error(getErrorMessage(error, "Không thể tải chi tiết ngôn ngữ."));
    } finally {
      setLoading(false);
    }
  }, [numericLanguageId]);

  React.useEffect(() => {
    void loadDetail();
  }, [loadDetail]);

  React.useEffect(() => {
    formRef.current?.setValue("resources", resources);
  }, [resources]);

  const validateResources = React.useCallback(() => {
    const seen = new Set<string>();

    for (const resource of resources) {
      const key = String(resource.key ?? "").trim();
      if (!key) {
        toast.error("Resource key không được để trống.");
        return false;
      }
      if (!ADMIN_RESOURCE_KEY_REGEX.test(key)) {
        toast.error(`Resource key "${key}" không đúng format admin.{module}.*.`);
        return false;
      }
      if (seen.has(key)) {
        toast.error(`Resource key "${key}" đang bị trùng.`);
        return false;
      }
      seen.add(key);
    }

    return true;
  }, [resources]);

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
      toast.success("Xuất XML thành công.");
    } catch (error: unknown) {
      toast.error(getErrorMessage(error, "Xuất XML thất bại."));
    }
  }, [numericLanguageId]);

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
    return <Alert severity="error">Thiếu mã ngôn ngữ.</Alert>;
  }

  if (loading) {
    return <Alert severity="info">Đang tải chi tiết ngôn ngữ...</Alert>;
  }

  if (!detail) {
    return <Alert severity="error">Không tìm thấy ngôn ngữ.</Alert>;
  }

  return (
    <Stack spacing={2}>
      <SectionCard
        title="Thông tin ngôn ngữ"
        extra={
          <IfPermission permissions={["languages.update"]}>
            <SafeButton variant="contained" startIcon={<SaveOutlinedIcon />} onClick={() => void handleSave()}>
              Lưu
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
        title="Admin Resources"
        extra={
          <Stack direction={{ xs: "column", sm: "row" }} spacing={1}>
            <IfPermission permissions={["languages.export"]}>
              <Button variant="outlined" startIcon={<DownloadOutlinedIcon />} onClick={() => void handleExport()}>
                Xuất XML
              </Button>
            </IfPermission>
            <IfPermission permissions={["languages.import"]}>
              <Button variant="outlined" startIcon={<CloudUploadOutlinedIcon />} onClick={handleOpenImportDialog}>
                Nhập XML
              </Button>
            </IfPermission>
          </Stack>
        }
      >
        <Stack spacing={2}>
          <Alert severity="info">
            Dùng resource keys theo chuẩn <strong>admin.{"{module}"}.*</strong>. Ví dụ: <strong>admin.settings.title</strong>.
          </Alert>
          <Divider />
          <LanguageResourcesEditor
            resources={resources}
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
