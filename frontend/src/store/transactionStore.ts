import { create } from "zustand";
import { Transaction, TransactionFilter, TransactionSummary, Category, Budget } from "../types/transaction";
import { transactionAPI, categoryAPI, budgetAPI } from "../api/transaction";

interface TransactionState {
  transactions: Transaction[];
  summary: TransactionSummary | null;
  categories: Category[];
  isLoading: boolean;
  filter: TransactionFilter;
  currentMonth: number;
  currentYear: number;
  budgets: Budget[];

  fetchTransactions: () => Promise<void>;
  fetchSummary: (month: number, year: number) => Promise<void>;
  fetchCategories: () => Promise<void>;
  setFilter: (filter: TransactionFilter) => void;
  deleteTransaction: (id: string) => Promise<void>;
  fetchBudgets: (month: number, year: number) => Promise<void>;

  // Helper: refresh semua data sekaligus
  refreshAll: () => Promise<void>;
}

export const useTransactionStore = create<TransactionState>((set, get) => ({
  transactions: [],
  summary: null,
  categories: [],
  isLoading: false,
  filter: {},
  currentMonth: new Date().getMonth() + 1,
  currentYear: new Date().getFullYear(),
  budgets: [],

  fetchTransactions: async () => {
    set({ isLoading: true });
    try {
      const data = await transactionAPI.getAll(get().filter);
      set({ transactions: data });
    } finally {
      set({ isLoading: false });
    }
  },

  fetchSummary: async (month, year) => {
    // Simpan month & year agar bisa dipakai refreshAll
    set({ currentMonth: month, currentYear: year });
    const data = await transactionAPI.getSummary(month, year);
    set({ summary: data });
  },

  fetchCategories: async () => {
    const data = await categoryAPI.getAll();
    set({ categories: data });
  },

  setFilter: (filter) => {
    set({ filter });
    get().fetchTransactions();
  },

  // ✅ Fix utama: hapus + refresh transaksi & summary sekaligus
  deleteTransaction: async (id) => {
    await transactionAPI.delete(id);
    await get().refreshAll();
  },

  fetchBudgets: async (month, year) => {
    const data = await budgetAPI.getByMonth(month, year);
    set({ budgets: data });
  },

  // ✅ Helper: panggil ini setiap kali data berubah
  refreshAll: async () => {
    const { currentMonth, currentYear, filter } = get();

    set({ isLoading: true });
    try {
      const [transactions, summary, budgets] = await Promise.all([
        transactionAPI.getAll(filter),
        transactionAPI.getSummary(currentMonth, currentYear),
        budgetAPI.getByMonth(currentMonth, currentYear),
      ]);
      set({ transactions, summary, budgets });
    } finally {
      set({ isLoading: false });
    }
  },
}));
