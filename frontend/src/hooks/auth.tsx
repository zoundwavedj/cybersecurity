import React, { createContext, useContext, useEffect, useState } from 'react';

const authContext = createContext<any>({});

export const ProvideAuth = ({ children }: any) => {
  const auth = useProvideAuth();
  return <authContext.Provider value={auth}>{children}</authContext.Provider>;
}

export const useAuth = (): {
  authenticated: boolean;
  login: (username: string, password: string) => Promise<boolean | void>;
  logout: () => Promise<boolean | void>;
} => {
  return useContext(authContext);
}

function useProvideAuth() {
  const [authenticated, setAuthenticated] = useState<boolean>(false);

  const login = async (username: string, password: string) => {
    return await (fetch("/login", {
      body: JSON.stringify({ username, password }),
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
          setAuthenticated(true);
          return true;
        } else {
          return false;
        }
      })
      .catch(err => {
        console.error(err);
        return false;
      }));
  };

  const logout = async () => {
    const accessToken = localStorage.getItem('accessToken');
    const refreshToken = localStorage.getItem('refreshToken');

    return await (fetch("/logout", {
      body: JSON.stringify({ accessToken, refreshToken }),
      headers: {
        'Content-Type': 'application/json',
      },
      method: 'POST',
    })
      .then(resp => resp.json())
      .then(data => {
        if (data.success) {
          setAuthenticated(false);
          localStorage.clear(); // Or clear only the tokens if there are other non-sensitive cached stuff
          return true;
        }

        return false;
      })
      .catch(err => {
        console.error(err);
        return false;
      }));
  };

  useEffect(() => {
    const accessToken = localStorage.getItem('accessToken');
    const refreshToken = localStorage.getItem('refreshToken');

    if (accessToken && refreshToken) {
      setAuthenticated(true);
    } else {
      setAuthenticated(false);
    }
  }, []);

  return {
    authenticated,
    login,
    logout,
  };
}