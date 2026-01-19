import axios from 'axios';
import { userStore } from '../store/user';

const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
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
export const getDownloadUrl = (att) => `http://localhost:8080/api/v1/mails/download?id=${att.id}&user_id=${getUserID()}`;
export const getPreviewUrl = (att) => `http://localhost:8080/api/v1/mails/download?id=${att.id}&user_id=${getUserID()}&disposition=inline`;
