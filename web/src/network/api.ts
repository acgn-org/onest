import axios, { AxiosError } from "axios";

export const baseUrl = "/api/";

const api = axios.create({
  baseURL: baseUrl,
});

function ErrHandler(err: AxiosError<any>): AxiosError {
  switch (true) {
    case err.name === "CanceledError":
      err.message = "";
      break;
    case err && err.response && err.response.data && err.response.data.msg:
      err.message = err.response.data.msg;
      break;
  }
  return err;
}

api.interceptors.response.use(undefined, (err: AxiosError) => {
  return Promise.reject(ErrHandler(err));
});

export default api;