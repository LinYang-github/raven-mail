import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
});

const USER_ID = 'user-123';

export const getInbox = (page = 1) => api.get(`/mails/inbox?page=${page}&user_id=${USER_ID}`);
export const getSent = (page = 1) => api.get(`/mails/sent?page=${page}&user_id=${USER_ID}`);
export const getMail = (id) => api.get(`/mails/${id}?user_id=${USER_ID}`);
export const sendMail = (formData) => api.post(`/mails/send?user_id=${USER_ID}`, formData, {
  headers: { 'Content-Type': 'multipart/form-data' }
});
export const getDownloadUrl = (path) => `http://localhost:8080/api/v1/mails/download?path=${encodeURIComponent(path)}&user_id=${USER_ID}`;
