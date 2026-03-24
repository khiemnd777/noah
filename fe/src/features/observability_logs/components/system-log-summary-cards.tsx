import type { ReactNode } from "react";
import ErrorOutlineRoundedIcon from "@mui/icons-material/ErrorOutlineRounded";
import WarningAmberRoundedIcon from "@mui/icons-material/WarningAmberRounded";
import { Box, Paper, Stack, Typography } from "@mui/material";
import { ResponsiveGrid } from "@shared/components/ui/responsive-grid";
import type { SystemLogsSummaryModel } from "@features/observability_logs/model/system-log.model";

type SystemLogSummaryCardsProps = {
  summary: SystemLogsSummaryModel;
};

type SummaryCardProps = {
  title: string;
  value: number;
  caption: string;
  icon: ReactNode;
};

function SummaryCard({ title, value, caption, icon }: SummaryCardProps) {
  return (
    <Paper elevation={0} sx={{ p: 2.5, borderRadius: 3, border: "1px solid", borderColor: "divider" }}>
      <Stack direction="row" spacing={2} alignItems="center">
        <Box sx={{ display: "flex", alignItems: "center", justifyContent: "center" }}>
          {icon}
        </Box>
        <Box>
          <Typography variant="body2" color="text.secondary">{title}</Typography>
          <Typography variant="h4">{value}</Typography>
          <Typography variant="caption" color="text.secondary">{caption}</Typography>
        </Box>
      </Stack>
    </Paper>
  );
}

export function SystemLogSummaryCards({ summary }: SystemLogSummaryCardsProps) {
  return (
    <ResponsiveGrid xs={1} sm={2} md={2} lg={2} xl={2}>
      <SummaryCard
        title="Warn Logs"
        value={summary.warnCount}
        caption="Theo bộ lọc hiện tại"
        icon={<WarningAmberRoundedIcon />}
      />
      <SummaryCard
        title="Error Logs"
        value={summary.errorCount}
        caption="Theo bộ lọc hiện tại"
        icon={<ErrorOutlineRoundedIcon />}
      />
    </ResponsiveGrid>
  );
}
