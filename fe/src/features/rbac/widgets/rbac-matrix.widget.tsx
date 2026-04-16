import { SectionCard } from "@root/shared/components/ui/section-card";
import { useI18n } from "@root/core/i18n/use-i18n";
import { registerSlot } from "@root/core/module/registry";
import { RBACMatrix } from "@features/rbac/components/rbac-matrix";

function RBACMatrixWidget() {
  const { t } = useI18n();

  return (
    <>
      <SectionCard title={t("admin.rbac.matrix.section_title")}>
        <RBACMatrix />
      </SectionCard>
    </>
  );
}

registerSlot({
  id: "rbac-matrix",
  name: "rbac:left",
  priority: 1,
  render: () => <RBACMatrixWidget />,
});
