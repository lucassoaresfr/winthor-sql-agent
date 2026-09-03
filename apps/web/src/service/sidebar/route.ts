import { connectDB } from "@/service/connect";
import {
  CreateChatRequest,
  CreateChatResponse,
  SaveMessageRequest,
  SaveMessageResponse,
  GetChatByIdResponse,
  ListUserChatsResponse,
} from "@/types/sidebar";

const api = connectDB();

// Função auxiliar para garantir o cabeçalho em todas as chamadas
const getAuthHeader = () => {
  const token = process.env.NEXT_PUBLIC_API_AUTH_TOKEN;
  return {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  };
};

export const chatService = {
  // 1. POST /api/v1/chats
  createChat: async (data: CreateChatRequest): Promise<CreateChatResponse> => {
    try {
      console.log("[chatService.createChat] Enviando:", data);
      const response = await api.post<CreateChatResponse>(
        "/chats",
        data,
        getAuthHeader(),
      );
      return response.data;
    } catch (error: any) {
      console.error(
        "[chatService.createChat] Erro:",
        error.response?.data || error.message,
      );
      throw error;
    }
  },

  // 2. POST /api/v1/chats/:id/messages
  addMessage: async (
    chatId: string,
    data: SaveMessageRequest,
  ): Promise<SaveMessageResponse> => {
    try {
      const response = await api.post<SaveMessageResponse>(
        `/chats/${chatId}/messages`,
        data,
        getAuthHeader(),
      );
      return response.data;
    } catch (error: any) {
      console.error(
        "[chatService.addMessage] Erro:",
        error.response?.data || error.message,
      );
      throw error;
    }
  },

  // 3. GET /api/v1/chats/:id
  getChatById: async (chatId: string): Promise<GetChatByIdResponse> => {
    try {
      const response = await api.get<GetChatByIdResponse>(
        `/chats/${chatId}`,
        getAuthHeader(),
      );
      return response.data;
    } catch (error: any) {
      console.error(
        "[chatService.getChatById] Erro:",
        error.response?.data || error.message,
      );
      throw error;
    }
  },

  // 4. GET /api/v1/chats/user/:user_id
  listUserChats: async (userId: string): Promise<ListUserChatsResponse> => {
    try {
      const response = await api.get<ListUserChatsResponse>(
        `/chats/user/${userId}`,
        getAuthHeader(),
      );
      return response.data;
    } catch (error: any) {
      console.error(
        "[chatService.listUserChats] Erro:",
        error.response?.data || error.message,
      );
      throw error;
    }
  },
};
