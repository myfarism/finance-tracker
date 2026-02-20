import { useEffect, useState } from "react";
import { useAuthStore } from "../store/authStore";
import { useTransactionStore } from "../store/transactionStore";
import { getCurrentMonthYear } from "../utils/format";
import SummaryCards from "../components/SummaryCards";
import TransactionChart from "../components/TransactionChart";
import TransactionList from "../components/TransactionList";
import TransactionForm from "../components/TransactionForm";
import { exportToCSV } from "../utils/export";
import BudgetSection from "../components/BudgetSection";


const MONTH_NAMES = [
  "", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
  "Juli", "Agustus", "September", "Oktober", "November", "Desember",
];

export default function DashboardPage() {
  const { user, logout } = useAuthStore();
  const {
    transactions, summary, isLoading,
    fetchTransactions, fetchSummary, fetchCategories, setFilter,
    budgets, fetchBudgets, categories,
  } = useTransactionStore();

  const [showForm, setShowForm] = useState(false);
  const [search, setSearch] = useState("");
  const [typeFilter, setTypeFilter] = useState("");
  const { month, year } = getCurrentMonthYear();

  useEffect(() => {
    fetchCategories();
    fetchTransactions();
    fetchSummary(month, year);
    fetchBudgets(month, year);
  }, []);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setFilter({ search, type: typeFilter });
  };

  const handleReset = () => {
    setSearch("");
    setTypeFilter("");
    setFilter({});
  };

  return (
    <div className="min-h-screen bg-slate-50">

      {/* Navbar */}
      <nav className="bg-white border-b border-slate-200 sticky top-0 z-20">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 h-14 flex items-center justify-between">
          <span className="text-base font-semibold text-slate-900 tracking-tight">
            finance<span className="text-indigo-600">.</span>
          </span>
          <div className="flex items-center gap-3">
            {/* Avatar + nama hanya di sm ke atas */}
            <span className="hidden sm:block text-sm text-slate-500">
              {user?.name}
            </span>
            {/* Avatar inisial di semua ukuran */}
            <div className="w-7 h-7 rounded-full bg-indigo-100 text-indigo-600 text-xs font-semibold flex items-center justify-center select-none">
              {user?.name?.charAt(0).toUpperCase()}
            </div>
            <button
              onClick={logout}
              className="text-sm text-slate-500 hover:text-slate-900 transition"
            >
              Keluar
            </button>
          </div>
        </div>
      </nav>

      {/* Page Header */}
      <div className="bg-white border-b border-slate-200">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 py-4 flex items-center justify-between">
            <div>
            <h2 className="text-base font-semibold text-slate-900">
                {MONTH_NAMES[month]} {year}
            </h2>
            <p className="text-xs text-slate-400 mt-0.5">
                Ringkasan keuangan bulan ini
            </p>
            </div>

            {/* Wrapper untuk kedua tombol */}
            <div className="flex items-center gap-3">
            <button
                onClick={() => setShowForm(true)}
                className="bg-slate-900 hover:bg-slate-800 active:bg-slate-950 text-white text-sm font-medium px-4 py-2 rounded-lg transition flex items-center gap-1.5"
            >
                <span className="text-base leading-none">+</span>
                <span className="hidden sm:inline">Tambah</span>
                <span className="sm:hidden">Baru</span>
            </button>

            <button
                onClick={() => exportToCSV(transactions)}
                className="border border-slate-200 bg-white text-slate-600 text-sm font-medium px-4 py-2 rounded-lg hover:bg-slate-50 transition"
            >
                Export CSV
            </button>
            </div>
        </div>
        </div>

      {/* Main */}
      <main className="max-w-5xl mx-auto px-4 sm:px-6 py-6 space-y-6">

        {/* Summary */}
        <SummaryCards summary={summary} />

        {/* Charts */}
        <TransactionChart transactions={transactions} />

        {/* Filter Bar */}
        <form
          onSubmit={handleSearch}
          className="flex flex-col sm:flex-row gap-2"
        >
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Cari transaksi..."
            className="flex-1 border border-slate-200 bg-white rounded-lg px-3 py-2 text-sm text-slate-900 placeholder:text-slate-400 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition"
          />
          <div className="flex gap-2">
            <select
              value={typeFilter}
              onChange={(e) => setTypeFilter(e.target.value)}
              className="flex-1 sm:flex-none border border-slate-200 bg-white rounded-lg px-3 py-2 text-sm text-slate-700 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition"
            >
              <option value="">Semua</option>
              <option value="income">Pemasukan</option>
              <option value="expense">Pengeluaran</option>
            </select>
            <button
              type="submit"
              className="bg-slate-900 text-white text-sm font-medium px-4 py-2 rounded-lg hover:bg-slate-800 transition"
            >
              Cari
            </button>
            <button
              type="button"
              onClick={handleReset}
              className="border border-slate-200 bg-white text-slate-600 text-sm font-medium px-4 py-2 rounded-lg hover:bg-slate-50 transition"
            >
              Reset
            </button>
          </div>
        </form>

        {/* List */}
        <TransactionList transactions={transactions} isLoading={isLoading} />

        <BudgetSection
          budgets={budgets}
          categories={categories}
          month={month}
          year={year}
        />

      </main>

      {showForm && <TransactionForm onClose={() => setShowForm(false)} />}
    </div>
  );
}
