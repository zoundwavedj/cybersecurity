import Axios from 'axios';

const instance = Axios.create();

instance.interceptors.request.use(async req => {
  const accessToken = localStorage.getItem('accessToken');

  req.headers = {
    'Authorization': `Bearer ${accessToken}`,
    'Accept': 'application/json',
    'Content-Type': 'application/json',
  };

  return req;
}, err => Promise.reject(err));

instance.interceptors.response.use(resp => resp,
  async err => {
    const req = err.config;

    if ((err.response.status === 401 || err.response.status === 403) && !req.retry) {
      req.retry = true;

      const refreshToken = localStorage.getItem('refreshToken');
      const accessToken = await refreshAccessToken(refreshToken);

      if (accessToken) {
        req.headers = `Bearer ${accessToken}`;

        return instance(req);
      }
    }

    return Promise.reject(err);
  });

const refreshAccessToken = async (refreshToken: string | null) => {
  return await (fetch('/refresh', {
    body: JSON.stringify({ token: refreshToken }),
    headers: {
      'Content-Type': 'application/json',
    },
    method: 'POST',
  })
    .then(resp => resp.json())
    .then(data => {
      if (data.accessToken && data.refreshToken) {
        localStorage.setItem('accessToken', data.accessToken);
        localStorage.setItem('refreshToken', data.refreshToken);

        return data.accessToken;
      }
    })
    .catch(err => console.error(err)));
};

export default instance;