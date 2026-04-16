export const prefixCurrency = "₫";

export function formatCurrency(value: number) {
  return new Intl.NumberFormat("vi-VN").format(Number.isFinite(value) ? value : 0);
}
