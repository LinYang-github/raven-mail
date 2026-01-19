import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
});

const USER_ID = 'user-123';

export const getInbox = (page = 1, query = '') => api.get(`/mails/inbox?page=${page}&user_id=${USER_ID}&q=${encodeURIComponent(query)}`);
export const getSent = (page = 1, query = '') => api.get(`/mails/sent?page=${page}&user_id=${USER_ID}&q=${encodeURIComponent(query)}`);
export const getMail = (id) => api.get(`/mails/${id}?user_id=${USER_ID}`);
export const deleteMail = (id) => api.delete(`/mails/${id}?user_id=${USER_ID}`);
export const sendMail = (formData) => api.post(`/mails/send?user_id=${USER_ID}`, formData, {
  headers: { 'Content-Type': 'multipart/form-data' }
});
export const getDownloadUrl = (att) => `http://localhost:8080/api/v1/mails/download?id=${att.id}&user_id=${USER_ID}`;
