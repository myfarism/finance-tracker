import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Budget, Category } from "../types/transaction";
import { budgetAPI } from "../api/transaction";
import { useTransactionStore } from "../store/transactionStore";
import { formatCurrency } from "../utils/format";

type FormData = {
  category_id: string;
  amount: number;
};

const schema = z.object({
  category_id: z.string().min(1, "Pilih kategori"),
  amount: z.number().min(1, "Masukkan nominal"),
});

interface Props {
  budgets: Budget[];
  categories: Category[];
  month: number;
  year: number;
}

export default function BudgetSection({ budgets, categories, month, year }: Props) {
  const { fetchBudgets } = useTransactionStore();
  const [showForm, setShowForm] = useState(false);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<FormData>({ resolver: zodResolver(schema) });

  const onSubmit = async (data: FormData) => {
    try {
      await budgetAPI.upsert({ ...data, month, year });
      await fetchBudgets(month, year);
      reset();
      setShowForm(false);
    } catch (err: any) {
      alert(err.response?.data?.message || "Gagal menyimpan budget");
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Hapus budget ini?")) return;
    await budgetAPI.delete(id);
    await fetchBudgets(month, year);
  };

  return (
    <div className="bg-white border border-slate-200 rounded-lg overflow-hidden">
      {/* Header */}
      <div className="px-5 py-4 border-b border-slate-100 flex items-center justify-between">
        <div>
          <h3 className="text-sm font-semibold text-slate-700">Budget Bulanan</h3>
          <p className="text-xs text-slate-400 mt-0.5">Batas pengeluaran per kategori</p>
        </div>
        <button
          onClick={() => setShowForm((v) => !v)}
          className="text-xs font-medium text-indigo-600 hover:text-indigo-700 border border-indigo-200 hover:border-indigo-300 px-3 py-1.5 rounded-lg transition"
        >
          {showForm ? "Batal" : "+ Atur Budget"}
        </button>
      </div>

      {/* Form Tambah Budget */}
      {showForm && (
        <div className="px-5 py-4 bg-slate-50 border-b border-slate-100">
          <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col sm:flex-row gap-3">
            <select
              {...register("category_id")}
              className="flex-1 border border-slate-200 bg-white rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition"
            >
              <option value="">Pilih kategori...</option>
              {categories.map((cat) => (
                <option key={cat.id} value={cat.id}>
                  {cat.icon} {cat.name}
                </option>
              ))}
            </select>

            <div className="relative flex-1">
              <span className="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-slate-400">
                Rp
              </span>
              <input
                {...register("amount", { valueAsNumber: true })}
                type="number"
                placeholder="500000"
                className="w-full border border-slate-200 bg-white rounded-lg pl-9 pr-3 py-2 text-sm focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition"
              />
            </div>

            <button
              type="submit"
              disabled={isSubmitting}
              className="bg-slate-900 hover:bg-slate-800 text-white text-sm font-medium px-4 py-2 rounded-lg transition disabled:opacity-40"
            >
              {isSubmitting ? "Menyimpan..." : "Simpan"}
            </button>
          </form>

          {(errors.category_id || errors.amount) && (
            <p className="text-xs text-red-500 mt-2">
              {errors.category_id?.message || errors.amount?.message}
            </p>
          )}
        </div>
      )}

      {/* Budget List */}
      {budgets.length === 0 ? (
        <div className="px-5 py-8 text-center">
          <p className="text-sm text-slate-500">Belum ada budget</p>
          <p className="text-xs text-slate-400 mt-1">
            Atur batas pengeluaran per kategori
          </p>
        </div>
      ) : (
        <ul className="divide-y divide-slate-50">
          {budgets.map((b) => {
            const pct = Math.min(
              b.amount > 0 ? (b.spent / b.amount) * 100 : 0,
              100
            );

            return (
              <li key={b.id} className="px-5 py-4 group hover:bg-slate-50 transition">
                {/* Row atas: kategori + jumlah + action */}
                <div className="flex items-center justify-between mb-2">
                  <div className="flex items-center gap-2">
                    <span className="text-base">{b.category?.icon ?? "📦"}</span>
                    <span className="text-sm font-medium text-slate-700">
                      {b.category?.name}
                    </span>
                    {b.is_over && (
                      <span className="text-xs font-medium text-red-500 bg-red-50 border border-red-100 px-2 py-0.5 rounded-full">
                        Over budget
                      </span>
                    )}
                  </div>
                  <div className="flex items-center gap-3">
                    <span className="text-xs text-slate-500 tabular-nums">
                      {formatCurrency(b.spent)}
                      <span className="text-slate-300 mx-1">/</span>
                      {formatCurrency(b.amount)}
                    </span>
                    <button
                      onClick={() => handleDelete(b.id)}
                      className="text-slate-300 hover:text-red-500 transition opacity-0 group-hover:opacity-100 text-xs"
                    >
                      ✕
                    </button>
                  </div>
                </div>

                {/* Progress bar */}
                <div className="h-1.5 bg-slate-100 rounded-full overflow-hidden">
                  <div
                    className={`h-full rounded-full transition-all duration-500 ${
                      b.is_over ? "bg-red-400" : pct > 80 ? "bg-amber-400" : "bg-indigo-500"
                    }`}
                    style={{ width: `${pct}%` }}
                  />
                </div>

                {/* Row bawah: sisa budget */}
                <div className="flex justify-between mt-1.5">
                  <span className="text-xs text-slate-400">
                    {pct.toFixed(0)}% terpakai
                  </span>
                  <span className={`text-xs font-medium ${
                    b.is_over ? "text-red-500" : "text-slate-500"
                  }`}>
                    {b.is_over
                      ? `Melebihi ${formatCurrency(Math.abs(b.remaining))}`
                      : `Sisa ${formatCurrency(b.remaining)}`}
                  </span>
                </div>
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
}
