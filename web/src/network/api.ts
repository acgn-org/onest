import axios, { AxiosError } from "axios";

export const baseUrl = "/api/";

const api = axios.create({
  baseURL: baseUrl,
});

api.interceptors.response.use(undefined, (err: AxiosError<any>) => {
  if (err.name === "CanceledError") return new Promise(() => {});
  if (err && err.response && err.response.data && err.response.data.msg)
    err.toString = () => err.response!.data.msg;
  return Promise.reject(err);
});

export default api;
