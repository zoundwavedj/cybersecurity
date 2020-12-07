import React, { FC, Fragment } from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import { CssBaseline } from '@material-ui/core';
import { Home } from './Home';
import { Login } from './Login';
import PrivateRoute from './utils/PrivateRoute';
import { ProvideAuth } from './hooks/auth';

const App: FC = () => {
  return (
    <Fragment>
      <CssBaseline />

      <ProvideAuth>
        <Router>
          <Switch>
            <Route exact path='/login'>
              <Login />
            </Route>
            <PrivateRoute exact path='/'>
              <Home />
            </PrivateRoute>
          </Switch>
        </Router>
      </ProvideAuth>
    </Fragment>
  );
}

export default App;
