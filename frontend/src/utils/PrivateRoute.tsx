import React, { FC } from 'react';
import { Redirect, Route, RouteProps } from 'react-router-dom';
import { useAuth } from '../hooks/auth';

const PrivateRoute: FC<RouteProps> = ({ children, ...props }) => {
  const auth = useAuth();

  return (
    <Route
      {...props}
      render={
        ({ location }) => auth.user?.accessToken && auth.user.refreshToken ?
          children : <Redirect to={{ pathname: "/login", state: { from: location } }} />
      }
    />
  );
};

export default PrivateRoute;