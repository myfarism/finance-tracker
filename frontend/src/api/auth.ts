import api from "./axios";
import { AuthResponse, LoginInput, RegisterInput } from "../types/auth";

export const authAPI = {
  register: async (data: RegisterInput): Promise<AuthResponse> => {
    const res = await api.post("/auth/register", data);
    return res.data.data;
  },

  verifyOTP: async (email: string, code: string): Promise<AuthResponse> => {
    const res = await api.post("/auth/verify-otp", { email, code });
    return res.data.data;
  },

  resendOTP: async (email: string): Promise<void> => {
    await api.post("/auth/resend-otp", { email });
  },

  login: async (data: LoginInput): Promise<AuthResponse> => {
    const res = await api.post("/auth/login", data);
    return res.data.data;
  },
};
