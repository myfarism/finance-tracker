import {
  BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer,
  Cell, PieChart, Pie,
} from "recharts";
import { Transaction } from "../types/transaction";
import { formatCurrency } from "../utils/format";

interface Props {
  transactions: Transaction[];
}

const CustomTooltip = ({ active, payload, label }: any) => {
  if (!active || !payload?.length) return null;
  return (
    <div className="bg-white border border-slate-200 rounded-lg px-3 py-2 text-xs shadow-sm">
      <p className="text-slate-400 mb-1">{label}</p>
      {payload.map((p: any) => (
        <p key={p.name} className="font-medium text-slate-800">
          {p.name}: {formatCurrency(p.value)}
        </p>
      ))}
    </div>
  );
};

const PIE_COLORS = ["#6366f1", "#818cf8", "#a5b4fc", "#c7d2fe", "#e0e7ff"];

export default function TransactionChart({ transactions }: Props) {
  const barData = Array.from({ length: 4 }, (_, i) => ({
    name: `W${i + 1}`,
    Masuk: transactions
      .filter((t) => {
        const day = new Date(t.date ?? "").getDate();
        return t.type === "income" && day >= i * 7 + 1 && day <= (i + 1) * 7;
      })
      .reduce((s, t) => s + t.amount, 0),
    Keluar: transactions
      .filter((t) => {
        const day = new Date(t.date ?? "").getDate();
        return t.type === "expense" && day >= i * 7 + 1 && day <= (i + 1) * 7;
      })
      .reduce((s, t) => s + t.amount, 0),
  }));

  const expenseMap = transactions
    .filter((t) => t.type === "expense")
    .reduce((acc, t) => {
      const name = t.category?.name ?? "Lainnya";
      acc[name] = (acc[name] ?? 0) + t.amount;
      return acc;
    }, {} as Record<string, number>);

  const pieData = Object.entries(expenseMap)
    .map(([name, value]) => ({ name, value }))
    .sort((a, b) => b.value - a.value)
    .slice(0, 5);

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      {/* Bar Chart */}
      <div className="bg-white border border-slate-200 rounded-lg p-5">
        <p className="text-sm font-medium text-slate-700 mb-4">Arus Kas per Minggu</p>
        <ResponsiveContainer width="100%" height={180}>
          <BarChart data={barData} barGap={3}>
            <XAxis
              dataKey="name"
              tick={{ fontSize: 11, fill: "#94a3b8" }}
              axisLine={false}
              tickLine={false}
            />
            <YAxis
              tickFormatter={(v) => `${(v / 1000).toFixed(0)}k`}
              tick={{ fontSize: 11, fill: "#94a3b8" }}
              axisLine={false}
              tickLine={false}
              width={36}
            />
            <Tooltip content={<CustomTooltip />} cursor={{ fill: "#f8fafc" }} />
            <Bar dataKey="Masuk" fill="#6366f1" radius={[3, 3, 0, 0]} maxBarSize={20} />
            <Bar dataKey="Keluar" fill="#e2e8f0" radius={[3, 3, 0, 0]} maxBarSize={20} />
          </BarChart>
        </ResponsiveContainer>

        {/* Legend */}
        <div className="flex gap-4 mt-3">
          {[{ label: "Masuk", color: "#6366f1" }, { label: "Keluar", color: "#e2e8f0" }].map((l) => (
            <div key={l.label} className="flex items-center gap-1.5">
              <div className="w-2.5 h-2.5 rounded-sm" style={{ backgroundColor: l.color }} />
              <span className="text-xs text-slate-500">{l.label}</span>
            </div>
          ))}
        </div>
      </div>

      {/* Pie Chart */}
      <div className="bg-white border border-slate-200 rounded-lg p-5">
        <p className="text-sm font-medium text-slate-700 mb-4">Pengeluaran per Kategori</p>
        {pieData.length === 0 ? (
          <div className="h-[180px] flex items-center justify-center">
            <p className="text-slate-400 text-sm">Belum ada data</p>
          </div>
        ) : (
          <div className="flex items-center gap-4">
            <div className="flex-shrink-0">
              <ResponsiveContainer width={140} height={140}>
                <PieChart>
                  <Pie
                    data={pieData}
                    cx="50%"
                    cy="50%"
                    innerRadius={40}
                    outerRadius={65}
                    dataKey="value"
                    strokeWidth={0}
                  >
                    {pieData.map((_, i) => (
                      <Cell key={`cell-${i}`} fill={PIE_COLORS[i % PIE_COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip content={<CustomTooltip />} />
                </PieChart>
              </ResponsiveContainer>
            </div>

            {/* Legend */}
            <div className="flex-1 space-y-2 min-w-0">
              {pieData.map((item, i) => (
                <div key={item.name} className="flex items-center gap-2 min-w-0">
                  <div
                    className="w-2 h-2 rounded-full flex-shrink-0"
                    style={{ backgroundColor: PIE_COLORS[i % PIE_COLORS.length] }}
                  />
                  <span className="text-xs text-slate-500 truncate">{item.name}</span>
                  <span className="text-xs text-slate-400 ml-auto flex-shrink-0">
                    {formatCurrency(item.value)}
                  </span>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
