import axios from 'axios';
import { userStore } from '../store/user';

const isDev = import.meta.env.DEV;
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || (isDev ? 'http://localhost:8080/api/v1' : '/api/v1');

const api = axios.create({
  baseURL: API_BASE_URL,
});

// 注入场次 ID 到所有请求头
api.interceptors.request.use(config => {
  config.headers['X-Session-ID'] = userStore.sessionId || 'default';
  return config;
});

// Helper to get current ID
const getUserID = () => userStore.id;

export const getInbox = (page = 1, query = '') => api.get(`/mails/inbox?page=${page}&user_id=${getUserID()}&q=${encodeURIComponent(query)}`);
export const getSent = (page = 1, query = '') => api.get(`/mails/sent?page=${page}&user_id=${getUserID()}&q=${encodeURIComponent(query)}`);
export const getMail = (id) => api.get(`/mails/${id}?user_id=${getUserID()}`);
export const deleteMail = (id) => api.delete(`/mails/${id}?user_id=${getUserID()}`);
export const sendMail = (formData) => api.post(`/mails/send?user_id=${getUserID()}`, formData, {
  headers: { 'Content-Type': 'multipart/form-data' }
});
export const triggerForceSave = (key) => api.post(`/onlyoffice/forcesave?key=${key}`);
export const deleteSession = (sessionId) => api.delete(`/sessions/${sessionId}`);
export const getDownloadUrl = (att) => `${API_BASE_URL}/mails/download?id=${att.id}&user_id=${getUserID()}`;
export const getPreviewUrl = (att) => `${API_BASE_URL}/mails/download?id=${att.id}&user_id=${getUserID()}&disposition=inline`;

// Chat APIs
export const sendChatMessage = (receiverId, content) => api.post(`/im/send?user_id=${getUserID()}`, { receiver_id: receiverId, content });
export const getChatHistory = (otherId) => api.get(`/im/history?user_id=${getUserID()}&other_id=${otherId}`);
export const markChatAsRead = (senderId) => api.post(`/im/read?user_id=${getUserID()}&sender_id=${senderId}`);
