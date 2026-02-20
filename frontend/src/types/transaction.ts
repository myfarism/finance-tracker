export type TransactionType = "income" | "expense";

export interface Category {
  id: string;
  name: string;
  icon: string;
}

export type TransactionType = "income" | "expense";

export interface Category {
  id: string;
  name: string;
  icon: string;
}

export interface Transaction {
  id: string;
  user_id: string;
  category_id: string;
  category: Category;
  type: TransactionType;
  amount: number;
  description: string;
  date: string | null;
  created_at: string | null;
  updated_at: string | null;
}


export interface CreateTransactionInput {
  category_id: string;
  type: TransactionType;
  amount: number;
  description?: string;
  date: string;
}

export interface UpdateTransactionInput {
  category_id?: string;
  type?: TransactionType;
  amount?: number;
  description?: string;
  date?: string;
}

export interface TransactionSummary {
  income: number;
  expense: number;
  balance: number;
  month: number;
  year: number;
}

export interface TransactionFilter {
  type?: string;
  category_id?: string;
  search?: string;
  start_date?: string;
  end_date?: string;
}

export interface Budget {
  id: string;
  category_id: string;
  category: Category;
  amount: number;
  month: number;
  year: number;
  spent: number;
  remaining: number;
  is_over: boolean;
}

export interface UpsertBudgetInput {
  category_id: string;
  amount: number;
  month: number;
  year: number;
}
