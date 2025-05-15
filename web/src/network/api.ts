import axios, { AxiosError } from "axios";

const BaseUrl = "/api/";

const api = axios.create({
  baseURL: BaseUrl,
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