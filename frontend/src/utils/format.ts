export const formatCurrency = (amount: number): string => {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
  }).format(amount);
};

export const formatDate = (dateStr: string | null | undefined): string => {
  // Guard: kalau kosong atau null
  if (!dateStr) return "-";

  const date = new Date(dateStr);

  // Guard: kalau hasil parse tidak valid
  if (isNaN(date.getTime())) return "-";

  // Guard: Golang zero time "0001-01-01"
  if (date.getFullYear() < 2000) return "-";

  return new Intl.DateTimeFormat("id-ID", {
    day: "numeric",
    month: "long",
    year: "numeric",
  }).format(date);
};

export const getCurrentMonthYear = () => {
  const now = new Date();
  return { month: now.getMonth() + 1, year: now.getFullYear() };
};
