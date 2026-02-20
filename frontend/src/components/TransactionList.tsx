import { useState } from "react";
import { Transaction } from "../types/transaction";
import { formatCurrency, formatDate } from "../utils/format";
import { useTransactionStore } from "../store/transactionStore";
import TransactionForm from "./TransactionForm";

interface Props {
  transactions: Transaction[];
  isLoading: boolean;
}

export default function TransactionList({ transactions, isLoading }: Props) {
  const { deleteTransaction } = useTransactionStore();
  const [editData, setEditData] = useState<Transaction | null>(null);

  const handleDelete = async (id: string) => {
    if (!confirm("Hapus transaksi ini?")) return;
    await deleteTransaction(id);
  };

  if (isLoading) {
    return (
      <div className="border border-slate-200 rounded-lg p-10 text-center bg-white">
        <p className="text-slate-400 text-sm">Memuat...</p>
      </div>
    );
  }

  if (transactions.length === 0) {
    return (
      <div className="border border-dashed border-slate-300 rounded-lg p-10 text-center bg-white">
        <p className="text-slate-500 text-sm font-medium">Tidak ada transaksi</p>
        <p className="text-slate-400 text-xs mt-1">Tambahkan transaksi pertama kamu</p>
      </div>
    );
  }

  return (
    <>
      <div className="bg-white border border-slate-200 rounded-lg overflow-hidden">
        {/* Header */}
        <div className="px-4 sm:px-5 py-3 border-b border-slate-100 bg-slate-50 flex items-center justify-between">
          <span className="text-xs font-medium text-slate-400 uppercase tracking-wider">
            Transaksi
          </span>
          <span className="text-xs text-slate-400">
            {transactions.length} item
          </span>
        </div>

        <ul className="divide-y divide-slate-50">
          {transactions.map((tx) => (
            <li
              key={tx.id}
              className="px-4 sm:px-5 py-3.5 hover:bg-slate-50 transition group"
            >
              <div className="flex items-center gap-3">
                {/* Icon */}
                <span className="text-base w-6 text-center flex-shrink-0">
                  {tx.category?.icon ?? "•"}
                </span>

                {/* Info */}
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium text-slate-800 truncate">
                    {tx.description || tx.category?.name}
                  </p>
                  <p className="text-xs text-slate-400 mt-0.5">
                    {tx.category?.name ?? "—"} · {formatDate(tx.date)}
                  </p>
                </div>

                {/* Amount + Actions */}
                <div className="flex items-center gap-2 flex-shrink-0">
                  <span className={`text-sm font-semibold tabular-nums ${
                    tx.type === "income" ? "text-indigo-600" : "text-slate-700"
                  }`}>
                    {tx.type === "income" ? "+" : "−"}
                    {formatCurrency(tx.amount)}
                  </span>

                  {/* Desktop: hover · Mobile: selalu tampil */}
                  <div className="flex gap-1 sm:opacity-0 sm:group-hover:opacity-100 transition-opacity">
                    <button
                      onClick={() => setEditData(tx)}
                      className="text-xs text-slate-400 hover:text-indigo-600 p-1 rounded hover:bg-indigo-50 transition"
                      title="Edit"
                    >
                      ✎
                    </button>
                    <button
                      onClick={() => handleDelete(tx.id)}
                      className="text-xs text-slate-400 hover:text-red-500 p-1 rounded hover:bg-red-50 transition"
                      title="Hapus"
                    >
                      ✕
                    </button>
                  </div>
                </div>
              </div>
            </li>
          ))}
        </ul>
      </div>

      {editData && (
        <TransactionForm onClose={() => setEditData(null)} editData={editData} />
      )}
    </>
  );
}
