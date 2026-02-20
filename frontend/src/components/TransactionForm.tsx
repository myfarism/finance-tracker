import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { transactionAPI } from "../api/transaction";
import { useTransactionStore } from "../store/transactionStore";
import { Transaction } from "../types/transaction";

type FormData = {
  category_id: string;
  type: "income" | "expense";
  amount: number;
  description?: string;
  date: string;
};

const schema = z.object({
  category_id: z.string().min(1, "Pilih kategori"),
  type: z.enum(["income", "expense"]),
  amount: z.number().min(1, "Nominal harus lebih dari 0"),
  description: z.string().optional(),
  date: z.string().min(1, "Pilih tanggal"),
});

interface Props {
  onClose: () => void;
  editData?: Transaction | null;
}

export default function TransactionForm({ onClose, editData }: Props) {
  const { categories, fetchTransactions, refreshAll } = useTransactionStore();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      type: "expense",
      date: new Date().toLocaleDateString("en-CA"),
    },
  });

  useEffect(() => {
    if (editData) {
      const rawDate = editData.date ?? "";
      const safeDate = rawDate.includes("T") ? rawDate.split("T")[0] : rawDate;
      reset({
        category_id: editData.category_id,
        type: editData.type,
        amount: editData.amount,
        description: editData.description ?? "",
        date: safeDate,
      });
    }
  }, [editData]);

  

  const onSubmit = async (data: FormData) => {
    try {
      if (editData) {
        await transactionAPI.update(editData.id, data);
      } else {
        await transactionAPI.create(data);
      }
      await refreshAll();
      onClose();
    } catch (err: any) {
      alert(err.response?.data?.message || "Terjadi kesalahan");
    }
  };

  return (
    // Overlay
    <div
      className="fixed inset-0 bg-black/40 z-50 flex items-end sm:items-center justify-center p-0 sm:p-4"
      onClick={(e) => e.target === e.currentTarget && onClose()}
    >
      {/* Panel — full width bottom sheet di mobile, popup di desktop */}
      <div className="bg-white w-full sm:max-w-md rounded-t-2xl sm:rounded-xl overflow-hidden">

        {/* Handle bar (mobile only) */}
        <div className="flex justify-center pt-3 pb-1 sm:hidden">
          <div className="w-10 h-1 bg-slate-200 rounded-full" />
        </div>

        {/* Header */}
        <div className="flex items-center justify-between px-5 py-4 border-b border-slate-100">
          <h2 className="text-sm font-semibold text-slate-800">
            {editData ? "Edit Transaksi" : "Transaksi Baru"}
          </h2>
          <button
            onClick={onClose}
            className="text-slate-400 hover:text-slate-600 transition text-lg leading-none"
          >
            ✕
          </button>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit(onSubmit)} className="px-5 py-4 space-y-4">

          {/* Tipe — Toggle style */}
          <div>
            <label className="text-xs font-medium text-slate-500 uppercase tracking-wider">
              Tipe
            </label>
            <div className="mt-1.5 grid grid-cols-2 gap-2">
              {(["expense", "income"] as const).map((t) => (
                <label
                  key={t}
                  className="flex items-center justify-center gap-2 border border-slate-200 rounded-lg py-2.5 cursor-pointer has-[:checked]:border-indigo-500 has-[:checked]:bg-indigo-50 transition"
                >
                  <input
                    {...register("type")}
                    type="radio"
                    value={t}
                    className="sr-only"
                  />
                  <span className="text-sm font-medium text-slate-700">
                    {t === "expense" ? "Pengeluaran" : "Pemasukan"}
                  </span>
                </label>
              ))}
            </div>
          </div>

          {/* Nominal */}
          <div>
            <label className="text-xs font-medium text-slate-500 uppercase tracking-wider">
              Nominal
            </label>
            <div className="mt-1.5 relative">
              <span className="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-slate-400">
                Rp
              </span>
              <input
                {...register("amount", { valueAsNumber: true })}
                type="number"
                placeholder="0"
                className="w-full border border-slate-200 bg-white rounded-lg pl-9 pr-3 py-2.5 text-sm text-slate-900 placeholder:text-slate-400 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition"
              />
            </div>
            {errors.amount && <p className="text-xs text-red-500 mt-1">{errors.amount.message}</p>}
          </div>

          {/* Kategori */}
          <div>
            <label className="text-xs font-medium text-slate-500 uppercase tracking-wider">
              Kategori
            </label>
            <select
              {...register("category_id")}
              className="mt-1.5 w-full border border-slate-200 bg-white rounded-lg px-3 py-2.5 text-sm text-slate-900 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition"
            >
              <option value="">Pilih kategori...</option>
              {categories.map((cat) => (
                <option key={cat.id} value={cat.id}>
                  {cat.icon} {cat.name}
                </option>
              ))}
            </select>
            {errors.category_id && <p className="text-xs text-red-500 mt-1">{errors.category_id.message}</p>}
          </div>

          {/* Deskripsi & Tanggal side by side di sm+ */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label className="text-xs font-medium text-slate-500 uppercase tracking-wider">
                Deskripsi
              </label>
              <input
                {...register("description")}
                type="text"
                placeholder="Opsional"
                className="mt-1.5 w-full border border-slate-200 bg-white rounded-lg px-3 py-2.5 text-sm text-slate-900 placeholder:text-slate-400 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition"
              />
            </div>

            <div>
              <label className="text-xs font-medium text-slate-500 uppercase tracking-wider">
                Tanggal
              </label>
              <input
                {...register("date")}
                type="date"
                className="mt-1.5 w-full border border-slate-200 bg-white rounded-lg px-3 py-2.5 text-sm text-slate-900 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition"
              />
              {errors.date && <p className="text-xs text-red-500 mt-1">{errors.date.message}</p>}
            </div>
          </div>

          {/* Actions */}
          <div className="flex gap-2 pt-1 pb-2">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 border border-slate-200 text-slate-600 text-sm font-medium py-2.5 rounded-lg hover:bg-slate-50 transition"
            >
              Batal
            </button>
            <button
              type="submit"
              disabled={isSubmitting}
              className="flex-1 bg-slate-900 hover:bg-slate-800 text-white text-sm font-medium py-2.5 rounded-lg transition disabled:opacity-40"
            >
              {isSubmitting ? "Menyimpan..." : "Simpan"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
