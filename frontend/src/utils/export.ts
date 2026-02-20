import { Transaction } from "../types/transaction";
import { formatCurrency, formatDate } from "./format";

export const exportToCSV = (transactions: Transaction[]) => {
  const headers = ["Tanggal", "Deskripsi", "Kategori", "Tipe", "Jumlah"];
  const rows = transactions.map((tx) => [
    formatDate(tx.date),
    tx.description || "-",
    tx.category?.name || "-",
    tx.type === "income" ? "Pemasukan" : "Pengeluaran",
    tx.amount.toString(),
  ]);

  const csv = [headers, ...rows]
    .map((row) => row.join(","))
    .join("\n");

  const blob = new Blob([csv], { type: "text/csv;charset=utf-8;" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = `transaksi-${new Date().toLocaleDateString("en-CA")}.csv`;
  a.click();
  URL.revokeObjectURL(url);
};
