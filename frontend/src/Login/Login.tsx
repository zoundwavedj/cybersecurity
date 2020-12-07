import React, { FC, FormEvent, Fragment, useState } from 'react';
import { Box, Button, Dialog, DialogActions, DialogContent, DialogTitle, Link, makeStyles, TextField, Typography } from '@material-ui/core';
import { useAuth } from '../hooks/auth';
import { Redirect, useHistory, useLocation } from 'react-router-dom';

const useStyles = makeStyles({
  link: {
    cursor: 'pointer',
  },
});

const Login: FC = () => {
  const auth = useAuth();
  const history = useHistory();
  const location = useLocation();
  const styles = useStyles();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [invalidLogin, setInvalidLogin] = useState(false);
  const [credsDialog, setCredsDialog] = useState(false);
  const [failureDialog, setFailureDialog] = useState(false);

  const login = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const { from } = (location.state as any) || { from: { pathname: '/' } };

    const success = await auth.login(username, password);

    if (success) {
      setInvalidLogin(false);
      history.replace(from);
    } else {
      setInvalidLogin(true);
    }
  };

  // THIS IS SUPER HACKY OOPS
  // Wait for hook's side effect to finish before attemping to render route :)
  if (auth.authenticated) {
    return <Redirect to={location.state as string} />
  }

  return (
    <Fragment>
      <Box
        display='flex'
        flexGrow={1}
        alignItems='center'
        justifyContent='center'>
        <Box
          border={1}
          borderColor='rgba(0,0,0,0.1)'
          borderRadius={8}
          display='flex'
          flexDirection='column'
          padding={2}>

          <form onSubmit={login}>
            <Box
              display='flex'
              flexDirection='column'
              marginBottom={2}>

              <Box marginBottom={1}>
                <TextField
                  label='Username'
                  onChange={(event) => setUsername(event.target.value)}
                  size='small'
                  variant='outlined' />
              </Box>

              <Box marginBottom={1}>
                <TextField
                  label='Password'
                  onChange={(event) => setPassword(event.target.value)}
                  size='small'
                  type='Password'
                  variant='outlined' />
              </Box>

              <Button fullWidth variant='contained' type='submit'>Login</Button>

              {invalidLogin &&
                <Box mt={1}>
                  <Typography color='error' align='center'>Invalid login!</Typography>
                </Box>
              }
            </Box>
          </form>

          <Typography align='center'>
            <Link
              className={styles.link}
              onClick={() => getSuperUserCreds(setUsername, setPassword, setCredsDialog, setFailureDialog)}>
              First time? Click here
            </Link>
          </Typography>
        </Box>
      </Box>

      {/* First-time creds dialog */}
      <Dialog
        fullWidth
        maxWidth='sm'
        open={credsDialog}>
        <DialogTitle>HERE YOU GO</DialogTitle>
        <DialogContent>
          <Typography>
            Username: {username}
          </Typography>
          <Typography>
            Password: {password}
          </Typography>
          <Box marginTop={2}>
            <Typography variant='subtitle2'>
              WARNING! This dialog will only appear once for the lifetime of this application.
              <br />
              Make sure to keep the credentials safe.
            </Typography>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button color='primary' onClick={() => setCredsDialog(false)}>DONE</Button>
        </DialogActions>
      </Dialog>

      {/* Failed retry dialog */}
      <Dialog
        fullWidth
        maxWidth='sm'
        open={failureDialog}>
        <DialogTitle>OOPS</DialogTitle>
        <DialogContent>
          Well, seems like you've already taken the creds before.
          <br />
          Told ya to keep em safe. Now what? :)
        </DialogContent>
        <DialogActions>
          <Button color='primary' onClick={() => setFailureDialog(false)}>I AM SAD</Button>
        </DialogActions>
      </Dialog>
    </Fragment >
  );
};

const getSuperUserCreds = (
  setUsername: React.Dispatch<React.SetStateAction<string>>,
  setPassword: React.Dispatch<React.SetStateAction<string>>,
  setCredsDialog: React.Dispatch<React.SetStateAction<boolean>>,
  setFailureDialog: React.Dispatch<React.SetStateAction<boolean>>
) => {
  fetch('/superuser', { method: 'GET' })
    .then(response => response.json())
    .then(data => {
      if (data.username && data.password) {
        setUsername(data.username);
        setPassword(data.password);
        setCredsDialog(true);
      } else {
        setFailureDialog(true);
      }
    })
    .catch(err => console.error(err));
}

export default Login;