import React, { FC, Fragment, useCallback, useEffect, useState } from 'react';
import { AppBar, Box, Toolbar, Typography, Button, TextField, Dialog, DialogActions, DialogContent, DialogTitle, TableContainer, Table, TableHead, TableBody, TableCell, TableRow } from '@material-ui/core';
import { useHistory } from 'react-router-dom';
import { useAuth } from '../hooks/auth';
import axiosInstance from '../utils/axios';

const Home: FC = () => {
  const auth = useAuth();
  const history = useHistory();
  const [name, setName] = useState("");
  const [dob, setDob] = useState("");
  const [email, setEmail] = useState("");
  const [ssn, setSsn] = useState("");
  const [users, setUsers] = useState<Array<string>>([]);
  const [newUserDialog, setNewUserDialog] = useState(false);
  const [viewUserDialog, setViewUserDialog] = useState(false);

  const logout = async () => {
    if (await auth.logout()) {
      history.replace('/login');
    }
  };

  const createUser = () => {
    axiosInstance.post('/user', { name, dob, email, ssn })
      .then(data => {
        if (data.status === 200) {
          listUsers();
          setNewUserDialog(false);
        }
      })
      .catch(err => {
        if (err.response.status === 401 || err.response.status === 403) {
          auth.logout();
        }
      });
  };

  const viewUser = (id: string) => {
    axiosInstance.get(`/user?id=${id}`)
      .then(data => {
        if (data.status === 200) {
          setName(data.data.name);
          setDob(data.data.dob);
          setEmail(data.data.email);
          setSsn(data.data.ssn);
          setViewUserDialog(true);
        }
      })
      .catch(err => {
        if (err.response.status === 401 || err.response.status === 403) {
          auth.logout();
        }
      });
  };

  const listUsers = useCallback(() => {
    axiosInstance.get('/users')
      .then(data => {
        if (data.status === 200) {
          setUsers(data.data.users);
        }
      })
      .catch(err => {
        if (err.response.status === 401 || err.response.status === 403) {
          auth.logout();
        }
      });
  }, [auth]);

  useEffect(() => {
    listUsers();
  }, [listUsers]);

  return (
    <Fragment>
      <AppBar position='static'>
        <Toolbar>
          <Box display='flex' flexGrow={1} justifyContent='space-between'>
            <Box display='flex' alignItems='center'>
              <Typography variant='h6'>
                Welcome!
              </Typography>
            </Box>

            <Button color="inherit" onClick={logout}>Log out</Button>
          </Box>
        </Toolbar>
      </AppBar>

      <Box display='flex' flex={1} flexDirection='column'>
        <Box display='flex' alignItems='center' marginX={2} marginTop={2}>
          {/* <TextField
            label='Search user'
            size='small'
            variant='outlined'
          /> */}

          <Button variant='outlined' color='primary' onClick={() => setNewUserDialog(true)}>New user</Button>
        </Box>

        <Box border={1} borderColor='rgba(0,0,0,0.1)' borderRadius={4} margin={2}>
          <TableContainer>
            <Table size='small'>
              <TableHead>
                <TableRow>
                  <TableCell>ID</TableCell>
                  <TableCell></TableCell>
                </TableRow>
              </TableHead>

              <TableBody>
                {users.map(user => (
                  <TableRow key={user}>
                    <TableCell>{user}</TableCell>
                    <TableCell align='right'>
                      <Button variant='contained' onClick={() => viewUser(user)}>View</Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </Box>
      </Box>

      {/* New user dialog */}
      <Dialog
        fullWidth
        maxWidth='sm'
        open={newUserDialog}
      >
        <DialogTitle>New User</DialogTitle>

        <DialogContent>
          <Box marginBottom={1}>
            <TextField
              fullWidth
              label='Name'
              onChange={event => setName(event.target.value)}
              size='small'
              variant='outlined'
            />
          </Box>

          <Box marginBottom={1}>
            <TextField
              fullWidth
              label='Date of birth'
              onChange={event => setDob(event.target.value)}
              size='small'
              variant='outlined'
            />
          </Box>

          <Box marginBottom={1}>
            <TextField
              fullWidth
              label='Email'
              onChange={event => setEmail(event.target.value)}
              size='small'
              variant='outlined'
            />
          </Box>

          <TextField
            fullWidth
            label='SSN'
            onChange={event => setSsn(event.target.value)}
            size='small'
            variant='outlined'
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setNewUserDialog(false)}>CANCEL</Button>
          <Button color='primary' onClick={createUser}>DONE</Button>
        </DialogActions>
      </Dialog>

      {/* View user dialog */}
      <Dialog
        fullWidth
        maxWidth='sm'
        open={viewUserDialog}
      >
        <DialogTitle>New User</DialogTitle>

        <DialogContent>
          <Box marginBottom={1}>
            <TextField
              disabled
              fullWidth
              label='Name'
              size='small'
              value={name}
              variant='outlined'
            />
          </Box>

          <Box marginBottom={1}>
            <TextField
              disabled
              fullWidth
              label='Date of birth'
              size='small'
              value={dob}
              variant='outlined'
            />
          </Box>

          <Box marginBottom={1}>
            <TextField
              disabled
              fullWidth
              label='Email'
              size='small'
              value={email}
              variant='outlined'
            />
          </Box>

          <TextField
            disabled
            fullWidth
            label='SSN'
            size='small'
            value={ssn}
            variant='outlined'
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => {
            setName("");
            setDob("");
            setEmail("");
            setSsn("");
            setViewUserDialog(false);
          }}>CLOSE</Button>
        </DialogActions>
      </Dialog>
    </Fragment>
  );
};

export default Home;