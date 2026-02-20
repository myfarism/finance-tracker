import api from "./axios";
import {
  Transaction,
  CreateTransactionInput,
  UpdateTransactionInput,
  TransactionSummary,
  TransactionFilter,
  Category,
  Budget,
  UpsertBudgetInput,
} from "../types/transaction";

export const categoryAPI = {
  getAll: async (): Promise<Category[]> => {
    const res = await api.get("/categories");
    return res.data.data;
  },
};

export const transactionAPI = {
  create: async (data: CreateTransactionInput): Promise<Transaction> => {
    const res = await api.post("/transactions", data);
    return res.data.data;
  },

  getAll: async (filter?: TransactionFilter): Promise<Transaction[]> => {
    const res = await api.get("/transactions", { params: filter });
    return res.data.data;
  },

  getByID: async (id: string): Promise<Transaction> => {
    const res = await api.get(`/transactions/${id}`);
    return res.data.data;
  },

  update: async (id: string, data: UpdateTransactionInput): Promise<Transaction> => {
    const res = await api.put(`/transactions/${id}`, data);
    return res.data.data;
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/transactions/${id}`);
  },

  getSummary: async (month: number, year: number): Promise<TransactionSummary> => {
    const res = await api.get("/transactions/summary", { params: { month, year } });
    return res.data.data;
  },
};

export const budgetAPI = {
  getByMonth: async (month: number, year: number): Promise<Budget[]> => {
    const res = await api.get("/budgets", { params: { month, year } });
    return res.data.data ?? [];
  },

  upsert: async (data: UpsertBudgetInput): Promise<Budget> => {
    const res = await api.post("/budgets", data);
    return res.data.data;
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/budgets/${id}`);
  },
};
