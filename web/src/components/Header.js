import React from 'react';
import { AppBar, Toolbar, Typography, Button, Box, IconButton } from '@mui/material';
import { Brightness4, Brightness7, AccountCircle } from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';

const Header = () => {
  const navigate = useNavigate();
  const [darkMode, setDarkMode] = React.useState(true);

  const toggleDarkMode = () => {
    setDarkMode(!darkMode);
    // In a real app, this would toggle the theme
  };

  return (
    <AppBar 
      position="fixed" 
      sx={{ 
        zIndex: (theme) => theme.zIndex.drawer + 1,
        background: 'linear-gradient(90deg, #1e3c72 0%, #2a5298 100%)',
      }}
    >
      <Toolbar>
        <Typography
          variant="h6"
          noWrap
          component="div"
          sx={{ 
            display: 'flex', 
            alignItems: 'center', 
            fontWeight: 700,
            cursor: 'pointer'
          }}
          onClick={() => navigate('/')}
        >
          <Box 
            component="img"
            src="/assets/logo.svg" 
            alt="TeaThis Logo" 
            sx={{ 
              height: 32, 
              mr: 1,
              display: { xs: 'none', sm: 'block' }
            }}
          />
          TeaThis
        </Typography>
        <Box sx={{ flexGrow: 1 }} />
        <IconButton color="inherit" onClick={toggleDarkMode}>
          {darkMode ? <Brightness7 /> : <Brightness4 />}
        </IconButton>
        <Button 
          color="inherit" 
          startIcon={<AccountCircle />}
          sx={{ ml: 2 }}
        >
          Account
        </Button>
      </Toolbar>
    </AppBar>
  );
};

export default Header;
