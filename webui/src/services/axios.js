import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5,
});

instance.interceptors.request.use(
	(config) => {
	  const token = sessionStorage.getItem('IDToken');
	  if (token) {
		config.headers.Authorization = "Bearer " + token;
	  }
	  return config;
	},
	(error) => Promise.reject(error),
);
  

export default instance;
