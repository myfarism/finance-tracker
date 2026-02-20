import { TransactionSummary } from "../types/transaction";
import { formatCurrency } from "../utils/format";

interface Props {
  summary: TransactionSummary | null;
}

export default function SummaryCards({ summary }: Props) {
  const balance = summary?.balance ?? 0;

  const items = [
    { label: "Pemasukan", value: summary?.income ?? 0, accent: false },
    { label: "Pengeluaran", value: summary?.expense ?? 0, accent: false },
    { label: "Saldo", value: balance, accent: true },
  ];

  return (
    // mobile: stack vertical, sm+: horizontal sejajar
    <div className="grid grid-cols-1 sm:grid-cols-3 gap-px bg-slate-200 border border-slate-200 rounded-lg overflow-hidden">
      {items.map((item) => (
        <div key={item.label} className="bg-white px-5 py-4">
          <p className="text-xs font-medium text-slate-400 uppercase tracking-wider mb-2">
            {item.label}
          </p>
          <p className={`text-lg sm:text-xl font-semibold truncate ${
            item.accent
              ? balance >= 0 ? "text-indigo-600" : "text-red-500"
              : "text-slate-900"
          }`}>
            {formatCurrency(item.value)}
          </p>
        </div>
      ))}
    </div>
  );
}
