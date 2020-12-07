import React, { createContext, useContext, useEffect, useState } from 'react';

export interface User {
  authenticated: boolean,
  accessToken: string,
  refreshToken: string,
};

const authContext = createContext<any>({});

export const ProvideAuth = ({ children }: any) => {
  const auth = useProvideAuth();
  return <authContext.Provider value={auth}>{children}</authContext.Provider>;
}

export const useAuth = (): {
  user: User | undefined;
  login: (username: string, password: string) => Promise<boolean | void>;
  logout: () => Promise<boolean | void>;
} => {
  return useContext(authContext);
}

function useProvideAuth() {
  const [user, setUser] = useState<User | undefined>(undefined);

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
          setUser({ authenticated: true, accessToken: data.accessToken, refreshToken: data.refreshToken });
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
    return await (fetch("/logout", {
      body: JSON.stringify({ accessToken: user?.accessToken, refreshToken: user?.refreshToken }),
      headers: {
        'Content-Type': 'application/json',
      },
      method: 'POST',
    })
      .then(resp => resp.json())
      .then(data => {
        if (data.success) {
          setUser(undefined);
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
      setUser({
        authenticated: true,
        accessToken,
        refreshToken,
      });
    } else {
      setUser(undefined);
    }
  }, []);

  return {
    user,
    login,
    logout,
  };
}