import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Link, useNavigate } from "react-router-dom";
import { authAPI } from "../api/auth";
import { useAuthStore } from "../store/authStore";

// --- Schema ---
const registerSchema = z.object({
  name: z.string().min(2, "Nama minimal 2 karakter"),
  email: z.string().email("Email tidak valid"),
  password: z.string().min(8, "Password minimal 8 karakter"),
});

const otpSchema = z.object({
  code: z
    .string()
    .length(6, "Kode harus 6 digit")
    .regex(/^\d+$/, "Hanya angka"),
});

type RegisterForm = z.infer<typeof registerSchema>;
type OTPForm = z.infer<typeof otpSchema>;

export default function RegisterPage() {
  const navigate = useNavigate();
  const setAuth = useAuthStore((s) => s.setAuth);

  // Step: "register" | "otp"
  const [step, setStep] = useState<"register" | "otp">("register");
  const [email, setEmail] = useState("");
  const [resendCooldown, setResendCooldown] = useState(0);

  // Form register
  const registerForm = useForm<RegisterForm>({
    resolver: zodResolver(registerSchema),
  });

  // Form OTP
  const otpForm = useForm<OTPForm>({
    resolver: zodResolver(otpSchema),
  });

  // Submit register → kirim OTP
  const onRegister = async (data: RegisterForm) => {
    try {
      await authAPI.register(data);
      setEmail(data.email);
      setStep("otp");
      startCooldown();
    } catch (err: any) {
      registerForm.setError("root", {
        message: err.response?.data?.message || "Registrasi gagal",
      });
    }
  };

  // Submit OTP → login otomatis
  const onVerifyOTP = async (data: OTPForm) => {
    try {
      const result = await authAPI.verifyOTP(email, data.code);
      setAuth(result.user, result.token);
      navigate("/dashboard");
    } catch (err: any) {
      otpForm.setError("root", {
        message: err.response?.data?.message || "Kode OTP salah",
      });
    }
  };

  // Resend OTP dengan cooldown 60 detik
  const startCooldown = () => {
    setResendCooldown(60);
    const interval = setInterval(() => {
      setResendCooldown((prev) => {
        if (prev <= 1) { clearInterval(interval); return 0; }
        return prev - 1;
      });
    }, 1000);
  };

  const handleResend = async () => {
    if (resendCooldown > 0) return;
    try {
      await authAPI.resendOTP(email);
      startCooldown();
    } catch (err: any) {
      alert(err.response?.data?.message || "Gagal mengirim ulang OTP");
    }
  };

  return (
    <div className="min-h-screen w-full bg-slate-50 flex items-center justify-center px-4 py-12">
      <div className="w-full max-w-sm">

        {/* Logo */}
        <div className="mb-8">
          <span className="text-2xl font-semibold text-slate-900 tracking-tight">
            finance<span className="text-indigo-600">.</span>
          </span>
          <p className="text-slate-500 text-sm mt-1">
            {step === "register" ? "Buat akun baru" : "Verifikasi email kamu"}
          </p>
        </div>

        {/* ── STEP 1: Register Form ── */}
        {step === "register" && (
          <form onSubmit={registerForm.handleSubmit(onRegister)} className="space-y-4">
            <div className="space-y-1">
              <label className="text-sm font-medium text-slate-700">Nama</label>
              <input
                {...registerForm.register("name")}
                placeholder="John Doe"
                className="w-full border border-slate-200 bg-white rounded-lg px-3 py-2.5 text-sm placeholder:text-slate-400 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition"
              />
              {registerForm.formState.errors.name && (
                <p className="text-xs text-red-500">
                  {registerForm.formState.errors.name.message}
                </p>
              )}
            </div>

            <div className="space-y-1">
              <label className="text-sm font-medium text-slate-700">Email</label>
              <input
                {...registerForm.register("email")}
                type="email"
                placeholder="nama@email.com"
                className="w-full border border-slate-200 bg-white rounded-lg px-3 py-2.5 text-sm placeholder:text-slate-400 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition"
              />
              {registerForm.formState.errors.email && (
                <p className="text-xs text-red-500">
                  {registerForm.formState.errors.email.message}
                </p>
              )}
            </div>

            <div className="space-y-1">
              <label className="text-sm font-medium text-slate-700">Password</label>
              <input
                {...registerForm.register("password")}
                type="password"
                placeholder="Minimal 8 karakter"
                className="w-full border border-slate-200 bg-white rounded-lg px-3 py-2.5 text-sm placeholder:text-slate-400 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition"
              />
              {registerForm.formState.errors.password && (
                <p className="text-xs text-red-500">
                  {registerForm.formState.errors.password.message}
                </p>
              )}
            </div>

            {/* Root error */}
            {registerForm.formState.errors.root && (
              <div className="bg-red-50 border border-red-200 rounded-lg px-3 py-2.5">
                <p className="text-xs text-red-600">
                  {registerForm.formState.errors.root.message}
                </p>
              </div>
            )}

            <button
              type="submit"
              disabled={registerForm.formState.isSubmitting}
              className="w-full bg-slate-900 hover:bg-slate-800 text-white text-sm font-medium py-2.5 rounded-lg transition disabled:opacity-40"
            >
              {registerForm.formState.isSubmitting ? "Mengirim..." : "Daftar"}
            </button>
          </form>
        )}

        {/* ── STEP 2: OTP Form ── */}
        {step === "otp" && (
          <div>
            {/* Info email */}
            <div className="bg-indigo-50 border border-indigo-100 rounded-lg px-4 py-3 mb-5">
              <p className="text-sm text-indigo-700">
                Kode OTP dikirim ke{" "}
                <span className="font-semibold">{email}</span>
              </p>
              <p className="text-xs text-indigo-500 mt-0.5">
                Berlaku selama 5 menit
              </p>
            </div>

            <form onSubmit={otpForm.handleSubmit(onVerifyOTP)} className="space-y-4">
              <div className="space-y-1">
                <label className="text-sm font-medium text-slate-700">
                  Kode OTP
                </label>
                <input
                  {...otpForm.register("code")}
                  type="text"
                  inputMode="numeric"
                  maxLength={6}
                  placeholder="123456"
                  className="w-full border border-slate-200 bg-white rounded-lg px-3 py-2.5 text-sm text-center tracking-[0.5em] font-mono placeholder:tracking-normal placeholder:text-slate-400 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition"
                />
                {otpForm.formState.errors.code && (
                  <p className="text-xs text-red-500">
                    {otpForm.formState.errors.code.message}
                  </p>
                )}
              </div>

              {/* Root error */}
              {otpForm.formState.errors.root && (
                <div className="bg-red-50 border border-red-200 rounded-lg px-3 py-2.5">
                  <p className="text-xs text-red-600">
                    {otpForm.formState.errors.root.message}
                  </p>
                </div>
              )}

              <button
                type="submit"
                disabled={otpForm.formState.isSubmitting}
                className="w-full bg-slate-900 hover:bg-slate-800 text-white text-sm font-medium py-2.5 rounded-lg transition disabled:opacity-40"
              >
                {otpForm.formState.isSubmitting ? "Memverifikasi..." : "Verifikasi"}
              </button>
            </form>

            {/* Resend + back */}
            <div className="flex items-center justify-between mt-4">
              <button
                onClick={() => setStep("register")}
                className="text-sm text-slate-400 hover:text-slate-600 transition"
              >
                ← Ganti email
              </button>
              <button
                onClick={handleResend}
                disabled={resendCooldown > 0}
                className="text-sm text-slate-500 hover:text-slate-900 disabled:text-slate-300 transition"
              >
                {resendCooldown > 0
                  ? `Kirim ulang (${resendCooldown}s)`
                  : "Kirim ulang"}
              </button>
            </div>
          </div>
        )}

        {/* Link ke login */}
        {step === "register" && (
          <p className="text-sm text-slate-500 mt-6">
            Sudah punya akun?{" "}
            <Link
              to="/login"
              className="text-slate-900 font-medium underline underline-offset-2"
            >
              Masuk
            </Link>
          </p>
        )}

      </div>
    </div>
  );
}
