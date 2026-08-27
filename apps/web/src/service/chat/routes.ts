import { ChatApiRequest, ChatApiResponse } from "@/types/chat";
import { connectDB } from "../connect";

export const ChatApi = async (
  data: ChatApiRequest,
): Promise<ChatApiResponse> => {
  try {
    const api = connectDB();
    const token = process.env.NEXT_PUBLIC_ORCHESTRATOR_TOKEN || "";

    const response = await api.post<ChatApiResponse>("/chat", data, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    return response.data;
  } catch (err: any) {
    console.error("[ChatApi Error]:", err);

    const errorMessage =
      err.response?.data?.error ||
      err.message ||
      "Erro de conexão com o servidor.";

    return {
      error: errorMessage,
    };
  }
};
